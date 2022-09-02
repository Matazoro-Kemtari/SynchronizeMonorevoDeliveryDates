package report

import "time"

type Sender interface {
	Send(
		tos []EmailAddress,
		ccs []EmailAddress,
		bccs []EmailAddress,
		from EmailAddress,
		subject string,
		editedPropositions []EditedProposition,
		prefixReport string,
		suffixReport string,
	) (string, error)
}

// 編集結果
type EditedProposition struct {
	WorkedNumber        string
	DET                 string
	Successful          bool
	DeliveryDate        time.Time
	UpdatedDeliveryDate time.Time
}

// 編集結果レポート
type EditedReport struct {
	EditedPropositions []EditedProposition
}

// 電子メールの名前とアドレス情報
type EmailAddress struct {
	Name    string
	Address string
}
