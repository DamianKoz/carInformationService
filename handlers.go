package main

import (
	"encoding/json"
	"net/http"
)

func (app *Config) HandleGetCars(w http.ResponseWriter, r *http.Request) {
	vehicleStatuses, err := app.getCarStatuses()
	if err != nil {
		SendJSONResponse(w, http.StatusBadRequest, JSONResponse{Error: true, Message: err.Error(), VehicleStatuses: vehicleStatuses})
	}
	SendJSONResponse(w, http.StatusOK, JSONResponse{Error: false, Message: "", VehicleStatuses: vehicleStatuses})
}

func (app *Config) getCarStatuses() ([]VehicleStatus, error) {
	vehicleStatuses, err := app.getAllCarStatuses()
	if err != nil {
		return nil, err
	}
	return vehicleStatuses, nil
	// return []VehicleStatus{
	// 	{VehicleID: 1,
	// 		FuelLevel:    65,
	// 		BatteryLevel: 40,
	// 		EngineStatus: "Normal",
	// 		SensorStatus: SensorStatus{
	// 			FrontCamera: "Operational",
	// 			RearCamera:  "Operational",
	// 			Radar:       "Operational",
	// 			Lidar:       "Operational",
	// 		},
	// 	},
	// }
}

func SendJSONResponse(w http.ResponseWriter, status int, data JSONResponse) {
	out, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	_, err = w.Write(out)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
