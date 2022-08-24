//go:build wireinject
// +build wireinject

package monorevo

import "github.com/google/wire"

var Set = wire.NewSet(
	NewPropositionFetchingUseCase,
	NewPropositionPostingUseCase,
	wire.Bind(new(FetchingExecutor), new(*PropositionFetchingUseCase)),
	wire.Bind(new(PostingExecutor), new(*PropositionPostingUseCase)),
)
