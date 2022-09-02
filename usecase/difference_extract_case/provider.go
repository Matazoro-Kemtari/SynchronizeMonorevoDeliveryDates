//go:build wireinject
// +build wireinject

package difference_extract_case

import "github.com/google/wire"

var Set = wire.NewSet(
	NewExtractingPropositionUseCase,
	wire.Bind(new(Executor), new(*PropositionExtractingUseCase)),
)
