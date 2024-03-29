package authorization

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/golovpeter/vk_intership_test_task/internal/common"
	"github.com/golovpeter/vk_intership_test_task/internal/repository/users"
	"github.com/sirupsen/logrus"
)

var skipURL = map[string]bool{
	"/v1/user/register": true,
	"/v1/user/login":    true,
	"/favicon.ico":      true,
	"/swagger.yaml":     true,
	"/docs":             true,
}

func AuthorizationMiddleware(
	logger *logrus.Logger,
	enforcer *casbin.Enforcer,
	repository users.Repository,
	next http.Handler,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if skipURL[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			logger.Error(common.InvalidCredentialsError)
			common.MakeErrorResponse(w, http.StatusUnauthorized, common.InvalidAuthHeader)
			return
		}

		spltToken := strings.Split(accessToken, " ")

		if spltToken[0] != "Bearer" {
			logger.Error(common.InvalidCredentialsError)
			common.MakeErrorResponse(w, http.StatusUnauthorized, common.InvalidAuthHeader)
			return
		}

		claims, err := common.GetTokenClaims(spltToken[1])
		if err != nil {
			logger.WithError(err).Error(err.Error())
			common.MakeErrorResponse(w, http.StatusUnauthorized, err)
			return
		}

		userRole, err := repository.GetUserRole(r.Context(), claims["Username"].(string))
		if err != nil {
			logger.Error(err)
			common.MakeErrorResponse(w, http.StatusUnauthorized, err)
			return
		}

		allow, err := enforcer.Enforce(userRole, r.URL.Path, r.Method)
		if err != nil {
			logger.Error(err)
			common.MakeErrorResponse(w, http.StatusUnauthorized, err)
			return
		}

		if !allow {
			logger.Error(common.AccessDeniedError)
			common.MakeErrorResponse(w, http.StatusUnauthorized, common.AccessDeniedError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
