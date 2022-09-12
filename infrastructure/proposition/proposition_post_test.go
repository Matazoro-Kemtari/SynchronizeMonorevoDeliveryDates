package proposition_test

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/infrastructure/proposition"
	"SynchronizeMonorevoDeliveryDates/usecase/appsetting_obtain_case"
	"fmt"
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
		DET:                 "1",
		DeliveryDate:        time.Date(2050, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedDeliveryDate: time.Date(2050, 1, 10, 0, 0, 0, 0, time.UTC),
		Code:                "22B-1",
	}
	pastCase := monorevo.DifferentProposition{
		WorkedNumber:        "22T-378",
		DET:                 "1",
		DeliveryDate:        time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedDeliveryDate: time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
		Code:                "22C-1",
	}
	d := time.Date(
		time.Now().Year(),
		time.Now().Month(),
		time.Now().Day(),
		0, 0, 0, 0, time.UTC)
	updatableCases := []monorevo.DifferentProposition{
		{
			WorkedNumber:        "99仮-1",
			DET:                 "1",
			DeliveryDate:        time.Date(2222, 8, 21, 0, 0, 0, 0, time.UTC),
			UpdatedDeliveryDate: time.Date(2222, 8, 22, 0, 0, 0, 0, time.UTC),
			Code:                "",
		},
	}

	appcnf := &appsetting_obtain_case.AppSettingDto{
		SandboxMode: appsetting_obtain_case.SandboxModeDto{
			Monorevo: false,
		},
	}

	cnf := proposition.TestMonorevoUserConfigCreate()

	type args struct {
		in0 []monorevo.DifferentProposition
	}
	tests := []struct {
		name    string
		p       *proposition.PropositionTable
		args    args
		want    []monorevo.UpdatedProposition
		wantErr bool
	}{
		{
			name: "異常系_存在しない作業Noはものレボ案件を更新するとエラーになること",
			p: proposition.NewPropositionTable(
				logger.Sugar(),
				appcnf,
				cnf,
			),
			args: args{
				[]monorevo.DifferentProposition{
					nonexisitentCase,
				},
			},
			want: []monorevo.UpdatedProposition{
				{
					WorkedNumber:        nonexisitentCase.WorkedNumber,
					DET:                 nonexisitentCase.DET,
					Successful:          false,
					Reason:              "ものレボ上で案件検索で失敗した",
					DeliveryDate:        nonexisitentCase.DeliveryDate,
					UpdatedDeliveryDate: nonexisitentCase.UpdatedDeliveryDate,
					Code:                nonexisitentCase.Code,
				},
			},
			wantErr: false,
		},
		{
			name: "異常系_納期を過去日で更新しようとするとエラーになること",
			p: proposition.NewPropositionTable(
				logger.Sugar(),
				appcnf,
				cnf,
			),
			args: args{
				[]monorevo.DifferentProposition{
					pastCase,
				},
			},
			want: []monorevo.UpdatedProposition{
				{
					WorkedNumber:        pastCase.WorkedNumber,
					DET:                 pastCase.DET,
					Successful:          false,
					Reason:              fmt.Sprintf("現在日(%v)より過去の納期(%v)は受付できない", d.Format("2006/01/02"), pastCase.UpdatedDeliveryDate.Format("2006/01/02")),
					DeliveryDate:        pastCase.DeliveryDate,
					UpdatedDeliveryDate: pastCase.UpdatedDeliveryDate,
					Code:                pastCase.Code,
				},
			},
			wantErr: false,
		},
		{
			name: "正常系_納期が更新できること",
			p: proposition.NewPropositionTable(
				logger.Sugar(),
				appcnf,
				cnf,
			),
			args: args{
				[]monorevo.DifferentProposition{
					updatableCases[0],
				},
			},
			want: []monorevo.UpdatedProposition{
				{
					WorkedNumber:        updatableCases[0].WorkedNumber,
					DET:                 updatableCases[0].DET,
					Successful:          true,
					Reason:              "",
					DeliveryDate:        updatableCases[0].DeliveryDate,
					UpdatedDeliveryDate: updatableCases[0].UpdatedDeliveryDate,
					Code:                updatableCases[0].Code,
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
