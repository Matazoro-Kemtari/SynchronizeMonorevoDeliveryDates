package monorevo

import (
	"time"
)

// ものレボから案件を操作する
type Fetcher interface {
	FetchAll() ([]Proposition, error)
}
type Poster interface {
	PostRange([]Proposition) error
}

type Proposition struct {
	WorkedNumber string
	DeliveryDate time.Time
}

func NewProposition(w string, d time.Time) *Proposition {
	return &Proposition{
		WorkedNumber: w,
		DeliveryDate: d,
	}
}

// テスト用Factoryメソッド
// 参考: https://shiimanblog.com/engineering/functional-options-pattern/
type Options struct {
	WorkedNumber string
	deliveryDate time.Time
}

type Option func(*Options)

func TestPropositionCreate(options ...Option) *Proposition {
	// デフォルト値設定
	opts := &Options{
		WorkedNumber: "99A-1234",
		deliveryDate: time.Now(),
	}
	return NewProposition(opts.WorkedNumber, opts.deliveryDate)
}
