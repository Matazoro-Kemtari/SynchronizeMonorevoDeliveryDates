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

// ものレボ案件のドメインモデル
type Proposition struct {
	WorkedNumber string
	Det          string
	DeliveryDate time.Time
}

func NewProposition(warkNumber string, det string, deliveryDate time.Time) *Proposition {
	return &Proposition{
		WorkedNumber: warkNumber,
		Det:          det,
		DeliveryDate: deliveryDate,
	}
}

// テスト用Factoryメソッド
// 参考: https://shiimanblog.com/engineering/functional-options-pattern/
type Options struct {
	WorkedNumber string
	Det          string
	DeliveryDate time.Time
}

type Option func(*Options)

func TestPropositionCreate(options ...Option) *Proposition {
	// デフォルト値設定
	opts := &Options{
		WorkedNumber: "99A-1234",
		Det:          "1",
		DeliveryDate: time.Now(),
	}
	return NewProposition(opts.WorkedNumber, opts.Det, opts.DeliveryDate)
}
