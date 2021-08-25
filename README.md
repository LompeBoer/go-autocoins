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
- The program **overrides the existing coin list in WickHunter**, no need to pause the bot. **v0.6**
- The program **does not** blacklist coins that are in open positions.
- The program is **read only** for **v1.0** you have to **manually** change the coin list in the WickHunter bot

## Instructions:
- Drop the executable file and the json settings file into the same folder with your bot. **Make sure the folder is not located on your desktop but is a dedicated folder elsewhere on your drive.**
- Make sure you have WickHunter bot version **v0.6.2** or higher.
- Define the following in autoCoins.json file
  - **version**: set this to 1 when using WickHunter bot v1.0 or set to 0 when using v0.6 (default = 1).
  - **max1hrPercent**: maximum 1hr price change percentage (default = 5).
  - **max4hrPercent**: maximum 4hr price change percentage (default = 5).
  - **max24hrPercent**: maximum 24hr price change percentage (default = 10).
  - **cooldownHrs**: the number of 1hr candles into the past to check for the price changes. Example: if the number is 4 (default), the bot will quarantine coins that had a 1hr price change more than defined in _max1hrPercent_ within the past X _cooldownHrs_ (default = 4). Note: cooldown only applies to 1hr changes, not to ATH or 24hr price changes.
  - **minAthPercent**: minimum proximity to ATH in percent (default = 5). Note: due to Binance limitations, the ATH is only pulled from the last 20 months, so it's not a true All Time High, but ATH-ish.
  - **minAge**: minimum coin age in days (default = 14).
  - **refresh**: the period in minutes of how often to check (recommended minimum 15 mins due to possibility of over-running your API limit) (default = 15).
  - **discord**: (optional) your discord webhook.
  - **mentionOnError**: use @here mention on Discord when an error occurs. (default = true)
  - **blacklist**: permanently blacklisted coins.
  - **googleApiKey**: (optional) Google API Key to access [WH Pairs list - STP Todd](https://docs.google.com/spreadsheets/d/1XWadBbVkbdi5Ub7bFhCcAhqpHiQXBETbeTg644pkTdI/). [[*](#google-docs-api)]
  - **marginAssets**: list of margin assets to trade. Use empty list to not filter out any pairs (default = ["USDT"]).
  - **filters**: this controls which filters are used
    - **blacklist**: Enable/disable the _blacklist_ (default = true).
    - **margin_assets**: Enable/disable checking the _marginAssets_  (default = true).
    - **google_sheet_permitted**: Enable/disable the Google Sheet _WH Pairs list_ permitted list (default = false).
    - **google_sheet_safe**: Enable/disable the Google Sheet _WH Pairs list_ safe list (default = false).
    - **wickhunter_db**: Enable/disable using the default WickHunter bot coin list (default = false).

- Double-click on the executable or run it from the terminal.

## Startup flags
You can supply flags at startup. These are optional.  
- **-config=path**: path to the config file (default = autoCoins.json).
- **-noconfig**: use default settings without a config file
- **-storage=path**: path to the storage file for WickHunter bot (default = storage.db).
- **-version**: prints the current go-autocoins version.
- **-pairs**: set pairs to permitted from the Google Sheet Pairs List and exits the program (Note: first close WH bot)
- **-safepairs**: set safe pairs to permitted from the Google Sheet Pairs List and exits the program (Note: first close WH bot)

## Difference with PowerShell autoCoins
### Missing
- Writing to a log file.
- Geo check at startup.
- Update available check.

### Fixes
- Errors that resulted in the PowerShell script stopping:
  - Does/should not crash, but instead will skip the run and retry again after set _refresh_.
  - Sends a message to the Discord WebHook with the error, with the possibility to mention @here so you can immediately intervene.
- When Binance adds a new coin without enough historical data it skips the coin.
- If the Binance API returns an error the program continues and does not crash. 

### Added
- Use Binance API weight limit to pause sending requests so you not over-run the API limit.
- Use Google Docs API to read [WH Pairs list - STP Todd](https://docs.google.com/spreadsheets/d/1XWadBbVkbdi5Ub7bFhCcAhqpHiQXBETbeTg644pkTdI/edit#gid=1034827699). [[*](#google-docs-api)]
- Specify which margin assets to trade.
- Reading of WH v1.0 storage.db file. Currently it is not possible for WH bot to pickup changes made by AutoCoins. So when using _version_ 1 it only sends messages to Discord and does **not** update the WH bot. 

## Google Docs API
**Update** You can leave this empty since v0.9.9 of AutoCoins
Using the Google Docs API requires an API Key.  
See https://developers.google.com/docs/api/how-tos/authorizing#APIKey on how to create the key.  
(Optional) Restrict the access to the `Google Sheets API` and restrict to the IP Address where AutoCoins is running.  
  
When using this functionality the program will only use the pairs specified by either the permitted or safe account column.

## Compile from source
- Download and install Go from https://golang.org/  
- Clone the repo using git or download the source.
- Windows: Run `go build -o bin/go-autocoins.exe cmd/autocoins/main.go` from the project directory.
- Linux/MacOS: Run `make build` from the project directory.

## Suggestions and issues
Please use the [issues](https://github.com/LompeBoer/go-autocoins/issues) page to request features or report bugs.