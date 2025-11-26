package routing

import (
	"fmt"

	"github.com/dontagr/pandora/internal/server/httpserver/validator"
	"github.com/dontagr/pandora/internal/server/service/handler"

	"github.com/dontagr/pandora/internal/server/httpserver"
	"github.com/dontagr/pandora/internal/server/service/jwt"
)

func InitRouting(
	server *httpserver.HTTPServer,
	jwt *jwt.JWTService,
	handler *handler.Handler,
) error {
	var err error
	jwtConfig := jwt.GetJWTEchoConfig()

	server.Master.Validator, err = validator.NewCustomValidator()
	if err != nil {
		return fmt.Errorf("failed create validator %v", err)
	}

	g := server.Master.Group("/api/user")
	g.POST("/register", handler.SignUp)
	g.POST("/login", handler.SignIn)

	s := server.Master.Group("/api/store")
	s.POST("/save", handler.Save, jwt.GetMiddleware(jwtConfig))
	s.POST("/load", handler.Load, jwt.GetMiddleware(jwtConfig))
	s.POST("/delete", handler.Delete, jwt.GetMiddleware(jwtConfig))
	s.POST("/list", handler.List, jwt.GetMiddleware(jwtConfig))

	return nil
}
