//go:build wireinject
// +build wireinject

package reportsetting

import (
	"SynchronizeMonorevoDeliveryDates/usecase/reportsetting"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewLoadableSetting,
	wire.Bind(new(reportsetting.SettingLoader), new(*LoadableSetting)),
)
