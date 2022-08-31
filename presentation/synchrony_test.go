package presentation_test

import (
	"SynchronizeMonorevoDeliveryDates/presentation"
	"SynchronizeMonorevoDeliveryDates/usecase/appsetting_obtain_case/mock_appsetting_obtain_case"
	"SynchronizeMonorevoDeliveryDates/usecase/difference_extract_case"
	"SynchronizeMonorevoDeliveryDates/usecase/difference_extract_case/mock_difference_extract_case"
	"SynchronizeMonorevoDeliveryDates/usecase/jobbook_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/jobbook_fetch_case/mock_jobbook_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_fetch_case/mock_proposition_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_post_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_post_case/mock_proposition_post_case"
	"SynchronizeMonorevoDeliveryDates/usecase/reportsetting_obtain_case/mock_reportsetting_obtain_case"
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
	mock_appCfgLoader := mock_appsetting_obtain_case.NewMockSettingLoader(ctrl)

	// reportsettingモック作成
	mock_repCfgLoader := mock_reportsetting_obtain_case.NewMockSettingLoader(ctrl)

	// webモック作成
	resWebFetches, mock_webFetcher := makeMockWebFetcher(ctrl)

	// DBモック作成
	resDbFetches, mock_dbFetcher := makeMockDbFetcher(ctrl)

	// 差分モック作成
	mock_diff := makeMockDifferent(resWebFetches, resDbFetches, ctrl)

	// 更新モック作成
	mock_post := makeMockWebPoster(resWebFetches, resDbFetches, ctrl)

	tests := []struct {
		name    string
		m       *presentation.SynchronizingDeliveryDate
		wantErr bool
	}{
		{
			name: "正常系_UseCaseを実行するとモックが実行されること",
			m: presentation.NewSynchronizingDeliveryDate(
				logger.Sugar(),
				mock_webFetcher,
				mock_dbFetcher,
				mock_diff,
				mock_post,
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

func makeMockWebPoster(resWebFetches []proposition_fetch_case.FetchedPropositionDto, resDbFetches []jobbook_fetch_case.JobBookDto, ctrl *gomock.Controller) *mock_proposition_post_case.MockPostingExecutor {
	postPrams := []proposition_post_case.PostingPropositionPram{}
	for i := 0; i < len(resWebFetches); i++ {
		postPrams = append(postPrams,
			proposition_post_case.PostingPropositionPram{
				WorkedNumber:        resWebFetches[i].WorkedNumber,
				Det:                 resWebFetches[i].Det,
				DeliveryDate:        resWebFetches[i].DeliveryDate,
				UpdatedDeliveryDate: resDbFetches[i].DeliveryDate,
			},
		)
	}
	resPosts := []proposition_post_case.PostedPropositionDto{}
	for i := 0; i < len(resWebFetches); i++ {
		resPosts = append(resPosts,
			proposition_post_case.PostedPropositionDto{
				WorkedNumber:        resWebFetches[i].WorkedNumber,
				Det:                 resWebFetches[i].Det,
				Successful:          true,
				DeliveryDate:        resWebFetches[i].DeliveryDate,
				UpdatedDeliveryDate: resDbFetches[i].DeliveryDate,
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
				Det:          pro.Det,
				DeliveryDate: pro.DeliveryDate,
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
			Det:                 resWebFetches[0].Det,
			DeliveryDate:        resWebFetches[0].DeliveryDate,
			UpdatedDeliveryDate: resDbFetches[0].DeliveryDate,
		},
		{
			WorkedNumber:        resWebFetches[1].WorkedNumber,
			Det:                 resWebFetches[1].Det,
			DeliveryDate:        resWebFetches[1].DeliveryDate,
			UpdatedDeliveryDate: resDbFetches[1].DeliveryDate,
		},
		{
			WorkedNumber:        resWebFetches[2].WorkedNumber,
			Det:                 resWebFetches[2].Det,
			DeliveryDate:        resWebFetches[2].DeliveryDate,
			UpdatedDeliveryDate: resDbFetches[2].DeliveryDate,
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
			Det:          "1",
			DeliveryDate: time.Now().AddDate(0, 0, -5),
		},
		{
			WorkedNumber: "88A-1234",
			Det:          "1",
			DeliveryDate: time.Now().AddDate(0, 0, -5),
		},
		{
			WorkedNumber: "77A-1234",
			Det:          "1",
			DeliveryDate: time.Now().AddDate(0, 0, -5),
		},
	}
	mock_webFetcher := mock_proposition_fetch_case.NewMockFetchingExecutor(ctrl)
	mock_webFetcher.EXPECT().Execute().Return(resWebFetches, nil)
	return resWebFetches, mock_webFetcher
}
