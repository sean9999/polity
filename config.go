package polity

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"

	"github.com/sean9999/go-oracle"
	realfs "github.com/sean9999/go-real-fs"
	"github.com/sean9999/polity/connection"
)

var ZeroConf CitizenConfig
var ErrInvalidConfig = errors.New("invalid config")

// a SelfConfig is an [oracle.Self] with an address
type SelfConfig struct {
	oracle.SelfConfig
	Address net.Addr `json:"addr"`
}

// a CitizenConfig is a SelfConfig, along with it's peers and a file handle
type CitizenConfig struct {
	connection connection.Connection
	handle     io.ReadWriter
	Self       SelfConfig            `json:"self"`
	Peers      map[string]peerConfig `json:"peers"`
}

func (cc *CitizenConfig) String() string {
	b, _ := json.MarshalIndent(cc, "", "\t")
	return string(b)
}

// save config to file or whatever the storage backend is
func (cc *CitizenConfig) Save() error {
	e := json.NewEncoder(cc.handle)
	e.SetIndent("", "\t")
	e.Encode(cc)
	return nil
}

func ConfigFrom(rw io.ReadWriter) (*CitizenConfig, error) {
	if rw == nil {
		return &ZeroConf, errors.New("nil reader")
	}
	jsonDecoder := json.NewDecoder(rw)
	var conf CitizenConfig
	err := jsonDecoder.Decode(&conf)
	if err != nil {
		return &ZeroConf, err
	}
	conf.handle = rw
	return &conf, nil
}

func ConfigFromFile(filesystem realfs.WritableFs, path string) (*CitizenConfig, error) {
	f, err := filesystem.OpenFile(path, os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	return ConfigFrom(f)
}
