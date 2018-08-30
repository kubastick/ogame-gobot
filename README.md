# Ogame GO-Bot
Simple OGame bot that interacts with your account.
## Motivation
I wanted to learn how to use selenium, so the first thing that came to my head is the good old ogame
## Build
```
cd $GOPATH/src
git clone https://github.com/kubastick/ogame-gobot OgameBot
cd OgameBot
go get ./...
go build OgameBot
```
## Download
https://github.com/kubastick/ogame-gobot/releases
## Configuration
Files:  
`config.toml` - General settings: Round time, optional slack bot id, etc.  
Example file content:
```
roundTime = 2_400_000 #Round time in miliseconds, see FAQ
headless = true #Use headless chrome (run in background)
seleniumAdress = "http://localhost:4444/wd/hub" #Selenium Chrome WebDriver adress
```
`./users/userfile.toml` - User data  
Example file content:
```
email = "user@email.com" #User email
password = "123456" #User password
server = "https://s155-pl.ogame.gameforge.com/game" #Server adress
serverButtonID = 11 #If your account have only 1 universe, leave this alone
```
Bot require Selenium Remote Chromedriver running (see config.toml)
## Abilities
- Upgrade Solar Plant if eneregy is too low
- Build first available building on first page
## How it works
In each round the bot tries to log in to the account and build available buildings
## FAQ
#### Can I get banned for using this?
Yes, I do not take any responsibility for this.
#### What is round
Round is the moment when the bot logs in to your account and checks the possibility of constructing buildings
#### Will the bot be able to attack players?
No, my motivation is not to spoil the game for other players, but to learn how to automate the browser

