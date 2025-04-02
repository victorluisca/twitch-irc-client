package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	client, err := NewIRCClientBuilder().
		WithNickname("your_twitch_nickname").
		WithPassword("oauth:your_oauth_token").
		WithCapabilities("twitch.tv/tags").
		Connect("irc.chat.twitch.tv:6667")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	if err := client.Join("#example_channel"); err != nil {
		fmt.Println(err)
		return
	}

	go client.ListenForMessages()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message != "" {
			if err := client.SendMessage(message); err != nil {
				fmt.Println(err)
			}
		}
	}
}
