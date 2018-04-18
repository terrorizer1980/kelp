package strategy

import (
	"github.com/lightyeario/kelp/alfonso/strategy/level"
	"github.com/lightyeario/kelp/alfonso/strategy/sideStrategy"
	"github.com/lightyeario/kelp/support"
	"github.com/stellar/go/clients/horizon"
)

// AutonomousConfig contains the configuration params for this Strategy
type AutonomousConfig struct {
	PRICE_TOLERANCE              float64 `valid:"-"`
	AMOUNT_TOLERANCE             float64 `valid:"-"`
	SPREAD                       float64 `valid:"-"` // this is the bid-ask spread (i.e. it is not the spread from the center price)
	PLATEAU_THRESHOLD_PERCENTAGE float64 `valid:"-"`
	MIN_AMOUNT_SPREAD            float64 `valid:"-"` // reduces the order size by this percentage resulting in a gain anytime 1 unit more than the first layer is consumed
	MAX_AMOUNT_SPREAD            float64 `valid:"-"` // reduces the order size by this percentage resulting in a gain anytime 1 unit more than the first layer is consumed
	MAX_LEVELS                   int16   `valid:"-"` // max number of levels to have on either side
	LEVEL_DENSITY                float64 `valid:"-"` // value between 0.0 to 1.0 used as a probability
}

// MakeAutonomousStrategy is a factory method for AutonomousStrategy
func MakeAutonomousStrategy(
	txButler *kelp.TxButler,
	assetBase *horizon.Asset,
	assetQuote *horizon.Asset,
	config *AutonomousConfig,
) Strategy {
	sellSideStrategy := sideStrategy.MakeSellSideStrategy(
		txButler,
		assetBase,
		assetQuote,
		level.MakeAutonomousLevelProvider(
			config.SPREAD,
			config.PLATEAU_THRESHOLD_PERCENTAGE,
			false,
			config.MIN_AMOUNT_SPREAD,
			config.MAX_AMOUNT_SPREAD,
			config.MAX_LEVELS,
			config.LEVEL_DENSITY),
		config.PRICE_TOLERANCE,
		config.AMOUNT_TOLERANCE,
		false,
	)
	// switch sides of base/quote here for buy side
	buySideStrategy := sideStrategy.MakeSellSideStrategy(
		txButler,
		assetQuote,
		assetBase,
		level.MakeAutonomousLevelProvider(
			config.SPREAD,
			config.PLATEAU_THRESHOLD_PERCENTAGE,
			true, // real base is passed in as quote so pass in true
			config.MIN_AMOUNT_SPREAD,
			config.MAX_AMOUNT_SPREAD,
			config.MAX_LEVELS,
			config.LEVEL_DENSITY),
		config.PRICE_TOLERANCE,
		config.AMOUNT_TOLERANCE,
		true,
	)

	return MakeComposeStrategy(
		assetBase,
		assetQuote,
		buySideStrategy,
		sellSideStrategy,
	)
}