package team

import (
	projectConfig "ApiChernoToRustat/config"
	"fmt"
	"log"
	"runtime"
	"strconv"
)

type CommonPlayerStat struct {
	MatchesInTeam                   int16 //Matches played | Кол-во игр
	StartingLineupAppearances       int16 //Starting lineup appearances | в стартовом составе
	SubstitutesIn                   int16 //Substitutes in | вышел на замену
	SubstitutesOut                  int16 //Substitutes out |  был заменен
	Goals                           int16 //Goals | Голы
	Shots                           int16 //Shots | Удары
	ShotsOnTarget                   int16 //Shots on target | Удары в створ ворот
	Assists                         int16 //Assists |  голевые передачи
	ShotsOnTargetAccuratePercentage int16 //Shots on target, % | Точность ударов
	MinutesPlayed                   int16 //Minutes played | Кол-во минут на поле
}
type PassesPlayerStatistic struct {
	Passes                                                int16 // Passes | кол-во передач в целом
	KeyPasses                                             int16 // Key passes | ключевые передачи
	StandardPasses                                        int16 // Assists | кол-во голевых передач
	Crosses                                               int16 // Crosses | навес
	KeyPassesAccurate                                     int16 // Key passes accurate | точность ключевых передач
	StandardPassesAccurate                                int16 // Передачи со стандартов точные | value_sum | передачи которые лостигают створа ворот
	CrossesAccurate                                       int16 // Crosses accurate, % | точность навесов
	PassesAccuratePercentage                              int16 // Passes accurate, % | Общая точность передач
	PassesToTheFinalThirdPartOfTheField                   int16 // Passes to the final third of the field | Передачи в финальной трети
	PassesToTheFinalThirdPartOfTheFieldAccuratePercentage int16 // Accurate passes into the final third of the pitch | Точность передач в финальной трети
	PassesIntoThePenaltyBox                               int16 // Passes into the penalty box | Передачи в штрафную
	PassesIntoThePenaltyBoxAccurate                       int16 // Passes into the penalty box accurate |  Точность передач в штрафную
	LongPasses                                            int16 // Forward passes | Длинные передачи
	LongPassesAccurate                                    int16 // Передачи длинные точные, %
	OverLongPasses                                        int16 // Передачи сверхдлинные точные
	OverLongPassesAccurate                                int16 // Передачи сверхдлинные точные, %
	StandardAccuratePercentage                            int16 // Передачи со стандартов точные, %
	PassesFromLegAccuratePercentage                       int16 // Передачи с игры ногами точные, %
	PassesInOneContactPercentage                          int16 // Передачи в одно касание точные, %
	ShortPassesAccuratePercentage                         int16 // Передачи короткие точные, %
	MediumPassesAccuratePercentage                        int16 // Передачи средние точные, %
	ForwardPassesAccuratePercentage                       int16 // accurate forward passes (180 degree angle)
}

type ChallengesPlayerStat struct {
	DefensiveChallenges              int16 // Defensive challenges | Единоборства в обороне
	DefensiveChallengesWonPercentage int16 // Defensive challenges won, % | Выйгранные единоборства в обороне
	AttackingChallenges              int16 // Attacking challenges | Единоборства в атаке
	AttackingChallengesWonPercentage int16 // Attacking challenges won, % | Выйгранные единоборства в атаке
	AirChallenges                    int16 // Air challenges | Единоборства в воздухе
	AirChallengesWonPercentage       int16 // Air challenges won, % | Выйгранные единоборства в воздухе
	Tackles                          int16 // Tackles | Отборы мяча
	TacklesAccuratePercentage        int16 // Tackles successful, % | Отборы мяча
	Dribbles                         int16 // Dribbles | Обводок противника
	DribblesSuccessfulPercentage     int16 // Dribbles successful, % | Успешных обводок противника
	Interceptions                    int16 // Interceptions | Перехваты
	Recoveries                       int16 // Recoveries |  Подборы
	YellowCards                      int16 // Yellow cards | Желтые карточки
	RedCards                         int16 // Red cards | Красные карточки
}

type GoalkeeperStat struct {
	OneOnOne                                      int16 // Goalkeeper - One-on-Ones
	OneOnOneAccuratePercentage                    int16 // Goalkeeper - One-on-Ones successful, %
	Saves                                         int16 // Saves
	SavesAccuratePercentage                       int16 // Shots saved, %
	OpponentsAttacksFromPenalty                   int16 // Вратарь - Атаки соперника со стандартов
	OpponentsAttacksFromPenaltyAccuratePercentage int16 // Вратарь - Атаки соперника со стандартов с ударами, %
	HandPasses                                    int16 // Goalkeeper – hand passes
	HandPassesAccuratePercentage                  int16 // Вратарь - передачи рукой точные, %
	ShotsOnTargetFixed                            int16 // Вратарь - удары в створ зафиксированные
	ShotsOnTargetFixedAccuratePercentage          int16 // Вратарь - удары в створ зафиксированные, %
	ConcededGoals                                 int16 //Goals conceded int64 // Вратарь - пропущенные мячи
}

type PlayerStat struct {
	Firstname             string
	Lastname              string
	Position              string
	CommonPlayerStat      CommonPlayerStat
	PassesPlayerStatistic PassesPlayerStatistic
	ChallengesPlayerStat  ChallengesPlayerStat
	GoalkeeperStat        GoalkeeperStat
}

type PlayerStatRecord struct {
	ID                 string
	CommonStatID       string
	PassesStatID       string
	InteractionsStatID string
	GoalkeeperStatID   string
}

