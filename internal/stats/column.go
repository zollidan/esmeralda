package stats

type column struct {
	header string
	value  func(r Row) interface{}
}

var Columns = []column{
	// базовая инфа
	{"число", func(r Row) interface{} { return r.Day }},
	{"месяц", func(r Row) interface{} { return int(r.Month) }},
	{"год", func(r Row) interface{} { return r.Year }},
	{"время", func(r Row) interface{} { return r.Time }},
	{"команда_1", func(r Row) interface{} { return r.HomeTeam }},
	{"команда_2", func(r Row) interface{} { return r.AwayTeam }},
	// {"бк", func(r Row) interface{} { return r.BK }},
	{"лига", func(r Row) interface{} { return r.League }},

	// H2H на поле хозяев — 15 матчей
	{"H2H15 всего", func(r Row) interface{} { return r.H2H15.Matches }},
	{"H2H15 победы К1", func(r Row) interface{} { return r.H2H15.Win1 }},
	{"H2H15 ничьи", func(r Row) interface{} { return r.H2H15.Draw }},
	{"H2H15 победы К2", func(r Row) interface{} { return r.H2H15.Win2 }},

	// H2H на поле хозяев — 25 матчей
	{"H2H25 всего", func(r Row) interface{} { return r.H2H25.Matches }},
	{"H2H25 победы К1", func(r Row) interface{} { return r.H2H25.Win1 }},
	{"H2H25 ничьи", func(r Row) interface{} { return r.H2H25.Draw }},
	{"H2H25 победы К2", func(r Row) interface{} { return r.H2H25.Win2 }},

	// домашние К1 — 25 матчей
	{"Дома К1 всего", func(r Row) interface{} { return r.Home25.Matches }},
	{"Дома К1 победы", func(r Row) interface{} { return r.Home25.Win1 }},
	{"Дома К1 ничьи", func(r Row) interface{} { return r.Home25.Draw }},
	{"Дома К1 поражения", func(r Row) interface{} { return r.Home25.Win2 }},

	// выездные К1 — 15 матчей
	{"Выезд К1-15 всего", func(r Row) interface{} { return r.Away15.Matches }},
	{"Выезд К1-15 победы", func(r Row) interface{} { return r.Away15.Win1 }},
	{"Выезд К1-15 ничьи", func(r Row) interface{} { return r.Away15.Draw }},
	{"Выезд К1-15 поражения", func(r Row) interface{} { return r.Away15.Win2 }},

	// выездные К1 — 25 матчей
	{"Выезд К1-25 всего", func(r Row) interface{} { return r.Away25.Matches }},
	{"Выезд К1-25 победы", func(r Row) interface{} { return r.Away25.Win1 }},
	{"Выезд К1-25 ничьи", func(r Row) interface{} { return r.Away25.Draw }},
	{"Выезд К1-25 поражения", func(r Row) interface{} { return r.Away25.Win2 }},

	// под вопросом
	// Выезд К2 — 25 матчей
	{"Выезд К2 всего", func(r Row) interface{} { return r.AwayK2_25.Matches }},
	{"Выезд К2 победы", func(r Row) interface{} { return r.AwayK2_25.Win1 }},
	{"Выезд К2 ничьи", func(r Row) interface{} { return r.AwayK2_25.Draw }},
	{"Выезд К2 поражения", func(r Row) interface{} { return r.AwayK2_25.Win2 }},

	// под вопросом
	// тоталы голов — выезд К2
	{"Выезд К2 гол >2.5 матчей", func(r Row) interface{} { return r.GoalsAway.Over25Matches }},
	{"Выезд К2 гол >2.5 итог", func(r Row) interface{} { return r.GoalsAway.Over25Total }},
	{"H2H гол >3 матчей", func(r Row) interface{} { return r.GoalsH2H.Over3Matches }},
	{"H2H гол >3 итог", func(r Row) interface{} { return r.GoalsH2H.Over3Total }},
	{"H2H гол >5 матчей", func(r Row) interface{} { return r.GoalsH2H.Over5Matches }},
	{"H2H гол >5 итог", func(r Row) interface{} { return r.GoalsH2H.Over5Total }},

	// тоталы голов — дома К1
	{"Дома гол >2.5 матчей", func(r Row) interface{} { return r.GoalsHome.Over25Matches }},
	{"Дома гол >2.5 итог", func(r Row) interface{} { return r.GoalsHome.Over25Total }},
	{"Дома гол >3 матчей", func(r Row) interface{} { return r.GoalsHome.Over3Matches }},
	{"Дома гол >3 итог", func(r Row) interface{} { return r.GoalsHome.Over3Total }},
	{"Дома гол >5 матчей", func(r Row) interface{} { return r.GoalsHome.Over5Matches }},
	{"Дома гол >5 итог", func(r Row) interface{} { return r.GoalsHome.Over5Total }},

	// тоталы — выезд К2
	{"Выезд К2 гол >2.5 матчей", func(r Row) interface{} { return r.GoalsAway.Over25Matches }},
	{"Выезд К2 гол >2.5 итог", func(r Row) interface{} { return r.GoalsAway.Over25Total }},
	{"Выезд К2 гол >3 матчей", func(r Row) interface{} { return r.GoalsAway.Over3Matches }},
	{"Выезд К2 гол >3 итог", func(r Row) interface{} { return r.GoalsAway.Over3Total }},
	{"Выезд К2 гол >5 матчей", func(r Row) interface{} { return r.GoalsAway.Over5Matches }},
	{"Выезд К2 гол >5 итог", func(r Row) interface{} { return r.GoalsAway.Over5Total }},

	// коэффициенты
	// {"фора", func(r Row) interface{} { return r.Fora }},
	// {"коэф 1", func(r Row) interface{} { return r.Coef1 }},
	// {"коэф 2", func(r Row) interface{} { return r.Coef2 }},
}

func Headers() []interface{} {
	h := make([]interface{}, len(Columns))
	for i, c := range Columns {
		h[i] = c.header
	}
	return h
}

func (r Row) ToSlice() []interface{} {
	s := make([]interface{}, len(Columns))
	for i, c := range Columns {
		s[i] = c.value(r)
	}
	return s
}