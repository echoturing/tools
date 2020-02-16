package utils

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
)

var (
	templateOnce           = sync.Once{}
	insertOrUpdateTemplate *template.Template
)

type InsertOrUpdateBody struct {
	Table                  string   `json:"table"`
	Columns                []string `json:"columns"`
	DuplicateUpdateColumns []string `json:"duplicateUpdateColumns"`
}

func (i *InsertOrUpdateBody) formatTemplateValues() map[string]interface{} {
	escape := WrapWith(i.Columns, "`")
	placeholders := make([]string, 0, len(i.Columns))
	updateStmts := make([]string, 0, len(i.Columns))
	for _, column := range i.Columns {
		placeholders = append(placeholders, "?")
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

func InsertOrUpdate(body *InsertOrUpdateBody) (string, error) {
	result := new(bytes.Buffer)
	err := getInsertOrUpdateTemplate().Execute(result, body.formatTemplateValues())
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
