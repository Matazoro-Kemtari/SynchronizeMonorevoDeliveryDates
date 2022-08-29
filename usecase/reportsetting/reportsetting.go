package reportsetting

type MailAddressDto struct {
	Email string
	Name  string
}

type ReportSettingDto struct {
	SenderAddress      MailAddressDto
	RecipientAddresses []MailAddressDto
	CcAddresses        []MailAddressDto
	BccAddresses       []MailAddressDto
	Subject            string
	PrefixReport       string
	SuffixReport       string
}

type SettingLoader interface {
	Load(path string) (*ReportSettingDto, error)
}
