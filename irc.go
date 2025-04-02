package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type IRCClient struct {
	password     string
	nickname     string
	capabilities string
	conn         net.Conn
	channel      string
}

type IRCClientBuilder struct {
	password     string
	nickname     string
	capabilities string
}

func NewIRCClientBuilder() *IRCClientBuilder {
	return &IRCClientBuilder{}
}

func (b *IRCClientBuilder) WithPassword(password string) *IRCClientBuilder {
	b.password = password
	return b
}

func (b *IRCClientBuilder) WithNickname(nickname string) *IRCClientBuilder {
	b.nickname = nickname
	return b
}

func (b *IRCClientBuilder) WithCapabilities(capabilities string) *IRCClientBuilder {
	b.capabilities = capabilities
	return b
}

func (b *IRCClientBuilder) Connect(addr string) (*IRCClient, error) {
	client := &IRCClient{
		password:     b.password,
		nickname:     b.nickname,
		capabilities: b.capabilities,
	}

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Twitch IRC: %w", err)
	}
	client.conn = conn
	fmt.Println("Connected to Twitch IRC.")

	if _, err := fmt.Fprintf(client.conn, "CAP REQ :%s\r\n", client.capabilities); err != nil {
		client.Close()
		return nil, fmt.Errorf("error sending CAP REQ command: %w", err)
	}

	if _, err := fmt.Fprintf(client.conn, "PASS %s\r\n", client.password); err != nil {
		client.Close()
		return nil, fmt.Errorf("error sending PASS command: %w", err)
	}

	if _, err := fmt.Fprintf(client.conn, "NICK %s\r\n", client.nickname); err != nil {
		client.Close()
		return nil, fmt.Errorf("error sending NICK command: %w", err)
	}

	return client, nil
}

func (c *IRCClient) Join(channel string) error {
	c.channel = channel
	if _, err := fmt.Fprintf(c.conn, "JOIN %s\r\n", c.channel); err != nil {
		return fmt.Errorf("error sending JOIN command: %w", err)
	}
	fmt.Printf("Joined channel: %s\n", c.channel)
	return nil
}

func (c *IRCClient) ListenForMessages() error {
	reader := bufio.NewReader(c.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading from server: %w", err)
		}
		fmt.Print(line)

		if strings.HasPrefix(line, "PING") {
			if _, err := fmt.Fprintf(c.conn, "PONG :%s\r\n", strings.TrimPrefix(line, "PING :")); err != nil {
				return fmt.Errorf("error sending PONG command: %w", err)
			}
		}
	}
}

func (c *IRCClient) SendMessage(message string) error {
	if _, err := fmt.Fprintf(c.conn, "PRIVMSG %s :%s\r\n", c.channel, message); err != nil {
		return fmt.Errorf("error sending PRIVMSG command: %w", err)
	}
	return nil
}

func (c *IRCClient) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
	fmt.Println("Disconnected from Twitch IRC.")
}
