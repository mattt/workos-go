package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dewski/workos/auditlog"
)

type user struct {
	Email      string
	DatabaseID int
}

func (u user) ToAuditableName() string {
	return u.Email
}

func (u user) ToAuditableID() string {
	return strconv.Itoa(u.DatabaseID)
}

func main() {
	auditlog.SetMetadata(map[string]interface{}{
		"environment": "development",
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, req *http.Request) {
		u := user{
			Email:      "me@domain.com",
			DatabaseID: 1,
		}

		event := auditlog.NewEventWithHTTP("user.login", auditlog.Create, req)
		event.SetActor(u)
		event.SetTarget(u)
		err := event.Publish()
		if err != nil {
			// Call out to sentry
			fmt.Println("Had a problem writing this event")
		}

		body, _ := json.Marshal(event)
		fmt.Fprintf(w, string(body))
	})

	http.HandleFunc("/event", func(w http.ResponseWriter, req *http.Request) {
		resp, err := auditlog.Find("someid")
		if err != nil {
			// Call out to sentry
			fmt.Println("Had a problem with request")
			fmt.Println(err)
		}

		body, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(body))
	})

	http.HandleFunc("/events", func(w http.ResponseWriter, req *http.Request) {
		resp, err := auditlog.FindAll(auditlog.EventsRequestParams{
			End:    time.Now(),
			Action: "user.login",
		})

		if err != nil {
			// Call out to sentry
			fmt.Println("Had a problem with request")
		}

		body, _ := json.Marshal(resp)
		fmt.Fprintf(w, string(body))
	})
	log.Fatal(http.ListenAndServe(":8081", nil))
}
