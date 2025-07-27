package subj

import (
	"fmt"
	"strings"
)

// A Subject is a type of message. It usually wants a response
type Subject string

// An Ack is functionally a Subject, but used as a response to a Subject and doesn't itself want a response
type Ack = Subject

// A Command is functionally a Subject, but used for sending commands to oneself.
type Command = Subject

func (subj1 Subject) Equals(subj2 any) bool {
	a := fmt.Sprintf("%s", subj1)
	b := fmt.Sprintf("%s", subj2)
	diff := strings.Compare(strings.ToUpper(a), strings.ToUpper(b))
	return diff == 0
}

const (
	FriendRequest          Subject = "Will you be my friend?"
	FriendRequestAccept    Subject = "Yes. I will be your friend"
	FriendsIntroduction    Subject = "You two should be friends"
	Hello                  Subject = "Hello. I'm alive"
	SoAndSoIsAlive         Subject = "So and so is alive"
	DumpThyself            Subject = "Dump Thyself"
	TellMeEverything       Subject = "Tell me everything you know"
	KillYourself           Subject = "Kill yourself, gracefully if possible"
	Sleep                  Subject = "Go to sleep. Stay alive but stop responding to messages"
	WakeUp                 Subject = "Wake from sleep"
	IWantToMeetYourFriends Subject = "i want to meet your friends"
	ByeBye                 Subject = "Good bye. I'm going away now"
)

const (
	HelloBack        Ack = "Hello. Glad you're still alive. I'm alive too"
	WelcomeBack      Ack = "Welcome back. I thought you were dead."
	ThisIsWhatIKnow  Ack = "This is what I know"
	RefuseToDie      Ack = "Fuck you. I won't die"
	HereAreMyFriends Ack = "Here are my friends"
	IHaveNoFriends   Ack = "I have no friends"
)

const (
	CmdBoot         Command = "Successfully Booted up"
	CmdBroadcast    Command = "Say hello to all your friends"
	CmdMakeFriends  Command = "ask all your friends who their friends are, and then make friends with them"
	CmdEveryoneDump Command = "tell your friends to dump themselves"
)

func ValidResponses(s Subject) []Subject {

	//	TODO: there is a lot of unnecessary allocation here. Just make a static map
	responses := make([]Subject, 0)

	switch s {
	case FriendRequest:
		responses = append(responses, FriendRequestAccept, FriendsIntroduction)
	case KillYourself:
		responses = append(responses, RefuseToDie)
	case Hello:
		responses = append(responses, HelloBack, WelcomeBack)
	case TellMeEverything:
		responses = append(responses, ThisIsWhatIKnow)
	case CmdMakeFriends:
		responses = append(responses, IWantToMeetYourFriends)
	case IWantToMeetYourFriends:
		responses = append(responses, HereAreMyFriends, IHaveNoFriends)
	}

	return responses
}
