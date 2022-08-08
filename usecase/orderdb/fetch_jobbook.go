package orderdb

import (
	"SynchronizeMonorevoDeliveryDates/domain/orderdb"
	"time"

	"go.uber.org/zap"
)

type JobBookDto struct {
	WorkedNumber string
	DeliveryDate time.Time
}

type Executer interface {
	Execute() ([]JobBookDto, error)
}

type FetchJobBookTable struct {
	sugar          *zap.SugaredLogger
	jobBookFetcher orderdb.JobBookFetcher
}

func NewFetchJobBookTable(
	sugar *zap.SugaredLogger,
	jobBookFetcher orderdb.JobBookFetcher,
) *FetchJobBookTable {
	return &FetchJobBookTable{
		sugar:          sugar,
		jobBookFetcher: jobBookFetcher,
	}
}

func (m *FetchJobBookTable) Execute() ([]JobBookDto, error) {
	jb, err := m.jobBookFetcher.FetchAll()
	if err != nil {
		m.sugar.Fatal("受注管理DBから作業台帳を取得できませんでした", err)
	}

	// 詰め替え
	dto := []JobBookDto{}
	for _, v := range jb {
		dto = append(
			dto,
			JobBookDto{
				WorkedNumber: v.WorkedNumber,
				DeliveryDate: v.DeliveryDate,
			},
		)
	}
	return dto, nil
}