func (p *Player) TakePlayerInfo(playersStruct *[]Player, config projectConfig.Config) (*[]Player, error) {
	var players []map[string]interface{}
	for _, player := range *playersStruct {
		rustatApi := fmt.Sprintf("http://feeds.rustatsport.ru/?tpl=11&user=%s&key=%s&player_id=%s&lang_id=1&format=json", config.LoginRustat, config.PasswordRustat, player.ID)

		graphQLObject := GraphQLRequest{}
		result, sendRequestErr := graphQLObject.SendGetRequest(rustatApi)
		if sendRequestErr != nil {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Не удалось отправить запрос по url: %s | file: %s | line: %d | error: %v", rustatApi, file, line-2, sendRequestErr)
		}

		data, dataOk := result["data"]
		if !dataOk {
			_, file, line, _ := runtime.Caller(0)
			log.Printf("Не удалось извлечь игрока из тела ответа - file: %s | line: %d", file, line-2)
			continue
		}

		rows, rowsOk := data.(map[string]interface{})["row"].([]interface{})
		if !rowsOk {
			_, file, line, _ := runtime.Caller(0)
			log.Printf("Не удалось получить map row из тела ответа - file: %s | line: %d", file, line-2)
			continue
		}

		for _, value := range rows {
			playerFromResp, playerOk := value.(map[string]interface{})
			if !playerOk {
				_, file, line, _ := runtime.Caller(0)
				log.Printf("Не удалось получить массив игрока - file: %s | line: %d", file, line-2)
				continue
			}
			players = append(players, playerFromResp)
		}

		if len(players) == 0 {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Не удалось вытащить информацию ни об одном игроке - file: %s | line: %d", file, line-1)
		}

	}
	for index, playerOld := range *playersStruct {
		for _, player := range players {
			if playerOld.ID == player["id"] {
				if player["birthday"] != nil {
					birthday := player["birthday"].(string)
					(*playersStruct)[index].Birthdate = &birthday
				}
				if player["gender_id"] != nil {
					genderID := player["gender_id"].(string)
					(*playersStruct)[index].GenderID = &genderID
				}
				if player["gender_name"] != nil {
					genderName := player["gender_name"].(string)
					(*playersStruct)[index].GenderName = &genderName
				}
				if player["nickname"] != nil {
					nickname := player["nickname"].(string)
					(*playersStruct)[index].Nickname = &nickname
				}
				if player["height"] != nil {
					height := player["height"].(string)
					(*playersStruct)[index].Height = &height
				}
				if player["weight"] != nil {
					weight := player["weight"].(string)
					(*playersStruct)[index].Weight = &weight
				}
				if player["weight"] != nil {
					weight := player["weight"].(string)
					(*playersStruct)[index].Weight = &weight
				}
				if player["season_id"] != nil {
					seasonID := player["season_id"].(string)
					(*playersStruct)[index].SeasonID = &seasonID
				}
				if player["contract_ending"] != nil {
					contractEnding := player["contract_ending"].(string)
					(*playersStruct)[index].ContractEnding = &contractEnding
				}
			}
		}
	}
	return playersStruct, nil
}

