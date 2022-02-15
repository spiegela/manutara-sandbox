package auth

// Authentication manager is adapted from Kubernetes dashboard:
// https://github.com/kubernetes/dashboard/blob/master/src/app/backend/auth/manager.go

import (
	"net/http"

	"github.com/spiegela/manutara/pkg/service/auth/api"
	clientAPI "github.com/spiegela/manutara/pkg/service/client/api"
	"k8s.io/apimachinery/pkg/api/errors"
	cmdAPI "k8s.io/client-go/tools/clientcmd/api"
)

// Implements AuthManager interface
type authManager struct {
	tokenManager  api.TokenManager
	clientManager clientAPI.ClientManager
}

// NewAuthManager generates a new AuthManager
func NewAuthManager(c clientAPI.ClientManager, t api.TokenManager) api.AuthManager {
	return &authManager{
		tokenManager:  t,
		clientManager: c,
	}
}

func (a authManager) Decrypt(jweToken string) (*cmdAPI.AuthInfo, error) {
	return a.tokenManager.Decrypt(jweToken)
}

func (a authManager) Refresh(jweToken string) (string, error) {
	return a.tokenManager.Refresh(jweToken)
}

func (a authManager) Login(spec *api.LoginSpec) (*api.AuthResponse, error) {
	authenticator, err := a.getAuthenticator(spec)
	if err != nil {
		return nil, err
	}

	authInfo, err := authenticator.GetAuthInfo()
	if err != nil {
		return nil, err
	}
	err = a.healthCheck(authInfo)

	nonCriticalError, criticalError := handleError(err)
	if criticalError != nil {
		return &api.AuthResponse{}, err
	} else if nonCriticalError.Code != 0 {
		return &api.AuthResponse{Error: nonCriticalError}, nil
	}
	token, err := a.tokenManager.Generate(authInfo)
	if err != nil {
		return nil, err
	}
	return &api.AuthResponse{JWEToken: token}, nil
}

func handleError(err error) (api.AuthError, error) {
	var nonCriticalError api.AuthError
	status, ok := err.(*errors.StatusError)
	if !ok {
		return nonCriticalError, err
	}
	if status.ErrStatus.Code == http.StatusForbidden ||
		status.ErrStatus.Code == http.StatusUnauthorized {
		return api.AuthError{
			Code:  int(status.ErrStatus.Code),
			Error: err,
		}, nil
	}
	return nonCriticalError, err
}

// Returns authenticator based on provided LoginSpec.
func (a authManager) getAuthenticator(spec *api.LoginSpec) (api.Authenticator, error) {
	return NewTokenAuthenticator(spec), nil
}

func (a authManager) healthCheck(authInfo cmdAPI.AuthInfo) error {
	return a.clientManager.HasAccess(authInfo)
}

var _ api.AuthManager = (*authManager)(nil)
