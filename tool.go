package tools

import (
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"reflect"
	"time"

	"fmt"

	"github.com/labstack/gommon/log"
)

//return the i is nil or not
func IsNil(object interface{}) bool {
	if object == nil {
		return true
	}
	value := reflect.ValueOf(object)
	kind := value.Kind()
	value.Interface()
	if kind >= reflect.Chan && kind <= reflect.Slice && value.IsNil() {
		return true
	}
	return false

}

var client *http.Client = &http.Client{Timeout: 10 * time.Second}

//http get to download a file and concurrent upload
func GetDownLoadAndPostUpload(downLoadUrl, uploadUrl, formField string) ([]byte, error) {
	var err error
	response, err := client.Get(downLoadUrl)
	if err != nil {
		return nil, err
	}
	r, w := io.Pipe()
	mpw := multipart.NewWriter(w)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				var err error
				switch r := r.(type) {
				case error:
					err = r
				default:
					err = fmt.Errorf("%v", r)
				}
				log.Error(err)
			}
		}()
		var part io.Writer
		defer w.Close()
		defer response.Body.Close()
		if part, err = mpw.CreateFormFile(formField, "filename"); err != nil {
			log.Fatal(err)
		}
		if _, err = io.Copy(part, response.Body); err != nil {
			log.Fatal(err)
		}
		if err = mpw.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	newResponse, err := http.Post(uploadUrl, mpw.FormDataContentType(), r)
	if err != nil {
		log.Fatal(err)
	}
	defer newResponse.Body.Close()
	responseBytes, err := ioutil.ReadAll(newResponse.Body)
	if err != nil {
		log.Fatal(err)
	}
	return responseBytes, nil
}
