package apierrors

import (
	"errors"
	"fmt"
)

type Code int

type APIerr struct {
	Code     Code   `json:"code"`
	Text     string `json:"text"`
	Info     string `json:"info"`
	Location string `json:"location"`
	Err      error  `json:"error"`
}

const (
	UnknownErrorCode             Code = 1000
	InternalErrorCode            Code = 1001
	InvalidInputErrorCode        Code = 1002
	UnsupportedMediaTypeCode     Code = 1003
	NotFoundErrorCode            Code = 1004
	ClientClosedRequestErrorCode Code = 1005
	EntityTooLargeErrorCode      Code = 1006
	UnauthorizedErrorCode        Code = 1007
	ForbiddenErrorCode           Code = 1008
)

func (ae *APIerr) WithLocation(location string) *APIerr {
	ae.Location = location
	return ae
}

func (ae *APIerr) AppendLocation(location string) *APIerr {
	if ae.Location == "" {
		ae.Location = location
	} else {
		ae.Location = fmt.Sprintf(ae.Location, "--->", location)
	}
	return ae
}

func (ae *APIerr) WithInfo(info string) *APIerr {
	ae.Info = info
	return ae
}

func (ae *APIerr) WithText(text string) *APIerr {
	ae.Text = text
	return ae
}

func (ae *APIerr) ToString() string {
	err := ""
	if ae.Err != nil {
		err = fmt.Sprint(" Error raw: ", ae.Err.Error(), ";")
	}

	info := ""
	if ae.Info != "" {
		info = fmt.Sprint(" Info: ", ae.Info, ";")
	}

	location := ""
	if ae.Location != "" {
		location = fmt.Sprint(" Location: ", ae.Location, ";")
	}

	return fmt.Sprintf("Error: %d: %s;%s%s%s", ae.Code, ae.Text, info, err, location)
}

func (ae *APIerr) ToErr(errsToJoin ...error) error {
	resultErr := ae.Err
	if len(errsToJoin) != 0 {
		resultErr = errors.Join(errsToJoin...)
	}
	return resultErr
}

func NewUnknownError(err error) *APIerr {
	return &APIerr{
		Code: UnknownErrorCode,
		Text: "Internal error",
		Err:  err,
	}
}

func NewInternalError(info string, err ...error) *APIerr {
	return &APIerr{
		Code: InternalErrorCode,
		Text: "Internal error",
		Info: info,
		Err:  getFirstErrOrNil(err...),
	}
}

func NewInvalidInputError(message string, err ...error) *APIerr {
	cause := getFirstErrOrNil(err...)
	apierr := &APIerr{
		Code: InvalidInputErrorCode,
		Text: message,
		Err:  cause,
	}

	if cause != nil {
		return apierr.WithText(cause.Error())
	}

	return apierr
}

func NewNotFoundError(info string, err ...error) *APIerr {
	return &APIerr{
		Code: NotFoundErrorCode,
		Text: "Resource not found",
		Info: info,
		Err:  getFirstErrOrNil(err...),
	}
}

func NewClientClosedRequestError() *APIerr {
	return &APIerr{
		Code: ClientClosedRequestErrorCode,
		Text: "User closed connection",
	}
}

func NewTooLargeError(info string) *APIerr {
	return &APIerr{
		Code: EntityTooLargeErrorCode,
		Text: "too large",
		Info: info,
	}
}

func NewUnauthorizedError(err ...error) *APIerr {
	return &APIerr{
		Code: UnauthorizedErrorCode,
		Text: "Unauthorized request received",
		Err:  getFirstErrOrNil(err...),
	}
}

func NewForbiddenError(err ...error) *APIerr {
	return &APIerr{
		Code: ForbiddenErrorCode,
		Text: "Request forbidden",
		Err:  getFirstErrOrNil(err...),
	}
}

func getFirstErrOrNil(err ...error) error {
	if len(err) == 0 {
		return nil
	}

	return err[0]
}

func ToAPIError(err error) *APIerr {
	var apierr *APIerr
	if errors.As(err, &apierr) {
		return apierr
	}

	return NewUnknownError(err)
}
