package twiliosendmail

import (
	"SynchronizeMonorevoDeliveryDates/domain/report"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

type SendGridConfig struct {
	apiKey string
}

func NewSendGridConfig() *SendGridConfig {
	return &SendGridConfig{
		apiKey: os.Getenv("API_KEY"),
	}
}

func TestSendGridConfigCreate(apiKey string) *SendGridConfig {
	return &SendGridConfig{
		apiKey: apiKey,
	}
}

type SendGridMail struct {
	sugar  *zap.SugaredLogger
	apiKey string
}

func NewSendGridMail(sugar *zap.SugaredLogger, cnf *SendGridConfig) report.EMailer {
	return &SendGridMail{
		sugar:  sugar,
		apiKey: cnf.apiKey,
	}
}

func (m *SendGridMail) Send(
	tos []report.EmailAddress,
	ccs []report.EmailAddress,
	bccs []report.EmailAddress,
	from report.EmailAddress,
	subject string,
	body string,
	replacements map[string]string,
) (string, error) {
	if m.apiKey == "" {
		m.sugar.Error("API KEYが設定されていません")
		return "", fmt.Errorf("API KEYが設定されていません")
	}
	if len(tos) == 0 || tos[0].Address == "" {
		m.sugar.Error("宛先が設定されていません")
		return "", fmt.Errorf("宛先が設定されていません")
	}
	if from.Address == "" {
		m.sugar.Error("差出人が設定されていません")
		return "", fmt.Errorf("差出人が設定されていません")
	}

	// メッセージの構築
	message := mail.NewV3Mail()

	// 宛先設定
	p := mail.NewPersonalization()
	for _, to := range tos {
		adr := mail.NewEmail(to.Name, to.Address)
		p.AddTos(adr)
		m.sugar.Info("宛先設定TO:", to.Name, to.Address)
	}
	for _, cc := range ccs {
		adr := mail.NewEmail(cc.Name, cc.Address)
		p.AddCCs(adr)
		m.sugar.Info("宛先設定CC:", cc.Name, cc.Address)
	}
	for _, bcc := range bccs {
		adr := mail.NewEmail(bcc.Name, bcc.Address)
		p.AddBCCs(adr)
		m.sugar.Info("宛先設定BCC:", bcc.Name, bcc.Address)
	}
	// 置換設定
	for key, value := range replacements {
		p.SetSubstitution("%"+key+"%", value)
	}
	message.AddPersonalizations(p)
	// 差出人
	f := mail.NewEmail(from.Name, from.Address)
	m.sugar.Info("差出人:", from.Name, from.Address)
	message.SetFrom(f)
	// 件名
	message.Subject = subject
	m.sugar.Info("件名:", subject)
	// 本文
	c := mail.NewContent("text/plain", body)
	message.AddContent(c)
	m.sugar.Info("本文:", body)
	// カスタムヘッダを指定 サンプルの例に倣って設定
	message.SetHeader("X-Sent-Using", "SendGrid-API")
	// SendGridのコンソールログで見分けるために設定
	file := getExecutableFileName()
	message.Categories = append(message.Categories, file)

	// メール設定 サンドボックスモード
	// TODO: 設定で可変にする
	var sb bool = true
	message.MailSettings = &mail.MailSettings{
		SandboxMode: &mail.Setting{Enable: &sb},
	}

	// メール送信
	client := sendgrid.NewSendClient(m.apiKey)
	response, err := client.Send(message)
	if err != nil {
		return "", fmt.Errorf("error: %v", err)
	}
	// 送信完了
	m.sugar.Info(response.StatusCode)
	m.sugar.Info(response.Body)
	m.sugar.Info(response.Headers)
	return response.Headers["Date"][0], nil
}

func getExecutableFileName() string {
	exe, _ := os.Executable()
	file := filepath.Base(exe)
	return file
}
