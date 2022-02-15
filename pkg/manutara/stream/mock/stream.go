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

package mock

import (
	"context"
	"github.com/spiegela/manutara/pkg/manutara/stream/api"
)

type StreamConfig struct {
	Error          error
	InitialContent map[string][]string
}

// NewStream creates a new mock stream client
func NewStream(c StreamConfig) *Stream {
	content := map[string][]string{}
	if c.InitialContent != nil {
		content = c.InitialContent
	}
	return &Stream{content: content, error: c.Error}
}

var _ api.Stream = (*Stream)(nil)

// Stream is a struct containing the state needed to read and write to Mock
// stream
type Stream struct {
	content map[string][]string
	error   error
}

// Peek returns the latest value without removing from the stream
func (m *Stream) Peek(key string) string {
	return m.content[key][len(m.content[key])-1]
}

// OpenWriter creates a Writer to an in memory stream
func (m *Stream) OpenWriter(key string) (api.Writer, error) {
	return &Writer{key: key, stream: m}, nil
}

var _ api.Writer = (*Writer)(nil)

// Writer is a mock stream writer that writes data to string array
type Writer struct {
	key    string
	stream *Stream
}

// Write adds content to an in-memory array to mock a stream and waits for the
// result
func (m *Writer) Write(ctx context.Context, message string) error {
	if m.stream.error != nil {
		return m.stream.error
	}
	content := m.stream.content[m.key]
	m.stream.content[m.key] = append(content, message)
	return nil
}

// WriteAsync adds content to an in-memory array to mock a stream and runs a
// callback upon completion
func (m *Writer) WriteAsync(ctx context.Context, message string, callback func(interface{}, error)) {
	content := m.stream.content[m.key]
	m.stream.content[m.key] = append(content, message)
	callback(message, m.stream.error)
}

// Close closes the connection without removing the data from the mock stream
func (m *Writer) Close() error {
	if m.stream.error != nil {
		return m.stream.error
	}
	return nil
}
