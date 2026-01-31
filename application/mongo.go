package application

import (
	"context"
	"dickobrazz/application/database"
	"dickobrazz/application/datetime"
	"dickobrazz/application/logging"
	"math"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const MaxSeasonsToShow = 14

type Cock struct {
	ID          string    `bson:"_id"`
	Size        int32     `bson:"size"`
	Nickname    string    `bson:"nickname"`
	UserID      int64     `bson:"user_id"`
	RequestedAt time.Time `bson:"requested_at"`
}

type UserCockRace struct { // Для аггрегаций только
	UserID    int64  `bson:"_id"`
	Nickname  string `bson:"nickname"`
	TotalSize int32  `bson:"total_size"`
}

type CockSeason struct {
	StartDate time.Time
	EndDate   time.Time
	IsActive  bool
	SeasonNum int
}

type SeasonWinner struct {
	UserID    int64  `bson:"_id"`
	Nickname  string `bson:"nickname"`
	TotalSize int32  `bson:"total_size"`
	Place     int    // 1, 2, 3
}

func InitializeMongoConnection(ctx context.Context, log *logging.Logger) *mongo.Client {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.F("MONGODB_URI does not have value, set it in .env file")
	}

	uri = os.ExpandEnv(uri)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetAppName("Dickobrazz").SetTimeout(10*time.Second))
	if err != nil {
		log.F("Failed to connect to MongoDB", logging.InnerError, err)
	}

	log.I("Successfully connected to MongoDB!")
	return client
}

func (app *Application) SaveCockToMongo(log *logging.Logger, cock *Cock) {
	collection := database.CollectionCocks(app.db)

	if _, err := collection.InsertOne(app.ctx, cock); err != nil {
		log.E("Failed to save cock to MongoDB", logging.InnerError, err)
	} else {
		log.I("Successfully saved cock to MongoDB")
	}
}

func (app *Application) GetUserProfile(log *logging.Logger, userID int64) *database.DocumentUserProfile {
	collection := database.CollectionUsers(app.db)

	var profile database.DocumentUserProfile
	if err := collection.FindOne(app.ctx, bson.D{{Key: "user_id", Value: userID}}).Decode(&profile); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil
		}
		log.E("Failed to get user profile", logging.InnerError, err)
		return nil
	}
	return &profile
}

func (app *Application) UpsertUserProfile(log *logging.Logger, userID int64, username string, isHidden bool) {
	collection := database.CollectionUsers(app.db)

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "user_id", Value: userID},
		{Key: "username", Value: username},
		{Key: "is_hidden", Value: isHidden},
		{Key: "updated_at", Value: datetime.NowTime()},
	}}}
	if _, err := collection.UpdateOne(app.ctx, bson.D{{Key: "user_id", Value: userID}}, update, options.Update().SetUpsert(true)); err != nil {
		log.E("Failed to upsert user profile", logging.InnerError, err)
	}
}

func (app *Application) AggregateCockSizes(log *logging.Logger) []UserCockRace {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineTopUsersBySize())
	if err != nil {
		log.E("Failed to aggregate cock sizes", logging.InnerError, err)
		return []UserCockRace{}
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil { // Кааааак же похуй...
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var results []UserCockRace
	if err = cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to parse aggregation results", logging.InnerError, err)
		return []UserCockRace{}
	}

	log.I("Successfully aggregated cock sizes")
	return results
}

func (app *Application) AggregateCockSizesForSeason(log *logging.Logger, season CockSeason) []UserCockRace {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineTopUsersInSeason(season.StartDate, season.EndDate))
	if err != nil {
		log.E("Failed to aggregate cock sizes for season", logging.InnerError, err)
		return []UserCockRace{}
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var results []UserCockRace
	if err = cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to parse season aggregation results", logging.InnerError, err)
		return []UserCockRace{}
	}

	log.I("Successfully aggregated cock sizes for season")
	return results
}

