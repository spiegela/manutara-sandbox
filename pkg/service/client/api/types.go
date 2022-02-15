package api

import (
	"crypto/rand"

	"github.com/sirupsen/logrus"

	"github.com/emicklei/go-restful"
	authAPI "github.com/spiegela/manutara/pkg/service/auth/api"
	"k8s.io/api/authorization/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// Client structures and interfaces are adapted from Kubernetes
// dashboard:
// https://github.com/kubernetes/dashboard/blob/master/src/app/backend/client/api/types.go

// ClientManager is responsible for initializing and creating clients to
// communicate with kubernetes apiserver on demand.
type ClientManager interface {
	Client(req *restful.Request) (kubernetes.Interface, error)

	CanI(req *restful.Request, ssar *v1.SelfSubjectAccessReview) bool

	Config(req *restful.Request) (*rest.Config, error)

	ClientCmdConfig(req *restful.Request) (clientcmd.ClientConfig, error)

	CSRFKey() string

	HasAccess(authInfo api.AuthInfo) error

	SetTokenManager(manager authAPI.TokenManager)
}

// CSRFTokenManager is responsible for generating, reading and updating token stored in a secret.
type CSRFTokenManager interface {
	// Token returns current csrf token used for csrf signing.
	Token() string
}

// GenerateCSRFKey generates random csrf key
func GenerateCSRFKey() string {
	bytes := make([]byte, 256)
	_, err := rand.Read(bytes)
	if err != nil {
		logrus.Fatal("could not generate csrf key")
	}

	return string(bytes)
}
