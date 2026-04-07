package partsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// Config — настройки адаптера
type Config struct {
	BaseURL string // обычно https://api.partsapi.ru/ или https://partsapi.ru/api (увидишь в тестере)
	APIKey  string
	Timeout time.Duration // например 10 * time.Second
}

// Adapter — основной адаптер
type Adapter struct {
	client  *http.Client
	baseURL string
	apiKey  string
}

// New создаёт адаптер
func New(cfg Config) *Adapter {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.partsapi.ru" // ← замени после теста в кабинете
	}

	return &Adapter{
		client:  &http.Client{Timeout: cfg.Timeout},
		baseURL: cfg.BaseURL,
		apiKey:  cfg.APIKey,
	}
}

// Part — пример структуры одной запчасти (расширь под реальный ответ из тестера)
type Part struct {
	Name      string `json:"name"`      // название группы (тормозная система и т.д.)
	Shortname string `json:"shortname"` // ?
	Parts     string `json:"parts"`     // строка вида "BOSCH|0986AB1234|ATE|605123|..."
	// иногда может быть массив вместо строки

	// добавь другие поля, которые вернёт API (цена, фото, кроссы и т.д.)
}

type GetPartsResponse struct {
	// общие поля (могут быть, а могут отсутствовать)
	Success bool   `json:"success,omitempty"`
	Error   string `json:"error,omitempty"`

	// основной массив — список групп с запчастями
	Data []Part
}

// GetPartsByVIN возвращает сырые данные из getPartsbyVIN
func (a *Adapter) GetPartsByVIN(ctx context.Context, vin, cat string, isOrigPart bool) ([]Part, error) {
	if len(vin) != 17 {
		return nil, fmt.Errorf("VIN должен быть ровно 17 символов")
	}

	params := url.Values{}
	params.Add("method", "getPartsbyVIN")
	params.Add("key", a.apiKey)
	params.Add("vin", vin)
	params.Add("lang", "ru")
	if isOrigPart {
		params.Add("type", "oem")
	}else {
        params.Add("type", "all") // или ""
    }
	if cat != "" {
		params.Add("cat", "1191")
	}

	fullURL := a.baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ошибка %d: %s", resp.StatusCode, body)
	}

	var groups []Part
    if err := json.NewDecoder(resp.Body).Decode(&groups); err != nil {
        return nil, fmt.Errorf("JSON decode failed: %w", err)
    }

	log.Println(groups)

	return groups, nil
}
