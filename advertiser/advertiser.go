package advertiser

import (
	"os"

	"github.com/hashicorp/mdns"
)

func Advertiser() {
	// Setup our service export
	host, _ := os.Hostname()
	info := []string{"My Polity service"}
	service, _ := mdns.NewMDNSService(host, "_polity._tcp", "", "", 8000, nil, info)

	// Create the mDNS server, defer shutdown
	server, _ := mdns.NewServer(&mdns.Config{Zone: service})
	defer server.Shutdown()
}
