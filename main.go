package main

import (
	"OgameBot/OgameController"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

var log = setupLogger() //Prepare logger

func main() {
	//Hello dear users ^^
	log.Info("Ogame bot. All rights reserved. Copyright 2018 kubastick.")
	log.Info("Getting file list from ./users directory")

	//Check for config
	if fileExists("./config.toml") {
		//TODO: config loading
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

		//Start users loop
		for i, _ := range users {
			processUser(&users[i])
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

func processUser(userData *user) {
	log.Info("Starting user processing")
	//Log user name
	log.Info(fmt.Sprintf("User email: %s", userData.Email))

	//Create controller
	controller := OgameController.NewOgameController(userData.Email, userData.Password, userData.Server, false)
	defer controller.Close()
	//TODO:Load headless from config (UP 1 LINE) ^

	//Login
	log.Info("Logging in")
	err := controller.LoginF()
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Fetching resources")
	err2 := controller.FetchResources()
	if err2 != nil {
		log.Error(err2)
		return
	}
	log.Info(fmt.Sprintf("Current resources: Metal %d, Crystal %d, Deuterium %d,Energy %d",
		controller.Metal,
		controller.Crystal,
		controller.Deuterium,
		controller.Energy))
}
