/*
    file:           configs/environment_loader.go
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

// init adalah fungsi untuk inisialisasi.
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
	JWTSecret              string
	JWTExpirationInSeconds int64
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

// LoadConfig adalah fungsi untuk memuat konfigurasi.
func LoadConfig() Config {
	dbHost := GetEnvFirst([]string{"DB_HOST"}, "")
	dbPort := GetEnvFirst([]string{"DB_PORT"}, "")
	dbAddress := ""
	if dbHost != "" || dbPort != "" {
		dbAddress = fmt.Sprintf("%s:%s", dbHost, dbPort)
	}

	return Config{
		PublicHost:             GetEnv("PUBLIC_HOST", ""),
		Port:                   GetEnv("PORT", ""),
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
		BservHost:              GetEnv("DBBSERV_HOST", ""),
		BservPort:              GetEnv("DBBSERV_PORT", ""),
		BservDatabase:          GetEnv("DBBSERV_DATABASE", ""),
		BservUser:              GetEnv("DBBSERV_USERNAME", ""),
		BservPassword:          GetEnv("DBBSERV_PASSWORD", ""),
		JWTSecret:              GetEnv("JWT_SECRET", ""),
		JWTExpirationInSeconds: GetEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 0),
	}
}

// LoadSimoConfig adalah fungsi untuk memuat SIMO konfigurasi.
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

// GetEnv adalah fungsi untuk mengambil env.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

// GetEnvFirst adalah fungsi untuk mengambil env first.
func GetEnvFirst(keys []string, fallback string) string {
	for _, key := range keys {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
	}
	return fallback
}

// GetEnvAsInt adalah fungsi untuk mengambil env as int.
func GetEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}

// GetEnvAsIntFirst adalah fungsi untuk mengambil env as int first.
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

// GetEnvAsBool adalah fungsi untuk mengambil env as bool.
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

// GetAdminPassword adalah fungsi untuk mengambil admin password.
func GetAdminPassword() string {
	return GetEnv("ADMIN_MENU_PASSWORD", "ecb123")
}

// ========== Optimization Configurations ==========

// GetGpioPollIntervalMs adalah fungsi untuk mengambil GPIO polling interval dalam milliseconds.
// Default: 250ms untuk backward compatibility.
func GetGpioPollIntervalMs() int64 {
	return GetEnvAsInt("GPIO_POLL_INTERVAL_MS", 250)
}

// GetGpioAdaptivePolling adalah fungsi untuk mengaktifkan adaptive polling.
// Jika true, polling akan melambat saat tidak ada perubahan state.
// Default: true untuk menghemat CPU.
func GetGpioAdaptivePolling() bool {
	return GetEnvAsBool("GPIO_ADAPTIVE_POLLING", true)
}

// GetSchedulerCleanMutexInterval adalah fungsi untuk mengambil CleanMutex interval dalam menit.
// Default: 10 menit.
func GetSchedulerCleanMutexInterval() int64 {
	return GetEnvAsInt("SCHEDULER_CLEAN_MUTEX_INTERVAL", 10)
}

// GetSchedulerPostDataInterval adalah fungsi untuk mengambil PostEcbData interval dalam menit.
// Default: 5 menit.
func GetSchedulerPostDataInterval() int64 {
	return GetEnvAsInt("SCHEDULER_POST_DATA_INTERVAL", 5)
}

// GetSchedulerSyncPoInterval adalah fungsi untuk mengambil SyncEcbPo interval dalam menit.
// Default: 10 menit.
func GetSchedulerSyncPoInterval() int64 {
	return GetEnvAsInt("SCHEDULER_SYNC_PO_INTERVAL", 10)
}

// IsRaspberryPi adalah fungsi untuk mendeteksi apakah aplikasi berjalan di Raspberry Pi.
// Bisa di-override dengan environment variable DB_IS_RASPBERRY_PI.
func IsRaspberryPi() bool {
	// Cek manual override dari env
	if val, ok := os.LookupEnv("DB_IS_RASPBERRY_PI"); ok {
		parsed, err := strconv.ParseBool(val)
		if err == nil {
			return parsed
		}
	}

	// Auto-detect: cek apakah running di ARM architecture
	// Bisa diperluas dengan cek /proc/cpuinfo atau model file
	if _, err := os.Stat("/proc/device-tree/model"); err == nil {
		return true
	}

	return false
}

// GetEnableResourceMonitor adalah fungsi untuk mengaktifkan resource monitoring dashboard.
// Default: false untuk tidak menambah overhead di production.
func GetEnableResourceMonitor() bool {
	return GetEnvAsBool("ENABLE_RESOURCE_MONITOR", false)
}

// GetResourceMonitorInterval adalah fungsi untuk mengambil monitoring update interval dalam detik.
// Default: 2 detik.
func GetResourceMonitorInterval() int64 {
	return GetEnvAsInt("RESOURCE_MONITOR_INTERVAL", 2)
}
