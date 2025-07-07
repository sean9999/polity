package polity

import (
	"fmt"
	stablemap "github.com/sean9999/go-stable-map"
)

// a KnowlegeBase contains facts about entities in Polity.
type KnowlegeBase[A Addresser] struct {
	Alives *stablemap.StableMap[*Peer[A], bool]
}

func NewKB[A Addresser]() KnowlegeBase[A] {
	kb := KnowlegeBase[A]{
		Alives: stablemap.New[*Peer[A], bool](),
	}
	return kb
}

func (kb *KnowlegeBase[A]) String() string {
	rows := `
Peer	Alive
`

	for k, v := range kb.Alives.Entries() {
		rows += fmt.Sprintf(`
%s	%v
`, k.Nickname(), v)

	}
	return rows
}
