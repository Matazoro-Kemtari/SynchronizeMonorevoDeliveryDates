package presentation

import (
	"SynchronizeMonorevoDeliveryDates/usecase/difference"
	"SynchronizeMonorevoDeliveryDates/usecase/difference/mock_difference"
	"SynchronizeMonorevoDeliveryDates/usecase/monorevo"
	"SynchronizeMonorevoDeliveryDates/usecase/monorevo/mock_monorevo"
	"SynchronizeMonorevoDeliveryDates/usecase/orderdb"
	"SynchronizeMonorevoDeliveryDates/usecase/orderdb/mock_orderdb"
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
		m       *SynchronizingDeliveryDate
		wantErr bool
	}{
		{
			name: "正常系_UseCaseを実行するとモックが実行されること",
			m: NewSynchronizingDeliveryDate(
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

func makeMockWebPoster(resWebFetches []monorevo.FetchedPropositionDto, resDbFetches []orderdb.JobBookDto, ctrl *gomock.Controller) *mock_monorevo.MockPoster {
	resPosts := []monorevo.PostedPropositionDto{
		{
			WorkedNumber:        resWebFetches[0].WorkedNumber,
			Det:                 resWebFetches[0].Det,
			Successful:          true,
			DeliveryDate:        resWebFetches[0].DeliveryDate,
			UpdatedDeliveryDate: resDbFetches[0].DeliveryDate,
		},
		{
			WorkedNumber:        resWebFetches[1].WorkedNumber,
			Det:                 resWebFetches[1].Det,
			Successful:          true,
			DeliveryDate:        resWebFetches[1].DeliveryDate,
			UpdatedDeliveryDate: resDbFetches[1].DeliveryDate,
		},
		{
			WorkedNumber:        resWebFetches[2].WorkedNumber,
			Det:                 resWebFetches[2].Det,
			Successful:          true,
			DeliveryDate:        resWebFetches[2].DeliveryDate,
			UpdatedDeliveryDate: resDbFetches[2].DeliveryDate,
		},
	}
	mock_post := mock_monorevo.NewMockPoster(ctrl)
	mock_post.EXPECT().PostRange(gomock.Any()).Return(resPosts, nil)
	return mock_post
}

func makeMockDifferent(resWebFetches []monorevo.FetchedPropositionDto, resDbFetches []orderdb.JobBookDto, ctrl *gomock.Controller) *mock_difference.MockExtractor {
	resDiffs := []difference.DifferentPropositionDto{
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
	mock_diff := mock_difference.NewMockExtractor(ctrl)
	mock_diff.EXPECT().Extract(gomock.Any()).Return(resDiffs)
	return mock_diff
}

func makeMockDbFetcher(ctrl *gomock.Controller) ([]orderdb.JobBookDto, *mock_orderdb.MockFetcher) {
	resDbFetches := []orderdb.JobBookDto{
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
	mock_dbFetcher := mock_orderdb.NewMockFetcher(ctrl)
	mock_dbFetcher.EXPECT().Fetch().Return(resDbFetches, nil)
	return resDbFetches, mock_dbFetcher
}

func makeMockWebFetcher(ctrl *gomock.Controller) ([]monorevo.FetchedPropositionDto, *mock_monorevo.MockFetcher) {
	resWebFetches := []monorevo.FetchedPropositionDto{
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
	mock_webFetcher := mock_monorevo.NewMockFetcher(ctrl)
	mock_webFetcher.EXPECT().Fetch().Return(resWebFetches, nil)
	return resWebFetches, mock_webFetcher
}
