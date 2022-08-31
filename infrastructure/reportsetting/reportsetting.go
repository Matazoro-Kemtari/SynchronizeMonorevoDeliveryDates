package reportsetting

import (
	"SynchronizeMonorevoDeliveryDates/usecase/reportsetting"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type MailAddress struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ReportSetting struct {
	SenderAddress      MailAddress   `json:"senderAddress"`
	ReplyToAddress     MailAddress   `json:"replyToAddress"`
	RecipientAddresses []MailAddress `json:"recipientAddresses"`
	CcAddresses        []MailAddress `json:"ccAddresses"`
	BccAddresses       []MailAddress `json:"bccAddresses"`
	Subject            string        `json:"subject"`
	PrefixReport       string        `json:"prefixReport"`
	SuffixReport       string        `json:"suffixReport"`
}

func (r *ReportSetting) ConvertToReportSettingDto() *reportsetting.ReportSettingDto {
	tos := []reportsetting.MailAddressDto{}
	for _, v := range r.RecipientAddresses {
		tos = append(tos, reportsetting.MailAddressDto{
			Email: v.Email,
			Name:  v.Name,
		})
	}
	ccs := []reportsetting.MailAddressDto{}
	for _, v := range r.CcAddresses {
		ccs = append(ccs, reportsetting.MailAddressDto{
			Email: v.Email,
			Name:  v.Name,
		})
	}
	bccs := []reportsetting.MailAddressDto{}
	for _, v := range r.BccAddresses {
		bccs = append(bccs, reportsetting.MailAddressDto{
			Email: v.Email,
			Name:  v.Name,
		})
	}
	return &reportsetting.ReportSettingDto{
		SenderAddress: reportsetting.MailAddressDto{
			Email: r.SenderAddress.Email,
			Name:  r.SenderAddress.Name,
		},
		ReplyToAddress: reportsetting.MailAddressDto{
			Email: r.ReplyToAddress.Email,
			Name:  r.ReplyToAddress.Name,
		},
		RecipientAddresses: tos,
		CcAddresses:        ccs,
		BccAddresses:       bccs,
		Subject:            r.Subject,
		PrefixReport:       r.PrefixReport,
		SuffixReport:       r.SuffixReport,
	}
}

type LoadableSetting struct {
	sugar *zap.SugaredLogger
}

func NewLoadableSetting(sugar *zap.SugaredLogger) *LoadableSetting {
	return &LoadableSetting{
		sugar: sugar,
	}
}

func (l *LoadableSetting) Load(path string) (*reportsetting.ReportSettingDto, error) {
	r, err := os.Open(path)
	if err != nil {
		msg := fmt.Sprintf("設定ファイルが開けませんでした path: %v error: %v", path, err)
		l.sugar.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	var setting ReportSetting
	if err := json.NewDecoder(r).Decode(&setting); err != nil {
		msg := fmt.Sprintf("jsonからGo構造体へデコードできませんでした path: %v error: %v", path, err)
		l.sugar.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// 詰め替えて返す
	tos := []reportsetting.MailAddressDto{}
	for _, v := range setting.RecipientAddresses {
		tos = append(tos, reportsetting.MailAddressDto{
			Email: v.Email,
			Name:  v.Name,
		})
	}
	ccs := []reportsetting.MailAddressDto{}
	for _, v := range setting.CcAddresses {
		ccs = append(ccs, reportsetting.MailAddressDto{
			Email: v.Email,
			Name:  v.Name,
		})
	}
	bccs := []reportsetting.MailAddressDto{}
	for _, v := range setting.BccAddresses {
		bccs = append(bccs, reportsetting.MailAddressDto{
			Email: v.Email,
			Name:  v.Name,
		})
	}
	return &reportsetting.ReportSettingDto{
		SenderAddress: reportsetting.MailAddressDto{
			Email: setting.SenderAddress.Email,
			Name:  setting.SenderAddress.Name,
		},
		RecipientAddresses: tos,
		CcAddresses:        ccs,
		BccAddresses:       bccs,
		Subject:            setting.Subject,
		PrefixReport:       setting.PrefixReport,
		SuffixReport:       setting.SuffixReport,
	}, nil
}

type Options struct {
	senderAddress      MailAddress
	replyToAddress     MailAddress
	recipientAddresses []MailAddress
	ccAddresses        []MailAddress
	bccAddresses       []MailAddress
	subject            string
	prefixReport       string
	suffixReport       string
}

type Option func(*Options)

func OptSenderAddress(address MailAddress) Option {
	return func(opts *Options) {
		opts.senderAddress = address
	}
}

func OptReplyToAddress(address MailAddress) Option {
	return func(opts *Options) {
		opts.replyToAddress = address
	}
}

func OptRecipientAddresses(addresses []MailAddress) Option {
	return func(opts *Options) {
		opts.recipientAddresses = addresses
	}
}

func OptCcAddresses(addresses []MailAddress) Option {
	return func(opts *Options) {
		opts.ccAddresses = addresses
	}
}

func OptBccAddresses(addresses []MailAddress) Option {
	return func(opts *Options) {
		opts.bccAddresses = addresses
	}
}

func OptSubject(subject string) Option {
	return func(opts *Options) {
		opts.subject = subject
	}
}

func OptPrefixReport(prefixReport string) Option {
	return func(opts *Options) {
		opts.prefixReport = prefixReport
	}
}

func OptSuffixReport(suffixReport string) Option {
	return func(opts *Options) {
		opts.suffixReport = suffixReport
	}
}

func TestReportSettingCreate(options ...Option) *ReportSetting {
	// デフォルト値
	opts := &Options{
		senderAddress: MailAddress{
			Email: "abc@example.com",
			Name:  "サンプル送信者",
		},
		replyToAddress: MailAddress{
			Email: "reply@example.com",
			Name:  "返信先",
		},
		recipientAddresses: []MailAddress{
			{Email: "to1@example.com", Name: "宛先1"},
			{Email: "to2@example.com", Name: "宛先2"},
		},
		ccAddresses: []MailAddress{
			{Email: "cc1@example.com", Name: "CC1"},
			{Email: "cc2@example.com", Name: "CC2"},
		},
		bccAddresses: []MailAddress{
			{Email: "bcc1@example.com", Name: "BCC1"},
			{Email: "bcc2@example.com", Name: "BCC2"},
		},
		subject:      "題名XXX",
		prefixReport: "接頭辞\n接頭辞\n接頭辞",
		suffixReport: "接尾辞\n接尾辞\n接尾辞",
	}

	for _, option := range options {
		option(opts)
	}

	return &ReportSetting{
		SenderAddress:      opts.senderAddress,
		ReplyToAddress:     opts.replyToAddress,
		RecipientAddresses: opts.recipientAddresses,
		CcAddresses:        opts.ccAddresses,
		BccAddresses:       opts.bccAddresses,
		Subject:            opts.subject,
		PrefixReport:       opts.prefixReport,
		SuffixReport:       opts.suffixReport,
	}
}
