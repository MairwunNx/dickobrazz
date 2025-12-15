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

// DocumentUserAchievement представляет достижение пользователя в MongoDB
type DocumentUserAchievement struct {
	UserID        int64     `bson:"user_id"`
	AchievementID string    `bson:"achievement_id"`
	Completed     bool      `bson:"completed"`
	CompletedAt   time.Time `bson:"completed_at,omitempty"`
	Progress      int       `bson:"progress"` // Текущий прогресс для ачивок с прогрессом
	LastCheckedAt time.Time `bson:"last_checked_at"`
}

// Achievement представляет определение достижения
type Achievement struct {
	ID          string
	Emoji       string
	Name        string
	Description string
	Respects    int // Количество кок-респектов за выполнение
	MaxProgress int // Максимальный прогресс (0 если без прогресса)
}

type DocumentAchievementCheck struct {
	// Количество дерганий
	TotalPulls []struct {
		Count int32 `bson:"count"`
	} `bson:"total_pulls"`

	// Сумма размеров
	TotalSize []struct {
		Total int32 `bson:"total"`
	} `bson:"total_size"`

	// Конкретные значения
	Sniper30cm []struct {
		Count int32 `bson:"count"`
	} `bson:"sniper_30cm"`

	HalfHundred50cm []struct {
		Count int32 `bson:"count"`
	} `bson:"half_hundred_50cm"`

	BeautifulNumbers []struct {
		Count int32 `bson:"count"`
	} `bson:"beautiful_numbers"`

	// Последовательности
	Recent10 []struct {
		Size        int32     `bson:"size"`
		RequestedAt time.Time `bson:"requested_at"`
	} `bson:"recent_10"`

	// Экстремумы
	MaxSize []struct {
		Max int32 `bson:"max"`
	} `bson:"max_size"`

	MinSize []struct {
		Min int32 `bson:"min"`
	} `bson:"min_size"`

	// Временные
	EarlyBird []struct {
		Count int32 `bson:"count"`
	} `bson:"early_bird"`

	Speedrunner []struct {
		Count int32 `bson:"count"`
	} `bson:"speedrunner"`

	// Праздничные
	Valentine []struct {
		Count int32 `bson:"count"`
	} `bson:"valentine"`

	NewYearGift []struct {
		Count int32 `bson:"count"`
	} `bson:"new_year_gift"`

	// Динамика
	Lightning []struct {
		Count int32 `bson:"count"`
	} `bson:"lightning"`

	// Последние 31 кок
	Last31 []struct {
		Size int32 `bson:"size"`
	} `bson:"last_31"`
}

type DocumentSeasonCount struct {
	Count int32 `bson:"count"`
}

type DocumentTravelerCheck struct {
	UniqueSizes int32 `bson:"unique_sizes"`
}

type DocumentMuscoviteCheck struct {
	Count int32 `bson:"count"`
}

type DocumentGlobalMaxMin struct {
	Max int32 `bson:"max"`
	Min int32 `bson:"min"`
}
