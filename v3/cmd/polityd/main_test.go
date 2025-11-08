package main

import (
	"bytes"
	_ "context"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v3/network/mem"
	"github.com/stretchr/testify/assert"
)

var (
	fallingDawnCli  hermeti.CLI[*appState]
	fallingDawnJoin string
)

var mother *mem.Network

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
	app := newTestApp(mother)
	cli := hermeti.NewCLI(&env, app)
	return *cli
}

func setup() hermeti.CLI[*appState] {
	mother = mem.NewNetwork()
	env := hermeti.TestEnv()
	err := env.MountDir("../../testdata")
	if err != nil {
		panic(err)
	}
	env.Args = []string{"polityd"}
	fallingDawnCli = createCitizen(1, env)
	return fallingDawnCli
}

func teardown() {
	alice := fallingDawnCli.App.me
	e := alice.Compose(nil, alice.Address())
	e.Letter.SetSubject("go away")
	e.Letter.PlainText = []byte("go away")
	alice.Send(nil, nil, e.Letter, e.Recipient)
}

func TestMain(m *testing.M) {

	fallingDawnCli = setup()

	go fallingDawnCli.Run()
	time.Sleep(250 * time.Millisecond)

	out, err := fallingDawnCli.Env.CaptureOutput()
	if err != nil {
		panic(err)
	}

	c, err := regexp.Compile(` -join=(.*)`)

	if err != nil {
		panic(err)
	}

	matches := c.FindAllStringSubmatch(out.String(), -1)

	if len(matches) < 1 || len(matches[0]) < 2 {
		panic("no matches")
	}

	fallingDawnJoin = matches[0][1]

	exitVal := m.Run()
	teardown()
	os.Exit(exitVal)
}

func TestCitizen_fallingDawn_boots(t *testing.T) {
	out := fallingDawnCli.Env.OutStream.(*bytes.Buffer)
	assert.Contains(t, out.String(), "falling-dawn")
	assert.Contains(t, out.String(), "a4e09292b651c278b9772c569f5fa9bb13d906b46ab68c9df9dc2b4409f8a2098a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5c")
	assert.Contains(t, out.String(), "polityd -join=")
	assert.Equal(t, "memnet://a4e09292b651c278b9772c569f5fa9bb13d906b46ab68c9df9dc2b4409f8a2098a88e3dd7409f195fd52db2d3cba5d72ca6709bf1d94121bf3748801b40f6f5c@memory", fallingDawnJoin)
}
