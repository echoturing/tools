package http

import (
	"context"
	"testing"

	"github.com/echoturing/log"
)

func TestGet(t *testing.T) {
	resp := make(map[string]interface{})
	err := Get(context.Background(), "http://test-dmp.yyuehd.com/api/v1/options/app/cat", &resp, nil)
	if err != nil {
		t.Error(err)
		return
	}
	log.Debug("resp", "resp", resp)
}
