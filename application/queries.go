package application

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

func (app *Application) HandleInlineQuery(log *Logger, query *tgbotapi.InlineQuery) {
	results := []any{
		app.InlineQueryCockSize(log, query),
		app.InlineQueryCockRace(log, query),
		app.InlineQueryCockRuler(log, query),
	}

	inlines := tgbotapi.InlineConfig{
		InlineQueryID: query.ID,
		IsPersonal:    true,
		CacheTime:     60,
		Results:       results,
	}

	if _, err := app.bot.Request(inlines); err != nil {
		app.log.E("Failed to send inline query", err)
	} else {
		app.log.I("Inline query successfully sent")
	}
}

func (app *Application) InlineQueryCockSize(log *Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	var size int

	if cached := app.GetCockSizeFromCache(log, query.From.ID); cached != nil {
		size = *cached
	} else {
		size = rand.Intn(61)

		cock := &Cock{
			ID:          uuid.NewString(),
			Size:        int32(size),
			Nickname:    query.From.UserName,
			UserID:      query.From.ID,
			RequestedAt: time.Now(),
		}

		app.SaveCockToCache(log, query.From.ID, size)
		app.SaveCockToMongo(log, cock)
	}

	emoji := EmojiFromSize(size)
	text := fmt.Sprintf("Мой кок: *%dсм* %s", size, emoji)
	return InitializeInlineQuery("Размер кока", text)
}

func (app *Application) InlineQueryCockRace(log *Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	cocks := app.AggregateCockSizes(log)
	text := app.GenerateCockRaceScoreboard(log, query.From.ID, cocks)
	return InitializeInlineQuery("Гонка коков", text)
}

func (app *Application) InlineQueryCockRuler(log *Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	cocks := app.GetCockSizesFromCache(log)

	sort.Slice(cocks, func(i, j int) bool {
		return cocks[i].size > cocks[j].size
	})

	text := app.GenerateCockRulerText(log, query.From.ID, cocks)
	return InitializeInlineQuery("Линейка коков", text)
}

func InitializeInlineQuery(title, message string) tgbotapi.InlineQueryResultArticle {
	return tgbotapi.NewInlineQueryResultArticleMarkdownV2(uuid.NewString(), title, message)
}

func (app *Application) GetCockSizeFromCache(log *Logger, userID int64) *int {
	key := GetCockCacheKey(userID)

	var cock UserCock
	if err := app.cache.Get(app.ctx, key, &cock); err != nil && err == redis.Nil {
		return nil
	} else {
		log.E("Failed to get cock size from redis", InnerError, err)
	}

	log.I("Successfully fetched cock from redis")
	return &cock.size
}

func (app *Application) SaveCockToCache(log *Logger, userID int64, size int) {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	ttl := time.Until(midnight)

	if err := app.cache.Set(&cache.Item{Ctx: app.ctx, Key: GetCockCacheKey(userID), Value: &UserCock{userId: userID, size: size}, TTL: ttl}); err != nil {
		log.E("Failed to save cock size to Redis", InnerError, err)
	} else {
		log.I("Successfully saved cock to redis")
	}
}

func (app *Application) SaveCockToMongo(log *Logger, cock *Cock) {
	collection := app.db.Database("dickbot_db").Collection("cocks")

	if _, err := collection.InsertOne(app.ctx, cock); err != nil {
		log.E("Failed to save cock to MongoDB", InnerError, err)
	} else {
		log.I("Successfully saved cock to MongoDB")
	}
}

func GetCockCacheKey(userID int64) string {
	return fmt.Sprintf("cock_size:%d", userID)
}

func (app *Application) AggregateCockSizes(log *Logger) []UserCockRace {
	collection := app.db.Database("dickbot_db").Collection("cocks")

	pipeline := mongo.Pipeline{
		{
			{Key: "$group", Value: bson.D{
				{Key: "_id", Value: "$user_id"},
				{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
				{Key: "nickname", Value: bson.D{{Key: "$first", Value: "$nickname"}}},
			}},
		},
		{{Key: "$sort", Value: bson.D{{Key: "total_size", Value: -1}}}},
		{{Key: "$limit", Value: 13}},
	}

	cursor, err := collection.Aggregate(app.ctx, pipeline)
	if err != nil {
		log.E("Failed to aggregate cock sizes", InnerError, err)
		panic(err)
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil { // Кааааак же похуй...
			log.E("Failed to close mongo cursor", InnerError, err)
		}
	}(cursor, app.ctx)

	var results []UserCockRace
	if err = cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to parse aggregation results", InnerError, err)
		panic(err)
	}

	log.I("Successfully aggregated cock sizes")
	return results
}

