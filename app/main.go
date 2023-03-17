package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

var gClient *openai.Client

// 先打开shadowsocks代理: clash for windows
func main() {

	fmt.Println("please ask:")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		question := scanner.Text()
		fmt.Printf("your question is :%v \n", question)
		fmt.Println(" please wait for resp ")
		// fmt.Println("echo reply is :", question)
		received, err := send(question)
		if err != nil {
			fmt.Println("send failed:", err)
		} else {
			fmt.Println("Success, receive resp:", received)
		}

	}
}

func send(content string) (received string, err error) {
	resp, err := gClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return
	}
	received = resp.Choices[0].Message.Content
	return
}

func init() {
	gClient = newClient()
}

func newClient() *openai.Client {
	sk := "<secret key>"
	config := openai.DefaultConfig(sk)
	proxyUrl, err := url.Parse("http://localhost:7890") //local proxy
	if err != nil {
		panic(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyUrl),
	}
	config.HTTPClient = &http.Client{
		Transport: transport,
	}
	client := openai.NewClientWithConfig(config)
	// client := openai.NewClient(sk) // no proxy
	return client
}
