package test

import (
	"fmt"
	app "main/application"
	"net/http"
	"os"
	"testing"
)

var testHandlers app.CarHandlers

func TestMain(m *testing.M) {

	var err error
	testHandlers, err = testSetup()

	if err != nil {
		fmt.Println("could not execute test setup. exiting...")
		os.Exit(0)
	}

	m.Run()

	teardown(testHandlers)
}

func TestFetch(t *testing.T) {

	for _, testcase := range accFetchTestCases {

		t.Logf("testing fetch use case: %s ", testcase.caseName)

		wr, req, err := SetupFETCHRequest(testcase.ID)
		if err != nil {
			t.Errorf("unexpected error running the test: %s \n error: %s", testcase.caseName, err)
		}

		testHandlers.Fetch(wr, req)
		err = assertStatusCode(wr.Code, testcase.expectedCode)
		if err != nil {
			t.Errorf(err.Error())
		}
	}
}

func TestCreateCar(t *testing.T) {

	for _, testcase := range accCreationTestCases {

		t.Logf("testing create use case: %s ", testcase.caseName)

		wr, req, err := SetupPOSTRequest(testcase.carDTO)

		if err != nil {
			t.Errorf("Unexpected error running the test: %s \n Error: %s", testcase.caseName, err)
		}

		testHandlers.Create(wr, req)

		err = assertStatusCode(wr.Code, testcase.expectedCode)
		if err != nil {
			t.Errorf(err.Error())
		}

		if testcase.expectedCode == http.StatusCreated {
			addNewIDToTearDown(wr.Result().Body)
		}
	}
}

func TestDeleteCar(t *testing.T) {

	for _, testcase := range accDeleteTestCases {

		t.Logf("testing delete use case: %s ", testcase.caseName)

		wr, req, err := SetupDELETERequest(testcase.ID, testcase.version)
		if err != nil {
			t.Errorf("unexpected error running the test: %s \n error: %s", testcase.caseName, err)
		}

		testHandlers.Delete(wr, req)

		err = assertStatusCode(wr.Code, testcase.expectedCode)
		if err != nil {
			t.Errorf(testcase.caseName + " " + err.Error())
		}
	}
}
func TestIntegration(t *testing.T) {

	//Car Create
	res, err := testHandlers.Service.Create(carForIntegration)
	if err != nil {
		t.Errorf("integration test failed. error: %s . statuscode %d :", err.Error(), res.StatusCode)
	}
	err = assertStatusCode(res.StatusCode, http.StatusCreated)
	if err != nil {
		t.Error(err.Error())
	}

	//Fetch Car created previosuly
	res, err = testHandlers.Service.Fetch(carForIntegration.Data.ID)
	if err != nil {
		t.Errorf("integration test failed. error: %s . statuscode %d :", err.Error(), res.StatusCode)
	}
	err = assertStatusCode(res.StatusCode, http.StatusOK)
	if err != nil {
		t.Error(err.Error())
	}

	//Delete Car
	res, err = testHandlers.Service.Delete(carForIntegration.Data.ID, "0")
	if err != nil {
		t.Errorf("integration test failed. error: %s . statuscode %d :", err.Error(), res.StatusCode)
	}
	err = assertStatusCode(res.StatusCode, http.StatusNoContent)
	if err != nil {
		t.Error(err.Error())
	}

	//Fetch Car Should Not Exist Now
	res, err = testHandlers.Service.Fetch(carForIntegration.Data.ID)
	if err == nil {
		t.Errorf("integration test failed. error: %s . statuscode %d :", err.Error(), res.StatusCode)
	}
	err = assertStatusCode(res.StatusCode, http.StatusNotFound)
	if err != nil {
		t.Error(err.Error())
	}

	//Should be able to create the car that was delete
	res, err = testHandlers.Service.Create(carForIntegration)
	if err != nil {
		t.Errorf("integration test failed. error: %s . statuscode %d :", err.Error(), res.StatusCode)
	}
	err = assertStatusCode(res.StatusCode, http.StatusCreated)
	if err != nil {
		t.Error(err.Error())
	}

}
