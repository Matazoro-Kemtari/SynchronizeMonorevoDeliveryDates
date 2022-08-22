package presentation

import (
	"SynchronizeMonorevoDeliveryDates/usecase/difference"
	"SynchronizeMonorevoDeliveryDates/usecase/monorevo"
	"SynchronizeMonorevoDeliveryDates/usecase/orderdb"

	"go.uber.org/zap"
)

type SynchronizingDeliveryDate struct {
	sugar      *zap.SugaredLogger
	webFetcher monorevo.Fetcher
	dbFetcher  orderdb.Fetcher
	extractor  difference.Extractor
	webPoster  monorevo.Poster
}

func NewSynchronizingDeliveryDate(
	sugar *zap.SugaredLogger,
	webFetcher monorevo.Fetcher,
	dbFetcher orderdb.Fetcher,
	extractor difference.Extractor,
	webPoster monorevo.Poster,

) *SynchronizingDeliveryDate {
	return &SynchronizingDeliveryDate{
		sugar:      sugar,
		webFetcher: webFetcher,
		dbFetcher:  dbFetcher,
		extractor:  extractor,
		webPoster:  webPoster,
	}
}

func (m *SynchronizingDeliveryDate) Synchronize() error {
	m.sugar.Info("ものレボから案件一覧を取得する")
	propositions, err := m.webFetcher.Fetch()
	if err != nil {
		m.sugar.Fatal("ものレボから案件一覧を取得で失敗しました", err)
	}
	m.sugar.Debug("propositions", propositions)

	m.sugar.Info("受注管理DBから作業情報を取得する")
	jobBooks, err := m.dbFetcher.Fetch()
	if err != nil {
		m.sugar.Fatal("受注管理DBから作業情報を取得で失敗しました", err)
	}
	m.sugar.Debug("jobBooks", jobBooks)

	// 詰め替え
	diffPram := m.convertDifferencePram(propositions, jobBooks)

	m.sugar.Info("比較差分を算出する")
	diff := m.extractor.Extract(diffPram)
	m.sugar.Debug("diff", diff)

	m.sugar.Info("ものレボへ案件一覧を送信する")
	posting := []monorevo.PostingPropositionPram{}
	for _, v := range diff {
		posting = append(posting,
			monorevo.PostingPropositionPram{
				WorkedNumber:        v.WorkedNumber,
				Det:                 v.Det,
				DeliveryDate:        v.DeliveryDate,
				UpdatedDeliveryDate: v.UpdatedDeliveryDate,
			})
	}
	posted, err := m.webPoster.PostRange(posting)
	if err != nil {
		m.sugar.Fatal("ものレボへ案件一覧を送信で失敗しました", err)
	}
	m.sugar.Debug("posted", posted)

	return nil
}

// 差分抽出パラメータへ詰め替え
func (*SynchronizingDeliveryDate) convertDifferencePram(propositions []monorevo.FetchedPropositionDto, jobBooks []orderdb.JobBookDto) difference.DifferenceSourcePram {
	diffPropositions := []difference.PropositionPram{}
	for _, pro := range propositions {
		diffPropositions = append(diffPropositions,
			difference.PropositionPram{
				WorkedNumber: pro.WorkedNumber,
				Det:          pro.Det,
				DeliveryDate: pro.DeliveryDate,
			},
		)
	}
	diffJobBooks := []difference.JobBookPram{}
	for _, job := range jobBooks {
		diffJobBooks = append(diffJobBooks,
			difference.JobBookPram{
				WorkedNumber: job.WorkedNumber,
				DeliveryDate: job.DeliveryDate,
			},
		)
	}
	diffPram := difference.DifferenceSourcePram{
		JobBooks:     diffJobBooks,
		Propositions: diffPropositions,
	}
	return diffPram
}
