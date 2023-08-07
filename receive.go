package main

func (n Node) Receive(bin []byte) (Envelope, error) {

	var e Envelope
	err := e.UnmarshalBinary(bin)
	return e, err

}
