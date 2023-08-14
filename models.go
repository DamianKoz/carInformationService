package main

import "fmt"

type JSONResponse struct {
	Error           bool            `json:"error"`
	Message         string          `json:"message"`
	VehicleStatuses []VehicleStatus `json:"vehicleStatus"`
}

type VehicleStatus struct {
	VehicleID    uint         `json:"vehicleID"`
	FuelLevel    int          `json:"fuelLevel"`
	BatteryLevel int          `json:"batteryLevel"`
	EngineStatus string       `json:"engineStatus"`
	SensorStatus SensorStatus `json:"sensorStatus"`
}

type SensorStatus struct {
	FrontCamera string `json:"frontCamera"`
	RearCamera  string `json:"rearCamera"`
	Radar       string `json:"radar"`
	Lidar       string `json:"lidar"`
}

func (app *Config) createTables() error {
	sqlStatement := `
	CREATE TABLE IF NOT EXISTS sensor_status (
		sensor_status_id SERIAL PRIMARY KEY,
		front_camera VARCHAR(100),
		rear_camera VARCHAR(100),
		radar VARCHAR(100),
		lidar VARCHAR(100)
	);
	
	CREATE TABLE IF NOT EXISTS vehicle_status (
		vehicle_id SERIAL PRIMARY KEY,
		fuel_level INT,
		battery_level INT,
		engine_status VARCHAR(50),
		sensor_status_id INT REFERENCES sensor_status(sensor_status_id)
	);
	
	`
	_, err := app.DB.Exec(sqlStatement)
	if err != nil {
		return err
	}
	return nil
}
func (app *Config) insert(vehicleStatus VehicleStatus) error {
	// Insert SensorStatus
	ssql := `
	INSERT INTO sensor_status (front_camera, rear_camera, radar, lidar)
	VALUES ($1, $2, $3, $4) RETURNING sensor_status_id;
	`
	var sensorStatusID int
	err := app.DB.QueryRow(ssql, vehicleStatus.SensorStatus.FrontCamera, vehicleStatus.SensorStatus.RearCamera, vehicleStatus.SensorStatus.Radar, vehicleStatus.SensorStatus.Lidar).Scan(&sensorStatusID)
	if err != nil {
		return err
	}

	// Insert VehicleStatus
	vsql := `
	INSERT INTO vehicle_status (fuel_level, battery_level, engine_status, sensor_status_id)
	VALUES ($1, $2, $3, $4);
	`
	_, err = app.DB.Exec(vsql, vehicleStatus.FuelLevel, vehicleStatus.BatteryLevel, vehicleStatus.EngineStatus, sensorStatusID)
	if err != nil {
		return err
	}
	return nil
}

func (app Config) getAllCarStatuses() ([]VehicleStatus, error) {
	sqlStatement := `
	SELECT vs.vehicle_id, vs.fuel_level, vs.battery_level, vs.engine_status,
	ss.front_camera, ss.rear_camera, ss.radar, ss.lidar
	FROM vehicle_status vs
	JOIN sensor_status ss ON vs.sensor_status_id = ss.sensor_status_id;
	`

	rows, err := app.DB.Query(sqlStatement)
	if err != nil {
		return []VehicleStatus{}, err
	}
	defer rows.Close()

	var vehicleStatuses []VehicleStatus

	for rows.Next() {
		fmt.Println("TEST")

		var vs VehicleStatus
		var ss SensorStatus

		err := rows.Scan(
			&vs.VehicleID,
			&vs.FuelLevel,
			&vs.BatteryLevel,
			&vs.EngineStatus,
			&ss.FrontCamera,
			&ss.RearCamera,
			&ss.Radar,
			&ss.Lidar,
		)
		if err != nil {
			panic(err)
		}

		vs.SensorStatus = ss
		vehicleStatuses = append(vehicleStatuses, vs)
		fmt.Printf("VehicleStatus: %v\n", vs)
	}

	if err := rows.Err(); err != nil {
		return vehicleStatuses, err
	}

	return vehicleStatuses, nil

}

func (app *Config) retrieve(vehicleID uint) (VehicleStatus, error) {
	sqlStatement := `
	SELECT vs.vehicle_id, vs.fuel_level, vs.battery_level, vs.engine_status,
	ss.front_camera, ss.rear_camera, ss.radar, ss.lidar
	FROM vehicle_status vs
	JOIN sensor_status ss ON vs.sensor_status_id = ss.sensor_status_id
	WHERE vs.vehicle_id = $1;
	`
	var vs VehicleStatus
	var ss SensorStatus

	err := app.DB.QueryRow(sqlStatement, vehicleID).Scan(
		&vs.VehicleID,
		&vs.FuelLevel,
		&vs.BatteryLevel,
		&vs.EngineStatus,
		&ss.FrontCamera,
		&ss.RearCamera,
		&ss.Radar,
		&ss.Lidar,
	)
	if err != nil {
		return vs, err
	}

	vs.SensorStatus = ss
	return vs, nil
}
