package achievement

// AchievementDef содержит клиентское определение достижения (ID + ключи локализации)
type AchievementDef struct {
	ID          string
	Name        string // ключ локализации
	Description string // ключ локализации
}
