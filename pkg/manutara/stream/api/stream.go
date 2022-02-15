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

package api

import (
	"context"
)

// Stream is an interface for a stream provider, which is the shared commit log
// for the CQRS system
type Stream interface {
	// OpenWriter opens a stream for writing. It takes a routing-key or filename
	// equivalent as an argument, and returns a struct implementing the Writer
	// interface
	OpenWriter(key string) (Writer, error)
}

// Writer is an interface for adding messages to the shared commit log
type Writer interface {
	// Write persists a message synchronously to the stream
	Write(ctx context.Context, message string) error

	// WriteAsync persists a message asynchronously to the stream, if the the
	// stream supports it. Otherwise, this is an alias for Write
	WriteAsync(ctx context.Context, message string, callback func(interface{}, error))

	// Close closes a stream to any further writes for this connection
	Close() error
}
