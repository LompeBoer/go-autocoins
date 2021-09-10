package autocoins

import (
	"encoding/json"
	"log"
	"os"
)

type SettingsFilters struct {
	Blacklist            bool `json:"blacklist"`
	MarginAssets         bool `json:"marginAssets"`
	GoogleSheetPermitted bool `json:"googleSheetPermitted"`
	GoogleSheetSafe      bool `json:"googleSheetSafe"`
	WickHunterDB         bool `json:"wickhunterDB"`
}

type Settings struct {
	configFilename string
	Version        int             `json:"version"`
	API            string          `json:"api"`
	Max1hrPercent  int             `json:"max1hrPercent"`
	Max4hrPercent  int             `json:"max4hrPercent"`
	Max24hrPercent int             `json:"max24hrPercent"`
	MinAthPercent  int             `json:"minAthPercent"`
	MinAge         int             `json:"minAge"`
	Refresh        int             `json:"refresh"`
	Proxy          string          `json:"proxy"`
	ProxyUser      string          `json:"proxyUser"`
	ProxyPass      string          `json:"proxyPass"`
	Discord        string          `json:"discord"`
	MentionOnError bool            `json:"mentionOnError"`
	BlackList      []string        `json:"blackList"`
	CooldownHours  int             `json:"cooldownHrs"`
	GoogleApiKey   string          `json:"googleApiKey"`
	MarginAssets   []string        `json:"marginAssets"`
	Filters        SettingsFilters `json:"filters"`
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
	s.configFilename = file
	return &s
}

func (s *Settings) ReloadConfig() *Settings {
	file := s.configFilename
	settings := LoadConfig(file)
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

	var settings []Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		log.Printf("Error unmarshal config file: %s\n", err.Error())
		return false
	}

	if len(settings) == 0 {
		log.Printf("No settings found in the config file '%s'\n", file)
		return false
	} else if len(settings) > 1 {
		log.Printf("Multiple settings found in config file '%s', using the first one\n", file)
	}

	*s = settings[0]

	log.Printf("Loaded settings from config file '%s'\n", file)

	return true
}

func (s *Settings) LoadDefaultConfig() {
	*s = Settings{
		Version:        1,
		API:            "http://localhost:5001",
		Max1hrPercent:  5,
		Max4hrPercent:  5,
		Max24hrPercent: 10,
		MinAthPercent:  5,
		MinAge:         14,
		Refresh:        15,
		Proxy:          "",
		ProxyUser:      "",
		ProxyPass:      "",
		Discord:        "",
		BlackList: []string{
			"BTCUSDT", "ETHUSDT", "YFIUSDT", "DEFIUSDT", "DOGEUSDT",
		},
		CooldownHours:  4,
		MarginAssets:   []string{"USDT"},
		MentionOnError: false,
		Filters: SettingsFilters{
			Blacklist:            true,
			MarginAssets:         true,
			GoogleSheetPermitted: false,
			GoogleSheetSafe:      false,
			WickHunterDB:         false,
		},
	}

	log.Println("Using default settings")
}

func (s *Settings) ValidateSettings() {
	if s.Version == 0 && len(s.MarginAssets) == 0 {
		s.MarginAssets = []string{"USDT"}
	}
	if s.Refresh < 1 {
		s.Refresh = 1
	}
	if s.Version == 1 && s.API == "" {
		// TODO: disable write.
		log.Println("No API URL set in config file.")
	}
}
