package connection

import (
	"slices"
	"strings"

	"v.io/x/lib/netstate"
)

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
