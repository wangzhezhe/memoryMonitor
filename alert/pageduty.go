package alert

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var (
	PageDutyTriggerUrl  = "https://events.pagerduty.com/generic/2010-04-15/create_event.json"
	PageDutyMemoryAlert = "cformMemoryAlert"
	DefaultKey          = "697b22fda6e34904b197fdd467602e41"
)

/*
   data_json={
     "service_key": key,
     "incident_key": event_name,
     "event_type": "trigger",
     "description": event_name,
     "client": APP_PAGEDUTY_CLIENT_NAME,
     #"client_url": os.getenv('APP_PAGEDUTY_CLIENT_NAME', 'defalut'),
     "details": {
       "data": data
     }
   }

*/
type AlertInfo struct {
	ServiceKey  string            `json:"service_key"`
	IncidentKey string            `json:"incident_key"`
	EventType   string            `json:"event_type"`
	Description string            `json:"description"`
	Client      string            `json:"client"`
	Details     map[string]string `json:"details"`
}

func Trigger(nodeIP string, containerID string, memoryPercentage int) error {
	message := "Cform production alert!!! " + "nodeIP:" + nodeIP + " ContainerID: " + containerID + " MemoryPercentage: " + strconv.Itoa(memoryPercentage) + "%"
	alertInfo := &AlertInfo{
		ServiceKey:  DefaultKey,
		IncidentKey: PageDutyMemoryAlert + "/" + nodeIP,
		EventType:   "trigger",
		Description: "cform memory alert: " + message,
	}
	info, err := json.Marshal(alertInfo)
	if err != nil {
		return nil
	}
	client := &http.Client{}
	reqest, err := http.NewRequest("POST", PageDutyTriggerUrl, strings.NewReader(string(info)))
	if err != nil {
		return err
	}
	response, err := client.Do(reqest)
	if err != nil {
		return err
	}
	returnbody, _ := ioutil.ReadAll(response.Body)
	log.Println(string(returnbody))

	return nil

}

func Resolve(nodeIP string, containerID string, memoryPercentage int) error {
	alertInfo := &AlertInfo{
		ServiceKey:  DefaultKey,
		IncidentKey: PageDutyMemoryAlert + "/" + nodeIP,
		EventType:   "resolve",
		Description: "cform memory alert resolved",
	}
	info, err := json.Marshal(alertInfo)
	if err != nil {
		return nil
	}
	client := &http.Client{}
	reqest, err := http.NewRequest("POST", PageDutyTriggerUrl, strings.NewReader(string(info)))
	if err != nil {
		return err
	}
	response, err := client.Do(reqest)
	if err != nil {
		return err
	}
	returnbody, _ := ioutil.ReadAll(response.Body)
	log.Println(string(returnbody))

	return nil

}