func (p *Player) TakePlayerStats(playersStruct *[]Player, config projectConfig.Config) (*[]PlayerStat, error) {
	var resultPlayerStats []PlayerStat

	for index, _ := range *playersStruct {
		log.Println("С рустат пришла статистика игрока - ", (*playersStruct)[index].Lastname)
		rustatApi := fmt.Sprintf("http://feeds.rustatsport.ru/?tpl=61&user=%s&key=%s&tournament_id=%s&player_id=%s&season_id=%s&date_start=&date_end=&lang_id=1&format=json", config.LoginRustat, config.PasswordRustat, config.TournamentIDForRustat, (*playersStruct)[index].ID, config.SeasonIDForRustat)

		graphQLObject := GraphQLRequest{}
		result, sendRequestErr := graphQLObject.SendGetRequest(rustatApi)
		if sendRequestErr != nil {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Не удалось отправить запрос по url: %s | file: %s | line: %d | error: %v", rustatApi, file, line-2, sendRequestErr)
		}

		data, dataOk := result["data"]
		if !dataOk {
			_, file, line, _ := runtime.Caller(0)
			log.Printf("Не удалось извлечь статистику игрока из тела ответа - file: %s | line: %d", file, line-2)
			continue
		}

		rows, rowsOk := data.(map[string]interface{})["row"].([]interface{})
		if !rowsOk {
			_, file, line, _ := runtime.Caller(0)
			log.Printf("Не удалось получить map row из тела ответа - file: %s | line: %d", file, line-2)
			continue
		}

		var params []map[string]interface{}
		for _, value := range rows {
			param, paramOk := value.(map[string]interface{})
			if !paramOk {
				_, file, line, _ := runtime.Caller(0)
				log.Printf("Не удалось получить параметры статистики игрока - file: %s | line: %d", file, line-2)
				continue
			}
			params = append(params, param)
		}
		// Если не удалось преобразовать ни одного игрока
		if len(params) == 0 {
			_, file, line, _ := runtime.Caller(0)
			return nil, fmt.Errorf("Не удалось получить ни одной статистики игрока - file: %s | line: %d", file, line-1)
		}

		//for _, value := range params {
		//	log.Println(value)
		//}
		//log.Println("----------------------")

		playerCommonStat := CommonPlayerStat{}
		playerPassesStat := PassesPlayerStatistic{}
		playerChallengesStat := ChallengesPlayerStat{}
		goalkeeperStat := GoalkeeperStat{}

		var playerStat PlayerStat

		for _, value := range params {
			// Кол-во игр в команде
			switch value["param_name"] {
			//------------------Общая статистика------------------
			case "Matches played",
				"Starting lineup appearances",
				"Substitutes in",
				"Substitute out",
				"Goals",
				"Shots",
				"Shots on target",
				"Assists",
				"Shots on target, %",
				"Minutes played",
				//-----------------Статистика передач---------------
				"Passes", "Key passes",
				"Передачи со стандартов",
				"Crosses",
				"Key passes accurate",
				"Передачи со стандартов точные",
				"Crosses accurate",
				"Passes accurate, %",
				"Passes to the final third of the field",
				"Accurate passes into the final third of the pitch",
				"Passes into the penalty box",
				"Passes into the penalty box accurate, %",
				"Wide shots",
				"Передачи длинные точные, %",
				"Передачи сверхдлинные",
				"Передачи сверхдлинные точные, %",
				"Передачи со стандартов точные, %",
				"Передачи с игры ногами точные, %",
				"Передачи в одно касание точные, %",
				"Передачи короткие точные, %",
				"Передачи средние точные, %",
				"accurate forward passes (180 degree angle)%",
				//---------------Статистика взаимодействий------------
				"Defensive challenges",
				"Defensive challenges won, %",
				"Attacking challenges",
				"Attacking challenges won, %",
				"Air challenges",
				"Air challenges won, %",
				"Tackles",
				"Tackles successful, %",
				"Dribbles", "Dribbles successful, %",
				"Interceptions",
				"Recoveries",
				"Yellow cards",
				"Red cards",
				//---------------Статистика вратаря-------------------
				"Goalkeeper - One-on-Ones",
				"Goalkeeper - One-on-Ones successful, %",
				"Saves", "Shots saved, %",
				"Вратарь - Атаки соперника со стандартов",
				"Вратарь - Атаки соперника со стандартов с ударами, %",
				"Goalkeeper – hand passes", "Вратарь - передачи рукой точные, %",
				"Вратарь - удары в створ зафиксированные",
				"Вратарь - удары в створ зафиксированные, %",
				"Goals conceded":
				valueStr, valueStrOk := value["value_sum"].(string)
				if valueStrOk {
					if valueFloat, parseErr := strconv.ParseFloat(valueStr, 32); parseErr == nil {
						switch value["param_name"] {
						//------------------Общая статистика------------------
						case "Matches played":
							playerCommonStat.MatchesInTeam = int16(valueFloat)
						case "Starting lineup appearances":
							playerCommonStat.StartingLineupAppearances = int16(valueFloat)
						case "Substitutes in":
							playerCommonStat.SubstitutesIn = int16(valueFloat)
						case "Substitute out":
							playerCommonStat.SubstitutesOut = int16(valueFloat)
						case "Goals":
							playerCommonStat.Goals = int16(valueFloat)
						case "Shots":
							playerCommonStat.Shots = int16(valueFloat)
						case "Shots on target":
							playerCommonStat.ShotsOnTarget = int16(valueFloat)
						case "Assists":
							playerCommonStat.Assists = int16(valueFloat)
						case "Shots on target, %":
							playerCommonStat.ShotsOnTargetAccuratePercentage = int16(valueFloat)
						case "Minutes played":
							playerCommonStat.MinutesPlayed = int16(valueFloat)
						//-----------------Статистика передач-----------------
						case "Passes":
							playerPassesStat.Passes = int16(valueFloat)
						case "Key passes":
							playerPassesStat.KeyPasses = int16(valueFloat)
						case "Передачи со стандартов":
							playerPassesStat.StandardPasses = int16(valueFloat)
						case "Crosses":
							playerPassesStat.Crosses = int16(valueFloat)
						case "Key passes accurate":
							playerPassesStat.KeyPassesAccurate = int16(valueFloat)
						case "Передачи со стандартов точные":
							playerPassesStat.StandardPassesAccurate = int16(valueFloat)
						case "Crosses accurate":
							playerPassesStat.CrossesAccurate = int16(valueFloat)
						case "Passes accurate, %":
							playerPassesStat.PassesAccuratePercentage = int16(valueFloat)
						case "Passes to the final third of the field":
							playerPassesStat.PassesToTheFinalThirdPartOfTheField = int16(valueFloat)
						case "Accurate passes into the final third of the pitch":
							playerPassesStat.PassesToTheFinalThirdPartOfTheFieldAccuratePercentage = int16(valueFloat)
						case "Passes into the penalty box":
							playerPassesStat.PassesIntoThePenaltyBox = int16(valueFloat)
						case "Passes into the penalty box accurate, %":
							playerPassesStat.PassesIntoThePenaltyBoxAccurate = int16(valueFloat)
						case "Wide shots":
							playerPassesStat.LongPasses = int16(valueFloat)
						case "Передачи длинные точные, %":
							playerPassesStat.LongPassesAccurate = int16(valueFloat)
						case "Передачи сверхдлинные":
							playerPassesStat.OverLongPasses = int16(valueFloat)
						case "Передачи сверхдлинные точные, %":
							playerPassesStat.OverLongPassesAccurate = int16(valueFloat)
						case "Передачи со стандартов точные, %":
							playerPassesStat.StandardAccuratePercentage = int16(valueFloat)
						case "Передачи с игры ногами точные, %":
							playerPassesStat.PassesFromLegAccuratePercentage = int16(valueFloat)
						case "Передачи в одно касание точные, %":
							playerPassesStat.PassesInOneContactPercentage = int16(valueFloat)
						case "Передачи короткие точные, %":
							playerPassesStat.ShortPassesAccuratePercentage = int16(valueFloat)
						case "Передачи средние точные, %":
							playerPassesStat.MediumPassesAccuratePercentage = int16(valueFloat)
						case "accurate forward passes (180 degree angle)%":
							playerPassesStat.ForwardPassesAccuratePercentage = int16(valueFloat)
						//---------------Статистика взаимодействий------------
						case "Defensive challenges":
							playerChallengesStat.DefensiveChallenges = int16(valueFloat)
						case "Defensive challenges won, %":
							playerChallengesStat.DefensiveChallengesWonPercentage = int16(valueFloat)
						case "Attacking challenges":
							playerChallengesStat.AttackingChallenges = int16(valueFloat)
						case "Attacking challenges won, %":
							playerChallengesStat.AttackingChallengesWonPercentage = int16(valueFloat)
						case "Air challenges":
							playerChallengesStat.AirChallenges = int16(valueFloat)
						case "Air challenges won, %":
							playerChallengesStat.AirChallengesWonPercentage = int16(valueFloat)
						case "Tackles":
							playerChallengesStat.Tackles = int16(valueFloat)
						case "Tackles successful, %":
							playerChallengesStat.TacklesAccuratePercentage = int16(valueFloat)
						case "Dribbles":
							playerChallengesStat.Dribbles = int16(valueFloat)
						case "Dribbles successful, %":
							playerChallengesStat.DribblesSuccessfulPercentage = int16(valueFloat)
						case "Interceptions":
							playerChallengesStat.Interceptions = int16(valueFloat)
						case "Recoveries":
							playerChallengesStat.Recoveries = int16(valueFloat)
						case "Yellow cards":
							playerChallengesStat.YellowCards = int16(valueFloat)
						case "Red cards":
							playerChallengesStat.RedCards = int16(valueFloat)
						//---------------Статистика вратаря-------------------
						case "Goalkeeper - One-on-Ones":
							goalkeeperStat.OneOnOne = int16(valueFloat)
						case "Goalkeeper - One-on-Ones successful, %":
							goalkeeperStat.OneOnOneAccuratePercentage = int16(valueFloat)
						case "Saves":
							goalkeeperStat.Saves = int16(valueFloat)
						case "Shots saved, %":
							goalkeeperStat.SavesAccuratePercentage = int16(valueFloat)
						case "Вратарь - Атаки соперника со стандартов":
							goalkeeperStat.OpponentsAttacksFromPenalty = int16(valueFloat)
						case "Вратарь - Атаки соперника со стандартов с ударами, %":
							goalkeeperStat.OpponentsAttacksFromPenaltyAccuratePercentage = int16(valueFloat)
						case "Goalkeeper – hand passes":
							goalkeeperStat.HandPasses = int16(valueFloat)
						case "Вратарь - передачи рукой точные, %":
							goalkeeperStat.HandPassesAccuratePercentage = int16(valueFloat)
						case "Вратарь - удары в створ зафиксированные":
							goalkeeperStat.ShotsOnTargetFixed = int16(valueFloat)
						case "Вратарь - удары в створ зафиксированные, %":
							goalkeeperStat.ShotsOnTargetFixedAccuratePercentage = int16(valueFloat)
						case "Goals conceded":
							goalkeeperStat.ConcededGoals = int16(valueFloat)
						}
					} else {
						_, file, line, _ := runtime.Caller(0)
						log.Printf("Не удалось преобразовать в int значение - file: %s | line: %d | error: %v", file, line-122, parseErr)
					}
				}
			}
		}
		playerPassesStat.PassesToTheFinalThirdPartOfTheFieldAccuratePercentage = int16(float32(playerPassesStat.PassesToTheFinalThirdPartOfTheFieldAccuratePercentage) / float32(playerPassesStat.PassesToTheFinalThirdPartOfTheField) * 100)
		playerStat.Firstname = (*playersStruct)[index].Firstname
		playerStat.Lastname = (*playersStruct)[index].Lastname
		playerStat.Position = (*playersStruct)[index].Position1Name
		playerStat.CommonPlayerStat = playerCommonStat
		playerStat.PassesPlayerStatistic = playerPassesStat
		playerStat.ChallengesPlayerStat = playerChallengesStat
		playerStat.GoalkeeperStat = goalkeeperStat
		resultPlayerStats = append(resultPlayerStats, playerStat)
	}
	return &resultPlayerStats, nil
}

