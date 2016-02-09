package data

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"time"
)

type GameDocument struct {
	Id         string    `bson:"_id,omitempty"`
	Invited    []string  `bson:"Invited"`
	Joined     []string  `bson:"Joined"`
	Owner      string    `bson:"Owner"`
	TargetTime int       `bson:"TargetTime"`
	IsStarted  bool      `bson:"IsStarted"`
	IsStopped  bool      `bson:"IsStopped"`
	Created    time.Time `bson:"Created"`
	Started    time.Time `bson:"Started"`
	Ended      time.Time `bson:"Ended"`
	EndedBy    string    `bson:"EndedBy"`
}

type Games struct {
	collection *mgo.Collection
}

func NewGamesService(collection *mgo.Collection) *Games {
	return &Games{collection}
}

func (t *Games) Create(ownerId string, players []string, targetTime int) (*GameDocument, error) {
	id := GenerateId()
	doc := GameDocument{id, players, []string{ownerId}, ownerId, targetTime, false, false, time.Now(), time.Time{}, time.Time{}, ""}
	err := t.collection.Insert(doc)
	return &doc, err
}

func (t *Games) Join(gameId string, userId string) error {
	err := t.collection.Update(bson.M{"$and": bson.M{"_id": bson.M{"$eq": gameId}, "Invited": bson.M{"$eq": userId}, "Joined": bson.M{"$ne": userId}}}, bson.M{"$push": bson.M{"Joined": userId}})
	return err
}

func (t *Games) Start(gameId string) (bool, error) {
	var doc *GameDocument
	err := t.collection.FindId(gameId).One(doc)
	if err == nil && doc != nil && len(doc.Joined) == len(doc.Invited) {
		err := t.collection.Update(bson.M{"_id": gameId}, bson.M{"$set": bson.M{"Started": time.Now(), "IsStarted": true}})
		return err == nil, err
	}
	return false, err
}

func (t *Games) Stop(gameId string, userId string) error {
	err := t.collection.Update(bson.M{"_id": gameId}, bson.M{"$set": bson.M{"Ended": time.Now(), "EndedBy": userId, "IsStopped": true}})
	return err
}
