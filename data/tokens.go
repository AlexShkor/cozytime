package data

import (
	"time"
   "strconv"
     "math/rand"
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

func (t *Tokens) Authorize(id string, code string) (string, error) {
	token := GenerateAccessToken()
    var doc TokenDocument
    err := t.collection.Find(bson.M{"_id": id}).One(&doc)
	if err != nil {
		return "", err
	}
    if code == doc.Code {
       updateError := t.collection.Update(bson.M{"_id": id}, bson.M{"Token": token})
       if updateError != nil {
		return "", err
	}
    }
	return token, err
}

func (t *Tokens) Create(phoneNumber string) (string, error) {
	id := GenerateId()
    code := strconv.Itoa(rand.Intn(999999))
	doc := TokenDocument{id, phoneNumber, "", code, time.Now()}
	err := t.collection.Insert(doc)

	return code, err
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

func (t *Tokens) FindFriends(phones []string) ([]TokenDocument, error) {
	var docs []TokenDocument
	err := t.collection.Find(bson.M{"$in": bson.M{"Phone": phones}}).All(&docs)
	if err != nil {
		return nil, err
	}
	return docs, err
}

func (t *Tokens) GetAll() ([]TokenDocument, error) {
	var docs []TokenDocument
	err := t.collection.Find(bson.M{}).All(&docs)
	if err != nil {
		return nil, err
	}
	return docs, nil
}
