package repository

import (
	"log"

	"github.com/KitaPDev/fogfarms-server/models"
	"github.com/KitaPDev/fogfarms-server/models/outputs"
	"github.com/KitaPDev/fogfarms-server/src/database"
	"github.com/KitaPDev/fogfarms-server/src/util/module"
	"github.com/lib/pq"
)

func GetLatestSensorData(moduleGroupID int) (map[string]*outputs.Dashboardoutput, error) {
	var moduleGroupIDs []int

	var dashboard = make(map[string]*outputs.Dashboardoutput)
	moduleGroupIDs = append(moduleGroupIDs, moduleGroupID)

	modules, err := module.GetModulesByModuleGroupIDs(moduleGroupIDs)
	if err != nil {
		return nil, err
	}
	var moduleIDs []int
	for _, m := range modules {
		moduleIDs = append(moduleIDs, m.ModuleID)
	}
	log.Println(moduleIDs)
	db := database.GetDB()

	sqlStatement := `select modulelabel,sensordata.moduleid,sensordata.timestamp,arrnutrientunittds,arrnutrientunitph,arrnutrientunitsolutiontemperature,arrgrowunitlux,arrgrowunithumidity,arrgrowunittemperature,nutrientamount from sensordata inner join (SELECT moduleid, max(timestamp) AS maxtime FROM sensordata GROUP BY moduleid) as maxTable on maxTable.moduleid=sensordata.moduleid AND sensordata.timestamp=maxTable.maxtime inner join (select moduleid, count(*) as nutrientamount from nutrientunit group by moduleid) AS nutrient on nutrient.moduleid = sensordata.moduleid AND nutrient.moduleid=maxTable.moduleid inner join module on module.moduleid=sensordata.moduleid where sensordata.moduleid = ANY($1);`
	rows, err := db.Query(sqlStatement, pq.Array(moduleIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sd models.SensorData
		var modulelabel string
		var nutrientAmount int
		err = rows.Scan(
			&modulelabel,
			&sd.ModuleID,
			&sd.TimeStamp,
			pq.Array(&sd.TDS),
			pq.Array(&sd.PH),
			pq.Array(&sd.SolutionTemperature),
			pq.Array(&sd.GrowUnitLux),
			pq.Array(&sd.GrowUnitHumidity),
			pq.Array(&sd.GrowUnitTemperature),
			&nutrientAmount,
		)
		if err != nil {
			return nil, err
		}
		log.Println(dashboard)
		log.Println(modulelabel)

		dashboard[modulelabel] = &outputs.Dashboardoutput{
			NutrientAmount: nutrientAmount,
			Sensordata:     sd,
		}

	}

	sqlStatement = `select moduleid,modulelabel,arrfogger,arrled,arrmixer,arrsolenoidvalve from module where moduleid=ANY($1);`
	rows, err = db.Query(sqlStatement, pq.Array(moduleIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var sd outputs.DashBoardModule
		var modulelabel string
		var moduleid int

		err = rows.Scan(
			&moduleid,
			&modulelabel,
			pq.Array(&sd.Fogger),
			pq.Array(&sd.LED),
			pq.Array(&sd.Mixer),
			pq.Array(&sd.SolenoidValve),
		)
		if err != nil {
			return nil, err
		}
		dashboard[modulelabel].Device = sd

	}
	return dashboard, nil
}

func RecordSensorData(moduleID int, tds []float64, ph []float64, solutionTemperature []float64,
	lux []float64, humidity []float64, temperature []float64) error {

	db := database.GetDB()

	sqlStatement :=
		`INSERT INTO SensorData (moduleid, arrnutrientunittds, arrnutrientunitph, arrnutrientunitsolutiontemperature, arrgrowunitlux, arrgrowunithumidity, arrgrowunittemperature)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.Query(sqlStatement, moduleID, pq.Array(tds), pq.Array(ph), pq.Array(solutionTemperature),
		pq.Array(lux), pq.Array(humidity), pq.Array(temperature))

	if err != nil {
		return err
	}

	return nil
}