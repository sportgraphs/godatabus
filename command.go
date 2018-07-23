package godatabus

type Command interface {
	Message
}

// BaseMessage contains the basic info that all commands should have
type BaseCommand struct {
	BaseMessage
}
