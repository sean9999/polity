package polity

import (
	"encoding/json"
)

type Peermap map[string]Peer

func (pm Peermap) MarshalJson() ([]byte, error) {
	m := map[string]map[string]string{}
	for k, v := range pm {
		m[k] = v.AsMap()
	}
	return json.Marshal(m)
}

func (pm *Peermap) UnmarshalJson(data []byte) error {
	var m map[string]map[string]string
	err := json.Unmarshal(data, &m)
	if err != nil {
		return err
	}
	ppm := *pm
	for nick, obj := range m {
		p, err := PeerFromHex([]byte(obj["pub"]))
		if err != nil {
			return err
		}
		ppm[nick] = p
	}
	return nil
}
