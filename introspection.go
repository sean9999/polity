package main

import "github.com/google/uuid"

func (n Node) Id() uuid.UUID {
	return n.id
}

func (n Node) Nickname() string {
	return n.nickname
}
