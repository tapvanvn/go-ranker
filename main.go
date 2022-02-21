package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/tapvanvn/go-ranker/goranker"
	engines "github.com/tapvanvn/godbengine"
	"github.com/tapvanvn/godbengine/engine"
	"github.com/tapvanvn/godbengine/engine/adapter"
	"github.com/tapvanvn/goranking"
	"github.com/tapvanvn/goutil"
)

var rankingSystem *goranking.RankingSystem = nil
var config *goranker.Config = nil

func main() {

	port := goutil.MustGetEnv("PORT")

	rootPath, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	configPath := getGeneralConfigPath(rootPath)

	jsonFile, err := os.Open(configPath)

	if err != nil {

		panic(err)
	}
	bytes, err := ioutil.ReadAll(jsonFile)

	if err != nil {

		panic(err)
	}
	bytes = goutil.TripJSONComment(bytes)

	config = &goranker.Config{}

	if err != nil {

		panic(err)
	}

	err = json.Unmarshal(bytes, config)

	if err != nil {

		panic(err)
	}

	goranker.SetConfig(config)

	engines.InitEngineFunc = startEngine
	_ = engines.GetEngine()

	rankingSystem = goranking.NewRankingSystem(config.TableSize)
	var pageSize = 1000
	var page = 0
	for {
		items, _ := goranker.GetRecordPage(page, pageSize)

		if len(items) == 0 {

			break
		}
		for _, item := range items {

			rankingSystem.PutScore(item.UserID, 0, uint64(item.Score))
			goranker.PutLastScore(item.UserID, item.Score)
		}
		page++
	}

	http.HandleFunc("/score/", func(w http.ResponseWriter, r *http.Request) {

		userIDString := strings.TrimSpace(r.URL.Path[7:])

		if r.Method == "POST" {

			frm := &goranker.FormPostScore{}

			if err := goutil.FromRequest(frm, r); err != nil {

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
			last := uint64(0)

			if testLast, err := goranker.GetLastScore(userIDString); err == nil {
				last = uint64(testLast)
			}

			rank := rankingSystem.PutScore(userIDString, last, uint64(frm.Score))

			record := &goranker.Record{
				UserID: userIDString,
				Score:  frm.Score,
				Rank:   rank,
			}

			go goranker.PutRecord(record)
			goranker.PutLastScore(userIDString, frm.Score)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte(strconv.FormatUint(uint64(rank), 10)))

			return

		} else if r.Method == "GET" {

			last, err := goranker.GetLastScore(userIDString)

			if err == nil {
				fmt.Printf("get last:%d user:%s\n", last, userIDString)
				rank := rankingSystem.GetScore(userIDString, uint64(last))
				rankingSystem.PrintDebug()
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(strconv.FormatUint(uint64(rank), 10)))
				return
			} else {
				w.WriteHeader(http.StatusNotFound)
				return
			}
		}
	})

	fmt.Println("run on port:", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func getGeneralConfigPath(rootPath string) string {

	configFile := goutil.GetEnv("CONFIG")

	if configFile == "" {

		configFile = "config.jsonc"
	}
	configPath := rootPath + "/config/" + configFile
	return configPath
}

//Start start engine
func startEngine(eng *engine.Engine) {

	//read redis define from env
	var memdb engine.MemPool = nil

	if config.MemDB != nil {

		if config.MemDB.Type == "redis" {

			redisConnectString := config.MemDB.ConnectionString
			fmt.Println("redis:", redisConnectString)
			redisPool := adapter.RedisPool{}

			err := redisPool.Init(redisConnectString)

			if err != nil {

				fmt.Println("cannot init redis")
			}
			memdb = &redisPool

		} else if config.MemDB.Type == "rediscluster" {
			redisConnectString := config.MemDB.ConnectionString
			fmt.Println("redis:", redisConnectString)
			redisPool := adapter.RedisClusterPool{}

			err := redisPool.Init(redisConnectString)

			if err != nil {

				fmt.Println("cannot init redis")
			}
			memdb = &redisPool
		}
	}
	var documentDB engine.DocumentPool = nil

	if config.DocumentDB != nil {

		connectString := config.DocumentDB.ConnectionString
		databaseName := config.DocumentDB.DatabaseName

		if config.DocumentDB.Provider == "mongodb" {

			mongoPool := &adapter.MongoPool{}
			err := mongoPool.InitWithDatabase(connectString, databaseName)

			if err != nil {

				log.Fatal("cannot init mongo")
			}
			documentDB = mongoPool

		} else {

			firestorePool := adapter.FirestorePool{}
			firestorePool.Init(connectString)
			documentDB = &firestorePool
		}
	}

	eng.Init(memdb, documentDB, nil)
}
