package main

import (
	"OgameBot/OgameController"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"time"
)

var log = setupLogger() //Prepare logger

func main() {
	//Hello dear users ^^
	log.Info("Ogame bot. All rights reserved. Copyright 2018 kubastick.")
	log.Info("Getting file list from ./users directory")

	//Check for config
	configData := config{Headless: true, RoundTime: 120000, SeleniumURL: "http://localhost:4444/wd/hub"}
	if fileExists("./config.toml") {
		data, err := ioutil.ReadFile("./config.toml")
		if err != nil {
			log.Fatal(err)
		}
		if _, err := toml.Decode(string(data), &configData); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Warning("Can't find config.toml using optimal defaults")
	}

	//Declare users
	var users []user

	//Load users
	if fileExists("./users") {
		//Load users
		files, err := ioutil.ReadDir("./users")
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {

			//Log to console
			log.Info(fmt.Sprintf("Reading file: %s", file.Name()))

			//Read file to ram
			data, err := ioutil.ReadFile("./users/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			var userData user
			if _, err := toml.Decode(string(data), &userData); err != nil {
				log.Fatal(err)
			}
			users = append(users, userData)
		}
		//Start infinite loop
		for {
			nextIterationTime := time.Now().Second()*1000 + configData.RoundTime

			//Start users loop
			for i := range users {
				processUser(&users[i], &configData)
			}

			waitTime := nextIterationTime - time.Now().Second()*1000
			if waitTime > 0 {
				log.Infof("Round finished, next in %d seconds", waitTime/1000)
				sleepTime := time.Duration(waitTime) * time.Millisecond
				time.Sleep(sleepTime)
			} else {
				log.Warningf("Consider increasing round time, you are running %d seconds behind", waitTime*-1)
			}
		}
	} else {
		log.Fatal("Directory ./users does not exists, exiting")
	}

	/*service, err := selenium.NewRemote(selenium.Capabilities{"browserName":"Chrome"},"http://localhost:4444/wd/hub")
	if(err!=nil) {
		logger.Fatal(err)
	}
	defer service.Quit()*/
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func processUser(userData *user, configData *config) {
	log.Info("Starting user processing")
	//Log user name
	log.Info(fmt.Sprintf("User email: %s", userData.Email))

	//Create controller
	controller := OgameController.NewOgameController(
		configData.SeleniumURL,
		userData.Email,
		userData.Password,
		userData.Server,
		userData.ServerButtonId,
		configData.Headless)
	defer controller.Close()
	//TODO:Load headless from config (UP 1 LINE) ^

	//Login
	log.Info("Logging in")
	err := controller.LoginF()
	if err != nil {
		log.Error(err)
		return
	}

	//Fetch resources
	log.Info("Fetching resources")
	err2 := controller.FetchResources()
	if err2 != nil {
		log.Error("Could not fetch resources:" + err2.Error())
		return
	}
	log.Info(fmt.Sprintf("Current resources: Metal %d, Crystal %d, Deuterium %d,Energy %d",
		controller.Metal,
		controller.Crystal,
		controller.Deuterium,
		controller.Energy))

	//Upgrade power plants first if energy is too low
	if controller.Energy < 0 {
		if controller.CanBuildBuilding(OgameController.MINING_BUILDINGS, OgameController.POWER_PLANT) {
			log.Info("Upgrading power plant because energy is too low")
			err := controller.BuildBuilding(OgameController.MINING_BUILDINGS, OgameController.POWER_PLANT)
			if err != nil {
				log.Warningf("Build is impossible: %s", err.Error())
			}
		}
		if controller.CanBuildBuilding(OgameController.MINING_BUILDINGS, OgameController.DEUTERIUM_POWER_PLANT) {
			log.Info("Upgrading deuterium power plant because energy is too low")
			err := controller.BuildBuilding(OgameController.MINING_BUILDINGS, OgameController.DEUTERIUM_POWER_PLANT)
			if err != nil {
				log.Warningf("Build is impossible: %s", err.Error())
			}
		}
	}

	//Build all buildings (temporary)
	for _, n := range []int{1, 2, 3, 4, 5, 7, 8, 9} {
		if controller.CanBuildBuilding(OgameController.MINING_BUILDINGS, n) {
			log.Info("Upgrading mining building")
			err := controller.BuildBuilding(OgameController.MINING_BUILDINGS, n)
			if err != nil {
				log.Warningf("Build is impossible: %s", err.Error())
			}
		}
	}

	for _, n := range []int{0, 1, 2, 3, 4, 5, 6, 7} {
		if controller.CanBuildBuilding(OgameController.STATION_BUILDINGS, n) {
			log.Info("Upgrading station building")
			err := controller.BuildBuilding(OgameController.STATION_BUILDINGS, n)
			if err != nil {
				log.Warningf("Build is impossible %s", err.Error())
			}
		}
	}

}
