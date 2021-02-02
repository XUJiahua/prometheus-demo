package card

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"prometheus-demo/chan/cybersource"
)

type Service interface {
	Auth(request *Request) (*Response, error)
	Capture(request *Request) (*Response, error)
	Refund(request *Request) (*Response, error)
}

func NewService(cybersourceService cybersource.Service) Service {
	return &service{
		chanService: cybersourceService,
	}
}

type service struct {
	chanService cybersource.Service
}

func (s service) call(request *Request) (*Response, error) {
	resp := &Response{
		Code:  CodeSuccess,
		TxnID: request.ID,
	}

	// business error
	err := s.chanService.Call()
	if err != nil {
		resp.Code = CodeChannelReject
		resp.Message = err.Error()
	}

	return resp, nil
}

func (s service) Auth(request *Request) (*Response, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate uuid")
	}
	request.ID = uid.String()

	// TODO: validation based on struct tag
	if request.CardNo == "" || request.Amount < 0 {
		return &Response{
			Code:    CodeInvalidRequest,
			Message: "card_no or amount invalid",
		}, nil
	}
	// more validation ...

	// invoke 3rd channel
	request.CardOp = OpAuth
	return s.call(request)
}

func (s service) Capture(request *Request) (*Response, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate uuid")
	}
	request.ID = uid.String()

	if request.AuthID == "" || request.Amount < 0 {
		return &Response{
			Code:    CodeInvalidRequest,
			Message: "auth_id or amount invalid",
		}, nil
	}
	// more validation ...

	// invoke 3rd channel
	request.CardOp = OpCapture
	return s.call(request)
}

func (s service) Refund(request *Request) (*Response, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate uuid")
	}
	request.ID = uid.String()

	if request.CaptureID == "" || request.Amount < 0 {
		return &Response{
			Code:    CodeInvalidRequest,
			Message: "capture_id or amount invalid",
		}, nil
	}
	// more validation ...

	// invoke 3rd channel
	request.CardOp = OpRefund
	return s.call(request)
}
