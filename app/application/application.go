package app

import (
	"fmt"
	"log"
	"main/repository"
	"main/service"
	"net/http"
)

func StartApplication() {

	SetupPathAndHandlers(CarHandlers{service.NewCarRestService(repository.NewCarRestRepository())})

	fmt.Println("Starting application on port :5050")
	if err := http.ListenAndServe(":5050", nil); err != nil {
		log.Fatal(err)
	}
}

func SetupPathAndHandlers(ah CarHandlers) {
	http.HandleFunc("/cars/", ah.HandleByMethod)
	http.HandleFunc("/cars", ah.Create)
}
