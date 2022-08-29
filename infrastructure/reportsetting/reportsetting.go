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
