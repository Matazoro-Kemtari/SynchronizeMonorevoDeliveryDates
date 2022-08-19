package proposition

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func TestPropositionTable_PostRange(t *testing.T) {
	err_read := godotenv.Load(`../../.env`)
	if err_read != nil {
		os.Exit(1)
	}

	logger, _ := zap.NewDevelopment()

	nonexisitentCase := monorevo.DifferentProposition{

		WorkedNumber:        "99A-9999",
		Det:                 "1",
		DeliveryDate:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedDeliveryDate: time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	pastCase := monorevo.DifferentProposition{
		WorkedNumber:        "22T-378",
		Det:                 "1",
		DeliveryDate:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedDeliveryDate: time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
	}
	updatableCases := []monorevo.DifferentProposition{
		{
			WorkedNumber:        "99仮-1",
			Det:                 "1",
			DeliveryDate:        time.Date(2222, 8, 21, 0, 0, 0, 0, time.UTC),
			UpdatedDeliveryDate: time.Date(2222, 8, 22, 0, 0, 0, 0, time.UTC),
		}, {
			WorkedNumber:        "99仮-1",
			Det:                 "2",
			DeliveryDate:        time.Date(2222, 10, 21, 0, 0, 0, 0, time.UTC),
			UpdatedDeliveryDate: time.Date(2222, 10, 22, 0, 0, 0, 0, time.UTC),
		},
	}
	type args struct {
		in0 []monorevo.DifferentProposition
	}
	tests := []struct {
		name    string
		p       *PropositionTable
		args    args
		want    []monorevo.UpdatedProposition
		wantErr bool
	}{
		{
			name: "異常系_存在しない作業Noはものレボ案件を更新するとエラーになること",
			p: NewPropositionTable(
				logger.Sugar(),
				os.Getenv("MONOREVO_COMPANY_ID"),
				os.Getenv("MONOREVO_USER_ID"),
				os.Getenv("MONOREVO_USER_PASSWORD"),
			),
			args: args{
				[]monorevo.DifferentProposition{
					nonexisitentCase,
				},
			},
			want: []monorevo.UpdatedProposition{
				{
					WorkedNumber:        nonexisitentCase.WorkedNumber,
					Det:                 nonexisitentCase.Det,
					Successful:          false,
					DeliveryDate:        nonexisitentCase.DeliveryDate,
					UpdatedDeliveryDate: nonexisitentCase.UpdatedDeliveryDate,
				},
			},
			wantErr: false,
		},
		{
			name: "異常系_納期を過去日で更新しようとするとエラーになること",
			p: NewPropositionTable(
				logger.Sugar(),
				os.Getenv("MONOREVO_COMPANY_ID"),
				os.Getenv("MONOREVO_USER_ID"),
				os.Getenv("MONOREVO_USER_PASSWORD"),
			),
			args: args{
				[]monorevo.DifferentProposition{
					pastCase,
				},
			},
			want: []monorevo.UpdatedProposition{
				{
					WorkedNumber:        pastCase.WorkedNumber,
					Det:                 pastCase.Det,
					Successful:          false,
					DeliveryDate:        pastCase.DeliveryDate,
					UpdatedDeliveryDate: pastCase.UpdatedDeliveryDate,
				},
			},
			wantErr: false,
		},
		{
			name: "正常系_納期が更新できること",
			p: NewPropositionTable(
				logger.Sugar(),
				os.Getenv("MONOREVO_COMPANY_ID"),
				os.Getenv("MONOREVO_USER_ID"),
				os.Getenv("MONOREVO_USER_PASSWORD"),
			),
			args: args{
				[]monorevo.DifferentProposition{
					updatableCases[0],
					updatableCases[1],
				},
			},
			want: []monorevo.UpdatedProposition{
				{
					WorkedNumber:        updatableCases[0].WorkedNumber,
					Det:                 updatableCases[0].Det,
					Successful:          true,
					DeliveryDate:        updatableCases[0].DeliveryDate,
					UpdatedDeliveryDate: updatableCases[0].UpdatedDeliveryDate,
				},
				{
					WorkedNumber:        updatableCases[1].WorkedNumber,
					Det:                 updatableCases[1].Det,
					Successful:          true,
					DeliveryDate:        updatableCases[1].DeliveryDate,
					UpdatedDeliveryDate: updatableCases[1].UpdatedDeliveryDate,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.PostRange(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("PropositionTable.PostRange() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(got) != len(tt.want) {
				t.Errorf("PropositionTable.PostRange()#len = %v, want#len %v", len(got), len(tt.want))
			}
			for i := 0; i < len(got); i++ {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("PropositionTable.PostRange() (index:%v)= %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
