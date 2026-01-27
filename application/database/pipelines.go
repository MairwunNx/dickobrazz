package database

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func PipelineDynamic(userId int64) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$facet", Value: bson.D{
			{Key: "individual_irk", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "user_total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
				}}},
				bson.D{{Key: "$lookup", Value: bson.D{
					{Key: "from", Value: "cocks"},
					{Key: "pipeline", Value: bson.A{
						bson.D{{Key: "$group", Value: bson.D{
							{Key: "_id", Value: nil},
							{Key: "global_total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
						}}},
					}},
					{Key: "as", Value: "global_data"},
				}}},
				bson.D{{Key: "$unwind", Value: "$global_data"}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "irk", Value: bson.D{{Key: "$round", Value: bson.A{
						bson.D{{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$lte", Value: bson.A{"$global_data.global_total_size", 0}}},
							0,
							bson.D{{Key: "$divide", Value: bson.A{
								bson.D{{Key: "$log10", Value: bson.D{{Key: "$add", Value: bson.A{1, "$user_total_size"}}}}},
								bson.D{{Key: "$log10", Value: bson.D{{Key: "$add", Value: bson.A{1, "$global_data.global_total_size"}}}}},
							}}},
						}}},
						3,
					}}}},
				}}},
			}},
			{Key: "individual_dominance", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "total_cock", Value: bson.D{{Key: "$sum", Value: "$size"}}},
					{Key: "total_user_cock", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$user_id", userId}}},
						"$size",
						0,
					}}}}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "dominance", Value: bson.D{{Key: "$round", Value: bson.A{
						bson.D{{Key: "$multiply", Value: bson.A{
							bson.D{{Key: "$cond", Value: bson.A{
								bson.D{{Key: "$eq", Value: bson.A{"$total_cock", 0}}},
								0,
								bson.D{{Key: "$divide", Value: bson.A{"$total_user_cock", "$total_cock"}}},
							}}},
							100,
						}}},
						1,
					}}}},
				}}},
			}},
			{Key: "individual_daily_growth", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 5}},
				bson.D{{Key: "$setWindowFields", Value: bson.D{
					{Key: "partitionBy", Value: "$user_id"},
					{Key: "sortBy", Value: bson.D{{Key: "requested_at", Value: -1}}},
					{Key: "output", Value: bson.D{
						{Key: "prev_size", Value: bson.D{{Key: "$shift", Value: bson.D{
							{Key: "output", Value: "$size"},
							{Key: "by", Value: 1},
						}}}},
					}},
				}}},
				bson.D{{Key: "$set", Value: bson.D{
					{Key: "growth", Value: bson.D{{Key: "$round", Value: bson.A{
						bson.D{{Key: "$subtract", Value: bson.A{"$size", "$prev_size"}}},
						1,
					}}}},
				}}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$user_id"},
					{Key: "average_daily_growth", Value: bson.D{{Key: "$avg", Value: "$growth"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "average", Value: bson.D{{Key: "$round", Value: bson.A{"$average_daily_growth", 1}}}},
				}}},
			}},
			{Key: "individual_daily_dynamics", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "requested_at", Value: 1},
					{Key: "size", Value: 1},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 2}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "curr_cock", Value: bson.D{{Key: "$first", Value: "$size"}}},
					{Key: "prev_cock", Value: bson.D{{Key: "$last", Value: "$size"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "yesterday_cock_change", Value: bson.D{{Key: "$subtract", Value: bson.A{"$curr_cock", "$prev_cock"}}}},
					{Key: "yesterday_cock_change_percent", Value: bson.D{{Key: "$round", Value: bson.A{
						bson.D{{Key: "$multiply", Value: bson.A{
							bson.D{{Key: "$divide", Value: bson.A{
								bson.D{{Key: "$subtract", Value: bson.A{"$curr_cock", "$prev_cock"}}},
								bson.D{{Key: "$max", Value: bson.A{"$prev_cock", 1}}},
							}}},
							100,
						}}},
						1,
					}}}},
				}}},
			}},
			{Key: "individual_cock_total", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "total", Value: bson.D{{Key: "$sum", Value: "$size"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "total", Value: 1},
				}}},
			}},
			{Key: "individual_cock_recent", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 5}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "avg_val", Value: bson.D{{Key: "$avg", Value: "$size"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "average", Value: bson.D{{Key: "$round", Value: bson.A{"$avg_val", 0}}}},
				}}},
			}},
			{Key: "overall", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
				}}},
			}},
			{Key: "overall_recent", Value: bson.A{
				// Группируем по пользователю и дню (в московском часовом поясе)
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: bson.D{
						{Key: "user_id", Value: "$user_id"},
						{Key: "year", Value: bson.D{
							{Key: "$year", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
						{Key: "month", Value: bson.D{
							{Key: "$month", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
						{Key: "day", Value: bson.D{
							{Key: "$dayOfMonth", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
					}},
					{Key: "last_cock_of_day", Value: bson.D{{Key: "$last", Value: "$size"}}},
					{Key: "day_timestamp", Value: bson.D{{Key: "$last", Value: "$requested_at"}}},
				}}},
				// Сортируем по дате (от новых к старым)
				bson.D{{Key: "$sort", Value: bson.D{{Key: "day_timestamp", Value: -1}}}},
				// Группируем по пользователю и берем последние 5 дней
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$_id.user_id"},
					{Key: "last_5_days", Value: bson.D{{Key: "$push", Value: bson.D{
						{Key: "size", Value: "$last_cock_of_day"},
						{Key: "timestamp", Value: "$day_timestamp"},
					}}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "last_5_days", Value: bson.D{{Key: "$slice", Value: bson.A{"$last_5_days", 5}}}},
				}}},
				// Разворачиваем массив для расчета статистики
				bson.D{{Key: "$unwind", Value: "$last_5_days"}},
				// Вычисляем среднее и медиану
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "median", Value: bson.D{{Key: "$median", Value: bson.D{
						{Key: "input", Value: "$last_5_days.size"},
						{Key: "method", Value: "approximate"},
					}}}},
					{Key: "average", Value: bson.D{{Key: "$avg", Value: "$last_5_days.size"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "median", Value: 1},
					{Key: "average", Value: bson.D{{Key: "$round", Value: bson.A{"$average", 0}}}},
				}}},
			}},
			{Key: "uniques", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$user_id"}}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "distribution", Value: bson.A{
				// Группируем по пользователю и дню (в московском часовом поясе)
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: bson.D{
						{Key: "user_id", Value: "$user_id"},
						{Key: "year", Value: bson.D{
							{Key: "$year", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
						{Key: "month", Value: bson.D{
							{Key: "$month", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
						{Key: "day", Value: bson.D{
							{Key: "$dayOfMonth", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
					}},
					{Key: "last_cock_of_day", Value: bson.D{{Key: "$last", Value: "$size"}}},
					{Key: "day_timestamp", Value: bson.D{{Key: "$last", Value: "$requested_at"}}},
				}}},
				// Сортируем по дате (от новых к старым)
				bson.D{{Key: "$sort", Value: bson.D{{Key: "day_timestamp", Value: -1}}}},
				// Группируем по пользователю и берем последние 5 дней
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$_id.user_id"},
					{Key: "last_5_days", Value: bson.D{{Key: "$push", Value: "$last_cock_of_day"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "last_5_days", Value: bson.D{{Key: "$slice", Value: bson.A{"$last_5_days", 5}}}},
				}}},
				// Разворачиваем массив
				bson.D{{Key: "$unwind", Value: "$last_5_days"}},
				// Считаем распределение
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "huge", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$gte", Value: bson.A{"$last_5_days", 19}}}, 1, 0,
					}}}}}},
					{Key: "little", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$lt", Value: bson.A{"$last_5_days", 19}}}, 1, 0,
					}}}}}},
				}}},
				bson.D{{Key: "$addFields", Value: bson.D{
					{Key: "total", Value: bson.D{{Key: "$add", Value: bson.A{"$huge", "$little"}}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "huge", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$total", 0}}},
						0,
						bson.D{{Key: "$multiply", Value: bson.A{bson.D{{Key: "$divide", Value: bson.A{"$huge", "$total"}}}, 100}}},
					}}}},
					{Key: "little", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$eq", Value: bson.A{"$total", 0}}},
						0,
						bson.D{{Key: "$multiply", Value: bson.A{bson.D{{Key: "$divide", Value: bson.A{"$little", "$total"}}}, 100}}},
					}}}},
				}}},
			}},
			{Key: "overall_growth_speed", Value: bson.A{
				// Группируем по дню и считаем сумму всех коков за каждый день
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: bson.D{
						{Key: "year", Value: bson.D{
							{Key: "$year", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
						{Key: "month", Value: bson.D{
							{Key: "$month", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
						{Key: "day", Value: bson.D{
							{Key: "$dayOfMonth", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
					}},
					{Key: "daily_total", Value: bson.D{{Key: "$sum", Value: "$size"}}},
					{Key: "day_timestamp", Value: bson.D{{Key: "$first", Value: "$requested_at"}}},
				}}},
				// Сортируем по дате (от новых к старым)
				bson.D{{Key: "$sort", Value: bson.D{{Key: "day_timestamp", Value: -1}}}},
				// Берем последние 5 дней
				bson.D{{Key: "$limit", Value: 5}},
				// Группируем все дни вместе
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "daily_totals", Value: bson.D{{Key: "$push", Value: "$daily_total"}}},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				}}},
				// Вычисляем скорость роста общего кока
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "growth_speed", Value: bson.D{{Key: "$round", Value: bson.A{
						bson.D{{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$gt", Value: bson.A{"$count", 0}}},
							// Скорость = средний дневной прирост за период
							bson.D{{Key: "$divide", Value: bson.A{
								bson.D{{Key: "$sum", Value: "$daily_totals"}},
								"$count",
							}}},
							0,
						}}},
						1,
					}}}},
				}}},
			}},
			{Key: "record", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: bson.D{
						{Key: "year", Value: bson.D{
							{Key: "$year", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
						{Key: "month", Value: bson.D{
							{Key: "$month", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
						{Key: "day", Value: bson.D{
							{Key: "$dayOfMonth", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
					}},
					{Key: "requested_at", Value: bson.D{{Key: "$first", Value: "$requested_at"}}},
					{Key: "total", Value: bson.D{{Key: "$sum", Value: "$size"}}},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "total", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 1}},
			}},
			{Key: "individual_record", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$sort", Value: bson.D{
					{Key: "size", Value: -1},
					{Key: "requested_at", Value: -1},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: 0},
					{Key: "requested_at", Value: 1},
					{Key: "total", Value: "$size"},
				}}},
			}},
			{Key: "total_cocks_count", Value: bson.A{
				bson.D{{Key: "$count", Value: "total_count"}},
			}},
			{Key: "individual_cocks_count", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$count", Value: "user_count"}},
			}},
			{Key: "individual_luck", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 5}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "avg_size", Value: bson.D{{Key: "$avg", Value: "$size"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "luck_coefficient", Value: bson.D{{Key: "$round", Value: bson.A{
						bson.D{{Key: "$divide", Value: bson.A{"$avg_size", 30.5}}},
						3,
					}}}},
				}}},
			}},
			{Key: "individual_volatility", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 5}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "sizes", Value: bson.D{{Key: "$push", Value: "$size"}}},
					{Key: "avg_size", Value: bson.D{{Key: "$avg", Value: "$size"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "volatility", Value: bson.D{{Key: "$round", Value: bson.A{
						bson.D{{Key: "$sqrt", Value: bson.D{{Key: "$avg", Value: bson.D{{Key: "$map", Value: bson.D{
							{Key: "input", Value: "$sizes"},
							{Key: "as", Value: "size"},
							{Key: "in", Value: bson.D{{Key: "$pow", Value: bson.A{
								bson.D{{Key: "$subtract", Value: bson.A{"$$size", "$avg_size"}}},
								2,
							}}}},
						}}}}}}},
						1,
					}}}},
				}}},
			}},
			{Key: "individual_five_cocks_dynamics", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 5}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "first_cock", Value: bson.D{{Key: "$first", Value: "$size"}}},
					{Key: "last_cock", Value: bson.D{{Key: "$last", Value: "$size"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "five_cocks_change", Value: bson.D{{Key: "$subtract", Value: bson.A{"$first_cock", "$last_cock"}}}},
					{Key: "five_cocks_change_percent", Value: bson.D{{Key: "$round", Value: bson.A{
						bson.D{{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$eq", Value: bson.A{"$last_cock", 0}}},
							0,
							bson.D{{Key: "$multiply", Value: bson.A{
								bson.D{{Key: "$divide", Value: bson.A{
									bson.D{{Key: "$subtract", Value: bson.A{"$first_cock", "$last_cock"}}},
									"$last_cock",
								}}},
								100,
							}}},
						}}},
						1,
					}}}},
				}}},
			}},
			{Key: "individual_growth_speed", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				// Группируем коки пользователя по дням
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: bson.D{
						{Key: "year", Value: bson.D{
							{Key: "$year", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
						{Key: "month", Value: bson.D{
							{Key: "$month", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
						{Key: "day", Value: bson.D{
							{Key: "$dayOfMonth", Value: bson.D{
								{Key: "date", Value: "$requested_at"},
								{Key: "timezone", Value: "Europe/Moscow"},
							}},
						}},
					}},
					{Key: "last_cock_of_day", Value: bson.D{{Key: "$last", Value: "$size"}}},
					{Key: "day_timestamp", Value: bson.D{{Key: "$last", Value: "$requested_at"}}},
				}}},
				// Сортируем по дате
				bson.D{{Key: "$sort", Value: bson.D{{Key: "day_timestamp", Value: -1}}}},
				// Берем последние 5 дней
				bson.D{{Key: "$limit", Value: 5}},
				// Группируем все дни вместе
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "daily_sizes", Value: bson.D{{Key: "$push", Value: "$last_cock_of_day"}}},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				}}},
				// Вычисляем скорость роста
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "growth_speed", Value: bson.D{{Key: "$round", Value: bson.A{
						bson.D{{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$gt", Value: bson.A{"$count", 0}}},
							// Скорость = средний дневной прирост за период
							bson.D{{Key: "$divide", Value: bson.A{
								bson.D{{Key: "$sum", Value: "$daily_sizes"}},
								"$count",
							}}},
							0,
						}}},
						1,
					}}}},
				}}},
			}},
			{Key: "individual_first_cock_date", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: 1}}}},
				bson.D{{Key: "$limit", Value: 1}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "first_date", Value: "$requested_at"},
				}}},
			}},
		}}},
	}
}

