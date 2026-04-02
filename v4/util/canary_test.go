package util

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCanary_happy(t *testing.T) {
	c := NewCanary(t.Context(), 1)
	defer c.Close()
	fmt.Fprintln(c, "hello world")
	data, err := c.Next()
	require.NoError(t, err)
	assert.Equal(t, "hello world\n", string(data))
}

func TestCanary_WatchFor_happy(t *testing.T) {
	ctx := t.Context()
	c := NewCanary(ctx, 16)
	defer c.Close()
	fmt.Fprintln(c, "foo")
	fmt.Fprintln(c, "bar")
	fmt.Fprintln(c, "bing")
	fmt.Fprintln(c, "bat")
	fmt.Fprintln(c, "barf")
	found := c.WatchFor(ctx, "bat")
	assert.True(t, found)
}

func TestCanary_with_cancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	c := NewCanary(ctx, 16)
	fmt.Fprintln(c, "foo")
	fmt.Fprintln(c, "bar")
	fmt.Fprintln(c, "bing")
	fmt.Fprintln(c, "bat")
	fmt.Fprintln(c, "barf")
	cancel()

	//	this is the sad path,
	//	but we also need to use cancellation
	//	because we need a way to say stop watching
	//	else it blocks forever
	found := c.WatchFor(ctx, "dink")
	assert.False(t, found)
	i, err := fmt.Fprintln(c, "plooo")
	assert.Zero(t, i)
	assert.Error(t, err)
}
