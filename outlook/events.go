package outlook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hiteshrepo/ms-graph-api/token"
)

type Event struct {
	ID      string    `json:"id"`
	Subject string    `json:"subject"`
	Start   EventTime `json:"start"`
	End     EventTime `json:"end"`
}

type EventTime struct {
	Datetime string `json:"dateTime"`
	Timezone string `json:"timeZone"`
}

func CreateEvent() {
	accessToken, err := token.GetTokenByPublicClientApp()
	if err != nil {
		panic(err)
	}
	eventURL := "https://graph.microsoft.com/v1.0/me/events"

	eventData := []byte(`
    {
        "subject": "Sample Event",
        "start": {
            "dateTime": "2023-09-25T12:00:00",
            "timeZone": "UTC"
        },
        "end": {
            "dateTime": "2023-09-25T13:00:00",
            "timeZone": "UTC"
        }
    }`)

	req, err := http.NewRequest("POST", eventURL, bytes.NewBuffer(eventData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		fmt.Println("Error creating calendar event:", resp.Status)
		return
	}

	fmt.Println("Calendar event created successfully")
}

func FetchAllEvents() ([]Event, error) {
	accessToken, err := token.GetTokenByPublicClientApp()
	if err != nil {
		panic(err)
	}

	eventsURL := "https://graph.microsoft.com/v1.0/me/events"

	req, err := http.NewRequest("GET", eventsURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error retrieving calendar events:", resp.Status)
		return nil, err
	}

	var response struct {
		Events []Event `json:"value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Println("Error decoding event response:", err)
		return nil, err
	}

	for _, event := range response.Events {
		fmt.Println("Subject:", event.Subject)
	}

	return response.Events, nil
}

func DeleteEvent(deleteEventId string) {

	allEvents, err := FetchAllEvents()
	if err != nil {
		panic(err)
	}

	if len(allEvents) == 0 {
		panic("no events")
	}

	anEvent := allEvents[0]

	accessToken, err := token.GetTokenByPublicClientApp()
	if err != nil {
		panic(err)
	}

	deleteURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/me/events/%s", anEvent.ID)

	req, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("Error deleting calendar event:", resp.Status)
		return
	}

	fmt.Println("Calendar event deleted successfully")
}
