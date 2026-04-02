package main

import (
	_ "context"
	"os"
	"testing"
	"testing/cryptotest"
	"time"

	"github.com/sean9999/polity/v4"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v4/network/mem"
	"github.com/stretchr/testify/assert"
)

// a test app uses the mem back-end
func newTestApp(env hermeti.Env) *polityd {
	mother := make(mem.Network)
	a := polityd{
		node: mother.Spawn(),
	}
	citizen := polity.NewCitizen(env.OutStream, a.node)
	a.me = citizen

	return &a
}

func createCitizen(t *testing.T, seed byte, env hermeti.Env) hermeti.CLI[*polityd] {
	cryptotest.SetGlobalRandom(t, uint64(seed))
	app := newTestApp(env)
	cli := hermeti.NewCLI(&env, app)
	return *cli
}

var env hermeti.Env

func setup() error {
	env = hermeti.TestEnv()
	err := env.MountDir("../../testdata")
	if err != nil {
		return err
	}
	env.Args = []string{"polityd"}
	return nil
}

func TestMain(m *testing.M) {

	err := setup()
	if err != nil {
		panic(err)
	}

	exitVal := m.Run()
	os.Exit(exitVal)
}

func TestCitizen_delicateStar_boots(t *testing.T) {
	aliceCli := createCitizen(t, 1, env)
	t.Cleanup(aliceCli.App.me.Shutdown)
	out, err := aliceCli.Env.CaptureOutput()
	if err != nil {
		panic(err)
	}
	go aliceCli.Run()
	time.Sleep(time.Second)
	assert.Contains(t, out.String(), "delicate-star")
	assert.Contains(t, out.String(), "ce083a23682c9d8d00430b0289dce6dd59fed5dcb906d9e66f1148ded2722043e1084cc90e5c218f3eaea876c59356842c618b5f7b67ba8a7296e1e329737ca8")
	assert.Contains(t, out.String(), "polityd -join=")
	aliceJoin := aliceCli.App.me.URL().String()
	assert.Equal(t, "memnet://ce083a23682c9d8d00430b0289dce6dd59fed5dcb906d9e66f1148ded2722043e1084cc90e5c218f3eaea876c59356842c618b5f7b67ba8a7296e1e329737ca8@memory", aliceJoin)
}
