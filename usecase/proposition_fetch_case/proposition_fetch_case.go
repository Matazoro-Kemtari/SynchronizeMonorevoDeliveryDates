package proposition_fetch_case

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"time"

	"go.uber.org/zap"
)

type FetchedPropositionDto struct {
	WorkedNumber string
	DET          string
	DeliveryDate time.Time
	Code         string
}

type FetchingExecutor interface {
	Execute() ([]FetchedPropositionDto, error)
}

type PropositionFetchingUseCase struct {
	sugar   *zap.SugaredLogger
	fetcher monorevo.MonorevoFetcher
}

func NewPropositionFetchingUseCase(
	sugar *zap.SugaredLogger,
	fetcher monorevo.MonorevoFetcher,
) *PropositionFetchingUseCase {
	return &PropositionFetchingUseCase{
		sugar:   sugar,
		fetcher: fetcher,
	}
}

func (m *PropositionFetchingUseCase) Execute() ([]FetchedPropositionDto, error) {
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
				DET:          v.DET,
				DeliveryDate: v.DeliveryDate,
				Code:         v.Code,
			},
		)
	}
	return cnv, nil
}
