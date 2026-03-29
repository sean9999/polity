package bonjour

import (
	"context"
	"strconv"

	"github.com/grandcat/zeroconf"
	"github.com/sean9999/go-oracle/v4/delphi"
	"github.com/sean9999/polity/v4"
	"github.com/sean9999/polity/v4/programs"
)

const (
	serviceName = "_polity._udp"
	domain      = "local."
)

var _ programs.Program = (*prog)(nil)

type prog struct {
	citizen  *polity.Citizen
	errs     chan error
	outbox   chan polity.Envelope
	inbox    chan polity.Envelope
	resolver *zeroconf.Resolver
	server   *zeroconf.Server
	adverts  chan string
}

func (h *prog) Init(citizen *polity.Citizen, inbox chan polity.Envelope, outbox chan polity.Envelope, errs chan error) error {

	// Discover all services on the network (e.g. _workstation._tcp)
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return err
	}

	port, err := strconv.Atoi(citizen.URL().Port())
	if err != nil {
		return err
	}

	server, err := advertise(citizen.KeyPair.PublicKey(), port)
	if err != nil {
		return err
	}

	h.server = server
	h.resolver = resolver
	h.citizen = citizen
	h.errs = errs
	h.outbox = outbox
	h.inbox = inbox
	h.adverts = make(chan string, 4)
	return nil
}

func (h *prog) Subjects() []polity.Subject {
	return []polity.Subject{
		polity.SubjBootUp,
	}
}

func (h *prog) Run(ctx context.Context) {

	_ = <-h.inbox

	for ad := range h.adverts {
		h.citizen.Log.Println("advertising", ad)

	}

}

func (h *prog) Shutdown() {
	close(h.adverts)
	h.server.Shutdown()
	h.citizen.Log.Println("prog shutdown")
}

func init() {
	programs.Register(new(prog))
}

func advertise(pubKey delphi.PublicKey, port int) (*zeroconf.Server, error) {

	kv := []string{
		"pubkey=" + pubKey.String(),
	}

	server, err := zeroconf.Register(pubKey.Nickname(), serviceName, domain, port, kv, nil)
	if err != nil {
		return nil, err
	}

	return server, nil
}
