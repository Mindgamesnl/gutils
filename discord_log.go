package gutils

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var apiTokens = make(map[string]DiscordLogger)

type DiscordLogger struct {
	ChannelName string
	Token string
	UserName string
	Avatar string
}

//Print Write to a discord channel
func (dl DiscordLogger) Print(message string) {
	go PrintToChannel(dl, message)
}

//AddDiscordChannel Create an async discord logger from a webhook
func AddDiscordChannel(name string, token string, userName string, avatar string) DiscordLogger {
	l := DiscordLogger{name, token, userName, avatar}
	apiTokens[name] = l
	return l
}

func PrintToChannel(l DiscordLogger, message string)  {
	var dict = map[string]string{}
	dict["username"] = l.UserName
	dict["avatar_url"] = l.Avatar
	dict["content"] = message

	data, err := json.Marshal(dict)
	if err != nil {
		return
	}

	_, err = http.Post(l.Token, "application/json", bytes.NewBuffer(data))
	if err != nil{
		return
	}
}