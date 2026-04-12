package main

import (
	"context"
	_ "context"
	"testing"
	"testing/cryptotest"
	"time"

	"sync"

	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v4"
	"github.com/sean9999/polity/v4/network/mem"
	"github.com/sean9999/polity/v4/programs"
	"github.com/sean9999/polity/v4/util"
	"github.com/stretchr/testify/assert"
)

// a test app uses the mem back-end
func newTestApp(env hermeti.Env) *appState {
	mother := make(mem.Network)
	a := appState{
		node: mother.Spawn(),
	}
	citizen := polity.NewCitizen(env.OutStream, a.node)
	a.me = citizen

	return &a
}

// refreshRegistry replaces every SurProgram wrapper in the global registry with
// a fresh one (new Inbox channel). Goroutines left over from a previous test
// still hold their old *SurProgram pointers, so when they eventually call
// Shutdown() → close(sp.Inbox) they touch their own orphaned channel and cannot
// double-close a channel owned by the current test.
//
// savedProgs captures the full program list on first call, before any test can
// deregister programs (e.g. bootup deregisters itself after running).
var (
	savedProgs []programs.Program
	saveOnce   sync.Once
)

func refreshRegistry() {
	saveOnce.Do(func() {
		for sp := range programs.Registry.Programs {
			savedProgs = append(savedProgs, sp.Program)
		}
	})
	var current []*programs.SurProgram
	for sp := range programs.Registry.Programs {
		current = append(current, sp)
	}
	for _, sp := range current {
		programs.Deregister(sp)
	}
	for _, p := range savedProgs {
		programs.Register(p)
	}
}

func createCitizen(t *testing.T, seed byte, env hermeti.Env) hermeti.CLI[*appState] {
	refreshRegistry()
	cryptotest.SetGlobalRandom(t, uint64(seed))
	app := newTestApp(env)
	cli := hermeti.NewCLI(&env, app)
	return *cli
}

func createEnv(t testing.TB) hermeti.Env {
	t.Helper()
	env := hermeti.TestEnv()
	env.Args = []string{"appState"}
	return env
}

func TestCitizen_delicateStar_boots(t *testing.T) {
	env := createEnv(t)
	alice := createCitizen(t, 1, env)
	//t.Cleanup(alice.App.me.Shutdown)
	out, err := alice.Env.CaptureOutput()
	if err != nil {
		panic(err)
	}
	go alice.Run()
	time.Sleep(time.Second)
	assert.Contains(t, out.String(), "delicate-star")
	assert.Contains(t, out.String(), "ce083a23682c9d8d00430b0289dce6dd59fed5dcb906d9e66f1148ded2722043e1084cc90e5c218f3eaea876c59356842c618b5f7b67ba8a7296e1e329737ca8")
	assert.Contains(t, out.String(), "appState -join=")
	join := alice.App.me.URL().String()
	assert.Equal(t, "memnet://ce083a23682c9d8d00430b0289dce6dd59fed5dcb906d9e66f1148ded2722043e1084cc90e5c218f3eaea876c59356842c618b5f7b67ba8a7296e1e329737ca8@memory", join)
}

func TestCitizen_output_happy(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 4*time.Second)
	defer cancel()
	canary := util.NewCanary(t.Context(), 16)
	env := createEnv(t)
	env.OutStream = canary
	bob := createCitizen(t, 2, env)
	go bob.Run()
	found := canary.WatchFor(ctx, "memnet")
	assert.True(t, found)
}

func TestCitizen_output_sad(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	canary := util.NewCanary(ctx, 16)
	env := createEnv(t)
	env.OutStream = canary
	carl := createCitizen(t, 3, env)
	//t.Cleanup(carl.App.me.Shutdown)
	go carl.Run()
	found := canary.WatchFor(ctx, "a little dab'll do ya")
	assert.False(t, found)
}

func TestCLI_no_args(t *testing.T) {
	assert.Panics(t, func() {
		bob := createCitizen(t, 0, hermeti.TestEnv())
		bob.Run()
	})
}
