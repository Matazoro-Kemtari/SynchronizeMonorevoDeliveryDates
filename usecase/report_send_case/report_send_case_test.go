package report_send_case_test

import (
	"SynchronizeMonorevoDeliveryDates/domain/report/mock_report"
	"SynchronizeMonorevoDeliveryDates/usecase/report_send_case"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestSendingReportUseCase_Execute(t *testing.T) {
	// logger生成
	logger, _ := zap.NewDevelopment()

	// モックコントローラーの生成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// レポート送信DIオブジェクト生成
	reportPram := report_send_case.TestReportPramCreate()
	mock_results := time.Now().Format("2006/01/02")
	mock_sender := mock_report.NewMockSender(ctrl)
	// EXPECTはctrl#Finishが呼び出される前に FetchAllを呼び出さなければエラーになる
	mock_sender.EXPECT().Send(
		report_send_case.ConvertToEmailAddresses(reportPram.Tos),
		report_send_case.ConvertToEmailAddresses(reportPram.CCs),
		report_send_case.ConvertToEmailAddresses(reportPram.BCCs),
		*reportPram.From.ConvertToEmailAddress(),
		*reportPram.ReplyTo.ConvertToEmailAddress(),
		reportPram.Subject,
		report_send_case.ConvertToEditedProposition(reportPram.EditedPropositions),
		reportPram.PrefixReport,
		reportPram.SuffixReport,
		reportPram.Replacements,
	).Return(mock_results, nil)

	type args struct {
		r report_send_case.ReportPram
	}
	tests := []struct {
		name    string
		s       *report_send_case.SendingReportUseCase
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "正常系_UseCaseを実行するとモックが実行されること",
			s: report_send_case.NewSendingReportUseCase(
				logger.Sugar(),
				mock_sender,
			),
			args: args{
				r: *reportPram,
			},
			want:    mock_results,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Execute(tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendingReportUseCase.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SendingReportUseCase.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
