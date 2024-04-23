package team

import (
	projectConfig "ApiChernoToRustat/config"
	"fmt"
	"log"
	"runtime"
	"strconv"
)

type CompositionOfTeam struct {
	TeamID      int64
	TeamName    string
	Composition Lineup
}

type Lineup struct {
	FirstLineup  []PlayersForLineup
	SecondLineup []PlayersForLineup
}

type PlayersForLineup struct {
	PlayerID  int64
	Firstname string
	Lastname  string
	Number    string
}

func TakeCompositionOfTeamForMatch(matchID string, config projectConfig.Config) ([]CompositionOfTeam, error) {
	rustatApi := fmt.Sprintf("http://feeds.rustatsport.ru/?tpl=39&user=%s&key=%s&match_id=%s&lang_id=0&format=json", config.LoginRustat, config.PasswordRustat, matchID)

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

	rows, rowsOk := data.(map[string]interface{})
	if !rowsOk {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось получить map row из тела ответа - file: %s | line: %d", file, line-2)
	}

	firstTeamComposition := rows["first_team"].([]interface{})[0].(map[string]interface{})["lineup"].([]interface{})[0].(map[string]interface{})["main"].([]interface{})[0].(map[string]interface{})["player"].([]interface{})
	secondTeamComposition := rows["second_team"].([]interface{})[0].(map[string]interface{})["lineup"].([]interface{})[0].(map[string]interface{})["main"].([]interface{})[0].(map[string]interface{})["player"].([]interface{})

	firstTeamName := rows["first_team"].([]interface{})[0].(map[string]interface{})["name"].(string)
	secondTeamName := rows["second_team"].([]interface{})[0].(map[string]interface{})["name"].(string)

	//Создали составы
	firstTeam, CreateCompositionFirstTeamErr := createComposition(firstTeamComposition, firstTeamName)
	if CreateCompositionFirstTeamErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось создать состав команды - file: %s | line: %d", file, line-2)
	}
	secondTeam, CreateCompositionSecondTeamErr := createComposition(secondTeamComposition, secondTeamName)
	if CreateCompositionSecondTeamErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("Не удалось создать состав команды - file: %s | line: %d", file, line-2)
	}
	//Всем выставили ID
	setID(firstTeam, config)
	setID(secondTeam, config)

	var matchComposition []CompositionOfTeam
	matchComposition = append(matchComposition, *firstTeam, *secondTeam)

	return matchComposition, nil
}

// Создаем состав команды
func createComposition(lineup []interface{}, teamName string) (*CompositionOfTeam, error) {
	teamFirstLineup := []PlayersForLineup{}
	teamSecondLineup := []PlayersForLineup{}

	for _, params := range lineup {
		playerPrams, playerPramsOk := params.(map[string]interface{})
		if !playerPramsOk {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Не удалось преобразовать инофрмацию об игроке для состава команды в map - file: %s | line: %d", file, line-2)
		}
		player := PlayersForLineup{}
		startingLineUp := false
		for key, value := range playerPrams {
			switch key {
			case "firstname", "lastname", "starting_lineup", "num":
				switch key {
				case "firstname":
					player.Firstname = value.(string)
				case "lastname":
					player.Lastname = value.(string)
				case "starting_lineup":
					if value == "1" {
						startingLineUp = true
					}
				case "num":
					player.Number = value.(string)
				}
			}
		}
		if startingLineUp == true {
			teamFirstLineup = append(teamFirstLineup, player)
		} else {
			teamSecondLineup = append(teamSecondLineup, player)
		}
	}

	teamLineUp := Lineup{
		FirstLineup:  teamFirstLineup,
		SecondLineup: teamSecondLineup,
	}

	team := CompositionOfTeam{
		TeamName:    teamName,
		Composition: teamLineUp,
	}

	return &team, nil
}

