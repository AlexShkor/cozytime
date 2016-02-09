package models

type CreateGame struct {
	Players    []string
	TargetTime int
}

type JoinGame struct {
	GameId string
}

type GameResponse struct {
	GameId string
}
