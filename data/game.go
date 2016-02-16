package data

import (
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
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
	Deleted    time.Time `bson:"Deleted"`
	IsDeleted  bool      `bson:"IsDeleted"`
}

type Games struct {
	collection *mgo.Collection
}

func NewGamesService(collection *mgo.Collection) *Games {
	return &Games{collection}
}

func (t *Games) Create(ownerId string, players []string, targetTime int) (*GameDocument, error) {
	id := GenerateId()
	doc := GameDocument{id, append(players, ownerId), []string{ownerId}, ownerId, targetTime, false, false, time.Now(), time.Time{}, time.Time{}, "", time.Time{}, false}
	err := t.collection.Insert(doc)
	return &doc, err
}

func (t *Games) Join(gameId string, userId string) error {
	err := t.collection.Update(bson.M{"$and": []bson.M{bson.M{"_id": bson.M{"$eq": gameId}}, bson.M{"Invited": bson.M{"$eq": userId}}, bson.M{"Joined": bson.M{"$ne": userId}}}}, bson.M{"$push": bson.M{"Joined": userId}})
	return err
}

func (t *Games) Leave(gameId string, userId string) error {
	err := t.collection.Update(bson.M{"$and": []bson.M{bson.M{"_id": bson.M{"$eq": gameId}}, bson.M{"Invited": bson.M{"$eq": userId}}, bson.M{"Joined": bson.M{"$eq": userId}}}}, bson.M{"$pull": bson.M{"Joined": userId}})
	return err
}

func (t *Games) Start(gameId string) (bool, error) {
	doc, err := t.Get(gameId)
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

func (t *Games) Get(gameId string) (*GameDocument, error) {
	var doc GameDocument
	err := t.collection.FindId(gameId).One(&doc)
	return &doc, err
}

func (t *Games) Delete(gameId string, userId string) error {
	err := t.collection.Update(bson.M{"$and": []bson.M{bson.M{"_id": bson.M{"$eq": gameId}}, bson.M{"IsStarted": bson.M{"$eq": false}}, bson.M{"OwnerId": bson.M{"$eq": userId}}}}, bson.M{"$set": bson.M{"Deleted": time.Now(), "IsDeleted": true}})
	return err
}

func (t *Games) GetAllForUser(userId string) ([]GameDocument, error) {
	var docs []GameDocument
	err := t.collection.Find(bson.M{"$or": []bson.M{bson.M{"OwnerId": bson.M{"$eq": userId}}, bson.M{"Invited": bson.M{"$eq": userId}}}}).All(&docs)
	return docs, err
}
