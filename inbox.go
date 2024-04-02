package main

func processEnvelope(me Node, inMail Envelope) {

	switch inMail.Message.Subject {
	case "will you be my friend?":
		if inMail.Verify() {

			affirmative := []byte("I will be your friend")

			msg := NewMessage("yes.", affirmative, inMail.Message.Id)
			go me.Spool(msg, inMail.From)
			me.AddFriend(inMail.From)
		}
	}

}
