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
