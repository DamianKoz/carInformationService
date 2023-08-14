package main

import (
	"github.com/go-chi/chi/v5"
)

func (app *Config) Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Route("/cars", func(r chi.Router) {
		r.Get("/", app.HandleGetCars)
	})
	// 	r.Post("/", createCar)

	// 	// Subrouter for one car
	// 	r.Route("/{carID}", func(r chi.Router) {
	// 		r.Get("/", getCar)
	// 		r.Put("/", updateCar)
	// 		r.Delete("/", deleteCar)
	// 	})
	// })

	return r
}
