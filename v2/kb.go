package polity

import (
	"fmt"
	stablemap "github.com/sean9999/go-stable-map"
)

// a KnowledgeBase contains facts about entities in Polity.
type KnowledgeBase[A Addresser] struct {
	Lives *stablemap.StableMap[*Peer[A], bool]
}

func NewKB[A Addresser]() KnowledgeBase[A] {
	kb := KnowledgeBase[A]{
		Lives: stablemap.New[*Peer[A], bool](),
	}
	return kb
}

func (kb *KnowledgeBase[A]) String() string {
	rows := `
Peer	Alive
`
	for k, v := range kb.Lives.Entries() {
		rows += fmt.Sprintf(`
%s	%v
`, k.Nickname(), v)

	}
	return rows
}
