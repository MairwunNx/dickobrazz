package achievements

import (
	"dickobrazz/src/entities/achievement"
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/formatting"
	"dickobrazz/src/shared/localization"
	"sort"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GenerateAchievementsText(
	localizationManager *localization.LocalizationManager,
	localizer *i18n.Localizer,
	allAchievements []achievement.AchievementDef,
	apiAchievements []api.AchievementData,
	page int,
	itemsPerPage int,
) string {
	apiMap := make(map[string]*api.AchievementData, len(apiAchievements))
	for i := range apiAchievements {
		apiMap[apiAchievements[i].ID] = &apiAchievements[i]
	}

	type AchievementWithStatus struct {
		Def         achievement.AchievementDef
		APIData     *api.AchievementData
		IsCompleted bool
	}

	achievementsWithStatus := make([]AchievementWithStatus, 0, len(allAchievements))

	for _, def := range allAchievements {
		apiData := apiMap[def.ID]
		isCompleted := apiData != nil && apiData.Completed

		achievementsWithStatus = append(achievementsWithStatus, AchievementWithStatus{
			Def:         def,
			APIData:     apiData,
			IsCompleted: isCompleted,
		})
	}

	sort.Slice(achievementsWithStatus, func(i, j int) bool {
		if achievementsWithStatus[i].IsCompleted != achievementsWithStatus[j].IsCompleted {
			return achievementsWithStatus[i].IsCompleted
		}
		return false
	})

	totalPages := (len(achievementsWithStatus) + itemsPerPage - 1) / itemsPerPage
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}

	startIdx := (page - 1) * itemsPerPage
	endIdx := startIdx + itemsPerPage
	if endIdx > len(achievementsWithStatus) {
		endIdx = len(achievementsWithStatus)
	}

	var lines []string
	for i := startIdx; i < endIdx; i++ {
		achStatus := achievementsWithStatus[i]
		line := FormatAchievementLine(localizationManager, localizer, achStatus.Def, achStatus.APIData, achStatus.IsCompleted)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

func FormatAchievementLine(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, def achievement.AchievementDef, apiData *api.AchievementData, isCompleted bool) string {
	escapedName := formatting.EscapeMarkdownV2(localizationManager.Localize(localizer, def.Name, nil))
	escapedDesc := formatting.EscapeMarkdownV2(localizationManager.Localize(localizer, def.Description, nil))

	emoji := "🔒"
	if apiData != nil {
		emoji = apiData.Emoji
	}

	if isCompleted {
		return localizationManager.Localize(localizer, localization.MsgAchievementCompleted, map[string]any{
			"Emoji":       emoji,
			"Name":        escapedName,
			"Description": escapedDesc,
		})
	} else if apiData != nil && apiData.Progress > 0 && apiData.MaxProgress > 0 {
		return localizationManager.Localize(localizer, localization.MsgAchievementInProgress, map[string]any{
			"Emoji":       emoji,
			"Name":        escapedName,
			"Progress":    apiData.Progress,
			"Max":         apiData.MaxProgress,
			"Description": escapedDesc,
		})
	}
	return localizationManager.Localize(localizer, localization.MsgAchievementNotCompleted, map[string]any{
		"Emoji":       emoji,
		"Name":        escapedName,
		"Description": escapedDesc,
	})
}
