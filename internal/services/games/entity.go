package games

import "github.com/google/uuid"

const (
	ErrMsgWrongDataGame = "Sorry, wrong data game"
)

type TypeGame string

const (
	TypeGameUnknown TypeGame = "unknown"
	TypeGameRevise  TypeGame = "revise"
)

func ToTypeGame(s string) TypeGame {
	switch s {
	case "revise":
		return TypeGameRevise
	default:
		return TypeGameUnknown
	}
}

type TypeStep string

const (
	StepChooseWord   TypeStep = "choose_word"
	StepCompleteWord TypeStep = "complete_word"
	StepAudio        TypeStep = "audio"
)

type (
	Game struct {
		TypeGame  TypeGame
		VocabID   uuid.UUID
		CountWord int
	}

	ReviseGameWord struct {
		Text       string
		Translates []string
		Examples   []string
		Right      int
		Wrong      int
	}

	ReviseGameData struct {
		VocabID uuid.UUID
		Steps   []struct {
			Type     string
			RightAns string
			WrongAns []string
			Chars    []string
		}
	}
)
