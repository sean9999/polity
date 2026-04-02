package polity

type Subject string

func (s Subject) String() string {
	return string(s)
}

const (
	SubjDieNow            Subject = "die now"
	SubjBootUp            Subject = "boot up"
	SubjLetsBeFriends     Subject = "let's be friends"
	SubjIamAlive          Subject = "i am alive"
	SubjTheseAreMyFriends Subject = "these are my friends"
	SubjHowAreYou         Subject = "how are you?"
	SubjStatusReport      Subject = "status report"
	SubjTrustUpdate       Subject = "trust update"
)
