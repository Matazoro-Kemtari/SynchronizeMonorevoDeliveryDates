package report_send_case

import (
	"SynchronizeMonorevoDeliveryDates/domain/report"
	"time"

	"go.uber.org/zap"
)

type EmailAddressPram struct {
	Name    string
	Address string
}

func (e *EmailAddressPram) ConvertToEmailAddress() *report.EmailAddress {
	return &report.EmailAddress{
		Name:    e.Name,
		Address: e.Address,
	}
}

type EmailAddressOptions struct {
	Name    string
	Address string
}

type EmailAddressOption func(*EmailAddressOptions)

func OptName(v string) EmailAddressOption {
	return func(opts *EmailAddressOptions) {
		opts.Name = v
	}
}

func OptAddress(v string) EmailAddressOption {
	return func(opts *EmailAddressOptions) {
		opts.Address = v
	}
}

func TestEmailAddressPramCreate(options ...EmailAddressOption) *EmailAddressPram {
	// デフォルト値
	opts := &EmailAddressOptions{
		Name:    "サンプルさん",
		Address: "foo@example.com",
	}

	for _, option := range options {
		option(opts)
	}

	return &EmailAddressPram{
		Name:    opts.Name,
		Address: opts.Address,
	}
}

type EditedPropositionPram struct {
	WorkedNumber        string
	DET                 string
	Successful          bool
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

type EditedPropositionOptions struct {
	WorkedNumber        string
	DET                 string
	Successful          bool
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

type EditedPropositionOption func(*EditedPropositionOptions)

func OptWorkedNumber(v string) EditedPropositionOption {
	return func(opts *EditedPropositionOptions) {
		opts.WorkedNumber = v
	}
}

func OptDET(v string) EditedPropositionOption {
	return func(opts *EditedPropositionOptions) {
		opts.DET = v
	}
}

func OptSuccessful(v bool) EditedPropositionOption {
	return func(opts *EditedPropositionOptions) {
		opts.Successful = v
	}
}

func OptDeliveryDate(v time.Time) EditedPropositionOption {
	return func(opts *EditedPropositionOptions) {
		opts.DeliveryDate = v
	}
}

func OptUpdatedDeliveryDate(v time.Time) EditedPropositionOption {
	return func(opts *EditedPropositionOptions) {
		opts.UpdatedDeliveryDate = v
	}
}

func (e *EditedPropositionPram) ConvertToEditedProposition() *report.EditedProposition {
	return &report.EditedProposition{
		WorkedNumber:        e.WorkedNumber,
		DET:                 e.DET,
		Successful:          e.Successful,
		DeliveryDate:        e.DeliveryDate,
		UpdatedDeliveryDate: e.UpdatedDeliveryDate,
	}
}

func TestEditedPropositionPramCreate(options ...EditedPropositionOption) *EditedPropositionPram {
	// デフォルト値
	opts := &EditedPropositionOptions{
		WorkedNumber:        "99A-1234",
		DET:                 "99",
		Successful:          false,
		DeliveryDate:        time.Date(2099, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedDeliveryDate: time.Date(2099, 1, 25, 0, 0, 0, 0, time.UTC),
	}

	for _, option := range options {
		option(opts)
	}

	return &EditedPropositionPram{
		WorkedNumber:        opts.WorkedNumber,
		DET:                 opts.DET,
		Successful:          opts.Successful,
		DeliveryDate:        opts.DeliveryDate,
		UpdatedDeliveryDate: opts.UpdatedDeliveryDate,
	}
}

type ReportPram struct {
	Tos                []EmailAddressPram
	CCs                []EmailAddressPram
	BCCs               []EmailAddressPram
	From               EmailAddressPram
	ReplyTo            EmailAddressPram
	Subject            string
	EditedPropositions []EditedPropositionPram
	PrefixReport       string
	SuffixReport       string
	Replacements       map[string]string
}

type ReportOptions struct {
	tos                []EmailAddressPram
	ccs                []EmailAddressPram
	bccs               []EmailAddressPram
	from               EmailAddressPram
	replyTo            EmailAddressPram
	subject            string
	editedPropositions []EditedPropositionPram
	prefixReport       string
	suffixReport       string
	replacements       map[string]string
}

type ReportOption func(*ReportOptions)

func OptTos(v []EmailAddressPram) ReportOption {
	return func(opts *ReportOptions) {
		opts.tos = v
	}
}

func OptCCs(v []EmailAddressPram) ReportOption {
	return func(opts *ReportOptions) {
		opts.ccs = v
	}
}

func OptBCCs(v []EmailAddressPram) ReportOption {
	return func(opts *ReportOptions) {
		opts.bccs = v
	}
}

func OptFrom(v EmailAddressPram) ReportOption {
	return func(opts *ReportOptions) {
		opts.from = v
	}
}

func OptReplyTo(v EmailAddressPram) ReportOption {
	return func(opts *ReportOptions) {
		opts.replyTo = v
	}
}

func OptSubject(v string) ReportOption {
	return func(opts *ReportOptions) {
		opts.subject = v
	}
}

func OptEditedPropositions(v []EditedPropositionPram) ReportOption {
	return func(opts *ReportOptions) {
		opts.editedPropositions = v
	}
}

func OptPrefixReport(v string) ReportOption {
	return func(opts *ReportOptions) {
		opts.prefixReport = v
	}
}

func OptSuffixReport(v string) ReportOption {
	return func(opts *ReportOptions) {
		opts.suffixReport = v
	}
}

func TestReportPramCreate(options ...ReportOption) *ReportPram {
	// デフォルト値
	opts := &ReportOptions{
		tos:                []EmailAddressPram{*TestEmailAddressPramCreate()},
		ccs:                []EmailAddressPram{},
		bccs:               []EmailAddressPram{},
		from:               *TestEmailAddressPramCreate(OptName("送信者"), OptAddress("testing@example.com")),
		replyTo:            *TestEmailAddressPramCreate(OptName("返信先"), OptAddress("testing-return@example.com")),
		subject:            "結果報告",
		editedPropositions: []EditedPropositionPram{*TestEditedPropositionPramCreate()},
		prefixReport:       "結果報告:接頭辞",
		suffixReport:       "結果報告:接尾辞",
		replacements:       map[string]string{"count": "999"},
	}

	for _, option := range options {
		option(opts)
	}

	return &ReportPram{
		Tos:                opts.tos,
		CCs:                opts.ccs,
		BCCs:               opts.bccs,
		From:               opts.from,
		ReplyTo:            opts.replyTo,
		Subject:            opts.subject,
		EditedPropositions: opts.editedPropositions,
		PrefixReport:       opts.prefixReport,
		SuffixReport:       opts.suffixReport,
		Replacements:       opts.replacements,
	}
}

type Executor interface {
	Execute(r ReportPram) (string, error)
}

type SendingReportUseCase struct {
	sugar  *zap.SugaredLogger
	sender report.Sender
}

func NewSendingReportUseCase(
	sugar *zap.SugaredLogger,
	sender report.Sender,
) *SendingReportUseCase {
	return &SendingReportUseCase{
		sugar:  sugar,
		sender: sender,
	}
}

func (s *SendingReportUseCase) Execute(r ReportPram) (string, error) {
	return s.sender.Send(
		ConvertToEmailAddresses(r.Tos),
		ConvertToEmailAddresses(r.CCs),
		ConvertToEmailAddresses(r.BCCs),
		*r.From.ConvertToEmailAddress(),
		*r.ReplyTo.ConvertToEmailAddress(),
		r.Subject,
		ConvertToEditedProposition(r.EditedPropositions),
		r.PrefixReport,
		r.SuffixReport,
		r.Replacements,
	)
}

func ConvertToEmailAddresses(e []EmailAddressPram) []report.EmailAddress {
	var ad []report.EmailAddress
	for _, v := range e {
		ad = append(ad, *v.ConvertToEmailAddress())
	}
	return ad
}

func ConvertToEditedProposition(e []EditedPropositionPram) []report.EditedProposition {
	var res []report.EditedProposition
	for _, v := range e {
		res = append(res, *v.ConvertToEditedProposition())
	}
	return res
}
