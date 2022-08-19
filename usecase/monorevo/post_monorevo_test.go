package monorevo

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/domain/monorevo/mock_monorevo"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestPropositionTable_PostRange(t *testing.T) {
	// logger生成
	logger, _ := zap.NewDevelopment()

	// PostRange戻り値
	mock_results := []monorevo.UpdatedProposition{
		*monorevo.TestUpdatedPropositionCreate(),
	}

	// モックコントローラーの生成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// ものレボDIオブジェクト生成
	mock_poster := mock_monorevo.NewMockPoster(ctrl)
	// EXPECTはctrl#Finishが呼び出される前に FetchAllを呼び出さなければエラーになる
	mock_poster.EXPECT().PostRange(gomock.Any()).Return(mock_results, nil)

	// UseCase戻り値
	results := []PostedPropositionDto{}
	for _, v := range mock_results {
		results = append(results,
			PostedPropositionDto{
				WorkedNumber:        v.WorkedNumber,
				Det:                 v.Det,
				Successful:          v.Successful,
				DeliveryDate:        v.DeliveryDate,
				UpdatedDeliveryDate: v.UpdatedDeliveryDate,
			},
		)
	}

	type args struct {
		p []PostingPropositionPram
	}
	tests := []struct {
		name    string
		m       *PropositionTable
		args    args
		want    []PostedPropositionDto
		wantErr bool
	}{
		{
			name: "正常系_UseCaseを実行するとモックが実行されること",
			m: NewPropositionTable(
				logger.Sugar(),
				nil,
				mock_poster,
			),
			args: args{
				p: []PostingPropositionPram{
					{
						WorkedNumber:        "99A-1234",
						Det:                 "1",
						DeliveryDate:        time.Now(),
						UpdatedDeliveryDate: time.Now(),
					},
				},
			},
			want:    results,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.PostRange(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("PropositionTable.PostRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PropositionTable.PostRange() = %v, want %v", got, tt.want)
			}
		})
	}
}
