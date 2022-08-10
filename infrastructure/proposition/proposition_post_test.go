package proposition

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"os"
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
		in0 []monorevo.Proposition
	}
	tests := []struct {
		name    string
		p       *PropositionTable
		args    args
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
				[]monorevo.Proposition{
					{
						WorkedNumber: "22T-378", //"99A-9999",
						DeliveryDate: time.Now(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.p.PostRange(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("PropositionTable.PostRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
