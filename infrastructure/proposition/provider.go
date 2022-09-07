//go:build wireinject
// +build wireinject

package proposition

import (
	"SynchronizeMonorevoDeliveryDates/domain/monorevo"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewPropositionTable,
	wire.Bind(new(monorevo.MonorevoFetcher), new(*PropositionTable)),
	wire.Bind(new(monorevo.MonorevoPoster), new(*PropositionTable)),
	NewMonorevoUserConfig,
)
