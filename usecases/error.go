package usecases

import (
	"errors"
	"net/http"
	"twitter-bff/domain/models"
)

func ErrorHandler(err error) (int, string) {
	if errors.Is(err, models.ErrInvalidArgument) {
		return http.StatusUnprocessableEntity, err.Error()
	}

	if errors.Is(err, models.ErrNotFound) {
		return http.StatusNotFound, err.Error()
	}

	if errors.Is(err, models.ErrInternal) {
		return http.StatusInternalServerError, err.Error()
	}

	return http.StatusInternalServerError, err.Error()
}
