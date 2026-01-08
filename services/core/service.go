/*
   file:           services/core/service.go
   description:    Kontroler inti untuk service
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package controllers

import (
	"fmt"
	"html"
	"net"
	"sort"
	"strings"
	"sync"

	"fyne.io/fyne/v2"

	"github.com/go-gorp/gorp"

	"go-ecb/app/types"
	"go-ecb/configs"
	"go-ecb/pkg/logging"
	"go-ecb/repository"
)

type Controller struct {
	DB         *gorp.DbMap
	SimoConfig configs.SimoConfig
	EnvConfig  configs.Config

	App                    fyne.App
	Window                 fyne.Window
	CurrentMenuId          int
	CurrentMenuAccessMode  int
	NavModes               []string
	navigationCache        *navigationCache
	navCacheErr            error
	navCacheOnce           sync.Once
	SimoMenu               string
	SimoBreadcrumb         string
	SimoSubMenus           []map[string]interface{}
	SimoCurrentMenuParents []int
	EcbLocation            string
	EcbLineNumber          int
	EcbLineType            string
	EcbLineIds             []string
	EcbWorkcenters         []string
	EcbTacktime            int64
	EcbStateDefault        string
	EcbMode                string
}

type navigationCache struct {
	byID     map[int]*types.Navigation
	byURL    map[string]*types.Navigation
	parent   map[int]int
	children map[int][]*types.Navigation
	roots    []*types.Navigation
}

// NewController adalah fungsi untuk baru pengendali.
func NewController(dbMap *gorp.DbMap, simoConfig configs.SimoConfig, envConfig configs.Config, app fyne.App, w fyne.Window) *Controller {
	c := &Controller{
		DB:         dbMap,
		SimoConfig: simoConfig,
		EnvConfig:  envConfig,
		App:        app,
		Window:     w,
	}
	c.init()
	return c
}

// init adalah fungsi untuk inisialisasi.
func (c *Controller) init() {
	c.SimoCurrentMenuParents = []int{}
	serverIPAddress := c.GetIPAddress()
	useConfig := true
	var ecbstation *types.EcbStation

	if serverIPAddress != "" {
		repo := repository.NewEcbStationRepository(c.DB)
		var err error
		ecbstation, err = repo.FindEcbStationByIP(serverIPAddress)
		if err != nil {
			logging.Logger().Warnf("Error querying Ecbstation: %v", err)
		} else if ecbstation == nil {
			logging.Logger().Infof("No Ecbstation found for IP: %s", serverIPAddress)
		} else {
			useConfig = false
		}
	}

	if useConfig {
		c.EcbLocation = c.SimoConfig.EcbLocation
		c.EcbLineType = c.SimoConfig.EcbLineType
		if c.EcbLineType == "sn-only-single" || c.EcbLineType == "refrig-single" || c.EcbLineType == "refrig-po-single" {
			c.EcbLineNumber = 1
		} else {
			c.EcbLineNumber = 2
		}
		c.EcbLineIds = strings.Split(c.SimoConfig.EcbLineIds, ",")
		c.EcbWorkcenters = strings.Split(c.SimoConfig.EcbWorkcenters, ",")
		c.EcbTacktime = c.SimoConfig.EcbTacktime
		c.EcbStateDefault = c.SimoConfig.EcbStateDefault

		c.EcbMode = c.SimoConfig.EcbMode
	} else {
		c.EcbLocation = ecbstation.Location
		c.EcbLineType = ecbstation.Linetype
		if c.EcbLineType == "sn-only-single" || c.EcbLineType == "refrig-single" || c.EcbLineType == "refrig-po-single" {
			c.EcbLineNumber = 1
		} else {
			c.EcbLineNumber = 2
		}
		c.EcbLineIds = strings.Split(ecbstation.Lineids, ",")
		c.EcbWorkcenters = strings.Split(ecbstation.Workcenters, ",")
		c.EcbTacktime = int64(ecbstation.Tacktime)
		c.EcbMode = ecbstation.Mode
	}
	if envLineType := strings.TrimSpace(c.SimoConfig.EcbLineType); envLineType != "" {
		c.EcbLineType = envLineType
		if c.EcbLineType == "sn-only-single" || c.EcbLineType == "refrig-single" || c.EcbLineType == "refrig-po-single" {
			c.EcbLineNumber = 1
		} else {
			c.EcbLineNumber = 2
		}
	}

	if err := c.ensureNavigationCache(); err != nil {
		logging.Logger().Errorf("[navigation] load failed: %v", err)
	}
}

// ensureNavigationCache adalah fungsi untuk memastikan navigasi cache.
func (c *Controller) ensureNavigationCache() error {
	c.navCacheOnce.Do(func() {
		c.navCacheErr = c.loadNavigationCache()
	})
	return c.navCacheErr
}

func (c *Controller) loadNavigationCache() error {
	if c.DB == nil {
		return fmt.Errorf("database connection is not configured")
	}

	repo := repository.NewNavigationRepository(c.DB)
	navigations, err := repo.GetAll()
	if err != nil {
		return err
	}

	cache := &navigationCache{
		byID:     make(map[int]*types.Navigation),
		byURL:    make(map[string]*types.Navigation),
		parent:   make(map[int]int),
		children: make(map[int][]*types.Navigation),
	}

	// Logic to populate navList from navigations slice
	var navList []*types.Navigation
	for _, nav := range navigations {
		// Use a local copy ensuring pointer safety if needed,
		// though nav is a value type here (Navigation struct).
		// Creating a pointer to a loop variable is safe in newer Go but explicit copy is safer.
		n := nav
		if n.Mode <= 0 {
			continue
		}
		cache.byID[n.ID] = &n
		cache.byURL[strings.ToLower(normalizePath(n.Url))] = &n
		navList = append(navList, &n)
	}

	sort.SliceStable(navList, func(i, j int) bool {
		if navList[i].Urutan == navList[j].Urutan {
			return navList[i].Title < navList[j].Title
		}
		return navList[i].Urutan < navList[j].Urutan
	})

	for _, nav := range navList {
		if nav.ParentId.Valid {
			parentID := int(nav.ParentId.Int64)
			cache.children[parentID] = append(cache.children[parentID], nav)
			cache.parent[nav.ID] = parentID
		} else {
			cache.roots = append(cache.roots, nav)
		}
	}

	for parentID, children := range cache.children {
		sort.SliceStable(children, func(i, j int) bool {
			if children[i].Urutan == children[j].Urutan {
				return children[i].Title < children[j].Title
			}
			return children[i].Urutan < children[j].Urutan
		})
		cache.children[parentID] = children
	}

	c.navigationCache = cache
	return nil
}

// normalizePath adalah fungsi untuk normalize jalur.
func normalizePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		path = "/"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	if len(path) > 1 && strings.HasSuffix(path, "/") {
		path = strings.TrimSuffix(path, "/")
	}
	return path
}

// assignCurrentMenu adalah fungsi untuk assign current menu.
func (c *Controller) assignCurrentMenu(nav *types.Navigation) {
	if nav == nil {
		return
	}
	c.CurrentMenuId = nav.ID
	c.SimoCurrentMenuParents = c.GetMenuParents(nav.ID)
	if breadcrumb, err := c.BuildBreadcrumb(); err == nil {
		c.SimoBreadcrumb = breadcrumb
	}
	if menu, err := c.BuildMainMenu(); err == nil {
		c.SimoMenu = menu
	}
	if subMenus, err := c.SubMenus(); err == nil {
		c.SimoSubMenus = subMenus
	}
}

// BuildMainMenu adalah fungsi untuk menyusun utama menu.
func (c *Controller) BuildMainMenu() (string, error) {
	if err := c.ensureNavigationCache(); err != nil {
		return "", err
	}
	c.CurrentMenuAccessMode = 0
	menu := c.buildMenu(c.navigationCache.roots)
	c.SimoMenu = menu
	return menu, nil
}

// buildMenu adalah fungsi untuk menyusun menu.
func (c *Controller) buildMenu(menu []*types.Navigation) string {
	if len(menu) == 0 {
		return ""
	}
	var builder strings.Builder
	for _, item := range menu {
		if item == nil || item.Mode <= 0 {
			continue
		}

		menuIcon := ""
		if item.Icon != "" {
			menuIcon = fmt.Sprintf(`<i class="fa %s"></i> `, html.EscapeString(item.Icon))
		}

		displayTitle := c.truncateTitle(item.Title, 10)
		isActive := ""
		if item.ID == c.CurrentMenuId {
			isActive = "active"
			c.CurrentMenuAccessMode = 1
		} else if c.containsParent(item.ID) {
			isActive = "active"
		}

		children := c.navigationCache.children[item.ID]
		if len(children) > 0 {
			builder.WriteString(fmt.Sprintf(`<li class="dropdown %s"><a href="%s" title="%s">%s%s <span class="caret"></span></a><ul class="dropdown-menu">`,
				isActive, fmt.Sprintf("/menu/%d", item.ID), html.EscapeString(item.Title), menuIcon, displayTitle))
			builder.WriteString(c.buildMenu(children))
			builder.WriteString("</ul></li>")
			continue
		}

		builder.WriteString(fmt.Sprintf(`<li class="%s"><a href="%s" title="%s">%s%s</a></li>`,
			isActive, fmt.Sprintf("/menu/%d", item.ID), html.EscapeString(item.Title), menuIcon, displayTitle))
	}
	return builder.String()
}

// BuildBreadcrumb adalah fungsi untuk menyusun breadcrumb.
func (c *Controller) BuildBreadcrumb() (string, error) {
	if err := c.ensureNavigationCache(); err != nil {
		return "", err
	}
	if c.CurrentMenuId == 0 {
		return "", nil
	}
	nav := c.navigationCache.byID[c.CurrentMenuId]
	if nav == nil {
		return "", nil
	}
	breadcrumb := fmt.Sprintf(`<li class="active">%s</li>`, html.EscapeString(nav.Title))
	if parentID, ok := c.navigationCache.parent[nav.ID]; ok {
		if parent := c.navigationCache.byID[parentID]; parent != nil {
			breadcrumb = c.buildBreadcrumbItem(parent, 0) + breadcrumb
		}
	}
	c.SimoBreadcrumb = breadcrumb
	return breadcrumb, nil
}

// buildBreadcrumbItem adalah fungsi untuk menyusun breadcrumb item.
func (c *Controller) buildBreadcrumbItem(item *types.Navigation, level int) string {
	if item == nil {
		return ""
	}
	breadcrumb := fmt.Sprintf(`<li><a href="%s">%s</a></li>`, fmt.Sprintf("/menu/%d", item.ID), html.EscapeString(item.Title))
	if parentID, ok := c.navigationCache.parent[item.ID]; ok {
		if parent := c.navigationCache.byID[parentID]; parent != nil {
			if level >= 3 {
				return "<li>..</li>" + breadcrumb
			}
			return c.buildBreadcrumbItem(parent, level+1) + breadcrumb
		}
	}
	return breadcrumb
}

// GetMenuParents adalah fungsi untuk mengambil menu parents.
func (c *Controller) GetMenuParents(id int) []int {
	if err := c.ensureNavigationCache(); err != nil {
		return nil
	}
	parentList := make([]int, 0)
	visited := make(map[int]struct{})
	current := id
	for {
		parentID, ok := c.navigationCache.parent[current]
		if !ok || parentID == 0 {
			break
		}
		if _, seen := visited[parentID]; seen {
			break
		}
		parentList = append(parentList, parentID)
		visited[parentID] = struct{}{}
		current = parentID
	}
	return parentList
}

// SubMenus adalah fungsi untuk sub menus.
func (c *Controller) SubMenus() ([]map[string]interface{}, error) {
	if err := c.ensureNavigationCache(); err != nil {
		return nil, err
	}
	if c.CurrentMenuId == 0 {
		return nil, nil
	}
	children := c.navigationCache.children[c.CurrentMenuId]
	submenus := make([]map[string]interface{}, 0, len(children))
	for _, submenu := range children {
		if submenu == nil || submenu.Mode <= 0 {
			continue
		}
		title := submenu.Title
		if len(c.navigationCache.children[submenu.ID]) > 0 {
			title = title + " >>"
		}
		submenus = append(submenus, map[string]interface{}{
			"id":     submenu.ID,
			"url":    fmt.Sprintf("/menu/%d", submenu.ID),
			"title":  title,
			"icon":   submenu.Icon,
			"urutan": submenu.Urutan,
		})
	}
	c.SimoSubMenus = submenus
	return submenus, nil
}

// containsParent adalah fungsi untuk contains parent.
func (c *Controller) containsParent(id int) bool {
	for _, parentID := range c.SimoCurrentMenuParents {
		if parentID == id {
			return true
		}
	}
	return false
}

// truncateTitle adalah fungsi untuk truncate title.
func (c *Controller) truncateTitle(title string, limit int) string {
	runes := []rune(title)
	if len(runes) <= limit {
		return title
	}
	return string(runes[:limit]) + ".."
}

// GetIPAddress adalah fungsi untuk mengambil ip address.
func (c *Controller) GetIPAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		logging.Logger().Errorf("Error enumerating network interfaces: %v", err)
		return ""
	}

	if ip := pickIPFromInterfaces(interfaces, false); ip != "" {
		return ip
	}
	return pickIPFromInterfaces(interfaces, true)
}

// isPrivateIPv4 adalah fungsi helper untuk mengecek private IPv4.
func isPrivateIPv4(ip net.IP) bool {
	if ip == nil {
		return false
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	switch {
	case ip4[0] == 10:
		return true
	case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
		return true
	case ip4[0] == 192 && ip4[1] == 168:
		return true
	default:
		return false
	}
}

// isLinkLocalIPv4 adalah fungsi helper untuk mengecek link-local IPv4.
func isLinkLocalIPv4(ip net.IP) bool {
	if ip == nil {
		return false
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	return ip4[0] == 169 && ip4[1] == 254
}

// pickIPFromInterfaces memilih IP dari interface yang tersedia dengan prioritas.
// Jika allowVirtual true, interface virtual ikut dipertimbangkan.
func pickIPFromInterfaces(interfaces []net.Interface, allowVirtual bool) string {
	var fallback string
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if !allowVirtual && isLikelyVirtualInterface(iface.Name) {
			continue
		}

		addrs, addrErr := iface.Addrs()
		if addrErr != nil {
			logging.Logger().Warnf("Error getting addresses for interface %s: %v", iface.Name, addrErr)
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ip4 := ip.To4()
			if ip4 == nil || ip4.IsLoopback() || isLinkLocalIPv4(ip4) {
				continue
			}

			if isPrivateIPv4(ip4) {
				return ip4.String()
			}

			if fallback == "" {
				fallback = ip4.String()
			}
		}
	}

	return fallback
}

// isLikelyVirtualInterface mencoba mendeteksi interface virtual berdasarkan nama umum.
func isLikelyVirtualInterface(name string) bool {
	lower := strings.ToLower(name)
	virtualHints := []string{
		"vmnet", "vmware", "vbox", "virtual", "veth", "v ethernet", "hyper-v", "loopback", "pseudo-interface", "tap", "tunnel", "ppp", "bridge",
	}
	for _, hint := range virtualHints {
		if strings.Contains(lower, hint) {
			return true
		}
	}
	return false
}
