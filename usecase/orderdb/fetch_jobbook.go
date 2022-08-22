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

type Fetcher interface {
	Fetch() ([]JobBookDto, error)
}

type JobBookTable struct {
	sugar          *zap.SugaredLogger
	jobBookFetcher orderdb.JobBookFetcher
}

func NewJobBookTable(
	sugar *zap.SugaredLogger,
	jobBookFetcher orderdb.JobBookFetcher,
) *JobBookTable {
	return &JobBookTable{
		sugar:          sugar,
		jobBookFetcher: jobBookFetcher,
	}
}

func (m *JobBookTable) Fetch() ([]JobBookDto, error) {
	job, err := m.jobBookFetcher.FetchAll()
	if err != nil {
		m.sugar.Fatal("受注管理DBから作業台帳を取得できませんでした", err)
	}

	// 詰め替え
	dto := []JobBookDto{}
	for _, v := range job {
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