func (app *Application) GetUserAggregatedCock(log *logging.Logger, userID int64) *UserCockRace {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineUserTotalSize(userID))
	if err != nil {
		log.E("Failed to aggregate user cock sizes", logging.InnerError, err)
		return nil
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var result UserCockRace
	if cursor.Next(app.ctx) {
		if err := cursor.Decode(&result); err != nil {
			log.E("Failed to decode aggregation result", logging.InnerError, err)
			return nil
		}
		return &result
	}

	log.I("No cocks found for user")
	return nil
}

// GetUserCocksCount возвращает количество коков пользователя за все время
func (app *Application) GetUserCocksCount(log *logging.Logger, userID int64) int {
	collection := database.CollectionCocks(app.db)

	count, err := collection.CountDocuments(app.ctx, bson.M{"user_id": userID})
	if err != nil {
		log.E("Failed to count user cocks", logging.InnerError, err)
		return 0
	}

	return int(count)
}

// GetTotalCockersCount возвращает общее количество уникальных участников за все время
func (app *Application) GetTotalCockersCount(log *logging.Logger) int {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineTotalCockersCount())
	if err != nil {
		log.E("Failed to count total cockers", logging.InnerError, err)
		return 0
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var result struct {
		Total int `bson:"total"`
	}
	if cursor.Next(app.ctx) {
		if err := cursor.Decode(&result); err != nil {
			log.E("Failed to decode count result", logging.InnerError, err)
			return 0
		}
		return result.Total
	}

	return 0
}

// GetSeasonCockersCount возвращает количество уникальных участников в сезоне
func (app *Application) GetSeasonCockersCount(log *logging.Logger, season CockSeason) int {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineSeasonCockersCount(season.StartDate, season.EndDate))
	if err != nil {
		log.E("Failed to count season cockers", logging.InnerError, err)
		return 0
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var result struct {
		Total int `bson:"total"`
	}
	if cursor.Next(app.ctx) {
		if err := cursor.Decode(&result); err != nil {
			log.E("Failed to decode count result", logging.InnerError, err)
			return 0
		}
		return result.Total
	}

	return 0
}

func (app *Application) GetUserPositionInLadder(log *logging.Logger, userID int64) int {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineUserPositionInLadder(userID))
	if err != nil {
		log.E("Failed to get user position in ladder", logging.InnerError, err)
		return 0
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var result struct {
		Position int `bson:"position"`
	}
	if cursor.Next(app.ctx) {
		if err := cursor.Decode(&result); err != nil {
			log.E("Failed to decode position result", logging.InnerError, err)
			return 0
		}
		return result.Position
	}

	return 0
}

func (app *Application) GetUserPositionInSeason(log *logging.Logger, userID int64, season CockSeason) int {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineUserPositionInSeason(userID, season.StartDate, season.EndDate))
	if err != nil {
		log.E("Failed to get user position in season", logging.InnerError, err)
		return 0
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var result struct {
		Position int `bson:"position"`
	}
	if cursor.Next(app.ctx) {
		if err := cursor.Decode(&result); err != nil {
			log.E("Failed to decode position result", logging.InnerError, err)
			return 0
		}
		return result.Position
	}

	return 0
}

func (app *Application) GetFirstCockDate(log *logging.Logger) *time.Time {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineFirstCockDate())
	if err != nil {
		log.E("Failed to get first cock date", logging.InnerError, err)
		return nil
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var result struct {
		FirstDate time.Time `bson:"first_date"`
	}
	if cursor.Next(app.ctx) {
		if err := cursor.Decode(&result); err != nil {
			log.E("Failed to decode first cock date", logging.InnerError, err)
			return nil
		}
		return &result.FirstDate
	}

	return nil
}

