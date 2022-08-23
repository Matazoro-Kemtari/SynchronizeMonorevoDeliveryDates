package monorevo

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"
	"SynchronizeMonorevoDeliveryDates/domain/monorevo/mock_monorevo"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestPropositionTable_Fetch(t *testing.T) {
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
	results := []FetchedPropositionDto{}
	for _, v := range mock_results {
		results = append(results, FetchedPropositionDto{
			WorkedNumber: v.WorkedNumber,
			Det:          v.Det,
			DeliveryDate: v.DeliveryDate,
		})
	}

	tests := []struct {
		name    string
		m       *PropositionTable
		want    []FetchedPropositionDto
		wantErr bool
	}{
		{
			name: "正常系_UseCaseを実行するとモックが実行されること",
			m: NewPropositionTable(
				logger.Sugar(),
				mock_fetcher,
				nil,
			),
			want:    results,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.m.Fetch()
			if (err != nil) != tt.wantErr {
				t.Errorf("PropositionTable.Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("len(PropositionTable.Fetch()) = %v, want %v", len(got), len(tt.want))
			}
			for i := 0; i < len(got); i++ {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("PropositionTable.Fetch()[%v] = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
