package main

import "encoding/json"

func (e Envelope) MarshalWireFormat() ([]byte, error) {
	// @todo: develop a more efficient wire format
	return json.Marshal(e)
}

func (e *Envelope) UnmarshalWireFormat(bytes []byte) error {
	// @todo: develop a more efficient wire format
	return json.Unmarshal(bytes, &e)
}
