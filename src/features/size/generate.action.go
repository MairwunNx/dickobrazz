package size

import (
	"context"
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/emoji"
	"dickobrazz/src/shared/geo"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type GenerateResult struct {
	Text    string
	Subtext string
}

type GenerateAction struct {
	api *api.APIClient
	loc *localization.LocalizationManager
}

func NewGenerateAction(apiClient *api.APIClient, loc *localization.LocalizationManager) *GenerateAction {
	return &GenerateAction{api: apiClient, loc: loc}
}

func (a *GenerateAction) Execute(ctx context.Context, log *logging.Logger, localizer *i18n.Localizer, userID int64, username string) (*GenerateResult, error) {
	cockData, err := a.api.GenerateCockSize(ctx, userID, username)
	if err != nil {
		log.E("Failed to generate cock size via API", logging.InnerError, err)
		return nil, err
	}

	size := cockData.Size
	emojiStr := emoji.EmojiFromSize(size)
	text := GenerateCockSizeText(a.loc, localizer, size, emojiStr)
	subtext := geo.GetRegionBySize(size)
	subtext = a.loc.Localize(localizer, subtext, nil)

	return &GenerateResult{Text: text, Subtext: subtext}, nil
}
