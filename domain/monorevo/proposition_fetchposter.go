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

// ものレボ案件差分
type DifferentProposition struct {
	WorkedNumber        string
	Det                 string
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

func NewDifferenceProposition(workNumber string, det string, deliveryDate time.Time, updatedDeliveryDate time.Time) *DifferentProposition {
	return &DifferentProposition{
		WorkedNumber:        workNumber,
		Det:                 det,
		DeliveryDate:        deliveryDate,
		UpdatedDeliveryDate: updatedDeliveryDate,
	}
}

// ものレボ案件編集結果
type UpdatedProposition struct {
	WorkedNumber        string
	Det                 string
	Successful          bool
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

func NewUpdatedProposition(
	workedNumber string,
	det string,
	successful bool,
	deliveryDate time.Time,
	updatedDeliveryDate time.Time,
) *UpdatedProposition {
	return &UpdatedProposition{
		WorkedNumber:        workedNumber,
		Det:                 det,
		Successful:          successful,
		DeliveryDate:        deliveryDate,
		UpdatedDeliveryDate: updatedDeliveryDate,
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
