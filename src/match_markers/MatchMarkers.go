package match_markers

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"log"
)

type MatchMarkers struct {
	ActionID                 string `json:"action_id"`
	ActionName               string `json:"action_name"`
	AttackFlangID            string `json:"attack_flang_id"`
	AttackFlangName          string `json:"attack_flang_name"`
	AttackNumber             string `json:"attack_number"`
	AttackStatusID           string `json:"attack_status_id"`
	AttackStatusName         string `json:"attack_status_name"`
	AttackTeamExternalID     string `json:"attack_team_external_id"`
	AttackTeamID             string `json:"attack_team_id"`
	AttackTeamName           string `json:"attack_team_name"`
	AttackTypeID             string `json:"attack_type_id"`
	AttackTypeName           string `json:"attack_type_name"`
	DL                       string `json:"dl"`
	Half                     string `json:"half"`
	ID                       string `json:"id"`
	PlayerExternalID         string `json:"player_external_id"`
	PlayerID                 string `json:"player_id"`
	PlayerName               string `json:"player_name"`
	PosX                     string `json:"pos_x"`
	PosY                     string `json:"pos_y"`
	PositionID               string `json:"position_id"`
	PositionName             string `json:"position_name"`
	PossessionID             string `json:"possession_id"`
	PossessionName           string `json:"possession_name"`
	PossessionNumber         string `json:"possession_number"`
	PossessionTeamExternalID string `json:"possession_team_external_id"`
	PossessionTeamID         string `json:"possession_team_id"`
	PossessionTeamName       string `json:"possession_team_name"`
	PossessionTime           string `json:"possession_time"`
	Second                   string `json:"second"`
	StandartID               string `json:"standart_id"`
	StandartName             string `json:"standart_name"`
	TeamExternalID           string `json:"team_external_id"`
	TeamID                   string `json:"team_id"`
	TeamName                 string `json:"team_name"`
	TS                       string `json:"ts"`
	ZoneID                   string `json:"zone_id"`
	ZoneName                 string `json:"zone_name"`
}

func (t *MatchMarkers) TakeMatchMarkersFromRustat() []MatchMarkers {
	var matchMarkers []MatchMarkers

	rustatAPI := "http://feeds.rustatsport.ru/link"

	client := resty.New()
	resp, err := client.R().Get(rustatAPI)
	if err != nil {
		log.Println("Ошибка отправки запроса", err)
		return matchMarkers
	}

	var jsonResponse map[string]interface{}
	err = json.Unmarshal(resp.Body(), &jsonResponse)
	if err != nil {
		log.Println("Не удалось превратить ответ в структуру", err)
		return matchMarkers
	}

	rows, ok := jsonResponse["data"].(map[string]interface{})["row"].([]interface{})
	if !ok {
		log.Println("Не удалось получить массив row")
		return matchMarkers
	}

	for _, row := range rows {
		rowMap, ok := row.(map[string]interface{})
		if !ok {
			log.Println("Не удалось преобразовать элемент row в map[string]interface{}")
			continue
		}

		match_markers := MatchMarkers{
			ActionID:                 getStringValue(rowMap["action_id"]),
			ActionName:               getStringValue(rowMap["action_name"]),
			AttackFlangID:            getStringValue(rowMap["attack_flang_id"]),
			AttackFlangName:          getStringValue(rowMap["attack_flang_name"]),
			AttackNumber:             getStringValue(rowMap["attack_number"]),
			AttackStatusID:           getStringValue(rowMap["attack_status_id"]),
			AttackStatusName:         getStringValue(rowMap["attack_status_name"]),
			AttackTeamExternalID:     getStringValue(rowMap["attack_team_external_id"]),
			AttackTeamID:             getStringValue(rowMap["attack_team_id"]),
			AttackTeamName:           getStringValue(rowMap["attack_team_name"]),
			AttackTypeID:             getStringValue(rowMap["attack_type_id"]),
			AttackTypeName:           getStringValue(rowMap["attack_type_name"]),
			DL:                       getStringValue(rowMap["dl"]),
			Half:                     getStringValue(rowMap["half"]),
			ID:                       getStringValue(rowMap["id"]),
			PlayerExternalID:         getStringValue(rowMap["player_external_id"]),
			PlayerID:                 getStringValue(rowMap["player_id"]),
			PlayerName:               getStringValue(rowMap["player_name"]),
			PosX:                     getStringValue(rowMap["pos_x"]),
			PosY:                     getStringValue(rowMap["pos_y"]),
			PositionID:               getStringValue(rowMap["position_id"]),
			PositionName:             getStringValue(rowMap["position_name"]),
			PossessionID:             getStringValue(rowMap["possession_id"]),
			PossessionName:           getStringValue(rowMap["possession_name"]),
			PossessionNumber:         getStringValue(rowMap["possession_number"]),
			PossessionTeamExternalID: getStringValue(rowMap["possession_team_external_id"]),
			PossessionTeamID:         getStringValue(rowMap["possession_team_id"]),
			PossessionTeamName:       getStringValue(rowMap["possession_team_name"]),
			PossessionTime:           getStringValue(rowMap["possession_time"]),
			Second:                   getStringValue(rowMap["second"]),
			StandartID:               getStringValue(rowMap["standart_id"]),
			StandartName:             getStringValue(rowMap["standart_name"]),
			TeamExternalID:           getStringValue(rowMap["team_external_id"]),
			TeamID:                   getStringValue(rowMap["team_id"]),
			TeamName:                 getStringValue(rowMap["team_name"]),
			TS:                       getStringValue(rowMap["ts"]),
			ZoneID:                   getStringValue(rowMap["zone_id"]),
			ZoneName:                 getStringValue(rowMap["zone_name"]),
		}

		matchMarkers = append(matchMarkers, match_markers)
	}

	return matchMarkers
}

func getStringValue(value interface{}) string {
	if value == nil {
		return ""
	}
	stringValue, ok := value.(string)
	if !ok {
		return ""
	}
	return stringValue
}
