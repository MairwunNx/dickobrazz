package application

import "dickobrazz/application/database"

// AllAchievements содержит все доступные достижения в игре
// Отсортированы по возрастанию респектов (33-2222)
var AllAchievements = []database.Achievement{
	// 33 респекта
	{
		ID:          "not_rubbed_yet",
		Emoji:       "🤏",
		Name:        "AchievementNameNotRubbedYet",
		Description: "AchievementDescNotRubbedYet",
		Respects:    33,
		MaxProgress: 10,
	},

	// 50 респектов
	{
		ID:          "half_hundred",
		Emoji:       "🌟",
		Name:        "AchievementNameHalfHundred",
		Description: "AchievementDescHalfHundred",
		Respects:    50,
		MaxProgress: 50,
	},
	{
		ID:          "diary",
		Emoji:       "📆",
		Name:        "AchievementNameDiary",
		Description: "AchievementDescDiary",
		Respects:    50,
		MaxProgress: 31,
	},

	// 90 респектов
	{
		ID:          "golden_hundred",
		Emoji:       "💯",
		Name:        "AchievementNameGoldenHundred",
		Description: "AchievementDescGoldenHundred",
		Respects:    90,
		MaxProgress: 100,
	},

	// 100 респектов
	{
		ID:          "skillful_hands",
		Emoji:       "💪",
		Name:        "AchievementNameSkillfulHands",
		Description: "AchievementDescSkillfulHands",
		Respects:    100,
		MaxProgress: 100,
	},

	// 135 респектов
	{
		ID:          "early_bird",
		Emoji:       "🌅",
		Name:        "AchievementNameEarlyBird",
		Description: "AchievementDescEarlyBird",
		Respects:    135,
		MaxProgress: 20,
	},

	// 200 респектов
	{
		ID:          "lightning",
		Emoji:       "⚡",
		Name:        "AchievementNameLightning",
		Description: "AchievementDescLightning",
		Respects:    200,
		MaxProgress: 1,
	},

	// 211 респектов
	{
		ID:          "sniper",
		Emoji:       "🎯",
		Name:        "AchievementNameSniper",
		Description: "AchievementDescSniper",
		Respects:    211,
		MaxProgress: 5,
	},

	// 222 респекта
	{
		ID:          "deja_vu",
		Emoji:       "🔄",
		Name:        "AchievementNameDejaVu",
		Description: "AchievementDescDejaVu",
		Respects:    222,
		MaxProgress: 2,
	},
	{
		ID:          "speedrunner",
		Emoji:       "⏱️",
		Name:        "AchievementNameSpeedrunner",
		Description: "AchievementDescSpeedrunner",
		Respects:    222,
		MaxProgress: 5,
	},
	{
		ID:          "midnight_puller",
		Emoji:       "☠️",
		Name:        "AchievementNameMidnightPuller",
		Description: "AchievementDescMidnightPuller",
		Respects:    222,
		MaxProgress: 10,
	},

	// 228 респектов
	{
		ID:          "rounder",
		Emoji:       "🔟",
		Name:        "AchievementNameRounder",
		Description: "AchievementDescRounder",
		Respects:    228,
		MaxProgress: 6,
	},
	{
		ID:          "everest",
		Emoji:       "🏔️",
		Name:        "AchievementNameEverest",
		Description: "AchievementDescEverest",
		Respects:    228,
		MaxProgress: 61,
	},
	{
		ID:          "mariana_trench",
		Emoji:       "🕳️",
		Name:        "AchievementNameMarianaTrench",
		Description: "AchievementDescMarianaTrench",
		Respects:    228,
		MaxProgress: 0,
	},

	// 233 респекта
	{
		ID:          "number_collector",
		Emoji:       "🔢",
		Name:        "AchievementNameNumberCollector",
		Description: "AchievementDescNumberCollector",
		Respects:    233,
		MaxProgress: 5,
	},

	// 300 респектов
	{
		ID:          "day_equals_size",
		Emoji:       "📅",
		Name:        "AchievementNameDayEqualsSize",
		Description: "AchievementDescDayEqualsSize",
		Respects:    300,
		MaxProgress: 0,
	},

	// 333 респекта
	{
		ID:          "solid_thousand",
		Emoji:       "💰",
		Name:        "AchievementNameSolidThousand",
		Description: "AchievementDescSolidThousand",
		Respects:    333,
		MaxProgress: 1000,
	},
	{
		ID:          "bull_trend",
		Emoji:       "📈",
		Name:        "AchievementNameBullTrend",
		Description: "AchievementDescBullTrend",
		Respects:    333,
		MaxProgress: 5,
	},
	{
		ID:          "bear_market",
		Emoji:       "📉",
		Name:        "AchievementNameBearMarket",
		Description: "AchievementDescBearMarket",
		Respects:    333,
		MaxProgress: 5,
	},

	// 500 респектов
	{
		ID:          "traveler",
		Emoji:       "🗺️",
		Name:        "AchievementNameTraveler",
		Description: "AchievementDescTraveler",
		Respects:    500,
		MaxProgress: 61,
	},

	// 555 респектов
	{
		ID:          "freeze",
		Emoji:       "❄️",
		Name:        "AchievementNameFreeze",
		Description: "AchievementDescFreeze",
		Respects:    555,
		MaxProgress: 5,
	},

	// 700 респектов
	{
		ID:          "five_k",
		Emoji:       "💎",
		Name:        "AchievementNameFiveK",
		Description: "AchievementDescFiveK",
		Respects:    700,
		MaxProgress: 5000,
	},
	{
		ID:          "oldtimer",
		Emoji:       "🗓️",
		Name:        "AchievementNameOldtimer",
		Description: "AchievementDescOldtimer",
		Respects:    700,
		MaxProgress: 3,
	},
	{
		ID:          "anniversary",
		Emoji:       "🎂",
		Name:        "AchievementNameAnniversary",
		Description: "AchievementDescAnniversary",
		Respects:    700,
		MaxProgress: 365,
	},

	// 777 респектов
	{
		ID:          "contrast_shower",
		Emoji:       "🚿",
		Name:        "AchievementNameContrastShower",
		Description: "AchievementDescContrastShower",
		Respects:    777,
		MaxProgress: 0,
	},
	{
		ID:          "pythagoras",
		Emoji:       "📐",
		Name:        "AchievementNamePythagoras",
		Description: "AchievementDescPythagoras",
		Respects:    777,
		MaxProgress: 1,
	},
	{
		ID:          "leet_speak",
		Emoji:       "💻",
		Name:        "AchievementNameLeetSpeak",
		Description: "AchievementDescLeetSpeak",
		Respects:    777,
		MaxProgress: 1,
	},

	// 800 респектов
	{
		ID:          "moscovite",
		Emoji:       "🏙️",
		Name:        "AchievementNameMoscovite",
		Description: "AchievementDescMoscovite",
		Respects:    800,
		MaxProgress: 5,
	},

	// 888 респектов
	{
		ID:          "hour_precision",
		Emoji:       "🕐",
		Name:        "AchievementNameHourPrecision",
		Description: "AchievementDescHourPrecision",
		Respects:    888,
		MaxProgress: 0,
	},

	// 900 респектов
	{
		ID:          "wonder_stranger",
		Emoji:       "💋",
		Name:        "AchievementNameWonderStranger",
		Description: "AchievementDescWonderStranger",
		Respects:    900,
		MaxProgress: 500,
	},

	// 999 респектов
	{
		ID:          "valentine",
		Emoji:       "💝",
		Name:        "AchievementNameValentine",
		Description: "AchievementDescValentine",
		Respects:    999,
		MaxProgress: 1,
	},
	{
		ID:          "new_year_gift",
		Emoji:       "🎄",
		Name:        "AchievementNameNewYearGift",
		Description: "AchievementDescNewYearGift",
		Respects:    999,
		MaxProgress: 1,
	},
	{
		ID:          "mens_solidarity",
		Emoji:       "🤝",
		Name:        "AchievementNameMensSolidarity",
		Description: "AchievementDescMensSolidarity",
		Respects:    999,
		MaxProgress: 1,
	},
	{
		ID:          "friday_13th",
		Emoji:       "☠️",
		Name:        "AchievementNameFriday13th",
		Description: "AchievementDescFriday13th",
		Respects:    999,
		MaxProgress: 1,
	},
	{
		ID:          "leap_cock",
		Emoji:       "📅",
		Name:        "AchievementNameLeapCock",
		Description: "AchievementDescLeapCock",
		Respects:    999,
		MaxProgress: 1,
	},

	// 1000 респектов
	{
		ID:          "turtle",
		Emoji:       "🐌",
		Name:        "AchievementNameTurtle",
		Description: "AchievementDescTurtle",
		Respects:    1000,
		MaxProgress: 10,
	},
	{
		ID:          "golden_cock",
		Emoji:       "👑",
		Name:        "AchievementNameGoldenCock",
		Description: "AchievementDescGoldenCock",
		Respects:    1000,
		MaxProgress: 10000,
	},
	{
		ID:          "sum_of_previous",
		Emoji:       "🎲",
		Name:        "AchievementNameSumOfPrevious",
		Description: "AchievementDescSumOfPrevious",
		Respects:    1000,
		MaxProgress: 0,
	},

	// 1222 респекта
	{
		ID:          "bazooka_hands",
		Emoji:       "💥",
		Name:        "AchievementNameBazookaHands",
		Description: "AchievementDescBazookaHands",
		Respects:    1222,
		MaxProgress: 1000,
	},

	// 1333 респекта
	{
		ID:          "triple",
		Emoji:       "🎰",
		Name:        "AchievementNameTriple",
		Description: "AchievementDescTriple",
		Respects:    1333,
		MaxProgress: 3,
	},
	{
		ID:          "veteran",
		Emoji:       "🗓️",
		Name:        "AchievementNameVeteran",
		Description: "AchievementDescVeteran",
		Respects:    1333,
		MaxProgress: 5,
	},
	{
		ID:          "minute_precision",
		Emoji:       "⏰",
		Name:        "AchievementNameMinutePrecision",
		Description: "AchievementDescMinutePrecision",
		Respects:    1333,
		MaxProgress: 0,
	},

	// 1777 респектов
	{
		ID:          "poker",
		Emoji:       "🎴",
		Name:        "AchievementNamePoker",
		Description: "AchievementDescPoker",
		Respects:    1777,
		MaxProgress: 4,
	},

	// 1888 респектов
	{
		ID:          "keeper",
		Emoji:       "🗓️",
		Name:        "AchievementNameKeeper",
		Description: "AchievementDescKeeper",
		Respects:    1888,
		MaxProgress: 10,
	},

	// 2000 респектов
	{
		ID:          "cosmic_cock",
		Emoji:       "🚀",
		Name:        "AchievementNameCosmicCock",
		Description: "AchievementDescCosmicCock",
		Respects:    2000,
		MaxProgress: 20000,
	},
	{
		ID:          "maximalist",
		Emoji:       "🔝",
		Name:        "AchievementNameMaximalist",
		Description: "AchievementDescMaximalist",
		Respects:    2000,
		MaxProgress: 10,
	},

	// 2222 респекта
	{
		ID:          "diamond_hands",
		Emoji:       "💎",
		Name:        "AchievementNameDiamondHands",
		Description: "AchievementDescDiamondHands",
		Respects:    2222,
		MaxProgress: 7,
	},
	{
		ID:          "diamond_eye",
		Emoji:       "💎",
		Name:        "AchievementNameDiamondEye",
		Description: "AchievementDescDiamondEye",
		Respects:    2222,
		MaxProgress: 5,
	},
	{
		ID:          "greek_myth",
		Emoji:       "⚡",
		Name:        "AchievementNameGreekMyth",
		Description: "AchievementDescGreekMyth",
		Respects:    2222,
		MaxProgress: 30000,
	},
	{
		ID:          "fibonacci_father",
		Emoji:       "🔢",
		Name:        "AchievementNameFibonacciFather",
		Description: "AchievementDescFibonacciFather",
		Respects:    2222,
		MaxProgress: 9,
	},
	{
		ID:          "annihilator_cannon",
		Emoji:       "☢️",
		Name:        "AchievementNameAnnihilatorCannon",
		Description: "AchievementDescAnnihilatorCannon",
		Respects:    2222,
		MaxProgress: 5000,
	},
}

// GetAchievementByID возвращает достижение по его ID
func GetAchievementByID(id string) *database.Achievement {
	for i := range AllAchievements {
		if AllAchievements[i].ID == id {
			return &AllAchievements[i]
		}
	}
	return nil
}
