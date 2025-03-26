package event

type Handler func(event Event)

type EventBus interface {
	Publish(event Event)
	Register(eventType Type, handler Handler)
}

type EventBussin struct {
	handlers map[Type][]Handler
}

func NewEventBussin() *EventBussin {
	return &EventBussin{
		handlers: map[Type][]Handler{},
	}
}

func (eb *EventBussin) Publish(ev Event) {
	for _, handler := range eb.handlers[ev.Type()] {
		handler(ev)
	}
}

func (eb *EventBussin) Register(eventType Type, handler Handler) {
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
}
