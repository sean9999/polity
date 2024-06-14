package polity3

import (
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

func (c *Citizen) Assert() Message {

	bodyAsText := "Hi.\nI'm\t%s.\nMy pubkey is\t%s.\nMy stable address is\t%s.\nA nonce I've never used before is\t%s.\n"
	body := fmt.Sprintf(bodyAsText, c.Nickname(), c.PublicKeyAsHex(), c.network.Address(), uuid.Must(uuid.NewRandom()))
	msg := c.Compose("I assert myself", []byte(body))
	msg.Plain.Headers["pubkey"] = string(c.PublicKeyAsHex())
	msg.Plain.Sign(rand.Reader, c.PrivateSigningKey())
	return msg

}
