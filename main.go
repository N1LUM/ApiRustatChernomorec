package main

import (
	projectConfig "ApiChernoToRustat/config"
	"ApiChernoToRustat/src/team"
	"fmt"
	"log"
	"runtime"
)

func main() {
	//Инициализируем наш кфг
	config := projectConfig.InitConfig()
	config.PrintConfig()

	//Вставка состава коннды на матч
	// Это я потом вынесу в yaml, сейчас пока данная функция в тесте
	matchIDS := []string{"272235",
		"343393",
		"343386",
		"343374",
		"343363",
		"343357",
		"326022",
		"280562",
		"343420",
		"343574",
		"343566",
		"343555",
		"343542",
		"343529",
		"343522",
		"343510",
		"343502",
		"343495",
		"343480",
		"343475",
		"491975",
		"343467",
		"343458",
		"487898",
		"343445",
		"343440",
		"343431",
		"474110"}

	for _, matchID := range matchIDS {
		result, _ := team.TakeCompositionOfTeamForMatch(matchID, config)
		log.Print(result)
		log.Println("-----------------------------------------------------------------------")
	}

	//Статистика игрока

	teamObject := team.Team{}

	//Вызываем метод API Рустата для того, чтобы получить команду по ID
	resultTeam, takeTeamFromRustatErr := teamObject.TakeTeamFromRustat(config)
	if takeTeamFromRustatErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Ошибка получения команды: file: %s | line: %d | error: %v", file, line-2, takeTeamFromRustatErr)
	}
	//Вызываем метод API Рустата для того, чтобы получить состав полученной нами ранее команды
	resultStruct, takeStructureOfTeamErr := teamObject.TakeStructureOfTeam(*resultTeam, config)
	if takeStructureOfTeamErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Ошибка получения состава команд: file: %s | line: %d | error: %v", file, line-2, takeStructureOfTeamErr)
	}

	playerObject := team.Player{}
	// Собираем общую инфу по игроку, имя, фамилия и т.д
	resultPlayer, takePlayerInfoErr := playerObject.TakePlayerInfo(resultStruct, config)
	if takePlayerInfoErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Не удалось получить инофрмацию об игроках: file: %s | line: %d | error: %v", file, line-2, takePlayerInfoErr)
	}
	// Собираем по нему статистику
	resultStat, takePlayerStatsErr := playerObject.TakePlayerStats(resultPlayer, config)
	if takePlayerStatsErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Не удалось получить статистику игрока: file: %s | line: %d | error: %v", file, line-2, takePlayerStatsErr)
	}
	//Смотрим что у нас получилось
	playerStatObject := team.PlayerStat{}
	playerStatObject.PrintPlayer(resultStat)

	for _, player := range *resultStat {
		//Ищем id записи игрока
		record, takeIdOfPlayerStatisticErr := team.TakeIdOfPlayerStatistic(player, config.SeasonID, config.LeagueID, config)
		if takeIdOfPlayerStatisticErr != nil {
			_, file, line, _ := runtime.Caller(0)
			log.Fatalf("Не удалось получить IDs статистик игроков: file: %s | line: %d | error: %v", file, line-2, takeIdOfPlayerStatisticErr)
		}
		//Если нашлось, закидываем в directus
		if record != nil {
			log.Println("ID Записи статистики в целом - ", record.ID,
				" | ID Общей статистики - ", record.CommonStatID,
				" | ID Статистики передач - ", record.PassesStatID,
				" | ID Статистики взаимодействий - ", record.InteractionsStatID,
				" | ID Статистики вратаря - ", record.GoalkeeperStatID)
			if sendPlayerAllStatisticErr := player.SendPlayerAllStatistic(&player, record, config); sendPlayerAllStatisticErr != nil {
				_, file, line, _ := runtime.Caller(0)
				log.Printf("Не удалось отправить статистику игрока: file: %s | line: %d | error: %v", file, line-1, sendPlayerAllStatisticErr)
			}
		} else {
			log.Println("У игрока "+player.Firstname+" "+player.Lastname, " нет статистики")
		}
	}
	//Турнирная таблица ------------------------------------------------------------------------------------------------
	tournamentTableTeamObject := team.TournamentTableTeam{}
	//Берем по API таблицу
	tournamentTable, takeTournamentTableErr := tournamentTableTeamObject.TakeTournamentTable(config)
	if takeTournamentTableErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Ошибка получения турнирной таблицы: file: %s | line: %d | error: %v", file, line-2, takeTournamentTableErr)
	}
	//Выводим результат
	if printTournamentTableErr := team.PrintTournamentTable(tournamentTable); printTournamentTableErr != nil {
		_, file, line, _ := runtime.Caller(0)
		fmt.Printf("Ошибка вывода турнирной таблицы: file: %s | line: %d | error: %v", file, line-1, printTournamentTableErr)
	}
	// Берем записи по нужным нам сезонам и лигам
	recordIDs, takeTeamRecordIDsFromRatingTableErr := team.TakeTeamRecordIDsFromRatingTable(config)
	if takeTeamRecordIDsFromRatingTableErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Не удалось получить записи: file: %s | line: %d | error: %v", file, line-2, takeTeamRecordIDsFromRatingTableErr)
	}
	//Создаем отношение id записи в директусе и id команды
	relations, takeRelationsBetweenRecordAndTeamStatErr := team.TakeRelationsBetweenRecordAndTeamStat(*recordIDs, config)
	if takeRelationsBetweenRecordAndTeamStatErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Не удалось создать связь команды и ее записи в Directus: file: %s | line: %d | error: %v", file, line-2, takeRelationsBetweenRecordAndTeamStatErr)
	}
	// Отправляем таблицу в директус
	if sendTableErr := team.SendTable(tournamentTable, *relations, config); sendTableErr != nil {
		_, file, line, _ := runtime.Caller(0)
		log.Fatalf("Не удалось отправить турнирную таблицу в Directus: file: %s | line: %d | error: %v", file, line-1, sendTableErr)
	}
}
