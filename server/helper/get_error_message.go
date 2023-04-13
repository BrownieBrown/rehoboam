package helper

import (
	"encoding/json"
	"net/http/httptest"
)

func GetErrorMessage(resp *httptest.ResponseRecorder) string {
	var respBody map[string]string
	json.Unmarshal(resp.Body.Bytes(), &respBody)
	return respBody["error"]
}
