package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"ogamebot/controller"
	"os"
	"path"
	"time"
)

var (
	log         = setupLogger() // Prepare logger
	userGenFlag = flag.Bool("generate-user", false, "generates new user config file")
)

const (
	usersDirName   = "./users"
	configFileName = "./config.toml"
)

func main() {

	// Parse flags
	parseArgs()

	// Hello dear users ^^
	log.Info("Ogame bot. MIT License. Copyright 2019 kubastick.")
	log.Infof("Getting file list from %s directory", usersDirName)

	// Check for config
	configData := config{Headless: true, RoundTime: 120000, SeleniumURL: "http://localhost:4444/wd/hub"}
	if fileExists(configFileName) {
		data, err := ioutil.ReadFile("./config.toml")
		if err != nil {
			log.Fatal(err)
		}
		if _, err := toml.Decode(string(data), &configData); err != nil {
			log.Fatal(err)
		}
	} else {
		log.Warningf("Can't find %s - creating and using new with optimal values", configFileName)
		file, err := os.Create(configFileName)
		if err != nil {
			log.Fatalf("Failed to create config file: %s", err.Error())
		}

		encoder := toml.NewEncoder(file)
		err = encoder.Encode(configData)
		if err != nil {
			log.Fatalf("Failed to encode config file data: %s", err.Error())
		}
	}

	// Declare users
	var users []user

	// Load users
	if fileExists(usersDirName) {
		// Load users
		files, err := ioutil.ReadDir(usersDirName)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {

			// Log to console
			log.Info(fmt.Sprintf("Reading file: %s", file.Name()))

			// Read file to ram
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
		// Start infinite loop
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
		log.Warningf("Directory %s does not exists, creating one and exiting", usersDirName)
		err := os.Mkdir(usersDirName, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create %s dir", usersDirName)
		}
	}

	/*service, err := selenium.NewRemote(selenium.Capabilities{"browserName":"Chrome"},"http://localhost:4444/wd/hub")
	if(err!=nil) {
		logger.Fatal(err)
	}
	defer service.Quit()*/
}

func parseArgs() {
	flag.Parse()
	if *userGenFlag {
		log.Info("Generating empty user file")
		err := generateEmptyUserFile("new-user")
		if err != nil {
			log.Fatalf("Failed to generate empty user file: %s", err.Error())
		}
		log.Info("User file generated - exiting")
		os.Exit(0)
	}
}

func generateEmptyUserFile(username string) error {
	userFile, err := os.Create(path.Join(usersDirName, username+".toml"))
	if err != nil {
		return err
	}

	encoder := toml.NewEncoder(userFile)
	err = encoder.Encode(&user{})
	if err != nil {
		return err
	}

	return nil
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func processUser(userData *user, configData *config) {
	log.Info("Starting user processing")
	// Log user name
	log.Info(fmt.Sprintf("User email: %s", userData.Email))

	// Create controller
	ogameController := controller.NewOgameController(
		configData.SeleniumURL,
		userData.Email,
		userData.Password,
		userData.Server,
		userData.ServerButtonId,
		configData.Headless)
	defer ogameController.Close()
	// TODO:Load headless from config (UP 1 LINE) ^

	// Login
	log.Info("Logging in")
	err := ogameController.LoginF()
	if err != nil {
		log.Error(err)
		return
	}

	// Fetch resources
	log.Info("Fetching resources")
	err2 := ogameController.FetchResources()
	if err2 != nil {
		log.Error("Could not fetch resources:" + err2.Error())
		return
	}
	log.Info(fmt.Sprintf("Current resources: Metal %d, Crystal %d, Deuterium %d,Energy %d",
		ogameController.Metal,
		ogameController.Crystal,
		ogameController.Deuterium,
		ogameController.Energy))

	// Upgrade power plants first if energy is too low
	if ogameController.Energy < 0 {
		if ogameController.CanBuildBuilding(controller.MINING_BUILDINGS, controller.POWER_PLANT) {
			log.Info("Upgrading power plant because energy is too low")
			err := ogameController.BuildBuilding(controller.MINING_BUILDINGS, controller.POWER_PLANT)
			if err != nil {
				log.Warningf("Build is impossible: %s", err.Error())
			}
		}
		if ogameController.CanBuildBuilding(controller.MINING_BUILDINGS, controller.DEUTERIUM_POWER_PLANT) {
			log.Info("Upgrading deuterium power plant because energy is too low")
			err := ogameController.BuildBuilding(controller.MINING_BUILDINGS, controller.DEUTERIUM_POWER_PLANT)
			if err != nil {
				log.Warningf("Build is impossible: %s", err.Error())
			}
		}
	} else { // Otherwise try building some buildings
		// Building all building is temporary
		for _, n := range []int{1, 2, 3, 4, 5, 7, 8, 9} {
			if ogameController.CanBuildBuilding(controller.MINING_BUILDINGS, n) {
				log.Info("Upgrading mining building")
				err := ogameController.BuildBuilding(controller.MINING_BUILDINGS, n)
				if err != nil {
					log.Warningf("Build is impossible: %s", err.Error())
				}
			}
		}

		for _, n := range []int{0, 1, 2, 3, 4, 5, 6, 7} {
			if ogameController.CanBuildBuilding(controller.STATION_BUILDINGS, n) {
				log.Info("Upgrading station building")
				err := ogameController.BuildBuilding(controller.STATION_BUILDINGS, n)
				if err != nil {
					log.Warningf("Build is impossible %s", err.Error())
				}
			}
		}
	}

}