func TakeIdOfPlayerStatistic(player PlayerStat, seasonID string, leagueID string, config projectConfig.Config) (*PlayerStatRecord, error) {
	query := `
       query($id:GraphQLStringOrFloat!, $seasonID: GraphQLStringOrFloat!, $leagueID: GraphQLStringOrFloat!){
			team_members(filter:{id:{_eq:$id}, player_statistic:{player_statistic_id:{id:{_nnull:true}}}}){
				id
				firstname
				lastname
				player_statistic{
					player_statistic_id(filter:{season:{id:{_eq:$seasonID}}, ligue:{id:{_eq:$leagueID}}}){
						id
						season{id}
						ligue{id}
						common_statistic{id}
						passes_statistic{id}
						interaction_statistic{id}
						goalkeeper_statistc{id}
					}
				}
			}
		}
  	`
	var id int64
	for key, value := range config.PlayerRelations {
		if player.Firstname+" "+player.Lastname == key {
			id = value
		}
	}

	variables := map[string]interface{}{
		"id":       id,
		"seasonID": seasonID,
		"leagueID": leagueID,
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

	data, dataOk := result["data"]
	if !dataOk {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("В теле ответа ничего нет | file:%s | line %d", file, line-2)
	}

	teamMembers, teamMembersOk := data.(map[string]interface{})["team_members"].([]interface{})
	if !teamMembersOk {
		_, file, line, _ := runtime.Caller(0)
		return nil, fmt.Errorf("В data, который был извлечен ранее, нет членов команды | file:%s | line %d", file, line-2)
	}

	var record PlayerStatRecord
	for _, value := range teamMembers {
		playerStatistics := value.(map[string]interface{})["player_statistic"].([]interface{})
		for _, statistic := range playerStatistics {

			statisticID := statistic.(map[string]interface{})["player_statistic_id"]

			if statisticID != nil {
				record.ID = statistic.(map[string]interface{})["player_statistic_id"].(map[string]interface{})["id"].(string)
				commonStat := statistic.(map[string]interface{})["player_statistic_id"].(map[string]interface{})["common_statistic"]
				passesStat := statistic.(map[string]interface{})["player_statistic_id"].(map[string]interface{})["passes_statistic"]
				interactionsStat := statistic.(map[string]interface{})["player_statistic_id"].(map[string]interface{})["interaction_statistic"]
				goalkeeperStat := statistic.(map[string]interface{})["player_statistic_id"].(map[string]interface{})["goalkeeper_statistc"]

				if commonStat != nil {
					commonStatID := statistic.(map[string]interface{})["player_statistic_id"].(map[string]interface{})["common_statistic"].(map[string]interface{})["id"].(string)
					record.CommonStatID = commonStatID
				}
				if passesStat != nil {
					passesStatID := statistic.(map[string]interface{})["player_statistic_id"].(map[string]interface{})["passes_statistic"].(map[string]interface{})["id"].(string)
					record.PassesStatID = passesStatID
				}
				if interactionsStat != nil {
					interactionsStatID := statistic.(map[string]interface{})["player_statistic_id"].(map[string]interface{})["interaction_statistic"].(map[string]interface{})["id"].(string)
					record.InteractionsStatID = interactionsStatID
				}
				if goalkeeperStat != nil {
					goalkeeperStatID := statistic.(map[string]interface{})["player_statistic_id"].(map[string]interface{})["goalkeeper_statistc"].(map[string]interface{})["id"].(string)
					record.GoalkeeperStatID = goalkeeperStatID
				}
			}
		}
	}
	return &record, nil
}

func (p *PlayerStat) SendPlayerAllStatistic(player *PlayerStat, record *PlayerStatRecord, config projectConfig.Config) error {
	if record.CommonStatID != "" {
		if sendPlayerCommonStatisticErr := player.SendPlayerCommonStatistic(*player, *record, config); sendPlayerCommonStatisticErr != nil {
			_, file, line, _ := runtime.Caller(0)
			return fmt.Errorf("Не удалось отправить общую статистику игрока: file: %s | line: %d | error: %v", file, line-1, sendPlayerCommonStatisticErr)
		}
	}
	if record.PassesStatID != "" {
		if sendPlayerPassesStatisticErr := player.SendPlayerPassesStatistic(*player, *record, config); sendPlayerPassesStatisticErr != nil {
			_, file, line, _ := runtime.Caller(0)
			return fmt.Errorf("Не удалось отправить статистику передач игрока: file: %s | line: %d | error: %v", file, line-1, sendPlayerPassesStatisticErr)
		}
	}
	if record.InteractionsStatID != "" {
		if sendPlayerInteractionStatisticErr := player.SendPlayerInteractionStatistic(*player, *record, config); sendPlayerInteractionStatisticErr != nil {
			_, file, line, _ := runtime.Caller(0)
			return fmt.Errorf("Не удалось отправить статистику взаимодействий игрока: file: %s | line: %d | error: %v", file, line-1, sendPlayerInteractionStatisticErr)
		}
	}
	if record.GoalkeeperStatID != "" {
		if sendPlayerGoalkeeperStatisticErr := player.SendPlayerGoalkeeperStatistic(*player, *record, config); sendPlayerGoalkeeperStatisticErr != nil {
			_, file, line, _ := runtime.Caller(0)
			return fmt.Errorf("Не удалось отправить статистику вратаря: file: %s | line: %d | error: %v", file, line-1, sendPlayerGoalkeeperStatisticErr)
		}
	}
	return nil
}

func (p *PlayerStat) SendPlayerCommonStatistic(player PlayerStat, record PlayerStatRecord, config projectConfig.Config) error {
	query := `
	 mutation($id: ID!, $data: update_overall_player_statistic_input!){
		update_overall_player_statistic_item(id: $id, data: $data) {
			id
			amount_of_matches
			starting_lineup_appearances
			substitutes_in
			substitutes_out
			goals_scored
			shots
			shots_on_goal
			assists
			amount_of_accurate_shots
			minutes_on_the_field
		}
	}`
	variables := map[string]interface{}{
		"id": record.CommonStatID,
		"data": map[string]interface{}{
			"amount_of_matches":           player.CommonPlayerStat.MatchesInTeam,
			"starting_lineup_appearances": player.CommonPlayerStat.StartingLineupAppearances,
			"substitutes_in":              player.CommonPlayerStat.SubstitutesIn,
			"substitutes_out":             player.CommonPlayerStat.SubstitutesOut,
			"goals_scored":                player.CommonPlayerStat.Goals,
			"shots":                       player.CommonPlayerStat.Shots,
			"shots_on_goal":               player.CommonPlayerStat.ShotsOnTarget,
			"assists":                     player.CommonPlayerStat.Assists,
			"amount_of_accurate_shots":    player.CommonPlayerStat.ShotsOnTargetAccuratePercentage,
			"minutes_on_the_field":        player.CommonPlayerStat.MinutesPlayed,
		},
	}

	graphqlRequest := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	result, sendRequestErr := graphqlRequest.SendPostRequest(graphqlRequest, config)
	if sendRequestErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Не удалось отправить graphQL запрос общей статистики игрока: file: %s | line: %d | error: %v", file, line-2, sendRequestErr)
	}

	log.Println("Статистика игрока "+player.Firstname+" "+player.Lastname+" успешно отправлена: ", result)

	return nil
}

