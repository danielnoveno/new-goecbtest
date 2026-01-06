/*
    file:           services/flash/notifier.go
    description:    Layanan backend untuk notifier
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package flash

import (
	"sync"

	"go-ecb/app/types"
)

type notifier struct {
	mu       sync.Mutex
	handlers map[int]func(types.FlashMessage)
	nextID   int
}

// NewNotifier adalah fungsi untuk baru notifier.
func NewNotifier() types.FlashNotifier {
	return &notifier{
		handlers: make(map[int]func(types.FlashMessage)),
	}
}

// Notify adalah fungsi untuk notify.
func (n *notifier) Notify(msg types.FlashMessage) {
	n.mu.Lock()
	handlers := make([]func(types.FlashMessage), 0, len(n.handlers))
	for _, handler := range n.handlers {
		handlers = append(handlers, handler)
	}
	n.mu.Unlock()

	for _, handler := range handlers {
		handler(msg)
	}
}

// Subscribe adalah fungsi untuk subscribe.
func (n *notifier) Subscribe(handler func(types.FlashMessage)) func() {
	n.mu.Lock()
	id := n.nextID
	n.nextID++
	n.handlers[id] = handler
	n.mu.Unlock()

	return func() {
		n.mu.Lock()
		delete(n.handlers, id)
		n.mu.Unlock()
	}
}
