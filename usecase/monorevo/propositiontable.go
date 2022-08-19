package monorevo

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"

	"go.uber.org/zap"
)

type PropositionTable struct {
	sugar   *zap.SugaredLogger
	fetcher monorevo.Fetcher
	Poster  monorevo.Poster
}

func NewPropositionTable(
	sugar *zap.SugaredLogger,
	fetcher monorevo.Fetcher,
	poster monorevo.Poster,
) *PropositionTable {
	return &PropositionTable{
		sugar:   sugar,
		fetcher: fetcher,
		Poster:  poster,
	}
}
