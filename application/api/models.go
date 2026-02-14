package api

// DataResponse — обёртка ответов API
type DataResponse[T any] struct {
	Data T `json:"data"`
}

// PageMeta — метаданные пагинации
type PageMeta struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// CockSizeData — ответ генерации размера
type CockSizeData struct {
	Size     int    `json:"size"`
	Hash     string `json:"hash"`
	Salt     string `json:"salt"`
	PulledAt string `json:"pulled_at"`
}

// LeaderboardEntry — запись дейли-лидерборда
type LeaderboardEntry struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Size     int    `json:"size"`
}

// UserNeighborhood — соседи пользователя в дейли-лидерборде
type UserNeighborhood struct {
	Above []LeaderboardEntry `json:"above"`
	Self  *LeaderboardEntry  `json:"self"`
	Below []LeaderboardEntry `json:"below"`
}

// CockRulerData — ответ дейли-лидерборда
type CockRulerData struct {
	Leaders           []LeaderboardEntry `json:"leaders"`
	TotalParticipants int                `json:"total_participants"`
	UserPosition      *int               `json:"user_position"`
	Neighborhood      UserNeighborhood   `json:"neighborhood"`
	Page              PageMeta           `json:"page"`
}

// RaceEntry — запись сезонного лидерборда
type RaceEntry struct {
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	TotalSize int    `json:"total_size"`
}

// RaceNeighborhood — соседи пользователя в сезонном лидерборде
type RaceNeighborhood struct {
	Above []RaceEntry `json:"above"`
	Self  *RaceEntry  `json:"self"`
	Below []RaceEntry `json:"below"`
}

