package logic

import (
	"database/sql"
	"db_diff/db"
	"db_diff/util"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func DatabaseDiff(data *db.CompareData) (dir string, err error) {
	sourceData := data.Source
	targetData := data.Target
	sourceDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", sourceData.User, sourceData.Psw, sourceData.Host, sourceData.Port, sourceData.Db)
	targetDSN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", targetData.User, targetData.Psw, targetData.Host, targetData.Port, targetData.Db)

	timeStamp := util.MilliFormat(time.Now().UnixMilli(), util.YyyyMMddhhmmss)

	sourceDB, err := sql.Open("mysql", sourceDSN)
	if err != nil {
		log.Printf("db open error: %v", err)
		return dir, err
	}
	defer util.Close(sourceDB)

	targetDB, err := sql.Open("mysql", targetDSN)
	if err != nil {
		log.Printf("db open error: %v", err)
		return dir, err
	}
	defer util.Close(targetDB)

	ddlDiff, err := structDiff(sourceDB, targetDB)
	if err != nil {
		log.Printf("struct diff error: %v", err)
		return dir, err
	}
	if ddlDiff != nil {
		dir, err = util.WriteLineToFile(ddlDiff, data.Common.Path, fmt.Sprintf("V%s__%s.sql", timeStamp, data.Common.Ddl))
		if err != nil {
			log.Printf("write file error: %v", err)
			return dir, err
		}
	}

	dmlDiff, err := dataDiff(sourceDB, targetDB)
	if err != nil {
		log.Printf("data diff error: %v", err)
		return dir, err
	}
	if dmlDiff != nil {
		dir, err = util.WriteLineToFile(dmlDiff, data.Common.Path, fmt.Sprintf("V%s__%s.sql", timeStamp, data.Common.Dml))
		if err != nil {
			log.Printf("write file error: %v", err)
			return dir, err
		}
	}
	return dir, err
}
