package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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

func TestGetCars(t *testing.T) {
	app := Config{}

	t.Run("returns vehicleStatuses of all cars", func(t *testing.T) {
		request, err := http.NewRequest(http.MethodGet, "/cars", nil)
		if err != nil {
			t.Errorf("Error creating request")
		}

		response := httptest.NewRecorder()

		handler := http.HandlerFunc(app.HandleGetCars)
		handler.ServeHTTP(response, request)

		want := JSONResponse{
			Error:   false,
			Message: "",
			VehicleStatuses: []VehicleStatus{
				{VehicleID: 1,
					FuelLevel:    65,
					BatteryLevel: 40,
					EngineStatus: "Normal",
					SensorStatus: SensorStatus{
						FrontCamera: "Operational",
						RearCamera:  "Operational",
						Radar:       "Operational",
						Lidar:       "Operational",
					},
				}},
		}

		var got JSONResponse
		json.Unmarshal(response.Body.Bytes(), &got)

		if got.VehicleStatuses[0] != want.VehicleStatuses[0] {
			t.Errorf("Got %+v, want %+v", got, want)
		}

	})

}
