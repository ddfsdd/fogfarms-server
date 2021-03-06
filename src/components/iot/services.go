package iot

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/KitaPDev/fogfarms-server/src/jsonhandler"
	"github.com/KitaPDev/fogfarms-server/src/util/module"
	"github.com/KitaPDev/fogfarms-server/src/util/modulegroup"
	"github.com/KitaPDev/fogfarms-server/src/util/sensordata"
)

func Update(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Token          string               `json:"token"`
		SensorData     map[string][]float64 `json:"sensor_data"`
		Mixers         []bool               `json:"mixers"`
		SolenoidValves []bool               `json:"solenoid_valves"`
		LEDs           []bool               `json:"leds"`
		Foggers        []bool               `json:"foggers"`
	}

	input := Input{}
	success := jsonhandler.DecodeJsonFromBody(w, r, &input)
	if !success {
		return
	}

	moduleID, err := module.GetModuleIDByToken(input.Token)
	if err != nil || moduleID == 0 {
		msg := "Error: Failed to Get ModuleID By Token"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	onAuto, err := modulegroup.GetOnAutoByModuleID(moduleID)
	if err != nil {
		msg := "Error: Failed to Get OnAuto By ModuleID"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	tds := input.SensorData["tds"]
	ph := input.SensorData["ph"]
	solutionTemp := input.SensorData["solution_temp"]
	lux := input.SensorData["grow_unit_lux"]
	humidity := input.SensorData["grow_unit_humidity"]
	temperature := input.SensorData["grow_unit_temp"]

	err = sensordata.RecordSensorData(moduleID, tds, ph, solutionTemp, lux, humidity, temperature)
	if err != nil {
		msg := "Error: Failed to Record Sensor Data"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	if onAuto {
		err = module.UpdateDeviceStatus(moduleID, input.Mixers, input.SolenoidValves, input.LEDs,
			input.Foggers)
		if err != nil {
			msg := "Error: Failed to Record Device Status"
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}
	}
	log.Println([]int{moduleID})
	modulegroup, err := modulegroup.GetModuleGroupsByModuleID(moduleID)
	log.Println(modulegroup)
	fogger, led, mixer, solenoidValve, err := module.GetDeviceStatus(moduleID)
	if err != nil {
		msg := "Error: Failed to Get Device Status"
		http.Error(w, msg, http.StatusInternalServerError)
		return
	}

	type Output struct {
		OnAuto         bool    `json:"on_auto"`
		Foggers        []bool  `json:"foggers"`
		LEDs           []bool  `json:"leds"`
		Mixers         []bool  `json:"mixers"`
		SolenoidValves []bool  `json:"solenoid_valves"`
		TDS            float64 `json:"tds"`
		PH             float64 `json:"ph"`
		Humidity       float64 `json:"humidity"`
		LightsOffHour  float64 `json:"lights_off_hour"`
		LightsOnHour   float64 `json:"lights_on_hour"`
	}

	output := Output{
		OnAuto:         onAuto,
		Foggers:        fogger,
		LEDs:           led,
		Mixers:         mixer,
		SolenoidValves: solenoidValve,
		Humidity:       modulegroup.Humidity,
		TDS:            modulegroup.TDS,
		PH:             modulegroup.PH,
		LightsOffHour:  modulegroup.LightsOffHour,
		LightsOnHour:   modulegroup.LightsOnHour,
	}

	jsonData, err := json.Marshal(output)
	if err != nil {
		msg := "Error: Failed to marshal JSON"
		http.Error(w, msg, http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
