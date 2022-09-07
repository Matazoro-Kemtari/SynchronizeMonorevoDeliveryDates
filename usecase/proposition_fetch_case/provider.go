//go:build wireinject
// +build wireinject

package proposition_fetch_case

import "github.com/google/wire"

var Set = wire.NewSet(
	NewPropositionFetchingUseCase,
	wire.Bind(new(FetchingExecutor), new(*PropositionFetchingUseCase)),
)
