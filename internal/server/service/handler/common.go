package handler

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/dontagr/pandora/internal/server/service/customerror"
	"github.com/dontagr/pandora/internal/server/service/secret"
	"github.com/dontagr/pandora/internal/server/service/user"

	"github.com/dontagr/pandora/internal/server/service/jwt"
)

type (
	Handler struct {
		log      *zap.SugaredLogger
		uService *user.Service
		sService *secret.Service
		jwt      *jwt.JWTService
	}
)

func NewHandler(
	uService *user.Service,
	sService *secret.Service,
	log *zap.SugaredLogger,
	jwtService *jwt.JWTService,
) *Handler {
	h := &Handler{
		log:      log,
		uService: uService,
		sService: sService,
		jwt:      jwtService,
	}

	return h
}

func (h *Handler) convertCustomErrorToServerCode(code int) int {
	switch code {
	case customerror.Internal:
		return http.StatusInternalServerError
	case customerror.Unprocessable:
		return http.StatusUnprocessableEntity
	case customerror.Payment:
		return http.StatusPaymentRequired
	case customerror.Unauthorized:
		return http.StatusUnauthorized
	case customerror.Conflict:
		return http.StatusConflict
	default:
		h.log.Warnf("undefind error code [%d]", code)
		return http.StatusInternalServerError
	}
}
