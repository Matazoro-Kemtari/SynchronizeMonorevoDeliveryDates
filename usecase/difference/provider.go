//go:build wireinject
// +build wireinject

package difference

import "github.com/google/wire"

var Set = wire.NewSet(
	NewExtractingPropositionUseCase,
	wire.Bind(new(Executor), new(*PropositionExtractingUseCase)),
)
