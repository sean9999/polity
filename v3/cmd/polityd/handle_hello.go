package main

import (
	"context"
	"fmt"

	"github.com/sean9999/go-oracle/v3/delphi"
	"github.com/sean9999/hermeti"
	"github.com/sean9999/polity/v3"
	"github.com/sean9999/polity/v3/subject"
)

func handleHello(env hermeti.Env, me *polity.Citizen, e polity.Envelope) {

	// validate Letter. if not valid. drop it and we're done.
	err := e.Letter.Verify(me.KeyPair)
	if err != nil {
		fmt.Fprintln(env.OutStream, "signature failed")
	}

	key, err := delphi.KeyFromString(e.Sender.User.Username())
	if err != nil {
		fmt.Fprintf(env.OutStream, "sender:\t%s\n", e.Sender.String())
	}
	pubKey := delphi.PublicKey(key)

	var subj subject.Subject

	//	do we know this person?
	peer := me.Peers.Get(pubKey)
	if peer != nil {
		//	yes. mark them as alive. tell my friends
		err := me.Profiles.SetAliveness(pubKey, true)
		if err != nil {
			return
		}
		subj = subject.SoAndSoIsAive
	} else {
		//	no. add them. tell my friends
		peer = polity.PeerFromKey(pubKey)
		me.Peers.Add(*peer, func() {
			_ = me.Profiles.SetAliveness(pubKey, true)
		})
		subj = subject.SoAndSoIsNew
	}

	letter := polity.NewLetter(nil)
	letter.SetSubject(subj.String())
	letter.PlainText = peer.Serialize()
	_ = letter.Sign(env.Randomness, me.KeyPair)

	ctx := context.Background()
	recipients := me.Peers.Minus(peer.PublicKey)
	err = me.Announce(ctx, env.Randomness, letter, recipients.URLs())

}
