package application

import "dickobrazz/application/database"

// AllAchievements содержит все доступные достижения в игре
var AllAchievements = []database.Achievement{
	// Точность и коллекционирование
	{
		ID:          "sniper",
		Emoji:       "🎯",
		Name:        "Снайпер",
		Description: "получил ровно 30см 5 раз",
		Respects:    30,
		MaxProgress: 5,
	},
	{
		ID:          "number_collector",
		Emoji:       "🔢",
		Name:        "Коллекционер чисел",
		Description: "получил все \"красивые\" числа (11, 22, 33, 44, 55)",
		Respects:    100,
		MaxProgress: 5,
	},
	{
		ID:          "half_hundred",
		Emoji:       "🌟",
		Name:        "Полсотни",
		Description: "получил ровно 50см",
		Respects:    50,
		MaxProgress: 0,
	},
	
	// Динамика
	{
		ID:          "bull_trend",
		Emoji:       "📈",
		Name:        "Бычий тренд",
		Description: "5 дней подряд рост кока",
		Respects:    50,
		MaxProgress: 0,
	},
	{
		ID:          "bear_market",
		Emoji:       "📉",
		Name:        "Медвежий рынок",
		Description: "5 дней подряд падение кока",
		Respects:    50,
		MaxProgress: 0,
	},
	{
		ID:          "lightning",
		Emoji:       "⚡",
		Name:        "Молния",
		Description: "рост на 50+см за один день",
		Respects:    100,
		MaxProgress: 0,
	},
	{
		ID:          "turtle",
		Emoji:       "🐌",
		Name:        "Черепаха",
		Description: "10 коков с изменением <5см",
		Respects:    30,
		MaxProgress: 10,
	},
	
	// Экстремальные значения
	{
		ID:          "everest",
		Emoji:       "🏔️",
		Name:        "Эверест",
		Description: "получил максимальный кок в системе",
		Respects:    333,
		MaxProgress: 0,
	},
	{
		ID:          "mariana_trench",
		Emoji:       "🕳️",
		Name:        "Марианская впадина",
		Description: "получил минимальный кок в системе",
		Respects:    333,
		MaxProgress: 0,
	},
	{
		ID:          "freeze",
		Emoji:       "❄️",
		Name:        "Мороз по коже",
		Description: "5 коков подряд с коком <20см",
		Respects:    30,
		MaxProgress: 0,
	},
	{
		ID:          "diamond_hands",
		Emoji:       "💎",
		Name:        "Алмазные руки",
		Description: "7 коков подряд 40+см",
		Respects:    100,
		MaxProgress: 0,
	},
	
	// Временные
	{
		ID:          "early_bird",
		Emoji:       "🌅",
		Name:        "Ранняя пташка",
		Description: "первый кок дня (до 6:00 МСК) 20 раз",
		Respects:    100,
		MaxProgress: 20,
	},
	{
		ID:          "speedrunner",
		Emoji:       "⏱️",
		Name:        "Спидраннер",
		Description: "получил кок за <30 сек после полуночи 5 раз",
		Respects:    100,
		MaxProgress: 5,
	},
	
	// Сезоны
	{
		ID:          "oldtimer",
		Emoji:       "🗓️",
		Name:        "Старожил",
		Description: "участвовал в 3+ сезонах",
		Respects:    100,
		MaxProgress: 0,
	},
	{
		ID:          "veteran",
		Emoji:       "🗓️",
		Name:        "Ветеран",
		Description: "участвовал в 5+ сезонах",
		Respects:    300,
		MaxProgress: 0,
	},
	{
		ID:          "keeper",
		Emoji:       "🗓️",
		Name:        "Хранитель",
		Description: "участвовал в 10+ сезонах",
		Respects:    1000,
		MaxProgress: 0,
	},
	
	// Последовательности
	{
		ID:          "triple",
		Emoji:       "🎰",
		Name:        "Тройка",
		Description: "получил один и тот же размер 3 раза подряд",
		Respects:    50,
		MaxProgress: 0,
	},
	{
		ID:          "deja_vu",
		Emoji:       "🔄",
		Name:        "Дежавю",
		Description: "получил одинаковый кок сегодня и вчера",
		Respects:    20,
		MaxProgress: 0,
	},
	{
		ID:          "poker",
		Emoji:       "🎴",
		Name:        "Покер",
		Description: "получил 4 одинаковых размера подряд",
		Respects:    100,
		MaxProgress: 0,
	},
	{
		ID:          "diamond_eye",
		Emoji:       "💎",
		Name:        "Глаз алмаз",
		Description: "5 одинаковых коков подряд",
		Respects:    500,
		MaxProgress: 0,
	},
	
	// Сложные коллекции
	{
		ID:          "rounder",
		Emoji:       "🔟",
		Name:        "Округлятор",
		Description: "получал 10, 20, 30, 40, 50, 60см за 31 кок",
		Respects:    200,
		MaxProgress: 6,
	},
	{
		ID:          "fibonacci_father",
		Emoji:       "🔢",
		Name:        "Отец фибоначчи",
		Description: "получил 1, 1, 2, 3, 5, 8, 13, 21, 34, 55см за 31 кок",
		Respects:    2222,
		MaxProgress: 10,
	},
	
	// География
	{
		ID:          "traveler",
		Emoji:       "🗺️",
		Name:        "Путешественник",
		Description: "прошел все регионы России по размерам",
		Respects:    500,
		MaxProgress: 0,
	},
	{
		ID:          "muscovite",
		Emoji:       "🏙️",
		Name:        "Москвич",
		Description: "получил размер \"Москва\" 5 раз за 31 день",
		Respects:    100,
		MaxProgress: 5,
	},
	
	// Праздничные
	{
		ID:          "valentine",
		Emoji:       "💝",
		Name:        "Валентинка",
		Description: "получить 14см кок 14 февраля",
		Respects:    50,
		MaxProgress: 0,
	},
	{
		ID:          "new_year_gift",
		Emoji:       "🎄",
		Name:        "Новогодний подарок",
		Description: "получить 60+см кок 31 декабря",
		Respects:    200,
		MaxProgress: 0,
	},
	
	// Накопление размера
	{
		ID:          "golden_hundred",
		Emoji:       "💯",
		Name:        "Золотая сотня",
		Description: "нарастить 100см кока",
		Respects:    20,
		MaxProgress: 0,
	},
	{
		ID:          "solid_thousand",
		Emoji:       "💰",
		Name:        "Четкий касарь",
		Description: "нарастить 1000см кока",
		Respects:    50,
		MaxProgress: 0,
	},
	{
		ID:          "five_k",
		Emoji:       "💎",
		Name:        "Пятикат",
		Description: "нарастить 5000см кока",
		Respects:    100,
		MaxProgress: 0,
	},
	{
		ID:          "golden_cock",
		Emoji:       "👑",
		Name:        "Золотой кок",
		Description: "нарастить 10000см",
		Respects:    300,
		MaxProgress: 0,
	},
	{
		ID:          "cosmic_cock",
		Emoji:       "🚀",
		Name:        "Космический кок",
		Description: "нарастить 20000см кок",
		Respects:    1000,
		MaxProgress: 0,
	},
	{
		ID:          "greek_myth",
		Emoji:       "⚡",
		Name:        "Миф древней греции",
		Description: "нарастить 30000см кок",
		Respects:    2222,
		MaxProgress: 0,
	},
	
	// Количество дерганий
	{
		ID:          "not_rubbed_yet",
		Emoji:       "🤏",
		Name:        "Еще не натерло",
		Description: "дернуть 10 раз кок",
		Respects:    20,
		MaxProgress: 0,
	},
	{
		ID:          "diary",
		Emoji:       "📆",
		Name:        "Ежедневник",
		Description: "дернул кок 31 раз",
		Respects:    30,
		MaxProgress: 0,
	},
	{
		ID:          "skillful_hands",
		Emoji:       "💪",
		Name:        "Очумелые ручки",
		Description: "дернуть 100 раз кок",
		Respects:    50,
		MaxProgress: 0,
	},
	{
		ID:          "wonder_stranger",
		Emoji:       "🔥",
		Name:        "Чудо незнакомка",
		Description: "дернуть 500 раз кок",
		Respects:    200,
		MaxProgress: 0,
	},
	{
		ID:          "bazooka_hands",
		Emoji:       "💥",
		Name:        "Руки базуки",
		Description: "дернуть 1000 раз кок",
		Respects:    500,
		MaxProgress: 0,
	},
	{
		ID:          "anniversary",
		Emoji:       "🎂",
		Name:        "Годовщина",
		Description: "дернул кок 365 раз",
		Respects:    500,
		MaxProgress: 0,
	},
	{
		ID:          "annihilator_cannon",
		Emoji:       "☢️",
		Name:        "Аннигиляторная пушка",
		Description: "дернуть кок 5000 раз",
		Respects:    2222,
		MaxProgress: 0,
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
