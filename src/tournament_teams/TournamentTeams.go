package tournament_teams

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
)

type TouranmentTeams struct {
	ID        string
	Name      string
	ShortName string
}

func (t *TouranmentTeams) TakeTouranmentTeamsFromRustat() []TouranmentTeams {
	var tournamentTeams []TouranmentTeams
	rustatApi := "http://feeds.rustatsport.ru/link"

	client := resty.New()
	resp, err := client.R().Get(rustatApi)
	if err != nil {
		log.Println("Ошибка отправки запроса", err)
		return tournamentTeams
	}

	var jsonMap map[string]interface{}
	err = json.Unmarshal(resp.Body(), &jsonMap)
	if err != nil {
		log.Println("Не удалось превратить ответ в структуру", err)
		return tournamentTeams
	}

	rows, ok := jsonMap["data"].(map[string]interface{})["row"].([]interface{})
	if !ok {
		log.Fatal("Не удалось получить массив row", err)
		return tournamentTeams
	}

	for _, row := range rows {
		// Преобразуем row в map[string]interface{}
		rowMap, ok := row.(map[string]interface{})
		if !ok {
			log.Println("Не удалось преобразовать элемент row в map[string]interface{}")
			continue
		}

		// Извлекаем значения из rowMap
		id, _ := rowMap["id"].(string)
		name, _ := rowMap["name"].(string)

		// Проверяем наличие поля short_name
		shortName, shortNameExists := rowMap["short_name"].(string)

		// Создаем экземпляр TouranmentTeams и добавляем его в список
		tournamentTeam := TouranmentTeams{
			ID:   id,
			Name: name,
		}

		// Если short_name существует, присваиваем его, иначе используем name
		if shortNameExists {
			tournamentTeam.ShortName = shortName
		} else {
			tournamentTeam.ShortName = name
		}

		tournamentTeams = append(tournamentTeams, tournamentTeam)
	}

	return tournamentTeams
}
