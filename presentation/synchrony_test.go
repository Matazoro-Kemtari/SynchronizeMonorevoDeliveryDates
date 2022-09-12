package presentation_test

import (
	"SynchronizeMonorevoDeliveryDates/presentation"
	"SynchronizeMonorevoDeliveryDates/usecase/appsetting_obtain_case"
	"SynchronizeMonorevoDeliveryDates/usecase/difference_extract_case"
	"SynchronizeMonorevoDeliveryDates/usecase/difference_extract_case/mock_difference_extract_case"
	"SynchronizeMonorevoDeliveryDates/usecase/jobbook_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/jobbook_fetch_case/mock_jobbook_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_fetch_case/mock_proposition_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_post_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_post_case/mock_proposition_post_case"
	"SynchronizeMonorevoDeliveryDates/usecase/report_send_case"
	"SynchronizeMonorevoDeliveryDates/usecase/report_send_case/mock_report_send_case"
	"SynchronizeMonorevoDeliveryDates/usecase/reportsetting_obtain_case"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

func TestSynchronizingDeliveryDate_Synchronize(t *testing.T) {
	// logger生成
	logger, _ := zap.NewDevelopment()

	// モックコントローラーの生成
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// appsettingモック作成
	appSetting := appsetting_obtain_case.TestAppSettingDtoCreate()

	// reportsettingモック作成
	reportSetting := reportsetting_obtain_case.TestReportSettingDtoCreate()

	// webモック作成
	resWebFetches, mock_webFetcher := makeMockWebFetcher(ctrl)

	// DBモック作成
	resDbFetches, mock_dbFetcher := makeMockDbFetcher(ctrl)

	// 差分モック作成
	mock_diff := makeMockDifferent(resWebFetches, resDbFetches, ctrl)

	// 更新モック作成
	mock_post := makeMockWebPoster(resWebFetches, resDbFetches, ctrl)

	// 報告モック作成
	mock_report := makeMockReportSender(reportSetting, resWebFetches, resDbFetches, ctrl)

	tests := []struct {
		name    string
		m       *presentation.SynchronizingDeliveryDate
		wantErr bool
	}{
		{
			name: "正常系_controllerを実行するとモックが実行されること",
			m: presentation.NewSynchronizingDeliveryDate(
				logger.Sugar(),
				appSetting,
				reportSetting,
				mock_webFetcher,
				mock_dbFetcher,
				mock_diff,
				mock_post,
				mock_report,
			),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Synchronize(); (err != nil) != tt.wantErr {
				t.Errorf("SynchronizingDeliveryDate.Synchronize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func convertToEmailAddresses(cfgs []reportsetting_obtain_case.MailAddressDto) []report_send_case.EmailAddressPram {
	var mails []report_send_case.EmailAddressPram
	for _, v := range cfgs {
		mails = append(mails, convertToEmailAddress(v))
	}
	return mails
}

func convertToEmailAddress(m reportsetting_obtain_case.MailAddressDto) report_send_case.EmailAddressPram {
	return report_send_case.EmailAddressPram{
		Name:    m.Name,
		Address: m.Email,
	}
}

func makeMockReportSender(
	setting *reportsetting_obtain_case.ReportSettingDto,
	resWebFetches []proposition_fetch_case.FetchedPropositionDto,
	resDbFetches []jobbook_fetch_case.JobBookDto,
	ctrl *gomock.Controller,
) *mock_report_send_case.MockExecutor {
	var editedPropositions []report_send_case.EditedPropositionPram
	for i := 0; i < len(resWebFetches); i++ {
		editedPropositions = append(editedPropositions,
			*report_send_case.TestEditedPropositionPramCreate(
				report_send_case.OptWorkedNumber(resWebFetches[i].WorkedNumber),
				report_send_case.OptDET(resWebFetches[i].DET),
				report_send_case.OptSuccessful(true),
				report_send_case.OptReason(""),
				report_send_case.OptDeliveryDate(resWebFetches[i].DeliveryDate),
				report_send_case.OptUpdatedDeliveryDate(resDbFetches[i].DeliveryDate),
				report_send_case.OptCode(resWebFetches[i].Code),
			),
		)
	}

	reportRes := time.Now().Format("2006/01/02")

	reportPram := *report_send_case.TestReportPramCreate(
		report_send_case.OptTos(convertToEmailAddresses(setting.RecipientAddresses)),
		report_send_case.OptCCs(convertToEmailAddresses(setting.CCAddresses)),
		report_send_case.OptBCCs(convertToEmailAddresses(setting.BCCAddresses)),
		report_send_case.OptFrom(convertToEmailAddress(setting.SenderAddress)),
		report_send_case.OptReplyTo(convertToEmailAddress(setting.ReplyToAddress)),
		report_send_case.OptSubject(setting.Subject),
		report_send_case.OptEditedPropositions(editedPropositions),
		report_send_case.OptPrefixReport(setting.PrefixReport),
		report_send_case.OptSuffixReport(setting.SuffixReport),
		report_send_case.OptReplacements(map[string]string{"count": fmt.Sprint(len(editedPropositions))}),
	)
	mock_report := mock_report_send_case.NewMockExecutor(ctrl)
	mock_report.EXPECT().Execute(reportPram).Return(reportRes, nil)

	return mock_report
}

func makeMockWebPoster(resWebFetches []proposition_fetch_case.FetchedPropositionDto, resDbFetches []jobbook_fetch_case.JobBookDto, ctrl *gomock.Controller) *mock_proposition_post_case.MockPostingExecutor {
	postPrams := []proposition_post_case.PostingPropositionPram{}
	for i := 0; i < len(resWebFetches); i++ {
		postPrams = append(postPrams,
			proposition_post_case.PostingPropositionPram{
				WorkedNumber:        resWebFetches[i].WorkedNumber,
				DET:                 resWebFetches[i].DET,
				DeliveryDate:        resWebFetches[i].DeliveryDate,
				UpdatedDeliveryDate: resDbFetches[i].DeliveryDate,
				Code:                resWebFetches[i].Code,
			},
		)
	}
	resPosts := []proposition_post_case.PostedPropositionDto{}
	for i := 0; i < len(resWebFetches); i++ {
		resPosts = append(resPosts,
			proposition_post_case.PostedPropositionDto{
				WorkedNumber:        resWebFetches[i].WorkedNumber,
				DET:                 resWebFetches[i].DET,
				Successful:          true,
				Reason:              "",
				DeliveryDate:        resWebFetches[i].DeliveryDate,
				UpdatedDeliveryDate: resDbFetches[i].DeliveryDate,
				Code:                resWebFetches[i].Code,
			},
		)
	}
	mock_post := mock_proposition_post_case.NewMockPostingExecutor(ctrl)
	mock_post.EXPECT().Execute(postPrams).Return(resPosts, nil)
	return mock_post
}

func makeMockDifferent(resWebFetches []proposition_fetch_case.FetchedPropositionDto, resDbFetches []jobbook_fetch_case.JobBookDto, ctrl *gomock.Controller) *mock_difference_extract_case.MockExecutor {
	diffPropositions := []difference_extract_case.PropositionPram{}
	for _, pro := range resWebFetches {
		diffPropositions = append(diffPropositions,
			difference_extract_case.PropositionPram{
				WorkedNumber: pro.WorkedNumber,
				DET:          pro.DET,
				DeliveryDate: pro.DeliveryDate,
				Code:         pro.Code,
			},
		)
	}
	diffJobBooks := []difference_extract_case.JobBookPram{}
	for _, job := range resDbFetches {
		diffJobBooks = append(diffJobBooks,
			difference_extract_case.JobBookPram{
				WorkedNumber: job.WorkedNumber,
				DeliveryDate: job.DeliveryDate,
			},
		)
	}

	diffPram := difference_extract_case.DifferenceSourcePram{
		JobBooks:     diffJobBooks,
		Propositions: diffPropositions,
	}
	resDiffs := []difference_extract_case.DifferentPropositionDto{
		{
			WorkedNumber:        resWebFetches[0].WorkedNumber,
			DET:                 resWebFetches[0].DET,
			DeliveryDate:        resWebFetches[0].DeliveryDate,
			UpdatedDeliveryDate: resDbFetches[0].DeliveryDate,
			Code:                resWebFetches[0].Code,
		},
		{
			WorkedNumber:        resWebFetches[1].WorkedNumber,
			DET:                 resWebFetches[1].DET,
			DeliveryDate:        resWebFetches[1].DeliveryDate,
			UpdatedDeliveryDate: resDbFetches[1].DeliveryDate,
			Code:                resWebFetches[1].Code,
		},
		{
			WorkedNumber:        resWebFetches[2].WorkedNumber,
			DET:                 resWebFetches[2].DET,
			DeliveryDate:        resWebFetches[2].DeliveryDate,
			UpdatedDeliveryDate: resDbFetches[2].DeliveryDate,
			Code:                resWebFetches[2].Code,
		},
	}
	mock_diff := mock_difference_extract_case.NewMockExecutor(ctrl)
	mock_diff.EXPECT().Execute(diffPram).Return(resDiffs)
	return mock_diff
}

func makeMockDbFetcher(ctrl *gomock.Controller) ([]jobbook_fetch_case.JobBookDto, *mock_jobbook_fetch_case.MockExecutor) {
	resDbFetches := []jobbook_fetch_case.JobBookDto{
		{
			WorkedNumber: "99A-1234",
			DeliveryDate: time.Now(),
		},
		{
			WorkedNumber: "88A-1234",
			DeliveryDate: time.Now(),
		},
		{
			WorkedNumber: "77A-1234",
			DeliveryDate: time.Now(),
		},
		{
			WorkedNumber: "66A-1234",
			DeliveryDate: time.Now(),
		},
	}
	mock_dbFetcher := mock_jobbook_fetch_case.NewMockExecutor(ctrl)
	mock_dbFetcher.EXPECT().Execute().Return(resDbFetches, nil)
	return resDbFetches, mock_dbFetcher
}

func makeMockWebFetcher(ctrl *gomock.Controller) ([]proposition_fetch_case.FetchedPropositionDto, *mock_proposition_fetch_case.MockFetchingExecutor) {
	resWebFetches := []proposition_fetch_case.FetchedPropositionDto{
		{
			WorkedNumber: "99A-1234",
			DET:          "1",
			DeliveryDate: time.Now().AddDate(0, 0, -5),
			Code:         "11A",
		},
		{
			WorkedNumber: "88A-1234",
			DET:          "1",
			DeliveryDate: time.Now().AddDate(0, 0, -5),
			Code:         "22B",
		},
		{
			WorkedNumber: "77A-1234",
			DET:          "1",
			DeliveryDate: time.Now().AddDate(0, 0, -5),
			Code:         "33C",
		},
	}
	mock_webFetcher := mock_proposition_fetch_case.NewMockFetchingExecutor(ctrl)
	mock_webFetcher.EXPECT().Execute().Return(resWebFetches, nil)
	return resWebFetches, mock_webFetcher
}
