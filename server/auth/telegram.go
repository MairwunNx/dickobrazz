package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"dickobrazz/server/models"
)

var errInvalidTelegramAuth = errors.New("invalid telegram auth data")

func ValidateTelegramAuth(payload models.TelegramAuthPayload, token string) error {
	dataCheck := buildDataCheckString(payload)
	secret := sha256.Sum256([]byte(token))
	hash := hmacSha256(secret[:], dataCheck)
	if !hmac.Equal([]byte(payload.Hash), []byte(hash)) {
		return errInvalidTelegramAuth
	}
	return nil
}

func ParseInitData(raw string) (models.TelegramAuthPayload, error) {
	values, err := url.ParseQuery(raw)
	if err != nil {
		return models.TelegramAuthPayload{}, err
	}

	getInt64 := func(key string) (int64, error) {
		value := strings.TrimSpace(values.Get(key))
		if value == "" {
			return 0, errors.New("missing " + key)
		}
		return strconv.ParseInt(value, 10, 64)
	}

	authDate, err := getInt64("auth_date")
	if err != nil {
		return models.TelegramAuthPayload{}, err
	}

	id, err := getInt64("id")
	if err != nil {
		return models.TelegramAuthPayload{}, err
	}

	payload := models.TelegramAuthPayload{
		ID:        id,
		FirstName: values.Get("first_name"),
		LastName:  values.Get("last_name"),
		Username:  values.Get("username"),
		PhotoURL:  values.Get("photo_url"),
		AuthDate:  authDate,
		Hash:      values.Get("hash"),
	}

	if payload.Hash == "" {
		return models.TelegramAuthPayload{}, errInvalidTelegramAuth
	}

	return payload, nil
}

func buildDataCheckString(payload models.TelegramAuthPayload) string {
	data := map[string]string{
		"auth_date":  strconv.FormatInt(payload.AuthDate, 10),
		"first_name": payload.FirstName,
		"id":         strconv.FormatInt(payload.ID, 10),
		"last_name":  payload.LastName,
		"photo_url":  payload.PhotoURL,
		"username":   payload.Username,
	}

	keys := make([]string, 0, len(data))
	for key, value := range data {
		if value == "" {
			continue
		}
		keys = append(keys, key)
	}
	sort.Strings(keys)

	lines := make([]string, 0, len(keys))
	for _, key := range keys {
		lines = append(lines, key+"="+data[key])
	}
	return strings.Join(lines, "\n")
}

func hmacSha256(secret []byte, data string) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
