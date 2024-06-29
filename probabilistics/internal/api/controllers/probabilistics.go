package controllers

// func loadJson(body io.Reader, dest interface{}) {
// 	var err error
// 	js, err := io.ReadAll(body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err = json.Unmarshal(js, dest)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func probCreate(w http.ResponseWriter, r *http.Request) {
// 	GetBasicInfo(r)
// 	var err error
// 	body := &service.ProbCreateBody{}
// 	loadJson(r.Body, body)

// 	tx, _ := postgres.Client.Begin()

// 	p := service.CreateProbabilistic(body)
// 	err = service.SaveProbabilistic(p, true, true, tx)
// 	wrapperKey := &wrapper.WrapperKey{
// 		Shared: wrapper.SharedProps{
// 			Key:     p.Meta().Key(),
// 			MaxCard: p.Filter().Meta().MaxCard(),
// 			MaxFp:   p.Filter().Meta().MaxFp(),
// 		},
// 		FilterType:   p.Filter().Meta().FilterType(),
// 		CardinalType: p.Cardinal().Meta().CardinalType(),
// 	}
// 	wrapper.GetWrapper().Add(*wrapperKey, p)

// 	if err != nil {
// 		panic(fmt.Sprint("Can't save probabilistic for some reason", err))
// 	}

// 	w.Write([]byte(fmt.Sprint("Created probabilistic", p)))
// }

// func probExists(w http.ResponseWriter, r *http.Request) {
// 	GetBasicInfo(r)
// }

// func probAdd(w http.ResponseWriter, r *http.Request) {
// 	GetBasicInfo(r)
// }

// func probCard(w http.ResponseWriter, r *http.Request) {
// 	GetBasicInfo(r)
// }

// type probabilistic struct {
// 	Create func(http.ResponseWriter, *http.Request)
// 	Exists func(http.ResponseWriter, *http.Request)
// 	Add    func(http.ResponseWriter, *http.Request)
// 	Card   func(http.ResponseWriter, *http.Request)
// }

// var Probabilistics = probabilistic{
// 	Create: probCreate,
// 	Exists: probExists,
// 	Add:    probAdd,
// 	Card:   probCard,
// }
