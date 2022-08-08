package difference

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
	DeliveryDate time.Time
}

type DifferenceSourcePram struct {
	JobBooks     []JobBookPram
	Propositions []PropositionPram
}

type DifferentPropositionDto struct {
	WorkedNumber        string
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

type Executer interface {
	Execute() []DifferentPropositionDto
}

type ExtractProposition struct {
	sugar     *zap.SugaredLogger
	extractor compare.Extractor
}

func NewExtractProposition(
	sugar *zap.SugaredLogger,
	extractor compare.Extractor,
) *ExtractProposition {
	return &ExtractProposition{
		sugar:     sugar,
		extractor: extractor,
	}
}

func (m *ExtractProposition) Execute(s DifferenceSourcePram) []DifferentPropositionDto {
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
			v.DeliveryDate,
		))
	}
	diff := m.extractor.ExtractForDeliveryDate(j, p)

	// DTOに詰め替え
	cnv := []DifferentPropositionDto{}
	for _, v := range diff {
		cnv = append(
			cnv,
			DifferentPropositionDto{
				WorkedNumber:        v.WorkedNumber,
				DeliveryDate:        v.DeliveryDate,
				UpdatedDeliveryDate: v.UpdatedDeliveryDate,
			},
		)
	}
	return cnv
}
