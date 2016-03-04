package data

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

type TokenDocument struct {
	Id          string    `bson:"_id,omitempty"`
	PhoneNumber string    `bson:"PhoneNumber"`
	Token       string    `bson:"Token"`
	Code        string    `bson:"Code"`
	Name        string    `bson:"Name"`
	DateTime    time.Time `bson:"DateTime"`
}

type Tokens struct {
	collection *mgo.Collection
}

func NewTokensService(collection *mgo.Collection) *Tokens {
	return &Tokens{collection}
}

func (t *Tokens) Authorize(phone string, code string) (*TokenDocument, error) {

	token := GenerateAccessToken()
	doc, err := t.FindByPhone(phone)
	if err != nil {
		return nil, err
	}
	if code == doc.Code {
		updateError := t.collection.Update(bson.M{"_id": doc.Id}, bson.M{"$set": bson.M{"Token": token, "Code": ""}})
		if updateError != nil {
			return nil, err
		}
	}
	doc.Token = token
	return doc, err
}

func (t *Tokens) Create(phoneNumber string) (string, error) {
	id := GenerateId()
	code := GenerateCode()
	doc := TokenDocument{id, phoneNumber, "", code, "", time.Now()}
	err := t.collection.Insert(doc)
	return code, err
}

func (t *Tokens) UpdateCode(id string, phoneNumber string) (string, error) {
	code := GenerateCode()
	updateError := t.collection.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"Code": code}})
	return code, updateError
}

func (t *Tokens) SetName(id string, name string) error {
	updateError := t.collection.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"Name": name}})
	return updateError
}

func GenerateCode() string {
	return strconv.Itoa(rand.Intn(999999))
}

func (t *Tokens) IsAuthorized(token string) (string, error) {
	var doc TokenDocument
	err := t.collection.Find(bson.M{"Token": token}).One(&doc)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	fmt.Println(doc)
	fmt.Println("USER ID:")
	fmt.Println(doc.Id)
	return doc.Id, nil
}

func (t *Tokens) FindByPhone(phone string) (*TokenDocument, error) {
	var doc TokenDocument
	err := t.collection.Find(bson.M{"PhoneNumber": phone}).One(&doc)
	if err != nil {
		return nil, err
	}
	return &doc, err
}


func (t *Tokens) FindByToken(token string) (*TokenDocument, error) {
	var doc TokenDocument
	err := t.collection.Find(bson.M{"Token": token}).One(&doc)
	if err != nil {
		return nil, err
	}
	return &doc, err
}

func (t *Tokens) FindFriends(phones []string) ([]TokenDocument, error) {
	var docs []TokenDocument
	err := t.collection.Find(bson.M{"PhoneNumber": bson.M{"$in": phones}}).All(&docs)
	if err != nil {
		return nil, err
	}
	return docs, err
}

func (t *Tokens) GetUsers(ids []string) ([]TokenDocument, error) {
	var docs []TokenDocument
	err := t.collection.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&docs)
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

func (t *Tokens) Get(userId string) (*TokenDocument, error) {
	var doc TokenDocument
	err := t.collection.FindId(userId).One(&doc)
	return &doc, err
}
