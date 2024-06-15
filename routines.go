package polity

import (
	"crypto/rand"
	"fmt"

	"github.com/google/uuid"
)

func (c *Citizen) Assert() Message {

	bodyAsText := "Hi.\nI'm %s.\nMy pubkey is %s.\nMy stable address is %s.\nA nonce I've never used before is %s.\n"
	body := fmt.Sprintf(bodyAsText, c.Nickname(), c.PublicKeyAsHex(), c.Network.Address(), uuid.Must(uuid.NewRandom()))
	msg := c.Compose(SubjAssertion, []byte(body))
	msg.Plain.Headers["pubkey"] = string(c.PublicKeyAsHex())
	msg.Plain.Sign(rand.Reader, c.PrivateSigningKey())
	return msg

}
