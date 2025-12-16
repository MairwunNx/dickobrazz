package application

import "dickobrazz/application/database"

// AllAchievements содержит все доступные достижения в игре
var AllAchievements = []database.Achievement{
	// Точность и коллекционирование
	{
		ID:          "sniper",
		Emoji:       "🎯",
		Name:        "Снайпер",
		Description: "получить ровно 30см пять раз",
		Respects:    30,
		MaxProgress: 5,
	},
	{
		ID:          "number_collector",
		Emoji:       "🔢",
		Name:        "Коллекционер чисел",
		Description: "получить все красивые числа (11, 22, 33, 44, 55)",
		Respects:    100,
		MaxProgress: 5,
	},
	{
		ID:          "half_hundred",
		Emoji:       "🌟",
		Name:        "Полсотни",
		Description: "получить кок 50см",
		Respects:    50,
		MaxProgress: 50,
	},
	
	// Динамика
	{
		ID:          "bull_trend",
		Emoji:       "📈",
		Name:        "Бычий тренд",
		Description: "рост кока 5 дней подряд",
		Respects:    50,
		MaxProgress: 5,
	},
	{
		ID:          "bear_market",
		Emoji:       "📉",
		Name:        "Медвежий рынок",
		Description: "падение кока 5 дней подряд",
		Respects:    50,
		MaxProgress: 5,
	},
	{
		ID:          "lightning",
		Emoji:       "⚡",
		Name:        "Молния",
		Description: "вырастить кок на 50см за день",
		Respects:    100,
		MaxProgress: 1,
	},
	{
		ID:          "turtle",
		Emoji:       "🐌",
		Name:        "Черепаха",
		Description: "10 коков подряд с изменением меньше 5см",
		Respects:    30,
		MaxProgress: 10,
	},
	
	// Экстремальные значения
	{
		ID:          "everest",
		Emoji:       "🏔️",
		Name:        "Эверест",
		Description: "получить максимальный кок среди всех",
		Respects:    333,
		MaxProgress: 61,
	},
	{
		ID:          "mariana_trench",
		Emoji:       "🕳️",
		Name:        "Марианская впадина",
		Description: "получить минимальный кок среди всех",
		Respects:    333,
		MaxProgress: 0,
	},
	{
		ID:          "freeze",
		Emoji:       "❄️",
		Name:        "Мороз по коже",
		Description: "5 коков подряд меньше 20см",
		Respects:    30,
		MaxProgress: 5,
	},
	{
		ID:          "diamond_hands",
		Emoji:       "💎",
		Name:        "Алмазные руки",
		Description: "7 коков подряд от 40см",
		Respects:    100,
		MaxProgress: 7,
	},
	
	// Временные
	{
		ID:          "early_bird",
		Emoji:       "🌅",
		Name:        "Ранняя пташка",
		Description: "дернуть кок до 6:00 МСК двадцать раз",
		Respects:    100,
		MaxProgress: 20,
	},
	{
		ID:          "speedrunner",
		Emoji:       "⏱️",
		Name:        "Спидраннер",
		Description: "дернуть кок за 30 секунд после полуночи пять раз",
		Respects:    100,
		MaxProgress: 5,
	},
	
	// Сезоны
	{
		ID:          "oldtimer",
		Emoji:       "🗓️",
		Name:        "Старожил",
		Description: "участвовать в 3 сезонах",
		Respects:    100,
		MaxProgress: 3,
	},
	{
		ID:          "veteran",
		Emoji:       "🗓️",
		Name:        "Ветеран",
		Description: "участвовать в 5 сезонах",
		Respects:    300,
		MaxProgress: 5,
	},
	{
		ID:          "keeper",
		Emoji:       "🗓️",
		Name:        "Хранитель",
		Description: "участвовать в 10 сезонах",
		Respects:    1000,
		MaxProgress: 10,
	},
	
	// Последовательности
	{
		ID:          "triple",
		Emoji:       "🎰",
		Name:        "Тройка",
		Description: "получить одинаковый размер 3 раза подряд",
		Respects:    50,
		MaxProgress: 3,
	},
	{
		ID:          "deja_vu",
		Emoji:       "🔄",
		Name:        "Дежавю",
		Description: "получить одинаковый кок два дня подряд",
		Respects:    20,
		MaxProgress: 2,
	},
	{
		ID:          "poker",
		Emoji:       "🎴",
		Name:        "Покер",
		Description: "получить одинаковый размер 4 раза подряд",
		Respects:    100,
		MaxProgress: 4,
	},
	{
		ID:          "diamond_eye",
		Emoji:       "💎",
		Name:        "Глаз алмаз",
		Description: "получить одинаковый размер 5 раз подряд",
		Respects:    500,
		MaxProgress: 5,
	},
	
	// Сложные коллекции
	{
		ID:          "rounder",
		Emoji:       "🔟",
		Name:        "Округлятор",
		Description: "получить все круглые числа (10, 20, 30, 40, 50, 60)",
		Respects:    200,
		MaxProgress: 6,
	},
	{
		ID:          "fibonacci_father",
		Emoji:       "🔢",
		Name:        "Отец фибоначчи",
		Description: "получить последовательность Фибоначчи (1, 2, 3, 5, 8, 13, 21, 34, 55)",
		Respects:    2222,
		MaxProgress: 9,
	},
	
	// География
	{
		ID:          "traveler",
		Emoji:       "🗺️",
		Name:        "Путешественник",
		Description: "получить все 61 размер кока (0-60см)",
		Respects:    500,
		MaxProgress: 61,
	},
	{
		ID:          "muscovite",
		Emoji:       "🏙️",
		Name:        "Москвич",
		Description: "получить кок 50см пять раз за месяц",
		Respects:    100,
		MaxProgress: 5,
	},
	
	// Праздничные
	{
		ID:          "valentine",
		Emoji:       "💝",
		Name:        "Валентинка",
		Description: "получить кок 14см в День Влюблённых",
		Respects:    50,
		MaxProgress: 1,
	},
	{
		ID:          "new_year_gift",
		Emoji:       "🎄",
		Name:        "Новогодний подарок",
		Description: "получить кок 60см+ в Новый Год",
		Respects:    200,
		MaxProgress: 1,
	},
	
	// Накопление размера
	{
		ID:          "golden_hundred",
		Emoji:       "💯",
		Name:        "Золотая сотня",
		Description: "нарастить 100см суммарно",
		Respects:    20,
		MaxProgress: 100,
	},
	{
		ID:          "solid_thousand",
		Emoji:       "💰",
		Name:        "Четкий касарь",
		Description: "нарастить 1000см суммарно",
		Respects:    50,
		MaxProgress: 1000,
	},
	{
		ID:          "five_k",
		Emoji:       "💎",
		Name:        "Пятикат",
		Description: "нарастить 5000см суммарно",
		Respects:    100,
		MaxProgress: 5000,
	},
	{
		ID:          "golden_cock",
		Emoji:       "👑",
		Name:        "Золотой кок",
		Description: "нарастить 10000см суммарно",
		Respects:    300,
		MaxProgress: 10000,
	},
	{
		ID:          "cosmic_cock",
		Emoji:       "🚀",
		Name:        "Космический кок",
		Description: "нарастить 20000см суммарно",
		Respects:    1000,
		MaxProgress: 20000,
	},
	{
		ID:          "greek_myth",
		Emoji:       "⚡",
		Name:        "Миф древней греции",
		Description: "нарастить 30000см суммарно",
		Respects:    2222,
		MaxProgress: 30000,
	},
	
	// Количество дерганий
	{
		ID:          "not_rubbed_yet",
		Emoji:       "🤏",
		Name:        "Еще не натерло",
		Description: "дернуть кок 10 раз",
		Respects:    20,
		MaxProgress: 10,
	},
	{
		ID:          "diary",
		Emoji:       "📆",
		Name:        "Ежедневник",
		Description: "дернуть кок 31 раз",
		Respects:    30,
		MaxProgress: 31,
	},
	{
		ID:          "skillful_hands",
		Emoji:       "💪",
		Name:        "Очумелые ручки",
		Description: "дернуть кок 100 раз",
		Respects:    50,
		MaxProgress: 100,
	},
	{
		ID:          "wonder_stranger",
		Emoji:       "💋",
		Name:        "Чудо незнакомка",
		Description: "дернуть кок 500 раз",
		Respects:    200,
		MaxProgress: 500,
	},
	{
		ID:          "bazooka_hands",
		Emoji:       "💥",
		Name:        "Руки базуки",
		Description: "дернуть кок 1000 раз",
		Respects:    500,
		MaxProgress: 1000,
	},
	{
		ID:          "anniversary",
		Emoji:       "🎂",
		Name:        "Годовщина",
		Description: "дернуть кок 365 раз",
		Respects:    500,
		MaxProgress: 365,
	},
	{
		ID:          "annihilator_cannon",
		Emoji:       "☢️",
		Name:        "Аннигиляторная пушка",
		Description: "дернуть кок 5000 раз",
		Respects:    2222,
		MaxProgress: 5000,
	},
	
	// Специальные совпадения
	{
		ID:          "sum_of_previous",
		Emoji:       "🎲",
		Name:        "Сумма предыдущих",
		Description: "получить кок равный сумме двух предыдущих",
		Respects:    1000,
		MaxProgress: 0,
	},
	{
		ID:          "minute_precision",
		Emoji:       "⏰",
		Name:        "Минутная точность",
		Description: "получить кок равный минутам времени",
		Respects:    1500,
		MaxProgress: 0,
	},
	{
		ID:          "hour_precision",
		Emoji:       "🕐",
		Name:        "Часовая точность",
		Description: "получить кок равный часу времени",
		Respects:    500,
		MaxProgress: 0,
	},
	{
		ID:          "day_equals_size",
		Emoji:       "📅",
		Name:        "День = Размер",
		Description: "получить кок равный дню месяца",
		Respects:    300,
		MaxProgress: 0,
	},
	{
		ID:          "contrast_shower",
		Emoji:       "🚿",
		Name:        "Контрастный душ",
		Description: "получить 0-3см сразу после 60+см",
		Respects:    800,
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
