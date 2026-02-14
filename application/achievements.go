package application

// AchievementDef содержит клиентское определение достижения (ID + ключи локализации)
type AchievementDef struct {
	ID          string
	Name        string // ключ локализации
	Description string // ключ локализации
}

// AllAchievements содержит все доступные достижения в игре (только клиентские данные)
// Бэкэнд возвращает emoji, respects, completed, progress, max_progress — маппинг по ID
var AllAchievements = []AchievementDef{
	{ID: "not_rubbed_yet", Name: "AchievementNameNotRubbedYet", Description: "AchievementDescNotRubbedYet"},
	{ID: "half_hundred", Name: "AchievementNameHalfHundred", Description: "AchievementDescHalfHundred"},
	{ID: "diary", Name: "AchievementNameDiary", Description: "AchievementDescDiary"},
	{ID: "golden_hundred", Name: "AchievementNameGoldenHundred", Description: "AchievementDescGoldenHundred"},
	{ID: "skillful_hands", Name: "AchievementNameSkillfulHands", Description: "AchievementDescSkillfulHands"},
	{ID: "early_bird", Name: "AchievementNameEarlyBird", Description: "AchievementDescEarlyBird"},
	{ID: "lightning", Name: "AchievementNameLightning", Description: "AchievementDescLightning"},
	{ID: "sniper", Name: "AchievementNameSniper", Description: "AchievementDescSniper"},
	{ID: "deja_vu", Name: "AchievementNameDejaVu", Description: "AchievementDescDejaVu"},
	{ID: "speedrunner", Name: "AchievementNameSpeedrunner", Description: "AchievementDescSpeedrunner"},
	{ID: "midnight_puller", Name: "AchievementNameMidnightPuller", Description: "AchievementDescMidnightPuller"},
	{ID: "rounder", Name: "AchievementNameRounder", Description: "AchievementDescRounder"},
	{ID: "everest", Name: "AchievementNameEverest", Description: "AchievementDescEverest"},
	{ID: "mariana_trench", Name: "AchievementNameMarianaTrench", Description: "AchievementDescMarianaTrench"},
	{ID: "number_collector", Name: "AchievementNameNumberCollector", Description: "AchievementDescNumberCollector"},
	{ID: "day_equals_size", Name: "AchievementNameDayEqualsSize", Description: "AchievementDescDayEqualsSize"},
	{ID: "solid_thousand", Name: "AchievementNameSolidThousand", Description: "AchievementDescSolidThousand"},
	{ID: "bull_trend", Name: "AchievementNameBullTrend", Description: "AchievementDescBullTrend"},
	{ID: "bear_market", Name: "AchievementNameBearMarket", Description: "AchievementDescBearMarket"},
	{ID: "traveler", Name: "AchievementNameTraveler", Description: "AchievementDescTraveler"},
	{ID: "freeze", Name: "AchievementNameFreeze", Description: "AchievementDescFreeze"},
	{ID: "five_k", Name: "AchievementNameFiveK", Description: "AchievementDescFiveK"},
	{ID: "oldtimer", Name: "AchievementNameOldtimer", Description: "AchievementDescOldtimer"},
	{ID: "anniversary", Name: "AchievementNameAnniversary", Description: "AchievementDescAnniversary"},
	{ID: "contrast_shower", Name: "AchievementNameContrastShower", Description: "AchievementDescContrastShower"},
	{ID: "pythagoras", Name: "AchievementNamePythagoras", Description: "AchievementDescPythagoras"},
	{ID: "leet_speak", Name: "AchievementNameLeetSpeak", Description: "AchievementDescLeetSpeak"},
	{ID: "moscovite", Name: "AchievementNameMoscovite", Description: "AchievementDescMoscovite"},
	{ID: "hour_precision", Name: "AchievementNameHourPrecision", Description: "AchievementDescHourPrecision"},
	{ID: "wonder_stranger", Name: "AchievementNameWonderStranger", Description: "AchievementDescWonderStranger"},
	{ID: "valentine", Name: "AchievementNameValentine", Description: "AchievementDescValentine"},
	{ID: "new_year_gift", Name: "AchievementNameNewYearGift", Description: "AchievementDescNewYearGift"},
	{ID: "mens_solidarity", Name: "AchievementNameMensSolidarity", Description: "AchievementDescMensSolidarity"},
	{ID: "friday_13th", Name: "AchievementNameFriday13th", Description: "AchievementDescFriday13th"},
	{ID: "leap_cock", Name: "AchievementNameLeapCock", Description: "AchievementDescLeapCock"},
	{ID: "turtle", Name: "AchievementNameTurtle", Description: "AchievementDescTurtle"},
	{ID: "golden_cock", Name: "AchievementNameGoldenCock", Description: "AchievementDescGoldenCock"},
	{ID: "sum_of_previous", Name: "AchievementNameSumOfPrevious", Description: "AchievementDescSumOfPrevious"},
	{ID: "bazooka_hands", Name: "AchievementNameBazookaHands", Description: "AchievementDescBazookaHands"},
	{ID: "triple", Name: "AchievementNameTriple", Description: "AchievementDescTriple"},
	{ID: "veteran", Name: "AchievementNameVeteran", Description: "AchievementDescVeteran"},
	{ID: "minute_precision", Name: "AchievementNameMinutePrecision", Description: "AchievementDescMinutePrecision"},
	{ID: "poker", Name: "AchievementNamePoker", Description: "AchievementDescPoker"},
	{ID: "keeper", Name: "AchievementNameKeeper", Description: "AchievementDescKeeper"},
	{ID: "cosmic_cock", Name: "AchievementNameCosmicCock", Description: "AchievementDescCosmicCock"},
	{ID: "maximalist", Name: "AchievementNameMaximalist", Description: "AchievementDescMaximalist"},
	{ID: "diamond_hands", Name: "AchievementNameDiamondHands", Description: "AchievementDescDiamondHands"},
	{ID: "diamond_eye", Name: "AchievementNameDiamondEye", Description: "AchievementDescDiamondEye"},
	{ID: "greek_myth", Name: "AchievementNameGreekMyth", Description: "AchievementDescGreekMyth"},
	{ID: "fibonacci_father", Name: "AchievementNameFibonacciFather", Description: "AchievementDescFibonacciFather"},
	{ID: "annihilator_cannon", Name: "AchievementNameAnnihilatorCannon", Description: "AchievementDescAnnihilatorCannon"},
}

// GetAchievementDefByID возвращает клиентское определение достижения по ID
func GetAchievementDefByID(id string) *AchievementDef {
	for i := range AllAchievements {
		if AllAchievements[i].ID == id {
			return &AllAchievements[i]
		}
	}
	return nil
}
