package proposition_post_case

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"time"

	"go.uber.org/zap"
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

type PostingExecutor interface {
	Execute([]PostingPropositionPram) ([]PostedPropositionDto, error)
}

type PropositionPostingUseCase struct {
	sugar  *zap.SugaredLogger
	Poster monorevo.MonorevoPoster
}

func NewPropositionPostingUseCase(
	sugar *zap.SugaredLogger,
	poster monorevo.MonorevoPoster,
) *PropositionPostingUseCase {
	return &PropositionPostingUseCase{
		sugar:  sugar,
		Poster: poster,
	}
}

func (m *PropositionPostingUseCase) Execute(p []PostingPropositionPram) ([]PostedPropositionDto, error) {
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
