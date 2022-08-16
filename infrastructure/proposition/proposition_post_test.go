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
					{
						WorkedNumber:        "99A-9999",
						Det:                 "1",
						DeliveryDate:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want: []monorevo.UpdatedProposition{
				{
					WorkedNumber:        "99A-9999",
					Det:                 "1",
					Successful:          false,
					DeliveryDate:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedDeliveryDate: time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
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
					{
						WorkedNumber:        "22T-378",
						Det:                 "1",
						DeliveryDate:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
						UpdatedDeliveryDate: time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
					},
				},
			},
			want: []monorevo.UpdatedProposition{
				{
					WorkedNumber:        "22T-378",
					Det:                 "1",
					Successful:          false,
					DeliveryDate:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedDeliveryDate: time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
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
					t.Errorf("PropositionTable.PostRange() (i:%v)= %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
