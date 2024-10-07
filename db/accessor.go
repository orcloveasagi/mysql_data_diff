package db

import (
	"db_diff/util"
	"encoding/json"
	"fmt"
	"log"
	"sync/atomic"
)

var allCompareData = make([]*CompareData, 0)
var maxId int64 = 0

const dbName = "db.json"
const dbPath = "db"

func init() {
	file, err := util.ReadFile(dbPath, dbName)
	if err != nil {
		util.WriteLineToFile([]string{"[]"}, dbPath, dbName)
		file, _ = util.ReadFile(dbPath, dbName)
	}
	err = json.Unmarshal(file, &allCompareData)
	if err != nil {
		log.Fatalf("db unmarshal error : %v", err)
	}
	for _, data := range allCompareData {
		if data.Common.Name == "" {
			data.Common.Name = fmt.Sprintf("%s => %s", data.Source.Db, data.Target.Db)
		}
		if data.Id > maxId {
			maxId = data.Id
		}
	}
}

func LoadAll() []*CompareData {
	return allCompareData
}

func Load(id int64) *CompareData {
	for _, data := range allCompareData {
		if data.Id == id {
			return data
		}
	}
	return nil
}

func WriteBack() error {
	marshal, err := json.Marshal(allCompareData)
	if err != nil {
		return err
	}
	err = util.WriteFile(marshal, dbPath, dbName)
	if err != nil {
		return err
	}
	return nil
}

func Delete(id int64) error {
	for i, data := range allCompareData {
		if data.Id == id {
			allCompareData = append(allCompareData[:i], allCompareData[i+1:]...)
			break
		}
	}
	return WriteBack()
}

func Update(new *CompareData) error {
	for i, old := range allCompareData {
		if old.Id == new.Id {
			allCompareData[i] = new
			return WriteBack()
		}
	}
	return fmt.Errorf("cannot update data not exist")
}

func Insert(data *CompareData) error {
	nextId := atomic.AddInt64(&maxId, 1)
	data.Id = nextId
	allCompareData = append(allCompareData, data)
	return WriteBack()
}
