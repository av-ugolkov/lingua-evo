package handler

import (
	entity "github.com/av-ugolkov/lingua-evo/internal/services/games"
)

func reviseGameToEntity(rs ReviseGameRq) entity.Game {
	return entity.Game{
		VocabID:   rs.VocabID,
		CountWord: rs.CountWord,
		TypeGame:  entity.ToTypeGame(rs.TypeGame),
	}
}
