package proposition_test

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/infrastructure/proposition"
	"SynchronizeMonorevoDeliveryDates/usecase/appsetting_obtain_case"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPropositionTable_FetchAll(t *testing.T) {
	err_read := godotenv.Load(`../../.env`)
	if err_read != nil {
		os.Exit(1)
	}

	logger, _ := zap.NewDevelopment()

	appcnf := &appsetting_obtain_case.AppSettingDto{
		SandboxMode: appsetting_obtain_case.SandboxModeDto{
			Monorevo: false,
		},
	}

	cnf := proposition.TestMonorevoUserConfigCreate()

	tests := []struct {
		name    string
		p       *proposition.PropositionTable
		want    *monorevo.Proposition
		wantErr bool
	}{
		{
			name: "正常系_ものレボから案件を取得できること",
			p: proposition.NewPropositionTable(
				logger.Sugar(),
				appcnf,
				cnf,
			),
			// 2022/9/9現在実データ上にある作業Noのため、将来的に消える可能性がある
			want: monorevo.TestPropositionCreate(
				monorevo.OptWorkedNumber("22T-260"),
				monorevo.OptDET("1"),
				monorevo.OptDeliveryDate(time.Date(2022, 12, 25, 0, 0, 0, 0, time.UTC)),
				monorevo.OptCode("31E"),
			),
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
			var foundIt bool
			for _, v := range got {
				if reflect.DeepEqual(v, *tt.want) {
					foundIt = true
					break
				}
			}
			assert.True(t, foundIt)
		})
	}
}
