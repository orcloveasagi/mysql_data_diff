package logic

import (
	"database/sql"
	"db_diff/util"
	"fmt"
	"log"
)

// 比较和同步数据库
func structDiff(sourceDB, targetDB *sql.DB) (result []string, err error) {
	sourceTables := Tables(sourceDB)
	targetTables := Tables(targetDB)

	result = append(result, EncodeHead, CloseForeignKeyCheck, EmptyString)
	diffFound := false
	// 检查源库和目标库的表
	for tableName := range sourceTables {
		if _, exists := targetTables[tableName]; !exists {
			// 源库存在，目标库不存在，生成创建表语句
			createTableDDL, err := getCreateTableSQL(sourceDB, tableName)
			if err != nil {
				return result, err
			}
			result = append(result, createTableDDL)
			result = append(result, EmptyString)
			diffFound = true
		}
	}

	for tableName := range targetTables {
		if _, exists := sourceTables[tableName]; !exists {
			// 源库不存在，目标库存在，生成删除表语句
			dropTableDDL := fmt.Sprintf("DROP TABLE IF EXISTS `%s`;", tableName)
			result = append(result, dropTableDDL)
			result = append(result, EmptyString)
			diffFound = true
		} else {
			// 源库和目标库都存在，检查表结构
			modifyTableDDL, err := compareTableStructure(sourceDB, targetDB, tableName)
			if err != nil {
				return result, err
			}
			if len(modifyTableDDL) == 0 {
				continue
			}
			result = append(result, modifyTableDDL...)
			result = append(result, EmptyString)
			diffFound = true
		}
	}

	if !diffFound {
		return nil, nil
	}
	result = append(result, OpenForeignKeyCheck)
	return result, err
}

// 获取数据库中的所有表
func Tables(db *sql.DB) map[string]struct{} {
	rows, err := util.Query(db, "SHOW TABLES")
	if err != nil {
		return nil
	}
	defer util.Close(rows)

	tables := make(map[string]struct{})
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			log.Fatalf("扫描表名时出错: %v", err)
		}
		tables[tableName] = struct{}{}
	}
	return tables
}

// 获取创建表的 SQL 语句
func getCreateTableSQL(db *sql.DB, tableName string) (string, error) {
	row, err := util.Query(db, fmt.Sprintf("SHOW CREATE TABLE `%s`;", tableName))
	if err != nil {
		return "", err
	}
	defer util.Close(row)

	var createSQL string
	row.Next()
	if err := row.Scan(&tableName, &createSQL); err != nil {
		log.Printf("获取创建表语句时出错: %v", err)
		return "", err
	}
	return createSQL + ";", nil
}

// 比较表结构
func compareTableStructure(sourceDB, targetDB *sql.DB, tableName string) (result []string, err error) {
	// 这里你可以添加更多的逻辑来比较字段类型、约束等
	sourceColumns, err := GetColumns(sourceDB, tableName)
	if err != nil {
		return result, err
	}
	targetColumns, err := GetColumns(targetDB, tableName)
	if err != nil {
		return result, err
	}

	// 比较字段
	for colName, colType := range sourceColumns {
		if targetType, exists := targetColumns[colName]; !exists {
			// 生成添加字段的 SQL 语句
			addColumnSQL := fmt.Sprintf("ALTER TABLE `%s` ADD COLUMN `%s` %s;", tableName, colName, colType)
			result = append(result, addColumnSQL)
		} else if colType != targetType {
			// 生成修改字段的 SQL 语句
			modifyColumnSQL := fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN `%s` %s;", tableName, colName, colType)
			result = append(result, modifyColumnSQL)
		}
	}

	for colName := range targetColumns {
		if _, exists := sourceColumns[colName]; !exists {
			// 生成删除字段的 SQL 语句
			dropColumnSQL := fmt.Sprintf("ALTER TABLE `%s` DROP COLUMN `%s`;", tableName, colName)
			result = append(result, dropColumnSQL)
		}
	}

	return result, err
}

// 获取表的所有列
func GetColumns(db *sql.DB, tableName string) (map[string]string, error) {
	rows, err := util.Query(db, fmt.Sprintf("SHOW COLUMNS FROM `%s`", tableName))
	if err != nil {
		return nil, err
	}
	defer util.Close(rows)

	columns := make(map[string]string)
	for rows.Next() {
		var field, colType, null, key, extra string
		var defaultVal sql.NullString
		if err := rows.Scan(&field, &colType, &null, &key, &defaultVal, &extra); err != nil {
			log.Fatalf("扫描列信息时出错: %v", err)
		}
		columns[field] = colType
	}
	return columns, nil
}
