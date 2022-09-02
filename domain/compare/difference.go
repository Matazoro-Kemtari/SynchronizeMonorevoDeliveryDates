package compare

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/domain/orderdb"
)

type Extractor interface {
	ExtractForDeliveryDate(
		j []orderdb.JobBook,
		p []monorevo.Proposition,
	) []monorevo.DifferentProposition
}

type Difference struct{}

func NewDifference() *Difference {
	return &Difference{}
}

// ものレボの納期と受注管理DBの納期を比較して 差分を返す
func (e Difference) ExtractForDeliveryDate(j []orderdb.JobBook, p []monorevo.Proposition) []monorevo.DifferentProposition {
	var diff []monorevo.DifferentProposition
	for _, job := range j {
		for _, pp := range p {
			if job.WorkedNumber == pp.WorkedNumber {
				if !job.DeliveryDate.Equal(pp.DeliveryDate) {
					diff = append(diff, *monorevo.NewDifferenceProposition(
						job.WorkedNumber,
						pp.DET,
						pp.DeliveryDate,
						job.DeliveryDate,
					))
				}
				break
			}
		}
	}
	if len(diff) == 0 {
		return nil
	}
	return diff
}
