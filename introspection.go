package main

import "github.com/google/uuid"

func (node Node) Info() Config {
	var c Config
	c.Id = node.id
	c.Nickname = node.nickname
	c.PublicKey = node.crypto.ed.pub
	return c
}

func (n Node) Id() uuid.UUID {
	return n.id
}

func (n Node) Nickname() string {
	return n.nickname
}
