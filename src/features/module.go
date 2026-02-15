package features

import (
	"dickobrazz/src/features/achievements"
	"dickobrazz/src/features/dynamics"
	"dickobrazz/src/features/health"
	"dickobrazz/src/features/help"
	"dickobrazz/src/features/ladder"
	"dickobrazz/src/features/privacy"
	"dickobrazz/src/features/promserver"
	"dickobrazz/src/features/race"
	"dickobrazz/src/features/ruler"
	"dickobrazz/src/features/seasons"
	"dickobrazz/src/features/size"

	"go.uber.org/fx"
)

var Module = fx.Module("features",
	size.Module,
	ruler.Module,
	ladder.Module,
	race.Module,
	dynamics.Module,
	seasons.Module,
	achievements.Module,
	help.Module,
	privacy.Module,
	health.Module,
	promserver.Module,
)
