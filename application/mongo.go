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
	seasons := app.GetAllSeasonsForStats(log)
	totalRespect := 0
	
	for _, season := range seasons {
		if !season.IsActive {
			respect := app.GetUserSeasonRespect(log, userID, season)
			totalRespect += respect
		}
	}
	
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