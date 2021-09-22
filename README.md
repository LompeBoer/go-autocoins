# go-autocoins 
[![build](https://github.com/LompeBoer/go-autocoins/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/LompeBoer/go-autocoins/actions/workflows/go.yml) 

Based on [autoCoins](https://github.com/daisy613/autoCoins) by [daisy613](https://github.com/daisy613)  
See her page for more great WH scripts.  

## Overview
- [What it does](#what-it-does)
- [Instructions](#instructions)
- [Startup flags](#startup-flags)
- [Difference with PowerShell autoCoins](#difference-with-powershell-autocoins)
- [Google Docs API](#google-docs-api)
- [Compile from source](#compile-from-source)
- [Suggestions and issues](#suggestions-and-issues)

## What it does:
- This program allows you to avoid most pumps/dumps by dynamically controlling the coin-list in WickHunter bot to blacklist\un-blacklist coins based on the following conditions:
  - combination of 1hr, 4hr and 24hr price percentage changes.
  - proximity to All Time High.
  - minimum coin age.
- The program **overrides the existing coin list in WickHunter**, no need to pause the bot.
- The program **does not** blacklist coins that are in open positions.

## Instructions:
- Download latest version from the [releases](https://github.com/LompeBoer/go-autocoins/releases) page.
  - Windows `go-autocoins.windows-amd64.zip`
  - Linux `go-autocoins.linux-amd64.tar.gz`
  - Mac OS X `go-autocoins.darwin-amd64.tar.gz`
- Extract the archive which contains the executable.
- Drop the executable file and the json settings file `autoCoins.json` into the same folder with your bot.
- Make sure you have WickHunter bot version **v1.1.4** or higher.
- Define the following in autoCoins.json file
  - **version**: set this to 1 when using WickHunter bot v1.1.4 or higher (default = 1).
  - **api**: use `http://localhost:5001` for Binance and `http://localhost:5000` for ByBit (default = `http://localhost:5001`)
  - **exchange**: the exchange to get data from. Only Binance is available for now (default = binance).
  - **autoCoins**:
    - **max1hrPercent**: maximum 1hr price change percentage (default = 5).
    - **max4hrPercent**: maximum 4hr price change percentage (default = 5).
    - **max24hrPercent**: maximum 24hr price change percentage (default = 10).
    - **cooldownHrs**: the number of 1hr candles into the past to check for the price changes. Example: if the number is 4 (default), the bot will quarantine coins that had a 1hr price change more than defined in _max1hrPercent_ within the past X _cooldownHrs_ (default = 4). Note: cooldown only applies to 1hr changes, not to ATH or 24hr price changes.
    - **minAthPercent**: minimum proximity to ATH in percent (default = 5). Note: due to Binance limitations, the ATH is only pulled from the last 20 months, so it's not a true All Time High, but ATH-ish.
    - **minAge**: minimum coin age in days (default = 14).
    - **refresh**: the period in minutes of how often to check (recommended minimum 15 mins due to possibility of over-running your API limit) (default = 15).
  - **filters**: this controls which filters are used
    - **blackList**: permanently blacklisted coins.
    - **excludeList**: coins on this list will not be quarantined. (default = [])
    - **marginAssets**: list of margin assets to trade, [read more](#margin-assets) (default = ["USDT"]).
    - **googleSheet**: [read more](#google-sheet)
      - **enabled**: enable/disable using the Google Sheet _WH Pairs list - STP Todd_ (default = false).
      - **safe**: use the column _SAFE ACCOUNT_ (default = false).
      - **whiteList**: overrides the sheet setting and treats these as being "safe". (default = [])
      - **apiKey**: (optional) Google API Key, [read more](#google-docs-api).
  - **discord**:
    - **webHook**: (optional) your discord webhook.
    - **mentionOnError**: use @here mention on Discord when an error occurs. (default = true)
  - **proxy**:
    - **address**: (optional) IP proxy and port to use (example "http://25.12.124.35:2763"). Leave blank if no proxy used ("").
    - **username**: (optional) proxy user.
    - **password**: (optional) proxy password.
- Make sure Wick Hunter bot is open.
- Double-click on the executable or run it from the terminal/commandprompt.

## Startup flags
You can supply flags at startup. These are optional.  
- **-config=path**: path to the config file (default = autoCoins.json).
- **-noconfig**: use default settings without a config file
- **-storage=path**: path to the storage file for WickHunter bot (default = storage.db).
- **-version**: prints the current go-autocoins version.
- **-pairs**: set pairs to permitted from the Google Sheet Pairs List and exits the program (Note: WH has to be running)
- **-safepairs**: set safe pairs to permitted from the Google Sheet Pairs List and exits the program (Note: WH has to be running)

## Filters
### WickHunter DB
Only coins in the WickHunter database will be used. 
If Binance has coins not present in WickHunter these new coins will not be processed.  
This filter is always active and can not be disabled.  

### Black List
This will permanently exclude coins from being traded.  

### Margin Assets
Only trade the margin asset specified.  
For example when using "USDT", then BTC**USDT** will be traded but BTC**BUSD** will not.  
Leave empty to trade all margin assets. `marginAssets: []`    
WickHunter itself only uses USDT so this is a redundant setting.  

### Google Sheet
The following Google Sheet is being used: [WH Pairs list - STP Todd](https://docs.google.com/spreadsheets/d/1XWadBbVkbdi5Ub7bFhCcAhqpHiQXBETbeTg644pkTdI/)  
  
When using this filter the program will only use the pairs specified by either the permitted or safe account column.  
Pairs added to _whitelist_ overrides the sheet setting and treats them as being "safe".  

### Exclude List
Coins on this list will not be quarantined.  

## Difference with PowerShell autoCoins
- There are changes to the `autoCoins.json` file. Make sure to update this file.  
- Wick Hunter has to be open.  

### Missing
- Writing to a log file.
- Geo check at startup.
- Support for v0.6.6

## Google Docs API  
Using the Google Docs API Key is optional. If you do not use one the functionality will be the same.  
  
See https://developers.google.com/docs/api/how-tos/authorizing#APIKey on how to create the key.  
(Optional) Restrict the access to the `Google Sheets API` and restrict to the IP Address where AutoCoins is running.  
  
## Compile from source
- Download and install Go from https://golang.org/  
- Clone the repo using git or download the source.
- Windows: Run `go build -o bin/go-autocoins.exe cmd/autocoins/*.go` from the project directory.
- Linux/MacOS: Run `make build` from the project directory.

## Suggestions and issues
Please use the [issues](https://github.com/LompeBoer/go-autocoins/issues) page to request features or report bugs.