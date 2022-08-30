package appsetting

type SandboxModeDto struct {
	Monorevo bool
	SendGrid bool
}

type AppSettingDto struct {
	SandboxMode SandboxModeDto
}

type SettingLoader interface {
	Load(path string) (*AppSettingDto, error)
}

type AppSettingDtoOptions struct {
	SandboxMode SandboxModeDto
}

type AppSettingDtoOption func(*AppSettingDtoOptions)

func OptSandboxMode(v SandboxModeDto) AppSettingDtoOption {
	return func(opts *AppSettingDtoOptions) {
		opts.SandboxMode = v
	}
}

func TestAppSettingDtoCreate(options ...AppSettingDtoOption) *AppSettingDto {
	// デフォルト値
	opts := &AppSettingDtoOptions{
		SandboxMode: *TestSandboxModeDtoCreate(),
	}

	for _, option := range options {
		option(opts)
	}

	return &AppSettingDto{
		SandboxMode: opts.SandboxMode,
	}
}

type SandboxModeDtoOptions struct {
	Monorevo bool
	SendGrid bool
}

type SandboxModeDtoOption func(*SandboxModeDtoOptions)

func OptSandboxModeMonorevo(v bool) SandboxModeDtoOption {
	return func(opts *SandboxModeDtoOptions) {
		opts.Monorevo = v
	}
}

func OptSandboxModeSendGrid(v bool) SandboxModeDtoOption {
	return func(opts *SandboxModeDtoOptions) {
		opts.SendGrid = v
	}
}

func TestSandboxModeDtoCreate(options ...SandboxModeDtoOption) *SandboxModeDto {
	// デフォルト値
	opts := &SandboxModeDtoOptions{
		Monorevo: true,
		SendGrid: true,
	}

	for _, option := range options {
		option(opts)
	}

	return &SandboxModeDto{
		Monorevo: opts.Monorevo,
		SendGrid: opts.SendGrid,
	}
}
