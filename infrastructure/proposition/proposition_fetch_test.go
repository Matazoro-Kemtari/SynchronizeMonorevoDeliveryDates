package proposition_test

import (
	"SynchronizeMonorevoDeliveryDates/infrastructure/proposition"
	"SynchronizeMonorevoDeliveryDates/usecase/appsetting_obtain_case"
	"os"
	"regexp"
	"testing"

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
		want    string
		wantErr bool
	}{
		{
			name: "正常系_ものレボから案件を取得できること",
			p: proposition.NewPropositionTable(
				logger.Sugar(),
				appcnf,
				cnf,
			),
			want:    `X?[0-9]{2}[A-Z]-[0-9]{1,4}`,
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
			assert.True(
				t,
				regexp.MustCompile(tt.want).
					Match(
						[]byte(got[0].WorkedNumber),
					),
			)
			assert.NotEmpty(t, got[0].DeliveryDate)
		})
	}
}
