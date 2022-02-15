package auth

import (
	authAPI "github.com/spiegela/manutara/pkg/service/auth/api"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Implements Authenticator interface
type tokenAuthenticator struct {
	token string
}

// GetAuthInfo implements Authenticator interface. See Authenticator for more
// information.
func (a *tokenAuthenticator) GetAuthInfo() (api.AuthInfo, error) {
	return api.AuthInfo{
		Token: a.token,
	}, nil
}

// NewTokenAuthenticator returns Authenticator based on LoginSpec.
func NewTokenAuthenticator(spec *authAPI.LoginSpec) authAPI.Authenticator {
	return &tokenAuthenticator{
		token: spec.Token,
	}
}
