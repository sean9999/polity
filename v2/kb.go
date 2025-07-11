package polity

import (
	"fmt"
	stablemap "github.com/sean9999/go-stable-map"
)

// A KnowledgeBase contains facts about entities in Polity.
type KnowledgeBase[A Addresser] struct {
	lives      *stablemap.StableMap[*Peer[A], bool]
	LiveEvents chan string
}

func NewKB[A Addresser]() KnowledgeBase[A] {
	kb := KnowledgeBase[A]{
		lives:      stablemap.New[*Peer[A], bool](),
		LiveEvents: make(chan string),
	}
	return kb
}

func (kb *KnowledgeBase[A]) UpdateAlives(p *Peer[A], alive bool) error {
	err := kb.lives.Set(p, alive)
	if err != nil {
		return err
	}
	adjective := "dead"
	if alive {
		adjective = "alive"
	}
	kb.LiveEvents <- fmt.Sprintf("I think %s\tis %s.", p.Nickname(), adjective)
	return nil
}

func (kb *KnowledgeBase[A]) String() string {
	rows := ""
	for k, v := range kb.lives.Entries() {
		rows += fmt.Sprintf("%s\tis\t%v\n", k.Nickname(), v)
	}
	return rows
}
