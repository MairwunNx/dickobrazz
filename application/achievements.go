package application

import "dickobrazz/application/database"

// AllAchievements содержит все доступные достижения в игре
// Отсортированы по возрастанию респектов (33-2222)
var AllAchievements = []database.Achievement{
	// 33 респекта
	{
		ID:          "not_rubbed_yet",
		Emoji:       "🤏",
		Name:        "Еще не натерло",
		Description: "дернуть кок 10 раз",
		Respects:    33,
		MaxProgress: 10,
	},
	
	// 50 респектов
	{
		ID:          "half_hundred",
		Emoji:       "🌟",
		Name:        "Полсотни",
		Description: "получить кок 50см",
		Respects:    50,
		MaxProgress: 50,
	},
	{
		ID:          "diary",
		Emoji:       "📆",
		Name:        "Ежедневник",
		Description: "дернуть кок 31 раз",
		Respects:    50,
		MaxProgress: 31,
	},
	
	// 90 респектов
	{
		ID:          "golden_hundred",
		Emoji:       "💯",
		Name:        "Золотая сотня",
		Description: "нарастить 100см суммарно",
		Respects:    90,
		MaxProgress: 100,
	},
	
	// 100 респектов
	{
		ID:          "skillful_hands",
		Emoji:       "💪",
		Name:        "Очумелые ручки",
		Description: "дернуть кок 100 раз",
		Respects:    100,
		MaxProgress: 100,
	},
	
	// 135 респектов
	{
		ID:          "early_bird",
		Emoji:       "🌅",
		Name:        "Ранняя пташка",
		Description: "дернуть кок до 6:00 МСК двадцать раз",
		Respects:    135,
		MaxProgress: 20,
	},
	
	// 200 респектов
	{
		ID:          "lightning",
		Emoji:       "⚡",
		Name:        "Молния",
		Description: "нарастить кок на 50см относительно предыдущего кока",
		Respects:    200,
		MaxProgress: 1,
	},
	
	// 211 респектов
	{
		ID:          "sniper",
		Emoji:       "🎯",
		Name:        "Снайпер",
		Description: "получить ровно 30см пять раз",
		Respects:    211,
		MaxProgress: 5,
	},
	
	// 222 респекта
	{
		ID:          "deja_vu",
		Emoji:       "🔄",
		Name:        "Дежавю",
		Description: "получить одинаковый кок два дня подряд",
		Respects:    222,
		MaxProgress: 2,
	},
	{
		ID:          "speedrunner",
		Emoji:       "⏱️",
		Name:        "Спидраннер",
		Description: "дернуть кок за 30 секунд после полуночи пять раз",
		Respects:    222,
		MaxProgress: 5,
	},
	{
		ID:          "midnight_puller",
		Emoji:       "🌙",
		Name:        "Полуночник",
		Description: "дернуть кок после 23:00 МСК пятьдесят раз",
		Respects:    222,
		MaxProgress: 50,
	},
	
	// 228 респектов
	{
		ID:          "rounder",
		Emoji:       "🔟",
		Name:        "Округлятор",
		Description: "получить все круглые числа (10, 20, 30, 40, 50, 60) за последние 31 коков",
		Respects:    228,
		MaxProgress: 6,
	},
	{
		ID:          "everest",
		Emoji:       "🏔️",
		Name:        "Эверест",
		Description: "получить максимальный кок среди всех",
		Respects:    228,
		MaxProgress: 61,
	},
	{
		ID:          "mariana_trench",
		Emoji:       "🕳️",
		Name:        "Марианская впадина",
		Description: "получить минимальный кок среди всех",
		Respects:    228,
		MaxProgress: 0,
	},
	
	// 233 респекта
	{
		ID:          "number_collector",
		Emoji:       "🔢",
		Name:        "Коллекционер чисел",
		Description: "получить все красивые числа (11, 22, 33, 44, 55)",
		Respects:    233,
		MaxProgress: 5,
	},
	
	// 300 респектов
	{
		ID:          "day_equals_size",
		Emoji:       "📅",
		Name:        "День = Размер",
		Description: "получить кок равный дню месяца",
		Respects:    300,
		MaxProgress: 0,
	},
	
	// 333 респекта
	{
		ID:          "solid_thousand",
		Emoji:       "💰",
		Name:        "Четкий касарь",
		Description: "нарастить 1000см суммарно",
		Respects:    333,
		MaxProgress: 1000,
	},
	{
		ID:          "bull_trend",
		Emoji:       "📈",
		Name:        "Бычий тренд",
		Description: "рост кока 5 дней подряд",
		Respects:    333,
		MaxProgress: 5,
	},
	{
		ID:          "bear_market",
		Emoji:       "📉",
		Name:        "Медвежий рынок",
		Description: "падение кока 5 дней подряд",
		Respects:    333,
		MaxProgress: 5,
	},
	
	// 500 респектов
	{
		ID:          "traveler",
		Emoji:       "🗺️",
		Name:        "Путешественник",
		Description: "получить все 61 размер кока (0-60см) за все время",
		Respects:    500,
		MaxProgress: 61,
	},
	
	// 555 респектов
	{
		ID:          "freeze",
		Emoji:       "❄️",
		Name:        "Мороз по коже",
		Description: "5 коков подряд меньше 20см",
		Respects:    555,
		MaxProgress: 5,
	},
	
	// 700 респектов
	{
		ID:          "five_k",
		Emoji:       "💎",
		Name:        "Пятикат",
		Description: "нарастить 5000см суммарно",
		Respects:    700,
		MaxProgress: 5000,
	},
	{
		ID:          "oldtimer",
		Emoji:       "🗓️",
		Name:        "Старожил",
		Description: "участвовать в 3 сезонах",
		Respects:    700,
		MaxProgress: 3,
	},
	{
		ID:          "anniversary",
		Emoji:       "🎂",
		Name:        "Годовщина",
		Description: "дернуть кок 365 раз",
		Respects:    700,
		MaxProgress: 365,
	},
	
	// 777 респектов
	{
		ID:          "contrast_shower",
		Emoji:       "🚿",
		Name:        "Контрастный душ",
		Description: "получить 0-3см сразу после 60+см",
		Respects:    777,
		MaxProgress: 0,
	},
	{
		ID:          "pythagoras",
		Emoji:       "📐",
		Name:        "Пифагор",
		Description: "получить три кока подряд, образующих пифагорову тройку (3-4-5, 5-12-13, 8-15-17 и т.д.)",
		Respects:    777,
		MaxProgress: 1,
	},
	{
		ID:          "leet_speak",
		Emoji:       "💻",
		Name:        "Leet speak",
		Description: "получить 13см и 37см подряд",
		Respects:    777,
		MaxProgress: 1,
	},
	
	// 800 респектов
	{
		ID:          "moscovite",
		Emoji:       "🏙️",
		Name:        "Москвич",
		Description: "получить кок 50см пять раз за месяц",
		Respects:    800,
		MaxProgress: 5,
	},
	
	// 888 респектов
	{
		ID:          "hour_precision",
		Emoji:       "🕐",
		Name:        "Часовая точность",
		Description: "получить кок равный часу времени",
		Respects:    888,
		MaxProgress: 0,
	},
	
	// 900 респектов
	{
		ID:          "wonder_stranger",
		Emoji:       "💋",
		Name:        "Чудо незнакомка",
		Description: "дернуть кок 500 раз",
		Respects:    900,
		MaxProgress: 500,
	},
	
	// 999 респектов
	{
		ID:          "valentine",
		Emoji:       "💝",
		Name:        "Валентинка",
		Description: "получить кок 14см в День Влюблённых",
		Respects:    999,
		MaxProgress: 1,
	},
	{
		ID:          "new_year_gift",
		Emoji:       "🎄",
		Name:        "Новогодний подарок",
		Description: "получить кок 60см+ в Новый Год",
		Respects:    999,
		MaxProgress: 1,
	},
	{
		ID:          "mens_solidarity",
		Emoji:       "🤝",
		Name:        "Мужская солидарность",
		Description: "получить кок 19см в Международный мужской день (19 ноября)",
		Respects:    999,
		MaxProgress: 1,
	},
	{
		ID:          "friday_13th",
		Emoji:       "☠️",
		Name:        "Пятница 13",
		Description: "получить кок 0см в пятницу 13-го",
		Respects:    999,
		MaxProgress: 1,
	},
	{
		ID:          "leap_cock",
		Emoji:       "📅",
		Name:        "Високосный кок",
		Description: "получить любой кок 29 февраля",
		Respects:    999,
		MaxProgress: 1,
	},
	
	// 1000 респектов
	{
		ID:          "turtle",
		Emoji:       "🐌",
		Name:        "Черепаха",
		Description: "10 коков подряд с изменением меньше 5см",
		Respects:    1000,
		MaxProgress: 10,
	},
	{
		ID:          "golden_cock",
		Emoji:       "👑",
		Name:        "Золотой кок",
		Description: "нарастить 10000см суммарно",
		Respects:    1000,
		MaxProgress: 10000,
	},
	{
		ID:          "sum_of_previous",
		Emoji:       "🎲",
		Name:        "Сумма предыдущих",
		Description: "получить кок равный сумме двух предыдущих",
		Respects:    1000,
		MaxProgress: 0,
	},
	
	// 1222 респекта
	{
		ID:          "bazooka_hands",
		Emoji:       "💥",
		Name:        "Руки базуки",
		Description: "дернуть кок 1000 раз",
		Respects:    1222,
		MaxProgress: 1000,
	},
	
	// 1333 респекта
	{
		ID:          "triple",
		Emoji:       "🎰",
		Name:        "Тройка",
		Description: "получить одинаковый размер 3 раза подряд",
		Respects:    1333,
		MaxProgress: 3,
	},
	{
		ID:          "veteran",
		Emoji:       "🗓️",
		Name:        "Ветеран",
		Description: "участвовать в 5 сезонах",
		Respects:    1333,
		MaxProgress: 5,
	},
	{
		ID:          "minute_precision",
		Emoji:       "⏰",
		Name:        "Минутная точность",
		Description: "получить кок равный минутам времени",
		Respects:    1333,
		MaxProgress: 0,
	},
	
	// 1777 респектов
	{
		ID:          "poker",
		Emoji:       "🎴",
		Name:        "Покер",
		Description: "получить одинаковый размер 4 раза подряд",
		Respects:    1777,
		MaxProgress: 4,
	},
	
	// 1888 респектов
	{
		ID:          "keeper",
		Emoji:       "🗓️",
		Name:        "Хранитель",
		Description: "участвовать в 10 сезонах",
		Respects:    1888,
		MaxProgress: 10,
	},
	
	// 2000 респектов
	{
		ID:          "cosmic_cock",
		Emoji:       "🚀",
		Name:        "Космический кок",
		Description: "нарастить 20000см суммарно",
		Respects:    2000,
		MaxProgress: 20000,
	},
	{
		ID:          "maximalist",
		Emoji:       "🔝",
		Name:        "Максималист",
		Description: "получить 61см десять раз",
		Respects:    2000,
		MaxProgress: 10,
	},
	
	// 2222 респекта
	{
		ID:          "diamond_hands",
		Emoji:       "💎",
		Name:        "Алмазные руки",
		Description: "7 коков подряд от 40см",
		Respects:    2222,
		MaxProgress: 7,
	},
	{
		ID:          "diamond_eye",
		Emoji:       "💎",
		Name:        "Глаз алмаз",
		Description: "получить одинаковый размер 5 раз подряд",
		Respects:    2222,
		MaxProgress: 5,
	},
	{
		ID:          "greek_myth",
		Emoji:       "⚡",
		Name:        "Миф древней греции",
		Description: "нарастить 30000см суммарно",
		Respects:    2222,
		MaxProgress: 30000,
	},
	{
		ID:          "fibonacci_father",
		Emoji:       "🔢",
		Name:        "Отец фибоначчи",
		Description: "получить последовательность Фибоначчи (1, 2, 3, 5, 8, 13, 21, 34, 55) за последние 31 коков",
		Respects:    2222,
		MaxProgress: 9,
	},
	{
		ID:          "annihilator_cannon",
		Emoji:       "☢️",
		Name:        "Аннигиляторная пушка",
		Description: "дернуть кок 5000 раз",
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
