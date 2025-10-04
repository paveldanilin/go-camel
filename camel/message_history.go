package camel

import "time"

// The MessageHistory maintains a list of all components that the message passed through.
// Every component that processes the message (including the originator) adds one entry to the list.
// The MessageHistory should be part of the message header because it contains system-specific control information.
// Keeping this information in the header separates it from the message body that contains application specific data.
type MessageHistory struct {
	time        time.Time
	elapsedTime int64
	routeName   string
	stepName    string
}

func (mh *MessageHistory) ElapsedTime() int64 {
	return mh.elapsedTime
}

func (mh *MessageHistory) Time() time.Time {
	return mh.time
}

func (mh *MessageHistory) RouteName() string {
	return ""
}

func (mh *MessageHistory) StepName() string {
	return mh.stepName
}

func (mh *MessageHistory) Message() *Message {
	return nil
}
