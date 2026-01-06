/*
    file:           services/about/service.go
    description:    Layanan tentang aplikasi untuk service
    created:        220711663@students.uajy.ac.id 04-11-2025
*/

package about

import (
	"go-ecb/configs"
	controllers "go-ecb/services/core"
	"go-ecb/views/ecb"

	"github.com/go-gorp/gorp"

	"fyne.io/fyne/v2"
)

type Controller struct {
	*controllers.Controller 
}

// NewController adalah fungsi untuk baru pengendali.
func NewController(dbMap *gorp.DbMap, simoConfig configs.SimoConfig, envConfig configs.Config, a fyne.App, w fyne.Window) *Controller {
	return &Controller{
		Controller: controllers.NewController(dbMap, simoConfig, envConfig, a, w),
	}
}

// Index adalah fungsi untuk indeks.
func (c *Controller) Index() fyne.CanvasObject {
	return ecb.AboutPage(
		"info",
		"ECB Test on desktop app based Raspberry PI3",
		"202509221",
		"220711663",
	)
}

// func (c *Controller) Index() fyne.CanvasObject {
// 	nav := c.GetMenuObject("about") 
// 	return ecb.AboutPage(
// 		nav.Icon,
// 		"ECB Test on desktop app based raspberry pi3",                                  
// 		"202509221",                                   
// 		"220711663",                 				   
// 	)
// }

// func (c *Controller) GetMenuObject(menuName string) *types.Navigation {
// 	if menuName == "about" {
// 		createdAt, _ := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")
// 		updatedAt, _ := time.Parse("2006-01-02 15:04:05", "2023-01-01 00:00:00")

// 		return &types.Navigation{
// 			ID:        1,
// 			ParentId:  sql.NullInt64{Int64: 0, Valid: true}, 
// 			Icon:      "info",                              
// 			Title:     "About",
// 			Url:       "/about",
// 			Mode:      1,
// 			Urutan:    1,
// 			CreatedAt: createdAt,
// 			UpdatedAt: updatedAt,
// 		}
// 	}
// 	return nil
// }
