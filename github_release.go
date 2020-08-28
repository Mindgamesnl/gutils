package gutils

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"errors"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//DownloadLatestRelease Download a new binary from a Gorelease on github.
//Example usage DownloadLatestRelease(token, "Mindgamesnl", "OpenAudioMc-GoRelay", "OpenAudioMc-GoRelay", "relay")
//will download the latest binary for this platform and save it as relay
func DownloadLatestRelease(githubToken string, githubUsername string, githubRepo string, targetBinaryName string, saveBinaryAs string) *error {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)


	release, _, err := client.Repositories.GetLatestRelease(ctx, githubUsername, githubRepo)

	if err != nil {
		return &err
	}

	var osName = ""
	goOs := runtime.GOOS
	switch goOs {
	case "windows":
		osName = "Windows"
	case "darwin":
		osName = "Darwin"
	case "linux":
		osName = "Linux"
	default:
		osName = "Fuck"
	}

	var arcName = runtime.GOARCH

	if arcName == "amd64" {
		arcName = "x86_64"
	} else {
		arcName = "i386"
	}

	target := osName + "_" + arcName + ".tar.gz"

	var found = false

	for i := range release.Assets {
		asset := release.Assets[i]
		if strings.HasSuffix(asset.GetName(), target) {
			found = true
			_ = downloadFile(target, "https://" +githubToken + ":@api.github.com/repos/" + githubUsername + "/" + githubRepo + "/releases/assets/"+strconv.FormatInt(asset.GetID(), 10))
			r, err := os.Open(target)
			if err != nil {
				return &err
			}
			time.Sleep(1 * time.Second)
			ExtractTarGz(r)
			_ = os.Rename(targetBinaryName, saveBinaryAs)
		}
	}

	if !found {
		errie := errors.New("Could not find release file")
		return &errie
	}
	return nil
}

// write as it downloads and not load the whole file into memory.
func downloadFile(filepath string, url string) error {

	// Get the data
	client := &http.Client{

	}

	resp, err := client.Get(url)
	// ...

	req, err := http.NewRequest("GET", url, nil)
	// ...
	req.Header.Add("Accept", `application/octet-stream`)
	resp, err = client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func ExtractTarGz(gzipStream io.Reader) *error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return &err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for true {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("ExtractTarGz: Next() failed: %s", err.Error())
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				log.Fatalf("ExtractTarGz: Mkdir() failed: %s", err.Error())
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				log.Fatalf("ExtractTarGz: Create() failed: %s", err.Error())
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				log.Fatalf("ExtractTarGz: Copy() failed: %s", err.Error())
			}
			_ = outFile.Close()

		default:
			log.Fatalf(
				"ExtractTarGz: uknown type: %s in %s",
				header.Typeflag,
				header.Name)
		}

	}

	return nil
}