func setID(team *CompositionOfTeam, config projectConfig.Config) error {
	for key, value := range config.TeamRelations {
		if team.TeamName == key {
			team.TeamID = value
		}
	}

	if team.TeamID == 2 {
		for index := range team.Composition.FirstLineup {
			player := &team.Composition.FirstLineup[index]
			for key, value := range config.PlayerRelations {
				if player.Firstname+" "+player.Lastname == key {
					player.PlayerID = value
				}
			}
		}
		for index := range team.Composition.SecondLineup {
			player := &team.Composition.SecondLineup[index]
			for key, value := range config.PlayerRelations {
				if player.Firstname+" "+player.Lastname == key {
					player.PlayerID = value
				}
			}
		}
	} else {
		log.Println(team.TeamName)
		if takePlayersIDErr := takePlayersID(&*team, config); takePlayersIDErr != nil {
			_, file, line, _ := runtime.Caller(0)
			return fmt.Errorf("Не удалось установить ID для игроков - file: %s | line: %d", file, line-2)
		}
	}

	return nil
}

func takePlayersID(team *CompositionOfTeam, config projectConfig.Config) error {
	query := `
      query(
		$team_name:String
	){
		team_members(filter:{team_id:{name:{_eq:$team_name}}}){
			id
			firstname
			lastname
		}
	}
    `
	variables := map[string]interface{}{
		"team_name": team.TeamName,
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

	data, dataOk := result["data"]
	if !dataOk {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Не удалось извлечь состав команды из тела ответа - file: %s | line: %d", file, line-2)
	}

	rows, rowsOk := data.(map[string]interface{})
	if !rowsOk {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Не удалось получить map row из тела ответа - file: %s | line: %d", file, line-2)
	}

	teamMembers := rows["team_members"].([]interface{})

	//Создаем из ответа мап ID - имя фамилия

	fullnameIDRealation, makefullnameIDRealationsErr := makefullnameIDRealations(teamMembers)
	if makefullnameIDRealationsErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Не удалось создать отношение ID - имя фамилия, соответственно не удастся продолжить дальнейшее выполнение - file: %s | line: %d", file, line-2)
	}

	log.Println(fullnameIDRealation)

	for index := range team.Composition.FirstLineup {
		playerLineup := &team.Composition.FirstLineup[index]
		for key, value := range fullnameIDRealation {
			if value == playerLineup.Firstname+" "+playerLineup.Lastname {
				var parseErr error
				playerLineup.PlayerID, parseErr = strconv.ParseInt(key, 10, 64)
				if parseErr != nil {
					_, file, line, _ := runtime.Caller(0)
					log.Println("Не удалось преобразовать string в int64, ID. - file: %s | line: %d", file, line-2)
					break
				}
			}
		}
	}
	for index := range team.Composition.SecondLineup {
		playerLineup := &team.Composition.SecondLineup[index]
		for key, value := range fullnameIDRealation {
			if value == playerLineup.Firstname+" "+playerLineup.Lastname {
				var parseErr error
				playerLineup.PlayerID, parseErr = strconv.ParseInt(key, 10, 64)
				if parseErr != nil {
					_, file, line, _ := runtime.Caller(0)
					log.Println("Не удалось преобразовать string в int64, ID. - file: %s | line: %d", file, line-2)
					break
				}
			}
		}
	}

	return nil
}

func makefullnameIDRealations(teamMembers []interface{}) (map[string]string, error) {
	result := make(map[string]string)

	for _, teamMember := range teamMembers {
		player, playerOk := teamMember.(map[string]interface{})
		if !playerOk {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Не удалось преобразовать игрока в map - file: %s | line: %d", file, line-2)
		}

		var id string
		var firstname string
		var lastname string

		for key, value := range player {
			switch key {
			case "id":
				id = value.(string)
			case "firstname":
				firstname = value.(string)
			case "lastname":
				lastname = value.(string)
			}
		}

		result[id] = firstname + " " + lastname
	}

	return result, nil

}
