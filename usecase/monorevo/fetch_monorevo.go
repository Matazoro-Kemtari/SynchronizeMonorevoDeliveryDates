package monorevo

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"time"

	"go.uber.org/zap"
)

type PropositionDto struct {
	WorkedNumber string
	DeliveryDate time.Time
}

type Executer interface {
	Execute() ([]PropositionDto, error)
}

type FetchPropositionTable struct {
	sugar   *zap.SugaredLogger
	fetcher monorevo.Fetcher
}

func NewFetchPropositionTable(
	sugar *zap.SugaredLogger,
	fetcher monorevo.Fetcher,
) *FetchPropositionTable {
	return &FetchPropositionTable{
		sugar:   sugar,
		fetcher: fetcher,
	}
}

func (m *FetchPropositionTable) Execute() ([]PropositionDto, error) {
	propositions, err := m.fetcher.FetchAll()
	if err != nil {
		m.sugar.Fatal("ものレボから案件一覧の取得に失敗しました", err)
	}

	// DTOに詰め替え
	dif := []PropositionDto{}
	for _, v := range propositions {
		dif = append(
			dif,
			PropositionDto{
				WorkedNumber: v.WorkedNumber,
				DeliveryDate: v.DeliveryDate,
			},
		)
	}
	return dif, nil
}
