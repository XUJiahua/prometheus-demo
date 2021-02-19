package card

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"prometheus-demo/chan/cybersource"
)

type handlerFunc func(*Request) (*Response, error)

type Service interface {
	// register middleware yourself
	WrapMiddlewares(handler handlerFunc) handlerFunc
	Auth(request *Request) (*Response, error)
	Capture(request *Request) (*Response, error)
	Refund(request *Request) (*Response, error)
}

type middlewareFunc func(handlerFunc) handlerFunc

type ServiceImpl struct {
	middlewares    []middlewareFunc
	OpCounter      *prometheus.CounterVec
	ChannelService cybersource.Service
}

func (s *ServiceImpl) RegisterMiddleware() Service {
	s.middlewares = append(s.middlewares,
		s.internalIDGenMiddleware,
		s.metricRecorderMiddleware)
	return s
}

func (s ServiceImpl) WrapMiddlewares(handler handlerFunc) handlerFunc {
	for _, middleware := range s.middlewares {
		handler = middleware(handler)
	}
	return handler
}

func (s ServiceImpl) internalIDGenMiddleware(handler handlerFunc) handlerFunc {
	return func(request *Request) (*Response, error) {
		uid, err := uuid.NewUUID()
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate uuid")
		}
		request.ID = uid.String()

		return handler(request)
	}
}

func (s ServiceImpl) metricRecorderMiddleware(handler handlerFunc) handlerFunc {
	return func(request *Request) (*Response, error) {
		resp, err := handler(request)

		// record
		if err != nil {
			s.OpCounter.WithLabelValues("server_error", string(request.CardOp)).Inc()
		} else {
			s.OpCounter.WithLabelValues(resp.Code, string(request.CardOp)).Inc()
		}

		return resp, err
	}
}

func (s ServiceImpl) call(request *Request) (*Response, error) {
	resp := &Response{
		Code:  CodeSuccess,
		TxnID: request.ID,
	}

	// business error
	err := s.ChannelService.Call(&cybersource.Request{
		// field mapping
		Op: string(request.CardOp),
	})
	if err != nil {
		resp.Code = CodeChannelReject
		resp.Message = err.Error()
	}

	return resp, nil
}

func (s ServiceImpl) Auth(request *Request) (*Response, error) {
	request.CardOp = OpAuth
	// TODO: validation based on struct tag
	if request.CardNo == "" || request.Amount < 0 {
		return &Response{
			Code:    CodeInvalidRequest,
			Message: "card_no or amount invalid",
		}, nil
	}
	// more validation ...

	// invoke 3rd channel
	return s.call(request)
}

func (s ServiceImpl) Capture(request *Request) (*Response, error) {
	request.CardOp = OpCapture
	if request.AuthID == "" || request.Amount < 0 {
		return &Response{
			Code:    CodeInvalidRequest,
			Message: "auth_id or amount invalid",
		}, nil
	}
	// more validation ...

	// invoke 3rd channel
	return s.call(request)
}

func (s ServiceImpl) Refund(request *Request) (*Response, error) {
	request.CardOp = OpRefund
	if request.CaptureID == "" || request.Amount < 0 {
		return &Response{
			Code:    CodeInvalidRequest,
			Message: "capture_id or amount invalid",
		}, nil
	}
	// more validation ...

	// invoke 3rd channel
	return s.call(request)
}
