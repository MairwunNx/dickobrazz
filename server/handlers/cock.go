package handlers

import (
	"dickobrazz/server/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CockSize(c *gin.Context) {
	respondOK(c, models.CockSizeResponse{
		Size: 0,
		Hash: "",
		Salt: "",
	})
}

func (h *Handler) CockRuler(c *gin.Context) {
	respondOK(c, models.CockRulerResponse{
		Leaders:           []models.LeaderboardEntry{},
		TotalParticipants: 0,
		UserPosition:      0,
		Neighborhood: models.UserNeighborhood{
			Above: []models.LeaderboardEntry{},
			Self:  nil,
			Below: []models.LeaderboardEntry{},
		},
		Page: paginationFromQuery(c),
	})
}

func (h *Handler) CockRace(c *gin.Context) {
	respondOK(c, models.CockRaceResponse{
		Season:            nil,
		Leaders:           []models.RaceEntry{},
		TotalParticipants: 0,
		UserPosition:      0,
		Neighborhood: models.RaceNeighborhood{
			Above: []models.RaceEntry{},
			Self:  nil,
			Below: []models.RaceEntry{},
		},
		Page: paginationFromQuery(c),
	})
}

func (h *Handler) CockDynamic(c *gin.Context) {
	respondOK(c, models.CockDynamicResponse{
		Overall: models.CockDynamicOverall{
			TotalSize:       0,
			UniqueUsers:     0,
			Recent:          models.CockDynamicRecentStat{},
			Distribution:    models.CockDynamicPercentile{},
			Record:          models.CockDynamicRecord{},
			TotalCocksCount: 0,
			GrowthSpeed:     0,
		},
		Individual: models.CockDynamicIndividual{
			TotalSize:          0,
			RecentAverage:      0,
			Irk:                0,
			Record:             models.CockDynamicRecord{},
			Dominance:          0,
			DailyGrowthAverage: 0,
			DailyDynamics:      models.CockDynamicDailyDynamics{},
			FiveCocksDynamics:  models.CockDynamicFiveCocksDynamics{},
			GrowthSpeed:        0,
			FirstCockDate:      "",
			LuckCoefficient:    0,
			Volatility:         0,
			CocksCount:         0,
		},
	})
}

func (h *Handler) CockSeasons(c *gin.Context) {
	respondOK(c, models.CockSeasonsResponse{
		Seasons: []models.SeasonWithWinners{},
		Page:    paginationFromQuery(c),
	})
}

func (h *Handler) CockAchievements(c *gin.Context) {
	respondOK(c, models.CockAchievementsResponse{
		Achievements: []models.Achievement{},
		Page:         paginationFromQuery(c),
	})
}

func (h *Handler) CockLadder(c *gin.Context) {
	respondOK(c, models.CockLadderResponse{
		Leaders:           []models.LeaderboardEntry{},
		TotalParticipants: 0,
		UserPosition:      0,
		Neighborhood: models.UserNeighborhood{
			Above: []models.LeaderboardEntry{},
			Self:  nil,
			Below: []models.LeaderboardEntry{},
		},
		Page: paginationFromQuery(c),
	})
}
