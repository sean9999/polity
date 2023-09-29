package main

import (
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type Config struct {
	Address    NodeAddress        `json:"address"`
	Id         uuid.UUID          `json:"id"`
	Nickname   string             `json:"nickname"`
	PublicKey  ed25519.PublicKey  `json:"pubkey"`
	PrivateKey ed25519.PrivateKey `json:"privkey"`
	Friends    []NodeAddress      `json:"friends"`
}

func (c Config) String() string {
	return fmt.Sprintf("nick:\t%s\npub:\t%s\naddr:\t%s", c.Nickname, c.PublicKey, c.Address)
}

// Config.Load loads a config from a file
func (c Config) Load(location string) (Config, error) {
	var config Config
	raw, err := os.ReadFile(location)
	if err != nil {
		return config, err
	}
	err = json.Unmarshal(raw, &config)
	return config, err
}

// Config.Save saves a config to file
func (c Config) Save(location string) error {
	jsonBytes, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	return os.WriteFile(location, jsonBytes, os.ModePerm)
}

func (node *Node) LoadConfig(location string) error {
	config, err := node.config.Load(location)
	if err != nil {
		return err
	}
	node.address = config.Address
	node.nickname = config.Nickname
	node.crypto.ed.pub = config.PublicKey
	node.crypto.ed.priv = config.PrivateKey
	node.id = config.Id
	node.friends = append(node.friends, config.Friends...)
	return nil
}

func (node Node) GetConfig() Config {
	var c Config

	c.Address = node.address
	c.Id = node.Id()
	c.Nickname = node.Nickname()
	c.PublicKey = node.crypto.ed.pub
	c.PrivateKey = node.crypto.ed.priv
	c.Friends = node.Friends()

	return c
}
