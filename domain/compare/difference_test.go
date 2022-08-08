package compare

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/domain/orderdb"
	"reflect"
	"testing"
	"time"
)

func TestDifference_ExtractForDeliveryDate(t *testing.T) {
	type args struct {
		j []orderdb.JobBook
		p []monorevo.Proposition
	}
	tests := []struct {
		name string
		e    *Difference
		args args
		want []DifferenceProposition
	}{
		{
			name: "正常系_作業Noが同じ注文の納期に差分が無いときはnilを返すこと",
			e:    NewDifference(),
			args: args{
				j: []orderdb.JobBook{
					{
						WorkedNumber: "99A-1",
						DeliveryDate: time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber: "99A-2",
						DeliveryDate: time.Date(3000, 2, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber: "99A-3",
						DeliveryDate: time.Date(3000, 3, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber: "99G-1",
						DeliveryDate: time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber: "99G-2",
						DeliveryDate: time.Date(3000, 2, 1, 0, 0, 0, 0, time.UTC),
					},
				},
				p: []monorevo.Proposition{
					{
						WorkedNumber: "99A-1",
						DeliveryDate: time.Date(3000, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber: "99A-2",
						DeliveryDate: time.Date(3000, 2, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber: "99A-3",
						DeliveryDate: time.Date(3000, 3, 1, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want: nil,
		},
		{
			name: "正常系_作業Noが同じ注文の納期に差がある案件の2つの納期を返すこと",
			e:    NewDifference(),
			args: args{
				j: []orderdb.JobBook{
					{
						WorkedNumber: "99B-1",
						DeliveryDate: time.Date(3000, 1, 10, 0, 0, 0, 0, time.UTC),
					},
					{ // 日付が違う(大きい)
						WorkedNumber: "99B-2",
						DeliveryDate: time.Date(3000, 2, 20, 0, 0, 0, 0, time.UTC),
					},
					{ // 日付が違う(小さい)
						WorkedNumber: "99B-3",
						DeliveryDate: time.Date(3000, 3, 9, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber: "99G-1",
						DeliveryDate: time.Date(3000, 1, 10, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber: "99G-2",
						DeliveryDate: time.Date(3000, 2, 10, 0, 0, 0, 0, time.UTC),
					},
				},
				p: []monorevo.Proposition{
					{
						WorkedNumber: "99B-1",
						DeliveryDate: time.Date(3000, 1, 10, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber: "99B-2",
						DeliveryDate: time.Date(3000, 2, 10, 0, 0, 0, 0, time.UTC),
					},
					{
						WorkedNumber: "99B-3",
						DeliveryDate: time.Date(3000, 3, 10, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want: []DifferenceProposition{
				{
					WorkedNumber:        "99B-2",
					DeliveryDate:        time.Date(3000, 2, 10, 0, 0, 0, 0, time.UTC),
					UpdatedDeliveryDate: time.Date(3000, 2, 20, 0, 0, 0, 0, time.UTC),
				},
				{
					WorkedNumber:        "99B-3",
					DeliveryDate:        time.Date(3000, 3, 10, 0, 0, 0, 0, time.UTC),
					UpdatedDeliveryDate: time.Date(3000, 3, 9, 0, 0, 0, 0, time.UTC),
				},
			},
		},
		{
			name: "正常系_受注管理DBがnilのときはnilを返すこと",
			e:    NewDifference(),
			args: args{
				j: nil,
				p: []monorevo.Proposition{
					{
						WorkedNumber: "12A-345",
						DeliveryDate: time.Now(),
					},
				},
			},
			want: nil,
		},
		{
			name: "正常系_ものレボがnilのときはnilを返すこと",
			e:    NewDifference(),
			args: args{
				j: []orderdb.JobBook{
					{
						WorkedNumber: "12A-345",
						DeliveryDate: time.Now(),
					},
				},
				p: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.ExtractForDeliveryDate(tt.args.j, tt.args.p); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Extract.DifferenceOfDeliveryDate() = %v, want %v", got, tt.want)
			}
		})
	}
}
