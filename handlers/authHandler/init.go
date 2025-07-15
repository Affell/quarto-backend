package authHandler

import "quarto/models"

func All(prefix string) (routes []models.Route) {
	routes = append(routes, models.Route{
		Path:    prefix + "/login",
		Method:  "POST",
		Handler: login,
	})

	routes = append(routes, models.Route{
		Path:    prefix + "/logout",
		Method:  "POST",
		Handler: logout,
	})

	routes = append(routes, models.Route{
		Path:    prefix + "/signup",
		Method:  "POST",
		Handler: signup,
	})

	routes = append(routes, models.Route{
		Path:    prefix + "/signout",
		Method:  "POST",
		Handler: signout,
	})

	routes = append(routes, models.Route{
		Path:    prefix + "/me",
		Method:  "GET",
		Handler: me,
	})

	routes = append(routes, models.Route{
		Path:    prefix + "/recover",
		Method:  "POST",
		Handler: recover,
	})

	routes = append(routes, models.Route{
		Path:    prefix + "/reset_password",
		Method:  "POST",
		Handler: resetPassword,
	})

	return
}
