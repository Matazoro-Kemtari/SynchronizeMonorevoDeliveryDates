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

func NewJobBook(
	w string,
	d time.Time,
) *JobBook {
	return &JobBook{
		WorkedNumber: w,
		DeliveryDate: d,
	}
}

// テスト用Factoryメソッド
type Options struct {
	WorkedNumber string
	DeliveryDate time.Time
}

type Option func(*Options)

func OptWorkedNumber(v string) Option {
	return func(opts *Options) {
		opts.WorkedNumber = v
	}
}

func OptDeliveryDate(v time.Time) Option {
	return func(opts *Options) {
		opts.DeliveryDate = v
	}
}

func TestJobBookCreate(options ...Option) *JobBook {
	// デフォルト値設定
	opts := &Options{
		WorkedNumber: "99A-1234",
		DeliveryDate: time.Now(),
	}

	for _, option := range options {
		option(opts)
	}
	return NewJobBook(opts.WorkedNumber, opts.DeliveryDate)
}
