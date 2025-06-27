package polity

import (
	"fmt"
	"net"

	stablemap "github.com/sean9999/go-stable-map"
)

// a KnowlegeBase contains facts about entities in Polity.
type KnowlegeBase[A net.Addr] struct {
	Alives *stablemap.StableMap[*Peer[A], bool]
}

func NewKB[A net.Addr]() KnowlegeBase[A] {
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
