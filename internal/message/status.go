package message

type SendingStatus int

const (
	PendingStatus SendingStatus = iota
	SuccessStatus
	FailedStatus
)
