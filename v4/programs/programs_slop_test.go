package programs

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/subject"
	"github.com/stretchr/testify/assert"
)

type mockProg struct {
	subjects       []subject.Subject
	initCalled     bool
	runCalled      bool
	shutdownCalled bool
}

func (m *mockProg) Init(citizen *polity.Citizen, inbox, outbox chan polity.Envelope, errs chan error) error {
	m.initCalled = true
	return nil
}
func (m *mockProg) Run(ctx context.Context) {
	m.runCalled = true
}
func (m *mockProg) Subjects() []subject.Subject {
	return m.subjects
}
func (m *mockProg) Shutdown() {
	m.shutdownCalled = true
}

func TestPrograms_Slop(t *testing.T) {
	p := &mockProg{subjects: []subject.Subject{"test"}}

	t.Run("Register and ProgramsThatHandle", func(t *testing.T) {
		Register(p)
		handlers := Registry.ProgramsThatHandle("test")
		assert.NotEmpty(t, handlers)
		assert.Equal(t, p, handlers[0].Program)
	})

	t.Run("SurProgram Init/Run/Shutdown", func(t *testing.T) {
		sp := &SurProgram{Program: p}
		citizen := polity.NewCitizen(nil, io.Discard, nil)
		outbox := make(chan polity.Envelope, 1)
		errs := make(chan error, 1)

		err := sp.Init(citizen, outbox, errs)
		assert.NoError(t, err)
		assert.True(t, p.initCalled)

		ctx, cancel := context.WithCancel(context.Background())
		go sp.Run(ctx)

		time.Sleep(10 * time.Millisecond)
		cancel()

		assert.True(t, p.runCalled)
		assert.True(t, p.shutdownCalled)
	})
}
