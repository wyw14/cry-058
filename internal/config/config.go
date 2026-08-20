package config

import "os"

type Config struct {
	HTTPAddr       string
	DefaultYear    int
	AttachmentRoot string
}

func Load() Config {
	year := 2026
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	root := os.Getenv("ATTACHMENT_ROOT")
	if root == "" {
		root = "./var/attachments"
	}
	return Config{HTTPAddr: addr, DefaultYear: year, AttachmentRoot: root}
}
