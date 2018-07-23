package godatabus

type Event interface {
	Message
}

type BaseEvent struct {
	BaseMessage
}