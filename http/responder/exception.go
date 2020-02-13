package responder

import (
	"net/http"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	gateway "github.com/grpc-ecosystem/grpc-gateway/runtime"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const errorDigit = 1000

type errorResp struct {
	code    codes.Code
	message string
	headers http.Header
}

func (e *errorResp) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	for k, v := range e.headers {
		for _, val := range v {
			rw.Header().Add(k, val)
		}
	}
	codeInt := int(e.code)
	if codeInt < 0 {
		rw.WriteHeader(http.StatusInternalServerError)
	} else if codeInt < errorDigit {
		rw.WriteHeader(gateway.HTTPStatusFromCode(e.code))
	} else {
		rw.WriteHeader(gateway.HTTPStatusFromCode(codes.Code(codeInt / errorDigit)))
	}

	if err := producer.Produce(rw, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    codeInt,
			"message": e.message,
		},
	}); err != nil {
		panic(err)
	}
}

// NotFoundError the error response when the response is not implemented
func NotFoundError(message string) middleware.Responder {
	return &errorResp{http.StatusNotFound, message, make(http.Header)}
}

// InternalServerError the error response when the response is internal error
func InternalServerError(message string) middleware.Responder {
	return &errorResp{http.StatusInternalServerError, message, make(http.Header)}
}

// BadRequestError the error response when the response is bad request
func BadRequestError(message string) middleware.Responder {
	return &errorResp{http.StatusBadRequest, message, make(http.Header)}
}

// UnauthorisedError the error response when the response is Unauthorised
func UnauthorisedError(message string) middleware.Responder {
	return &errorResp{http.StatusUnauthorized, message, make(http.Header)}
}

// PaymentRequiredError the error response when the response is PaymentRequiredError
func PaymentRequiredError(message string) middleware.Responder {
	return &errorResp{http.StatusPaymentRequired, message, make(http.Header)}
}

// ForbiddenError the error response when the response is Forbidden
func ForbiddenError(message string) middleware.Responder {
	return &errorResp{http.StatusForbidden, message, make(http.Header)}
}

func PreConditionFailedError(message string) middleware.Responder {
	return &errorResp{http.StatusPreconditionFailed, message, make(http.Header)}
}

func MapGrpcError(err error) middleware.Responder {
	s, ok := status.FromError(err)
	if !ok {
		return InternalServerError("")
	}
	if s.Code() == codes.OK {
		log.Errorf("Handling Ok as an error! %v", s.Message())
		return InternalServerError("Server error.")
	}
	return &errorResp{s.Code(), s.Message(), make(http.Header)}
}

func MakeErrResp(code codes.Code, message string) middleware.Responder {
	return &errorResp{code, message, make(http.Header)}
}
