package polity

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func (c *Citizen) Assert() Message {

	//	cryptographic proof that I am me
	bodyAsText := "Hi.\nI'm %s.\nMy pubkey is %s.\nMy stable address is %s.\nA nonce I've never used before is %s.\n"
	body := fmt.Sprintf(bodyAsText, c.Nickname(), c.PublicKeyAsHex(), c.Connection.LocalAddr(), uuid.Must(uuid.NewRandom()))
	msg := c.Compose(SubjAssertion, []byte(body))
	msg.ThreadId = msg.Id
	msg.Plain.Headers["pubkey"] = string(c.PublicKeyAsHex())
	msg.Plain.Sign(rand.Reader, c.PrivateSigningKey())
	return msg

}

func (c *Citizen) Howdee() (Message, error) {

	//	these are my friends. Who are your friends?
	j, err := json.MarshalIndent(c.Peers(), "", "\t")
	if err != nil {
		return NoMessage, err
	}
	return c.Compose(SubjWhoDoYouKnow, j), nil

}
