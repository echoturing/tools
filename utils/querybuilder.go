package utils

import (
	"bytes"
	"fmt"
	"reflect"
	"sync"
	"text/template"
	"time"
)

var (
	templateOnce           = sync.Once{}
	insertOrUpdateTemplate *template.Template
)

func NewInsertOrUpdateBody(table string, columns, duplicateUpdateColumns []string, valueStruct interface{}) *insertOrUpdateBody {
	return &insertOrUpdateBody{
		Table:                  table,
		Columns:                columns,
		DuplicateUpdateColumns: duplicateUpdateColumns,
		ValueStruct:            valueStruct,
	}
}

type insertOrUpdateBody struct {
	Table                  string      `json:"table"`
	Columns                []string    `json:"columns"`
	DuplicateUpdateColumns []string    `json:"duplicateUpdateColumns"`
	ValueStruct            interface{} `json:"valueStruct"`
}

func (i *insertOrUpdateBody) formatTemplateValues() map[string]interface{} {
	escape := WrapWith(i.Columns, "`")
	placeholders := make([]string, 0, len(i.Columns))
	updateStmts := make([]string, 0, len(i.Columns))
	for range i.Columns {
		placeholders = append(placeholders, "?")

	}
	for _, column := range i.DuplicateUpdateColumns {
		updateStmts = append(updateStmts, fmt.Sprintf("`%s`=values(`%s`)", column, column))
	}
	return map[string]interface{}{
		"table":        i.Table,
		"columns":      StringJoin(escape, ","),
		"placeholders": StringJoin(placeholders, ","),
		"updateStmt":   StringJoin(updateStmts, ","),
	}
}

func getInsertOrUpdateTemplate() *template.Template {
	templateOnce.Do(func() {
		tmp := template.New("insertOrUpdate")
		insertOrUpdateTemplate, _ = tmp.Parse(
			"insert into `{{.table}}` " +
				" ({{.columns}}) " +
				" values " +
				" ({{.placeholders}}) " +
				" on duplicate key update " +
				" {{.updateStmt}}",
		)
	})
	return insertOrUpdateTemplate
}

// getDBTagSimpleValues 根据dbKeys获取data对应的值，并且把值按简单类型(数字，字符串，布尔)列表返回
func getDBTagSimpleValues(data interface{}, dbKeys []string) []interface{} {
	v := reflect.ValueOf(data)
	t := v.Kind()
	var s reflect.Value
	if t == reflect.Ptr {
		s = v.Elem()
	} else if t == reflect.Struct {
		s = v
	}
	typeOfs := s.Type()
	midMap := map[string]interface{}{}
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		interValue := f.Interface()
		switch interValue.(type) {
		case time.Time:
			interValue = interValue.(time.Time).Format(time.RFC3339)
		}
		midMap[typeOfs.Field(i).Tag.Get("db")] = interValue
	}
	result := make([]interface{}, 0, len(dbKeys))
	for _, key := range dbKeys {
		if value, ok := midMap[key]; ok {
			result = append(result, value)
		}
	}
	return result
}

func InsertOrUpdate(body *insertOrUpdateBody) (queryStr string, values []interface{}, err error) {
	result := new(bytes.Buffer)
	err = getInsertOrUpdateTemplate().Execute(result, body.formatTemplateValues())
	if err != nil {
		return "", nil, err
	}
	values = getDBTagSimpleValues(body.ValueStruct, body.Columns)
	return result.String(), values, nil
}
