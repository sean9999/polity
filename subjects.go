package polity

// Subjects are well-known values that indicate a type of Message and how it should be handled
type Subject string

func (subj Subject) String() string {
	return string(subj)
}

const (
	NoSubject          Subject = ""
	SubjAssertion      Subject = "I assert myself"
	SubjKillYourself   Subject = "kill yourself"
	SubjGoProverb      Subject = "a go proverb is"
	SubjGenericMsg     Subject = "generic message"
	SubjStartMarcoPolo Subject = "do you want to play marco polo?"
	SubjMarco          Subject = "marco!"
	SubjPolo           Subject = "polo!"
	SubjHelloSelf      Subject = "hello to myself"
	SubjHowdee         Subject = "hi, how are you?"
	SubjGoodThanks     Subject = "i am operating within normal parameters"
	SubjStatusReport   Subject = "this is my status report"
	SubjWhoDoYouKnow   Subject = "who do you know?"
	SubjImBack         Subject = "I'm back. What's new?"
	SubjWelcomeBack    Subject = "welcome back. Here's what's new"
	SubjectNewFriend   Subject = "I have a new friend."
)
