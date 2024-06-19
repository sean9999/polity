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
	SubjHowAreYou      Subject = "how are you?"
	SubjGoodThanks     Subject = "good, thanks"
	SubjWhoDoYouKnow   Subject = "who do you know?"
)
