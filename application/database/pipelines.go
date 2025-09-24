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
			{Key: "individual_cock", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "total", Value: bson.D{{Key: "$sum", Value: "$size"}}},
					{Key: "avg_val", Value: bson.D{{Key: "$avg", Value: "$size"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "total", Value: 1},
					{Key: "average", Value: bson.D{{Key: "$round", Value: bson.A{"$avg_val", 0}}}},
				}}},
			}},
			{Key: "overall", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
					{Key: "median", Value: bson.D{{Key: "$median", Value: bson.D{
						{Key: "input", Value: "$size"},
						{Key: "method", Value: "approximate"},
					}}}},
					{Key: "average", Value: bson.D{{Key: "$avg", Value: "$size"}}},
				}}},
				bson.D{{Key: "$project", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "size", Value: 1},
					{Key: "median", Value: 1},
					{Key: "average", Value: bson.D{{Key: "$round", Value: bson.A{"$average", 0}}}},
				}}},
				bson.D{{Key: "$limit", Value: 1}},
			}},
			{Key: "uniques", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$user_id"}}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "distribution", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "huge", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$gte", Value: bson.A{"$size", 19}}}, 1, 0,
					}}}}}},
					{Key: "little", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$lt", Value: bson.A{"$size", 19}}}, 1, 0,
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
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$requested_at"},
					{Key: "requested_at", Value: bson.D{{Key: "$first", Value: "$requested_at"}}},
					{Key: "total", Value: bson.D{{Key: "$first", Value: "$size"}}},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "total", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 1}},
			}},
			{Key: "total_cocks_count", Value: bson.A{
				bson.D{{Key: "$count", Value: "total_count"}},
			}},
			{Key: "individual_cocks_count", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: userId}}}},
				bson.D{{Key: "$count", Value: "user_count"}},
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
