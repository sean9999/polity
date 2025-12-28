package main

// bootUp sends a message to self, including a handy dandy join code.
//func bootUp(a *appState, _ hermeti.Env, outbox chan polity.Envelope) {
//	e := a.me.Compose(nil, a.me.Node.Address())
//	e.Letter.SetSubject(subject.BootUp)
//	greeting := fmt.Sprintf("hi! i'm %s. You can join me with:\n\npolityd -join=%s\n", a.me.Oracle.NickName(), a.me.Node.Address())
//	e.Letter.PlainText = []byte(greeting)
//	outbox <- *e
//}
