package goranker

import (
	"fmt"

	engines "github.com/tapvanvn/godbengine"
	"github.com/tapvanvn/godbengine/engine"
)

func getMemKeyOfUser(userID string) string {

	return fmt.Sprintf("%s_%s", config.MemDB.Prefix, userID)
}

func GetLastScore(userID string) (int64, error) {

	eng := engines.GetEngine()
	memDB := eng.GetMemPool()
	return memDB.GetInt(getMemKeyOfUser(userID))
}

func PutLastScore(userID string, score int64) error {

	eng := engines.GetEngine()
	memDB := eng.GetMemPool()
	return memDB.SetInt(getMemKeyOfUser(userID), score)
}

func PutRecord(record *Record) error {

	eng := engines.GetEngine()
	docDB := eng.GetDocumentPool()
	return docDB.Put(config.DocumentDB.CollectionName, record)
}

func GetRecordPage(page int, pageSize int) ([]*Record, error) {

	eng := engines.GetEngine()
	docDB := eng.GetDocumentPool()
	query := engine.MakeDBQuery(config.DocumentDB.CollectionName, false)
	query.Sort("Rank", false)
	query.Paging(page, pageSize)
	result := docDB.Query(query)

	defer result.Close()

	if result.Count() == 0 {

		return nil, nil

	} else if result.Error() != nil {

		return nil, result.Error()
	}
	documents := make([]*Record, 0)
	for {
		document := &Record{}
		if err := result.Next(document); err != nil {
			break
		}
		documents = append(documents, document)
	}
	return documents, nil
}
