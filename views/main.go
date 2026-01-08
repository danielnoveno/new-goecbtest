/*
   file:           views/main.go
   description:    Antarmuka Fyne untuk main (Responsive Header Update)
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package views

import (
	// "errors"
	"fmt"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-gorp/gorp"
	_ "golang.org/x/image/webp"

	"go-ecb/app/types"
	"go-ecb/configs"
	"go-ecb/pkg/logging"
	"go-ecb/repository"
	"go-ecb/services/about"

	// flashsvc "go-ecb/services/flash"
	"go-ecb/services/maintenance"
	"go-ecb/services/setting"
	"go-ecb/services/system"
	"go-ecb/utils"
	"go-ecb/views/ecb"
	double "go-ecb/views/ecb/double"
	navigation "go-ecb/views/ecb/navigation"
	single "go-ecb/views/ecb/single"
	customtheme "go-ecb/views/theme"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	xlayout "fyne.io/x/fyne/layout"
)

type headerPalette struct {
	HeaderStart color.Color
	HeaderEnd   color.Color
	Accent      color.Color
}

// BuildMainWindow adalah fungsi untuk menyusun utama window.
func BuildMainWindow(a fyne.App, w fyne.Window, dbmap *gorp.DbMap, pinController maintenance.PinController) *types.AppState {

	var setBody func(fyne.CanvasObject)
	var body fyne.CanvasObject
	state := &types.AppState{
		Window:  w,
		DbMap:   dbmap,
		SetBody: func(obj fyne.CanvasObject) { setBody(obj) },
	}
	// flashService := flashsvc.NewNotifier()
	// state.Flash = flashService
	// flashService.Subscribe(func(msg types.FlashMessage) {
	// 	if state.Window == nil {
	// 		return
	// 	}
	// 	switch msg.Level {
	// 	case types.FlashLevelError:
	// 		dialog.ShowError(errors.New(msg.Body), state.Window)
	// 	default:
	// 		dialog.ShowInformation(msg.Title, msg.Body, state.Window)
	// 	}
	// })
	var applyIndicator func()
	var setMode func(string)

	var childMu sync.Mutex
	var activeChild fyne.Window

	state.TryOpenWindow = func(title string, setup func(fyne.Window)) bool {
		childMu.Lock()
		// if activeChild != nil {
		// 	current := activeChild
		// 	childMu.Unlock()

		// 	dialog.ShowInformation("Form masih terbuka, tutup terlebih dahulu",
		// 		"Selesaikan atau tutup form sekarang sebelum membuat halmman baru.", w)

		// 	current.RequestFocus()
		// 	return false
		// }

		childWin := a.NewWindow(title)
		activeChild = childWin
		childMu.Unlock()

		childWin.SetOnClosed(func() {
			childMu.Lock()
			if activeChild == childWin {
				activeChild = nil
			}
			childMu.Unlock()
		})

		setup(childWin)
		childWin.Show()
		childWin.RequestFocus()
		return true
	}

	const (
		headerHeight = 64
	)

	locText := widget.NewLabel("Location: -")
	locText.TextStyle = fyne.TextStyle{Bold: true}

	modeBinding := binding.NewString()
	if err := modeBinding.Set("Mode: -"); err != nil {
		logging.Logger().Warnf("Warning: failed to initialize mode binding: %v\n", err)
	}
	modeText := widget.NewLabelWithData(modeBinding)
	modeText.TextStyle = fyne.TextStyle{Bold: true}

	formatLineTypeValue := func(value string) string {
		lineType := strings.TrimSpace(value)
		if lineType == "" {
			return "-"
		}
		return strings.ToUpper(lineType)
	}
	lineTypeBinding := binding.NewString()
	if err := lineTypeBinding.Set(formatLineTypeValue("")); err != nil {
		logging.Logger().Warnf("Warning: failed to initialize line type binding: %v\n", err)
	}
	lineTypeText := widget.NewLabelWithData(lineTypeBinding)
	lineTypeText.TextStyle = fyne.TextStyle{Bold: true}

	clockBinding := binding.NewString()
	if err := clockBinding.Set(time.Now().Format("15:04:05")); err != nil {
		logging.Logger().Warnf("Warning: failed to initialize clock binding: %v\n", err)
	}
	clockText := widget.NewLabelWithData(clockBinding)
	clockText.Alignment = fyne.TextAlignTrailing
	clockText.TextStyle = fyne.TextStyle{Bold: true}

	type navThemeOption struct {
		palette customtheme.Palette
		header  headerPalette
	}

	fallbackPalette := customtheme.DefaultPalette()
	defaultHeader := headerPalette{
		HeaderStart: fallbackPalette.HeaderColor,
		HeaderEnd:   fallbackPalette.HeaderColor,
		Accent:      fallbackPalette.Accent,
	}

	menuIconTint := fallbackPalette.Text
	var refreshMenuIcons func()

	loadLogoResource := func(relPath string) fyne.Resource {
		path := utils.ResolvePath(relPath)
		res, err := fyne.LoadResourceFromPath(path)
		if err == nil {
			return res
		}
		logging.Logger().Warnf("Warning: gagal memuat logo dari %s: %v\n", path, err)
		return theme.FyneLogo()
	}
	expandedLogoResource := loadLogoResource(filepath.Join("assets", "Logo-polytron.webp"))
	collapsedLogoResource := loadLogoResource(filepath.Join("assets", "logo-nb.webp"))
	logoImage := canvas.NewImageFromResource(expandedLogoResource)
	logoImage.FillMode = canvas.ImageFillContain
	logoImage.SetMinSize(fyne.NewSize(140, 40))

	sep := func() fyne.CanvasObject {
		r := canvas.NewRectangle(theme.DisabledColor())
		r.SetMinSize(fyne.NewSize(1, 16))
		return container.NewCenter(r)
	}

	toggleBtn := widget.NewButtonWithIcon("", theme.MenuIcon(), nil)
	toggleBtn.Importance = widget.LowImportance

	logoBackground := canvas.NewRectangle(color.Transparent)
	logoBackground.SetMinSize(fyne.NewSize(220, headerHeight))

	logoBlock := container.NewStack(
		logoBackground,
		container.NewCenter(logoImage),
	)

	setMode = func(mode string) {
		mode = strings.TrimSpace(mode)
		if mode == "" {
			mode = "LIVE"
		}
		state.Mode = mode
		if err := modeBinding.Set(state.Mode); err != nil {
			logging.Logger().Warnf("Warning: failed to set mode binding: %v\n", err)
		}
	}
	ecb.RegisterModeChangeHandler(setMode)

	themeEntries := map[string]navThemeOption{}
	var themeLabels []string

	themeRepo := repository.NewThemeRepository(dbmap)
	dbThemes, err := themeRepo.GetThemes()
	if err != nil {
		logging.Logger().Warnf("Warning: gagal memuat tema: %v\n", err)
	}

	if len(dbThemes) == 0 {
		themeEntries[fallbackPalette.Name] = navThemeOption{
			palette: fallbackPalette,
			header:  defaultHeader,
		}
		themeLabels = append(themeLabels, fallbackPalette.Name)
	} else {
		for _, rec := range dbThemes {
			pal := customtheme.PaletteFromRecord(rec)
			label := pal.Name
			if label == "" {
				label = fmt.Sprintf("Theme %d", rec.ID)
			}
			if _, exists := themeEntries[label]; exists {
				continue
			}
			start := customtheme.ColorFromHex(rec.HeaderStart, fallbackPalette.HeaderColor)
			end := customtheme.ColorFromHex(rec.HeaderEnd, start)
			themeEntries[label] = navThemeOption{
				palette: pal,
				header: headerPalette{
					HeaderStart: start,
					HeaderEnd:   end,
					Accent:      pal.Accent,
				},
			}
			themeLabels = append(themeLabels, label)
		}
	}
	if len(themeLabels) == 0 {
		themeEntries[fallbackPalette.Name] = navThemeOption{
			palette: fallbackPalette,
			header:  defaultHeader,
		}
		themeLabels = append(themeLabels, fallbackPalette.Name)
	}

	var selectedTheme string
	if len(themeLabels) > 0 {
		selectedTheme = themeLabels[0]
	}

	var applyThemeSelection func(string)
	themeSelect := widget.NewSelect(themeLabels, func(name string) {
		if applyThemeSelection != nil {
			applyThemeSelection(name)
		}
	})
	themeSelect.PlaceHolder = "Theme"

	headerBg := canvas.NewRectangle(defaultHeader.HeaderStart)
	infoRow := container.NewHBox(
		locText,
		sep(),
		modeText,
		lineTypeText,
		sep(),
		clockText,
	)

	infoWrapper := container.NewVBox(layout.NewSpacer(), container.NewHBox(infoRow), layout.NewSpacer())

	responsiveInfoCenter := xlayout.NewResponsiveLayout(
		xlayout.Responsive(container.NewPadded(infoWrapper), 1, 0.9, 0.8, 0.8),
	)
	centerBlock := container.NewHBox(
		layout.NewSpacer(),
		container.NewPadded(responsiveInfoCenter),
		layout.NewSpacer(),
	)

	headerContent := container.NewBorder(
		nil,
		nil,
		container.NewVBox(layout.NewSpacer(), container.NewPadded(toggleBtn), layout.NewSpacer()),
		container.NewVBox(layout.NewSpacer(), container.NewPadded(themeSelect), layout.NewSpacer()),
		centerBlock,
	)

	header := container.NewStack(
		headerBg,
		container.NewPadded(headerContent),
	)

	ticker := time.NewTicker(time.Second)
	go func() {
		for t := range ticker.C {
			if err := clockBinding.Set(t.Format("15:04:05")); err != nil {
				logging.Logger().Warnf("Warning: gagal memperbarui jam (%v)\n", err)
			}
		}
	}()

	simoConfig := configs.LoadSimoConfig()
	if err := lineTypeBinding.Set(formatLineTypeValue(simoConfig.EcbLineType)); err != nil {
		logging.Logger().Warnf("Warning: failed to update line type binding: %v\n", err)
	}
	envConfig := configs.LoadConfig()
	adminMenuPassword := configs.GetAdminPassword()
	w.SetTitle(fmt.Sprintf("%s %s", simoConfig.Title, simoConfig.Version))

	if cfgTheme := strings.TrimSpace(simoConfig.Theme); cfgTheme != "" {
		if _, ok := themeEntries[cfgTheme]; ok {
			selectedTheme = cfgTheme
		}
	}

	aboutController := about.NewController(dbmap, simoConfig, envConfig, a, w)
	settingController := setting.NewSettingController(dbmap, simoConfig, envConfig, a, w)

	const (
		sidebarExpandedWidth  = 220
		sidebarCollapsedWidth = 50
	)

	sidebarCollapsed := true

	sidebarBg := canvas.NewRectangle(fallbackPalette.HeaderColor)

	menuList := container.NewVBox()
	sidebarFooter := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	sidebarContent := container.NewBorder(
		container.NewVBox(logoBlock, widget.NewSeparator()),
		container.NewPadded(sidebarFooter),
		nil,
		nil,
		menuList,
	)

	var menuButtons []*widget.Button
	indicatorRects := []*canvas.Rectangle{}
	var sizingRects []*canvas.Rectangle
	selectedMenuIndex := 0
	var setActiveMenu func(int)
	var selectMenu func(int)

	applyIndicator = func() {
		for i, r := range indicatorRects {
			if i == selectedMenuIndex {
				r.FillColor = color.NRGBA{R: 255, G: 255, B: 255, A: 30}
			} else {
				r.FillColor = color.Transparent
			}
			r.Refresh()
		}
	}

	menuRenderers := map[string]func() fyne.CanvasObject{
		"ecb-test": func() fyne.CanvasObject {
			switch strings.ToLower(strings.TrimSpace(simoConfig.EcbLineType)) {
			case "sn-only-double":
				return double.SnOnlyDoubleScreen(w, dbmap, simoConfig)
			case "refrig-single":
				return single.RefrigSingleScreen(w, dbmap, simoConfig)
			case "refrig-po-single":
				return single.RefrigPoSingleScreen(w, dbmap, simoConfig)
			case "refrig-double":
				return double.RefrigDoubleScreen(w, dbmap, simoConfig)
			case "refrig-po-double":
				return double.RefrigPoDoubleScreen(w, dbmap, simoConfig)
			default:
				return single.SnOnlySingleScreen(w, dbmap, simoConfig)
			}
		},
		"settings": func() fyne.CanvasObject {
			d := settingController.GetSettingPageData()
			return ecb.SettingPage(
				w,
				d,
				settingController.UpdateECBSettings,
				settingController.UpdateMasterData,
			)
		},
		"maintenance": func() fyne.CanvasObject {
			return ecb.MaintenanceUI(setMode, state.Mode, w, pinController)
		},
		"about": func() fyne.CanvasObject {
			return aboutController.Index()
		},
		"shutdown": func() fyne.CanvasObject {
			info := widget.NewLabel("Kirim perintah shutdown hanya jika perangkat fisik terkoneksi.")
			info.Wrapping = fyne.TextWrapWord
			runSystemCommandAsync(w, system.Shutdown)

			return container.NewVBox(
				widget.NewLabelWithStyle("Shutdown", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				info,
				layout.NewSpacer(),
			)
		},
		"reboot": func() fyne.CanvasObject {
			info := widget.NewLabel("Reboot akan me-reset sesi dan subsystem.")
			info.Wrapping = fyne.TextWrapWord
			runSystemCommandAsync(w, system.Reboot)

			return container.NewVBox(
				widget.NewLabelWithStyle("Reboot", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
				info,
				layout.NewSpacer(),
			)
		},
	}

	protectedMenuRoutes := map[string]struct{}{
		"settings":    {},
		"maintenance": {},
	}

	navigationRepo := repository.NewNavigationRepository(dbmap)
	rootNavigations, err := navigationRepo.FindRootNavigations()
	if err != nil {
		dialog.ShowError(err, w)
	}
	menu := buildMenuFromNavigationItems(rootNavigations, menuRenderers)
	content := container.NewMax()
	contentBackground := canvas.NewRectangle(fallbackPalette.Background)
	contentPanel := container.NewMax(contentBackground, container.NewPadded(content))
	contentScroll := container.NewVScroll(contentPanel)

	footerBg := canvas.NewRectangle(fallbackPalette.Footer)
	footerBg.SetMinSize(fyne.NewSize(0, 32))

	footerLabel := canvas.NewText(
		"PT. Hartono Istana Teknologi © 2025",
		fallbackPalette.Text,
	)
	footerLabel.TextSize = 10
	footerLabel.Alignment = fyne.TextAlignCenter

	footer := container.NewStack(
		footerBg,
		container.NewCenter(footerLabel),
	)

	applyThemeSelection = func(name string) {
		entry, ok := themeEntries[name]
		if !ok {
			return
		}
		menuIconTint = entry.palette.Text

		headerBg.FillColor = entry.header.HeaderStart
		headerBg.Refresh()

		sidebarBg.FillColor = entry.palette.HeaderColor
		sidebarBg.Refresh()

		logoBackground.FillColor = color.Transparent
		logoBackground.Refresh()

		contentBackground.FillColor = entry.palette.Background
		contentBackground.Refresh()

		footerBg.FillColor = entry.palette.HeaderColor
		footerBg.Refresh()

		footerLabel.Color = entry.palette.Text
		footerLabel.Refresh()

		a.Settings().SetTheme(customtheme.New(entry.palette))

		locText.Refresh()
		modeText.Refresh()
		lineTypeText.Refresh()
		clockText.Refresh()
		toggleBtn.Refresh()
		if applyIndicator != nil {
			applyIndicator()
		}
		if refreshMenuIcons != nil {
			refreshMenuIcons()
		}
	}
	if selectedTheme != "" {
		applyThemeSelection(selectedTheme)
		themeSelect.SetSelected(selectedTheme)
	}

	contentColumn := container.NewBorder(
		header,
		footer,
		nil,
		nil,
		contentScroll,
	)

	state.Menu = menu
	menuButtons = make([]*widget.Button, len(menu))
	indicatorRects = make([]*canvas.Rectangle, len(menu))
	sizingRects = make([]*canvas.Rectangle, len(menu))

	for idx := range menu {
		i := idx

		indicator := canvas.NewRectangle(color.Transparent)
		indicatorRects[i] = indicator

		btn := widget.NewButtonWithIcon(menu[i].Title, resolveMenuIcon(menu[i].Icon, menuIconTint), func() {
			selectMenu(i)
		})
		btn.Alignment = widget.ButtonAlignLeading
		btn.Importance = widget.LowImportance

		menuButtons[i] = btn

		sizingRect := canvas.NewRectangle(color.Transparent)
		sizingRect.SetMinSize(fyne.NewSize(sidebarExpandedWidth, 42))
		sizingRects[i] = sizingRect
		btnWrapper := container.NewStack(indicator, sizingRect, btn)
		menuList.Add(container.NewPadded(btnWrapper))
	}

	refreshMenuIcons = func() {
		for idx, btn := range menuButtons {
			if btn == nil {
				continue
			}
			btn.SetIcon(resolveMenuIcon(menu[idx].Icon, menuIconTint))
			btn.Refresh()
		}
	}
	refreshMenuIcons()

	setActiveMenu = func(i int) {
		if i < 0 || i >= len(menu) {
			return
		}
		selectedMenuIndex = i
		applyIndicator()

		content.Objects = []fyne.CanvasObject{menu[i].Show()}
		content.Refresh()
	}
	promptMenuPassword := func(i int) {
		if i < 0 || i >= len(menu) {
			return
		}
		passwordEntry := widget.NewPasswordEntry()
		content := container.NewVBox(
			widget.NewLabel("Masukkan password administrator:"),
			passwordEntry,
		)

		confirmDialog := dialog.NewCustomConfirm(fmt.Sprintf("Password untuk %s", menu[i].Title),
			"Buka", "Batal", content,
			func(ok bool) {
				if !ok {
					return
				}
				if strings.TrimSpace(passwordEntry.Text) != adminMenuPassword {
					dialog.ShowError(fmt.Errorf("password salah untuk menu %s", menu[i].Title), w)
					return
				}
				setActiveMenu(i)
			}, w)

		passwordEntry.OnSubmitted = func(_ string) {
			confirmDialog.Confirm()
		}

		confirmDialog.Show()
		if canvas := w.Canvas(); canvas != nil {
			canvas.Focus(passwordEntry)
		}
	}

	selectMenu = func(i int) {
		if i < 0 || i >= len(menu) {
			return
		}
		key := strings.ToLower(strings.TrimSpace(menu[i].Key))
		if _, guard := protectedMenuRoutes[key]; guard {
			promptMenuPassword(i)
			return
		}
		setActiveMenu(i)
	}

	navigation.RegisterNavigationHandler(func(route string) {
		target := strings.ToLower(strings.TrimSpace(route))
		if target == "" {
			return
		}
		for idx := range menu {
			if strings.ToLower(strings.TrimSpace(menu[idx].Key)) == target {
				selectMenu(idx)
				return
			}
		}
	})

	if len(menu) > 0 {
		selectMenu(0)
	}

	setMenuButtonsWidth := func(collapsed bool) {
		targetWidth := sidebarExpandedWidth
		if collapsed {
			targetWidth = sidebarCollapsedWidth
		}
		for _, rect := range sizingRects {
			rect.SetMinSize(fyne.NewSize(float32(targetWidth), 42))
			rect.Refresh()
		}
	}

	sidebar := container.NewStack(
		sidebarBg,
		sidebarContent,
	)

	body = container.NewBorder(nil, nil, sidebar, nil, contentColumn)
	root := body
	mainStack := container.NewStack(root)

	syncSidebarAppearance := func() {
		var logoSize fyne.Size
		if sidebarCollapsed {
			logoImage.Resource = collapsedLogoResource
			logoSize = fyne.NewSize(sidebarCollapsedWidth-12, 32)

			logoBackground.SetMinSize(fyne.NewSize(sidebarCollapsedWidth, headerHeight))
			toggleBtn.SetIcon(theme.MenuIcon())
			sidebarBg.SetMinSize(fyne.NewSize(sidebarCollapsedWidth, 0))
			setMenuButtonsWidth(true)

			for _, btn := range menuButtons {
				btn.SetText("")
				btn.Alignment = widget.ButtonAlignCenter
			}
		} else {
			logoImage.Resource = expandedLogoResource
			logoSize = fyne.NewSize(140, 40)

			logoBackground.SetMinSize(fyne.NewSize(sidebarExpandedWidth, headerHeight))
			toggleBtn.SetIcon(theme.MenuIcon())
			sidebarBg.SetMinSize(fyne.NewSize(sidebarExpandedWidth, 0))
			setMenuButtonsWidth(false)

			for i, btn := range menuButtons {
				btn.SetText(state.Menu[i].Title)
				btn.Alignment = widget.ButtonAlignLeading
			}
		}

		logoImage.SetMinSize(logoSize)
		logoImage.Refresh()
		sidebarBg.Refresh()
		applyIndicator()
		logoBlock.Refresh()
	}

	syncSidebarAppearance()

	toggleBtn.OnTapped = func() {
		sidebarCollapsed = !sidebarCollapsed
		syncSidebarAppearance()
	}
	setBody = func(obj fyne.CanvasObject) {
		content.Objects = []fyne.CanvasObject{obj}
		content.Refresh()
	}
	state.Location = strings.TrimSpace(simoConfig.EcbLocation)
	if state.Location == "" {
		state.Location = "-"
	}
	initialMode := simoConfig.EcbMode
	if strings.TrimSpace(initialMode) == "" {
		initialMode = "LIVE"
	}
	setMode(initialMode)

	locText.SetText(state.Location)

	w.SetContent(mainStack)
	w.Resize(fyne.NewSize(1024, 700))
	startWindowSizeLogger(w)

	var closeDialogVisible bool
	w.SetCloseIntercept(func() {
		if closeDialogVisible {
			return
		}
		closeDialogVisible = true

		dialog.ShowConfirm("Exit confirmation", "Are you sure you want to exit the application?", func(confirm bool) {
			closeDialogVisible = false
			if confirm {
				w.SetCloseIntercept(nil)
				w.Close()
			}
		}, w)
	})

	if canvas := w.Canvas(); canvas != nil {
		prevTyped := canvas.OnTypedKey()
		canvas.SetOnTypedKey(func(ev *fyne.KeyEvent) {
			if ev != nil {
				entryFocused := false
				if focused := canvas.Focused(); focused != nil {
					_, entryFocused = focused.(*widget.Entry)
				}

				switch ev.Name {
				case fyne.KeyReturn, fyne.KeyEnter:
					if activateDialogPrimaryAction(w) {
						return
					}
				case fyne.KeyY:
					if !entryFocused && activateDialogActionByLabel(w, "yes") {
						return
					}
				case fyne.KeyB:
					if !entryFocused && activateDialogActionByLabel(w, "buka") {
						return
					}
				}
			}
			if prevTyped != nil {
				prevTyped(ev)
			}
		})
	}

	return state
}

func buildMenuFromNavigationItems(items []*types.Navigation, renderers map[string]func() fyne.CanvasObject) []types.MenuItem {
	menus := make([]types.MenuItem, 0, len(items))
	for _, nav := range items {
		if nav == nil {
			continue
		}
		routeKey := strings.TrimSpace(nav.Route)
		if routeKey == "" {
			routeKey = strings.TrimSpace(nav.Url)
		}
		renderer := renderers[routeKey]
		if renderer == nil {
			renderer = func(title string) func() fyne.CanvasObject {
				return func() fyne.CanvasObject {
					return container.NewCenter(widget.NewLabel(fmt.Sprintf("Menu \"%s\" belum tersedia.", title)))
				}
			}(nav.Title)
		}
		menus = append(menus, types.MenuItem{
			Title: nav.Title,
			Key:   routeKey,
			Icon:  nav.Icon,
			Show:  renderer,
			NavID: nav.ID,
		})
	}
	return menus
}

func resolveMenuIcon(name string, tint color.Color) fyne.Resource {
	switch strings.ToLower(name) {

	case "info", "about":
		return theme.InfoIcon()

	case "settings":
		return theme.SettingsIcon()

	case "list":
		return theme.ListIcon()

	case "performance":
		return loadCustomMenuIcon(filepath.Join("assets", "maintenance.svg"), tint)

	case "shutdown":
		return theme.CancelIcon()

	case "reboot", "refresh":
		return theme.ViewRefreshIcon()

	case "bolt":
		return loadCustomMenuIcon(filepath.Join("assets", "bolt.svg"), tint)

	case "power-off":
		return loadCustomMenuIcon(filepath.Join("assets", "power-off.svg"), tint)

	default:
		return theme.MenuIcon()
	}
}

var tintedMenuIconCache = map[string]fyne.Resource{}

func loadCustomMenuIcon(relPath string, tint color.Color) fyne.Resource {
	if tint == nil {
		tint = color.Black
	}
	hexColor := colorToHex(tint)
	cacheKey := fmt.Sprintf("%s|%s", relPath, hexColor)
	resolvedPath := utils.ResolvePath(relPath)
	if res, ok := tintedMenuIconCache[cacheKey]; ok {
		logging.Logger().Debugf("Debug: menu icon %s (tint %s) diambil dari cache [%s]\n", relPath, hexColor, resolvedPath)
		return res
	}

	data, err := os.ReadFile(resolvedPath)
	if err != nil {
		logging.Logger().Warnf("Warning: gagal memuat menu icon %s (%s): %v\n", relPath, resolvedPath, err)
		res := theme.MenuIcon()
		tintedMenuIconCache[cacheKey] = res
		return res
	}

	svgContent := sanitizeSVGForTint(string(data))
	svgContent = strings.ReplaceAll(svgContent, "#000000", hexColor)

	resourceName := fmt.Sprintf("%s-%s", filepath.Base(relPath), hexColor)
	res := fyne.NewStaticResource(resourceName, []byte(svgContent))
	logging.Logger().Infof("Info: menu icon %s dimuat (%d bytes) dengan tint %s dari %s\n", relPath, len(svgContent), hexColor, resolvedPath)
	tintedMenuIconCache[cacheKey] = res
	return res
}

func colorToHex(c color.Color) string {
	if c == nil {
		return "#000000"
	}
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02X%02X%02X", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

func sanitizeSVGForTint(svg string) string {
	svg = strings.ReplaceAll(svg, ` style="--darkreader-inline-stroke: var(--darkreader-text-000000, #e8e6e3);"`, "")
	svg = strings.ReplaceAll(svg, ` style="--darkreader-inline-fill: var(--darkreader-background-000000, #000000);"`, "")
	svg = strings.ReplaceAll(svg, ` data-darkreader-inline-stroke=""`, "")
	svg = strings.ReplaceAll(svg, ` data-darkreader-inline-fill=""`, "")
	return svg
}

func startWindowSizeLogger(w fyne.Window) {
	if w == nil {
		return
	}
	canvas := w.Canvas()
	if canvas == nil {
		return
	}

	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		last := canvas.Size()
		logging.Logger().Infof("[window] size initial %.0fx%.0f", last.Width, last.Height)

		for range ticker.C {
			current := canvas.Size()
			if current.Width != last.Width || current.Height != last.Height {
				last = current
				logging.Logger().Infof("[window] size changed %.0fx%.0f", current.Width, current.Height)
			}
		}
	}()
}

// runSystemCommandAsync adalah fungsi untuk menjalankan system command async.
func runSystemCommandAsync(w fyne.Window, action func() error) {
	go func() {
		if err := action(); err != nil {
			fyne.Do(func() {
				dialog.ShowError(err, w)
			})
		}
	}()
}
