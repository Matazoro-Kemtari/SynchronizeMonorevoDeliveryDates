//go:build wireinject
// +build wireinject

package report_send_case

import "github.com/google/wire"

var Set = wire.NewSet(
	NewSendingReportUseCase,
	wire.Bind(new(Executor), new(*SendingReportUseCase)),
)
