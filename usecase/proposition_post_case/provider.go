//go:build wireinject
// +build wireinject

package proposition_post_case

import "github.com/google/wire"

var Set = wire.NewSet(
	NewPropositionPostingUseCase,
	wire.Bind(new(PostingExecutor), new(*PropositionPostingUseCase)),
)
