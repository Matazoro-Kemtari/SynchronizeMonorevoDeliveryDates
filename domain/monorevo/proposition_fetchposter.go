package monorevo

import (
	"time"
)

// ものレボから案件を取得する
type MonorevoFetcher interface {
	FetchAll() ([]Proposition, error)
}

// ものレボ案件のドメインモデル
type Proposition struct {
	WorkedNumber string
	DET          string
	DeliveryDate time.Time
	Code         string
}

func NewProposition(warkNumber string, det string, deliveryDate time.Time, code string) *Proposition {
	return &Proposition{
		WorkedNumber: warkNumber,
		DET:          det,
		DeliveryDate: deliveryDate,
		Code:         code,
	}
}

// ものレボに案件を更新する
type MonorevoPoster interface {
	PostRange([]DifferentProposition) ([]UpdatedProposition, error)
}

// ものレボ案件差分
type DifferentProposition struct {
	WorkedNumber        string
	DET                 string
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
	Code                string
}

func NewDifferenceProposition(workNumber string, det string, deliveryDate time.Time, updatedDeliveryDate time.Time, code string) *DifferentProposition {
	return &DifferentProposition{
		WorkedNumber:        workNumber,
		DET:                 det,
		DeliveryDate:        deliveryDate,
		UpdatedDeliveryDate: updatedDeliveryDate,
		Code:                code,
	}
}

// ものレボ案件編集結果
type UpdatedProposition struct {
	WorkedNumber        string
	DET                 string
	Successful          bool
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
	Code                string
}

func NewUpdatedProposition(
	workedNumber string,
	det string,
	successful bool,
	deliveryDate time.Time,
	updatedDeliveryDate time.Time,
	code string,
) *UpdatedProposition {
	return &UpdatedProposition{
		WorkedNumber:        workedNumber,
		DET:                 det,
		Successful:          successful,
		DeliveryDate:        deliveryDate,
		UpdatedDeliveryDate: updatedDeliveryDate,
		Code:                code,
	}
}

// テスト用Factoryメソッド
// 参考: https://shiimanblog.com/engineering/functional-options-pattern/
type PropositionOptions struct {
	WorkedNumber string
	DET          string
	DeliveryDate time.Time
	Code         string
}

type PropositionOption func(*PropositionOptions)

func OptWorkedNumber(v string) PropositionOption {
	return func(opts *PropositionOptions) {
		opts.WorkedNumber = v
	}
}

func OptDET(v string) PropositionOption {
	return func(opts *PropositionOptions) {
		opts.DET = v
	}
}

func OptDeliveryDate(v time.Time) PropositionOption {
	return func(opts *PropositionOptions) {
		opts.DeliveryDate = v
	}
}

func OptCode(v string) PropositionOption {
	return func(opts *PropositionOptions) {
		opts.Code = v
	}
}

func TestPropositionCreate(options ...PropositionOption) *Proposition {
	// デフォルト値設定
	opts := &PropositionOptions{
		WorkedNumber: "99A-1234",
		DET:          "1",
		DeliveryDate: time.Now(),
		Code:         "99A-1",
	}

	for _, option := range options {
		option(opts)
	}

	return NewProposition(opts.WorkedNumber, opts.DET, opts.DeliveryDate, opts.Code)
}

type UpdatedPropositionOptions struct {
	WorkedNumber        string
	DET                 string
	Successful          bool
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
	Code                string
}

type UpdatedPropositionOption func(*UpdatedPropositionOptions)

func TestUpdatedPropositionCreate(options ...UpdatedPropositionOption) *UpdatedProposition {
	// デフォルト値設定
	opts := &UpdatedPropositionOptions{
		WorkedNumber:        "99A-1234",
		DET:                 "1",
		Successful:          true,
		DeliveryDate:        time.Now(),
		UpdatedDeliveryDate: time.Now(),
		Code:                "99A-1",
	}

	for _, option := range options {
		option(opts)
	}

	return NewUpdatedProposition(
		opts.WorkedNumber,
		opts.DET,
		opts.Successful,
		opts.DeliveryDate,
		opts.UpdatedDeliveryDate,
		opts.Code,
	)
}