func (p *PlayerStat) SendPlayerPassesStatistic(player PlayerStat, record PlayerStatRecord, config projectConfig.Config) error {
	query := `
	 mutation($id: ID!, $data: update_second_section_of_player_statistics_input!){
		update_second_section_of_player_statistics_item(id: $id, data: $data) {
			id
			name_of_player
			lastname_of_player
			total_transmissions
			key_transmissions
			key_accurate
			passes_from_standard_position
			passes_from_standard_position_accurate
			a_cross
			a_cross_accurate
			average_accuracy
			penalty_kicks
			successful_penalties
			final_third_strikes
			final_third_successful
			long_passes
			successful_long_passes
			overlong_passes
			successful_overlong_passes
			passes_from_standard_position_percent
			passes_from_legs_percent
			passes_in_one_touch_percent
			short_passes_percent
			medium_passes_percent
			passes_in_180_degree_percent
		}
	}
 `
	variables := map[string]interface{}{
		"id": record.PassesStatID,
		"data": map[string]interface{}{
			"total_transmissions":                    player.PassesPlayerStatistic.Passes,
			"key_transmissions":                      player.PassesPlayerStatistic.KeyPasses,
			"key_accurate":                           player.PassesPlayerStatistic.KeyPassesAccurate,
			"passes_from_standard_position":          player.PassesPlayerStatistic.StandardPasses,
			"passes_from_standard_position_accurate": player.PassesPlayerStatistic.StandardPassesAccurate,
			"a_cross":                                player.PassesPlayerStatistic.Crosses,
			"a_cross_accurate":                       player.PassesPlayerStatistic.CrossesAccurate,
			"average_accuracy":                       player.PassesPlayerStatistic.PassesAccuratePercentage,
			"penalty_kicks":                          player.PassesPlayerStatistic.PassesIntoThePenaltyBox,
			"successful_penalties":                   player.PassesPlayerStatistic.PassesIntoThePenaltyBoxAccurate,
			"final_third_strikes":                    player.PassesPlayerStatistic.PassesToTheFinalThirdPartOfTheField,
			"final_third_successful":                 player.PassesPlayerStatistic.PassesToTheFinalThirdPartOfTheFieldAccuratePercentage,
			"long_passes":                            player.PassesPlayerStatistic.LongPasses,
			"successful_long_passes":                 player.PassesPlayerStatistic.LongPassesAccurate,
			"overlong_passes":                        player.PassesPlayerStatistic.OverLongPasses,
			"successful_overlong_passes":             player.PassesPlayerStatistic.OverLongPassesAccurate,
			"passes_from_standard_position_percent":  player.PassesPlayerStatistic.StandardAccuratePercentage,
			"passes_from_legs_percent":               player.PassesPlayerStatistic.PassesFromLegAccuratePercentage,
			"passes_in_one_touch_percent":            player.PassesPlayerStatistic.PassesInOneContactPercentage,
			"short_passes_percent":                   player.PassesPlayerStatistic.ShortPassesAccuratePercentage,
			"medium_passes_percent":                  player.PassesPlayerStatistic.MediumPassesAccuratePercentage,
			"passes_in_180_degree_percent":           player.PassesPlayerStatistic.ForwardPassesAccuratePercentage,
		},
	}

	graphqlRequest := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	result, sendRequestErr := graphqlRequest.SendPostRequest(graphqlRequest, config)
	if sendRequestErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Не удалось отправить graphQL запрос статистики передач игрока: file: %s | line: %d | error: %v", file, line-2, sendRequestErr)
	}

	log.Println("Статистика игрока "+player.Firstname+" "+player.Lastname+" успешно отправлена: ", result)

	return nil
}

