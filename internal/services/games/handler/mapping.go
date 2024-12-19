package handler

import (
	entity "github.com/av-ugolkov/lingua-evo/internal/services/games"
)

func reviseGameFromRsToEntity(rs ReviseGameRq) entity.ReviseGame {
	return entity.ReviseGame{
		VocabID:   rs.VocabID,
		CountWord: rs.CountWord,
		TypeGame:  rs.TypeGame,
	}
}
