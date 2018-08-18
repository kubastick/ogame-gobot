package OgameController

import (
	"fmt"
	"github.com/tebeka/selenium"
	"github.com/tebeka/selenium/chrome"
	"log"
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
	Login    string
	Password string
	Server   string
	Headless bool

	//Driver
	driver selenium.WebDriver

	//Resources
	Metal     int
	Crystal   int
	Deuterium int
	Energy    int
}

func NewOgameController(login string, password string, server string, headless bool) OgameController {

	//Create OgameController object
	controller := OgameController{
		Login:    login,
		Password: password,
		Server:   server,
		Headless: headless,
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
	universumButton, err := o.driver.FindElement(selenium.ByXPATH, "//*[@id=\"accountlist\"]/div/div[1]/div[2]/div/div/div[11]/button")
	universumButton.Click()
	//Wait for JS to load
	time.Sleep(1000 * time.Second)
	//Login to cosmic server
	o.driver.Get(fmt.Sprintf(MAIN_PAGE, o.Server))
	o.closeOtherTabs()

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
}

func (o *OgameController) Close() {
	o.driver.Close()
}

func (o *OgameController) FetchResources() {

}

func (o *OgameController) getIfAnother(url string) {

}
