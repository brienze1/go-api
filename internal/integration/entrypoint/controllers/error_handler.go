package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/brienze1/notes-api/internal/integration/adapters"
	"log/slog"

	"github.com/brienze1/notes-api/internal/domain/exceptions"
)

type errorHandler struct {
}

func (e *errorHandler) HandlePanic(ctx context.Context, recovered any) (response []byte, statusCode int) {
	if recovered != nil {
		err := exceptions.NewInternalServerError(fmt.Sprintf("panic: %v", recovered))
		return e.HandleError(ctx, err)
	}
	return
}

func (e *errorHandler) HandleError(ctx context.Context, err error) (response []byte, statusCode int) {
	var errParsed *exceptions.ErrorType
	if !errors.As(err, &errParsed) {
		errParsed = exceptions.NewInternalServerError(err.Error())
	}

	slog.ErrorContext(ctx, "errorHandler.HandleError", slog.String("errorDetails", string(errParsed.JSON())))

	/*
	   if errParsed.StatusCode == http.StatusInternalServerError {
	       slack notification?
	   }
	*/

	return errParsed.JSON(), errParsed.StatusCode
}

func NewErrorHandler() adapters.ErrorHandler {
	return &errorHandler{}
}
