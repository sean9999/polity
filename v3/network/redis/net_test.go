package redis

import (
	"encoding/json"
	"testing"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type dog struct {
	Barks bool
	Name  string
}

func (d dog) MarshalBinary() ([]byte, error) {
	return json.Marshal(d)
}

func (d *dog) UnmarshalBinary(b []byte) error {
	return json.Unmarshal(b, d)
}

func TestNetwork(t *testing.T) {
	n := new(Network)
	err := n.Up(t.Context())
	assert.NoError(t, err)
	assert.NotNil(t, n.rdb)
}

func TestNetwork_Up(t *testing.T) {
	n := new(Network)
	err := n.Up(t.Context())
	require.NoError(t, err)

	kp := delphi.KeyPair{}

	node := n.Spawn()
	err = node.Connect(t.Context(), kp)
	require.NoError(t, err)

	t.Run("pubsub", func(t *testing.T) {

		key := node.addr.String()

		go func() {
			err := node.rdb.Publish(t.Context(), key, "hi").Err()
			assert.NoError(t, err)
		}()

		msg := <-node.inbox
		assert.Equal(t, "hi", msg.Payload)

	})

	t.Run("pubsub auto marshal", func(t *testing.T) {

		go func() {
			err := node.rdb.Publish(t.Context(), node.addr.String(), &dog{Name: "fido", Barks: true}).Err()
			assert.NoError(t, err)
		}()

		data := <-node.inbox

		dog := &dog{}
		err = dog.UnmarshalBinary([]byte(data.Payload))
		require.NoError(t, err)
		assert.Equal(t, dog.Name, "fido")

	})

}
