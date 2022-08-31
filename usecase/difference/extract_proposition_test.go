package difference_test

import (
	"SynchronizeMonorevoDeliveryDates/domain/compare/mock_compare"
	"SynchronizeMonorevoDeliveryDates/usecase/difference"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestPropositionExtractingUseCase_Execute(t *testing.T) {
	// logger生成
	logger, _ := zap.NewDevelopment()

	// モックコントローラーの生成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 納期差分抽出DIオブジェクト生成
	mock_diff := mock_compare.NewMockExtractor(ctrl)
	// EXPECTはctrl#Finishが呼び出される前に FetchAllを呼び出さなければエラーになる
	mock_diff.EXPECT().ExtractForDeliveryDate(gomock.Any(), gomock.Any()).Return(nil)

	type args struct {
		s difference.DifferenceSourcePram
	}
	tests := []struct {
		name string
		m    *difference.PropositionExtractingUseCase
		args args
		want []difference.DifferentPropositionDto
	}{
		{
			name: "正常系_UseCaseを実行するとモックが実行されること",
			m:    difference.NewExtractingPropositionUseCase(logger.Sugar(), mock_diff),
			args: args{},
			want: []difference.DifferentPropositionDto{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Execute(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PropositionExtractingUseCase.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