func (app *Application) GetAllSeasons(log *logging.Logger) []CockSeason {
	firstCockDate := app.GetFirstCockDate(log)
	if firstCockDate == nil {
		log.I("No cocks found in database")
		return []CockSeason{}
	}

	moscowLoc := datetime.NowLocation()

	// Нормализуем дату первого кока к началу дня (00:00:00) в московской локации
	year, month, day := firstCockDate.Year(), firstCockDate.Month(), firstCockDate.Day()
	normalizedFirstDate := time.Date(year, month, day, 0, 0, 0, 0, moscowLoc)

	var seasons []CockSeason
	currentDate := normalizedFirstDate
	seasonNum := 1
	now := datetime.NowTime()

	// Нормализуем текущую дату к началу дня для корректного сравнения
	nowYear, nowMonth, nowDay := now.Year(), now.Month(), now.Day()
	normalizedNow := time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, moscowLoc)

	for currentDate.Before(normalizedNow) || currentDate.Equal(normalizedNow) {
		// Каждый сезон длится 3 месяца
		endDate := currentDate.AddDate(0, 3, 0)
		isActive := (normalizedNow.After(currentDate) || normalizedNow.Equal(currentDate)) && normalizedNow.Before(endDate)

		seasons = append(seasons, CockSeason{
			StartDate: currentDate,
			EndDate:   endDate,
			IsActive:  isActive,
			SeasonNum: seasonNum,
		})

		currentDate = endDate
		seasonNum++
	}

	if len(seasons) > MaxSeasonsToShow {
		seasons = seasons[len(seasons)-MaxSeasonsToShow:]
	}

	return seasons
}

func (app *Application) GetAllSeasonsCount(log *logging.Logger) int {
	firstCockDate := app.GetFirstCockDate(log)
	if firstCockDate == nil {
		return 0
	}

	moscowLoc := datetime.NowLocation()

	// Нормализуем дату первого кока к началу дня (00:00:00) в московской локации
	year, month, day := firstCockDate.Year(), firstCockDate.Month(), firstCockDate.Day()
	normalizedFirstDate := time.Date(year, month, day, 0, 0, 0, 0, moscowLoc)

	currentDate := normalizedFirstDate
	count := 0
	now := datetime.NowTime()

	// Нормализуем текущую дату к началу дня для корректного сравнения
	nowYear, nowMonth, nowDay := now.Year(), now.Month(), now.Day()
	normalizedNow := time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, moscowLoc)

	for currentDate.Before(normalizedNow) || currentDate.Equal(normalizedNow) {
		endDate := currentDate.AddDate(0, 3, 0)
		count++
		currentDate = endDate
	}

	return count
}

func (app *Application) GetUserSeasonWins(log *logging.Logger, userID int64) int {
	seasons := app.GetAllSeasonsForStats(log)
	wins := 0

	for _, season := range seasons {
		if !season.IsActive {
			winners := app.GetSeasonWinners(log, season)
			for _, winner := range winners {
				if winner.UserID == userID && winner.Place <= 3 {
					wins++
					break
				}
			}
		}
	}

	return wins
}

func (app *Application) GetUserCockRespect(log *logging.Logger, userID int64) int {
	// Респекты из сезонов
	seasons := app.GetAllSeasonsForStats(log)
	totalRespect := 0

	for _, season := range seasons {
		if !season.IsActive {
			respect := app.GetUserSeasonRespect(log, userID, season)
			totalRespect += respect
		}
	}

	// Добавляем респекты из достижений
	achievementRespects := app.GetUserAchievementRespects(log, userID)
	totalRespect += achievementRespects

	return totalRespect
}

func (app *Application) GetUserSeasonRespect(log *logging.Logger, userID int64, season CockSeason) int {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineAllUsersInSeason(season.StartDate, season.EndDate))
	if err != nil {
		log.E("Failed to get season ranking for respect calculation", logging.InnerError, err)
		return 0
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var results []UserCockRace
	if err = cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to parse season ranking", logging.InnerError, err)
		return 0
	}

	for position, user := range results {
		if user.UserID == userID {
			place := position + 1
			return CalculateCockRespect(place)
		}
	}

	return 0
}

