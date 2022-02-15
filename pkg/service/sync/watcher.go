package secret

import (
	"sync"

	"k8s.io/apimachinery/pkg/watch"
)

// Implements watch.Interface
type PollWatcher struct {
	eventChan chan watch.Event
	stopped   bool
	sync.Mutex
}

// Stop stops poll watcher and closes event channel.
func (p *PollWatcher) Stop() {
	p.Lock()
	defer p.Unlock()
	if !p.stopped {
		close(p.eventChan)
		p.stopped = true
	}
}

// IsStopped returns whether or not watcher was stopped.
func (p *PollWatcher) IsStopped() bool {
	return p.stopped
}

// ResultChan returns result channel that user can watch for incoming events.
func (p *PollWatcher) ResultChan() <-chan watch.Event {
	p.Lock()
	defer p.Unlock()
	return p.eventChan
}

// NewPollWatcher creates instance of PollWatcher.
func NewPollWatcher() *PollWatcher {
	return &PollWatcher{
		eventChan: make(chan watch.Event),
		stopped:   false,
	}
}
