package OgameController

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	//Pages
	MAIN_PAGE    = "%s/index.php?"
	MINING_PAGE  = "%s/index.php?page=resources"
	STATION_PAGE = "%s/index.php?page=station"

	//Timeouts
	CHECK_TIMEOUT = time.Second * 4
	FIND_TIMEOUT  = time.Second * 10
)

type OgameController struct {
	//Constructors
	Login          string
	Password       string
	Server         string
	Headless       bool
	ServerButtonID int
	//Driver
	driver selenium.WebDriver

	//Resources
	Metal     int
	Crystal   int
	Deuterium int
	Energy    int
}

func NewOgameController(seleniumAdress string, login string, password string, server string, serverButtonID int, headless bool) OgameController {

	//Create OgameController object
	controller := OgameController{
		Login:          login,
		Password:       password,
		Server:         server,
		Headless:       headless,
		ServerButtonID: serverButtonID,
	}
	//Add chrome specific options
	options := chrome.Capabilities{}
	if headless {
		options.Args = []string{"--headless"}
	}

	//Create capabilities object
	capabilities := selenium.Capabilities{"browserName": "Chrome"}
	capabilities.AddChrome(options)

	//Start new driver and bind it to controller
	driver, err := selenium.NewRemote(capabilities, "http://localhost:4444/wd/hub")
	if err != nil {
		println("FATAL: Could not connect to selenium server")
		log.Fatal(err)
	}
	controller.driver = driver

	//Resize window
	windowHandle, _ := controller.driver.CurrentWindowHandle()
	controller.driver.ResizeWindow(windowHandle, 1366, 768)

	//Set default timeout
	controller.driver.SetImplicitWaitTimeout(FIND_TIMEOUT)
	return controller
}

func (o *OgameController) LoginF() error {
	o.driver.Get(`https://pl.ogame.gameforge.com/`)

	//Press login tab
	loginTab, err := o.driver.FindElement(selenium.ByID, "ui-id-1")
	loginTab.Click()
	//Dismiss cookie alert
	cookieCloseButton, err := o.driver.FindElement(selenium.ByID, "accept_btn")
	cookieCloseButton.Click()
	//Send login text
	usernameLogin, err := o.driver.FindElement(selenium.ByID, "usernameLogin")
	usernameLogin.SendKeys(o.Login)
	//Send password text
	password, err := o.driver.FindElement(selenium.ByID, "passwordLogin")
	password.SendKeys(o.Password)
	//Click submit button
	submitButton, err := o.driver.FindElement(selenium.ByID, "loginSubmit")
	submitButton.Click()
	//Uni button
	universumButton, err := o.driver.FindElement(selenium.ByXPATH, fmt.Sprintf("//*[@id=\"accountlist\"]/div/div[1]/div[2]/div/div/div[%d]/button", o.ServerButtonID))
	universumButton.Click()
	//Wait for JS to load
	time.Sleep(1 * time.Second)
	//Login to cosmic server
	o.driver.Get(fmt.Sprintf(MAIN_PAGE, o.Server))
	o.closeOtherTabs() //causing segmentation fault
	return err
}

func (o *OgameController) closeOtherTabs() {
	mainWindow, _ := o.driver.CurrentWindowHandle()
	allWindows, _ := o.driver.WindowHandles()
	for _, handle := range allWindows {
		if handle != mainWindow {
			o.driver.CloseWindow(handle)
		}
	}
	mainWindow, _ = o.driver.CurrentWindowHandle() //This is weird, but it prevent's errors
	o.driver.SwitchWindow(mainWindow)
}

func (o *OgameController) Close() {
	allWindows, _ := o.driver.WindowHandles()
	for _, handle := range allWindows {
		o.driver.CloseWindow(handle)
	}
	o.driver.Close()
}

func (o *OgameController) FetchResources() error {
	o.getIfAnother(fmt.Sprintf(MAIN_PAGE, o.Server))

	//Metal
	metalText, err := o.getResourceText("resources_metal")
	o.Metal = o.parseResourceText(metalText)
	//Crystal
	crystalText, err := o.getResourceText("resources_crystal")
	o.Crystal = o.parseResourceText(crystalText)
	//Deuterium
	deuteriumText, err := o.getResourceText("resources_deuterium")
	o.Deuterium = o.parseResourceText(deuteriumText)
	//Energy
	energyText, err := o.getResourceText("resources_energy")
	o.Energy = o.parseResourceText(energyText)
	return err
}

func (o *OgameController) getIfAnother(url string) {
	currentURL, _ := o.driver.CurrentURL()
	if currentURL != url {
		o.driver.Get(url)
	}
}

func (o *OgameController) parseResourceText(text string) int {
	//Remove whitespaces
	text2 := strings.Replace(text, " ", "", -1)
	//Remove dots
	text3 := strings.Replace(text2, ".", "", -1)
	//return
	result, _ := strconv.Atoi(text3)
	return result
}

func (o *OgameController) getResourceText(id string) (string, error) {
	element, err := o.driver.FindElement(selenium.ByID, id)
	elementText, err := element.Text()
	return elementText, err
}

func (o *OgameController) CanBuildBuilding(category int, building int) bool {
	if category == 0 {
		o.getIfAnother(fmt.Sprintf(MINING_PAGE, o.Server))
	} else if category == 1 {
		o.getIfAnother(fmt.Sprintf(STATION_PAGE, o.Server))
	} else {
		panic("Paramater outside range")
	}
	defer o.driver.SetImplicitWaitTimeout(FIND_TIMEOUT)
	o.driver.SetImplicitWaitTimeout(CHECK_TIMEOUT)
	//Click building button
	buildingButton, _ := o.driver.FindElement(selenium.ByID, fmt.Sprintf("button%d", building))
	buildingButton.Click()
	_, err := o.driver.FindElement(selenium.ByClassName, "build-it")
	if err != nil {
		return false
	}
	return true
}

func (o *OgameController) BuildBuilding(category int, building int) error {
	if category == 0 {
		o.getIfAnother(fmt.Sprintf(MINING_PAGE, o.Server))
	} else if category == 1 {
		o.getIfAnother(fmt.Sprintf(STATION_PAGE, o.Server))
	} else {
		panic("Parameter outside range")
	}
	defer o.driver.SetImplicitWaitTimeout(FIND_TIMEOUT)
	o.driver.SetImplicitWaitTimeout(CHECK_TIMEOUT)
	//Click building button
	buildingButton, err := o.driver.FindElement(selenium.ByID, fmt.Sprintf("button%d", building))
	buildingButton.Click()
	buildButton, err := o.driver.FindElement(selenium.ByClassName, "build-it")
	buildButton.Click()
	return err
}
