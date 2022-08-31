//go:build wireinject
// +build wireinject

package twiliosendmail

import "github.com/google/wire"

var Set = wire.NewSet(
	NewSendGridMail,
	// TODO: ここまだ
	wire.Bind(new(), new()),
)