// RaceSeasonInfo — информация о текущем сезоне
type RaceSeasonInfo struct {
	SeasonNum int    `json:"season_num"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// CockRaceData — ответ сезонного лидерборда
type CockRaceData struct {
	Season            RaceSeasonInfo   `json:"season"`
	Leaders           []RaceEntry      `json:"leaders"`
	TotalParticipants int              `json:"total_participants"`
	UserPosition      *int             `json:"user_position"`
	Neighborhood      RaceNeighborhood `json:"neighborhood"`
	Page              PageMeta         `json:"page"`
}

// LadderEntry — запись вечного лидерборда
type LadderEntry struct {
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	TotalSize int    `json:"total_size"`
}

// LadderNeighborhood — соседи пользователя в вечном лидерборде
type LadderNeighborhood struct {
	Above []LadderEntry `json:"above"`
	Self  *LadderEntry  `json:"self"`
	Below []LadderEntry `json:"below"`
}

// CockLadderData — ответ вечного лидерборда
type CockLadderData struct {
	Leaders           []LadderEntry      `json:"leaders"`
	TotalParticipants int                `json:"total_participants"`
	UserPosition      *int               `json:"user_position"`
	Neighborhood      LadderNeighborhood `json:"neighborhood"`
	Page              PageMeta           `json:"page"`
}

// CockDynamicRecentStat — средние/медианные значения
type CockDynamicRecentStat struct {
	Average float64 `json:"average"`
	Median  float64 `json:"median"`
}

// CockDynamicPercentile — распределение
type CockDynamicPercentile struct {
	HugePercent   float64 `json:"huge_percent"`
	LittlePercent float64 `json:"little_percent"`
}

// CockDynamicGlobalRecord — рекорд общий
type CockDynamicGlobalRecord struct {
	RequestedAt *string `json:"requested_at"`
	Total       int     `json:"total"`
}

// CockDynamicGlobalData — общая динамика коков
type CockDynamicGlobalData struct {
	TotalSize       int                     `json:"total_size"`
	UniqueUsers     int                     `json:"unique_users"`
	Recent          CockDynamicRecentStat   `json:"recent"`
	Distribution    CockDynamicPercentile   `json:"distribution"`
	Record          CockDynamicGlobalRecord `json:"record"`
	TotalCocksCount int                     `json:"total_cocks_count"`
	GrowthSpeed     float64                 `json:"growth_speed"`
}

// CockDynamicDailyDynamics — дневная динамика
type CockDynamicDailyDynamics struct {
	YesterdayCockChange        int     `json:"yesterday_cock_change"`
	YesterdayCockChangePercent float64 `json:"yesterday_cock_change_percent"`
}

// CockDynamicFiveCocksDynamics — динамика за 5 коков
type CockDynamicFiveCocksDynamics struct {
	FiveCocksChange        int     `json:"five_cocks_change"`
	FiveCocksChangePercent float64 `json:"five_cocks_change_percent"`
}

// CockDynamicPersonalRecord — персональный рекорд
type CockDynamicPersonalRecord struct {
	RequestedAt *string `json:"requested_at"`
	Size        int     `json:"size"`
}

// CockDynamicPersonalData — персональная динамика коков
type CockDynamicPersonalData struct {
	TotalSize          int                          `json:"total_size"`
	RecentAverage      float64                      `json:"recent_average"`
	Irk                float64                      `json:"irk"`
	Record             CockDynamicPersonalRecord    `json:"record"`
	Dominance          float64                      `json:"dominance"`
	DailyGrowthAverage float64                      `json:"daily_growth_average"`
	DailyDynamics      CockDynamicDailyDynamics     `json:"daily_dynamics"`
	FiveCocksDynamics  CockDynamicFiveCocksDynamics `json:"five_cocks_dynamics"`
	GrowthSpeed        float64                      `json:"growth_speed"`
	FirstCockDate      *string                      `json:"first_cock_date"`
	LuckCoefficient    float64                      `json:"luck_coefficient"`
	Volatility         float64                      `json:"volatility"`
	CocksCount         int                          `json:"cocks_count"`
}

// AchievementData — достижение пользователя
type AchievementData struct {
	ID          string `json:"id"`
	Emoji       string `json:"emoji"`
	Respects    int    `json:"respects"`
	Completed   bool   `json:"completed"`
	Progress    int    `json:"progress"`
	MaxProgress int    `json:"max_progress"`
}

// CockAchievementsData — ответ достижений
type CockAchievementsData struct {
	Achievements            []AchievementData `json:"achievements"`
	AchievementsTotal       int               `json:"achievements_total"`
	AchievementsDone        int               `json:"achievements_done"`
	AchievementsDonePercent float64           `json:"achievements_done_percent"`
}

// SeasonWinner — победитель сезона
type SeasonWinner struct {
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	TotalSize int    `json:"total_size"`
	Place     int    `json:"place"`
}

// SeasonNeighborhood — соседи в сезоне
type SeasonNeighborhood struct {
	Above []SeasonWinner `json:"above"`
	Self  *SeasonWinner  `json:"self"`
	Below []SeasonWinner `json:"below"`
}

// SeasonWithWinners — сезон с победителями
type SeasonWithWinners struct {
	SeasonNum         int                `json:"season_num"`
	StartDate         string             `json:"start_date"`
	EndDate           string             `json:"end_date"`
	IsActive          bool               `json:"is_active"`
	Winners           []SeasonWinner     `json:"winners"`
	TotalParticipants int                `json:"total_participants"`
	UserPosition      *int               `json:"user_position"`
	Neighborhood      SeasonNeighborhood `json:"neighborhood"`
}

// CockSeasonsData — ответ сезонов
type CockSeasonsData struct {
	Seasons        []SeasonWithWinners `json:"seasons"`
	Page           PageMeta            `json:"page"`
	UserSeasonWins int                 `json:"user_season_wins"`
}

// RespectData — ответ респектов
type RespectData struct {
	SeasonRespect      float64 `json:"season_respect"`
	AchievementRespect float64 `json:"achievement_respect"`
	TotalRespect       float64 `json:"total_respect"`
}

// UpdatePrivacyPayload — запрос обновления приватности
type UpdatePrivacyPayload struct {
	IsHidden bool `json:"is_hidden"`
}

// UserProfile — профиль пользователя
type UserProfile struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	IsHidden  bool    `json:"is_hidden"`
	CreatedAt *string `json:"created_at"`
}
