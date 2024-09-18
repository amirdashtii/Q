package models

type PaymentRequest struct {
	Action      string `json:"action"`
	TerminalId  string `json:"TerminalId"`
	Amount      int64  `json:"Amount"`
	ResNum      string `json:"ResNum"`
	RedirectUrl string `json:"RedirectUrl"`
	CellNumber  string `json:"CellNumber"`
}

type Response struct {
	Status    int    `json:"status"`
	Token     string `json:"token"`
	ErrorCode string `json:"errorCode"`
	ErrorDesc string `json:"errorDesc"`
}

type PaymentReceipt struct {
	MID       string `json:"MID" form:"MID"`
	State     string `json:"State" form:"State"`
	Status    string `json:"Status" form:"Status"`
	RRN       string `json:"RRN" form:"RRN"`
	RefNum    string `json:"RefNum" form:"RefNum"`
	ResNum    string `json:"ResNum" form:"ResNum"`
	TraceNo   string `json:"TraceNo" form:"TraceNo"`
	Amount    string `json:"Amount" form:"Amount"`
	Wage      string `json:"Wage" form:"Wage"`
	CID       string `json:"CID" form:"CID"`
	SecurePan string `json:"SecurePan" form:"SecurePan"`
	Token     string `json:"Token" form:"Token"`
}

type VerifyTransactionRequest struct {
	MID    string `json:"TerminalNumber"`
	RefNum string `json:"RefNum"`
}

type TransactionDetail struct {
	AffectiveAmount float64
	HashedPan       string
	MaskedPan       string
	OrginalAmount   float64
	RRN             string
	RefNum          string
	StraceDate      string
	StraceNo        string
	TerminalNumber  float64
}

type Transaction struct {
	ResultCode        int
	ResultDescription string
	Success           bool
	TransactionDetail TransactionDetail
}

var BankErrorMessage = map[string]string{
	"1":  "canceled by user",
	"2":  "ok",
	"3":  "failed",
	"4":  "session is null",
	"5":  "invalid parameters",
	"8":  "merchant ip address is invalid",
	"10": "token not found",
	"11": "terminal not found",
	"12": "multi settle policy errors",
	"21": "token required",
}
