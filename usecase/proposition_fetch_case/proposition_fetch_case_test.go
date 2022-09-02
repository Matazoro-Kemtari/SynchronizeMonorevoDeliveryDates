package proposition_fetch_case_test

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/domain/monorevo/mock_monorevo"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_fetch_case"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestPropositionFetchingUseCase_Execute(t *testing.T) {
	// logger生成
	logger, _ := zap.NewDevelopment()

	// FetchAll戻り値
	mock_results := []monorevo.Proposition{}
	mock_pro := monorevo.TestPropositionCreate()
	mock_results = append(mock_results, *mock_pro)

	// モックコントローラーの生成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// ものレボDIオブジェクト生成
	mock_fetcher := mock_monorevo.NewMockMonorevoFetcher(ctrl)
	// EXPECTはctrl#Finishが呼び出される前に FetchAllを呼び出さなければエラーになる
	mock_fetcher.EXPECT().FetchAll().Return(mock_results, nil)

	// UseCase戻り値
	results := []proposition_fetch_case.FetchedPropositionDto{}
	for _, v := range mock_results {
		results = append(results, proposition_fetch_case.FetchedPropositionDto{
			WorkedNumber: v.WorkedNumber,
			DET:          v.DET,
			DeliveryDate: v.DeliveryDate,
		})
	}

	tests := []struct {
		name    string
		m       *proposition_fetch_case.PropositionFetchingUseCase
		want    []proposition_fetch_case.FetchedPropositionDto
		wantErr bool
	}{
		{
			name: "正常系_UseCaseを実行するとモックが実行されること",
			m: proposition_fetch_case.NewPropositionFetchingUseCase(
				logger.Sugar(),
				mock_fetcher,
			),
			want:    results,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("PropositionFetchingUseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("len(PropositionFetchingUseCase.Execute()) = %v, want %v", len(got), len(tt.want))
			}
			for i := 0; i < len(got); i++ {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("PropositionFetchingUseCase.Execute()[%v] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
