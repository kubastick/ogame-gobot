package OgameController

import (
	"github.com/tebeka/selenium"
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

type ogameController struct {
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