func PipelineTopUsersBySize() mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "nickname", Value: bson.D{{Key: "$first", Value: "$nickname"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "total_size", Value: -1}}}},
		{{Key: "$limit", Value: 13}},
	}
}

func PipelineTotalCockersCount() mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
		}}},
		{{Key: "$count", Value: "total"}},
	}
}

func PipelineUserPositionInLadder(userID int64) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "total_size", Value: -1}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "users", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "user_id", Value: "$_id"},
				{Key: "total_size", Value: "$total_size"},
			}}}},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$users"},
			{Key: "includeArrayIndex", Value: "position"},
		}}},
		{{Key: "$match", Value: bson.D{{Key: "users.user_id", Value: userID}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "position", Value: bson.D{{Key: "$add", Value: bson.A{"$position", 1}}}},
		}}},
	}
}

// PipelineUsersAroundPositionInLadder возвращает 3 пользователей около указанной позиции
func PipelineUsersAroundPositionInLadder(position int) mongo.Pipeline {
	skip := position - 2 // Берем 1 до, текущий, 1 после
	if skip < 0 {
		skip = 0
	}

	return mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "nickname", Value: bson.D{{Key: "$first", Value: "$nickname"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "total_size", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: 3}},
	}
}

func PipelineUserTotalSize(userID int64) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "nickname", Value: bson.D{{Key: "$first", Value: "$nickname"}}},
		}}},
	}
}

