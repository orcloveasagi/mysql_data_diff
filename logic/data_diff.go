package logic

import (
	"database/sql"
	"db_diff/util"
	"fmt"
	"log"
	"strings"
)

type table struct {
	rows    map[string]map[string]string
	columns []string
	pkName  string
}

func dataDiff(sourceDB, targetDB *sql.DB) (result []string, err error) {
	tables := Tables(sourceDB)

	result = append(result, EncodeHead, CloseForeignKeyCheck, EmptyString)
	diffFound := false
	for tableName := range tables {
		diff, err := diffOneTable(sourceDB, targetDB, tableName)
		if err != nil {
			return nil, err
		}
		if len(diff) == 0 {
			continue
		}
		result = append(result, diff...)
		result = append(result, EmptyString)
		diffFound = true
	}
	if !diffFound {
		return nil, nil
	}
	result = append(result, OpenForeignKeyCheck)
	return result, nil
}

// 逐行比较表数据
func diffOneTable(sourceDB, targetDB *sql.DB, tableName string) (result []string, err error) {
	sourceTable, err := readTable(sourceDB, tableName)
	if err != nil {
		return nil, err
	}
	targetTable, err := readTable(targetDB, tableName)
	if err != nil {
		return nil, err
	}

	for pk, sourceRow := range sourceTable.rows {
		targetRow, exists := targetTable.rows[pk]
		if !exists {
			insertSQL := generateInsertSQL(tableName, sourceRow, sourceTable.columns)
			result = append(result, insertSQL)
			continue
		}
		diff := diffFields(sourceRow, targetRow)
		if len(diff) > 0 {
			updateSQL := generateUpdateSQL(tableName, diff, sourceTable.pkName, pk, sourceTable.columns)
			result = append(result, updateSQL)
			continue
		}
	}

	for pk := range targetTable.rows {
		_, exists := sourceTable.rows[pk]
		if !exists {
			deleteSQL := generateDeleteSQL(tableName, sourceTable.pkName, pk)
			result = append(result, deleteSQL)
		}
	}
	return result, nil
}

func readTable(db *sql.DB, tbName string) (table table, err error) {
	var primaryKeyName string
	primaryKeyRow, err := util.Query(db, fmt.Sprintf(`
		SELECT COLUMN_NAME
		FROM information_schema.KEY_COLUMN_USAGE
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = '%s' AND CONSTRAINT_NAME = 'PRIMARY'`, tbName))
	if err != nil {
		return table, err
	}
	defer util.Close(primaryKeyRow)
	if !primaryKeyRow.Next() { //table not found
		return table, err
	}
	scanErr := primaryKeyRow.Scan(&primaryKeyName)
	if scanErr != nil {
		log.Printf("scan error: %v", scanErr)
		return table, scanErr
	}

	dataRows, err := util.Query(db, fmt.Sprintf("SELECT * FROM `%s`", tbName))
	if err != nil {
		return table, err
	}
	defer util.Close(dataRows)

	colTypes, err := GetColumns(db, tbName)
	if err != nil {
		return table, err
	}
	// 获取列名
	sourceColumns, err := dataRows.Columns()
	if err != nil {
		log.Fatalf("获取数据表列名数组时出错: %v", err)
	}

	// 逐行读取数据
	rowsMapByPrimaryKey := make(map[string]map[string]string)
	for dataRows.Next() {
		row := make([]any, len(sourceColumns))
		for i := range row {
			row[i] = new(any)
		}
		if err := dataRows.Scan(row...); err != nil {
			log.Printf("scan table data error %v", err)
			return table, err
		}
		rowMap := make(map[string]string)
		for i, value := range row {
			colName := sourceColumns[i]
			colType := colTypes[colName]
			rowMap[colName] = util.MapMySQL(value, colType)
		}
		keyValue := rowMap[primaryKeyName]
		rowsMapByPrimaryKey[keyValue] = rowMap
	}
	table.rows = rowsMapByPrimaryKey
	table.pkName = primaryKeyName
	table.columns = sourceColumns
	return table, nil
}

// 比较两行数据是否相等
func diffFields(sourceRow, targetRow map[string]string) map[string]string {
	diff := make(map[string]string)
	for name, sourceValue := range sourceRow {
		targetValue, exists := targetRow[name]
		if !exists {
			diff[name] = sourceValue
			continue
		}
		if targetValue != sourceValue {
			diff[name] = sourceValue
			continue
		}
	}
	return diff
}

// 生成插入数据的 SQL 语句
func generateInsertSQL(tableName string, row map[string]string, columns []string) string {
	setClause := ""
	for _, name := range columns {
		value, exists := row[name]
		if !exists {
			continue
		}
		if value == "NULL" {
			continue
		}
		setClause += fmt.Sprintf("`%s`=%s,", name, value)
	}
	setClause = strings.TrimSuffix(setClause, ",")

	return fmt.Sprintf("INSERT INTO `%s` SET %s;", tableName, setClause)
}

// 生成删除数据的 SQL 语句
func generateDeleteSQL(tableName string, pkName string, pkValue string) string {
	return fmt.Sprintf("DELETE FROM `%s` WHERE `%s` = %s;", tableName, pkName, pkValue)
}

// 生成更新数据的 SQL 语句
func generateUpdateSQL(tableName string, diff map[string]string, pkName string, pkValue any, columns []string) string {
	whereClause := fmt.Sprintf("`%s`='%v'", pkName, pkValue)
	setClause := ""
	for _, name := range columns {
		value, exists := diff[name]
		if !exists {
			continue
		}
		setClause += fmt.Sprintf("`%s`=%s,", name, value)
	}
	setClause = strings.TrimSuffix(setClause, ",")

	return fmt.Sprintf("UPDATE `%s` SET %s WHERE %s;", tableName, setClause, whereClause)
}
