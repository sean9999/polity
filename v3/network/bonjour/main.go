package main

import (
	"fmt"
	"net"
	"sync"

	"github.com/hashicorp/mdns"
)

func main() {

	info := []string{"My awesome service"}

	port := 12345

	ips := []net.IP{
		net.IPv4(10, 0, 0, 45),
		net.IPv4(10, 0, 0, 68),
	}

	service, err := mdns.NewMDNSService("billie", "_polity._udp", "local.", "", port, ips, info)
	if err != nil {
		panic(err)
	}

	server, err := mdns.NewServer(&mdns.Config{Zone: service})
	if err != nil {
		panic(err)
	}
	defer server.Shutdown()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go listen(wg)

	wg.Wait()

}

func listen(wg *sync.WaitGroup) {
	// Make a channel for results and start listening
	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			fmt.Printf("Got new entry: %v\n", entry)
		}
	}()

	// Start the lookup
	mdns.Lookup("_polity._udpx", entriesCh)
	close(entriesCh)
	wg.Done()
}
