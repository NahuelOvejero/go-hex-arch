package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"main/domain"
	"net/http"
	"os"
)

var host = os.Getenv("ACC_API_HOST")

// implements CarRepositoryPort - Structural Typing
type CarRestRepository struct {
}

func NewCarRestRepository() CarRestRepository {
	return CarRestRepository{}
}

func (as CarRestRepository) Fetch(id string) (http.Response, error) {

	resp, err := carApiFetch(id)

	if err != nil {
		return *resp, err
	}

	if resp.StatusCode != http.StatusOK {
		return *resp, errors.New(respBodyErrorString(*resp))
	}

	return *resp, nil
}

func (as CarRestRepository) Create(accData domain.CarDTO) (http.Response, error) {

	resp, err := carApiPost(accData)

	if err != nil {
		return *resp, err
	}

	if resp.StatusCode != http.StatusCreated {
		return *resp, errors.New(respBodyErrorString(*resp))
	}

	return *resp, nil
}

func (as CarRestRepository) Delete(id, version string) (http.Response, error) {

	resp, err := carApiDelete(id, version)

	if err != nil {
		return *resp, err
	}

	return *resp, nil
}

//Helpers

func carApiFetch(id string) (*http.Response, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/vaulta/api/cars/%s", host, id), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func carApiPost(reqBody domain.CarDTO) (*http.Response, error) {

	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/vaulta/api/cars", host), bytes.NewBuffer(jsonReqBody))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func carApiDelete(id, version string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/vaulta/api/cars/%s?version=%s", host, id, version), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func respBodyErrorString(resp http.Response) string {

	var errResp errorResponse

	errDecode := json.NewDecoder(resp.Body).Decode(&errResp)
	if errDecode != nil {
		return "error performing the request"
	}

	if errResp.ErrorMessage != "" {
		return errResp.ErrorMessage
	}
	return "error performing the request"
}

type errorResponse struct {
	ErrorMessage string `json:"error_message,omitempty"`
}
