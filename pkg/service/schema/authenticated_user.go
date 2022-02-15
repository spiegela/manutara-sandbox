package schema

import (
	"encoding/json"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
)

// AuthenticatedUser contains the fields for the currently authenticated user
type AuthenticatedUser struct {
	Protected string `json:"protected"`

	AdditionalAuthData string `json:"aad"`

	EncryptedKey string `json:"encrypted_key"`

	InitializationVector string `json:"iv"`

	CipherText string `json:"ciphertext"`

	Tag string `json:"tag"`
}

// AuthenticatedUserType is a GraphQL Type defined an authenticated user
func AuthenticatedUserType(nodeDef *relay.NodeDefinitions) *graphql.Object {
	namespaceType := NamespaceType(nodeDef)
	return graphql.NewObject(graphql.ObjectConfig{
		Name:        "AuthenticatedUser",
		Description: "A Kubernetes user that has been authenticated to the API",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("AuthenticatedUser", nil),
			"protected": &graphql.Field{
				Name: "Protected",
				Type: graphql.String,
			},
			"aad": &graphql.Field{
				Name: "AdditionalAuthData",
				Type: graphql.String,
			},
			"encryptedKey": &graphql.Field{
				Name: "EncryptedKey",
				Type: graphql.String,
			},
			"iv": &graphql.Field{
				Name: "InitializationVector",
				Type: graphql.String,
			},
			"ciphertext": &graphql.Field{
				Name: "CipherText",
				Type: graphql.String,
			},
			"tag": &graphql.Field{
				Name: "Tag",
				Type: graphql.String,
			},
			"namespaces": &graphql.Field{
				Type: connectionDefinitionFor("Namespace", namespaceType).ConnectionType,
				Args: relay.ConnectionArgs,
				Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
					var emptyList = []interface{}{}
					args := relay.NewConnectionArguments(p.Args)
					r, ok := p.Info.RootValue.(map[string]interface{})
					if !ok {
						return emptyList, fmt.Errorf("no root value assigned to the request")
					}
					c, ok := r["client"].(client.Client)
					if !ok {
						return emptyList, fmt.Errorf("no authenticated Kubernetes client for request")
					}
					resp, err := GetNamespaces(c, p)
					if err != nil {
						return emptyList, err
					}
					return relay.ConnectionFromArray(resp, args), nil
				},
				Description: "Returns a list of Kubernetes namespaces",
			},
		},
		Interfaces: []*graphql.Interface{
			nodeDef.NodeInterface,
		},
	})
}

// GetAuthenticatedUser returns the current user from the http context
func GetAuthenticatedUser(p graphql.ResolveParams) (interface{}, error) {
	root, ok := p.Source.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("no root object containing a token found")
	}
	token, ok := root["token"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid token found")
	}
	var authenticatedUser = &AuthenticatedUser{}
	err := json.Unmarshal([]byte(token), authenticatedUser)
	if err != nil {
		return nil, fmt.Errorf("unable to deserialize token: %s", err)
	}
	return authenticatedUser, nil
}
