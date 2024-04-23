package team

import (
	projectConfig "ApiChernoToRustat/config"
	"fmt"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"text/tabwriter"
)

type Team struct {
	ID           string
	Name         string
	PhotoURL     string
	PhotoPath    string
	BaseYear     string
	WebSiteURL   string
	SeasonID     string
	TeamTypeID   string
	TeamTypeName string
	CountryID    string
	CountryName  string
	GenderID     string
	GenderName   string
	TS           string
	Index        *string
}

type Player struct {
	ID                     string
	ExternalID             string
	Firstname              string
	Lastname               string
	ClubID                 string
	ClubTeamExternalID     string
	ClubName               string
	IsNationalTeam         string
	NationalTeamID         *string
	NationalTeamExternalID *string
	NationalTeamName       *string
	Position1Name          string
	Position1ID            string
	Position2Name          *string
	Position2ID            *string
	Position3Name          *string
	Position3ID            *string
	FootName               string
	FootID                 string
	ClubNumber             string
	NationalNumber         *string
	CountryID              string
	CountryName            string
	TS                     string
	Birthdate              *string
	ContractEnding         *string
	GenderID               *string
	GenderName             *string
	SeasonID               *string
	Nickname               *string
	Weight                 *string
	Height                 *string
}

type TournamentTableTeam struct {
	TeamID         int64
	TeamName       string
	Position       int16
	Games          int16
	Winnings       int16
	Draws          int16
	Losses         int16
	Goals          int16
	GoalsConceded  int16
	DifferentGoals int16
	Points         int16
}

func (t *Team) TakeTeamFromRustat(config projectConfig.Config) (*Team, error) {
	team := Team{}
	rustatApi := fmt.Sprintf("http://feeds.rustatsport.ru/?tpl=12&user=%s&key=%s&team_id=%s&lang_id=1&format=json", config.LoginRustat, config.PasswordRustat, config.TeamIDFromRustat)

	graphQLObject := GraphQLRequest{}
	result, sendRequestErr := graphQLObject.SendGetRequest(rustatApi)
	if sendRequestErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось отправить graphQL запрос на url: %s | file: %s | line: %d | error: %v", rustatApi, file, line-2, sendRequestErr)
	}

	data, dataOk := result["data"]
	if !dataOk {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось извлечь команду из тела ответа - file: %s | line: %d", file, line-2)
	}

	rows, rowsOk := data.(map[string]interface{})["row"].([]interface{})[0].(map[string]interface{})
	if !rowsOk {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось получить map row из тела ответа - file: %s | line: %d", file, line-2)
	}

	team = Team{
		ID:           rows["id"].(string),
		Name:         rows["name"].(string),
		PhotoURL:     rows["photo"].(string),
		SeasonID:     rows["season_id"].(string),
		BaseYear:     rows["base_year"].(string),
		WebSiteURL:   rows["web_site"].(string),
		TeamTypeID:   rows["team_type_id"].(string),
		TeamTypeName: rows["team_type_name"].(string),
		CountryID:    rows["country_id"].(string),
		CountryName:  rows["country_name"].(string),
		GenderID:     rows["gender_id"].(string),
		GenderName:   rows["gender_name"].(string),
		TS:           rows["ts"].(string),
	}
	if rows["index"] != nil {
		team.Index = rows["index"].(*string)
	} else {
		team.Index = nil
	}
	return &team, nil
}

