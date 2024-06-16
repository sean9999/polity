package polity

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sean9999/go-oracle"
)

func (c *Citizen) StartMarcoPolo() Message {

	gameId := uuid.Must(uuid.NewRandom()).String()
	body := fmt.Sprintf("do you want to play marco polo with gameId %s?", gameId)
	msg := c.Compose(SubjStartMarcoPolo, []byte(body))
	msg.Plain.Headers["gameId"] = gameId
	err := c.Sign(msg.Plain)
	if err != nil {
		panic(err)
	}
	return msg

}

func (c *Citizen) Marco(msg Message) Message {

	myMsg := new(oracle.PlainText)
	myMsg.Clone(msg.Plain)
	response := Message{
		Plain: myMsg,
	}
	return response

}
