package goranker

import "github.com/tapvanvn/goranking"

type Record struct {
	UserID string         `json:"UserID" bson:"UserID"`
	Score  int64          `json:"Score" bson:"Score"`
	Rank   goranking.Rank `json:"Rank" bson:"Rank"`
}

func (doc *Record) GetID() string {
	return doc.UserID
}