func (t *Team) TakeStructureOfTeam(team Team, config projectConfig.Config) (*[]Player, error) {
	rustatApi := fmt.Sprintf("http://feeds.rustatsport.ru/?tpl=5&user=%s&key=%s&team_id=%s&season_id=%s&lang_id=1&format=json", config.LoginRustat, config.PasswordRustat, team.ID, team.SeasonID)

	graphQLObject := GraphQLRequest{}
	result, sendRequestErr := graphQLObject.SendGetRequest(rustatApi)
	if sendRequestErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось отправить graphQL запрос на url: %s | file: %s | line: %d | error: %v", rustatApi, file, line-2, sendRequestErr)
	}

	data, dataOk := result["data"]
	if !dataOk {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось извлечь состав команды из тела ответа - file: %s | line: %d", file, line-2)
	}

	rows, jsonMapOk := data.(map[string]interface{})["row"].([]interface{})
	if !jsonMapOk {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось получить map row из тела ответа - file: %s | line: %d", file, line-2)
	}

	var players []map[string]interface{}
	for _, value := range rows {
		player, valueOk := value.(map[string]interface{})
		if !valueOk {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Не удалось преобразовать статистику игрока в map - file: %s | line: %d", file, line-2)
		}
		players = append(players, player)
	}

	var resultPlayers []Player

	for _, value := range players {
		player := Player{
			ID:             value["id"].(string),
			Firstname:      value["firstname"].(string),
			Lastname:       value["lastname"].(string),
			ClubID:         value["club_team_id"].(string),
			IsNationalTeam: value["is_national_team"].(string),
			Position1Name:  value["position1_name"].(string),
			Position1ID:    value["position1_id"].(string),
			ClubNumber:     value["club_number"].(string),
			CountryID:      value["country1_id"].(string),
			CountryName:    value["country1_name"].(string),
			TS:             value["ts"].(string),
		}
		if value["external_id"] != nil {
			player.ExternalID = value["external_id"].(string)
		}
		if value["club_team_external_id"] != nil {
			player.ClubTeamExternalID = value["club_team_external_id"].(string)
		}
		if value["club_team_name"] != nil {
			player.ClubName = value["club_team_name"].(string)
		}
		if value["foot_name"] != nil {
			player.FootName = value["foot_name"].(string)
		}
		if value["foot_id"] != nil {
			player.FootID = value["foot_id"].(string)
		}
		if value["national_team_id"] != nil {
			val := value["national_team_id"].(string)
			player.NationalTeamID = &val
		}
		if value["national_team_external_id"] != nil {
			val := value["national_team_external_id"].(string)
			player.NationalTeamExternalID = &val
		}
		if value["national_team_name"] != nil {
			val := value["national_team_name"].(string)
			player.NationalTeamName = &val
		}
		if value["position2_id"] != nil {
			val := value["position2_id"].(string)
			player.Position2ID = &val
		}
		if value["position2_name"] != nil {
			val := value["position2_name"].(string)
			player.Position2Name = &val
		}
		if value["position3_id"] != nil {
			val := value["position3_id"].(string)
			player.Position3ID = &val
		}
		if value["position3_name"] != nil {
			val := value["position3_name"].(string)
			player.Position3Name = &val
		}
		if value["national_number"] != nil {
			val := value["national_number"].(string)
			player.NationalNumber = &val
		}

		resultPlayers = append(resultPlayers, player)
	}

	return &resultPlayers, nil

}

