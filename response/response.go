package response

import (
	"encoding/json"
	"net/http"

	"github.com/auremsinistram/go-errors"
	"github.com/labstack/echo/v5"
)

type Response interface {
	MakeSuccess(data any)
	MakeFailure(code int, message string)
}

type ErrorData interface {
	Status() int
	Message() string
}

func Send[T ErrorData](
	ctx *echo.Context,
	res Response,
	data any,
	err error,
	errRes []byte,
	errData map[int]T,
) error {
	status, bytes := pack(
		res,
		data,
		err,
		errRes,
		errData,
	)

	if e := ctx.JSONBlob(status, bytes); e != nil {
		return errors.Wrap(e, "response - Send - #1")
	}

	return err
}

func pack[T ErrorData](
	res Response,
	data any,
	err error,
	errRes []byte,
	errData map[int]T,
) (int, []byte) {
	if err == nil {
		res.MakeSuccess(data)

		bytes, e := json.Marshal(res)
		if e != nil {
			return http.StatusInternalServerError, errRes
		}

		return http.StatusOK, bytes
	}

	code, ok := errors.GetCode(err)
	if !ok {
		return http.StatusInternalServerError, errRes
	}

	value, ok := errData[code]
	if !ok {
		return http.StatusInternalServerError, errRes
	}

	res.MakeFailure(code, value.Message())

	bytes, e := json.Marshal(res)
	if e != nil {
		return http.StatusInternalServerError, errRes
	}

	return value.Status(), bytes
}
