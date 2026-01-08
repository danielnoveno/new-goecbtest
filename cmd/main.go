/*
   file:           cmd/main.go
   description:    Entrypoint desktop untuk main
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package main

import (
	"embed"
	"os"
	"path/filepath"

	// "runtime/debug"
	"strings"

	"go-ecb/configs"
	"go-ecb/db"
	"go-ecb/pkg/logging"
	"go-ecb/pkg/monitor"
	"go-ecb/repository"
	maintenance "go-ecb/services/maintenance"

	// "go-ecb/services/adminer"
	"go-ecb/services/gpio"
	"go-ecb/services/station"
	scheduler "go-ecb/task"
	"go-ecb/utils"
	"go-ecb/views"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"github.com/gofrs/flock"
	_ "golang.org/x/image/webp"
)

var translations embed.FS

func main() {
	// defer func() {
	// 	if r := recover(); r != nil {
	// 		log.Printf("[panic] main recovered: %v\n%s", r, debug.Stack())
	// 	}
	// }()

	simoConfig := configs.LoadSimoConfig()
	logger := logging.Init(simoConfig.AppDebug)
	logger.Debugw("logger initialized for debugging", "debug", simoConfig.AppDebug)
	// defer logging.Sync()

	dbmap, err := db.InitDb()
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbmap.Db.Close()

	monitor.PrintSystemInfo()
	monitor.Start()

	// Initialize Station Service
	stationRepo := repository.NewEcbStationRepository(dbmap)
	stationService := station.NewStationService(stationRepo)
	if _, err := stationService.Initialize(); err != nil {
		logger.Errorf("Failed to initialize station: %v", err)
	}

	lockPath := filepath.Join(os.TempDir(), "go-ecb.lock")
	appLock := flock.New(lockPath)
	locked, err := appLock.TryLock()
	if err != nil {
		logger.Fatalf("Failed to lock app: %v", err)
	}
	if !locked {
		logger.Warnf("Aplikasi sudah berjalan")
		return
	}
	defer appLock.Unlock()

	a := app.New()

	// iconPath := strings.TrimSpace(simoConfig.Icon)
	// if iconPath == "" {
	iconPath := filepath.Join("assets", "logo-nb.webp")
	// }
	iconPath = utils.ResolvePath(iconPath)
	icon, err := fyne.LoadResourceFromPath(iconPath)
	if err != nil {
		// log.Printf("Warning: error laod icon %s: %v", iconPath, err)
		icon = theme.FyneLogo()
	}
	a.SetIcon(icon)
	windowTitle := strings.TrimSpace(simoConfig.Title)
	if windowTitle == "" {
		windowTitle = "Go ECB Test"
	}
	// if simoConfig.Version != "" {
	// 	windowTitle = windowTitle + " " + simoConfig.Version
	// }
	w := a.NewWindow(windowTitle)

	pinService := maintenance.NewPinConfigService(repository.NewEcbConfigRepository(dbmap))
	if pinService == nil {
		logger.Warnf("maintenance pin service disabled")
	} else if err := pinService.Refresh(); err != nil {
		logger.Errorf("Failed to load maintenance pins: %v", err)
	}

	gpio.InitializeHardware()
	gpio.StartEcbStateUpdater(dbmap)
	
	state := views.BuildMainWindow(a, w, dbmap, pinService)
	scheduler.Start()
	_ = state

	// w.SetFullScreen(true)
	w.SetMaster()
	w.ShowAndRun()
}
