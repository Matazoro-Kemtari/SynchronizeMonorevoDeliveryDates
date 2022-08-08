package orderdb

import (
	"SynchronizeMonorevoDeliveryDates/domain/database"
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
	jobBookFetcher database.JobBookFetcher
}

func NewFetchJobBookTable(sugar *zap.SugaredLogger) *FetchJobBookTable {
	return &FetchJobBookTable{
		sugar: sugar,
	}
}

func (m *FetchJobBookTable) Execute() ([]JobBookDto, error) {
	return nil, nil
}
