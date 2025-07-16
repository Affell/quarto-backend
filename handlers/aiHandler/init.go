package aiHandler

import "quarto/models"

func All(prefix string) []models.Route {
	return []models.Route{
		{
			Method:  "POST",
			Path:    prefix + "/solve",
			Handler: solve,
		},
	}
}
