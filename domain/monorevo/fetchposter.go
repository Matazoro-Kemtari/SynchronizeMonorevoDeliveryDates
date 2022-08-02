package monorevo

import (
	"time"
)

// ものレボから案件を操作する
type FetchPoster interface {
	FetchAll() ([]Proposition, error)
	PostRange([]Proposition) error
}

type Proposition struct {
	WorkedNumber string
	DeliveryDate time.Time
}
