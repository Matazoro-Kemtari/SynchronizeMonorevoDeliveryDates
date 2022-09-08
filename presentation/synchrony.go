package presentation

import (
	"SynchronizeMonorevoDeliveryDates/usecase/appsetting_obtain_case"
	"SynchronizeMonorevoDeliveryDates/usecase/difference_extract_case"
	"SynchronizeMonorevoDeliveryDates/usecase/jobbook_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_fetch_case"
	"SynchronizeMonorevoDeliveryDates/usecase/proposition_post_case"
	"SynchronizeMonorevoDeliveryDates/usecase/report_send_case"
	"SynchronizeMonorevoDeliveryDates/usecase/reportsetting_obtain_case"
	"strconv"

	"go.uber.org/zap"
)

type SynchronizingDeliveryDate struct {
	sugar         *zap.SugaredLogger
	appSetting    *appsetting_obtain_case.AppSettingDto
	reportSetting *reportsetting_obtain_case.ReportSettingDto
	webFetcher    proposition_fetch_case.FetchingExecutor
	dbFetcher     jobbook_fetch_case.Executor
	extractor     difference_extract_case.Executor
	webPoster     proposition_post_case.PostingExecutor
	reportSender  report_send_case.Executor
}

func NewSynchronizingDeliveryDate(
	sugar *zap.SugaredLogger,
	appSetting *appsetting_obtain_case.AppSettingDto,
	reportSetting *reportsetting_obtain_case.ReportSettingDto,
	webFetcher proposition_fetch_case.FetchingExecutor,
	dbFetcher jobbook_fetch_case.Executor,
	extractor difference_extract_case.Executor,
	webPoster proposition_post_case.PostingExecutor,
	reportSender report_send_case.Executor,
) *SynchronizingDeliveryDate {
	return &SynchronizingDeliveryDate{
		sugar:         sugar,
		appSetting:    appSetting,
		reportSetting: reportSetting,
		webFetcher:    webFetcher,
		dbFetcher:     dbFetcher,
		extractor:     extractor,
		webPoster:     webPoster,
		reportSender:  reportSender,
	}
}

func (m *SynchronizingDeliveryDate) Synchronize() error {
	m.sugar.Info("ものレボから案件一覧を取得する")
	propositions, err := m.webFetcher.Execute()
	if err != nil {
		m.sugar.Fatal("ものレボから案件一覧を取得で失敗しました", err)
	}
	m.sugar.Debug("propositions", propositions)

	m.sugar.Info("受注管理DBから作業情報を取得する")
	jobBooks, err := m.dbFetcher.Execute()
	if err != nil {
		m.sugar.Fatal("受注管理DBから作業情報を取得で失敗しました", err)
	}
	m.sugar.Debug("jobBooks", jobBooks)

	// 詰め替え
	diffPram := convertToDifferencePram(propositions, jobBooks)

	m.sugar.Info("比較差分を算出する")
	diff := m.extractor.Execute(diffPram)
	m.sugar.Debug("diff", diff)

	var posted []proposition_post_case.PostedPropositionDto
	if diff != nil {
		m.sugar.Info("ものレボへ案件一覧を送信する")
		posting := convertToPostPrams(diff)
		var err error
		posted, err = m.webPoster.Execute(posting)
		if err != nil {
			m.sugar.Fatal("ものレボへ案件一覧を送信で失敗しました", err)
		}
		m.sugar.Debug("posted", posted)
	}

	// 詰め替え
	reportPram := m.convertToReportPram(posted)

	sent, err := m.reportSender.Execute(reportPram)
	if err != nil {
		m.sugar.Fatal("結果報告で失敗しました", err)
	}
	m.sugar.Debug("sent", sent)

	return nil
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

func convertToEditedPropositionPrams(p []proposition_post_case.PostedPropositionDto) []report_send_case.EditedPropositionPram {
	var params []report_send_case.EditedPropositionPram
	for _, v := range p {
		params = append(params, report_send_case.EditedPropositionPram{
			WorkedNumber:        v.WorkedNumber,
			DET:                 v.DET,
			Successful:          v.Successful,
			DeliveryDate:        v.DeliveryDate,
			UpdatedDeliveryDate: v.UpdatedDeliveryDate,
		})
	}
	return params
}

func (m *SynchronizingDeliveryDate) convertToReportPram(p []proposition_post_case.PostedPropositionDto) report_send_case.ReportPram {
	return report_send_case.ReportPram{
		Tos:                convertToEmailAddresses(m.reportSetting.RecipientAddresses),
		CCs:                convertToEmailAddresses(m.reportSetting.CCAddresses),
		BCCs:               convertToEmailAddresses(m.reportSetting.BCCAddresses),
		From:               convertToEmailAddress(m.reportSetting.SenderAddress),
		ReplyTo:            convertToEmailAddress(m.reportSetting.ReplyToAddress),
		Subject:            m.reportSetting.Subject,
		EditedPropositions: convertToEditedPropositionPrams(p),
		PrefixReport:       m.reportSetting.PrefixReport,
		SuffixReport:       m.reportSetting.SuffixReport,
		Replacements:       map[string]string{"count": strconv.Itoa(len(p))},
	}
}

// ものレボ更新パラメータへ詰め替え
func convertToPostPrams(diff []difference_extract_case.DifferentPropositionDto) []proposition_post_case.PostingPropositionPram {
	posting := []proposition_post_case.PostingPropositionPram{}
	for _, v := range diff {
		posting = append(posting,
			proposition_post_case.PostingPropositionPram{
				WorkedNumber:        v.WorkedNumber,
				DET:                 v.DET,
				DeliveryDate:        v.DeliveryDate,
				UpdatedDeliveryDate: v.UpdatedDeliveryDate,
			})
	}
	return posting
}

// 差分抽出パラメータへ詰め替え
func convertToDifferencePram(propositions []proposition_fetch_case.FetchedPropositionDto, jobBooks []jobbook_fetch_case.JobBookDto) difference_extract_case.DifferenceSourcePram {
	diffPropositions := []difference_extract_case.PropositionPram{}
	for _, pro := range propositions {
		diffPropositions = append(diffPropositions,
			difference_extract_case.PropositionPram{
				WorkedNumber: pro.WorkedNumber,
				DET:          pro.DET,
				DeliveryDate: pro.DeliveryDate,
			},
		)
	}
	diffJobBooks := []difference_extract_case.JobBookPram{}
	for _, job := range jobBooks {
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
	return diffPram
}
