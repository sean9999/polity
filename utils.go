package main

import "crypto/rand"

var randy = rand.Reader

func barfOn(err error) {
	if err != nil {
		panic(err)
	}
}
