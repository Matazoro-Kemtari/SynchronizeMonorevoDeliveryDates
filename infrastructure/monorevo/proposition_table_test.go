package monorevo

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func TestPropositionTable_FetchAll(t *testing.T) {
	err_read := godotenv.Load(`../../.env`)
	if err_read != nil {
		os.Exit(1)
	}

	logger, _ := zap.NewDevelopment()

	tests := []struct {
		name    string
		p       *PropositionTable
		want    [][]monorevo.Proposition
		wantErr bool
	}{
		{
			name: "正常系_ブラウザが起動すること",
			p: NewPropositionTable(
				logger.Sugar(),
				os.Getenv("MONOREVO_COMPANY_ID"),
				os.Getenv("MONOREVO_USER_ID"),
				os.Getenv("MONOREVO_USER_PASSWORD"),
			),
			want:    [][]monorevo.Proposition{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.FetchAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("PropositionTable.FetchAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PropositionTable.FetchAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
