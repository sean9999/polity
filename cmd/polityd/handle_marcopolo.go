package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/sean9999/polity"
)

func turn(txt string) (string, int, error) {
	allLines := strings.Split(txt, "\n")
	lines := []string{}
	for _, str := range allLines {
		if len(str) > 0 {
			lines = append(lines, str)
		}
	}
	if len(lines) != 2 {
		return "", 0, errors.New("text is not two lines")
	}
	who := lines[0]
	num, err := strconv.Atoi(lines[1])
	return who, num, err
}

func nextTurn(who string, m int) string {
	nextWho := "marco"
	if who == "marco" {
		nextWho = "polo"
	}
	return fmt.Sprintf("%s\n%d\n", nextWho, m+1)
}

func handleMarco(me *polity.Citizen, msg polity.Message) error {
	turn, m, err := turn(msg.Body())
	if err != nil {
		return err
	}
	if m < 1024 {
		response := me.Compose(polity.SubjStartMarcoPolo, []byte(nextTurn(turn, m)))
		me.Send(response, msg.Sender().Address())
	}
	return nil
}
