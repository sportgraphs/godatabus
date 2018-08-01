package godatabus

type RoutingKeyResolver interface {
	Resolve(message Message) string
}

type SimpleRoutingKeyResolver struct {
}

func NewSimpleRoutingKeyResolver() *SimpleRoutingKeyResolver {
	return &SimpleRoutingKeyResolver{}
}

func (s *SimpleRoutingKeyResolver) Resolve(message Message) string {
	return message.GetType()
}