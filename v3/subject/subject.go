package subject

type Subject string

func (s Subject) String() string {
	return string(s)
}

func From(s string) Subject {
	return Subject(s)
}

const (
	BootUp              = "boot up and join string"
	IWantToBeYourFriend = "I want to be your friend"
	SoAndSoIsAive       = "So and so is alive"
	SoAndSoIsNew        = "So and so is new"
	HereAreMyFriends    = "Here are my friends"
	IamAlive            = "I'm alive. "
	DieNow              = "go away"
)
