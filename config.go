package main

import (
	"crypto/ed25519"
	"encoding/json"
	"net"
	"os"
)

type Config struct {
	Address    net.Addr           `json:"address"`
	Id         string             `json:"id"`
	Nickname   string             `json:"nickname"`
	PublicKey  ed25519.PublicKey  `json:"pubkey"`
	PrivateKey ed25519.PrivateKey `json:"privkey"`
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
	jsonBytes, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(location, jsonBytes, os.ModePerm)
}

func (node Node) LoadConfig(location string) error {
	config, err := node.config.Load(location)
	if err != nil {
		return err
	}
	node.address = config.Address
	node.nickname = config.Nickname
	node.crypto.ed.pub = config.PublicKey
	node.crypto.ed.priv = config.PrivateKey
	node.id = config.Id
	return nil
}
