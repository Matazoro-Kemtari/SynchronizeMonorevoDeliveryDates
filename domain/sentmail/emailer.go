package sentmail

type EMailer interface {
	Send(
		tos []EmailAddress,
		ccs []EmailAddress,
		bccs []EmailAddress,
		from EmailAddress,
		subject string,
		body string,
		replacements map[string]string,
	) (string, error)
}

// 電子メールの名前とアドレス情報
type EmailAddress struct {
	Name    string
	Address string
}
