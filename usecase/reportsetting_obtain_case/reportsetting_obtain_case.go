package reportsetting_obtain_case

type MailAddressDto struct {
	Email string
	Name  string
}

type ReportSettingDto struct {
	SenderAddress      MailAddressDto
	ReplyToAddress     MailAddressDto
	RecipientAddresses []MailAddressDto
	CCAddresses        []MailAddressDto
	BCCAddresses       []MailAddressDto
	Subject            string
	PrefixReport       string
	SuffixReport       string
}

type SettingLoader interface {
	Load(path string) (*ReportSettingDto, error)
}
