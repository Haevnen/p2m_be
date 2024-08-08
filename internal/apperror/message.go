package apperror

import (
	"context"
	"fmt"
	"net/http"

	api "github.com/Haevnen/p2m_be/internal/app/p2m_api/gen/api"
)

type definedErrorDetail struct {
	httpStatus int
	resType    string
	errCode    string
	msg        string
}

var (
	errMessageMap = map[error]definedErrorDetail{

		context.Canceled:  {httpStatus: http.StatusBadRequest, resType: string(api.ValidationFailed), errCode: errCode4999, msg: "The operation was canceled."},
		ErrInternalServer: {httpStatus: http.StatusInternalServerError, resType: string(api.InternalError), errCode: errCodeInternalServerError, msg: "An Unexpected Error has occurred."},
	}
)

func (d definedErrorDetail) detailsJA(params ...any) []string {
	return []string{fmt.Sprintf(d.msg, params...)}
}

const (
	errCode4999                = "ERR_4999"
	errCodeInternalServerError = "ERR_9999"
)
