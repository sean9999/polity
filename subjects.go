package polity

// Subjects are well-known values that indicate a type of Message and how it should be handled
type Subject string

func (subj Subject) String() string {
	return string(subj)
}

const (
	SubjAssertion    Subject = "I assert myself"
	SubjKillYourself Subject = "kill yourself"
	SubjGoProverb    Subject = "a go proverb is"
	SubjGenericMsg   Subject = "generic message"
)
