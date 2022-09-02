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

type MailAddressOptions struct {
	Email string
	Name  string
}

type MailAddressOption func(*MailAddressOptions)

func OptEmail(v string) MailAddressOption {
	return func(opts *MailAddressOptions) {
		opts.Email = v
	}
}

func OptName(v string) MailAddressOption {
	return func(opts *MailAddressOptions) {
		opts.Name = v
	}
}

func TestMailAddressDtoCreate(options ...MailAddressOption) *MailAddressDto {
	// デフォルト値
	opts := &MailAddressOptions{
		Email: "abc@example.com",
		Name:  "サンプル",
	}

	for _, option := range options {
		option(opts)
	}

	return &MailAddressDto{
		Email: opts.Email,
		Name:  opts.Name,
	}
}

type ReportSettingOptions struct {
	SenderAddress      MailAddressDto
	ReplyToAddress     MailAddressDto
	RecipientAddresses []MailAddressDto
	CCAddresses        []MailAddressDto
	BCCAddresses       []MailAddressDto
	Subject            string
	PrefixReport       string
	SuffixReport       string
}

type ReportSettingOption func(*ReportSettingOptions)

func OptSenderAddress(v MailAddressDto) ReportSettingOption {
	return func(opts *ReportSettingOptions) {
		opts.SenderAddress = v
	}
}

func OptReplyToAddress(v MailAddressDto) ReportSettingOption {
	return func(opts *ReportSettingOptions) {
		opts.ReplyToAddress = v
	}
}

func OptRecipientAddresses(v []MailAddressDto) ReportSettingOption {
	return func(opts *ReportSettingOptions) {
		opts.RecipientAddresses = v
	}
}

func OptCCAddresses(v []MailAddressDto) ReportSettingOption {
	return func(opts *ReportSettingOptions) {
		opts.CCAddresses = v
	}
}

func OptBCCAddresses(v []MailAddressDto) ReportSettingOption {
	return func(opts *ReportSettingOptions) {
		opts.BCCAddresses = v
	}
}

func OptSubject(v string) ReportSettingOption {
	return func(opts *ReportSettingOptions) {
		opts.Subject = v
	}
}

func OptPrefixReport(v string) ReportSettingOption {
	return func(opts *ReportSettingOptions) {
		opts.PrefixReport = v
	}
}

func OptSuffixReport(v string) ReportSettingOption {
	return func(opts *ReportSettingOptions) {
		opts.SuffixReport = v
	}
}

func TestReportSettingDtoCreate(options ...ReportSettingOption) *ReportSettingDto {
	// デフォルト値
	opts := &ReportSettingOptions{
		SenderAddress:      *TestMailAddressDtoCreate(),
		ReplyToAddress:     *TestMailAddressDtoCreate(),
		RecipientAddresses: []MailAddressDto{*TestMailAddressDtoCreate()},
		CCAddresses:        []MailAddressDto{*TestMailAddressDtoCreate()},
		BCCAddresses:       []MailAddressDto{*TestMailAddressDtoCreate()},
		Subject:            "題名XXX",
		PrefixReport:       "接頭辞\n接頭辞\n接頭辞",
		SuffixReport:       "接尾辞\n接尾辞\n接尾辞",
	}

	for _, option := range options {
		option(opts)
	}

	return &ReportSettingDto{
		SenderAddress:      opts.SenderAddress,
		ReplyToAddress:     opts.ReplyToAddress,
		RecipientAddresses: opts.RecipientAddresses,
		CCAddresses:        opts.CCAddresses,
		BCCAddresses:       opts.BCCAddresses,
		Subject:            opts.Subject,
		PrefixReport:       opts.PrefixReport,
		SuffixReport:       opts.SuffixReport,
	}
}
