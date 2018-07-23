package godatabus

type Message interface {
	GetType() string
	IsValid() bool
}

// BaseMessage contains the basic info that all commands should have
type BaseMessage struct {
}

// GetType returns the command type
func (b BaseMessage) GetType() string {
	return "base"
}

// IsValid checks validates the message
func (b BaseMessage) IsValid() bool {
	return true
}
