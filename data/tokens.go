package data

import (
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type TokenDocument struct {
	Id          string    `bson:"_id,omitempty"`
	PhoneNumber string    `bson:"PhoneNumber"`
	Token       string    `bson:"Token"`
    Code       string    `bson:"Code"`
	DateTime    time.Time `bson:"DateTime"`
}

type Tokens struct {
	collection *mgo.Collection
}

func NewTokensService(collection *mgo.Collection) *Tokens {
	return &Tokens{collection}
}

func (t *Tokens) Authorize(phoneNumber string) (string, error) {
	id := GenerateId()
	token := GenerateAccessToken()

	doc := TokenDocument{id, phoneNumber, token, "", time.Now()}
	err := t.collection.Insert(doc)

	return token, err
}

func (t *Tokens) IsAuthorized(token string) (string, error) {
	var doc TokenDocument
	err := t.collection.Find(bson.M{"Token": token}).One(&doc)
	if err != nil {
		return "", err
	}
	return doc.PhoneNumber, err
}

func (t *Tokens) FindByPhone(phone string) (*TokenDocument, error) {
	var doc *TokenDocument
	err := t.collection.Find(bson.M{"Phone": phone}).One(doc)
	if err != nil {
		return nil, err
	}
	return doc, err
}

func (t *Tokens) GetAll() ([]TokenDocument, error) {
	var docs []TokenDocument
	err := t.collection.Find(bson.M{}).All(&docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}
