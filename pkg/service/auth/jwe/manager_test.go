package jwe_test

import (
	"testing"

	"github.com/spiegela/manutara/pkg/service/auth/jwe"

	authAPI "github.com/spiegela/manutara/pkg/service/auth/api"
	"k8s.io/client-go/tools/clientcmd/api"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestJwe(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JWE Manager Spec")
}

const Token = "TEST TOKEN"

var _ = Describe("JWE Manager", func() {
	Context("#Generate/#Decrypt", func() {
		var (
			genErr       error
			decErr       error
			tokenManager authAPI.TokenManager
		)

		BeforeEach(func() {
			kh := jwe.NewRSAKeyHolder()
			tokenManager = jwe.NewJWETokenManager(kh)

		})

		Context("when successful", func() {

			var (
				authInfoOut *api.AuthInfo
				token       string
			)

			BeforeEach(func() {
				authInfoIn := api.NewAuthInfo()
				authInfoIn.Token = Token
				token, genErr = tokenManager.Generate(*authInfoIn)
				authInfoOut, decErr = tokenManager.Decrypt(token)
			})

			It("should not error", func() {
				Expect(genErr).ToNot(HaveOccurred())
			})

			It("should generate a token string", func() {
				Expect(token).ToNot(BeEmpty())
			})

			It("should not error when decrypting the token", func() {
				Expect(decErr).ToNot(HaveOccurred())
			})

			It("should have the original info after decryption", func() {
				Expect(authInfoOut.Token).To(Equal(Token))
			})
		})

	})
})
