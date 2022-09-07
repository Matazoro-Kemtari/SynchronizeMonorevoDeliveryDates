//go:build wireinject
// +build wireinject

package jobbook

import (
	"SynchronizeMonorevoDeliveryDates/domain/orderdb"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewRepository,
	wire.Bind(new(orderdb.JobBookFetcher), new(*Repository)),
	NewOrderDbConfig,
)
