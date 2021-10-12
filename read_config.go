package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

const TIME_LAYOUT = "2006-01-02"
var daysOfWeek = map[string]time.Weekday {
	"mon": time.Monday,
	"tue": time.Tuesday,
	"wed": time.Wednesday,
	"thu": time.Thursday,
	"fri": time.Friday,
	"sat": time.Saturday,
	"sun": time.Sunday,
}

func ReadJson(fileLocation string) *ConfigPlan {
	jsonFile, err := os.Open(fileLocation)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var configPlan ConfigPlan
	byteJson, _ := ioutil.ReadAll(jsonFile)
	jsonFile.Close()
	json.Unmarshal(byteJson, &configPlan)
	// read weekdays into internal time format
	configPlan.Weekdays = readWeekdays(&byteJson)
	if isConfigValid(&configPlan) {
		return &configPlan
	}
	log.Fatal("config keys are not valid")
	return nil
}

func readWeekdays(byteJson *[]byte) map[time.Weekday]string {
	type PlanConfigWeekdaysOnly struct {
		Weekdays map[string]string `json:"weekdays"`
	}
	var configPlanWeekdaysOnly PlanConfigWeekdaysOnly
	json.Unmarshal(*byteJson, &configPlanWeekdaysOnly)
	weekdayMap := make(map[time.Weekday]string)
	for key, value := range configPlanWeekdaysOnly.Weekdays {
		_, found := Find(getWeekdays(daysOfWeek), key)
		if !found {
			log.Fatal("weekdays are not valid")
		}
		weekdayMap[daysOfWeek[key]] = value
	}
	return weekdayMap
}

func isConfigValid(configPlan *ConfigPlan) bool {
	// check if startDate is a valid date
	_, err := time.Parse(TIME_LAYOUT, configPlan.StartDate)
	if err != nil {
		return false
	}
	// check if weekdays use valid splits
	splitNames := make([]string, len(configPlan.Splits))
	for idx, split := range configPlan.Splits {
		splitNames[idx] = split.Name
	}
	for _, val := range configPlan.Weekdays {
		_, found := Find(splitNames, val)
		if !found {
			return false
		}
	}
	// check if exercise identifiers in splits are defined
	availableExercises := make([]string, len(configPlan.Exercises))
	for idx, exercise := range configPlan.Exercises {
		availableExercises[idx] = exercise.Name
	}
	for _, split := range configPlan.Splits {
		for _, execution := range split.Executions {
			for _, exercise := range execution.Variations {
				_, found := Find(availableExercises, exercise)
				if !found {
					return false
				}
			}
		}
	}
	return true
}

func getWeekdays(m map[string]time.Weekday) []string {
	keys := make([]string, len(m))
	i := 0
	for key := range m {
		keys[i] = key
		i++
	}
	return keys
}