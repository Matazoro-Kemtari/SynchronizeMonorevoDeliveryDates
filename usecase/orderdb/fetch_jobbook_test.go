package orderdb

import (
	"SynchronizeMonorevoDeliveryDates/domain/orderdb"
	"SynchronizeMonorevoDeliveryDates/domain/orderdb/mock_orderdb"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestFetchJobBookTable_Execute(t *testing.T) {
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
	results := []JobBookDto{}
	for _, v := range mock_results {
		results = append(results, JobBookDto{
			v.WorkedNumber,
			v.DeliveryDate,
		})
	}

	tests := []struct {
		name    string
		m       *FetchJobBookTable
		want    []JobBookDto
		wantErr bool
	}{
		{
			name: "正常系_UseCaseを実行するとFetcherが実行されること",
			m: NewFetchJobBookTable(
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
				t.Errorf("FetchJobBookTable.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchJobBookTable.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}