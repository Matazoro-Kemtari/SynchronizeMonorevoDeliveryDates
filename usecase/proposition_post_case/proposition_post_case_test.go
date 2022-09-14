package proposition_post_case_test

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/domain/monorevo/mock_monorevo"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_post_case"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestPropositionPostingUseCase_Execute(t *testing.T) {
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
	mock_poster := mock_monorevo.NewMockMonorevoPoster(ctrl)
	// EXPECTはctrl#Finishが呼び出される前に FetchAllを呼び出さなければエラーになる
	mock_poster.EXPECT().PostRange(gomock.Any()).Return(mock_results, nil)

	// UseCase戻り値
	results := []proposition_post_case.PostedPropositionDto{}
	for _, v := range mock_results {
		results = append(results,
			proposition_post_case.PostedPropositionDto{
				WorkedNumber:        v.WorkedNumber,
				DET:                 v.DET,
				Successful:          v.Successful,
				Reason:              v.Reason,
				DeliveryDate:        v.DeliveryDate,
				UpdatedDeliveryDate: v.UpdatedDeliveryDate,
				Code:                v.Code,
			},
		)
	}

	type args struct {
		p []proposition_post_case.PostingPropositionPram
	}
	tests := []struct {
		name    string
		m       *proposition_post_case.PropositionPostingUseCase
		args    args
		want    []proposition_post_case.PostedPropositionDto
		wantErr bool
	}{
		{
			name: "正常系_UseCaseを実行するとモックが実行されること",
			m: proposition_post_case.NewPropositionPostingUseCase(
				logger.Sugar(),
				mock_poster,
			),
			args: args{
				p: []proposition_post_case.PostingPropositionPram{
					{
						WorkedNumber:        "99A-1234",
						DET:                 "1",
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
			got, err := tt.m.Execute(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("PropositionPostingUseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PropositionPostingUseCase.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
