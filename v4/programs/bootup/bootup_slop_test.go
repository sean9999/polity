package bootup

import (
	"io"
	"net/url"
	"testing"

	"github.com/sean9999/polity/v4"
	"github.com/stretchr/testify/assert"
)

type mockNode struct{ polity.Node }

func (m mockNode) URL() *url.URL {
	u, _ := url.Parse("test://localhost")
	return u
}

func TestBootup_Slop(t *testing.T) {
	c := polity.NewCitizen(io.Discard, mockNode{})
	p := new(prog)

	t.Run("Init", func(t *testing.T) {
		err := p.Init(c, nil, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, c, p.c)

		err = p.Init(nil, nil, nil, nil)
		assert.Error(t, err)
	})

	t.Run("Subjects", func(t *testing.T) {
		subs := p.Subjects()
		assert.NotEmpty(t, subs)
	})

	t.Run("Run and Shutdown", func(t *testing.T) {
		p.Init(c, nil, nil, nil)
		p.Run(t.Context())
		p.Shutdown()
	})

	t.Run("bootupGreeting", func(t *testing.T) {
		msg := bootupGreeting("alice", "memnet://peer@memory")
		assert.Contains(t, msg, "polityd -join=memnet://peer@memory")
	})
}
