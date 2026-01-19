package audit

// Observer - интерфейс наблюдателя.
type Observer interface {
	Notify(event Event)
}

// Publisher - подписчики.
type Publisher struct {
	observers []Observer
}

// NewPublisher - конструктор.
func NewPublisher() *Publisher {
	return &Publisher{}
}

// Subscribe - подписка на события.
func (p *Publisher) Subscribe(o Observer) {
	p.observers = append(p.observers, o)
}

// Publish - публикация события.
func (p *Publisher) Publish(event Event) {
	for _, o := range p.observers {
		o.Notify(event)
	}
}
