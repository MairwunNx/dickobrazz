package database

import "time"

type DocumentCockDynamic struct {
	IndividualCockTotal []struct {
		Total int `bson:"total"`
	} `bson:"individual_cock_total"`

	IndividualCockRecent []struct {
		Average int `bson:"average"`
	} `bson:"individual_cock_recent"`

	IndividualIrk []struct {
		Irk float64 `bson:"irk"`
	} `bson:"individual_irk"`

	IndividualRecord []struct {
		RequestedAt time.Time `bson:"requested_at"`
		Total       int       `bson:"total"`
	} `bson:"individual_record"`

	IndividualDominance []struct {
		Dominance float64 `bson:"dominance"`
	} `bson:"individual_dominance"`

	IndividualDailyGrowth []struct {
		Average float64 `bson:"average"`
	} `bson:"individual_daily_growth"`

	IndividualDailyDynamics []struct {
		YesterdayCockChange        int     `bson:"yesterday_cock_change"`
		YesterdayCockChangePercent float64 `bson:"yesterday_cock_change_percent"`
	} `bson:"individual_daily_dynamics"`

	Overall []struct {
		Size int `bson:"size"`
	} `bson:"overall"`

	OverallRecent []struct {
		Average int `bson:"average"`
		Median  int `bson:"median"`
	} `bson:"overall_recent"`

	Uniques []struct {
		Count int `bson:"count"`
	} `bson:"uniques"`

	Distribution []struct {
		HugePercent   float64 `bson:"huge"`
		LittlePercent float64 `bson:"little"`
	} `bson:"distribution"`

	Record []struct {
		RequestedAt time.Time `bson:"requested_at"`
		Total       int       `bson:"total"`
	} `bson:"record"`

	TotalCocksCount []struct {
		TotalCount int `bson:"total_count"`
	} `bson:"total_cocks_count"`

	IndividualCocksCount []struct {
		UserCount int `bson:"user_count"`
	} `bson:"individual_cocks_count"`

	IndividualLuck []struct {
		LuckCoefficient float64 `bson:"luck_coefficient"`
	} `bson:"individual_luck"`

	IndividualVolatility []struct {
		Volatility float64 `bson:"volatility"`
	} `bson:"individual_volatility"`

	IndividualFiveCocksDynamics []struct {
		FiveCocksChange        int     `bson:"five_cocks_change"`
		FiveCocksChangePercent float64 `bson:"five_cocks_change_percent"`
	} `bson:"individual_five_cocks_dynamics"`

	IndividualGrowthSpeed []struct {
		GrowthSpeed float64 `bson:"growth_speed"`
	} `bson:"individual_growth_speed"`
}
