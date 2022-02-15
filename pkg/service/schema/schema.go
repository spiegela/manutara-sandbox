package schema

import (
	"errors"

	"github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
)

// GenerateSchema returns GraphQL schema for this application
func GenerateSchema(cfg *rest.Config) *graphql.Schema {
	schema, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery(nodeDef),
	})
	if err != nil {
		logrus.Fatal(err)
	}
	return &schema
}
