package stream

import (
	"github.com/spiegela/manutara/pkg/manutara/stream/api"
	"github.com/spiegela/manutara/pkg/manutara/stream/mock"
	"github.com/spiegela/manutara/pkg/manutara/stream/nats"
	"log"
	"net/url"
)

// NewStream will parse a URI string and return a Stream
func NewStream(uri string) (api.Stream, string) {
	url, err := url.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}
	switch url.Scheme {
	case "mock":
		s := mock.NewStream(mock.StreamConfig{})
		return s, url.Path
	case "nats":
		s, err := nats.NewStream(nats.StreamConfig{
			ClusterName: url.Host,
			URI:         uri,
		})
		if err != nil {
			log.Fatal(err)
		}
		return s, url.Path
	default:
		log.Fatalf("unsupported stream scheme: %s", url.Scheme)
	}
	return nil, ""
}
