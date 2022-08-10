package compare

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/domain/orderdb"
	"time"
)

type DifferentProposition struct {
	WorkedNumber        string
	Det                 string
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

type Extractor interface {
	ExtractForDeliveryDate(
		j []orderdb.JobBook,
		p []monorevo.Proposition,
	) []DifferentProposition
}

type Difference struct{}

func NewDifference() *Difference {
	return &Difference{}
}

func NewDifferenceProposition(workNumber string, det string, deliveryDate time.Time, updatedDeliveryDate time.Time) *DifferentProposition {
	return &DifferentProposition{
		WorkedNumber:        workNumber,
		Det:                 det,
		DeliveryDate:        deliveryDate,
		UpdatedDeliveryDate: updatedDeliveryDate,
	}
}

// ものレボの納期と受注管理DBの納期を比較して 差分を返す
func (e Difference) ExtractForDeliveryDate(j []orderdb.JobBook, p []monorevo.Proposition) []DifferentProposition {
	var diff []DifferentProposition
	for _, job := range j {
		for _, pp := range p {
			if job.WorkedNumber == pp.WorkedNumber {
				if !job.DeliveryDate.Equal(pp.DeliveryDate) {
					diff = append(diff, *NewDifferenceProposition(
						job.WorkedNumber,
						pp.Det,
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
