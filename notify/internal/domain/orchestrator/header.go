package orchestrator

import "time"

type Header struct {
	TransactionId       string
	TransactionDateTime time.Time
	OrderType           string
	OrderService        string
	Retries             int
	ResponseCode        int
	ResponseMessage     string
}
