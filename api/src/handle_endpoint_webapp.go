package src

import (
	"encoding/json"
	"net/http"
)

// does user with unid x exist ? returns boolean
func HandleEndpointDoesUserExistByUnid(app *App) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// get the {unid} value from the HTTP route
		unid := r.PathValue("unid")

		// basic validation
		if unid == "" {
			WriteErrs(w, http.StatusBadRequest, "unid is required")
			return
		}

		// fetch the user through sqlc
		doesUserExist, err := app.Queries.DoesUserExistByUnid(
			r.Context(),
			unid,
		)

		// unexpected database error
		if err != nil {
			WriteErrs(w, http.StatusInternalServerError, "something went wrong")
			return
		}

		// return domain data only
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{
			"result": doesUserExist,
		})

	}
}