func PipelineFirstCockDate() mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: 1}}}},
		{{Key: "$limit", Value: 1}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "first_date", Value: "$requested_at"},
		}}},
	}
}

func PipelineSeasonWinners(startDate, endDate time.Time) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "requested_at", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lt", Value: endDate},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "nickname", Value: bson.D{{Key: "$first", Value: "$nickname"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "total_size", Value: -1}}}},
		{{Key: "$limit", Value: 3}},
	}
}

func PipelineTopUsersInSeason(startDate, endDate time.Time) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "requested_at", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lt", Value: endDate},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "nickname", Value: bson.D{{Key: "$first", Value: "$nickname"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "total_size", Value: -1}}}},
		{{Key: "$limit", Value: 13}},
	}
}

func PipelineSeasonCockersCount(startDate, endDate time.Time) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "requested_at", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lt", Value: endDate},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
		}}},
		{{Key: "$count", Value: "total"}},
	}
}

func PipelineUserPositionInSeason(userID int64, startDate, endDate time.Time) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "requested_at", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lt", Value: endDate},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "total_size", Value: -1}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "users", Value: bson.D{{Key: "$push", Value: bson.D{
				{Key: "user_id", Value: "$_id"},
				{Key: "total_size", Value: "$total_size"},
			}}}},
		}}},
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$users"},
			{Key: "includeArrayIndex", Value: "position"},
		}}},
		{{Key: "$match", Value: bson.D{{Key: "users.user_id", Value: userID}}}},
		{{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "position", Value: bson.D{{Key: "$add", Value: bson.A{"$position", 1}}}},
		}}},
	}
}

