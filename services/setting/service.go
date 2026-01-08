/*
   file:           services/setting/service.go
   description:    Layanan pengaturan untuk service
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package setting

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-gorp/gorp"

	"go-ecb/app/types"
	"go-ecb/configs"
	"go-ecb/pkg/logging"
	"go-ecb/repository"
	controllers "go-ecb/services/core"

	"fyne.io/fyne/v2"
)

type SettingController struct {
	*controllers.Controller
	ecbConfigStore repository.EcbConfigStore
}

// NewSettingController adalah fungsi untuk baru pengaturan pengendali.
func NewSettingController(db *gorp.DbMap, simoConfig configs.SimoConfig, envConfig configs.Config, app fyne.App, w fyne.Window) *SettingController {
	return &SettingController{
		Controller:     controllers.NewController(db, simoConfig, envConfig, app, w),
		ecbConfigStore: repository.NewEcbConfigRepository(db),
	}
}

type SettingPageData struct {
	CurrentMenuIcon        string
	CurrentMenuTitle       string
	CurrentMenuDescription string
	ServerIPAddress        string
	Simo3IPAddress         string
	UseWLAN                string
}

// GetSettingPageData adalah fungsi untuk mengambil pengaturan page data.
func (c *SettingController) GetSettingPageData() SettingPageData {
	serverIPAddress := c.Controller.GetIPAddress()

	// Initialize with default/fallback values.
	// These values will be used only if they are not found in the database.
	simo3IPAddress := "10.30.1.5"
	useWLAN := "no"

	configs, err := c.ecbConfigStore.FindEcbConfigsBySection("settings")
	if err != nil {
		logging.Logger().Warnf("Error finding EcbConfigs for section 'settings': %v", err)
	} else {
		for _, config := range configs {
			switch config.Variable {
			case "ipsimo3":
				simo3IPAddress = config.Value
			case "usewlan":
				useWLAN = config.Value
			}
		}
	}
	currentMenuIcon := "⚙️" 
	currentMenuTitle := "Setting ECB Station"
	currentMenuDescription := "Menu ini digunakan untuk menampilkan dan mengatur setting ECB Station ini."

	return SettingPageData{
		CurrentMenuIcon:        currentMenuIcon,
		CurrentMenuTitle:       currentMenuTitle,
		CurrentMenuDescription: currentMenuDescription,
		ServerIPAddress:        serverIPAddress,
		Simo3IPAddress:         simo3IPAddress,
		UseWLAN:                useWLAN,
	}
}

// UpdateECBSettings adalah fungsi untuk memperbarui ecb pengaturan.
func (c *SettingController) UpdateECBSettings(settings types.ECBSetting) error {
	c.saveConfigVariable("settings", "ipsimo3", settings.Simo3IPAddress)
	c.saveConfigVariable("settings", "iplocal", settings.ServerIPAddress)
	c.saveConfigVariable("settings", "usewlan", settings.UseWLAN)

	err := c.touchSettingFile(settings.ServerIPAddress, settings.Simo3IPAddress, settings.UseWLAN)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

// saveConfigVariable adalah fungsi untuk menyimpan konfigurasi variable.
func (c *SettingController) saveConfigVariable(section, variable, value string) {
	config, err := c.ecbConfigStore.FindEcbConfigBySectionAndVariable(section, variable)
	if err != nil && err.Error() == "ecbconfig not found" {
		newConfig := &types.EcbConfig{
			Section:   section,
			Variable:  variable,
			Value:     value,
			Ordering:  "0",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if createErr := c.ecbConfigStore.CreateEcbConfig(newConfig); createErr != nil {
			logging.Logger().Errorf("Error creating EcbConfig %s/%s: %v", section, variable, createErr)
		}
	} else if err != nil {
		logging.Logger().Errorf("Error finding EcbConfig %s/%s: %v", section, variable, err)
	} else {
		config.Value = value
		config.UpdatedAt = time.Now()
		if updateErr := c.ecbConfigStore.UpdateEcbConfig(config); updateErr != nil {
			logging.Logger().Errorf("Error updating EcbConfig %s/%s: %v", section, variable, updateErr)
		}
	}
}

// touchSettingFile adalah fungsi untuk touch pengaturan file.
func (c *SettingController) touchSettingFile(iplocal, ipsimo3, usewlan string) error {
	filePath := filepath.Join("storage", "app", "files", "changesetting")
	content := []byte(iplocal + "-" + ipsimo3 + "-" + usewlan)

	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	err := os.WriteFile(filePath, content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}
	logging.Logger().Infof("Simulated file touch: %s with content: %s", filePath, string(content))
	return nil
}

// UpdateMasterData adalah fungsi untuk memperbarui master data.
func (c *SettingController) UpdateMasterData() error {
	logging.Logger().Infof("Master data update initiated (placeholder).")
	return nil
}
