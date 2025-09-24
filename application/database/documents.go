package database

import "time"

type DocumentCockDynamic struct {
	IndividualCock []struct {
		Total   int `bson:"total"`
		Average int `bson:"average"`
	} `bson:"individual_cock"`

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
		Size    int `bson:"size"`
		Average int `bson:"average"`
		Median  int `bson:"median"`
	} `bson:"overall"`

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
}
