package polity

import (
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/sean9999/go-oracle"
)

var ZeroConf CitizenConfig
var ErrInvalidConfig = errors.New("invalid config")

type SelfConfig struct {
	oracle.Self
	Address string `toml:"addr" json:"addr"`
}

type CitizenConfig struct {
	handle io.ReadWriter
	Self   SelfConfig          `toml:"self" json:"self"`
	Peers  []map[string]string `toml:"peer" json:"peer"`
}

func (cc *CitizenConfig) String() string {
	b, _ := json.MarshalIndent(cc, "", "\t")
	return string(b)
}

// save config to file or whatever the storage backend is
func (cc *CitizenConfig) Save() error {
	e := toml.NewEncoder(cc.handle)
	e.Encode(cc)
	return nil
}

func ConfigFrom(rw io.ReadWriter) (*CitizenConfig, error) {
	if rw == nil {
		return &ZeroConf, errors.New("nil reader")
	}
	tomlDecoder := toml.NewDecoder(rw)
	var conf CitizenConfig
	_, err := tomlDecoder.Decode(&conf)
	if err != nil {
		return &ZeroConf, err
	}
	conf.handle = rw
	return &conf, nil
}

func ConfigFromFile(path string) (*CitizenConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return ConfigFrom(f)
}
