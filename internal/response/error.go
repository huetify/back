package response

import (
	"encoding/json"
	"fmt"
	"github.com/ermos/chalk"
	"github.com/huetify/back/internal/dbm"
	"net/http"
	"os"
	"reflect"
)

type errorContent struct {
	Type	int 	`json:"type"`
	Message string 	`json:"message"`
}

var errorMask = "an error has occured"

// Print JSON error from type error
func Error (db *dbm.Instance, w http.ResponseWriter, status int, errID int, err interface{}) bool{
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if db != nil {
		_ = db.Rollback()
		_ = db.Close()
	}

	var errorText string
	t := reflect.ValueOf(err).Type().String()
	if t == "string" {
		errorText = err.(string)
	} else {
		errorText = err.(error).Error()
	}

	var errorValue string
	if status == http.StatusInternalServerError && os.Getenv("HUETIFY_DEBUG") != "true" {
		errorValue = errorMask
	}else{
		errorValue = errorText
	}

	if status == http.StatusInternalServerError {
		fmt.Println(chalk.Redf("[ Error ] %s", errorText))
	}

	failed := json.NewEncoder(w).Encode(&errorContent{ Type: errID, Message: errorValue })
	if failed != nil {
		return false
	}

	return true
}

// Print JSON success from type structure
func Success (db *dbm.Instance, w http.ResponseWriter, status int, success interface{}) bool{
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if db != nil {
		_ = db.Commit()
		_ = db.Close()
	}

	if isNil(success) {
		_, _ = fmt.Fprintf(w, "[]")
		return true
	}

	failed := json.NewEncoder(w).Encode(success)
	if failed != nil {
		return false
	}

	return true
}

// Print nothing with 204 status
func NoContent (db *dbm.Instance, w http.ResponseWriter) bool {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)

	if db != nil {
		_ = db.Commit()
		_ = db.Close()
	}

	return true
}

func isNil(i interface{}) bool {
	if reflect.ValueOf(i).Kind() == reflect.Slice {
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