func (p *PlayerStat) SendPlayerInteractionStatistic(player PlayerStat, record PlayerStatRecord, config projectConfig.Config) error {
	query := `
	 mutation($id: ID!, $data: update_player_interaction_input!){
		update_player_interaction_item(id: $id, data: $data) {
			id
			name_of_player
			lastname_of_player
			defensive_duels
			successful_defensive_duels
			offensive_duels
			successful_offensive_duels
			duels_upstairs
			successful_duels_upstairs
			duels_for_the_ball
			successful_duels_for_the_ball
			dribbling
			successful_dribbling
			interceptions_balls
			selections
			yellow_cards
			red_cards
		}
	}
 `
	variables := map[string]interface{}{
		"id": record.InteractionsStatID,
		"data": map[string]interface{}{
			"defensive_duels":               player.ChallengesPlayerStat.DefensiveChallenges,
			"successful_defensive_duels":    player.ChallengesPlayerStat.DefensiveChallengesWonPercentage,
			"offensive_duels":               player.ChallengesPlayerStat.AttackingChallenges,
			"successful_offensive_duels":    player.ChallengesPlayerStat.AttackingChallengesWonPercentage,
			"duels_upstairs":                player.ChallengesPlayerStat.AirChallenges,
			"successful_duels_upstairs":     player.ChallengesPlayerStat.AirChallengesWonPercentage,
			"duels_for_the_ball":            player.ChallengesPlayerStat.Tackles,
			"successful_duels_for_the_ball": player.ChallengesPlayerStat.TacklesAccuratePercentage,
			"dribbling":                     player.ChallengesPlayerStat.Dribbles,
			"successful_dribbling":          player.ChallengesPlayerStat.DribblesSuccessfulPercentage,
			"interceptions_balls":           player.ChallengesPlayerStat.Interceptions,
			"selections":                    player.ChallengesPlayerStat.Recoveries,
			"yellow_cards":                  player.ChallengesPlayerStat.YellowCards,
			"red_cards":                     player.ChallengesPlayerStat.RedCards,
		},
	}

	graphqlRequest := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	result, sendRequestErr := graphqlRequest.SendPostRequest(graphqlRequest, config)
	if sendRequestErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Не удалось отправить graphQL запрос татистики взаимодействий игрока: file: %s | line: %d | error: %v", file, line-2, sendRequestErr)
	}

	log.Println("Статистика игрока "+player.Firstname+" "+player.Lastname+" успешно отправлена: ", result)

	return nil
}

