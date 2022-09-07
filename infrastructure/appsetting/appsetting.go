package appsetting

import (
	"SynchronizeMonorevoDeliveryDates/usecase/appsetting_obtain_case"
	"encoding/json"
	"fmt"
	"os"

	"go.uber.org/zap"
)

type SandboxMode struct {
	Monorevo bool `json:"monorevo"`
	SendGrid bool `json:"sendgrid"`
}

type AppSetting struct {
	SandboxMode SandboxMode `json:"sandboxmode"`
}

func (m *AppSetting) ConvertToAppSettingDto() *appsetting_obtain_case.AppSettingDto {
	return &appsetting_obtain_case.AppSettingDto{
		SandboxMode: appsetting_obtain_case.SandboxModeDto{
			Monorevo: m.SandboxMode.Monorevo,
			SendGrid: m.SandboxMode.SendGrid,
		},
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

func (l *LoadableSetting) Load(path string) (*appsetting_obtain_case.AppSettingDto, error) {
	r, err := os.Open(path)
	if err != nil {
		msg := fmt.Sprintf("設定ファイルが開けませんでした path: %v error: %v", path, err)
		l.sugar.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	var setting AppSetting
	if err := json.NewDecoder(r).Decode(&setting); err != nil {
		msg := fmt.Sprintf("jsonからGo構造体へデコードできませんでした path: %v error: %v", path, err)
		l.sugar.Error(msg)
		return nil, fmt.Errorf(msg)
	}

	// 詰め替えて返す
	return setting.ConvertToAppSettingDto(), nil
}

type Options struct {
	sandboxMode SandboxMode
}

type Option func(*Options)

func OptSandboxMode(v SandboxMode) Option {
	return func(opts *Options) {
		opts.sandboxMode = v
	}
}

func TestAppSettingCreate(options ...Option) *AppSetting {
	// デフォルト値
	opts := &Options{
		sandboxMode: SandboxMode{
			Monorevo: true,
			SendGrid: true,
		},
	}

	for _, option := range options {
		option(opts)
	}

	return &AppSetting{
		SandboxMode: SandboxMode{},
	}
}
