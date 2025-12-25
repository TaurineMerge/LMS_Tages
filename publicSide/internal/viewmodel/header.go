package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/routing"
)

type HeaderViewModel struct {
	HomeRoute       string
	CategoriesRoute string
	ProfileRoute    string
	LoginRoute      string
	LogoutRoute     string
	RegRoute        string
}

func NewHeader() *HeaderViewModel {
	return &HeaderViewModel{
		HomeRoute: routing.RouteHome,
		CategoriesRoute: routing.RouteCategories,
		ProfileRoute: routing.ExternalServiceRouteProfile,
		LoginRoute: routing.RouteLogin,
		LogoutRoute: routing.RouteLogout,
		RegRoute: routing.RouteReg,
	}
}