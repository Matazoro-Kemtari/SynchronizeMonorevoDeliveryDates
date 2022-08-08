package monorevo

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/domain/monorevo/mock_monorevo"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestFetchMonoRevoPropositionTable_Execute(t *testing.T) {
	// logger生成
	logger, _ := zap.NewDevelopment()

	// FetchAll戻り値
	moc_results := []monorevo.Proposition{}
	moc_pro := monorevo.TestPropositionCreate()
	moc_results = append(moc_results, *moc_pro)

	// モックコントローラーの生成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// ものレボDIオブジェクト生成
	mock_fetcher := mock_monorevo.NewMockFetcher(ctrl)
	// EXPECTはctrl#Finishが呼び出される前に FetchAllを呼び出さなければエラーになる
	mock_fetcher.EXPECT().FetchAll().Return(moc_results, nil)

	// UseCase戻り値
	result := []PropositionDto{}
	for _, v := range moc_results {
		result = append(result, PropositionDto{
			v.WorkedNumber,
			v.DeliveryDate,
		})
	}

	tests := []struct {
		name    string
		m       *FetchMonoRevoPropositionTable
		want    []PropositionDto
		wantErr bool
	}{
		{
			name: "正常系_UseCaseを実行するとFetcherが実行されること",
			m: NewFetchMonoRevoPropositionTable(
				logger.Sugar(),
				mock_fetcher,
			),
			want:    result,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchMonoRevoPropositionTable.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchMonoRevoPropositionTable.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
