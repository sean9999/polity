package main

func processEnvelope(me Node, inMail Envelope) {

	switch inMail.Message.Subject {
	case "will you be my friend?":
		if inMail.Verify() {
			msg := NewMessage("yes.", "I will be your friend.", inMail.Message.Id)
			go me.Spool(msg, inMail.From)
			me.AddFriend(inMail.From)
		}
	}

}
