package jobbook_fetch_case_test

import (
	"SynchronizeMonorevoDeliveryDates/domain/orderdb"
	"SynchronizeMonorevoDeliveryDates/domain/orderdb/mock_orderdb"
	"SynchronizeMonorevoDeliveryDates/usecase/jobbook_fetch_case"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestJobBookFetchingUseCase_Execute(t *testing.T) {
	// logger生成
	logger, _ := zap.NewDevelopment()

	// FetchAll戻り値
	mock_results := []orderdb.JobBook{}
	mock_job := orderdb.TestJobBookCreate()
	mock_results = append(mock_results, *mock_job)

	// モックコントローラーの生成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 受注管理DB DIオブジェクト生成
	mock_fetcher := mock_orderdb.NewMockJobBookFetcher(ctrl)
	// EXPECTはctrl#Finishが呼び出される前に FetchAllを呼び出さなければエラーになる
	mock_fetcher.EXPECT().FetchAll().Return(mock_results, nil)

	// UseCase戻り値
	results := []jobbook_fetch_case.JobBookDto{}
	for _, v := range mock_results {
		results = append(results, jobbook_fetch_case.JobBookDto{
			v.WorkedNumber,
			v.DeliveryDate,
		})
	}

	tests := []struct {
		name    string
		m       *jobbook_fetch_case.JobBookFetchingUseCase
		want    []jobbook_fetch_case.JobBookDto
		wantErr bool
	}{
		{
			name: "正常系_UseCaseを実行するとモックが実行されること",
			m: jobbook_fetch_case.NewJobBookFetchingUseCase(
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
				t.Errorf("JobBookFetchingUseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JobBookFetchingUseCase.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
