package calendar

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
)

type Calendar struct {
	ID             string
	MatchDate      string
	TournamentID   string
	TournamentName string
	TournamentRank string
	SeasonID       string
	SeasonName     string
	RoundID        string
	RoundName      string
	Team1ID        string
	Team1Name      string
	Team2ID        string
	Team2Name      string
	Team1Score     string
	Team2Score     string
	StatusID       string
	StatusName     string
	StadiumID      string
	StadiumName    string
	Duration       string
	MatchName      string
	GroupID        string
}

func (t *Calendar) TakeCalendarFromRustat() []Calendar {
	var calendar []Calendar
	rustatApi := "http://feeds.rustatsport.ru/link"

	client := resty.New()
	resp, err := client.R().Get(rustatApi)
	if err != nil {
		log.Println("Ошибка отправки запроса", err)
		return calendar
	}

	var jsonData map[string]interface{}
	err = json.Unmarshal(resp.Body(), &jsonData)
	if err != nil {
		log.Println("Не удалось превратить ответ в структуру", err)
		return calendar
	}
	data, dataExists := jsonData["data"].(map[string]interface{})
	if !dataExists {
		log.Println("Отсутствует поле 'data' в JSON-ответе")
		return calendar
	}

	var params []map[string]interface{}
	for _, value := range jsonData["data"].(map[string]interface{})["row"].([]interface{}) {
		param, ok := value.(map[string]interface{})
		if !ok {
			log.Fatal("Не удалось получить массив row", err)
		}
		params = append(params, param)
	}
	for _, value := range params {
		log.Println(value)
	}

	rows, rowsExist := data["row"].([]interface{})
	if !rowsExist {
		log.Println("Отсутствует поле 'row' в 'data' JSON-ответе")
		return calendar
	}

	for _, row := range rows {
		matchData := Calendar{
			ID:             interfaceToString(row.(map[string]interface{})["id"]),
			MatchDate:      interfaceToString(row.(map[string]interface{})["match_date"]),
			TournamentID:   interfaceToString(row.(map[string]interface{})["tournament_id"]),
			TournamentName: interfaceToString(row.(map[string]interface{})["tournament_name"]),
			TournamentRank: interfaceToString(row.(map[string]interface{})["tournament_rank"]),
			SeasonID:       interfaceToString(row.(map[string]interface{})["season_id"]),
			SeasonName:     interfaceToString(row.(map[string]interface{})["season_name"]),
			RoundID:        interfaceToString(row.(map[string]interface{})["round_id"]),
			RoundName:      interfaceToString(row.(map[string]interface{})["round_name"]),
			Team1ID:        interfaceToString(row.(map[string]interface{})["team1_id"]),
			Team1Name:      interfaceToString(row.(map[string]interface{})["team1_name"]),
			Team2ID:        interfaceToString(row.(map[string]interface{})["team2_id"]),
			Team2Name:      interfaceToString(row.(map[string]interface{})["team2_name"]),
			Team1Score:     interfaceToString(row.(map[string]interface{})["team1_score"]),
			Team2Score:     interfaceToString(row.(map[string]interface{})["team2_score"]),
			StatusID:       interfaceToString(row.(map[string]interface{})["status_id"]),
			StatusName:     interfaceToString(row.(map[string]interface{})["status_name"]),
			StadiumID:      interfaceToString(row.(map[string]interface{})["stadium_id"]),
			StadiumName:    interfaceToString(row.(map[string]interface{})["stadium_name"]),
			Duration:       interfaceToString(row.(map[string]interface{})["duration"]),
			MatchName:      interfaceToString(row.(map[string]interface{})["match_name"]),
			GroupID:        interfaceToString(row.(map[string]interface{})["group_id"]),
		}
		calendar = append(calendar, matchData)
	}

	return calendar
}

func interfaceToString(val interface{}) string {
	if val == nil {
		return "" // Возвращаем пустую строку, если значение nil
	}
	if str, ok := val.(string); ok {
		return str // Преобразовываем значение к строке, если это возможно
	}
	return "" // Возвращаем пустую строку, если значение не является строкой
}
