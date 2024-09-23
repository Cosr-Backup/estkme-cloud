package config

import (
	"errors"
	"net"
	"os"
	"strings"
)

type Config struct {
	ListenAddress string
	Prompt        string
	Verbose       bool
}

var C = &Config{}

var (
	ErrPromptTooLong = errors.New("prompt message is too long (max: 100 characters)")
	ErrInvalidPrompt = errors.New("prompt message contains non-printable ASCII characters")
)

func (c *Config) IsValid() error {
	if _, err := net.ResolveTCPAddr("tcp", c.ListenAddress); err != nil {
		return err
	}
	if len(c.Prompt) > 100 {
		return ErrPromptTooLong
	}
	// Prompt message is only allowed contain printable ASCII characters
	for _, r := range c.Prompt {
		if r < 32 || r > 126 {
			return ErrInvalidPrompt
		}
	}
	return nil
}

func (c *Config) GetPrompt() []byte {
	if c.Prompt != "" {
		return []byte("!! Prompt !! \n" + strings.Replace(c.Prompt, "_br_", "\n", -1))
	}
	return []byte{}
}

func (c *Config) LoadEnv() {
	if os.Getenv("ESTKME_CLOUD_LISTEN_ADDRESS") != "" {
		c.ListenAddress = os.Getenv("ESTKME_CLOUD_LISTEN_ADDRESS")
	}
	if os.Getenv("ESTKME_CLOUD_ADVERTISING") != "" {
		c.Prompt = os.Getenv("ESTKME_CLOUD_ADVERTISING")
	}
	if os.Getenv("ESTKME_CLOUD_VERBOSE") != "" {
		c.Verbose = true
	}
}
