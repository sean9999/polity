package main

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type consoleLogger interface {
	Log(...string)
	Info(...string)
	Warn(...string)
	Error(error)
	Debug(...string)
}

type consoleMachine struct {
	stdOut io.Writer
	stdErr io.Writer
}

func (c consoleMachine) Log(strs ...string) {
	for _, str := range strs {
		log.Log().Msg(str)
	}
}

func (c consoleMachine) Debug(strs ...string) {
	for _, str := range strs {
		log.Debug().Msg(str)
	}
}

func (c consoleMachine) Info(strs ...string) {
	for _, str := range strs {
		log.Info().Msg(str)
	}
}

func (c consoleMachine) Warn(strs ...string) {
	for _, str := range strs {
		log.Warn().Msg(str)
	}
}

func (c consoleMachine) Error(e error) {
	log.Err(e)
}

func NewConsole(w1 io.Writer, w2 io.Writer) consoleMachine {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"

	c := consoleMachine{
		stdOut: w1,
		stdErr: w2,
	}
	return c
}

func LogEnvelope(whichway string, e Envelope) {

	log.Info().
		Str("From", e.From.Username()).
		Str("To", e.From.Username()).
		Str("Subject", e.Message.Subject).
		Str("Body", e.Message.Body).
		Msg(whichway)

}

var console consoleLogger = NewConsole(os.Stdout, os.Stderr)
