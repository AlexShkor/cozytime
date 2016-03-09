package models

import (
	"time"
)


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

type GameDto struct {
    Id         string    
	Invited    []string  
	Joined     []string  
	Owner      string    
	TargetTime int       
	IsStarted  bool      
	IsStopped  bool     
	Created    time.Time 
	Started    time.Time
	Ended      time.Time 
	EndedBy    string   
	Deleted    time.Time
	IsDeleted  bool      
    Users []UserDto
}
