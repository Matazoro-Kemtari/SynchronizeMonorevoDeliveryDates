package monorevo

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"time"
)

type PostingPropositionPram struct {
	WorkedNumber        string
	Det                 string
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

type PostedPropositionDto struct {
	WorkedNumber        string
	Det                 string
	Successful          bool
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

type Poster interface {
	PostRange([]PostingPropositionPram) ([]PostedPropositionDto, error)
}

func (m *PropositionTable) PostRange(p []PostingPropositionPram) ([]PostedPropositionDto, error) {
	diffs := []monorevo.DifferentProposition{}
	for _, v := range p {
		diffs = append(
			diffs,
			*monorevo.NewDifferenceProposition(
				v.WorkedNumber,
				v.Det,
				v.DeliveryDate,
				v.UpdatedDeliveryDate,
			),
		)
	}
	res, err := m.Poster.PostRange(diffs)
	if err != nil {
		m.sugar.Fatalf("ものレボへ案件の更新に失敗しました error: %v", err)
	}

	// DTOに詰め替え
	cnv := []PostedPropositionDto{}
	for _, v := range res {
		cnv = append(
			cnv,
			PostedPropositionDto{
				WorkedNumber:        v.WorkedNumber,
				Det:                 v.Det,
				Successful:          v.Successful,
				DeliveryDate:        v.DeliveryDate,
				UpdatedDeliveryDate: v.UpdatedDeliveryDate,
			},
		)
	}
	return cnv, nil
}
