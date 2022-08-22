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
	Det          string
	DeliveryDate time.Time
}

type DifferenceSourcePram struct {
	JobBooks     []JobBookPram
	Propositions []PropositionPram
}

type DifferentPropositionDto struct {
	WorkedNumber        string
	Det                 string
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

type Extractor interface {
	Extract(DifferenceSourcePram) []DifferentPropositionDto
}

type ExtractingProposition struct {
	sugar     *zap.SugaredLogger
	extractor compare.Extractor
}

func NewExtractingProposition(
	sugar *zap.SugaredLogger,
	extractor compare.Extractor,
) *ExtractingProposition {
	return &ExtractingProposition{
		sugar:     sugar,
		extractor: extractor,
	}
}

func (m *ExtractingProposition) Extract(s DifferenceSourcePram) []DifferentPropositionDto {
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
			v.Det,
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
				Det:                 v.Det,
				DeliveryDate:        v.DeliveryDate,
				UpdatedDeliveryDate: v.UpdatedDeliveryDate,
			},
		)
	}
	return cnv
}
