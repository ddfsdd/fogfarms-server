package dashboard

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func MakeHTTPHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/dashboard", populateDashboard).
		Methods("POST").
		Schemes("http")

	router.HandleFunc("/dashboard/toggle_auto", toggleAuto).
		Methods("POST").
		Schemes("http")

	router.HandleFunc("/dashboard/set_env_param", setEnvironmentParameters).
		Methods("POST").
		Schemes("http")

	router.HandleFunc("/dashboard/reset_timer", resetTimer).
		Methods("POST").
		Schemes("http")

	router.HandleFunc("/dashboard/update_device_status", updateDeviceStatus).
		Methods("POST").
		Schemes("http")
	ahandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://localhost:3000", "https://25.22.245.97:3000"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	}).Handler(router)
	return ahandler
}

func populateDashboard(w http.ResponseWriter, r *http.Request) {
	PopulateDashboard(w, r)
}

func toggleAuto(w http.ResponseWriter, r *http.Request) {
	ToggleAuto(w, r)
}

func setEnvironmentParameters(w http.ResponseWriter, r *http.Request) {
	SetEnvironmentParameters(w, r)
}

func resetTimer(w http.ResponseWriter, r *http.Request) {
	ResetTimer(w, r)
}

func updateDeviceStatus(w http.ResponseWriter, r *http.Request) {
	UpdateDeviceStatus(w, r)
}
