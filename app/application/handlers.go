package app

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"path"

	"main/domain"
	"main/service"

	"github.com/google/uuid"
)

type CarHandlers struct {
	Service service.CarRestService
}

func (ah *CarHandlers) Fetch(w http.ResponseWriter, r *http.Request) {

	error_msg, code := validateFetchRequest(r)
	if error_msg != "" {
		http.Error(w, error_msg, code)
		return
	}

	id := path.Base(r.URL.Path)

	resp, err := ah.Service.Fetch(id)
	if err != nil {
		http.Error(w, err.Error(), resp.StatusCode)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, err.Error(), resp.StatusCode)
		return
	}

	bodyRespByte, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		http.Error(w, respBodyToString(resp.Body), http.StatusInternalServerError)
		return
	}

	w.Write(bodyRespByte)
}

func (ah *CarHandlers) Create(w http.ResponseWriter, r *http.Request) {

	reqBody, code, err := validateRequestPOST(r)
	if err != nil {
		http.Error(w, err.Error(), code)
	}

	resp, err := ah.Service.Create(reqBody)
	if err != nil {
		http.Error(w, err.Error(), resp.StatusCode)
		return
	}

	bodyRespByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, string(bodyRespByte), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(bodyRespByte)
}

func (ah *CarHandlers) Delete(w http.ResponseWriter, r *http.Request) {

	id, version, error_msg, code := validateRequestDELETE(r)
	if error_msg != "" {
		http.Error(w, error_msg, code)
		return
	}

	resp, err := ah.Service.Delete(id, version)
	if err != nil {
		http.Error(w, err.Error(), resp.StatusCode)
		return
	}

	if resp.StatusCode == http.StatusConflict {
		http.Error(w, "specified version is incorrect", resp.StatusCode)
	}

	w.WriteHeader(resp.StatusCode)
}

func (ah *CarHandlers) HandleByMethod(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		ah.Fetch(w, r)
	case http.MethodDelete:
		ah.Delete(w, r)
	default:
		http.Error(w, "method not allowed. did you mean to POST to /cars?", http.StatusMethodNotAllowed)
	}
}

// HELPERS

func validateRequestPOST(r *http.Request) (domain.CarDTO, int, error) {

	var reqBody domain.CarDTO

	if r.Method != http.MethodPost {
		return reqBody, http.StatusMethodNotAllowed, errors.New("only post method allowed. did you mean to get on /cars/{id}?")
	}

	errDecode := json.NewDecoder(r.Body).Decode(&reqBody)
	if errDecode != nil {
		return reqBody, http.StatusBadRequest, errors.New("Invalid Car JSON data \n " + errDecode.Error())
	}

	if reqBody.Data.ID == "" {
		reqBody.Data.ID = uuid.NewString()
	}

	return reqBody, 0, nil
}

func validateRequestDELETE(r *http.Request) (string, string, string, int) {

	version := r.URL.Query().Get("version")
	if version == "" {
		return "", version, "expected \"version\" as a query parameter to delete an car", http.StatusBadRequest
	}

	id := path.Base(r.URL.Path)
	if id == "" {
		return id, version, "Can't get an car for empty id", http.StatusBadRequest
	}

	return id, version, "", 0
}

func validateFetchRequest(r *http.Request) (string, int) {

	id := path.Base(r.URL.Path)
	if id == "" {
		return "Can't get an car for empty id", http.StatusBadRequest
	}
	//Check to avoid /cars/foo/boo/random/{uiid}
	if (path.Dir(r.URL.Path)) != "/cars" {
		return "404 page not found - did you mean: /cars/" + id, http.StatusNotFound
	}

	return "", 0
}

func respBodyToString(respBody io.ReadCloser) string {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(respBody)
	if err != nil {
		return "Error performing the request"
	}
	bodyString := buf.String()
	return bodyString
}
