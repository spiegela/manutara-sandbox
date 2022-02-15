package api

// Authentication structures and interfaces are adapted from Kubernetes
// dashboard:
// https://github.com/kubernetes/dashboard/blob/master/src/app/backend/auth/api/types.go

import (
	"time"

	"k8s.io/client-go/tools/clientcmd/api"
)

// DefaultTokenTTL is the default number of seconds that an authentication token
// is valid.
const DefaultTokenTTL = 900

// Authenticator represents authentication methods supported by GraphQL server.
// Currently supported types are:
//    - Token based - Any bearer token accepted by apiserver
type Authenticator interface {
	// GetAuthInfo returns filled AuthInfo structure that can be used for K8S
	// api client creation.
	GetAuthInfo() (api.AuthInfo, error)
}

// LoginSpec is extracted from request coming from the frontend during login
// request. It contains all the information required to authenticate user.
type LoginSpec struct {
	// Token is the bearer token for authentication to the kubernetes cluster.
	Token string `json:"token"`
}

// AuthError is returned if the backend response for a login/refresh receives
// a non-critical error such as 401 or 403
type AuthError struct {
	Code  int
	Error error `json:"message"`
}

// AuthResponse is returned from our backend as a response for login/refresh
// requests. It contains generated JWToken and a list of non-critical errors
// such as 'Failed authentication'.
type AuthResponse struct {
	// JWEToken is a token generated during login request that contains AuthInfo
	// data in the payload.
	JWEToken string `json:"jweToken"`

	// Error is a non-critical errors that happened during login request.
	Error AuthError `json:"error"`
}

// AuthManager is used for user authentication management.
type AuthManager interface {
	// Refresh takes valid token that hasn't expired yet and returns a new one
	// with expiration time set to TokenTTL. In case provided token has expired,
	// token expiration error is returned.
	Refresh(string) (string, error)

	// Login authenticates user based on provided LoginSpec and returns
	// AuthResponse. AuthResponse contains generated token and list of
	// non-critical errors such as 'Failed authentication'.
	Login(*LoginSpec) (*AuthResponse, error)

	// Decrypt returns the decrypted authentication information from the JWE
	Decrypt(string) (*api.AuthInfo, error)
}

// TokenRefreshSpec contains token that is required by token refresh operation.
type TokenRefreshSpec struct {
	// JWToken is a token generated during login request that contains AuthInfo
	// data in the payload.
	JWToken string `json:"jwToken"`
}

// TokenManager is responsible for generating and decrypting tokens used for
// authorization. Authorization is handled by K8S apiserver. Token contains
// AuthInfo structure used to create K8S api client.
type TokenManager interface {
	// Generate secure token based on AuthInfo structure and save it tokens'
	// payload.
	Generate(api.AuthInfo) (string, error)

	// Decrypt generated token and return AuthInfo structure that will be used
	// for K8S api client creation.
	Decrypt(string) (*api.AuthInfo, error)

	// Refresh returns refreshed token based on provided token. In case provided
	// token has expired, token expiration error is returned.
	Refresh(string) (string, error)

	// SetTokenTTL sets expiration time (in seconds) of generated tokens.
	SetTokenTTL(time.Duration)
}
