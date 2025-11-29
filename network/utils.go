package network

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"

	"v.io/x/lib/netstate"
)

var ErrNetworkUp = errors.New("can't bring network up")
var ErrConnection = errors.New("couldn't create connection")

var ErrNotImplemented = errors.New("not implemented")

func isNotPublic(a netstate.Address) bool {
	return !netstate.IsPublicUnicastIPv6(a)
}

func isNotLoopBack(a netstate.Address) bool {
	isLoop := strings.Contains(a.Interface().Flags().String(), "loopback")
	return !isLoop
}

func isBroadCast(a netstate.Address) bool {
	isBroad := strings.Contains(a.Interface().Flags().String(), "broadcast")
	return isBroad
}

func isLinkLocalAndRoutable(a netstate.Address) bool {
	return isNotPublic(a) && isBroadCast(a) && isNotLoopBack(a)
}

func isNotWeird(a netstate.Address) bool {
	weirdDevices := []string{"ap1", "llw0", "awdl0"}
	ifaceName := a.Interface().Name()
	return !slices.Contains(weirdDevices, ifaceName)
}

type mac [6]byte

func uint64ToMac(i uint64) mac {
	var m mac
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	copy(m[:], b)
	return m
}

func join(s []string) string {
	return strings.Join(s, "")
}

func (m mac) String() string {
	str := hex.EncodeToString(m[:])
	letters := strings.Split(str, "")
	return fmt.Sprintf("%s:%s:%s:%s:%s:%s", join(letters[0:2]), join(letters[2:4]), join(letters[4:6]), join(letters[6:8]), join(letters[8:10]), join(letters[10:12]))
}

func (m mac) Postfix() string {
	str := hex.EncodeToString(m[:])
	letters := strings.Split(str, "")
	return fmt.Sprintf("%s:%s:%s", join(letters[0:4]), join(letters[4:8]), join(letters[8:12]))
}
