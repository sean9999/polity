package subj

type Subject string

const (
	RefuseToDie         Subject = "Fuck you. I won't die"
	KillYourself        Subject = "Kill yourself"
	FriendRequest       Subject = "Will you be my friend?"
	FriendRequestAccept Subject = "Yes. I will be your friend"
	Hello               Subject = "Hello. I'm alive"
	HelloBack           Subject = "Hello. I'm alive too"
	SoAndSoIsAlive      Subject = "So and so is alive"
	Boot                Subject = "Successfully Booted up"
	DumpThyself         Subject = "Dump Thyself"
	TellMeEverything    Subject = "Tell me everything you know"
	ThisIsWhatIKnow     Subject = "This is what I know"
)

func ValidResponses(s Subject) []Subject {
	responses := make([]Subject, 0)

	switch s {
	case FriendRequest:
		responses = append(responses, FriendRequestAccept)
	case KillYourself:
		responses = append(responses, RefuseToDie)
	case Hello:
		responses = append(responses, HelloBack)
	case TellMeEverything:
		responses = append(responses, ThisIsWhatIKnow)
	}

	return responses
}
