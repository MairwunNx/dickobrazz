package models

type CockSizeResponse struct {
	Size int    `json:"size"`
	Hash string `json:"hash"`
	Salt string `json:"salt"`
}

type LeaderboardEntry struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Size     int    `json:"size"`
}

type CockRulerResponse struct {
	Leaders           []LeaderboardEntry `json:"leaders"`
	TotalParticipants int                `json:"total_participants"`
	UserPosition      int                `json:"user_position"`
	Neighborhood      UserNeighborhood   `json:"neighborhood"`
	Page              PageMeta           `json:"page"`
}

type RaceEntry struct {
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	TotalSize int    `json:"total_size"`
}

type CockRaceResponse struct {
	Season            *SeasonInfo      `json:"season,omitempty"`
	Leaders           []RaceEntry      `json:"leaders"`
	TotalParticipants int              `json:"total_participants"`
	UserPosition      int              `json:"user_position"`
	Neighborhood      RaceNeighborhood `json:"neighborhood"`
	Page              PageMeta         `json:"page"`
}

type CockDynamicResponse struct {
	Overall    CockDynamicOverall    `json:"overall"`
	Individual CockDynamicIndividual `json:"individual"`
}

type SeasonInfo struct {
	SeasonNum int    `json:"season_num"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	IsActive  bool   `json:"is_active"`
}

type CockSeasonsResponse struct {
	Seasons []SeasonWithWinners `json:"seasons"`
	Page    PageMeta            `json:"page"`
}

type Achievement struct {
	ID          string `json:"id"`
	Emoji       string `json:"emoji,omitempty"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Respects    int    `json:"respects,omitempty"`
	Completed   bool   `json:"completed"`
	Progress    int    `json:"progress,omitempty"`
	MaxProgress int    `json:"max_progress,omitempty"`
}

type CockAchievementsResponse struct {
	Achievements []Achievement `json:"achievements"`
	Page         PageMeta      `json:"page"`
}

type CockLadderResponse struct {
	Leaders           []LeaderboardEntry `json:"leaders"`
	TotalParticipants int                `json:"total_participants"`
	UserPosition      int                `json:"user_position"`
	Neighborhood      UserNeighborhood   `json:"neighborhood"`
	Page              PageMeta           `json:"page"`
}

type UserNeighborhood struct {
	Above []LeaderboardEntry `json:"above"`
	Self  *LeaderboardEntry  `json:"self,omitempty"`
	Below []LeaderboardEntry `json:"below"`
}

type RaceNeighborhood struct {
	Above []RaceEntry `json:"above"`
	Self  *RaceEntry  `json:"self,omitempty"`
	Below []RaceEntry `json:"below"`
}

type CockDynamicOverall struct {
	TotalSize       int                   `json:"total_size"`
	UniqueUsers     int                   `json:"unique_users"`
	Recent          CockDynamicRecentStat `json:"recent"`
	Distribution    CockDynamicPercentile `json:"distribution"`
	Record          CockDynamicRecord     `json:"record"`
	TotalCocksCount int                   `json:"total_cocks_count"`
	GrowthSpeed     float64               `json:"growth_speed"`
}

type CockDynamicIndividual struct {
	TotalSize          int                          `json:"total_size"`
	RecentAverage      int                          `json:"recent_average"`
	Irk                float64                      `json:"irk"`
	Record             CockDynamicRecord            `json:"record"`
	Dominance          float64                      `json:"dominance"`
	DailyGrowthAverage float64                      `json:"daily_growth_average"`
	DailyDynamics      CockDynamicDailyDynamics     `json:"daily_dynamics"`
	FiveCocksDynamics  CockDynamicFiveCocksDynamics `json:"five_cocks_dynamics"`
	GrowthSpeed        float64                      `json:"growth_speed"`
	FirstCockDate      string                       `json:"first_cock_date"`
	LuckCoefficient    float64                      `json:"luck_coefficient"`
	Volatility         float64                      `json:"volatility"`
	CocksCount         int                          `json:"cocks_count"`
}

type CockDynamicRecentStat struct {
	Average int `json:"average"`
	Median  int `json:"median"`
}

type CockDynamicPercentile struct {
	HugePercent   float64 `json:"huge_percent"`
	LittlePercent float64 `json:"little_percent"`
}

type CockDynamicRecord struct {
	RequestedAt string `json:"requested_at"`
	Total       int    `json:"total"`
}

type CockDynamicDailyDynamics struct {
	YesterdayCockChange        int     `json:"yesterday_cock_change"`
	YesterdayCockChangePercent float64 `json:"yesterday_cock_change_percent"`
}

type CockDynamicFiveCocksDynamics struct {
	FiveCocksChange        int     `json:"five_cocks_change"`
	FiveCocksChangePercent float64 `json:"five_cocks_change_percent"`
}

type SeasonWinner struct {
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	TotalSize int    `json:"total_size"`
	Place     int    `json:"place"`
}

type SeasonWithWinners struct {
	SeasonInfo
	Winners []SeasonWinner `json:"winners"`
}
