package sappress

import (
    "encoding/json"
    "fmt"
    "os"
)

const configPath = "config.json"

type Config struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Token    string `json:"token"`
    UserKey  string `json:"userkey"`
}

// LoadConfig loads the config from config.json
func LoadConfig() (*Config, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, err
    }

    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}

// SaveConfig saves the entire config object to file
func SaveConfig(cfg *Config) error {
    data, err := json.MarshalIndent(cfg, "", "  ")
    if err != nil {
        return err
    }
    return os.WriteFile(configPath, data, 0644)
}

// UpdateConfigField updates a single field (like "token") and saves the config
func UpdateConfigField(field string, value string) error {
    cfg, err := LoadConfig()
    if err != nil {
        return err
    }

    switch field {
    case "email":
        cfg.Email = value
    case "password":
        cfg.Password = value
    case "token":
        cfg.Token = value
    case "userkey":
        cfg.UserKey = value
    default:
        return fmt.Errorf("unknown field: %s", field)
    }

    return SaveConfig(cfg)
}
