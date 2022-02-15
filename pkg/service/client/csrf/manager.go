package csrf

import (
	"github.com/spiegela/manutara/pkg/service/client/api"
	"k8s.io/client-go/kubernetes"
)

// Implements CSRFTokenManager interface.
type csrfTokenManager struct {
	token  string
	client kubernetes.Interface
}

func (c *csrfTokenManager) init() {
	// TODO:store CSRF key in secret for later retrieval
	c.token = api.GenerateCSRFKey()
}

// Token implements CSRFTokenManager interface.
func (c *csrfTokenManager) Token() string {
	return c.token
}

// NewCSRFTokenManager creates and initializes new instace of csrf token manager.
func NewCSRFTokenManager(client kubernetes.Interface) api.CSRFTokenManager {
	manager := &csrfTokenManager{client: client}
	manager.init()

	return manager
}
