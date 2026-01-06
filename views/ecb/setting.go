/*
   file:           views/ecb/setting.go
   description:    Layar ECB untuk setting
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package ecb

import (
	"fmt"
	"image/color"
	"strings"

	"go-ecb/app/types"
	"go-ecb/services/setting"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// SettingPage adalah fungsi untuk pengaturan page.
func SettingPage(w fyne.Window, data setting.SettingPageData, onSave func(types.ECBSetting) error, onUpdateMaster func() error) fyne.CanvasObject {
	headerIcon := canvas.NewText(data.CurrentMenuIcon, color.Black)
	headerIcon.TextStyle.Bold = true
	headerTitle := canvas.NewText(data.CurrentMenuTitle, color.Black)
	headerTitle.TextStyle.Bold = true

	header := container.NewHBox(
		headerIcon,
		widget.NewLabel("|"),
		headerTitle,
	)

	description := widget.NewLabel(data.CurrentMenuDescription)
	description.Wrapping = fyne.TextWrapWord

	currentIPLabel := widget.NewLabel(fmt.Sprintf("Current Machine IP Address: %s", data.ServerIPAddress))
	currentIPLabel.Wrapping = fyne.TextWrapWord

	updateMasterDataButton := widget.NewButton("update master data", func() {
		if onUpdateMaster == nil {
			dialog.ShowInformation("Info", "Fitur master data belum aktif.", w)
			return
		}
		if err := onUpdateMaster(); err != nil {
			dialog.ShowError(err, w)
			return
		}
		dialog.ShowInformation("Master data", "Proses sinkron master data masuk ke antrean lokal.", w)
	})

	formServerEntry := widget.NewEntry()
	formServerEntry.SetText(data.ServerIPAddress)

	formSimoEntry := widget.NewEntry()
	formSimoEntry.SetText(data.Simo3IPAddress)

	useWLANCheck := widget.NewCheck("Gunakan koneksi WiFi (Raspberry Pi3)", nil)
	useWLANCheck.SetChecked(strings.EqualFold(data.UseWLAN, "yes"))

	saveButton := widget.NewButton("Simpan setting", func() {
		if onSave == nil {
			dialog.ShowInformation("Info", "Fungsi penyimpanan belum tersedia.", w)
			return
		}

		serverValue := strings.TrimSpace(formServerEntry.Text)
		simoValue := strings.TrimSpace(formSimoEntry.Text)
		useWLANValue := "no"
		if useWLANCheck.Checked {
			useWLANValue = "yes"
		}

		reviewText := fmt.Sprintf("Review:\nIP Address mesin ini:%s\nIP Address server SIMO3:%s\nmenggunakan WLAN:%s", serverValue, simoValue, useWLANValue)
		dialog.ShowConfirm("Review setting", reviewText+"\nApakah setting akan disimpan?", func(save bool) {
			if !save {
				dialog.ShowInformation("Setting dibatalkan", "setting dibatalkan.", w)
				return
			}

			if err := onSave(types.ECBSetting{
				ServerIPAddress: serverValue,
				Simo3IPAddress:  simoValue,
				UseWLAN:         useWLANValue,
			}); err != nil {
				dialog.ShowError(err, w)
				return
			}
			dialog.ShowInformation("Setting diproses", "setting segera diproses setelah proses restart beberapa saat lagi.", w)
		}, w)
	})

	formCard := widget.NewCard(
		"Setting mesin",
		"Isi form berikut untuk memperbarui konfigurasi ECB Station.",
		container.NewVBox(
			widget.NewForm(
				widget.NewFormItem("IP Address mesin ini", formServerEntry),
				widget.NewFormItem("IP Address server SIMO3", formSimoEntry),
			),
			useWLANCheck,
			widget.NewLabel("Jika tidak menggunakan Raspberry Pi3 atau tidak ingin WiFi, pastikan kabel LAN tersambung."),
		),
	)

	buttonContainer := container.NewHBox(updateMasterDataButton, saveButton)

	bodyCard := widget.NewCard(
		"",
		"",
		container.NewVBox(description, widget.NewSeparator(), currentIPLabel, widget.NewSeparator(), formCard, widget.NewSeparator(), buttonContainer),
	)

	return container.NewBorder(
		nil,
		nil,
		nil,
		nil,
		container.NewVBox(header, bodyCard),
	)
}
