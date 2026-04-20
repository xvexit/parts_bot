package partsapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Config — настройки адаптера
type Config struct {
	BaseURL string // обычно https://api.partsapi.ru/ или https://partsapi.ru/api (увидишь в тестере)
	APIKey  string
	// APIKeyVIN используется для метода VINdecode.
	// Если не задан, используется APIKey.
	APIKeyVIN string
	// APIKeyTree используется для метода getSearchTree.
	// Если не задан, используется APIKey.
	APIKeyTree string
	Timeout    time.Duration // например 10 * time.Second
}

// Adapter — основной адаптер
type Adapter struct {
	client     *http.Client
	baseURL    string
	apiKey     string
	vinAPIKey  string
	treeAPIKey string
}

// New создаёт адаптер
func New(cfg Config) *Adapter {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.partsapi.ru" // ← замени после теста в кабинете
	}
	if cfg.APIKeyVIN == "" {
		cfg.APIKeyVIN = cfg.APIKey
	}
	if cfg.APIKeyTree == "" {
		cfg.APIKeyTree = cfg.APIKey
	}

	return &Adapter{
		client:     &http.Client{Timeout: cfg.Timeout},
		baseURL:    cfg.BaseURL,
		apiKey:     cfg.APIKey,
		vinAPIKey:  cfg.APIKeyVIN,
		treeAPIKey: cfg.APIKeyTree,
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

type SearchTreeNode struct {
	ID    string `json:"id"`
	Text  string `json:"text"`
	Path  string `json:"path"`
	Level int    `json:"level"`
}

var vinRegexp = regexp.MustCompile(`^[A-HJ-NPR-Z0-9]{17}$`)

// GetPartsByVIN возвращает сырые данные из getPartsbyVIN
func (a *Adapter) GetPartsByVIN(ctx context.Context, vin, cat string, isOrigPart bool) ([]Part, error) {
	vin, err := normalizeAndValidateVIN(vin)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("method", "getPartsbyVIN")
	params.Add("key", a.apiKey)
	params.Add("vin", vin)
	params.Add("lang", "ru")
	if isOrigPart {
		params.Add("type", "oem")
	} else {
		params.Add("type", "all") // или ""
	}
	if cat != "" {
		params.Add("cat", cat)
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

// GetSearchTreeByVIN возвращает упрощенный плоский список узлов из getSearchTree.
func (a *Adapter) GetSearchTreeByVIN(ctx context.Context, vin string) ([]SearchTreeNode, error) {
	vin, err := normalizeAndValidateVIN(vin)
	if err != nil {
		return nil, err
	}
	carID, err := a.resolveCarIDByVIN(ctx, vin)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("method", "getSearchTree")
	params.Add("key", a.treeAPIKey)
	params.Add("lang", "16")
	params.Add("carId", carID)
	params.Add("carType", "PC")

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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	rawRows, err := decodeSearchTreeRows(body)
	if err != nil {
		return nil, err
	}

	nodes := flattenTreeRows(rawRows)
	if len(nodes) == 0 {
		return nil, fmt.Errorf("дерево пустое для VIN")
	}
	return nodes, nil
}

func (a *Adapter) resolveCarIDByVIN(ctx context.Context, vin string) (string, error) {
	keys := uniqueNonEmpty([]string{a.vinAPIKey, a.treeAPIKey, a.apiKey})
	var lastErr error
	for _, key := range keys {
		carID, err := a.resolveCarIDByVINWithKey(ctx, vin, key)
		if err == nil {
			return carID, nil
		}
		lastErr = err
	}
	if lastErr != nil {
		return "", lastErr
	}
	return "", fmt.Errorf("VINdecode: не задан API ключ")
}

func (a *Adapter) resolveCarIDByVINWithKey(ctx context.Context, vin, key string) (string, error) {
	params := url.Values{}
	params.Add("method", "VINdecode")
	params.Add("key", key)
	params.Add("vin", vin)
	params.Add("lang", "ru")

	fullURL := a.baseURL + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("VINdecode read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("VINdecode ошибка %d: %s", resp.StatusCode, truncateForLog(body, 400))
	}

	var obj map[string]any
	if err := json.Unmarshal(body, &obj); err != nil {
		return "", fmt.Errorf("VINdecode JSON decode failed: %w; body: %s", err, truncateForLog(body, 400))
	}

	if msg := extractErrorMessage(obj); msg != "" {
		return "", fmt.Errorf("VINdecode error: %s", msg)
	}

	carID := strings.TrimSpace(valueAsString(obj["carId"]))
	if carID == "" {
		carID = strings.TrimSpace(valueAsString(obj["TecDocExternalId"]))
	}
	if carID == "" {
		return "", fmt.Errorf("VINdecode не вернул carId: %s", truncateForLog(body, 400))
	}

	return carID, nil
}

func uniqueNonEmpty(values []string) []string {
	out := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, v := range values {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}

func decodeSearchTreeRows(body []byte) ([]map[string]any, error) {
	var rows []map[string]any
	if err := json.Unmarshal(body, &rows); err == nil {
		return rows, nil
	}

	var obj map[string]any
	if err := json.Unmarshal(body, &obj); err != nil {
		return nil, fmt.Errorf("JSON decode failed: %w; body: %s", err, truncateForLog(body, 400))
	}

	if msg := extractErrorMessage(obj); msg != "" {
		return nil, fmt.Errorf("partsapi getSearchTree error: %s", msg)
	}

	if dataRows := anyToRows(obj["data"]); len(dataRows) > 0 {
		return dataRows, nil
	}
	if dataRows := anyToRows(obj["items"]); len(dataRows) > 0 {
		return dataRows, nil
	}

	return nil, fmt.Errorf("неожиданный формат getSearchTree: %s", truncateForLog(body, 400))
}

func anyToRows(v any) []map[string]any {
	list, ok := v.([]any)
	if !ok || len(list) == 0 {
		return nil
	}
	rows := make([]map[string]any, 0, len(list))
	for _, item := range list {
		row, ok := item.(map[string]any)
		if !ok {
			continue
		}
		rows = append(rows, row)
	}
	return rows
}

func extractErrorMessage(obj map[string]any) string {
	keys := []string{"error", "message", "msg", "status", "description"}
	for _, k := range keys {
		if s := strings.TrimSpace(valueAsString(obj[k])); s != "" {
			low := strings.ToLower(s)
			if low == "ok" || low == "success" {
				continue
			}
			return s
		}
	}
	return ""
}

func truncateForLog(body []byte, max int) string {
	s := strings.TrimSpace(string(body))
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}

func normalizeAndValidateVIN(vin string) (string, error) {
	normalized := strings.ToUpper(strings.TrimSpace(vin))
	if len(normalized) != 17 {
		return "", fmt.Errorf("некорректный VIN: должно быть ровно 17 символов")
	}
	if !vinRegexp.MatchString(normalized) {
		return "", fmt.Errorf("некорректный VIN: допустимы латинские буквы и цифры, без I/O/Q")
	}
	return normalized, nil
}

func flattenTreeRows(rows []map[string]any) []SearchTreeNode {
	result := make([]SearchTreeNode, 0, len(rows)*3)
	seen := make(map[string]struct{})

	for _, row := range rows {
		parts := make([]string, 0, 4)
		rootText := valueAsString(row["ROOT_NODE_TEXT"])
		if rootText != "" {
			parts = append(parts, rootText)
		}

		for level := 1; level <= 5; level++ {
			textKey := "NODE_" + strconv.Itoa(level) + "_TEXT"
			idKey := "NODE_" + strconv.Itoa(level) + "_STR_ID"
			text := strings.TrimSpace(valueAsString(row[textKey]))
			id := strings.TrimSpace(valueAsString(row[idKey]))
			if text == "" {
				continue
			}

			parts = append(parts, text)
			if id == "" {
				continue
			}

			path := strings.Join(parts, " / ")
			if _, exists := seen[id]; exists {
				continue
			}
			seen[id] = struct{}{}

			result = append(result, SearchTreeNode{
				ID:    id,
				Text:  text,
				Path:  path,
				Level: level,
			})
		}
	}

	return result
}

func valueAsString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case float64:
		return strconv.FormatInt(int64(t), 10)
	case int:
		return strconv.Itoa(t)
	case int64:
		return strconv.FormatInt(t, 10)
	case json.Number:
		return t.String()
	default:
		return fmt.Sprintf("%v", t)
	}
}
