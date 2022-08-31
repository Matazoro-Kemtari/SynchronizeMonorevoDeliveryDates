package twiliosendmail

import (
	"SynchronizeMonorevoDeliveryDates/domain/report"
	"SynchronizeMonorevoDeliveryDates/usecase/appsetting"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

type SendGridConfig struct {
	ApiKey string
}

func NewSendGridConfig() *SendGridConfig {
	return &SendGridConfig{
		ApiKey: os.Getenv("API_KEY"),
	}
}

type Options struct {
	apiKey string
}
type Option func(*Options)

func OptApiKey(key string) Option {
	return func(opts *Options) {
		opts.apiKey = key
	}
}

func TestSendGridConfigCreate(options ...Option) *SendGridConfig {
	// デフォルト値
	opts := &Options{
		apiKey: "ABCDEFG",
	}

	for _, option := range options {
		option(opts)
	}

	return &SendGridConfig{
		ApiKey: opts.apiKey,
	}
}

type SendGridMail struct {
	sugar       *zap.SugaredLogger
	apiKey      string
	sandboxMode bool
}

func NewSendGridMail(
	sugar *zap.SugaredLogger,
	appcnf *appsetting.AppSettingDto,
	cnf *SendGridConfig,
) *SendGridMail {
	return &SendGridMail{
		sugar:       sugar,
		apiKey:      cnf.ApiKey,
		sandboxMode: appcnf.SandboxMode.SendGrid,
	}
}

func (m *SendGridMail) Send(
	tos []report.EmailAddress,
	ccs []report.EmailAddress,
	bccs []report.EmailAddress,
	from report.EmailAddress,
	replyTo report.EmailAddress,
	subject string,
	editedPropositions []report.EditedProposition,
	prefixReport string,
	suffixReport string,
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
	if len(editedPropositions) < 1 {
		m.sugar.Error("編集結果がありません")
		return "", fmt.Errorf("編集結果がありません")
	}

	// メッセージの構築
	message := mail.NewV3Mail()

	p := mail.NewPersonalization()

	// 宛先設定
	m.addRecipients(p, tos, ccs, bccs)

	message.AddPersonalizations(p)

	// 差出人
	f := mail.NewEmail(from.Name, from.Address)
	m.sugar.Info("差出人:", from.Name, from.Address)
	message.SetFrom(f)

	// 返信先
	reply := mail.NewEmail(replyTo.Name, replyTo.Address)
	m.sugar.Info("返信先:", replyTo.Name, replyTo.Address)
	message.SetReplyTo(reply)

	// 件名
	message.Subject = subject
	m.sugar.Info("件名:", subject)

	// 本文
	plain := m.makePlainText(editedPropositions, prefixReport, suffixReport)
	c := mail.NewContent("text/plain", plain)
	message.AddContent(c)
	m.sugar.Info("本文:", plain)
	html := m.makeHtmlText(editedPropositions, prefixReport, suffixReport)
	h := mail.NewContent("text/html", html)
	message.AddContent(h)

	// カスタムヘッダを指定 サンプルの例に倣って設定
	message.SetHeader("X-Sent-Using", "SendGrid-API")
	// SendGridのコンソールログで見分けるために設定
	file := getExecutableFileName()
	message.Categories = append(message.Categories, file)

	if m.sandboxMode {
		// メール設定 サンドボックスモードON
		message.MailSettings = &mail.MailSettings{
			SandboxMode: &mail.Setting{Enable: &m.sandboxMode},
		}
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

func SuccessfulStr(b bool) string {
	r := "失敗"
	if b {
		r = "成功"
	}
	return r
}

func (m *SendGridMail) makePlainText(
	editedPropositions []report.EditedProposition,
	prefixReport string,
	suffixReport string,
) string {
	body := replaceLf(prefixReport) + "\n"
	body += "作業No\tDET番号\t成否\t変更前納期\t変更後納期\n"
	for _, v := range editedPropositions {
		body += fmt.Sprintf(
			"%v\t%v\t%v\t%v\t%v\n",
			v.WorkedNumber,
			v.Det,
			SuccessfulStr(v.Successful),
			v.DeliveryDate.Format("2006/01/02"),
			v.UpdatedDeliveryDate.Format("2006/01/02"),
		)
	}
	body += replaceLf(suffixReport)

	return body
}

func (m *SendGridMail) makeHtmlText(
	editedPropositions []report.EditedProposition,
	prefixReport string,
	suffixReport string,
) string {
	body := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN" "http://www.w3.org/TR/html4/loose.dtd">` + "\n"
	body += "<html><head>"
	body += `<meta http-equiv="content-type" content="text/html; charset=UTF-8">`
	body += "</head><body>\n"
	body += fmt.Sprintf("<p>%v</p>\n", replaceBr(prefixReport))
	body += "<table><tr><th>作業No</th><th>DET番号</th><th>成否</th><th>変更前納期</th><th>変更後納期</th></tr>"
	for _, v := range editedPropositions {
		body += fmt.Sprintf(
			"<tr><td>%v</td><td>%v</td><td>%v</td><td>%v</td><td>%v</td></tr>\n",
			v.WorkedNumber,
			v.Det,
			SuccessfulStr(v.Successful),
			v.DeliveryDate.Format("2006/01/02"),
			v.UpdatedDeliveryDate.Format("2006/01/02"),
		)
	}
	body += "</table>"
	body += fmt.Sprintf("<p>%v</p>\n", replaceBr(suffixReport))
	body += "</body></html>"
	return body
}

func replaceLf(txt string) string {
	txt = strings.ReplaceAll(txt, "\r\n", "\n")
	txt = strings.ReplaceAll(txt, "\r", "\n")
	return txt
}

func replaceBr(txt string) string {
	txt = replaceLf(txt)
	txt = strings.ReplaceAll(txt, "\n", "<br>\n")
	return txt
}

func (m *SendGridMail) addRecipients(p *mail.Personalization, tos []report.EmailAddress, ccs []report.EmailAddress, bccs []report.EmailAddress) {
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
}

func getExecutableFileName() string {
	exe, _ := os.Executable()
	file := filepath.Base(exe)
	return file
}
