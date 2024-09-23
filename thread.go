package polity

import "github.com/google/uuid"

type ThreadLifecycleState uint8

const (
	LifecycleStateUninitialized ThreadLifecycleState = iota
	LifecycleStateOpen
	LifecycleStateClosed
)

// A thread is a logical grouping of ordered messages.
// It has a definite beginning and definite end.
type Thread struct {
	Id             uuid.UUID
	Messages       []Message
	LifecycleState ThreadLifecycleState
}

var NoThread = Thread{}

// constructor
func NewThread() (Thread, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return NoThread, err
	}
	messages := make([]Message, 0)
	t := Thread{
		Id:             id,
		Messages:       messages,
		LifecycleState: LifecycleStateOpen,
	}
	return t, nil
}
