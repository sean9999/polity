package main

import (
	"testing"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v3"
	"github.com/stretchr/testify/assert"
)

func TestPolity_Slop(t *testing.T) {
	t.Run("Init errors", func(t *testing.T) {
		a := &app{node: nil}
		env := hermeti.TestEnv()
		err := a.Init(&env)
		assert.Error(t, err)
	})

	t.Run("Init success", func(t *testing.T) {
		// Mock node is enough here
		a := &app{node: new(mockNode)}
		env := hermeti.TestEnv()
		env.Args = []string{"polity"}
		err := a.Init(&env)
		assert.NoError(t, err)
	})

	t.Run("Run panics", func(t *testing.T) {
		a := &app{}
		env := hermeti.TestEnv()
		assert.Panics(t, func() {
			a.Run(env)
		})
	})
}

type mockNode struct{ polity.Node }
