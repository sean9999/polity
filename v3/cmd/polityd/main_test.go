package main

import (
	_ "context"
	"os"
	"testing"
	"time"

	"github.com/sean9999/hermeti"
	"github.com/stretchr/testify/assert"
)

var (
	aliceCli  hermeti.CLI[*appState]
	aliceJoin string
	//bobCli    hermeti.CLI[*appState]
	//bobJoin   string
)

type deterministicRandomness byte

func (d deterministicRandomness) Read(p []byte) (int, error) {
	for i := range p {
		p[i] = byte(d)
	}
	return len(p), nil
}

func createCitizen(seed byte, env hermeti.Env) hermeti.CLI[*appState] {
	randy := deterministicRandomness(seed)
	env.Randomness = randy
	app := newTestApp()
	cli := hermeti.NewCLI(&env, app)
	return *cli
}

func setup() error {

	env := hermeti.TestEnv()

	err := env.MountDir("../../testdata")
	if err != nil {
		return err
	}
	env.Args = []string{"polityd"}
	aliceCli = createCitizen(1, env)
	//bobCli = createCitizen(2, env)
	return nil
}

func teardown() {
	alice := aliceCli.App.me
	alice.Shutdown()
}

func TestMain(m *testing.M) {

	err := setup()
	if err != nil {
		panic(err)
	}

	go aliceCli.Run()
	time.Sleep(250 * time.Millisecond)

	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func TestCitizen_fallingDawn_boots(t *testing.T) {
	out, err := aliceCli.Env.CaptureOutput()
	if err != nil {
		panic(err)
	}

	assert.Contains(t, out.String(), "falling-dawn")
	assert.Contains(t, out.String(), "a4e09292b651c278b9772c569f5fa9bb13d906b46ab68c9df9dc2b4409f8a2098a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5c")
	assert.Contains(t, out.String(), "polityd -join=")
	aliceJoin = aliceCli.App.me.URL().String()
	assert.Equal(t, "memnet://a4e09292b651c278b9772c569f5fa9bb13d906b46ab68c9df9dc2b4409f8a2098a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5c@memory", aliceJoin)
}
