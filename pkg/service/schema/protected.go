package schema

import (
	"fmt"

	"github.com/graphql-go/relay"

	"github.com/spiegela/manutara/pkg/service/auth/api"

	"github.com/graphql-go/graphql"
)

func protectedResolver(resolver func(p graphql.ResolveParams) (interface{}, error)) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		root, ok := p.Source.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("source is not a schema root object")
		}
		if root["errors"] != nil {
			errorList, ok := root["errors"].([]api.AuthError)
			if !ok {
				return nil, fmt.Errorf("unable to check for authentication errors")
			}
			if len(errorList) > 0 {
				return nil, fmt.Errorf("authentication errors found")
			}
		}
		if root["token"] == nil {
			return nil, fmt.Errorf("source does not contain a JWE token")
		}
		_, ok = root["token"].(string)
		if !ok {
			return nil, fmt.Errorf("source does not contain a JWE token")
		}
		// TODO:should we validate the token and expiry here?
		return resolver(p)
	}
}

func protectedConnectionResolver(resolver func(p graphql.ResolveParams) ([]interface{}, error)) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		_, err := protectedResolver(func(p graphql.ResolveParams) (i interface{}, e error) {
			return nil, e
		})(p)
		if err != nil {
			return nil, err
		}
		resp, err := resolver(p)
		if err != nil {
			return nil, err
		}
		args := relay.NewConnectionArguments(p.Args)
		return relay.ConnectionFromArray(resp, args), nil
	}
}
