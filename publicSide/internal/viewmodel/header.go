// Package viewmodel содержит структуры, которые используются для передачи данных в шаблоны (views).
package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
)

// HeaderViewModel предоставляет все необходимые URL-адреса для навигации в шапке сайта.
type HeaderViewModel struct {
	HomeRoute       string
	CategoriesRoute string
	ProfileRoute    string
	LoginRoute      string
	LogoutRoute     string
	RegRoute        string
}

// NewHeader создает новую модель представления для шапки сайта.
func NewHeader() *HeaderViewModel {
	return &HeaderViewModel{
		HomeRoute:       routing.RouteHome,
		CategoriesRoute: routing.RouteCategories,
		ProfileRoute:    routing.ExternalServiceRouteProfile,
		LoginRoute:      routing.RouteLogin,
		LogoutRoute:     routing.RouteLogout,
		RegRoute:        routing.RouteReg,
	}
}
