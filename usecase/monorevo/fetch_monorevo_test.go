package monorevo

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/domain/monorevo/mock_monorevo"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestFetchPropositionTable_Execute(t *testing.T) {
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
	mock_fetcher := mock_monorevo.NewMockFetcher(ctrl)
	// EXPECTはctrl#Finishが呼び出される前に FetchAllを呼び出さなければエラーになる
	mock_fetcher.EXPECT().FetchAll().Return(mock_results, nil)

	// UseCase戻り値
	results := []PropositionDto{}
	for _, v := range mock_results {
		results = append(results, PropositionDto{
			v.WorkedNumber,
			v.DeliveryDate,
		})
	}

	tests := []struct {
		name    string
		m       *FetchPropositionTable
		want    []PropositionDto
		wantErr bool
	}{
		{
			name: "正常系_UseCaseを実行するとモックが実行されること",
			m: NewFetchPropositionTable(
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
				t.Errorf("FetchPropositionTable.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchPropositionTable.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
