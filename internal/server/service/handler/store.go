package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dontagr/pandora/internal/models"
)

func (h *Handler) Save(c echo.Context) error {
	requestStoreSave, echoError := h.getRequestSave(c)
	if echoError != nil {
		return echoError
	}

	user := h.jwt.GetUser(c)

	fmt.Println(user)
	fmt.Println(requestStoreSave)

	return c.JSON(http.StatusOK, "Пользователь успешно проверен")
}

func (h *Handler) getRequestSave(c echo.Context) (*models.RequestStoreSave, *echo.HTTPError) {
	requestStoreSave := &models.RequestStoreSave{}
	if err := c.Bind(requestStoreSave); err != nil {
		h.log.Errorf("request failed: %v", err)

		return nil, echo.NewHTTPError(http.StatusBadRequest, "Неверный формат запроса")
	}

	return requestStoreSave, nil
}
