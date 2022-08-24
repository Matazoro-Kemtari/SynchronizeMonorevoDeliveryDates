package monorevo

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"

	"go.uber.org/zap"
)

type PropositionTable struct {
	sugar   *zap.SugaredLogger
	fetcher monorevo.MonorevoFetcher
	Poster  monorevo.MonorevoPoster
}

func NewPropositionTable(
	sugar *zap.SugaredLogger,
	fetcher monorevo.MonorevoFetcher,
	poster monorevo.MonorevoPoster,
) *PropositionTable {
	return &PropositionTable{
		sugar:   sugar,
		fetcher: fetcher,
		Poster:  poster,
	}
}
