package api

// Методы APIClient для вызова Dickobrazz Backend API (Cock, Profile).
//
// Авторизация через CSOT: x-internal-token задаётся в NewAPIClient,
// x-internal-user-id и x-internal-user-name передаются в setUserHeaders.
// Для protected endpoints userID обязателен; username рекомендуется для персонализированных данных.

import (
	"context"
	"fmt"
)

// GenerateCockSize вызывает POST /api/v1/cock/size (стянуть кокич).
// Protected endpoint. Генерирует размер на сегодня и возвращает hash/salt для верификации.
// userID, username — для CSOT-контекста (обязательны).
func (c *APIClient) GenerateCockSize(ctx context.Context, userID int64, username string) (*CockSizeData, error) {
	var result DataResponse[CockSizeData]
	req := c.client.R().SetContext(ctx).SetResult(&result)
	setUserHeaders(req, userID, username)

	resp, err := req.Post("/api/v1/cock/size")
	if err != nil {
		return nil, fmt.Errorf("generate cock size: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetCockRuler вызывает GET /api/v1/cock/ruler (дейли-топ лидерборд, кок-рулетка).
// Публичный endpoint. userID/username опциональны — при передаче заполняются user_position и neighborhood.
// limit — 13–50, page — ≥1.
func (c *APIClient) GetCockRuler(ctx context.Context, userID int64, username string, limit, page int) (*CockRulerData, error) {
	var result DataResponse[CockRulerData]
	req := c.client.R().SetContext(ctx).SetResult(&result).
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetQueryParam("page", fmt.Sprintf("%d", page))
	setUserHeaders(req, userID, username)

	resp, err := req.Get("/api/v1/cock/ruler")
	if err != nil {
		return nil, fmt.Errorf("get cock ruler: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetCockRace вызывает GET /api/v1/cock/race (гонка коков, лидерборд сезона).
// Публичный endpoint. Персонализация при userID/username. limit — 13–50, page — ≥1.
func (c *APIClient) GetCockRace(ctx context.Context, userID int64, username string, limit, page int) (*CockRaceData, error) {
	var result DataResponse[CockRaceData]
	req := c.client.R().SetContext(ctx).SetResult(&result).
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetQueryParam("page", fmt.Sprintf("%d", page))
	setUserHeaders(req, userID, username)

	resp, err := req.Get("/api/v1/cock/race")
	if err != nil {
		return nil, fmt.Errorf("get cock race: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetCockLadder вызывает GET /api/v1/cock/ladder (ладдер коков, вечный лидерборд).
// Публичный endpoint. Персонализация при userID/username. limit — 13–50, page — ≥1.
func (c *APIClient) GetCockLadder(ctx context.Context, userID int64, username string, limit, page int) (*CockLadderData, error) {
	var result DataResponse[CockLadderData]
	req := c.client.R().SetContext(ctx).SetResult(&result).
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetQueryParam("page", fmt.Sprintf("%d", page))
	setUserHeaders(req, userID, username)

	resp, err := req.Get("/api/v1/cock/ladder")
	if err != nil {
		return nil, fmt.Errorf("get cock ladder: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetCockDynamicGlobal вызывает GET /api/v1/cock/dynamic/global.
// Публичный endpoint. Общая динамика без авторизации (user-контекст не влияет на ответ).
func (c *APIClient) GetCockDynamicGlobal(ctx context.Context) (*CockDynamicGlobalData, error) {
	var result DataResponse[CockDynamicGlobalData]
	req := c.client.R().SetContext(ctx).SetResult(&result)

	resp, err := req.Get("/api/v1/cock/dynamic/global")
	if err != nil {
		return nil, fmt.Errorf("get cock dynamic global: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetCockDynamicPersonal вызывает GET /api/v1/cock/dynamic/personal.
// Protected endpoint. Персональная динамика коков (irk, dominance, record и т.д.).
func (c *APIClient) GetCockDynamicPersonal(ctx context.Context, userID int64, username string) (*CockDynamicPersonalData, error) {
	var result DataResponse[CockDynamicPersonalData]
	req := c.client.R().SetContext(ctx).SetResult(&result)
	setUserHeaders(req, userID, username)

	resp, err := req.Get("/api/v1/cock/dynamic/personal")
	if err != nil {
		return nil, fmt.Errorf("get cock dynamic personal: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetCockAchievements вызывает GET /api/v1/cock/achievements.
// Protected endpoint. Список достижений и прогресс пользователя.
func (c *APIClient) GetCockAchievements(ctx context.Context, userID int64, username string) (*CockAchievementsData, error) {
	var result DataResponse[CockAchievementsData]
	req := c.client.R().SetContext(ctx).SetResult(&result)
	setUserHeaders(req, userID, username)

	resp, err := req.Get("/api/v1/cock/achievements")
	if err != nil {
		return nil, fmt.Errorf("get cock achievements: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetCockSeasons вызывает GET /api/v1/cock/seasons.
// Публичный endpoint. Сезоны с победителями (топ-3). user_season_wins при userID/username.
// limit — 13–50, page — ≥1.
func (c *APIClient) GetCockSeasons(ctx context.Context, userID int64, username string, limit, page int) (*CockSeasonsData, error) {
	var result DataResponse[CockSeasonsData]
	req := c.client.R().SetContext(ctx).SetResult(&result).
		SetQueryParam("limit", fmt.Sprintf("%d", limit)).
		SetQueryParam("page", fmt.Sprintf("%d", page))
	setUserHeaders(req, userID, username)

	resp, err := req.Get("/api/v1/cock/seasons")
	if err != nil {
		return nil, fmt.Errorf("get cock seasons: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetCockRespects вызывает GET /api/v1/cock/respects.
// Protected endpoint. Cock-Respect™: season_respect, achievement_respect, total_respect.
func (c *APIClient) GetCockRespects(ctx context.Context, userID int64, username string) (*RespectData, error) {
	var result DataResponse[RespectData]
	req := c.client.R().SetContext(ctx).SetResult(&result)
	setUserHeaders(req, userID, username)

	resp, err := req.Get("/api/v1/cock/respects")
	if err != nil {
		return nil, fmt.Errorf("get cock respects: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// UpdatePrivacy вызывает PATCH /api/v1/me/privacy.
// Protected endpoint. Обновляет is_hidden (скрывает/показывает пользователя в лидербордах).
func (c *APIClient) UpdatePrivacy(ctx context.Context, userID int64, username string, isHidden bool) (*UserProfile, error) {
	var result DataResponse[UserProfile]
	req := c.client.R().SetContext(ctx).SetResult(&result).
		SetBody(&UpdatePrivacyPayload{IsHidden: isHidden})
	setUserHeaders(req, userID, username)

	resp, err := req.Patch("/api/v1/me/privacy")
	if err != nil {
		return nil, fmt.Errorf("update privacy: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result.Data, nil
}

// GetProfile вызывает GET /api/v1/me.
// Protected endpoint. Профиль текущего пользователя (id, username, is_hidden, created_at).
func (c *APIClient) GetProfile(ctx context.Context, userID int64, username string) (*UserProfile, error) {
	var result DataResponse[UserProfile]
	req := c.client.R().SetContext(ctx).SetResult(&result)
	setUserHeaders(req, userID, username)

	resp, err := req.Get("/api/v1/me")
	if err != nil {
		return nil, fmt.Errorf("get profile: %w", err)
	}
	if err := checkResponse(resp); err != nil {
		return nil, err
	}
	return &result.Data, nil
}
