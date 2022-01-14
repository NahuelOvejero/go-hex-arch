package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	app "main/application"
	"main/domain"
	"main/repository"
	"main/service"
	"net/http"
	"net/http/httptest"

	"github.com/google/uuid"
)

// tables for tests - use cases

var accCreationTestCases = []struct {
	caseName     string
	carDTO       domain.CarDTO
	expectedCode int
}{
	{
		caseName:     "VALID ACCOUNT - ALL REQUIRED FIELDS",
		carDTO:       validCar,
		expectedCode: http.StatusCreated,
	},
	{
		caseName:     "VALID ACCOUNT - UIID SHOULD BE GENERATED",
		carDTO:       validCarWithoutID,
		expectedCode: http.StatusCreated,
	},
	{
		caseName:     "INVALID ACCOUNT - REQUIRED FIELDS MISSING",
		carDTO:       incompleteCar,
		expectedCode: http.StatusBadRequest,
	},
}

var accFetchTestCases = []struct {
	caseName     string
	ID           string
	expectedCode int
}{
	{
		caseName:     "VALID ACCOUNT ID",
		ID:           carToFetch.Data.ID,
		expectedCode: http.StatusOK,
	},
	{
		caseName:     "NOT FOUND ACCOUNT ID",
		ID:           uuid.NewString(),
		expectedCode: http.StatusNotFound,
	},
	{
		caseName:     "INVALID ACCOUNT ID",
		ID:           "123-notuiid-random-rubish",
		expectedCode: http.StatusBadRequest,
	},
}

var accDeleteTestCases = []struct {
	caseName     string
	ID           string
	version      *int64
	expectedCode int
}{
	{
		caseName:     "VALID ID AND VERSION",
		ID:           carToDelete.Data.ID,
		version:      intAddress(0),
		expectedCode: http.StatusNoContent,
	},
	{
		caseName:     "NOT FOUND ID",
		ID:           uuid.NewString(),
		version:      intAddress(0),
		expectedCode: http.StatusNotFound,
	},
	{
		caseName:     "INVALID VERSION",
		ID:           carToFetch.Data.ID,
		version:      intAddress(99),
		expectedCode: http.StatusConflict,
	},
	{
		caseName:     "INVALID ID",
		ID:           "000-000-000-000-000-000-000",
		version:      intAddress(0),
		expectedCode: http.StatusBadRequest,
	},
}

var validCar = domain.CarDTO{
	Data: &domain.CarData{
		Type:    "cars",
		ID:      uuid.NewString(),
		Version: intAddress(0)},
}

var carToFetch = domain.CarDTO{
	Data: &domain.CarData{
		Type:    "cars",
		ID:      "ab0bd6f5-c3f5-44b2-b677-acd23cdde73c",
		Version: intAddress(0)},
}

var carToDelete = domain.CarDTO{
	Data: &domain.CarData{
		Type:    "cars",
		ID:      "6ba7b814-9dad-11d1-80b4-00c04fd430c8",
		Version: intAddress(0)},
}

var carToDeleteInvalidVersion = domain.CarDTO{
	Data: &domain.CarData{
		Type:    "cars",
		ID:      "fa2dabd9-2de7-428a-a21f-beca3ab8b3d6",
		Version: intAddress(3)},
}

var carForIntegration = domain.CarDTO{
	Data: &domain.CarData{
		Type:    "cars",
		ID:      "1df124da-9a3b-4080-bd19-236d1b0b37d8",
		Version: intAddress(0)},
}

var validCarWithoutID = domain.CarDTO{
	Data: &domain.CarData{
		Type:    "cars",
		Version: intAddress(0)},
}

var incompleteCar = domain.CarDTO{
	Data: &domain.CarData{
		Version: intAddress(0)},
}

// HELPERS

func intAddress(x int64) *int64 {
	return &x
}

func BodyParseFromDTO(dto domain.CarDTO) (*bytes.Buffer, error) {

	jsonReqBody, err := json.Marshal(dto)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(jsonReqBody), nil
}

func SetupPOSTRequest(dto domain.CarDTO) (*httptest.ResponseRecorder, *http.Request, error) {

	postBody, err := BodyParseFromDTO(dto)
	if err != nil {
		return nil, nil, err
	}

	wr := httptest.NewRecorder()

	req, err := http.NewRequest("POST", "/cars", postBody)
	if err != nil {
		return nil, nil, err
	}

	return wr, req, nil
}

func SetupFETCHRequest(id string) (*httptest.ResponseRecorder, *http.Request, error) {

	req, err := http.NewRequest("GET", "/cars/"+id, nil)

	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	wr := httptest.NewRecorder()

	return wr, req, nil
}

func SetupDELETERequest(id string, version *int64) (*httptest.ResponseRecorder, *http.Request, error) {

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/cars/%s?version=%d", id, *version), nil)

	if err != nil {
		return nil, nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	wr := httptest.NewRecorder()

	return wr, req, nil
}

func assertStatusCode(a, b int) error {
	if a != b {
		return fmt.Errorf("expected code %d , got %d", b, a)
	}
	return nil
}

// SETUP
func testSetup() (app.CarHandlers, error) {

	ah := app.CarHandlers{service.NewCarRestService(repository.NewCarRestRepository())}

	wr, req, err := SetupPOSTRequest(carToDelete)
	if err != nil {
		return ah, err
	}
	ah.Create(wr, req)

	wr, req, err = SetupPOSTRequest(carToFetch)
	if err != nil {
		return ah, err
	}
	ah.Create(wr, req)

	return ah, nil
}

// TEARDOWN
var accIDTearDown = []string{carToFetch.Data.ID, carToDeleteInvalidVersion.Data.ID, carForIntegration.Data.ID}

func addNewIDToTearDown(body io.ReadCloser) {
	var reqBody domain.CarDTO
	errDecode := json.NewDecoder(body).Decode(&reqBody)

	if errDecode != nil {
		fmt.Printf("something went wrong adding the ID to teardown list %s", errDecode)
	}

	accIDTearDown = append(accIDTearDown, reqBody.Data.ID)
}

func teardown(testHandlers app.CarHandlers) {

	for _, ID := range accIDTearDown {

		_, err := testHandlers.Service.Delete(ID, "0")

		if err != nil {
			fmt.Printf("something went wrong deleting test car %s", ID)
		}
	}
}
