package access

type Type uint8

const (
	Private     Type = 0
	Subscribers Type = 1
	Public      Type = 2
)

type Status uint8

const (
	Forbidden Status = 0
	Read      Status = 1
	Edit      Status = 2
)
