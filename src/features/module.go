package features

import (
	"dickobrazz/src/features/achievements"
	"dickobrazz/src/features/dynamics"
	"dickobrazz/src/features/help"
	"dickobrazz/src/features/ladder"
	"dickobrazz/src/features/privacy"
	"dickobrazz/src/features/race"
	"dickobrazz/src/features/ruler"
	"dickobrazz/src/features/seasons"
	"dickobrazz/src/features/size"

	"go.uber.org/fx"
)

var Module = fx.Module("features",
	fx.Provide(
		// Size
		size.NewGenerateAction,
		size.NewHandler,

		// Ruler
		ruler.NewGetAction,
		ruler.NewHandler,

		// Ladder
		ladder.NewGetAction,
		ladder.NewHandler,

		// Race
		race.NewGetAction,
		race.NewHandler,

		// Dynamics
		dynamics.NewGetAction,
		dynamics.NewHandler,

		// Seasons
		seasons.NewGetAction,
		seasons.NewHandler,
		seasons.NewCallbackHandler,

		// Achievements
		achievements.NewGetAction,
		achievements.NewHandler,
		achievements.NewCallbackHandler,

		// Help
		help.NewHandler,

		// Privacy
		privacy.NewHandler,
	),
)