func CalculateCockRespect(place int) int {
	if place <= 0 {
		return 0 // За то что не вошел в сезончик *trollface.jpg*
	}

	if place == 1 {
		return 1488 // Почему бы и нет, *заслужено.gif*
	}

	// Базовая идея: чем дальше место — тем ниже респект, но не линейно
	score := int(3000 / math.Pow(float64(place), 1.2))

	if score < 1 {
		return 1 // Минимум — 1 очко за выход в сезон
	}

	return score
}

func (app *Application) GetAllSeasonsForStats(log *logging.Logger) []CockSeason {
	firstCockDate := app.GetFirstCockDate(log)
	if firstCockDate == nil {
		log.I("No cocks found in database")
		return []CockSeason{}
	}

	moscowLoc := datetime.NowLocation()

	// Нормализуем дату первого кока к началу дня (00:00:00) в московской локации
	year, month, day := firstCockDate.Year(), firstCockDate.Month(), firstCockDate.Day()
	normalizedFirstDate := time.Date(year, month, day, 0, 0, 0, 0, moscowLoc)

	var seasons []CockSeason
	currentDate := normalizedFirstDate
	seasonNum := 1
	now := datetime.NowTime()

	// Нормализуем текущую дату к началу дня для корректного сравнения
	nowYear, nowMonth, nowDay := now.Year(), now.Month(), now.Day()
	normalizedNow := time.Date(nowYear, nowMonth, nowDay, 0, 0, 0, 0, moscowLoc)

	for currentDate.Before(normalizedNow) || currentDate.Equal(normalizedNow) {
		endDate := currentDate.AddDate(0, 3, 0)
		isActive := (normalizedNow.After(currentDate) || normalizedNow.Equal(currentDate)) && normalizedNow.Before(endDate)

		seasons = append(seasons, CockSeason{
			StartDate: currentDate,
			EndDate:   endDate,
			IsActive:  isActive,
			SeasonNum: seasonNum,
		})

		currentDate = endDate
		seasonNum++
	}

	return seasons
}

func (app *Application) GetCurrentSeason(log *logging.Logger) *CockSeason {
	seasons := app.GetAllSeasons(log)
	for _, season := range seasons {
		if season.IsActive {
			return &season
		}
	}
	return nil
}

func (app *Application) GetSeasonWinners(log *logging.Logger, season CockSeason) []SeasonWinner {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineSeasonWinners(season.StartDate, season.EndDate))
	if err != nil {
		log.E("Failed to get season winners", logging.InnerError, err)
		return []SeasonWinner{}
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var results []SeasonWinner
	if err = cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to parse season winners", logging.InnerError, err)
		return []SeasonWinner{}
	}

	// Добавляем места (1, 2, 3)
	for i := range results {
		results[i].Place = i + 1
	}

	return results
}

func (app *Application) GetUsersAroundPositionInLadder(log *logging.Logger, position int) []UserCockRace {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineUsersAroundPositionInLadder(position))
	if err != nil {
		log.E("Failed to get users around position in ladder", logging.InnerError, err)
		return []UserCockRace{}
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var results []UserCockRace
	if err = cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to parse users around position", logging.InnerError, err)
		return []UserCockRace{}
	}

	return results
}

func (app *Application) GetUsersAroundPositionInSeason(log *logging.Logger, position int, season CockSeason) []UserCockRace {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineUsersAroundPositionInSeason(position, season.StartDate, season.EndDate))
	if err != nil {
		log.E("Failed to get users around position in season", logging.InnerError, err)
		return []UserCockRace{}
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var results []UserCockRace
	if err = cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to parse users around position in season", logging.InnerError, err)
		return []UserCockRace{}
	}

	return results
}