func (p *PlayerStat) SendPlayerGoalkeeperStatistic(player PlayerStat, record PlayerStatRecord, config projectConfig.Config) error {
	query := `
	 mutation($id: ID!, $data: update_goalkeeper_statistic_input!){
		update_goalkeeper_statistic_item(id: $id, data: $data) {
			id
			name_of_player
			last_name_of_player
			one_on_one
			one_on_one_percent
			saves
			saves_percent
			opponent_attacks_from_penalty
			opponent_attacks_from_penalty_percent
			hand_passes
			hand_passes_percent
			shots_on_target_fixed
			shots_on_target_fixed_percent
			conceded_goals
			red_cards
			yellow_cards  
		}
	}
 `
	variables := map[string]interface{}{
		"id": record.GoalkeeperStatID,
		"data": map[string]interface{}{
			"one_on_one":                            player.GoalkeeperStat.OneOnOne,
			"one_on_one_percent":                    player.GoalkeeperStat.OneOnOneAccuratePercentage,
			"saves":                                 player.GoalkeeperStat.Saves,
			"saves_percent":                         player.GoalkeeperStat.SavesAccuratePercentage,
			"opponent_attacks_from_penalty":         player.GoalkeeperStat.OpponentsAttacksFromPenalty,
			"opponent_attacks_from_penalty_percent": player.GoalkeeperStat.OpponentsAttacksFromPenaltyAccuratePercentage,
			"hand_passes":                           player.GoalkeeperStat.HandPasses,
			"hand_passes_percent":                   player.GoalkeeperStat.HandPassesAccuratePercentage,
			"shots_on_target_fixed":                 player.GoalkeeperStat.ShotsOnTargetFixed,
			"shots_on_target_fixed_percent":         player.GoalkeeperStat.ShotsOnTargetFixedAccuratePercentage,
			"conceded_goals":                        player.GoalkeeperStat.ConcededGoals,
			"yellow_cards":                          player.ChallengesPlayerStat.YellowCards,
			"red_cards":                             player.ChallengesPlayerStat.RedCards,
		},
	}

	graphqlRequest := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	result, sendRequestErr := graphqlRequest.SendPostRequest(graphqlRequest, config)
	if sendRequestErr != nil {
		_, file, line, _ := runtime.Caller(0)
		return fmt.Errorf("Не удалось отправить graphQL запрос статистики вратаря: file: %s | line: %d | error: %v", file, line-2, sendRequestErr)
	}

	log.Println("Статистика игрока "+player.Firstname+" "+player.Lastname+" успешно отправлена: ", result)

	return nil
}

