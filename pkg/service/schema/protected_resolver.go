package schema

import (
	"fmt"

	"github.com/spiegela/manutara/pkg/service/auth/api"

	"github.com/graphql-go/graphql"
)

func protectedResovler(resolver func(p graphql.ResolveParams) (interface{}, error)) graphql.FieldResolveFn {
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
		// TODO should we validate the token and expiry here?
		return resolver(p)
	}
}
