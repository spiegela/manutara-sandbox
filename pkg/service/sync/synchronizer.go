package secret

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/runtime"

	"k8s.io/apimachinery/pkg/types"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"

	"k8s.io/apimachinery/pkg/util/wait"

	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/api/core/v1"
)

type ActionHandlerFunction func(runtime.Object)

type synchronizer struct {
	namespace      string
	errChan        chan error
	watcher        *PollWatcher
	client         k8sclient.Client
	mux            sync.Mutex
	actionHandlers map[watch.EventType][]ActionHandlerFunction
}

func NewSynchronizer(client k8sclient.Client, namespace string) Synchronizer {
	return &synchronizer{
		client:    client,
		namespace: namespace,
	}
}

// Time interval between which secret should be resynchronized.
const secretSyncPeriod = 5 * time.Minute

func (s *synchronizer) Name() string {
	//return fmt.Sprintf("%s-%s", s.name, s.namespace)
	return fmt.Sprintf("%s", s.namespace)
}

// Create a watch in a new channel subscribing to secret events
func (s *synchronizer) Start() {
	s.errChan = make(chan error)
	watcher := s.watch(s.namespace, s.name, secretSyncPeriod)
	go func() {
		log.Printf("Starting secret synchronizer for %s in namespace %s", s.name, s.namespace)
		defer watcher.Stop()
		defer close(s.errChan)
		for {
			select {
			case ev, ok := <-watcher.ResultChan():
				if !ok {
					s.errChan <- fmt.Errorf("%s watch ended with timeout", s.Name())
					return
				}
				if err := s.handleEvent(ev); err != nil {
					s.errChan <- err
					return
				}
			}
		}
	}()
}

// Error returns error channel. Any error that happens during running
// synchronizer will be send to this channel.
func (s *synchronizer) Error() chan error {
	return s.errChan
}

func (s *synchronizer) watch(namespace string, name string, interval time.Duration) watch.Interface {
	stopCh := make(chan struct{})
	go wait.Until(func() {
		if s.watcher.IsStopped() {
			close(stopCh)
			return
		}

		s.watcher.eventChan <- s.getSecretEvent()
	}, interval, stopCh)
	return s.watcher
}

func (s *synchronizer) handleEvent(event watch.Event) error {
	for _, handler := range s.actionHandlers[event.Type] {
		handler(event.Object)
	}
	switch event.Type {
	case watch.Added:
		secret, ok := event.Object.(*v1.Secret)
		if !ok {
			return errors.New(fmt.Sprintf("Expected secret got %s", reflect.TypeOf(event.Object)))
		}
		s.update(*secret)
	case watch.Modified:
		secret, ok := event.Object.(*v1.Secret)
		if !ok {
			return errors.New(fmt.Sprintf("Expected secret got %s", reflect.TypeOf(event.Object)))
		}
		s.update(*secret)
	case watch.Deleted:
		s.mux.Lock()
		s.secret = nil
		s.mux.Unlock()
	case watch.Error:
		return &k8serrors.UnexpectedObjectError{Object: event.Object}
	}
	return nil
}

func (s *synchronizer) update(secret v1.Secret) {
	if reflect.DeepEqual(s.secret, &secret) {
		// Skip update if existing object is the same as new one
		return
	}
	s.mux.Lock()
	s.secret = &secret
	s.mux.Unlock()
}

func (s *synchronizer) getSecretEvent() (event watch.Event) {
	secret := &v1.Secret{}
	event = watch.Event{
		Object: secret,
		Type:   watch.Added,
	}
	err := s.client.Get(context.TODO(), types.NamespacedName{
		Namespace: s.namespace,
		Name:      s.name,
	}, secret)
	if err != nil {
		event.Type = watch.Error
	}
	// In case it was never created we can still mark it as deleted and let
	// secret be recreated.
	if k8serrors.IsNotFound(err) {
		event.Type = watch.Deleted
	}
	return
}

func (s *synchronizer) Refresh() {
	s.mux.Lock()
	defer s.mux.Unlock()

	secret := v1.Secret{}
	err := s.client.Get(context.TODO(), types.NamespacedName{
		Namespace: s.namespace,
		Name:      s.name,
	}, &secret)
	if err != nil {
		log.Printf("Secret synchronizer %s failed to refresh secret", s.Name())
		return
	}

	s.secret = &secret
}
