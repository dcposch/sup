package sup

import (
	"fmt"
	"log"
	"net/http"
	"time"
    "encoding/json"
)

//
// USER ROUTE
//

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
        getStatusHandler(w, r)
    } else if r.Method == "POST" {
		postStatusHandler(w, r)
	}
}


// GET /api/status/<user> to get statuses
func getStatusHandler(w http.ResponseWriter, r *http.Request) {
    user := r.URL.Path[len("/api/status/"):]
    statuses, err := GetStatuses(user)
    if err != nil {
        msg := fmt.Sprintf("Couldn't load %s: %v", user, err)
        log.Println(msg)
        http.Error(w, msg, http.StatusBadRequest)
    }
    jsonBytes, err := json.Marshal(statuses)
    if err != nil {
        msg := fmt.Sprintf("Invalid data for %s: %v", user, err)
        log.Println(msg)
        http.Error(w, msg, http.StatusInternalServerError)
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonBytes)
}


// POST /api/status to create a new status
func postStatusHandler(w http.ResponseWriter, r *http.Request) {
	status := new(Status)
    status.IP = r.Header.Get("X-Real-IP")
    if status.IP == "" {
        status.IP = r.RemoteAddr
    }
    status.CreateTime = time.Now()
    json.NewDecoder(r.Body).Decode(status)
	validateUser(status.User)
	validateTags(status.Tags)

	log.Printf("Woot! New status %s %s %s %s\n", status.IP, status.User, status.Tags, status.Description)

    err := SaveStatus(status)
	if err != nil {
        msg := fmt.Sprintf("Couldn't save status :( %v", err)
        http.Error(w, msg, http.StatusInternalServerError)
	}
}
