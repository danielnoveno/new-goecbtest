/*
    file:           configs/environment.go
    description:    Loader konfigurasi untuk env
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package configs

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

func init() {
	envCandidates := []string{".env", filepath.Join("..", ".env")}
	for _, envPath := range envCandidates {
		if err := godotenv.Load(envPath); err == nil {
			break
		}
	}
}

type Config struct {
	PublicHost             string
	Port                   string
	DBUser                 string
	DBPassword             string
	DBAddress              string
	DBName                 string
	DBHost                 string
	DBPort                 string
	SimoprdHost            string
	SimoprdPort            string
	SimoprdDatabase        string
	SimoprdUser            string
	SimoprdPassword        string
	BservHost              string
	BservPort              string
	BservDatabase          string
	BservUser              string
	BservPassword          string
}

type SimoConfig struct {
	AppEnv          string
	AppDebug        bool
	AppKey          string
	Title           string
	Version         string
	Description     string
	Icon            string
	EcbLocation     string
	EcbLineType     string
	EcbLineIds      string
	EcbWorkcenters  string
	EcbTacktime     int64
	EcbStateDefault string
	Theme           string
	EcbMode         string
}

func LoadConfig() Config {
	dbHost := GetEnvFirst([]string{"DB_HOST"}, "")
	dbPort := GetEnvFirst([]string{"DB_PORT"}, "")
	dbAddress := ""
	if dbHost != "" || dbPort != "" {
		dbAddress = fmt.Sprintf("%s:%s", dbHost, dbPort)
	}

	return Config{
		PublicHost:             GetEnvFirst([]string{"PUBLIC_HOST"}, ""),
		Port:                   GetEnvFirst([]string{"PORT"}, ""),
		DBUser:                 GetEnvFirst([]string{"DB_USERNAME", "DB_USER"}, ""),
		DBPassword:             GetEnvFirst([]string{"DB_PASSWORD"}, ""),
		DBHost:                 dbHost,
		DBPort:                 dbPort,
		DBAddress:              dbAddress,
		DBName:                 GetEnvFirst([]string{"DB_DATABASE", "DB_NAME"}, ""),
		SimoprdHost:            GetEnvFirst([]string{"DBSIMOPRD_HOST"}, ""),
		SimoprdPort:            GetEnvFirst([]string{"DBSIMOPRD_PORT"}, ""),
		SimoprdDatabase:        GetEnvFirst([]string{"DBSIMOPRD_DATABASE"}, ""),
		SimoprdUser:            GetEnvFirst([]string{"DBSIMOPRD_USERNAME"}, ""),
		SimoprdPassword:        GetEnvFirst([]string{"DBSIMOPRD_PASSWORD"}, ""),
		BservHost:              GetEnvFirst([]string{"DBBSERV_HOST"}, ""),
		BservPort:              GetEnvFirst([]string{"DBBSERV_PORT"}, ""),
		BservDatabase:          GetEnvFirst([]string{"DBBSERV_DATABASE"}, ""),
		BservUser:              GetEnvFirst([]string{"DBBSERV_USERNAME"}, ""),
		BservPassword:          GetEnvFirst([]string{"DBBSERV_PASSWORD"}, ""),
	}
}

func LoadSimoConfig() SimoConfig {
	return SimoConfig{
		AppEnv:          GetEnvFirst([]string{"APP_ENV", "SIMO_ENV"}, ""),
		AppDebug:        GetEnvAsBool("APP_DEBUG", false),
		AppKey:          GetEnvFirst([]string{"APP_KEY"}, ""),
		Title:           GetEnvFirst([]string{"APP_TITLE", "SIMO_TITLE"}, ""),
		Version:         GetEnvFirst([]string{"APP_VERSION", "SIMO_VERSION"}, ""),
		Description:     GetEnvFirst([]string{"APP_DESCRIPTION", "SIMO_DESCRIPTION"}, ""),
		Icon:            GetEnvFirst([]string{"APP_ICON", "SIMO_ICON"}, ""),
		EcbLocation:     GetEnvFirst([]string{"ECB_LOCATION", "SIMO_ECBLOCATION"}, ""),
		EcbLineType:     GetEnvFirst([]string{"ECB_LINE_TYPE", "SIMO_ECBLINETYPE"}, ""),
		EcbLineIds:      GetEnvFirst([]string{"ECB_LINEID", "ECB_LINEIDS", "SIMO_ECBLINEIDS"}, ""),
		EcbWorkcenters:  GetEnvFirst([]string{"ECB_WORKCENTER", "ECB_WORKCENTERS", "SIMO_ECBWORKCENTERS"}, ""),
		EcbTacktime:     GetEnvAsIntFirst([]string{"ECB_TACKTIME", "SIMO_ECBTACKTIME"}, 0),
		EcbStateDefault: GetEnvFirst([]string{"ECB_STATE_DEFAULT", "SIMO_ECBSTATEDEFAULT"}, ""),
		Theme:           GetEnvFirst([]string{"APP_THEME_DEFAULT", "SIMO_THEME"}, ""),
		EcbMode:         GetEnvFirst([]string{"ECB_MODE", "SIMO_ECBMODE"}, ""),
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

func GetEnvFirst(keys []string, fallback string) string {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
	}
	return fallback
}

func GetEnvAsIntFirst(keys []string, fallback int64) int64 {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			i, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fallback
			}
			return i
		}
	}
	return fallback
}

func GetEnvAsBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		return parsed
	}
	return fallback
}

func GetAdminPassword() string {
	return GetEnvFirst([]string{"ADMIN_MENU_PASSWORD"}, "ecb123")
}

func GetGpioPollIntervalMs() int64 {
	return GetEnvAsIntFirst([]string{"GPIO_POLL_INTERVAL_MS"}, 250)
}

func GetGpioAdaptivePolling() bool {
	return GetEnvAsBool("GPIO_ADAPTIVE_POLLING", true)
}

func GetSchedulerCleanMutexInterval() int64 {
	return GetEnvAsIntFirst([]string{"SCHEDULER_CLEAN_MUTEX_INTERVAL"}, 10)
}

func GetSchedulerPostDataInterval() int64 {
	return GetEnvAsIntFirst([]string{"SCHEDULER_POST_DATA_INTERVAL"}, 5)
}

func GetSchedulerSyncPoInterval() int64 {
	return GetEnvAsIntFirst([]string{"SCHEDULER_SYNC_PO_INTERVAL"}, 10)
}

func IsRaspberryPi() bool {
	if val, ok := os.LookupEnv("DB_IS_RASPBERRY_PI"); ok {
		parsed, err := strconv.ParseBool(val)
		if err == nil {
			return parsed
		}
	}

	if _, err := os.Stat("/proc/device-tree/model"); err == nil {
		return true
	}

	return false
}

func GetEnableResourceMonitor() bool {
	return GetEnvAsBool("ENABLE_RESOURCE_MONITOR", false)
}

func GetResourceMonitorInterval() int64 {
	return GetEnvAsIntFirst([]string{"RESOURCE_MONITOR_INTERVAL"}, 2)
}
