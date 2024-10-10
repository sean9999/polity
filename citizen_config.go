package polity

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	realfs "github.com/sean9999/go-real-fs"
)

var ZeroConf CitizenConfig
var ErrInvalidConfig = errors.New("invalid config")

// a SelfConfig is an [oracle.Self] with an address
type SelfConfig struct {
	Nickname   string     `json:"nick"`
	PublicKey  string     `json:"pub"`
	Privatekey string     `json:"priv"`
	Addresses  AddressMap `json:"addrs"`
}

// a CitizenConfig is a SelfConfig, along with it's peers and a file handle
type CitizenConfig struct {
	Self  SelfConfig   `json:"self"`
	Peers *AddressBook `json:"peers"`
}

func (cc *CitizenConfig) String() string {
	b, _ := json.MarshalIndent(cc, "", "\t")
	return string(b)
}

// write the CitizenConfig
func (cc *CitizenConfig) Export(w io.Writer) error {
	e := json.NewEncoder(w)
	e.SetIndent("", "\t")
	return e.Encode(cc)
}

func ConfigFrom(rw io.ReadWriter) (*CitizenConfig, error) {
	if rw == nil {
		return &ZeroConf, errors.New("nil reader")
	}
	jsonDecoder := json.NewDecoder(rw)
	conf := CitizenConfig{
		Self:  SelfConfig{},
		Peers: &AddressBook{},
	}
	err := jsonDecoder.Decode(&conf)
	if err != nil {
		return &ZeroConf, err
	}
	return &conf, nil
}

func ConfigFromFile(filesystem realfs.WritableFs, path string) (*CitizenConfig, error) {
	f, err := filesystem.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	return ConfigFrom(f)
}
