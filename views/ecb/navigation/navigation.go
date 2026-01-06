/*
   file:           views/ecb/navigation/navigation.go
   description:    Navigation helpers for ECB views
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package navigation

var (
	navigationHandler func(route string)
)

// RegisterNavigationHandler saves a navigation callback so other components can request a route change.
func RegisterNavigationHandler(fn func(route string)) {
	navigationHandler = fn
}

// NavigateToRoute triggers a registered navigation handler. Returns true when a handler was invoked.
func NavigateToRoute(route string) bool {
	if navigationHandler == nil || route == "" {
		return false
	}
	navigationHandler(route)
	return true
}
