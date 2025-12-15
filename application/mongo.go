package application

import (
	"context"
	"dickobrazz/application/database"
	"dickobrazz/application/logging"
	"math"
	"os"
	"time"

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

func (app *Application) AggregateCockSizes(log *logging.Logger) []UserCockRace {
	collection := database.CollectionCocks(app.db)

	cursor, err := collection.Aggregate(app.ctx, database.PipelineTopUsersBySize())
	if err != nil {
		log.E("Failed to aggregate cock sizes", logging.InnerError, err)
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
	
	var seasons []CockSeason
	currentDate := *firstCockDate
	seasonNum := 1
	now := time.Now()
	
	for currentDate.Before(now) {
		// Каждый сезон длится 3 месяца
		endDate := currentDate.AddDate(0, 3, 0)
		isActive := now.After(currentDate) && now.Before(endDate)
		
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
	
	currentDate := *firstCockDate
	count := 0
	now := time.Now()
	
	for currentDate.Before(now) {
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
	
	var seasons []CockSeason
	currentDate := *firstCockDate
	seasonNum := 1
	now := time.Now()
	
	for currentDate.Before(now) {
		endDate := currentDate.AddDate(0, 3, 0)
		isActive := now.After(currentDate) && now.Before(endDate)
		
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
	if userID != 362695653 {
		return 0
	}
	
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
	if userID != 362695653 {
		return
	}

	// Получаем текущие достижения пользователя
	userAchievements := app.GetUserAchievements(log, userID)
	log.I("Ach: Current user achievements count", "count", len(userAchievements))

	// Проверяем, когда последний раз проверяли достижения
	now := time.Now()
	// for _, ach := range userAchievements {
	// 	if !ach.LastCheckedAt.IsZero() {
	// 		moscowTime := ach.LastCheckedAt.In(time.FixedZone("MSK", 3*60*60))
	// 		todayMoscow := now.In(time.FixedZone("MSK", 3*60*60))
			
	// 		log.I("Ach: Last check time", "last_checked", moscowTime, "today", todayMoscow, "achievement_id", ach.AchievementID)
			
	// 		// Если уже проверяли сегодня, выходим
	// 		if moscowTime.Year() == todayMoscow.Year() &&
	// 			moscowTime.Month() == todayMoscow.Month() &&
	// 			moscowTime.Day() == todayMoscow.Day() {
	// 			log.I("Achievements already checked today")
	// 			return
	// 		}
	// 	}
	// 	break // Достаточно проверить одну запись
	// }
	
	log.I("Ach: Starting achievements check")

	// Запускаем пайплайн проверки
	collection := database.CollectionCocks(app.db)
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

	var results []map[string]interface{}
	if err = cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to parse achievement check results", logging.InnerError, err)
		return
	}

	if len(results) == 0 {
		log.E("No results from achievement check pipeline")
		return
	}

	data := results[0]

	// Получаем глобальные максимум и минимум
	globalCursor, err := collection.Aggregate(app.ctx, database.PipelineGlobalMaxMin())
	if err != nil {
		log.E("Failed to get global max/min", logging.InnerError, err)
		return
	}
	defer globalCursor.Close(app.ctx)

	var globalResults []map[string]interface{}
	if err = globalCursor.All(app.ctx, &globalResults); err != nil {
		log.E("Failed to parse global max/min", logging.InnerError, err)
		return
	}

	var globalMax, globalMin int32
	if len(globalResults) > 0 {
		if max, ok := globalResults[0]["max"].(int32); ok {
			globalMax = max
		}
		if min, ok := globalResults[0]["min"].(int32); ok {
			globalMin = min
		}
	}

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
			log.I("Ach: Achievement already completed, skipping", "achievement_id", achID)
			return
		}

		log.I("Ach: Updating achievement", "achievement_id", achID, "completed", completed, "progress", progress)

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
			log.I("Ach: Achievement completed!", "achievement_id", achID)
		}

		opts := options.Update().SetUpsert(true)
		_, err := achievementCollection.UpdateOne(app.ctx, filter, update, opts)
		if err != nil {
			log.E("Failed to update achievement", "achievement_id", achID, logging.InnerError, err)
		}
	}

	// Проверяем достижения по количеству дерганий
	if totalPulls, ok := data["total_pulls"].([]interface{}); ok && len(totalPulls) > 0 {
		if pullData, ok := totalPulls[0].(map[string]interface{}); ok {
			if count, ok := pullData["count"].(int32); ok {
				log.I("Ach: Total pulls from pipeline", "count", count)
				updateAchievement("not_rubbed_yet", count >= 10, int(count))
				updateAchievement("diary", count >= 31, int(count))
				updateAchievement("skillful_hands", count >= 100, int(count))
				updateAchievement("anniversary", count >= 365, int(count))
				updateAchievement("wonder_stranger", count >= 500, int(count))
				updateAchievement("bazooka_hands", count >= 1000, int(count))
				updateAchievement("annihilator_cannon", count >= 5000, int(count))
			}
		}
	} else {
		log.W("Ach: No total_pulls data in pipeline results")
	}

	// Проверяем достижения по накопленному размеру
	if totalSize, ok := data["total_size"].([]interface{}); ok && len(totalSize) > 0 {
		if sizeData, ok := totalSize[0].(map[string]interface{}); ok {
			if total, ok := sizeData["total"].(int32); ok {
				log.I("Ach: Total size from pipeline", "total", total)
				updateAchievement("golden_hundred", total >= 100, int(total))
				updateAchievement("solid_thousand", total >= 1000, int(total))
				updateAchievement("five_k", total >= 5000, int(total))
				updateAchievement("golden_cock", total >= 10000, int(total))
				updateAchievement("cosmic_cock", total >= 20000, int(total))
				updateAchievement("greek_myth", total >= 30000, int(total))
			}
		}
	} else {
		log.W("Ach: No total_size data in pipeline results")
	}

	// Проверяем достижение "Снайпер" (30см 5 раз)
	if sniper, ok := data["sniper_30cm"].([]interface{}); ok && len(sniper) > 0 {
		if sniperData, ok := sniper[0].(map[string]interface{}); ok {
			if count, ok := sniperData["count"].(int32); ok {
				updateAchievement("sniper", count >= 5, int(count))
			}
		}
	}

	// Проверяем достижение "Полсотни" (50см)
	if halfHundred, ok := data["half_hundred_50cm"].([]interface{}); ok && len(halfHundred) > 0 {
		if hhData, ok := halfHundred[0].(map[string]interface{}); ok {
			if count, ok := hhData["count"].(int32); ok {
				updateAchievement("half_hundred", count >= 1, int(count))
			}
		}
	}

	// Проверяем "Коллекционер чисел"
	if beautifulNumbers, ok := data["beautiful_numbers"].([]interface{}); ok && len(beautifulNumbers) > 0 {
		if bnData, ok := beautifulNumbers[0].(map[string]interface{}); ok {
			if count, ok := bnData["count"].(int32); ok {
				updateAchievement("number_collector", count >= 5, int(count))
			}
		}
	}

	// Проверяем экстремумы (Эверест и Марианская впадина)
	if maxSize, ok := data["max_size"].([]interface{}); ok && len(maxSize) > 0 {
		if maxData, ok := maxSize[0].(map[string]interface{}); ok {
			if userMax, ok := maxData["max"].(int32); ok {
				updateAchievement("everest", userMax == globalMax && globalMax > 0, int(userMax))
			}
		}
	}

	if minSize, ok := data["min_size"].([]interface{}); ok && len(minSize) > 0 {
		if minData, ok := minSize[0].(map[string]interface{}); ok {
			if userMin, ok := minData["min"].(int32); ok {
				updateAchievement("mariana_trench", userMin == globalMin && globalMin >= 0, int(userMin))
			}
		}
	}

	// Проверяем временные достижения
	if earlyBird, ok := data["early_bird"].([]interface{}); ok && len(earlyBird) > 0 {
		if ebData, ok := earlyBird[0].(map[string]interface{}); ok {
			if count, ok := ebData["count"].(int32); ok {
				updateAchievement("early_bird", count >= 20, int(count))
			}
		}
	}

	if speedrunner, ok := data["speedrunner"].([]interface{}); ok && len(speedrunner) > 0 {
		if srData, ok := speedrunner[0].(map[string]interface{}); ok {
			if count, ok := srData["count"].(int32); ok {
				updateAchievement("speedrunner", count >= 5, int(count))
			}
		}
	}

	// Проверяем праздничные
	if valentine, ok := data["valentine"].([]interface{}); ok && len(valentine) > 0 {
		if vData, ok := valentine[0].(map[string]interface{}); ok {
			if count, ok := vData["count"].(int32); ok {
				updateAchievement("valentine", count >= 1, int(count))
			}
		}
	}

	if newYearGift, ok := data["new_year_gift"].([]interface{}); ok && len(newYearGift) > 0 {
		if nygData, ok := newYearGift[0].(map[string]interface{}); ok {
			if count, ok := nygData["count"].(int32); ok {
				updateAchievement("new_year_gift", count >= 1, int(count))
			}
		}
	}

	// Проверяем "Молния"
	if lightning, ok := data["lightning"].([]interface{}); ok && len(lightning) > 0 {
		if lData, ok := lightning[0].(map[string]interface{}); ok {
			if count, ok := lData["count"].(int32); ok {
				updateAchievement("lightning", count >= 1, int(count))
			}
		}
	}

	// Проверяем последовательности в последних 10 коках
	if recent10, ok := data["recent_10"].([]interface{}); ok && len(recent10) > 0 {
		sizes := make([]int32, 0, len(recent10))
		for _, item := range recent10 {
			if doc, ok := item.(map[string]interface{}); ok {
				if size, ok := doc["size"].(int32); ok {
					sizes = append(sizes, size)
				}
			}
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
	if last31, ok := data["last_31"].([]interface{}); ok && len(last31) >= 31 {
		sizes := make(map[int32]bool)
		for _, item := range last31 {
			if doc, ok := item.(map[string]interface{}); ok {
				if size, ok := doc["size"].(int32); ok {
					sizes[size] = true
				}
			}
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
		
		var seasonResults []map[string]interface{}
		if err = seasonCursor.All(app.ctx, &seasonResults); err == nil && len(seasonResults) > 0 {
			if count, ok := seasonResults[0]["count"].(int32); ok {
				updateAchievement("oldtimer", count >= 3, int(count))
				updateAchievement("veteran", count >= 5, int(count))
				updateAchievement("keeper", count >= 10, int(count))
			}
		}
	}

	// Проверяем достижение "Путешественник" (все 61 размер: 0-60)
	travelerCursor, err := collection.Aggregate(app.ctx, database.PipelineCheckTraveler(userID))
	if err != nil {
		log.E("Failed to check traveler achievement", logging.InnerError, err)
	} else {
		defer travelerCursor.Close(app.ctx)
		
		var travelerResults []map[string]interface{}
		if err = travelerCursor.All(app.ctx, &travelerResults); err == nil && len(travelerResults) > 0 {
			if uniqueSizes, ok := travelerResults[0]["unique_sizes"].(int32); ok {
				updateAchievement("traveler", uniqueSizes >= 61, int(uniqueSizes))
			}
		}
	}

	// Проверяем достижение "Москвич" (размер 50см 5 раз за последние 31 день)
	thirtyOneDaysAgo := now.AddDate(0, 0, -31)
	muscoviteCursor, err := collection.Aggregate(app.ctx, database.PipelineCheckMuscovite(userID, thirtyOneDaysAgo))
	if err != nil {
		log.E("Failed to check muscovite achievement", logging.InnerError, err)
	} else {
		defer muscoviteCursor.Close(app.ctx)
		
		var muscoviteResults []map[string]interface{}
		if err = muscoviteCursor.All(app.ctx, &muscoviteResults); err == nil && len(muscoviteResults) > 0 {
			if count, ok := muscoviteResults[0]["count"].(int32); ok {
				updateAchievement("muscovite", count >= 5, int(count))
			}
		}
	}

	log.I("Successfully checked and updated achievements")
}