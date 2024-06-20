package main

import (
	httphandlers "amartha-test/entities/http_handlers"
	"amartha-test/handlers"
	usecase "amartha-test/usecases"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func getRoutes(modules module) *chi.Mux {
	router := chi.NewRouter()

	router.Route("/", func(path chi.Router) {
		path.Get("/reconciliation", modules.httpHandler.TransactionHandler.HandleReconciliation)
	})

	return router
}

func main() {

	transactionsUsecase := usecase.NewTransactionUsecase(usecase.TransactionUsecase{})

	transactionsHandler := handlers.NewTransactionHandler(handlers.TransactionHandler{
		TransactionUsecase: transactionsUsecase,
	})

	modules := loadModules(httphandlers.Handlers{
		TransactionHandler: transactionsHandler,
	})

    router := getRoutes(modules)

	log.Fatal(http.ListenAndServe(":8000", router))
}