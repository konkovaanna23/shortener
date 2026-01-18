package audit

import "fmt"

type Observer interface {
	Notify(event Event)
}

type Publisher struct {
	observers []Observer
}

func NewPublisher() *Publisher {
	return &Publisher{}
}

func (p *Publisher) Subscribe(o Observer) {
	p.observers = append(p.observers, o)
}

func (p *Publisher) Publish(event Event) {
	fmt.Println("Запрос на публикацию")
	fmt.Println("Подписчиков", len(p.observers))
	for _, o := range p.observers {
		o.Notify(event)
	}
}
