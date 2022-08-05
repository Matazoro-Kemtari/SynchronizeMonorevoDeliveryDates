package compare

import (
	"SynchronizeMonorevoDeliveryDates/domain/database"
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"time"
)

type Extract struct{}

func NewExtract() *Extract {
	return &Extract{}
}

type DifferenceProposition struct {
	WorkedNumber        string
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

func NewDifferenceProposition(workNumber string, deliveryDate time.Time, updatedDeliveryDate time.Time) *DifferenceProposition {
	return &DifferenceProposition{
		WorkedNumber:        workNumber,
		DeliveryDate:        deliveryDate,
		UpdatedDeliveryDate: updatedDeliveryDate,
	}
}

// ものレボの納期と受注管理DBの納期を比較して 差分を返す
func (e Extract) DifferenceOfDeliveryDate(j []database.JobBook, p []monorevo.Proposition) []DifferenceProposition {
	var diff []DifferenceProposition
	for _, job := range j {
		for _, pp := range p {
			if job.WorkedNumber == pp.WorkedNumber {
				if !job.DeliveryDate.Equal(pp.DeliveryDate) {
					diff = append(diff, *NewDifferenceProposition(
						job.WorkedNumber,
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
