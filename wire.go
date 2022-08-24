//go:build wireinject
// +build wireinject

package main

import (
	"SynchronizeMonorevoDeliveryDates/domain/compare"
	"SynchronizeMonorevoDeliveryDates/infrastructure/jobbook"
	"SynchronizeMonorevoDeliveryDates/infrastructure/proposition"
	"SynchronizeMonorevoDeliveryDates/presentation"
	"SynchronizeMonorevoDeliveryDates/usecase/difference"
	"SynchronizeMonorevoDeliveryDates/usecase/monorevo"
	"SynchronizeMonorevoDeliveryDates/usecase/orderdb"

	"github.com/google/wire"
	"go.uber.org/zap"
)

func InitializeSynchronize(log *zap.SugaredLogger) *presentation.SynchronizingDeliveryDate {
	wire.Build(
		presentation.Set,
		difference.Set,
		monorevo.Set,
		orderdb.Set,
		jobbook.Set,
		proposition.Set,
		compare.Set,
	)
	return nil
}
