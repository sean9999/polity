package main

import (
	"crypto/rand"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/google/uuid"
	"github.com/sean9999/go-delphi"
	"github.com/sean9999/polity/v2"
)

var NoUUID uuid.UUID

type peerToJoin struct {
	peer *polity.Peer[*net.UDPAddr]
}

func (p *peerToJoin) String() string {
	if p.peer != nil {
		return p.peer.String()
	}
	return ""
}

func (p *peerToJoin) Set(s string) error {
	u, err := url.Parse(fmt.Sprintf("%s://%s", "udp", s))

	if err != nil {
		return err
	}
	port, err := strconv.Atoi(u.Port())
	if err != nil {
		return err
	}

	ip := net.ParseIP(u.Hostname())

	addr := &net.UDPAddr{
		IP:   ip,
		Port: port,
	}
	peer := polity.NewPeer[*net.UDPAddr]()
	peer.Addr = addr

	//	TODO: ensure this is valid hex, of the right size, and marshals into a valid public key
	hexStr := u.User.Username()
	peer.Peer.Peer = delphi.KeyFromHex(hexStr)
	p.peer = peer
	return nil
}

func main() {
	done := make(chan error)
	friend := new(peerToJoin)
	f := flag.NewFlagSet("fset", flag.ContinueOnError)
	f.Var(friend, "join", "node to join")
	err := f.Parse(os.Args[1:])

	if err != nil {
		done <- err
		return
	}

	p, err := polity.NewPrincipal(rand.Reader, new(polity.LocalUDP4Net))
	if err != nil {
		panic(err)
	}

	go func() {
		for e := range p.Inbox {
			msg := e.Message
			subj := e.Message.Subject
			var body string
			if msg.Encrypted() {
				//	TODO: decrypt message.
				body = fmt.Sprintf("%x", msg.CipherText)
			} else {
				body = string(e.Message.PlainText)
			}

			//	log out message
			color.Magenta("\n#\t%s", string(subj))
			color.Cyan("MsgId:\t%s\n", e.ID)
			color.Cyan("Thread:\t%s\n", e.Thread)
			color.Blue("Signed:\t%v\n", msg.Verify())
			color.Blue("Enc:\t%v\n", msg.Encrypted())
			color.Green("From:\t%s@%s\n", e.Message.SenderKey.Nickname(), e.Sender.Addr.String())
			color.Green("To:\t%s@%s\n", e.Message.RecipientKey.Nickname(), e.Recipient.Addr.String())
			fmt.Println(body)

			// if subj.Equals("friend request") {
			// 	p.Peers.Set(e.Message.SenderKey.Nickname(), )
			// }

			if subj.Equals("DIE NOW") {
				time.Sleep(time.Second * 1)
				close(p.Inbox)
			}
		}
		done <- errors.New("goodbye!")
	}()

	message := fmt.Sprintf(`WOWZA!
I'm %s.
To join me, do:
polityd -join %s
`, p.Nickname(), p.AsPeer().String())

	e := p.Compose([]byte(message), p.AsPeer(), polity.NilId)
	e.Subject("boot up")
	_, err = p.Send(e)

	if err != nil {
		done <- err
	}

	if friend.peer != nil {
		j := p.Compose([]byte("i want to join you"), friend.peer, polity.MessageId(NoUUID))
		j.Subject("friend request")
		//	a friend request must be signed
		err = j.Message.Sign(rand.Reader, p)
		if err != nil {
			done <- err
		}
		_, err = p.Send(j)

		time.Sleep(time.Second * 5)
		j = p.Compose([]byte("i want you to die"), friend.peer, polity.NilId)
		j.Subject("die now")
		err = j.Message.Sign(rand.Reader, p)
		if err != nil {
			done <- err
		}
		_, err = p.Send(j)

	}

	err = <-done
	fmt.Fprintln(os.Stderr, err)

}
