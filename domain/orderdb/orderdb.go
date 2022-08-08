package orderdb

import (
	"time"
)

// 受注管理からM作業台帳を問い合わせする
type JobBookFetcher interface {
	FetchAll() ([]JobBook, error)
}

type JobBook struct {
	WorkedNumber string
	DeliveryDate time.Time
}
