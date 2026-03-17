package stats

import (
	"github.com/zollidan/esmeralda/internal/models"
)

type Column struct {
	Header string
	Value  func(g *models.Game) interface{}
}

var Columns = []Column{
	// Базовая информация
	{"число", func(g *models.Game) interface{} { return g.Day }},
	{"месяц", func(g *models.Game) interface{} { return g.Month }},
	{"год", func(g *models.Game) interface{} { return g.Year }},
	{"время", func(g *models.Game) interface{} { return g.Time }},
	{"команда_1", func(g *models.Game) interface{} { return g.HomeTeam }},
	{"команда_2", func(g *models.Game) interface{} { return g.AwayTeam }},
	{"лига", func(g *models.Game) interface{} { return g.League }},

	{"Сумма очных игр на поле команда 1", func(g *models.Game) interface{} { return g.H2HHomeMatches }},
	{"Побед на своем поле в очных играх на поле команда 1", func(g *models.Game) interface{} { return g.H2HHomeWin1 }},
	{"Ничьи в очных играх на поле команда 1", func(g *models.Game) interface{} { return g.H2HHomeDraw }},
	{"Поражений на своем поле в очных играх на поле команда 1", func(g *models.Game) interface{} { return g.H2HHomeWin2 }},

	{"общее количество матчей дома первой команды", func(g *models.Game) interface{} { return g.HomeTeamHomeMatches }},
	{"победа на своем поле", func(g *models.Game) interface{} { return g.HomeTeamHomeWin }},
	{"ничья на своем поле", func(g *models.Game) interface{} { return g.HomeTeamHomeDraw }},
	{"поражение на своем поле", func(g *models.Game) interface{} { return g.HomeTeamHomeLoss }},

	{"Общее количество матчей в гостях второй команды", func(g *models.Game) interface{} { return g.AwayTeamAwayMatches }},
	{"Поражения команды 2 в гостях", func(g *models.Game) interface{} { return g.AwayTeamAwayLoss }},
	{"Игра в гостях. Ничьи", func(g *models.Game) interface{} { return g.AwayTeamAwayDraw }},
	{"Победы команды 2 в гостях", func(g *models.Game) interface{} { return g.AwayTeamAwayWin }},

	{"Очные встречи. Свое поле. Тотал 2.5 Б/М", func(g *models.Game) interface{} { return g.H2HHomeOver25Total }},
	{"Очные встречи. Свое поле. Тотал 2.5 Б", func(g *models.Game) interface{} { return g.H2HHomeOver25Over }},
	{"Очные встречи. Свое поле. Тотал 2.5 М", func(g *models.Game) interface{} { return g.H2HHomeOver25Under }},

	{"Все встречи. Свое поле. Тотал 2.5 Б/М", func(g *models.Game) interface{} { return g.HomeAllOver25Total }},
	{"Все встречи. Свое поле. Тотал 2.5 Б", func(g *models.Game) interface{} { return g.HomeAllOver25Over }},
	{"Все встречи. Свое поле. Тотал 2.5 М", func(g *models.Game) interface{} { return g.HomeAllOver25Under }},

	{"Все встречи. Гостевое поле. Тотал 2.5 Б/М", func(g *models.Game) interface{} { return g.AwayAllOver25Total }},
	{"Все встречи. Гостевое поле. Тотал 2.5 Б", func(g *models.Game) interface{} { return g.AwayAllOver25Over }},
	{"Все встречи. Гостевое поле. Тотал 2.5 М", func(g *models.Game) interface{} { return g.AwayAllOver25Under }},

	{"Кол-во игр очн. (25) команда 1 и команда 2 на любом поле", func(g *models.Game) interface{} { return g.H2HAnyGames25 }},
	{"Сумма мячей очн. (25)", func(g *models.Game) interface{} { return g.H2HAnyGoals25 }},
	{"Кол-во игр очн. (5) команда 1 и команда 2 на любом поле", func(g *models.Game) interface{} { return g.H2HAnyGames5 }},
	{"Сумма мячей очн. (5)", func(g *models.Game) interface{} { return g.H2HAnyGoals5 }},
	{"Кол-во игр очн. (3) команда 1 и команда 2 на любом поле", func(g *models.Game) interface{} { return g.H2HAnyGames3 }},
	{"Сумма мячей очн. (3)", func(g *models.Game) interface{} { return g.H2HAnyGoals3 }},

	{"количество игр хозяев (25) на любом поле", func(g *models.Game) interface{} { return g.HomeAnyGames25 }},
	{"сумма мячей хозяев (25) на любом поле", func(g *models.Game) interface{} { return g.HomeAnyGoals25 }},
	{"количество игр хозяев (5) на любом поле", func(g *models.Game) interface{} { return g.HomeAnyGames5 }},
	{"сумма мячей хозяев (5) на любом поле", func(g *models.Game) interface{} { return g.HomeAnyGoals5 }},
	{"количество игр хозяев (3) на любом поле", func(g *models.Game) interface{} { return g.HomeAnyGames3 }},
	{"сумма мячей хозяев (3) на любом поле", func(g *models.Game) interface{} { return g.HomeAnyGoals3 }},

	{"количество игр гостей (25) на любом поле", func(g *models.Game) interface{} { return g.AwayAnyGames25 }},
	{"сумма мячей гостей (25) на любом поле", func(g *models.Game) interface{} { return g.AwayAnyGoals25 }},
	{"количество игр гостей (5) на любом поле", func(g *models.Game) interface{} { return g.AwayAnyGames5 }},
	{"сумма мячей гостей (5) на любом поле", func(g *models.Game) interface{} { return g.AwayAnyGoals5 }},
	{"количество игр гостей (3) на любом поле", func(g *models.Game) interface{} { return g.AwayAnyGames3 }},
	{"сумма мячей гостей (3) на любом поле", func(g *models.Game) interface{} { return g.AwayAnyGoals3 }},

	{"Кол-во игр очн. (25) команда 1 и команда 2 на поле хозяев", func(g *models.Game) interface{} { return g.H2HHomeGames25 }},
	{"Сумма мячей очн. (25) команда 1 и команда 2 на поле хозяев", func(g *models.Game) interface{} { return g.H2HHomeGoals25 }},
	{"Кол-во игр очн. (5) команда 1 и команда 2 на поле хозяев", func(g *models.Game) interface{} { return g.H2HHomeGames5 }},
	{"Сумма мячей очн. (5) команда 1 и команда 2 на поле хозяев", func(g *models.Game) interface{} { return g.H2HHomeGoals5 }},
	{"Кол-во игр очн. (3) команда 1 и команда 2 на поле хозяев", func(g *models.Game) interface{} { return g.H2HHomeGames3 }},
	{"Сумма мячей очн. (3) команда 1 и команда 2 на поле хозяев", func(g *models.Game) interface{} { return g.H2HHomeGoals3 }},

	{"количество игр хозяев (25) на поле хозяев", func(g *models.Game) interface{} { return g.HomeHomeGames25 }},
	{"сумма мячей хозяев (25) на поле хозяев", func(g *models.Game) interface{} { return g.HomeHomeGoals25 }},
	{"количество игр хозяев (5) на поле хозяев", func(g *models.Game) interface{} { return g.HomeHomeGames5 }},
	{"сумма мячей хозяев (5) на поле хозяев", func(g *models.Game) interface{} { return g.HomeHomeGoals5 }},
	{"количество игр хозяев (3) на поле хозяев", func(g *models.Game) interface{} { return g.HomeHomeGames3 }},
	{"сумма мячей хозяев (3) на поле хозяев", func(g *models.Game) interface{} { return g.HomeHomeGoals3 }},

	{"количество игр гостей (25) на поле гостей", func(g *models.Game) interface{} { return g.AwayAwayGames25 }},
	{"сумма мячей гостей (25) на поле гостей", func(g *models.Game) interface{} { return g.AwayAwayGoals25 }},
	{"количество игр гостей (5) на поле гостей", func(g *models.Game) interface{} { return g.AwayAwayGames5 }},
	{"сумма мячей гостей (5) на поле гостей", func(g *models.Game) interface{} { return g.AwayAwayGoals5 }},
	{"количество игр гостей (3) на поле гостей", func(g *models.Game) interface{} { return g.AwayAwayGames3 }},
	{"сумма мячей гостей (3) на поле гостей", func(g *models.Game) interface{} { return g.AwayAwayGoals3 }},
}

// Headers возвращает заголовки столбцов
func Headers() []string {
	headers := make([]string, len(Columns))
	for i, col := range Columns {
		headers[i] = col.Header
	}
	return headers
}

// GetValue возвращает значение для ячейки из игры по индексу колонки
func GetValue(g *models.Game, colIdx int) interface{} {
	if colIdx < 0 || colIdx >= len(Columns) {
		return ""
	}
	return Columns[colIdx].Value(g)
}

// ToSlice преобразует Game в слайс значений для Excel экспорта
// (без расширенных данных из JSON - они лежат отдельно)
func GameToSlice(g *models.Game) []interface{} {
	slice := make([]interface{}, len(Columns))
	for i, col := range Columns {
		slice[i] = col.Value(g)
	}
	return slice
}
