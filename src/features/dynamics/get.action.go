package dynamics

import (
	"context"
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/datetime"
	"dickobrazz/src/shared/formatting"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type GetAction struct {
	api *api.APIClient
	loc *localization.LocalizationManager
}

func NewGetAction(apiClient *api.APIClient, loc *localization.LocalizationManager) *GetAction {
	return &GetAction{api: apiClient, loc: loc}
}

func (a *GetAction) Execute(ctx context.Context, log *logging.Logger, localizer *i18n.Localizer, userID int64, username string) (string, error) {
	type globalResult struct {
		data *api.CockDynamicGlobalData
		err  error
	}
	type personalResult struct {
		data *api.CockDynamicPersonalData
		err  error
	}
	type seasonsResult struct {
		data *api.CockSeasonsData
		err  error
	}
	type respectsResult struct {
		data *api.RespectData
		err  error
	}

	globalCh := make(chan globalResult, 1)
	personalCh := make(chan personalResult, 1)
	seasonsCh := make(chan seasonsResult, 1)
	respectsCh := make(chan respectsResult, 1)

	go func() {
		data, err := a.api.GetCockDynamicGlobal(ctx)
		globalCh <- globalResult{data, err}
	}()
	go func() {
		data, err := a.api.GetCockDynamicPersonal(ctx, userID, username)
		personalCh <- personalResult{data, err}
	}()
	go func() {
		data, err := a.api.GetCockSeasons(ctx, userID, username, 15, 1)
		seasonsCh <- seasonsResult{data, err}
	}()
	go func() {
		data, err := a.api.GetCockRespects(ctx, userID, username)
		respectsCh <- respectsResult{data, err}
	}()

	globalRes := <-globalCh
	personalRes := <-personalCh
	seasonsRes := <-seasonsCh
	respectsRes := <-respectsCh

	if globalRes.err != nil {
		log.E("Failed to get global dynamic", logging.InnerError, globalRes.err)
		return "", globalRes.err
	}
	if personalRes.err != nil {
		log.E("Failed to get personal dynamic", logging.InnerError, personalRes.err)
		text := a.loc.Localize(localizer, "MsgCockDynamicNoData", nil)
		return text, nil
	}

	global := globalRes.data
	personal := personalRes.data

	userCockRespect := 0
	if respectsRes.err == nil && respectsRes.data != nil {
		userCockRespect = int(respectsRes.data.TotalRespect)
	}
	userSeasonWins := 0
	if seasonsRes.err == nil && seasonsRes.data != nil {
		userSeasonWins = seasonsRes.data.UserSeasonWins
	}

	individualRecordTotal := personal.Record.Size
	individualRecordDate := datetime.NowTime()
	if personal.Record.RequestedAt != nil {
		individualRecordDate = personal.Record.RequestedAt.Time
	}

	overallRecordTotal := global.Record.Total
	overallRecordDate := datetime.NowTime()
	if global.Record.RequestedAt != nil {
		overallRecordDate = global.Record.RequestedAt.Time
	}

	var userPullingPeriod string
	if personal.FirstCockDate != nil {
		userPullingPeriod = formatting.FormatUserPullingPeriod(a.loc, localizer, personal.FirstCockDate.Time, datetime.NowTime())
	} else {
		userPullingPeriod = a.loc.Localize(localizer, localization.MsgUserPullingRecently, nil)
	}

	text := NewMsgCockDynamicsTemplate(
		a.loc, localizer,
		global.TotalSize, global.UniqueUsers,
		int(global.Recent.Average), int(global.Recent.Median),
		personal.TotalSize, int(personal.RecentAverage), personal.Irk,
		individualRecordTotal, individualRecordDate,
		personal.DailyDynamics.YesterdayCockChangePercent, personal.DailyDynamics.YesterdayCockChange,
		personal.FiveCocksDynamics.FiveCocksChangePercent, personal.FiveCocksDynamics.FiveCocksChange,
		global.Distribution.HugePercent, global.Distribution.LittlePercent,
		overallRecordDate, overallRecordTotal,
		personal.Dominance,
		userSeasonWins, userCockRespect,
		global.TotalCocksCount, personal.CocksCount,
		personal.LuckCoefficient, personal.Volatility,
		personal.GrowthSpeed, global.GrowthSpeed,
		userPullingPeriod,
	)

	return text, nil
}
