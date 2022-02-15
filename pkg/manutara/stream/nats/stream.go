/*
Copyright 2019 Aaron Spiegel.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nats

import (
	"context"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/rs/xid"
	"github.com/spiegela/manutara/pkg/manutara/stream/api"
	"log"
	"sync"
)

// Stream implements the stream.Stream interface
type Stream struct {
	// ClusterName is the NATS streaming cluster name
	ClusterName string

	// URI is the connection string representing the NATS client connection
	URL stan.Option
}

// StreamConfig provides the configuration parameters needed for a nats stream
type StreamConfig struct {
	// ClusterName is the NATS streaming cluster name
	ClusterName string

	// URI is the connection string representing the NATS client connection
	URI string
}

// NewStream creates a new nats-streaming client stream
func NewStream(config StreamConfig) (*Stream, error) {
	return &Stream{
		URL: stan.NatsURL(config.URI),
	}, nil
}

func (s *Stream) OpenWriter(subject string) (api.Writer, error) {
	id := xid.New()
	conn, err := stan.Connect(id.String(), s.ClusterName, s.URL)
	if err != nil {
		return nil, err
	}
	var mut sync.Mutex
	return Writer{conn, subject, mut}, nil
}

var _ api.Stream = (*Stream)(nil)

type Writer struct {
	conn    stan.Conn
	subject string
	mut     sync.Mutex
}

func (w Writer) Write(ctx context.Context, message string) error {
	err := w.conn.Publish(w.subject, []byte(message))
	if err != nil {
		return err
	}
	return nil
}

func (w Writer) WriteAsync(_ context.Context, message string, callback func(interface{}, error)) {
	w.mut.Lock()
	guid, err := w.conn.PublishAsync(w.subject, []byte(message), func(lguid string, err error) {
		callback(lguid, err)
	})
	if err != nil {
		log.Fatalf("Error during async publish: %v\n", err)
	}
	w.mut.Unlock()
	if guid == "" {
		log.Fatal("Expected non-empty guid to be returned.")
	}
}

func (w Writer) Close() error {
	if err := w.conn.Close(); err != nil {
		return err
	}
	return nil
}

var _ api.Writer = (*Writer)(nil)
