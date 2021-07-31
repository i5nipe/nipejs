package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	Bot     = os.Getenv("NIPEKEY")
	Chat_id = os.Getenv("NIPEJSCHAT")
)

func Deuboa(mensagem string) {

	params := url.Values{}
	params.Add("chat_id", Chat_id)
	params.Add("text", mensagem)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.telegram.org/bot%v/sendMessage", Bot), body)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}
