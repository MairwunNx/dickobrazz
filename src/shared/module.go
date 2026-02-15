package shared

import (
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/collector"
	"dickobrazz/src/shared/config"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"
	"dickobrazz/src/shared/metrics"

	"go.uber.org/fx"
)

var Module = fx.Module("shared",
	logging.Module,
	config.Module,
	localization.Module,
	api.Module,
	metrics.Module,
	collector.Module,
)
