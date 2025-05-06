package database

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm/logger"

	// Configs
	cfg "hintword.com/api/app/configs"

	// Gorm
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// MysqlDB is the mysql connection handle
	MysqlDB *gorm.DB
)

func setupGORMCallbacks(db *gorm.DB) error {
	handleConnectionError := func(tx *gorm.DB) {
		if tx.Error != nil && (strings.Contains(tx.Error.Error(), "connect: cannot assign requested address") ||
			strings.Contains(tx.Error.Error(), "connect: connection refused") ||
			strings.Contains(tx.Error.Error(), "invalid connection")) {
			log.Println(strings.Repeat("!", 40))
			log.Println("😔 Database connection lost")
			log.Println(strings.Repeat("!", 40))
			log.Fatal(tx.Error)
		}
	}

	callbacks := []struct {
		name     string
		callback func(name string, fn func(*gorm.DB)) error
	}{
		{"gorm:raw", db.Callback().Raw().After("gorm:raw").Register},
		{"gorm:query", db.Callback().Query().After("gorm:query").Register},
		{"gorm:create", db.Callback().Create().After("gorm:create").Register},
		{"gorm:update", db.Callback().Update().After("gorm:update").Register},
		{"gorm:delete", db.Callback().Delete().After("gorm:delete").Register},
	}

	for _, cb := range callbacks {
		if err := cb.callback("check_connection_error", handleConnectionError); err != nil {
			return fmt.Errorf("failed to register callback for %s: %w", cb.name, err)
		}
	}

	return nil
}

func ConnectMysql() {
	dsn := cfg.GetConfig().Mysql.GetMysqlConnectionInfo()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	db.Set("gorm:auto_preload", true)
	if err != nil {
		log.Println(strings.Repeat("!", 40))
		log.Println("😏 Could Not Establish Mysql DB Connection")
		log.Println(strings.Repeat("!", 40))
		log.Fatal(err)
	}

	//// Set maximum number of open connections
	//sqlDB, err := db.DB()
	//if err != nil {
	//	log.Println(strings.Repeat("!", 40))
	//	log.Println("😏 Set maximum number of open connections")
	//	log.Println(strings.Repeat("!", 40))
	//	log.Fatal(err)
	//}
	//sqlDB.SetMaxOpenConns(10)

	// Set maximum number of idle connections
	//sqlDB.SetMaxIdleConns(5)

	// Setup GORM callbacks
	if err = setupGORMCallbacks(db); err != nil {
		log.Println(strings.Repeat("!", 40))
		log.Println("😏 Failed to setup GORM callbacks")
		log.Println(strings.Repeat("!", 40))
		log.Fatal(err)
	}

	MysqlDB = db

	log.Println(strings.Repeat("-", 40))
	log.Println("😀 Connected To Mysql DB")
	log.Println(strings.Repeat("-", 40))
}

type WhereCondition struct {
	Key            string               `json:"key"`
	Condition      string               `json:"condition"`
	Value          interface{}          `json:"value"`
	GroupCondition *GroupWhereCondition `json:"group_condition"`
	SubQuery       *SubQueryCondition   `json:"sub_query"`
}

type GroupWhereCondition struct {
	Condition string           `json:"condition"`
	Values    []WhereCondition `json:"values"`
}

type SubQueryCondition struct {
	TableName  string           `json:"table_name"`
	Model      interface{}      `json:"model"`
	FieldName  string           `json:"field_name"`
	Conditions []WhereCondition `json:"conditions"`
}

func ConditionBuilder(db *gorm.DB, whereConditions *[]WhereCondition, orConditions *[]WhereCondition, accessConditions *[]WhereCondition) *gorm.DB {
	for _, condition := range *whereConditions {
		if condition.GroupCondition != nil {
			var groupWhereClause string
			var groupArgs []interface{}
			for _, subCondition := range condition.GroupCondition.Values {
				if strings.HasPrefix(subCondition.Condition, "IS") {
					groupWhereClause += fmt.Sprintf("%s %s %s %s ", subCondition.Key, subCondition.Condition, subCondition.Value, condition.GroupCondition.Condition)
				} else {
					groupWhereClause += fmt.Sprintf("%s %s ? %s ", subCondition.Key, subCondition.Condition, condition.GroupCondition.Condition)
					groupArgs = append(groupArgs, subCondition.Value)
				}
			}

			if len(groupWhereClause) > 0 {
				groupWhereClause = strings.TrimSuffix(groupWhereClause, fmt.Sprintf(" %s ", condition.GroupCondition.Condition))
				db = db.Where(groupWhereClause, groupArgs...)
			}
		} else if condition.SubQuery != nil {
			// Process subQuery recursively
			subQuery := buildSubQuery(db, condition.SubQuery)
			// Modify the main query based on the subQuery result
			db = db.Where(fmt.Sprintf("%s %s (%s)", condition.Key, condition.Condition, subQuery))
		} else {
			if strings.HasPrefix(condition.Condition, "IS") {
				db = db.Where(fmt.Sprintf("%s %s %s", condition.Key, condition.Condition, condition.Value))
			} else {
				db = db.Where(fmt.Sprintf("%s %s ?", condition.Key, condition.Condition), condition.Value)
			}
		}
	}

	if orConditions != nil {
		var orWhereClause string
		var orArgs []interface{}
		for _, condition := range *orConditions {
			if condition.SubQuery != nil {
				orWhereClause += fmt.Sprintf("%s %s (%s) OR ",
					condition.Key, condition.Condition, buildSubQuery(db, condition.SubQuery))
			} else if strings.HasPrefix(condition.Condition, "IS") {
				orWhereClause += fmt.Sprintf("%s %s %s OR ", condition.Key, condition.Condition, condition.Value)
			} else {
				orWhereClause += fmt.Sprintf("%s %s ? OR ", condition.Key, condition.Condition)
				orArgs = append(orArgs, condition.Value)
			}
		}

		// Remove the last " AND " or " OR " from orWhereClause
		if len(orWhereClause) > 0 {
			orWhereClause = orWhereClause[:len(orWhereClause)-4]
			db = db.Where(orWhereClause, orArgs...)
		}
	}

	if accessConditions != nil {
		var accessWhereClause string
		var accArgs []interface{}
		for _, condition := range *accessConditions {
			if condition.SubQuery != nil {
				subQuery := buildSubQuery(db, condition.SubQuery)
				accessWhereClause += fmt.Sprintf("%s %s (%s) OR ", condition.Key, condition.Condition, subQuery)
			} else if strings.HasPrefix(condition.Condition, "IS") {
				accessWhereClause += fmt.Sprintf("%s %s %s OR ", condition.Key, condition.Condition, condition.Value)
			} else {
				accessWhereClause += fmt.Sprintf("%s %s ? OR ", condition.Key, condition.Condition)
				accArgs = append(accArgs, condition.Value)
			}
		}

		// Remove the last " AND " or " OR " from accessWhereClause
		if len(accessWhereClause) > 0 {
			accessWhereClause = accessWhereClause[:len(accessWhereClause)-4]
			db = db.Where(accessWhereClause, accArgs...)
		}
	}

	return db
}

func buildSubQuery(db *gorm.DB, subQuery *SubQueryCondition) string {
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tx = tx.Table(subQuery.TableName).Select(subQuery.FieldName)
		//.Where(fmt.Sprintf("%s %s ?", coreModel.UserRoleMapColumns.UserID, "="), userModel.ID)
		tx = ConditionBuilder(tx, &subQuery.Conditions, nil, nil)

		return tx.Find(subQuery.Model)
	})
	return sql
}
