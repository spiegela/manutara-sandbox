package login

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/spiegela/manutara/pkg/service/auth"
	authAPI "github.com/spiegela/manutara/pkg/service/auth/api"
	clientAPI "github.com/spiegela/manutara/pkg/service/client/api"
	"golang.org/x/text/language"
)

const (
	// InvalidAuthToken is the i18n message name for an error that an invalid
	// token was provided
	InvalidAuthToken = "invalid_auth_token"

	// LoginSuccessful is the i18n message name for a successful login
	LoginSuccess = "login_success"

	// NoAuthToken is the i18n message name for an error that no authentication
	// token was provided
	NoAuthToken = "no_auth_token"

	// RequestError is the i18n message name for an error that the server was
	// unable to process the request due to an invalid request
	RequestError = "request_error"

	// ServerError is the i18n message name for an error that the server was
	// unable to process the request due to an internal error
	ServerError = "server_error"
)

type Handler struct {
	AuthManager authAPI.AuthManager
	i18n        *i18n.Bundle
}

func NewHandler(client clientAPI.ClientManager, token authAPI.TokenManager) *Handler {
	authManager := auth.NewAuthManager(client, token)
	bundle := &i18n.Bundle{DefaultLanguage: language.English}
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		logrus.Fatal("unable to get executable location")
	}
	for _, lang := range []string{"en"} {
		filepath := path.Join(path.Dir(filename), fmt.Sprintf("../../i18n/%s.toml", lang))
		bundle.MustLoadMessageFile(filepath)
	}
	return &Handler{
		AuthManager: authManager,
		i18n:        bundle,
	}
}

// ServeHTTP provides an entrypoint into executing login requests.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{}
	lang := r.Header.Get("Accept-Language")
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) == 0 {
		resp := createResponse("error",
			h.localize(lang, NoAuthToken, data), "")
		h.writeResp(w, http.StatusUnauthorized, resp)
		return
	}
	splitToken := strings.Split(bearerToken, " ")
	if len(splitToken) < 2 || splitToken[0] != "Bearer" || len(splitToken[1]) == 0 {
		resp := createResponse("error",
			h.localize(lang, InvalidAuthToken, data), "")
		h.writeResp(w, http.StatusUnauthorized, resp)
		return
	}
	loginResp, err := h.AuthManager.Login(&authAPI.LoginSpec{Token: splitToken[1]})
	if err != nil {
		data["Request"] = "login"
		data["Error"] = err.Error()
		resp := createResponse("error",
			h.localize(lang, ServerError, data), "")
		h.writeResp(w, http.StatusInternalServerError, resp)
		return
	} else if loginResp == nil {
		data["Request"] = "login"
		data["Error"] = "Login response from authentication service was empty"
		resp := createResponse("error",
			h.localize(lang, ServerError, data), "")
		h.writeResp(w, http.StatusInternalServerError, resp)
		return
	} else if loginResp.Error.Code != 0 {
		data["Request"] = "login"
		data["Error"] = loginResp.Error.Error
		resp := createResponse("error",
			h.localize(lang, ServerError, data), "")
		h.writeResp(w, http.StatusInternalServerError, resp)
		return
	}
	if len(loginResp.JWEToken) == 0 {
		data["Request"] = "login"
		data["Error"] = "Login service failed to return a valid token"
		resp := createResponse("error",
			h.localize(lang, ServerError, data), "")
		h.writeResp(w, http.StatusInternalServerError, resp)
		return
	}
	resp := createResponse("info",
		h.localize(lang, LoginSuccess, data), loginResp.JWEToken)
	h.writeResp(w, http.StatusOK, resp)
}

func (h *Handler) writeResp(w http.ResponseWriter, statusCode int, resp []byte) {
	w.WriteHeader(statusCode)
	_, err := w.Write(resp)
	if err != nil {
		logrus.Fatal(err)
	}
}

func createResponse(level string, message string, token string) []byte {
	resp, err := json.Marshal(map[string]string{
		"alertType":    level,
		"alertMessage": message,
		"token":        token,
	})
	if err != nil {
		logrus.Fatal(err)
	}
	return resp
}

func (h *Handler) localize(language string, message string, data map[string]interface{}) string {
	return i18n.NewLocalizer(h.i18n, language).MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID: message,
		},
		TemplateData: data,
	})
}
