package lirc

import (
	"log"
	"path/filepath"
)

type remoteButton struct {
	remote string
	button string
}

// Handle is a function that can be registered to handle an lirc Event
type Handle func(LircEvent)

// Handle registers a new event handle
func (l *LircRouter) Handle(remote string, button string, handle Handle) {
	var rb remoteButton

	if remote == "" {
		rb.remote = "*"
	} else {
		rb.remote = remote
	}

	if button == "" {
		rb.button = "*"
	} else {
		rb.button = button
	}

	if l.handlers == nil {
		l.handlers = make(map[remoteButton]Handle)
	}

	l.handlers[rb] = handle
}

func (l *LircRouter) Run() {
	var rb remoteButton

	for {
		event := <-l.receive
		match := 0

		// Check for exakt match
		rb.remote = event.remote
		rb.button = event.button
		if h, ok := l.handlers[rb]; ok {
			h(event)
			continue
		}

		// Check for pattern matches
		for k, h := range l.handlers {
			remote_matched, _ := filepath.Match(k.remote, event.remote)
			button_matched, _ := filepath.Match(k.button, event.button)

			if remote_matched && button_matched {
				h(event)
				match = 1
			}
		}

		if match == 0 {
			log.Println("No match for ", event)
		}
	}
}
