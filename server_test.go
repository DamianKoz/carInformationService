package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestSendJSONResponse(t *testing.T) {
	response := httptest.NewRecorder()

	data := JSONResponse{
		Error:           false,
		Message:         "",
		VehicleStatuses: make([]VehicleStatus, 1),
	}

	SendJSONResponse(response, http.StatusOK, data)

	if response.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, response.Code)
	}
	var got JSONResponse
	err := json.Unmarshal(response.Body.Bytes(), &got)
	if err != nil {
		t.Errorf("Could not unmarshal json response: %s", response.Body.String())
	}

	if got.VehicleStatuses[0] != data.VehicleStatuses[0] {
		t.Errorf("Got %+v, want %+v", got, data)
	}
}

func TestHandleGetCars(t *testing.T) {
	db, err := connectToDB()

	app := Config{
		DB: db,
	}

	app.createTables()

	req, err := http.NewRequest("GET", "/cars", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	w := httptest.NewRecorder()
	app.HandleGetCars(w, req)

	res := w.Result()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK, got %v", res.StatusCode)
	}

	var jsonResponse JSONResponse
	err = json.NewDecoder(res.Body).Decode(&jsonResponse)
	if err != nil {
		t.Errorf("could not decode response: %v", err)
	}

	if jsonResponse.Error {
		t.Errorf("expected error to be false, got true")
	}
}

// func TestGetCars(t *testing.T) {
// 	app := Config{}

// 	t.Run("returns vehicleStatuses of all cars", func(t *testing.T) {
// 		request, err := http.NewRequest(http.MethodGet, "/cars", nil)
// 		if err != nil {
// 			t.Errorf("Error creating request")
// 		}

// 		response := httptest.NewRecorder()

// 		handler := http.HandlerFunc(app.HandleGetCars)
// 		handler.ServeHTTP(response, request)

// 		want := JSONResponse{
// 			Error:   false,
// 			Message: "",
// 			VehicleStatuses: []VehicleStatus{
// 				{VehicleID: 1,
// 					FuelLevel:    65,
// 					BatteryLevel: 40,
// 					EngineStatus: "Normal",
// 					SensorStatus: SensorStatus{
// 						FrontCamera: "Operational",
// 						RearCamera:  "Operational",
// 						Radar:       "Operational",
// 						Lidar:       "Operational",
// 					},
// 				}},
// 		}

// 		var got JSONResponse
// 		json.Unmarshal(response.Body.Bytes(), &got)

// 		if got.VehicleStatuses[0] != want.VehicleStatuses[0] {
// 			t.Errorf("Got %+v, want %+v", got, want)
// 		}

// 	})

// }
