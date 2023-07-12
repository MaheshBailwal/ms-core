package persistence

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/microsoft/go-mssqldb"
)

const schema = "plum"

type fieldMetadData struct {
	name         string
	dbField      string
	wrapInQuotes bool
	isIdentity   bool
	isPK         bool
	value        string
}

func CreateInsertQuery(q interface{}) string {
	fieldsMetaData := getFieldsMetaData(q)
	var fieldsName = ""
	var fieldsValue = ""
	var output = ""

	for _, field := range fieldsMetaData {
		if field.isIdentity {
			output = fmt.Sprintf("Output Inserted.%s", field.dbField)
			continue
		}

		fieldsName = fmt.Sprintf("%s %s,", fieldsName, field.dbField)
		if field.wrapInQuotes {
			fieldsValue = fmt.Sprintf("%s '%s',", fieldsValue, field.value)
		} else {
			fieldsValue = fmt.Sprintf("%s %s,", fieldsValue, field.value)
		}
	}

	fieldsName = strings.Trim(fieldsName, ",")
	fieldsValue = strings.Trim(fieldsValue, ",")
	table := reflect.TypeOf(q).Name()
	query := fmt.Sprintf("insert into [%s].[%s] (%s) %s values(%s)", schema, strings.ToLower(table), fieldsName, output, fieldsValue)
	return query
}

func CreateUpdateQuery(q interface{}) string {
	fieldsMetaData := getFieldsMetaData(q)
	var updateFields = ""
	var where = ""

	for _, field := range fieldsMetaData {
		if field.isPK {
			where = fmt.Sprintf(" where  %s=%s", field.dbField, field.value)
			continue
		}

		if field.wrapInQuotes {
			updateFields = fmt.Sprintf("%s %s='%s',", updateFields, field.dbField, field.value)
		} else {
			updateFields = fmt.Sprintf("%s %s=%s,", updateFields, field.dbField, field.value)
		}

	}

	updateFields = strings.Trim(updateFields, ",")
	table := reflect.TypeOf(q).Name()
	query := fmt.Sprintf("update [%s].[%s] SET %s %s", schema, strings.ToLower(table), updateFields, where)
	return query
}

func CreateReadAllQuery[T any](filters QueryFilter) string {
	var entity T
	table := reflect.ValueOf(&entity).Elem().Type().Name()
	fieldsMetaData := getFieldsMetaData(entity)
	filterField := ""
	sortField := ""

	var fieldsName = ""
	for _, field := range fieldsMetaData {
		fieldsName = fmt.Sprintf("%s %s as %s,", fieldsName, field.dbField, field.name)

		val, ok := filters.Filters[strings.ToLower(field.name)]
		if ok {
			if field.wrapInQuotes {
				filterField = fmt.Sprintf(" %s %s='%s' AND", filterField, field.dbField, val)
			} else {
				filterField = fmt.Sprintf(" %s %s=%s AND", filterField, field.dbField, val)
			}
		}

		val, ok = filters.Sort[strings.ToLower(field.name)]
		if ok {
			sortField = fmt.Sprintf("%s %s %s,", sortField, field.dbField, val)
		}
	}

	fieldsName = strings.Trim(fieldsName, ",")

	if filterField != "" {
		filterField = fmt.Sprintf("where %s", strings.Trim(filterField, "AND"))
	}

	if sortField != "" {
		sortField = fmt.Sprintf("order by %s", strings.Trim(sortField, ","))
	}

	return fmt.Sprintf("select %s from  [%s].[%s] %s %s",
		strings.ToLower(fieldsName), schema, strings.ToLower(table), filterField, sortField)
}

func CreateReadByIdQuery[T any](id int64) string {
	var entity T
	table := reflect.ValueOf(&entity).Elem().Type().Name()
	fieldsMetaData := getFieldsMetaData(entity)
	where := ""
	fieldsName := ""

	for _, field := range fieldsMetaData {
		fieldsName = fmt.Sprintf("%s %s as %s,", fieldsName, field.dbField, field.name)

		if field.isPK {
			where = fmt.Sprintf("where  %s=%d", field.dbField, id)
		}
	}

	fieldsName = strings.Trim(fieldsName, ",")

	return fmt.Sprintf("select %s from  [%s].[%s] %s",
		strings.ToLower(fieldsName), schema, strings.ToLower(table), where)
}

func CreateDeleteQuery[T any](id int64) string {
	var entity T
	table := reflect.ValueOf(&entity).Elem().Type().Name()
	fieldsMetaData := getFieldsMetaData(entity)
	where := ""

	for _, field := range fieldsMetaData {
		if field.isPK {
			where = fmt.Sprintf("where  %s=%d", field.dbField, id)
			break
		}
	}

	return fmt.Sprintf("delete  from  [%s].[%s] %s", schema, strings.ToLower(table), where)
}

func Map[T any](rows sqlx.Rows) []T {
	var entity T
	var entities []T

	for rows.Next() {
		err := rows.StructScan(&entity)
		if err != nil {
			log.Fatalln(err)
		}
		entities = append(entities, entity)
	}
	return entities
}

func getFieldsMetaData(q interface{}) []fieldMetadData {
	//TODO: apply cache
	var fieldsMetaData []fieldMetadData
	if reflect.ValueOf(q).Kind() == reflect.Struct {
		v := reflect.ValueOf(q)

		for i := 0; i < v.NumField(); i++ {
			metaData := fieldMetadData{}
			var stype = reflect.TypeOf(q)
			metaData.value, metaData.wrapInQuotes = getFieldValue(q, i)
			metaData.dbField = strings.ToLower(stype.Field(i).Tag.Get("dbfield"))
			metaData.name = stype.Field(i).Name
			if metaData.dbField == "" {
				metaData.dbField = stype.Field(i).Name
			}

			if stype.Field(i).Tag.Get("identity") == "true" {
				metaData.isIdentity = true
				metaData.isPK = true
			}
			if stype.Field(i).Tag.Get("pk") == "true" {
				metaData.isPK = true
			}

			fieldsMetaData = append(fieldsMetaData, metaData)
		}
	}
	return fieldsMetaData
}

func getFieldValue(q interface{}, indx int) (string, bool) {
	v := reflect.ValueOf(q)
	field := v.Field(indx)
	switch field.Kind() {
	case reflect.Int:
		return fmt.Sprintf("%d", field.Int()), false
	case reflect.Int64:
		return fmt.Sprintf("%d", field.Int()), false
	case reflect.String:
		return field.String(), true
	case reflect.Bool:
		return getFieldValueEx(q, indx)
	case reflect.Struct:
		return getFieldValueEx(q, indx)
	case reflect.Array:
		return getFieldValueEx(q, indx)
	default:
		fmt.Println("Unsupported type", field.Kind().String())
	}

	return "", false
}

func getFieldValueEx(v interface{}, field int) (string, bool) {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).Field(field)
	//fmt.Println("f->", v)
	fieldValue := f.Interface()
	switch v := fieldValue.(type) {
	case bool:
		if v {
			return "1", false
		}
		return "0", false
	case time.Time:
		return v.Format(time.DateTime), true
	case uuid.UUID:
		return v.String(), true
	case sql.NullInt64:
		if v.Valid {
			return fmt.Sprintf("%d", v.Int64), false
		} else {
			return "Null", false
		}
	case sql.NullString:
		if v.Valid {
			return fmt.Sprintf("%s", v.String), true
		} else {
			return "Null", false
		}
	default:
		return "", false
	}
}