// PipelineUsersAroundPositionInSeason возвращает 3 пользователей около указанной позиции в сезоне
func PipelineUsersAroundPositionInSeason(position int, startDate, endDate time.Time) mongo.Pipeline {
	skip := position - 2 // Берем 1 до, текущий, 1 после
	if skip < 0 {
		skip = 0
	}

	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "requested_at", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lt", Value: endDate},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "nickname", Value: bson.D{{Key: "$first", Value: "$nickname"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "total_size", Value: -1}}}},
		{{Key: "$skip", Value: skip}},
		{{Key: "$limit", Value: 3}},
	}
}

func PipelineAllUsersInSeason(startDate, endDate time.Time) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "requested_at", Value: bson.D{
				{Key: "$gte", Value: startDate},
				{Key: "$lt", Value: endDate},
			}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "nickname", Value: bson.D{{Key: "$first", Value: "$nickname"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "total_size", Value: -1}}}},
	}
}

// PipelineCheckAchievements проверяет все условия достижений для пользователя
func PipelineCheckAchievements(userId int64) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
		{{Key: "$facet", Value: bson.D{
			// 1. Простой счетчик записей (количество дерганий)
			{Key: "total_pulls", Value: bson.A{
				bson.D{{Key: "$count", Value: "count"}},
			}},

			// 2. Сумма размеров (накопленный размер)
			{Key: "total_size", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "total", Value: bson.D{{Key: "$sum", Value: "$size"}}},
				}}},
			}},

			// 3. Конкретные значения
			{Key: "sniper_30cm", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "size", Value: 30}}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "half_hundred_50cm", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "size", Value: 50}}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "maximalist_61cm", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "size", Value: 61}}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "beautiful_numbers", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "size", Value: bson.D{{Key: "$in", Value: bson.A{11, 22, 33, 44, 55}}}}}}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$size"},
				}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},

			// 4. Последовательности (последние записи)
			{Key: "recent_10", Value: bson.A{
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 10}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: 1}}}},
			}},

			// 5. Экстремумы (максимум и минимум)
			{Key: "max_size", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "max", Value: bson.D{{Key: "$max", Value: "$size"}}},
				}}},
			}},
			{Key: "min_size", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "min", Value: bson.D{{Key: "$min", Value: "$size"}}},
				}}},
			}},

			// 6. Временные (ранняя пташка - до 6:00 МСК)
			{Key: "early_bird", Value: bson.A{
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "hour", Value: bson.D{{Key: "$hour", Value: bson.D{
						{Key: "date", Value: "$requested_at"},
						{Key: "timezone", Value: "Europe/Moscow"},
					}}}},
				}}},
				bson.D{{Key: "$match", Value: bson.D{{Key: "hour", Value: bson.D{{Key: "$lt", Value: 6}}}}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},

			// 7. Спидраннер (<30 сек после полуночи)
			{Key: "speedrunner", Value: bson.A{
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "hour", Value: bson.D{{Key: "$hour", Value: bson.D{
						{Key: "date", Value: "$requested_at"},
						{Key: "timezone", Value: "Europe/Moscow"},
					}}}},
					{Key: "minute", Value: bson.D{{Key: "$minute", Value: bson.D{
						{Key: "date", Value: "$requested_at"},
						{Key: "timezone", Value: "Europe/Moscow"},
					}}}},
					{Key: "second", Value: bson.D{{Key: "$second", Value: bson.D{
						{Key: "date", Value: "$requested_at"},
						{Key: "timezone", Value: "Europe/Moscow"},
					}}}},
				}}},
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "hour", Value: 0},
					{Key: "minute", Value: 0},
					{Key: "second", Value: bson.D{{Key: "$lt", Value: 30}}},
				}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},

			// 7.1. На волоске от смерти (после 23:00 МСК)
			{Key: "midnight_puller", Value: bson.A{
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "hour", Value: bson.D{{Key: "$hour", Value: bson.D{
						{Key: "date", Value: "$requested_at"},
						{Key: "timezone", Value: "Europe/Moscow"},
					}}}},
				}}},
				bson.D{{Key: "$match", Value: bson.D{{Key: "hour", Value: 23}}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},

			// 8. Праздничные
			{Key: "valentine", Value: bson.A{
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "size", Value: 1},
					{Key: "month", Value: bson.D{
						{Key: "$month", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
					{Key: "day", Value: bson.D{
						{Key: "$dayOfMonth", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
				}}},
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "month", Value: 2},
					{Key: "day", Value: 14},
					{Key: "size", Value: 14},
				}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "new_year_gift", Value: bson.A{
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "size", Value: 1},
					{Key: "month", Value: bson.D{
						{Key: "$month", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
					{Key: "day", Value: bson.D{
						{Key: "$dayOfMonth", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
				}}},
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "month", Value: 12},
					{Key: "day", Value: 31},
					{Key: "size", Value: bson.D{{Key: "$gte", Value: 60}}},
				}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "mens_solidarity", Value: bson.A{
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "size", Value: 1},
					{Key: "month", Value: bson.D{
						{Key: "$month", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
					{Key: "day", Value: bson.D{
						{Key: "$dayOfMonth", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
				}}},
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "month", Value: 11},
					{Key: "day", Value: 19},
					{Key: "size", Value: 19},
				}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "friday_13th", Value: bson.A{
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "size", Value: 1},
					{Key: "day", Value: bson.D{
						{Key: "$dayOfMonth", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
					{Key: "day_of_week", Value: bson.D{
						{Key: "$dayOfWeek", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
				}}},
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "day", Value: 13},
					{Key: "day_of_week", Value: 6}, // Пятница
					{Key: "size", Value: 0},
				}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "leap_cock", Value: bson.A{
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "month", Value: bson.D{
						{Key: "$month", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
					{Key: "day", Value: bson.D{
						{Key: "$dayOfMonth", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
				}}},
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "month", Value: 2},
					{Key: "day", Value: 29},
				}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},

			// 9. Динамика (молния - рост на 50+см)
			{Key: "lightning", Value: bson.A{
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: 1}}}},
				bson.D{{Key: "$setWindowFields", Value: bson.D{
					{Key: "sortBy", Value: bson.D{{Key: "requested_at", Value: 1}}},
					{Key: "output", Value: bson.D{
						{Key: "prev_size", Value: bson.D{
							{Key: "$shift", Value: bson.D{
								{Key: "output", Value: "$size"},
								{Key: "by", Value: -1},
							}},
						}},
					}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "growth", Value: bson.D{{Key: "$subtract", Value: bson.A{"$size", bson.D{{Key: "$ifNull", Value: bson.A{"$prev_size", "$size"}}}}}}},
				}}},
				bson.D{{Key: "$match", Value: bson.D{{Key: "growth", Value: bson.D{{Key: "$gte", Value: 50}}}}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},

			// 10. Последние 31 кок для сложных коллекций
			{Key: "last_31", Value: bson.A{
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 31}},
			}},

			// 11. Последние 3 кока для проверки специальных совпадений
			{Key: "recent_3", Value: bson.A{
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 3}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: 1}}}},
			}},
		}}},
	}
}

