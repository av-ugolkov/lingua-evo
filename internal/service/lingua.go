package service

type Lingua struct {
	ws *WordsService
}

func NewLinguaService(ws *WordsService) *Lingua {
	return &Lingua{
		ws: ws,
	}
}
