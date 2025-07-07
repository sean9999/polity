package main

import (
	"bytes"
	"encoding/pem"
	"fmt"
	"io"
	"os"

	"github.com/sean9999/hermeti"
)

// a pemBag is a loose collection of PEM encoded data
type pemBag map[string][]*pem.Block

// a subCommand is a function that runs against an [*app]
type subCommand func(hermeti.Env, *app)

// a functionMap organizes subCommands
type functionMap map[string]subCommand

func stdinHasData(e *hermeti.Env) bool {
	r := e.InStream
	f, ok := r.(*os.File)
	if ok {
		fi, err := f.Stat()
		if err != nil {
			return false
		}
		// Return true if stdin is not a character device (i.e., data is being piped in)
		return (fi.Mode() & os.ModeCharDevice) == 0
	} else {
		//	if we can rewind, do that
		r, ok := r.(io.ReadSeeker)
		if ok {
			i, err := r.Read(make([]byte, 1))
			if err != nil {
				return false
			}
			if i == 0 {
				return false
			}
			r.Seek(0, 0)
			return true
		} else {
			//	we can't rewind, so let's create a buffer
			buf := new(bytes.Buffer)
			io.Copy(buf, r)
			e.InStream = buf
			return buf.Len() > 0
		}
	}
}

// Bagify reads in an [io.Reader], collecting PEMs into a pemBag.
// It exits when the byte-stream can no longer be decoded into a PEM.
func (a *app) bagify(r io.Reader, ptr *pemBag) error {
	bag := *ptr
	pemBytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("could not bagify input. %w", err)
	}
	block, rest := pem.Decode(pemBytes)
	for {
		if block == nil {
			break
		}
		bag[block.Type] = append(bag[block.Type], block)
		block, rest = pem.Decode(rest)
	}
	return nil
}
