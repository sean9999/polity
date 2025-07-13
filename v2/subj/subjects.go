package subj

import (
	"fmt"
	"strings"
)

type Subject string

func (subj Subject) Equals(thing any) bool {
	a := fmt.Sprintf("%s", subj)
	b := fmt.Sprintf("%s", thing)
	diff := strings.Compare(strings.ToUpper(a), strings.ToUpper(b))
	return diff == 0
}

const (
	Boot                Subject = "Successfully Booted up"
	FriendRequest       Subject = "Will you be my friend?"
	FriendRequestAccept Subject = "Yes. I will be your friend"
	FriendIntroduction  Subject = "You two should be friends"
	Hello               Subject = "Hello. I'm alive"
	HelloBack           Subject = "Hello. Glad you're still alive. I'm alive too"
	WelcomeBack         Subject = "Welcome back. I thought you were dead."
	SoAndSoIsAlive      Subject = "So and so is alive"
	DumpThyself         Subject = "Dump Thyself"
	TellMeEverything    Subject = "Tell me everything you know"
	ThisIsWhatIKnow     Subject = "This is what I know"
	RefuseToDie         Subject = "Fuck you. I won't die"
	KillYourself        Subject = "Kill yourself"
)

func ValidResponses(s Subject) []Subject {

	//	TODO: there is a lot of unnecessary allocation here. Just make a static map
	responses := make([]Subject, 0)

	switch s {
	case FriendRequest:
		responses = append(responses, FriendRequestAccept, FriendIntroduction)
	case KillYourself:
		responses = append(responses, RefuseToDie)
	case Hello:
		responses = append(responses, HelloBack, WelcomeBack)
	case TellMeEverything:
		responses = append(responses, ThisIsWhatIKnow)
	}

	return responses
}
