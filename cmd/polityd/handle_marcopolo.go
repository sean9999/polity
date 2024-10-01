package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sean9999/go-flargs"
	"github.com/sean9999/polity"
)

type MarcoPoloBody struct {
	GameId uuid.UUID
	Who    polity.Subject
	Num    int
}

func (mpb MarcoPoloBody) Serialize() string {
	return fmt.Sprintf("%s\n%s\n%d\n", &mpb.GameId, mpb.Who, mpb.Num)
}

func (mpb *MarcoPoloBody) Deserialize(txt string) error {
	allLines := strings.Split(txt, "\n")
	lines := []string{}
	for _, str := range allLines {
		if len(str) > 0 {
			lines = append(lines, str)
		}
	}
	if len(lines) != 3 {
		return errors.New("text is not three lines")
	}
	num, err := strconv.Atoi(lines[2])
	if err != nil {
		return err
	}
	mpb.GameId = uuid.MustParse(lines[0])
	mpb.Who = polity.Subject(lines[1])
	mpb.Num = num
	return nil
}

func (previous MarcoPoloBody) Next() MarcoPoloBody {
	alt := polity.SubjMarco
	if previous.Who == alt {
		alt = polity.SubjPolo
	}
	nxt := MarcoPoloBody{
		GameId: previous.GameId,
		Who:    alt,
		Num:    previous.Num + 1,
	}
	return nxt
}

func NewMarcoPoloGame() MarcoPoloBody {
	id, _ := uuid.NewRandom()
	return MarcoPoloBody{
		GameId: id,
		Who:    polity.SubjStartMarcoPolo,
		Num:    0,
	}
}

func handleMarco(env *flargs.Environment, me *polity.Citizen, msg polity.Message) error {
	upperBound := 4096

	b := new(MarcoPoloBody)
	err := b.Deserialize(msg.Body())
	if err != nil {
		return err
	}

	c := b.Next()

	var startTime time.Time
	var stopime time.Time
	if c.Num == 1 {
		startTime = time.Now()
	}
	if c.Num < upperBound {
		response := me.Compose(c.Who, []byte(c.Serialize()))
		me.Send(response, msg.Sender())
	} else {
		//	mister even prints out
		stopime = time.Now()
		totalDuration := time.Duration(stopime.Nanosecond() - startTime.Nanosecond())
		averageDuration := time.Duration(totalDuration / time.Duration(upperBound))

		//	print out
		fmt.Println("game id:\t", c.GameId)
		fmt.Println("num hops:\t", upperBound)
		fmt.Println("total time:\t", totalDuration.String())
		fmt.Println("average hop:\t", averageDuration.String())

		if c.Num == upperBound {
			//	mister odd prints out
			response := me.Compose(polity.SubjStartMarcoPolo, []byte(c.Serialize()))
			me.Send(response, msg.Sender())
		}
	}
	return nil
}
