//go:build wireinject
// +build wireinject

package twiliosendmail

import (
	"SynchronizeMonorevoDeliveryDates/domain/report"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewSendGridMail,
	wire.Bind(new(report.Sender), new(*SendGridMail)),
	NewSendGridConfig,
)
