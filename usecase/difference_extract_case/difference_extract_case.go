package difference_extract_case

import (
	"SynchronizeMonorevoDeliveryDates/domain/compare"
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/domain/orderdb"
	"time"

	"go.uber.org/zap"
)

type JobBookPram struct {
	WorkedNumber string
	DeliveryDate time.Time
}

type PropositionPram struct {
	WorkedNumber string
	DET          string
	DeliveryDate time.Time
}

type DifferenceSourcePram struct {
	JobBooks     []JobBookPram
	Propositions []PropositionPram
}

type DifferentPropositionDto struct {
	WorkedNumber        string
	DET                 string
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

type Executor interface {
	Execute(DifferenceSourcePram) []DifferentPropositionDto
}

type PropositionExtractingUseCase struct {
	sugar     *zap.SugaredLogger
	extractor compare.Extractor
}

func NewExtractingPropositionUseCase(
	sugar *zap.SugaredLogger,
	extractor compare.Extractor,
) *PropositionExtractingUseCase {
	return &PropositionExtractingUseCase{
		sugar:     sugar,
		extractor: extractor,
	}
}

func (m *PropositionExtractingUseCase) Execute(s DifferenceSourcePram) []DifferentPropositionDto {
	j := []orderdb.JobBook{}
	for _, v := range s.JobBooks {
		j = append(j, *orderdb.NewJobBook(
			v.WorkedNumber,
			v.DeliveryDate,
		))
	}
	p := []monorevo.Proposition{}
	for _, v := range s.Propositions {
		p = append(p, *monorevo.NewProposition(
			v.WorkedNumber,
			v.DET,
			v.DeliveryDate,
		))
	}
	diff := m.extractor.ExtractForDeliveryDate(j, p)
	if diff == nil {
		return nil
	}

	// DTOに詰め替え
	cnv := []DifferentPropositionDto{}
	for _, v := range diff {
		cnv = append(
			cnv,
			DifferentPropositionDto{
				WorkedNumber:        v.WorkedNumber,
				DET:                 v.DET,
				DeliveryDate:        v.DeliveryDate,
				UpdatedDeliveryDate: v.UpdatedDeliveryDate,
			},
		)
	}
	return cnv
}
