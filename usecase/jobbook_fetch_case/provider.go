//go:build wireinject
// +build wireinject

package jobbook_fetch_case

import "github.com/google/wire"

var Set = wire.NewSet(
	NewJobBookFetchingUseCase,
	wire.Bind(new(Executor), new(*JobBookFetchingUseCase)),
)
