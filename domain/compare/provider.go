package compare

import "github.com/google/wire"

var Set = wire.NewSet(
	NewDifference,
	wire.Bind(new(Extractor), new(*Difference)),
)
