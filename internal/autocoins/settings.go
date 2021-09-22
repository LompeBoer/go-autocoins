package autocoins

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

type SettingsAutoCoins struct {
	Max1hrPercent  int `json:"max1hrPercent"`
	Max4hrPercent  int `json:"max4hrPercent"`
	Max24hrPercent int `json:"max24hrPercent"`
	CooldownHours  int `json:"cooldownHrs"`
	MinAthPercent  int `json:"minAthPercent"`
	MinAge         int `json:"minAge"`
	Refresh        int `json:"refresh"`
}

type SettingsFilterGoogleSheet struct {
	Enabled   bool     `json:"enabled"`
	Safe      bool     `json:"safe"`
	WhiteList []string `json:"whiteList"`
	APIKey    string   `json:"apiKey"`
}

type SettingsFilters struct {
	BlackList    []string                  `json:"blackList"`
	ExcludeList  []string                  `json:"excludeList"`
	MarginAssets []string                  `json:"marginAssets"`
	GoogleSheet  SettingsFilterGoogleSheet `json:"googleSheet"`
	WickHunterDB bool                      `json:"wickHunterDB"`
}

type SettingsDiscord struct {
	WebHook        string `json:"webHook"`
	MentionOnError bool   `json:"mentionOnError"`
}

type SettingsProxy struct {
	Address  string `json:"address"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Settings struct {
	configFilename string
	Version        int               `json:"version"`
	API            string            `json:"api"`
	Exchange       string            `json:"exchange"`
	AutoCoins      SettingsAutoCoins `json:"autoCoins"`
	Filters        SettingsFilters   `json:"filters"`
	Discord        SettingsDiscord   `json:"discord"`
	Proxy          SettingsProxy     `json:"proxy"`
}

func LoadConfig(file string) *Settings {
	s := Settings{}
	r := false
	if file != "" {
		r = s.LoadConfigFile(file)
	}
	if !r {
		s.LoadDefaultConfig()
	}

	s.ValidateSettings()
	s.PostProcess()
	s.configFilename = file
	return &s
}

func (s *Settings) ReloadConfig() *Settings {
	file := s.configFilename
	settings := LoadConfig(file)
	settings.ValidateSettings()
	settings.PostProcess()
	settings.configFilename = file
	return settings
}

func (s *Settings) LoadConfigFile(file string) bool {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		log.Printf("Config file '%s' does not exist.\n", file)
		return false
	}

	data, err := os.ReadFile(file)
	if err != nil {
		log.Printf("Error loading config file: %s\n", err.Error())
		return false
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		log.Printf("Error unmarshal config file: %s\n", err.Error())
		return false
	}

	*s = settings

	log.Printf("Loaded settings from config file '%s'\n", file)

	return true
}

func (s *Settings) LoadDefaultConfig() {
	*s = Settings{
		Version:  1,
		API:      "http://localhost:5001",
		Exchange: "binance",
		AutoCoins: SettingsAutoCoins{
			Max1hrPercent:  5,
			Max4hrPercent:  5,
			Max24hrPercent: 10,
			CooldownHours:  4,
			MinAthPercent:  5,
			MinAge:         14,
			Refresh:        15,
		},
		Filters: SettingsFilters{
			BlackList: []string{
				"BTCUSDT", "ETHUSDT", "YFIUSDT", "DEFIUSDT", "DOGEUSDT",
			},
			ExcludeList:  []string{},
			MarginAssets: []string{"USDT"},
			GoogleSheet: SettingsFilterGoogleSheet{
				Enabled:   false,
				Safe:      false,
				WhiteList: []string{},
				APIKey:    "",
			},
			WickHunterDB: true,
		},
		Discord: SettingsDiscord{
			WebHook:        "",
			MentionOnError: false,
		},
		Proxy: SettingsProxy{
			Address:  "",
			Username: "",
			Password: "",
		},
	}

	log.Println("Using default settings")
}

func (s *Settings) ValidateSettings() {
	if s.AutoCoins.Refresh < 1 {
		s.AutoCoins.Refresh = 1
	}
	if s.API == "" {
		log.Fatal("No API URL set in config file.")
	}

	s.Filters.WickHunterDB = true
}

func (s *Settings) PostProcess() {
	sort.Strings(s.Filters.BlackList)
	sort.Strings(s.Filters.ExcludeList)
	sort.Strings(s.Filters.GoogleSheet.WhiteList)
	sort.Strings(s.Filters.MarginAssets)
}
