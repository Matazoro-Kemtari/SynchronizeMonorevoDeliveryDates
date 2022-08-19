package monorevo

import (
	"time"
)

type FetchedPropositionDto struct {
	WorkedNumber string
	DeliveryDate time.Time
}

type Fetcher interface {
	Fetch() ([]FetchedPropositionDto, error)
}

func (m *PropositionTable) Fetch() ([]FetchedPropositionDto, error) {
	propositions, err := m.fetcher.FetchAll()
	if err != nil {
		m.sugar.Fatalf("ものレボから案件一覧の取得に失敗しました error: %v", err)
	}

	// DTOに詰め替え
	cnv := []FetchedPropositionDto{}
	for _, v := range propositions {
		cnv = append(
			cnv,
			FetchedPropositionDto{
				WorkedNumber: v.WorkedNumber,
				DeliveryDate: v.DeliveryDate,
			},
		)
	}
	return cnv, nil
}
