package polity

import (
	"encoding/pem"
)

type PemBag map[string][]pem.Block

func (pb *PemBag) Add(key string, block pem.Block) {
	// Ensure the map is initialized so the zero value of PemBag works
	if *pb == nil {
		*pb = make(PemBag)
	}
	b := *pb
	b[key] = append(b[key], block)
}

func (pb *PemBag) Get(key string) ([]pem.Block, bool) {
	m := *pb
	thing, ok := m[key]
	return thing, ok
}

func (pb *PemBag) Size() int {
	n := 0
	for _, blocks := range *pb {
		n += len(blocks)
	}
	return n
}

func cumulativeWrite(bytesWritten int, blocks []*pem.Block, data []byte) (int, []*pem.Block, []byte) {
	block, rest := pem.Decode(data)
	if block == nil {
		return bytesWritten, blocks, nil
	}
	blocks = append(blocks, block)
	bytesWritten = bytesWritten + len(data) - len(rest)
	return cumulativeWrite(bytesWritten, blocks, rest)
}

func (pb *PemBag) Write(p []byte) (int, error) {
	bytesWritten, blocks, _ := cumulativeWrite(0, make([]*pem.Block, 0), p)
	for _, block := range blocks {
		pb.Add(block.Type, *block)
	}
	return bytesWritten, nil
}