func (t TournamentTableTeam) TakeTournamentTable(config projectConfig.Config) ([]TournamentTableTeam, error) {
	tournamentTable := []TournamentTableTeam{}
	rustatApi := fmt.Sprintf("http://feeds.rustatsport.ru/?tpl=46&user=%s&key=%s&tournament_id=%s&season_id=%s&lang_id=1&format=json", config.LoginRustat, config.PasswordRustat, config.TournamentIDForRustat, config.SeasonIDForRustat)

	graphQLObject := GraphQLRequest{}
	result, sendRequestErr := graphQLObject.SendGetRequest(rustatApi)
	if sendRequestErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось отправить graphQL запрос на url: %s | file: %s | line: %d | error: %v", rustatApi, file, line-2, sendRequestErr)
	}

	data, dataOk := result["data"]
	if !dataOk {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось извлечь турнирную таблицу из тела ответа - file: %s | line: %d", file, line-2)
	}

	rows, jsonMapOk := data.(map[string]interface{})["row"].([]interface{})
	if !jsonMapOk {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось получить map row из тела ответа - file: %s | line: %d", file, line-2)
	}

	var params []map[string]interface{}
	for _, value := range rows {
		param, paramOk := value.(map[string]interface{})
		if !paramOk {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Не удалось преобразовать статистику команды в турнирной таблице в map - file: %s | line: %d", file, line-2)
		}

		params = append(params, param)
	}

	for _, param := range params {
		team := TournamentTableTeam{
			Position: 999,
		}

		for key, value := range param {
			switch key {
			// Ищем совпадение названий команды с ответа и с конфига, если совпало даем id команде, под которым она находится в директусе
			case "team_name":
				team.TeamName = value.(string)
				for teamName, teamID := range config.TeamRelations {
					if value == teamName {
						team.TeamID = teamID
					}
				}
			case "position", "total", "won", "draw", "lost", "goals_for", "goals_against", "goals_diff", "points":
				valueStr, valueStrOk := value.(string)
				if valueStrOk {
					if valueFloat, parseErr := strconv.ParseFloat(valueStr, 32); parseErr == nil {
						switch key {
						case "position":
							team.Position = int16(valueFloat)
						case "total":
							team.Games = int16(valueFloat)
						case "won":
							team.Winnings = int16(valueFloat)
						case "draw":
							team.Draws = int16(valueFloat)
						case "lost":
							team.Losses = int16(valueFloat)
						case "goals_for":
							team.Goals = int16(valueFloat)
						case "goals_against":
							team.GoalsConceded = int16(valueFloat)
						case "goals_diff":
							team.DifferentGoals = int16(valueFloat)
						case "points":
							team.Points = int16(valueFloat)
						}
					} else {
						_, file, line, _ := runtime.Caller(0)
						log.Printf("Не удалось преобразовать в int значение: file: %s | line: %d | error: %v", file, line-22, parseErr)
					}
				}
			}
		}
		tournamentTable = append(tournamentTable, team)
	}
	return tournamentTable, nil
}

// Берем id записей стат команд из нужной нам турнирной таблицы
func TakeTeamRecordIDsFromRatingTable(config projectConfig.Config) (*[]string, error) {
	query := `
       query da(
			$season_id:GraphQLStringOrFloat!,
			$league_id:GraphQLStringOrFloat!,
		){
					ratings(filter:{
						season_id: { id: { _eq: $season_id}},
						league_id: { id: { _eq: $league_id}},
					}){
						season_id{
							id
						}
						league_id{
							id
						}
						teams{
							id
						}
					}
		}
    `
	variables := map[string]interface{}{
		"season_id": config.SeasonID,
		"league_id": config.LeagueID,
	}
	graphqlRequest := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	result, sendRequestErr := graphqlRequest.SendPostRequest(graphqlRequest, config)
	if sendRequestErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось отправить graphQL запрос: file: %s | line: %d | error: %v", file, line-2, sendRequestErr)
	}

	ids := []string{}
	ratingsOfTeam, ratingsOfTeamOk := result["data"].(map[string]interface{})["ratings"].([]interface{})[0].(map[string]interface{})["teams"].([]interface{})
	if !ratingsOfTeamOk {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось получить ID сезона и лиги - file: %s | line: %d", file, line-2)
	}

	for _, value := range ratingsOfTeam {
		id, idOk := value.(map[string]interface{})["id"].(string)
		if !idOk {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Не удалось преобразовать ответ в ID - file: %s | line: %d", file, line-2)
		}
		ids = append(ids, id)
	}
	return &ids, nil
}

