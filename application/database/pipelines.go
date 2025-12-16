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
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$user_id"},
					{Key: "cocks", Value: bson.D{{Key: "$push", Value: bson.D{
						{Key: "size", Value: "$size"},
						{Key: "requested_at", Value: "$requested_at"},
					}}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "last_cocks", Value: bson.D{{Key: "$slice", Value: bson.A{
						bson.D{{Key: "$sortArray", Value: bson.D{
							{Key: "input", Value: "$cocks"},
							{Key: "sortBy", Value: bson.D{{Key: "requested_at", Value: -1}}},
						}}},
						5,
					}}}},
				}}},
				bson.D{{Key: "$unwind", Value: "$last_cocks"}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "median", Value: bson.D{{Key: "$median", Value: bson.D{
						{Key: "input", Value: "$last_cocks.size"},
						{Key: "method", Value: "approximate"},
					}}}},
					{Key: "average", Value: bson.D{{Key: "$avg", Value: "$last_cocks.size"}}},
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
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$user_id"},
					{Key: "cocks", Value: bson.D{{Key: "$push", Value: bson.D{
						{Key: "size", Value: "$size"},
						{Key: "requested_at", Value: "$requested_at"},
					}}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "last_cocks", Value: bson.D{{Key: "$slice", Value: bson.A{
						bson.D{{Key: "$sortArray", Value: bson.D{
							{Key: "input", Value: "$cocks"},
							{Key: "sortBy", Value: bson.D{{Key: "requested_at", Value: -1}}},
						}}},
						5,
					}}}},
				}}},
				bson.D{{Key: "$unwind", Value: "$last_cocks"}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "huge", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$gte", Value: bson.A{"$last_cocks.size", 19}}}, 1, 0,
					}}}}}},
					{Key: "little", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$lt", Value: bson.A{"$last_cocks.size", 19}}}, 1, 0,
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
			{Key: "record", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: bson.D{
						{Key: "year", Value: bson.D{{Key: "$year", Value: "$requested_at"}}},
						{Key: "month", Value: bson.D{{Key: "$month", Value: "$requested_at"}}},
						{Key: "day", Value: bson.D{{Key: "$dayOfMonth", Value: "$requested_at"}}},
					}},
					{Key: "requested_at", Value: bson.D{{Key: "$first", Value: "$requested_at"}}},
					{Key: "total", Value: bson.D{{Key: "$sum", Value: "$size"}}},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "total", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 1}},
			}},
			{Key: "individual_record", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "size", Value: -1}}}},
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
				bson.D{{Key: "$match", Value: bson.D{{Key: "prev_size", Value: bson.D{{Key: "$ne", Value: nil}}}}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "growth", Value: bson.D{{Key: "$abs", Value: bson.D{{Key: "$subtract", Value: bson.A{"$size", "$prev_size"}}}}}},
				}}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "avg_growth_speed", Value: bson.D{{Key: "$avg", Value: "$growth"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "growth_speed", Value: bson.D{{Key: "$round", Value: bson.A{"$avg_growth_speed", 1}}}},
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
			
			// 8. Праздничные
			{Key: "valentine", Value: bson.A{
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "size", Value: 1},
					{Key: "month", Value: bson.D{{Key: "$month", Value: "$requested_at"}}},
					{Key: "day", Value: bson.D{{Key: "$dayOfMonth", Value: "$requested_at"}}},
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
					{Key: "month", Value: bson.D{{Key: "$month", Value: "$requested_at"}}},
					{Key: "day", Value: bson.D{{Key: "$dayOfMonth", Value: "$requested_at"}}},
				}}},
				bson.D{{Key: "$match", Value: bson.D{
					{Key: "month", Value: 12},
					{Key: "day", Value: 31},
					{Key: "size", Value: bson.D{{Key: "$gte", Value: 60}}},
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
func PipelineCountSeasons(userId int64) mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$season_name"},
		}}},
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