// GetUserAchievements получает все достижения пользователя
func (app *Application) GetUserAchievements(log *logging.Logger, userID int64) map[string]*database.DocumentUserAchievement {
	collection := database.CollectionAchievements(app.db)

	cursor, err := collection.Find(app.ctx, map[string]interface{}{"user_id": userID})
	if err != nil {
		log.E("Failed to get user achievements", logging.InnerError, err)
		return make(map[string]*database.DocumentUserAchievement)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var achievements []*database.DocumentUserAchievement
	if err = cursor.All(app.ctx, &achievements); err != nil {
		log.E("Failed to parse user achievements", logging.InnerError, err)
		return make(map[string]*database.DocumentUserAchievement)
	}

	// Преобразуем в map для быстрого доступа
	achievementMap := make(map[string]*database.DocumentUserAchievement)
	for _, ach := range achievements {
		achievementMap[ach.AchievementID] = ach
	}

	return achievementMap
}

// GetUserAchievementRespects подсчитывает общее количество кок-респектов из достижений
func (app *Application) GetUserAchievementRespects(log *logging.Logger, userID int64) int {
	// Проверка только для тестового пользователя
	// if userID != 362695653 {
	// 	return 0
	// }

	userAchievements := app.GetUserAchievements(log, userID)
	totalRespects := 0

	for achID, userAch := range userAchievements {
		if userAch.Completed {
			ach := GetAchievementByID(achID)
			if ach != nil {
				totalRespects += ach.Respects
			}
		}
	}

	return totalRespects
}

// CheckAndUpdateAchievements проверяет и обновляет достижения пользователя (только для mairwunnx)
func (app *Application) CheckAndUpdateAchievements(log *logging.Logger, userID int64) {
	// Проверка только для тестового пользователя
	// if userID != 362695653 {
	// 	return
	// }

	// Получаем текущие достижения пользователя
	userAchievements := app.GetUserAchievements(log, userID)

	// Проверяем, когда последний раз проверяли достижения
	now := datetime.NowTime()
	for _, ach := range userAchievements {
		if !ach.LastCheckedAt.IsZero() {
			moscowTime := ach.LastCheckedAt.In(datetime.NowLocation())
			todayMoscow := now

			// Если уже проверяли сегодня, выходим
			if moscowTime.Year() == todayMoscow.Year() &&
				moscowTime.Month() == todayMoscow.Month() &&
				moscowTime.Day() == todayMoscow.Day() {
				log.I("Achievements already checked today")
				return
			}
		}
		break // Достаточно проверить одну запись
	}

	collection := database.CollectionCocks(app.db)

	// Запускаем пайплайн проверки
	cursor, err := collection.Aggregate(app.ctx, database.PipelineCheckAchievements(userID))
	if err != nil {
		log.E("Failed to check achievements", logging.InnerError, err)
		return
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var results []database.DocumentAchievementCheck
	if err = cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to parse achievement check results", logging.InnerError, err)
		return
	}

	if len(results) == 0 {
		log.E("No results from achievement check pipeline")
		return
	}

	data := results[0]

	const absoluteMax int32 = 61
	const absoluteMin int32 = 0

	// Обновляем достижения
	achievementCollection := database.CollectionAchievements(app.db)

	// Функция для обновления/создания достижения
	updateAchievement := func(achID string, completed bool, progress int) {
		filter := map[string]interface{}{
			"user_id":        userID,
			"achievement_id": achID,
		}

		existingAch := userAchievements[achID]
		if existingAch != nil && existingAch.Completed {
			// Уже выполнено, не обновляем
			return
		}

		update := map[string]interface{}{
			"$set": map[string]interface{}{
				"user_id":         userID,
				"achievement_id":  achID,
				"completed":       completed,
				"progress":        progress,
				"last_checked_at": now,
			},
		}

		if completed && (existingAch == nil || !existingAch.Completed) {
			update["$set"].(map[string]interface{})["completed_at"] = now
		}

		opts := options.Update().SetUpsert(true)
		_, err := achievementCollection.UpdateOne(app.ctx, filter, update, opts)
		if err != nil {
			log.E("Failed to update achievement", "achievement_id", achID, logging.InnerError, err)
		}
	}

	// Проверяем достижения по количеству дерганий
	if len(data.TotalPulls) > 0 {
		count := data.TotalPulls[0].Count
		updateAchievement("not_rubbed_yet", count >= 10, int(count))
		updateAchievement("diary", count >= 31, int(count))
		updateAchievement("skillful_hands", count >= 100, int(count))
		updateAchievement("anniversary", count >= 365, int(count))
		updateAchievement("wonder_stranger", count >= 500, int(count))
		updateAchievement("bazooka_hands", count >= 1000, int(count))
		updateAchievement("annihilator_cannon", count >= 5000, int(count))
	}

	// Проверяем достижения по накопленному размеру
	if len(data.TotalSize) > 0 {
		total := data.TotalSize[0].Total
		updateAchievement("golden_hundred", total >= 100, int(total))
		updateAchievement("solid_thousand", total >= 1000, int(total))
		updateAchievement("five_k", total >= 5000, int(total))
		updateAchievement("golden_cock", total >= 10000, int(total))
		updateAchievement("cosmic_cock", total >= 20000, int(total))
		updateAchievement("greek_myth", total >= 30000, int(total))
	}

	// Проверяем достижение "Снайпер" (30см 5 раз)
	if len(data.Sniper30cm) > 0 {
		count := data.Sniper30cm[0].Count
		updateAchievement("sniper", count >= 5, int(count))
	}

	// Проверяем достижение "Полсотни" (50см)
	if len(data.HalfHundred50cm) > 0 {
		count := data.HalfHundred50cm[0].Count
		updateAchievement("half_hundred", count >= 1, int(count))
	}

	// Проверяем достижение "Максималист" (61см)
	if len(data.Maximalist61cm) > 0 {
		count := data.Maximalist61cm[0].Count
		updateAchievement("maximalist", count >= 10, int(count))
	}

	// Проверяем "Коллекционер чисел"
	if len(data.BeautifulNumbers) > 0 {
		count := data.BeautifulNumbers[0].Count
		updateAchievement("number_collector", count >= 5, int(count))
	}

	// Проверяем экстремумы (Эверест и Марианская впадина)
	// Эверест - получить максимальный кок среди всех возможных (61см)
	if len(data.MaxSize) > 0 {
		userMax := data.MaxSize[0].Max
		updateAchievement("everest", userMax == absoluteMax, int(userMax))
	}

	// Марианская впадина - получить минимальный кок среди всех возможных (0см)
	if len(data.MinSize) > 0 {
		userMin := data.MinSize[0].Min
		updateAchievement("mariana_trench", userMin == absoluteMin, int(userMin))
	}

	// Проверяем временные достижения
	if len(data.EarlyBird) > 0 {
		count := data.EarlyBird[0].Count
		updateAchievement("early_bird", count >= 20, int(count))
	}

	if len(data.Speedrunner) > 0 {
		count := data.Speedrunner[0].Count
		updateAchievement("speedrunner", count >= 5, int(count))
	}

	if len(data.MidnightPuller) > 0 {
		count := data.MidnightPuller[0].Count
		updateAchievement("midnight_puller", count >= 10, int(count))
	}

	// Проверяем праздничные
	if len(data.Valentine) > 0 {
		count := data.Valentine[0].Count
		updateAchievement("valentine", count >= 1, int(count))
	}

	if len(data.NewYearGift) > 0 {
		count := data.NewYearGift[0].Count
		updateAchievement("new_year_gift", count >= 1, int(count))
	}

	if len(data.MensSolidarity) > 0 {
		count := data.MensSolidarity[0].Count
		updateAchievement("mens_solidarity", count >= 1, int(count))
	}

	if len(data.Friday13th) > 0 {
		count := data.Friday13th[0].Count
		updateAchievement("friday_13th", count >= 1, int(count))
	}

	if len(data.LeapCock) > 0 {
		count := data.LeapCock[0].Count
		updateAchievement("leap_cock", count >= 1, int(count))
	}

	// Проверяем "Молния"
	if len(data.Lightning) > 0 {
		count := data.Lightning[0].Count
		updateAchievement("lightning", count >= 1, int(count))
	}

	// Проверяем последовательности в последних 10 коках
	if len(data.Recent10) > 0 {
		sizes := make([]int32, 0, len(data.Recent10))
		for _, item := range data.Recent10 {
			sizes = append(sizes, item.Size)
		}

		// Проверяем последовательности одинаковых
		if len(sizes) >= 2 {
			updateAchievement("deja_vu", sizes[len(sizes)-1] == sizes[len(sizes)-2], 0)
		}

		// Проверяем тройки, покер, глаз алмаз
		for i := 0; i <= len(sizes)-3; i++ {
			if sizes[i] == sizes[i+1] && sizes[i+1] == sizes[i+2] {
				updateAchievement("triple", true, 0)
				if i <= len(sizes)-4 && sizes[i+2] == sizes[i+3] {
					updateAchievement("poker", true, 0)
					if i <= len(sizes)-5 && sizes[i+3] == sizes[i+4] {
						updateAchievement("diamond_eye", true, 0)
					}
				}
			}
		}

		// Проверяем тренды (рост/падение 5 дней)
		if len(sizes) >= 5 {
			allGrowth := true
			allDecline := true
			for i := 1; i < 5; i++ {
				if sizes[i] <= sizes[i-1] {
					allGrowth = false
				}
				if sizes[i] >= sizes[i-1] {
					allDecline = false
				}
			}
			updateAchievement("bull_trend", allGrowth, 0)
			updateAchievement("bear_market", allDecline, 0)
		}

		// Проверяем "Мороз по коже" (5 коков <20см)
		if len(sizes) >= 5 {
			allFrozen := true
			for i := 0; i < 5; i++ {
				if sizes[len(sizes)-5+i] >= 20 {
					allFrozen = false
					break
				}
			}
			updateAchievement("freeze", allFrozen, 0)
		}

		// Проверяем "Алмазные руки" (7 коков 40+см)
		if len(sizes) >= 7 {
			allDiamond := true
			for i := 0; i < 7; i++ {
				if sizes[len(sizes)-7+i] < 40 {
					allDiamond = false
					break
				}
			}
			updateAchievement("diamond_hands", allDiamond, 0)
		}

		// Проверяем "Черепаха" (10 коков с изменением <5см)
		if len(sizes) >= 10 {
			turtleCount := 0
			for i := 1; i < 10; i++ {
				diff := sizes[i] - sizes[i-1]
				if diff < 0 {
					diff = -diff
				}
				if diff < 5 {
					turtleCount++
				}
			}
			updateAchievement("turtle", turtleCount >= 9, turtleCount)
		}
	}

	// Проверяем сложные коллекции в последних 31 коке
	if len(data.Last31) >= 31 {
		sizes := make(map[int32]bool)
		for _, item := range data.Last31 {
			sizes[item.Size] = true
		}

		// Проверяем "Округлятор"
		rounderCount := 0
		for _, val := range []int32{10, 20, 30, 40, 50, 60} {
			if sizes[val] {
				rounderCount++
			}
		}
		updateAchievement("rounder", rounderCount == 6, rounderCount)

		// Проверяем "Отец фибоначчи"
		fibonacciCount := 0
		for _, val := range []int32{1, 2, 3, 5, 8, 13, 21, 34, 55} {
			if sizes[val] {
				fibonacciCount++
			}
		}
		updateAchievement("fibonacci_father", fibonacciCount == 9, fibonacciCount)
	}

	// Проверяем достижения по сезонам
	seasonCursor, err := collection.Aggregate(app.ctx, database.PipelineCountSeasons(userID))
	if err != nil {
		log.E("Failed to get season count", logging.InnerError, err)
	} else {
		defer seasonCursor.Close(app.ctx)

		var seasonResults []database.DocumentSeasonCount
		if err = seasonCursor.All(app.ctx, &seasonResults); err == nil && len(seasonResults) > 0 {
			count := seasonResults[0].Count
			updateAchievement("oldtimer", count >= 3, int(count))
			updateAchievement("veteran", count >= 5, int(count))
			updateAchievement("keeper", count >= 10, int(count))
		}
	}

	// Проверяем достижение "Путешественник" (все 61 размер: 0-60)
	travelerCursor, err := collection.Aggregate(app.ctx, database.PipelineCheckTraveler(userID))
	if err != nil {
		log.E("Failed to check traveler achievement", logging.InnerError, err)
	} else {
		defer travelerCursor.Close(app.ctx)

		var travelerResults []database.DocumentTravelerCheck
		if err = travelerCursor.All(app.ctx, &travelerResults); err == nil && len(travelerResults) > 0 {
			uniqueSizes := travelerResults[0].UniqueSizes
			updateAchievement("traveler", uniqueSizes >= 61, int(uniqueSizes))
		}
	}

	// Проверяем достижение "Москвич" (размер 50см 5 раз за последние 31 день)
	thirtyOneDaysAgo := now.AddDate(0, 0, -31)
	muscoviteCursor, err := collection.Aggregate(app.ctx, database.PipelineCheckMuscovite(userID, thirtyOneDaysAgo))
	if err != nil {
		log.E("Failed to check muscovite achievement", logging.InnerError, err)
	} else {
		defer muscoviteCursor.Close(app.ctx)

		var muscoviteResults []database.DocumentMuscoviteCheck
		if err = muscoviteCursor.All(app.ctx, &muscoviteResults); err == nil && len(muscoviteResults) > 0 {
			count := muscoviteResults[0].Count
			updateAchievement("muscovite", count >= 5, int(count))
		}
	}

	// Проверяем специальные совпадения в последних 3 коках
	if len(data.Recent3) >= 3 {
		recent := data.Recent3

		// 1. Сумма предыдущих: recent[2] == recent[0] + recent[1]
		if recent[2].Size == recent[0].Size+recent[1].Size {
			updateAchievement("sum_of_previous", true, 0)
		}

		// 2. Контрастный душ: после 60+ получить 0-3
		if recent[1].Size >= 60 && recent[2].Size <= 3 {
			updateAchievement("contrast_shower", true, 0)
		}

		// 3. Пифагор: три кока подряд образуют пифагорову тройку
		a, b, c := recent[0].Size, recent[1].Size, recent[2].Size
		// Проверяем все варианты перестановок (a²+b²=c², a²+c²=b², b²+c²=a²)
		isPythagorean := (a*a+b*b == c*c) || (a*a+c*c == b*b) || (b*b+c*c == a*a)
		if isPythagorean {
			updateAchievement("pythagoras", true, 1)
		}

		// 4. Leet speak (1337): 13см и 37см подряд
		if recent[1].Size == 13 && recent[2].Size == 37 {
			updateAchievement("leet_speak", true, 1)
		}

		// Проверяем последний кок на совпадения с временем
		lastCock := recent[2]
		moscowTime := lastCock.RequestedAt.In(datetime.NowLocation())

		hour := moscowTime.Hour()
		minute := moscowTime.Minute()
		day := moscowTime.Day()

		// 5. Минутная точность: размер == минуты (например 24см в xx:24)
		if int32(minute) == lastCock.Size {
			updateAchievement("minute_precision", true, 0)
		}

		// 6. Часовая точность: размер == час (например 11см в 11:xx)
		if int32(hour) == lastCock.Size {
			updateAchievement("hour_precision", true, 0)
		}

		// 7. День = Размер: размер == день месяца (например 15см 15 числа)
		if int32(day) == lastCock.Size {
			updateAchievement("day_equals_size", true, 0)
		}
	}

	log.I("Successfully checked and updated achievements")
}
