package godatabus

type EventBus struct {
	MessageBusSupportingMiddleware
}

func NewEventBus() *EventBus {
	return &EventBus{}
}

type CommandBus struct {
	MessageBusSupportingMiddleware
}

func NewComanndBus() *CommandBus {
	return &CommandBus{}
}