// Делаем мап связей между записью и командой
func TakeRelationsBetweenRecordAndTeamStat(recordIDs []string, config projectConfig.Config) (*map[int]int64, error) {
	query := `
       query($ids: [GraphQLStringOrFloat]!){
			teams_stats(
				filter:{id:{_in: $ids}}
			){
				id
				team_id{
					id
					name
				}
				games
				winnings
				losses
				points
			}
		}
   `
	variables := map[string]interface{}{
		"ids": recordIDs,
	}
	graphqlRequest := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	result, sendRequestErr := graphqlRequest.SendPostRequest(graphqlRequest, config)
	if sendRequestErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось отправить graphQL запрос: file: %s | line: %d | error: %v", file, line-2, sendRequestErr)
	}

	//Начинаем создавать связи команда - ее запись
	relations := map[int]int64{}
	teamStats, teamStatsOk := result["data"].(map[string]interface{})["teams_stats"].([]interface{})
	if !teamStatsOk {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось достать записи с рейтингами команд - file: %s | line: %d", file, line-2)
	}

	for _, value := range teamStats {
		keyStr, keyStrOk := value.(map[string]interface{})["id"].(string)
		if !keyStrOk {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Не удалось преобразовать ID в string - file: %s | line: %d", file, line-2)
		}
		itemStr, itemStrOk := value.(map[string]interface{})["team_id"].(map[string]interface{})["id"].(string)
		if !itemStrOk {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Не удалось преобразовать teamID в string - file: %s | line: %d", file, line-2)
		}
		key, atoiErr := strconv.Atoi(keyStr)
		if atoiErr != nil {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Ошибка при конвертации ключа в int: file: %s | line: %d | error: %v", file, line-2, atoiErr)
		}
		item, parseErr := strconv.ParseInt(itemStr, 10, 64)
		if parseErr != nil {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Ошибка при конвертации ключа в int64: file: %s | line: %d | error: %v", file, line-2, parseErr)
		}

		relations[key] = item
	}
	return &relations, nil
}

func SendTable(table []TournamentTableTeam, relations map[int]int64, config projectConfig.Config) error {
	for _, team := range table {
		id := 0
		for key, value := range relations {
			if value == team.TeamID {
				id = key
			}
		}
		query := `
		   mutation($id: ID!, $data: update_teams_stats_input!){
				update_teams_stats_item(id: $id, data: $data) {
					id
					games
					winnings
					draws
					losses
					goals_scored
					conceded_goals
					difference_g_c
					points
				}
			}
   		`

		variables := map[string]interface{}{
			"id": id,
			"data": map[string]interface{}{
				"games":          team.Games,
				"position":       team.Position,
				"winnings":       team.Winnings,
				"draws":          team.Draws,
				"losses":         team.Losses,
				"goals_scored":   team.Goals,
				"conceded_goals": team.GoalsConceded,
				"difference_g_c": team.DifferentGoals,
				"points":         team.Points,
			},
		}

		graphqlRequest := GraphQLRequest{
			Query:     query,
			Variables: variables,
		}

		result, sendRequestErr := graphqlRequest.SendPostRequest(graphqlRequest, config)
		if sendRequestErr != nil {
			_, file, line, _ := runtime.Caller(0)
			return fmt.Errorf("Не удалось отправить graphQL запрос: file: %s | line: %d | error: %v", file, line-2, sendRequestErr)
		}

		log.Println("Ответ от сервера Directus:", result)
	}
	return nil
}

type ByPosition []TournamentTableTeam

func (a ByPosition) Len() int           { return len(a) }
func (a ByPosition) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPosition) Less(i, j int) bool { return a[i].Position < a[j].Position }

func PrintTournamentTable(teams []TournamentTableTeam) error {
	sort.Sort(ByPosition(teams))

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	log.SetOutput(w)

	// Header
	header := "TeamID\tTeamName\tPosition\tGames\tWinnings\tDraws\tLosses\tGoals\tGoalsConceded\tDifferentGoals\tPoints"
	log.Println(header)

	// Data rows
	for _, team := range teams {
		row := ""
		row += fmt.Sprintf("%v\t", team.TeamID)
		row += fmt.Sprintf("%v\t", team.TeamName)
		row += fmt.Sprintf("%v\t", team.Position)
		row += fmt.Sprintf("%v\t", team.Games)
		row += fmt.Sprintf("%v\t", team.Winnings)
		row += fmt.Sprintf("%v\t", team.Draws)
		row += fmt.Sprintf("%v\t", team.Losses)
		row += fmt.Sprintf("%v\t", team.Goals)
		row += fmt.Sprintf("%v\t", team.GoalsConceded)
		row += fmt.Sprintf("%v\t", team.DifferentGoals)
		row += fmt.Sprintf("%v\t", team.Points)
		log.Println(row)
	}

	if flushErr := w.Flush(); flushErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Ошибка очищения буфера: file: %s | line: %d | error: %v", file, line-1, flushErr)
	}

	return nil
}
