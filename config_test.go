package polity

import (
	"testing"

	realfs "github.com/sean9999/go-real-fs"
)

var realFileSytem = realfs.NewWritable()

func TestSave(t *testing.T) {

	// pwd, err := os.Getwd()
	// if err != nil {
	// 	t.Error(err)
	// }

	conf, err := ConfigFromFile(realFileSytem, "testdata/falling-wave.toml")
	if err != nil {
		t.Error(err)
	}

	if conf.Self.Nickname != "falling-wave" {
		t.Errorf("expected falling-wave, got %q", conf.Self.Nickname)
	}

	want := "asdfsdf"
	got := conf.String()
	if got != want {
		t.Errorf("got %s but wanted %s", got, want)
	}

}