// PipelineGlobalMaxMin возвращает глобальные максимум и минимум
func PipelineGlobalMaxMin() mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "max", Value: bson.D{{Key: "$max", Value: "$size"}}},
			{Key: "min", Value: bson.D{{Key: "$min", Value: "$size"}}},
		}}},
	}
}

// PipelineCountSeasons подсчитывает количество уникальных сезонов для пользователя
// Сезон = 3-месячный период с момента первого кока в системе
// Логика: вычисляем номер сезона как разницу в месяцах / 3 (округленную вниз)
func PipelineCountSeasons(userId int64) mongo.Pipeline {
	return mongo.Pipeline{
		// Получаем дату первого кока в системе (в московском часовом поясе)
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "first_cock_date", Value: bson.D{{Key: "$min", Value: "$requested_at"}}},
		}}},
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "cocks"},
			{Key: "let", Value: bson.D{{Key: "first_date", Value: "$first_cock_date"}}},
			{Key: "pipeline", Value: bson.A{
				// Берем только коки пользователя
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				// Вычисляем номер сезона для каждого кока
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "requested_at", Value: 1},
					// Извлекаем год и месяц из requested_at (в московском часовом поясе)
					{Key: "year", Value: bson.D{
						{Key: "$year", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
					{Key: "month", Value: bson.D{
						{Key: "$month", Value: bson.D{
							{Key: "date", Value: "$requested_at"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
					// Извлекаем год и месяц из first_date
					{Key: "first_year", Value: bson.D{
						{Key: "$year", Value: bson.D{
							{Key: "date", Value: "$$first_date"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
					{Key: "first_month", Value: bson.D{
						{Key: "$month", Value: bson.D{
							{Key: "date", Value: "$$first_date"},
							{Key: "timezone", Value: "Europe/Moscow"},
						}},
					}},
				}}},
				// Вычисляем разницу в месяцах и номер сезона
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "months_diff", Value: bson.D{
						{Key: "$subtract", Value: bson.A{
							// Месяцы от начала эпохи для requested_at
							bson.D{{Key: "$add", Value: bson.A{
								bson.D{{Key: "$multiply", Value: bson.A{"$year", 12}}},
								"$month",
							}}},
							// Месяцы от начала эпохи для first_date
							bson.D{{Key: "$add", Value: bson.A{
								bson.D{{Key: "$multiply", Value: bson.A{"$first_year", 12}}},
								"$first_month",
							}}},
						}},
					}},
				}}},
				// Номер сезона = разница в месяцах / 3 (округленная вниз)
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "season_number", Value: bson.D{
						{Key: "$floor", Value: bson.D{
							{Key: "$divide", Value: bson.A{"$months_diff", 3}},
						}},
					}},
				}}},
				// Группируем по номеру сезона
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$season_number"},
				}}},
			}},
			{Key: "as", Value: "seasons"},
		}}},
		{{Key: "$unwind", Value: "$seasons"}},
		{{Key: "$count", Value: "count"}},
	}
}

// PipelineCheckTraveler проверяет, получил ли пользователь все размеры (регионы)
func PipelineCheckTraveler(userId int64) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$size"},
		}}},
		{{Key: "$count", Value: "unique_sizes"}},
	}
}

// PipelineCheckMuscovite проверяет, получил ли пользователь размер 50см 5 раз за последние 31 день
func PipelineCheckMuscovite(userId int64, startDate time.Time) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{
			{Key: "user_id", Value: userId},
			{Key: "size", Value: 50},
			{Key: "requested_at", Value: bson.D{{Key: "$gte", Value: startDate}}},
		}}},
		{{Key: "$count", Value: "count"}},
	}
}
