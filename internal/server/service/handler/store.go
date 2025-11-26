package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dontagr/pandora/internal/models"
	"github.com/dontagr/pandora/internal/server/service/customerror"
)

func (h *Handler) Sync(c echo.Context) error {
	requestSyncNode, echoError := h.getRequestSyncNode(c)
	if echoError != nil {
		if echoError.Err != nil {
			h.log.Infof(echoError.Error())
		}

		return echo.NewHTTPError(h.convertCustomErrorToServerCode(echoError.Code), echoError.Message)
	}

	user := h.jwt.GetUser(c)

	fmt.Println(requestSyncNode)
	fmt.Println(user)

	err := h.sService.SaveSecret(user.ID, requestSyncNode.Kind, requestSyncNode.Label, requestSyncNode.Data, requestSyncNode.Meta, requestSyncNode.Version)
	if err != nil {
		h.log.Infof(err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка")
	}

	return c.JSON(http.StatusOK, "Данные сохранены")

	return c.JSON(http.StatusOK, "ok")
}

func (h *Handler) Load(c echo.Context) error {
	requestStoreLoad, echoError := h.getRequestLoad(c)
	if echoError != nil {
		if echoError.Err != nil {
			h.log.Infof(echoError.Error())
		}

		return echo.NewHTTPError(h.convertCustomErrorToServerCode(echoError.Code), echoError.Message)
	}

	user := h.jwt.GetUser(c)
	secret, err := h.sService.GetSecret(user.ID, requestStoreLoad.Kind, requestStoreLoad.Label)
	if err != nil {
		h.log.Infof(err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка")
	}
	if secret == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Cекрет не найден")
	}

	return c.JSON(http.StatusOK, h.sService.GetResponceStoreLoad(secret))
}

func (h *Handler) List(c echo.Context) error {
	requestStoreList, echoError := h.getRequestList(c)
	if echoError != nil {
		if echoError.Err != nil {
			h.log.Infof(echoError.Error())
		}

		return echo.NewHTTPError(h.convertCustomErrorToServerCode(echoError.Code), echoError.Message)
	}

	user := h.jwt.GetUser(c)

	secret, err := h.sService.GetListSecret(user.ID, requestStoreList.Kind)
	if err != nil {
		h.log.Infof(err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка")
	}

	return c.JSON(http.StatusOK, h.sService.GetResponceStoreList(secret))
}

func (h *Handler) Delete(c echo.Context) error {
	requestStoreLoad, echoError := h.getRequestLoad(c)
	if echoError != nil {
		if echoError.Err != nil {
			h.log.Infof(echoError.Error())
		}

		return echo.NewHTTPError(h.convertCustomErrorToServerCode(echoError.Code), echoError.Message)
	}

	user := h.jwt.GetUser(c)
	secret, err := h.sService.GetSecret(user.ID, requestStoreLoad.Kind, requestStoreLoad.Label)
	if err != nil {
		h.log.Infof(err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка")
	}
	if secret == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Cекрет не найден")
	}

	err = h.sService.DeleteSecret(user.ID, requestStoreLoad.Kind, requestStoreLoad.Label)
	if err != nil {
		h.log.Infof(err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка")
	}

	return c.JSON(http.StatusOK, "ok")
}

func (h *Handler) Save(c echo.Context) error {
	requestStoreSave, echoError := h.getRequestSave(c)
	if echoError != nil {
		if echoError.Err != nil {
			h.log.Infof(echoError.Error())
		}

		return echo.NewHTTPError(h.convertCustomErrorToServerCode(echoError.Code), echoError.Message)
	}

	user := h.jwt.GetUser(c)

	data, err := requestStoreSave.GetData()
	if err != nil {
		h.log.Infof(err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, "проблема с сообщением")
	}

	err = h.sService.SaveSecret(user.ID, requestStoreSave.Kind, requestStoreSave.Label, data, requestStoreSave.Meta, 1)
	if err != nil {
		h.log.Infof(err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, "Внутренняя ошибка")
	}

	return c.JSON(http.StatusOK, "Данные сохранены")
}

func (h *Handler) getRequestSave(c echo.Context) (*models.RequestStoreSave, *customerror.CustomError) {
	requestStoreSave := &models.RequestStoreSave{}
	if err := c.Bind(requestStoreSave); err != nil {
		return nil, customerror.NewCustomError(customerror.BadRequest, "Неверный формат запроса", fmt.Errorf("request failed: %v", err))
	}

	return requestStoreSave, nil
}

func (h *Handler) getRequestList(c echo.Context) (*models.RequestStoreList, *customerror.CustomError) {
	requestStoreList := &models.RequestStoreList{}
	if err := c.Bind(requestStoreList); err != nil {
		return nil, customerror.NewCustomError(customerror.BadRequest, "Неверный формат запроса", fmt.Errorf("request failed: %v", err))
	}

	return requestStoreList, nil
}

func (h *Handler) getRequestLoad(c echo.Context) (*models.RequestStoreLoad, *customerror.CustomError) {
	requestStoreLoad := &models.RequestStoreLoad{}
	if err := c.Bind(requestStoreLoad); err != nil {
		return nil, customerror.NewCustomError(customerror.BadRequest, "Неверный формат запроса", fmt.Errorf("request failed: %v", err))
	}

	return requestStoreLoad, nil
}

func (h *Handler) getRequestSyncNode(c echo.Context) (*models.SyncNode, *customerror.CustomError) {
	syncNode := &models.SyncNode{}
	if err := c.Bind(syncNode); err != nil {
		return nil, customerror.NewCustomError(customerror.BadRequest, "Неверный формат запроса", fmt.Errorf("request failed: %v", err))
	}

	return syncNode, nil
}
