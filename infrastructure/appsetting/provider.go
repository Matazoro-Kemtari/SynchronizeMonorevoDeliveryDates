//go:build wireinject
// +build wireinject

package appsetting

import (
	"SynchronizeMonorevoDeliveryDates/usecase/appsetting"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewLoadableSetting,
	wire.Bind(new(appsetting.SettingLoader), new(*LoadableSetting)),
)