func (p *PlayerStat) PrintPlayer(players *[]PlayerStat) {
	for index, player := range *players {
		log.Println(index)

		fmt.Println("Player Name | Имя игрока:", player.Firstname, player.Lastname)
		fmt.Println("Matches Played | Количество матчей:", player.CommonPlayerStat.MatchesInTeam)
		fmt.Println("Starting Lineup Appearances | Начал в стартовом составе:", player.CommonPlayerStat.StartingLineupAppearances)
		fmt.Println("Substitutes In | Вышел на замену:", player.CommonPlayerStat.SubstitutesIn)
		fmt.Println("Substitutes Out | Был заменён:", player.CommonPlayerStat.SubstitutesOut)
		fmt.Println("Position | Позиция игрока, для понятия, вратарь, чи не:", player.Position)
		fmt.Println("Goals | Голы:", player.CommonPlayerStat.Goals)
		fmt.Println("Shots | Удары:", player.CommonPlayerStat.Shots)
		fmt.Println("Shots on Target | Удары в створ:", player.CommonPlayerStat.ShotsOnTarget)
		fmt.Println("Assists | Голевые моменты:", player.CommonPlayerStat.Assists)
		fmt.Println("Shots on Target Accuracy, % | Точность ударов в створ:", player.CommonPlayerStat.ShotsOnTargetAccuratePercentage)
		fmt.Println("Minutes Played | Минуты на поле:", player.CommonPlayerStat.MinutesPlayed)
		fmt.Println("Passes | Передачи:", player.PassesPlayerStatistic.Passes)
		fmt.Println("Key Passes | Ключевые передачи:", player.PassesPlayerStatistic.KeyPasses)
		fmt.Println("Key Pass Accuracy | Точные ключевые передачи:", player.PassesPlayerStatistic.KeyPassesAccurate)
		fmt.Println("Standards Passes | Передачи со стандартов:", player.PassesPlayerStatistic.StandardPasses)
		fmt.Println("Standards Passes Accurate | Точность передач со стандартов:", player.PassesPlayerStatistic.StandardPassesAccurate)
		fmt.Println("Crosses | Навесы:", player.PassesPlayerStatistic.Crosses)
		fmt.Println("Crosses Accuracy | Точные навесы:", player.PassesPlayerStatistic.CrossesAccurate)
		fmt.Println("Pass Accuracy, % | Общая точность передач:", player.PassesPlayerStatistic.PassesAccuratePercentage)
		fmt.Println("Passes into the Penalty Box | Передачи в штрафную:", player.PassesPlayerStatistic.PassesIntoThePenaltyBox)
		fmt.Println("Pass Accuracy into the Penalty Box, % | Точность передач в штрафную:", player.PassesPlayerStatistic.PassesIntoThePenaltyBoxAccurate)
		fmt.Println("Passes to the Final Third of the Field | Передачи в финальную треть поля:", player.PassesPlayerStatistic.PassesToTheFinalThirdPartOfTheField)
		fmt.Println("Pass Accuracy to the Final Third of the Field, % | Точность передач в финальную треть:", player.PassesPlayerStatistic.PassesToTheFinalThirdPartOfTheFieldAccuratePercentage)
		fmt.Println("Long Passes | Длинные передачи:", player.PassesPlayerStatistic.LongPasses)
		fmt.Println("Long Pass Accuracy, % | Точность длинных передач:", player.PassesPlayerStatistic.LongPassesAccurate)
		fmt.Println("Over Long Passes | Сверхдлинные передачи:", player.PassesPlayerStatistic.OverLongPasses)
		fmt.Println("Over Long Pass Accuracy, % | Точность сверхдлинных передач:", player.PassesPlayerStatistic.OverLongPassesAccurate)
		fmt.Println("StandardAccuratePercentage % | Передачи со стандартов:", player.PassesPlayerStatistic.StandardAccuratePercentage)
		fmt.Println("PassesFromLegAccuratePercentage % | Передачи с игры ногами:", player.PassesPlayerStatistic.PassesFromLegAccuratePercentage)
		fmt.Println("PassesInOneContactPercentage % | Передачи в одно касание:", player.PassesPlayerStatistic.PassesInOneContactPercentage)
		fmt.Println("ShortPassesAccuratePercentage % | Короткие передачи:", player.PassesPlayerStatistic.ShortPassesAccuratePercentage)
		fmt.Println("MediumPassesAccuratePercentage % | Средние передачи:", player.PassesPlayerStatistic.MediumPassesAccuratePercentage)
		fmt.Println("ForwardPassesAccuratePercentage % | Передачи вперед под углом 180 градусов:", player.PassesPlayerStatistic.ForwardPassesAccuratePercentage)
		fmt.Println("Defensive Challenges | Единоборства в обороне:", player.ChallengesPlayerStat.DefensiveChallenges)
		fmt.Println("Defensive Challenges Won, % | Выйгранные единоборства в обороне:", player.ChallengesPlayerStat.DefensiveChallengesWonPercentage)
		fmt.Println("Attacking Challenges | Единоборства в атаке:", player.ChallengesPlayerStat.AttackingChallenges)
		fmt.Println("Attacking Challenges Won, % | Выйгранные единоборства в атаке:", player.ChallengesPlayerStat.AttackingChallengesWonPercentage)
		fmt.Println("Air Challenges | Единоборства в воздухе:", player.ChallengesPlayerStat.AirChallenges)
		fmt.Println("Air Challenges Won, % | Выйгранные единоборства в воздухе:", player.ChallengesPlayerStat.AirChallengesWonPercentage)
		fmt.Println("Dribbles | Обводки противника:", player.ChallengesPlayerStat.Dribbles)
		fmt.Println("Dribble Success, % | Успешность обводок:", player.ChallengesPlayerStat.DribblesSuccessfulPercentage)
		fmt.Println("Tackles | Отборы мяча:", player.ChallengesPlayerStat.Tackles)
		fmt.Println("Tackle Success, % | Успешность отборов:", player.ChallengesPlayerStat.TacklesAccuratePercentage)
		fmt.Println("Interceptions | Перехваты мяча:", player.ChallengesPlayerStat.Interceptions)
		fmt.Println("Recoveries | Подборы мяча:", player.ChallengesPlayerStat.Recoveries)
		fmt.Println("Yellow Cards | Жёлтые карточки:", player.ChallengesPlayerStat.YellowCards)
		fmt.Println("Red Cards | Красные карточки:", player.ChallengesPlayerStat.RedCards)
		fmt.Println("OneOnOne | Один на один:", player.GoalkeeperStat.OneOnOne)
		fmt.Println("OneOnOneAccuratePercentage | Один на один %:", player.GoalkeeperStat.OneOnOneAccuratePercentage)
		fmt.Println("Saves | Сэйвы", player.GoalkeeperStat.Saves)
		fmt.Println("SavesAccuratePercentage | Сэйвы %:", player.GoalkeeperStat.SavesAccuratePercentage)
		fmt.Println("OpponentsAttacksFromPenalty | Атаки соперника со штрафных и свободных:", player.GoalkeeperStat.OpponentsAttacksFromPenalty)
		fmt.Println("OpponentsAttacksFromPenaltyAccuratePercentage | Атаки соперника со штрафных и свободных %:", player.GoalkeeperStat.OpponentsAttacksFromPenaltyAccuratePercentage)
		fmt.Println("HandPasses | Передачи рукой", player.GoalkeeperStat.HandPasses)
		fmt.Println("HandPassesAccuratePercentage | Передачи рукой %:", player.GoalkeeperStat.HandPassesAccuratePercentage)
		fmt.Println("ShotsOnTargetFixed | Удары в створ зафиксированные", player.GoalkeeperStat.ShotsOnTargetFixed)
		fmt.Println("ShotsOnTargetFixedAccuratePercentage | Удары в створ зафиксированные %:", player.GoalkeeperStat.ShotsOnTargetFixedAccuratePercentage)
		fmt.Println("ConcededGoals | Пропущенные мячи", player.GoalkeeperStat.ConcededGoals)
		fmt.Println("------------------------------------------------------------------------------------------------")
	}
}
