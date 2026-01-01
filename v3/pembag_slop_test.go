package polity

import (
	"encoding/pem"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPemBag_Slop(t *testing.T) {
	var pb PemBag
	
	block1 := pem.Block{Type: "TEST1", Bytes: []byte("data1")}
	block2 := pem.Block{Type: "TEST1", Bytes: []byte("data2")}
	block3 := pem.Block{Type: "TEST2", Bytes: []byte("data3")}

	// Test Add
	pb.Add("TEST1", block1)
	assert.NotNil(t, pb)
	assert.Equal(t, 1, pb.Size())

	pb.Add("TEST1", block2)
	assert.Equal(t, 2, pb.Size())

	pb.Add("TEST2", block3)
	assert.Equal(t, 3, pb.Size())

	// Test Get
	blocks1, ok1 := pb.Get("TEST1")
	assert.True(t, ok1)
	assert.Equal(t, 2, len(blocks1))

	blocks2, ok2 := pb.Get("TEST2")
	assert.True(t, ok2)
	assert.Equal(t, 1, len(blocks2))

	blocks3, ok3 := pb.Get("NONEXISTENT")
	assert.False(t, ok3)
	assert.Nil(t, blocks3)

	// Test Write
	var pb2 PemBag
	data := pem.EncodeToMemory(&block1)
	data = append(data, pem.EncodeToMemory(&block3)...)
	data = append(data, []byte("some junk at the end")...)

	n, err := pb2.Write(data)
	assert.NoError(t, err)
	assert.Equal(t, len(data)-len("some junk at the end"), n)
	assert.Equal(t, 2, pb2.Size())
	
	b1, _ := pb2.Get("TEST1")
	assert.Equal(t, 1, len(b1))
	assert.Equal(t, block1.Bytes, b1[0].Bytes)

	b2, _ := pb2.Get("TEST2")
	assert.Equal(t, 1, len(b2))
	assert.Equal(t, block3.Bytes, b2[0].Bytes)
}
