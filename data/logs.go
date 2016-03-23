package data

import (
	"time"

	"labix.org/v2/mgo"
)

type LogDocument struct {
	Id          string    `bson:"_id,omitempty"`
	Client string    `bson:"Client"`
	User       string    `bson:"User"`
	Level        string    `bson:"Level"`
	Message        string    `bson:"Message"`
	ServerDateTime    time.Time `bson:"DateTime"`
    ClientDateTime    time.Time `bson:"DateTime"`
}

type Logs struct {
	collection *mgo.Collection
}

func NewLogsService(collection *mgo.Collection) *Logs {
	return &Logs{collection}
}


func (t *Tokens) Save(doc *LogDocument)  error {
	doc.Id = GenerateId()
	err := t.collection.Insert(doc)
	return err
}