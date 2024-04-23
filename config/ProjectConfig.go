package projectConfig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Config struct {
	DirectusURL           string           `yaml:"directusURL"`
	Email                 string           `yaml:"email"`
	Password              string           `yaml:"password"`
	LoginRustat           string           `yaml:"loginRustat"`
	PasswordRustat        string           `yaml:"passwordRustat"`
	TeamIDFromRustat      string           `yaml:"teamIDFromRustat"`
	TeamRelations         map[string]int64 `yaml:"teamRelations"`
	PlayerRelations       map[string]int64 `yaml:"playerRelations"`
	Token                 string
	SeasonID              string
	LeagueID              string
	SeasonIDForRustat     string `yaml:"seasonIDForRustat"`
	TournamentIDForRustat string `yaml:"tournamentIDForRustat"`
}

func InitConfig() Config {
	yamlFile, readCfgFileErr := ioutil.ReadFile("config/config.yml")
	if readCfgFileErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Ошибка чтения файла config.yml - file: %s | line: %d | error: %v", file, line-2, readCfgFileErr)
	}

	config := Config{}

	unmarshalErr := yaml.Unmarshal(yamlFile, &config)
	if unmarshalErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Ошибка разбора YAML - file: %s | line: %d | error: %v", file, line-1, unmarshalErr)
	}

	if loginErr := config.LoginToDirectus(); loginErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Ошибка авторизации и получения токена - file: %s | line: %d | error: %v", file, line-1, loginErr)
	}

	if takeGlobalSettingsErr := config.TakeSeasonAndLeagueIDsFromGlobalSettings(); takeGlobalSettingsErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Ошибка получения ID сезона и лиги из глобальных настроек Directus - file: %s | line: %d | error: %v", file, line-1, takeGlobalSettingsErr)
	}

	if config.Token == "" {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Токен не был получен, дальнейшее выполнение без него не имеет смысла - file: %s | line: %d | error: %v", file, line-1)
	}

	return config
}

func (c *Config) TakeSeasonAndLeagueIDsFromGlobalSettings() error {
	query := `
       query{
			global_settings{
				season_id{
					id
				}
				league_id{
					id
				}
			}
		}
    `

	graphqlRequest := GraphQLRequest{
		Query: query,
	}

	jsonData, marshalErr := json.Marshal(graphqlRequest)
	if marshalErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Ошибка marshal тела запроса: file: %s | line: %d | error: %v", file, line-2, marshalErr)
	}

	reader := bytes.NewReader(jsonData)

	url := fmt.Sprintf("https://%s/graphql", c.DirectusURL)
	req, createReqErr := http.NewRequest("POST", url, reader)
	if createReqErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Ошибка создания запроса с url: %s | file: %s | line: %d | error: %v", url, file, line-2, createReqErr)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}

	resp, reqErr := client.Do(req)
	if reqErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Ошибка отправки запроса по url: %s | file: %s | line: %d | error: %v", url, file, line-2, reqErr)
	}

	var result map[string]interface{}
	if decoderErr := json.NewDecoder(resp.Body).Decode(&result); decoderErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Ошибка при decode тела ответа по url:  %s | file: %s | line: %d | error: %v", url, file, line-1, decoderErr)
	}

	if closeErr := resp.Body.Close(); closeErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Ошибка при попытке закрыть тело ответа по url: %s | file: %s | line: %d | error: %v", url, file, line-1, closeErr)
	}

	data, ok := result["data"]
	if !ok {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Не удалось получить ID сезона и лиги - file: %s | line: %d", file, line-2)
	}

	seasonID := data.(map[string]interface{})["global_settings"].(map[string]interface{})["season_id"].(map[string]interface{})["id"].(string)

	leagueID := data.(map[string]interface{})["global_settings"].(map[string]interface{})["league_id"].(map[string]interface{})["id"].(string)

	c.SeasonID = seasonID
	c.LeagueID = leagueID

	return nil
}

func (c *Config) LoginToDirectus() error {
	url := fmt.Sprintf("https://%s/auth/login", c.DirectusURL)
	user := User{
		Email:    c.Email,
		Password: c.Password,
	}

	jsonData, marshalErr := json.Marshal(user)
	if marshalErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Не удалось конвертировать сущность в json: file: %s | line: %d | error: %v", file, line-2, marshalErr)
	}

	reader := bytes.NewReader(jsonData)

	resp, createReqErr := http.Post(url, "application/json", reader)
	if createReqErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Ошибка при выполнении запроса: file: %s | line: %d | error: %v", file, line-2, createReqErr)
	}

	var result map[string]interface{}
	if decoderErr := json.NewDecoder(resp.Body).Decode(&result); decoderErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Ошибка декодирования ответа: file: %s | line: %d | error: %v", file, line-2, decoderErr)
	}

	if closeErr := resp.Body.Close(); closeErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Ошибка при закрытии тела ответа: file: %s | line: %d | error: %v", file, line-1, closeErr)
	}
	data, ok := result["data"]
	if !ok {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Не удалось получить токен доступа: file: %s | line: %d", file, line-2)
	}

	accessToken, ok := data.(map[string]interface{})["access_token"].(string)
	if !ok {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Токен доступа имеет неверный формат: file: %s | line: %d | error: %v", file, line-2)
	}

	c.Token = accessToken

	return nil
}
func (c *Config) PrintConfig() {
	log.Println("URL на который будут отправляться запросы - ", c.DirectusURL, "\n",
		"Email для авторизации -", c.Email, "\n",
		"Password для авторизации -", c.Password, "\n",
		"Токен авторизации -", c.Token)
}
