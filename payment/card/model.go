package card

type Op string

var (
	OpAuth    Op = "auth"
	OpCapture Op = "capture"
	OpRefund  Op = "refund"
)

type Request struct {
	// internal ID
	ID string `json:"-"`
	// internal use
	CardOp Op      `json:"-"`
	CardNo string  `json:"card_no"`
	Amount float64 `json:"amount"`
	// return from Auth Response
	AuthID string `json:"auth_id"`
	// return from Capture Response
	CaptureID string `json:"capture_id"`
}

type Response struct {
	// common field
	Code string `json:"code"`
	// common field
	Message string `json:"message"`
	// ID of Auth/Capture/Refund Transaction
	TxnID string `json:"txn_id"`
}

// codes
var (
	CodeSuccess        = "00"
	CodeInvalidRequest = "10"
	CodeChannelReject  = "20"
)
