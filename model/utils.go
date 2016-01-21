package model

import (
	"encoding/json"
	"io"
	"time"
)

func GetMillis() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type AppError struct {
	Message       string `json:"message"`        // Message to be display to the end user without debugging information
	DetailedError string `json:"detailed_error"` // Internal error string to help the developer
	RequestId     string `json:"request_id"`     // The RequestId that's also set in the header
	StatusCode    int    `json:"status_code"`    // The http status code
	Where         string `json:"-"`              // The function where it happened in the form of Struct.Func
}

//Error custom error implementation
func (er *AppError) Error() string {
	return er.Where + ": " + er.Message + ", " + er.DetailedError
}

//NewAppError return custom error
func NewAppError(where string, message string, details string) *AppError {
	ap := &AppError{}
	ap.Message = message
	ap.Where = where
	ap.DetailedError = details
	ap.StatusCode = 500
	return ap
}

//MapToJson will encode map to a json string
func MapToJson(objmap map[string]string) string {
	if b, err := json.Marshal(objmap); err != nil {
		return ""
	} else {
		return string(b)
	}
}

//MapFromJson will decode key/value json to map
func MapFromJson(data io.Reader) map[string]string {
	decoder := json.NewDecoder(data)

	var objmap map[string]string
	if err := decoder.Decode(&objmap); err != nil {
		return make(map[string]string)
	} else {
		return objmap
	}
}

//ArrayToJson will encode map to a json string
func ArrayToJson(objmap []string) string {
	if b, err := json.Marshal(objmap); err != nil {
		return ""
	} else {
		return string(b)
	}
}

//ArrayFromJson will decode array of strings from json
func ArrayFromJson(data io.Reader) []string {
	decoder := json.NewDecoder(data)

	var objarr []string
	if err := decoder.Decode(&objarr); err != nil {
		return make([]string, 0)
	} else {
		return objarr
	}
}

func BoolMapFromJson(data io.Reader) map[string]bool {
	decoder := json.NewDecoder(data)

	var objmap map[string]bool
	if err := decoder.Decode(&objmap); err != nil {
		return make(map[string]bool)
	} else {
		return objmap
	}
}

func BoolMapToJson(objmap map[string]bool) string {
	if b, err := json.Marshal(objmap); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func StringInterfaceToJson(objmap map[string]interface{}) string {
	if b, err := json.Marshal(objmap); err != nil {
		return ""
	} else {
		return string(b)
	}
}

func StringInterfaceFromJson(data io.Reader) map[string]interface{} {
	decoder := json.NewDecoder(data)

	var objmap map[string]interface{}
	if err := decoder.Decode(&objmap); err != nil {
		return make(map[string]interface{})
	} else {
		return objmap
	}
}
