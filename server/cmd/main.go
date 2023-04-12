package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/heesu0/rfc-translator/internal/rfc"
)

func readRFCHandler(w http.ResponseWriter, _ *http.Request) {
	rfcNumber := 5246
	rfcDocument := rfc.NewDocument()
	contents, err := rfcDocument.GetText(rfcNumber)
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "plain/text")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write([]byte(contents))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func setupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/rfc", readRFCHandler).Methods("GET")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../client/build/")))

	return router
}

func main() {
	router := setupRouter()
	httpServer := http.Server{
		Addr:         ":8089",
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	done := make(chan struct{})
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
		<-signalChan

		if err := httpServer.Close(); err != nil {
			fmt.Println(err)
		}
		close(done)
	}()

	fmt.Println("Server started on http://localhost:8089")

	if err := httpServer.ListenAndServe(); err != nil {
		fmt.Println(err)
	}

	<-done
}