func (app *Application) GenerateCockRaceScoreboard(log *Logger, userID int64, sizes []UserCockRace) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, user := range sizes {
		isCurrentUser := user.UserID == userID
		emoji := GetPlaceEmoji(index + 1)

		if isCurrentUser {
			isUserInScoreboard = true
		}

		var scoreboardLine string
		if isCurrentUser {
			scoreboardLine = fmt.Sprintf("➡️ %s %d. @%s %d", emoji, index+1, user.Nickname, user.TotalSize)
		} else {
			scoreboardLine = fmt.Sprintf("%s %d. @%s %d", emoji, index+1, user.Nickname, user.TotalSize)
		}

		if index < 3 {
			winners = append(winners, scoreboardLine)
		} else {
			others = append(others, scoreboardLine)
		}
	}

	if !isUserInScoreboard {
		if cock := app.GetUserAggregatedCock(log, userID); cock != nil {
			others = append(others, fmt.Sprintf("➡️ 🥀 @%s %d", cock.Nickname, cock.TotalSize))
		} else {
			others = append(others, "➡️ 🥀 **Тебе соболезнуем... потому что не смотрел на кок!**")
		}
	}

	return fmt.Sprintf(
		`**Участники гонки коков:**

Победители в номинации:
%s

🥀 Остальным, соболезнуем:
%s

*В гонке коков легко участвовать, просто запрашивай свой кок каждый день, все коки сбрасываются каждый день по МСК*`,
		strings.Join(winners, "\n"),
		strings.Join(others, "\n"),
	)
}

func (app *Application) GetUserAggregatedCock(log *Logger, userID int64) *UserCockRace {
	collection := app.db.Database("dickbot_db").Collection("cocks")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "nickname", Value: bson.D{{Key: "$first", Value: "$nickname"}}},
		}}},
	}

	cursor, err := collection.Aggregate(app.ctx, pipeline)
	if err != nil {
		log.E("Failed to aggregate user cock sizes", InnerError, err)
		return nil
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", InnerError, err)
		}
	}(cursor, app.ctx)

	var result UserCockRace
	if cursor.Next(app.ctx) {
		if err := cursor.Decode(&result); err != nil {
			log.E("Failed to decode aggregation result", InnerError, err)
			return nil
		}
		return &result
	}

	log.I("No cocks found for user")
	return nil
}

func (app *Application) GetCockSizesFromCache(log *Logger) []UserCock {
	var cockSizes []UserCock

	iter := app.redis.Scan(app.ctx, 0, "cock_size:*", 0).Iterator()
	for iter.Next(app.ctx) {
		key := iter.Val()
		var cock UserCock
		if err := app.redis.Get(app.ctx, key).Scan(&cock); err != nil {
			if err == redis.Nil {
				continue
			}
			log.E("Failed to fetch cock from Redis", InnerError, err)
			panic(err)
		}

		userID, _ := strconv.ParseInt(strings.TrimPrefix(key, "cock_size:"), 10, 64)
		cock.userId = userID
		cockSizes = append(cockSizes, cock)
	}

	if err := iter.Err(); err != nil {
		log.E("Failed to iterate over cock keys in Redis", InnerError, err)
		panic(err)
	}

	log.I("Successfully fetched all cock sizes from Redis")
	return cockSizes
}

func (app *Application) GenerateCockRulerText(log *Logger, userID int64, cocks []UserCock) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, cock := range cocks {
		isCurrentUser := cock.userId == userID
		emoji := GetPlaceEmoji(index + 1)
		line := fmt.Sprintf("%s %d. @%s %dсм", emoji, index+1, cock.userName, cock.size)

		if isCurrentUser {
			isUserInScoreboard = true
			line = fmt.Sprintf("➡️ %s", line)
		}

		if index < 10 {
			winners = append(winners, line)
		} else {
			others = append(others, line)
		}
	}

	if !isUserInScoreboard {
		if userCock := app.GetCockSizeFromCache(log, userID); userCock != nil {
			others = append(others, fmt.Sprintf("➡️ 🥀 @%d %dсм", userID, *userCock))
		} else {
			others = append(others, "➡️ 🥀 **Тебе соболезнуем... потому что не смотрел на кок!**")
		}
	}

	return fmt.Sprintf(
		`**Линейка коков:**

Победители в номинации:
%s

🥀 Остальным, соболезнуем:
%s

*Линейка коков -- чистый рандом, сегодня ты бог, завтра ты лох. Все коки сбрасываются каждые сутки!*`,
		strings.Join(winners, "\n"),
		strings.Join(others, "\n"),
	)
}

func GetPlaceEmoji(place int) string {
	switch place {
	case 1:
		return "🥇"
	case 2:
		return "🥈"
	case 3:
		return "🥉"
	default:
		return "🤧"
	}
}
