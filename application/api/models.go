package api

import "dickobrazz/application/datetime"

// DataResponse — обёртка ответов API.
// Все успешные ответы возвращают данные в поле data.
// Соответствует паттерну { data: T } в OpenAPI spec.
type DataResponse[T any] struct {
	Data T `json:"data"`
}

// PageMeta — метаданные пагинации.
// Используется в лидербордах (ruler, race, ladder) и сезонах.
// limit — лимит записей на странице (13–50), page — номер страницы (≥1).
type PageMeta struct {
	Limit      int `json:"limit"`
	Page       int `json:"page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// CockSizeData — ответ POST /api/v1/cock/size (стянуть кокич).
// Protected endpoint. size — сгенерированный размер в см, hash/salt — для верификации.
type CockSizeData struct {
	Size     int    `json:"size"`
	Hash     string `json:"hash"`
	Salt     string `json:"salt"`
	PulledAt string `json:"pulled_at"`
}

// LeaderboardEntry — запись дейли-лидерборда (кок-рулетка).
// Одна строка топа за текущий день. size — размер за день.
type LeaderboardEntry struct {
	UserID   int64  `json:"user_id"`
	Nickname string `json:"nickname"`
	Size     int    `json:"size"`
}

// UserNeighborhood — соседи пользователя в дейли-лидерборде.
// above — записи выше пользователя, self — позиция пользователя (nil если не участвовал),
// below — записи ниже. Заполняется только при user-контексте (cookie/bearer/CSOT).
type UserNeighborhood struct {
	Above []LeaderboardEntry `json:"above"`
	Self  *LeaderboardEntry  `json:"self"`
	Below []LeaderboardEntry `json:"below"`
}

// CockRulerData — ответ GET /api/v1/cock/ruler (дейли-топ лидерборд).
// Публичный endpoint. user_position и neighborhood заполняются при авторизации.
type CockRulerData struct {
	Leaders           []LeaderboardEntry `json:"leaders"`
	TotalParticipants int                `json:"total_participants"`
	UserPosition      *int               `json:"user_position"`
	Neighborhood      UserNeighborhood   `json:"neighborhood"`
	Page              PageMeta           `json:"page"`
}

// RaceEntry — запись сезонного лидерборда (гонка коков).
// total_size — сумма размеров за текущий сезон (3 месяца).
type RaceEntry struct {
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	TotalSize int    `json:"total_size"`
}

// RaceNeighborhood — соседи пользователя в сезонном лидерборде.
// Структура аналогична UserNeighborhood, но для сезонного топа.
type RaceNeighborhood struct {
	Above []RaceEntry `json:"above"`
	Self  *RaceEntry  `json:"self"`
	Below []RaceEntry `json:"below"`
}

// RaceSeasonInfo — информация о текущем сезоне гонки.
// season_num — номер сезона, start_date/end_date — границы сезона (ISO 8601).
type RaceSeasonInfo struct {
	SeasonNum int                    `json:"season_num"`
	StartDate datetime.LocalDateTime `json:"start_date"`
	EndDate   datetime.LocalDateTime `json:"end_date"`
}

// CockRaceData — ответ GET /api/v1/cock/race (лидерборд сезона).
// Публичный endpoint. Персонализация при user-контексте.
type CockRaceData struct {
	Season            RaceSeasonInfo   `json:"season"`
	Leaders           []RaceEntry      `json:"leaders"`
	TotalParticipants int              `json:"total_participants"`
	UserPosition      *int             `json:"user_position"`
	Neighborhood      RaceNeighborhood `json:"neighborhood"`
	Page              PageMeta         `json:"page"`
}

// LadderEntry — запись вечного лидерборда (ладдер коков).
// total_size — сумма размеров за всё время.
type LadderEntry struct {
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	TotalSize int    `json:"total_size"`
}

// LadderNeighborhood — соседи пользователя в вечном лидерборде.
type LadderNeighborhood struct {
	Above []LadderEntry `json:"above"`
	Self  *LadderEntry  `json:"self"`
	Below []LadderEntry `json:"below"`
}

// CockLadderData — ответ GET /api/v1/cock/ladder (общий лидерборд).
// Публичный endpoint. Персонализация при user-контексте.
type CockLadderData struct {
	Leaders           []LadderEntry      `json:"leaders"`
	TotalParticipants int                `json:"total_participants"`
	UserPosition      *int               `json:"user_position"`
	Neighborhood      LadderNeighborhood `json:"neighborhood"`
	Page              PageMeta           `json:"page"`
}

// CockDynamicRecentStat — средние и медианные значения за последние измерения.
// Используется в CockDynamicGlobalData (общая динамика).
type CockDynamicRecentStat struct {
	Average float64 `json:"average"`
	Median  float64 `json:"median"`
}

// CockDynamicPercentile — распределение по перцентилям.
// huge_percent — доля «огромных», little_percent — доля «маленьких».
type CockDynamicPercentile struct {
	HugePercent   float64 `json:"huge_percent"`
	LittlePercent float64 `json:"little_percent"`
}

// CockDynamicGlobalRecord — глобальный рекорд (максимальный total за один день).
// requested_at — дата рекорда (nullable).
type CockDynamicGlobalRecord struct {
	RequestedAt *datetime.LocalDateTime `json:"requested_at"`
	Total       int                     `json:"total"`
}

// CockDynamicGlobalData — ответ GET /api/v1/cock/dynamic/global.
// Общая динамика коков. Публичный endpoint, авторизация не требуется.
// Соответствует CockDynamicOverall в spec.
type CockDynamicGlobalData struct {
	TotalSize       int                     `json:"total_size"`
	UniqueUsers     int                     `json:"unique_users"`
	Recent          CockDynamicRecentStat   `json:"recent"`
	Distribution    CockDynamicPercentile   `json:"distribution"`
	Record          CockDynamicGlobalRecord `json:"record"`
	TotalCocksCount int                     `json:"total_cocks_count"`
	GrowthSpeed     float64                 `json:"growth_speed"`
}

// CockDynamicDailyDynamics — дневная динамика пользователя.
// Изменение относительно вчера (абсолютное и в процентах).
type CockDynamicDailyDynamics struct {
	YesterdayCockChange        int     `json:"yesterday_cock_change"`
	YesterdayCockChangePercent float64 `json:"yesterday_cock_change_percent"`
}

// CockDynamicFiveCocksDynamics — динамика за последние 5 измерений.
type CockDynamicFiveCocksDynamics struct {
	FiveCocksChange        int     `json:"five_cocks_change"`
	FiveCocksChangePercent float64 `json:"five_cocks_change_percent"`
}

// CockDynamicPersonalRecord — персональный рекорд пользователя.
// requested_at — дата рекорда (nullable), size — максимальный размер.
type CockDynamicPersonalRecord struct {
	RequestedAt *datetime.LocalDateTime `json:"requested_at"`
	Size        int                     `json:"size"`
}

// CockDynamicPersonalData — ответ GET /api/v1/cock/dynamic/personal.
// Персональная динамика коков. Protected endpoint.
// Соответствует CockDynamicIndividual в spec.
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
	FirstCockDate      *datetime.LocalDateTime      `json:"first_cock_date"`
	LuckCoefficient    float64                      `json:"luck_coefficient"`
	Volatility         float64                      `json:"volatility"`
	CocksCount         int                          `json:"cocks_count"`
}

// AchievementData — достижение пользователя (Achievement в spec).
// progress/max_progress — прогресс выполнения, completed — достигнуто ли.
type AchievementData struct {
	ID          string `json:"id"`
	Emoji       string `json:"emoji"`
	Respects    int    `json:"respects"`
	Completed   bool   `json:"completed"`
	Progress    int    `json:"progress"`
	MaxProgress int    `json:"max_progress"`
}

// CockAchievementsData — ответ GET /api/v1/cock/achievements.
// Protected endpoint. achievements_done_percent — процент выполненных ачивок.
type CockAchievementsData struct {
	Achievements            []AchievementData `json:"achievements"`
	AchievementsTotal       int               `json:"achievements_total"`
	AchievementsDone        int               `json:"achievements_done"`
	AchievementsDonePercent float64           `json:"achievements_done_percent"`
}

// SeasonWinner — победитель сезона (топ-3 или топ-N).
// place — место (1, 2, 3 и т.д.).
type SeasonWinner struct {
	UserID    int64  `json:"user_id"`
	Nickname  string `json:"nickname"`
	TotalSize int    `json:"total_size"`
	Place     int    `json:"place"`
}

// SeasonNeighborhood — соседи пользователя в сезоне.
// above/self/below — победители выше, сам пользователь, ниже.
type SeasonNeighborhood struct {
	Above []SeasonWinner `json:"above"`
	Self  *SeasonWinner  `json:"self"`
	Below []SeasonWinner `json:"below"`
}

// SeasonWithWinners — сезон с победителями (SeasonInfo + winners в spec).
// is_active — текущий сезон. user_position/neighborhood — при user-контексте.
type SeasonWithWinners struct {
	SeasonNum         int                    `json:"season_num"`
	StartDate         datetime.LocalDateTime `json:"start_date"`
	EndDate           datetime.LocalDateTime `json:"end_date"`
	IsActive          bool                   `json:"is_active"`
	Winners           []SeasonWinner         `json:"winners"`
	TotalParticipants int                    `json:"total_participants"`
	UserPosition      *int                   `json:"user_position"`
	Neighborhood      SeasonNeighborhood     `json:"neighborhood"`
}

// CockSeasonsData — ответ GET /api/v1/cock/seasons.
// Публичный endpoint. user_season_wins — число побед пользователя в сезонах.
type CockSeasonsData struct {
	Seasons        []SeasonWithWinners `json:"seasons"`
	Page           PageMeta            `json:"page"`
	UserSeasonWins int                 `json:"user_season_wins"`
}

// RespectData — ответ GET /api/v1/cock/respects.
// Cock-Respect™: детализация респектов. Protected endpoint.
type RespectData struct {
	SeasonRespect      float64 `json:"season_respect"`
	AchievementRespect float64 `json:"achievement_respect"`
	TotalRespect       float64 `json:"total_respect"`
}

// UpdatePrivacyPayload — тело PATCH /api/v1/me/privacy.
// required: is_hidden. Скрывает пользователя в лидербордах.
type UpdatePrivacyPayload struct {
	IsHidden bool `json:"is_hidden"`
}

// UserProfile — профиль пользователя (GET /api/v1/me, PATCH /api/v1/me/privacy).
// created_at — nullable. is_hidden — скрыт ли в лидербордах.
type UserProfile struct {
	ID        int64                   `json:"id"`
	Username  string                  `json:"username"`
	IsHidden  bool                    `json:"is_hidden"`
	CreatedAt *datetime.LocalDateTime `json:"created_at"`
}
