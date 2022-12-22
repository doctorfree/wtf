package nextbus

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/rivo/tview"
	"github.com/doctorfree/wtf/logger"
	"github.com/doctorfree/wtf/view"
)

// Widget is the container for your module's data
type Widget struct {
	view.TextWidget

	settings *Settings
}

// NewWidget creates and returns an instance of Widget
func NewWidget(tviewApp *tview.Application, redrawChan chan bool, pages *tview.Pages, settings *Settings) *Widget {
	widget := Widget{
		TextWidget: view.NewTextWidget(tviewApp, redrawChan, pages, settings.common),

		settings: settings,
	}
	return &widget
}

/* -------------------- Exported Functions -------------------- */

// Refresh updates the onscreen contents of the widget
func (widget *Widget) Refresh() {
	// The last call should always be to the display function
	widget.display()
}

/* -------------------- Unexported Functions -------------------- */

func (widget *Widget) content() string {
	return getNextBus(widget.settings.agency, widget.settings.route, widget.settings.stopID)
}

type AutoGenerated struct {
	Copyright   string      `json:"copyright"`
	Predictions Predictions `json:"predictions"`
}

type Prediction struct {
	AffectedByLayover string `json:"affectedByLayover"`
	Seconds           string `json:"seconds"`
	TripTag           string `json:"tripTag"`
	Minutes           string `json:"minutes"`
	IsDeparture       string `json:"isDeparture"`
	Block             string `json:"block"`
	DirTag            string `json:"dirTag"`
	Branch            string `json:"branch"`
	EpochTime         string `json:"epochTime"`
	Vehicle           string `json:"vehicle"`
}

type Direction struct {
	PredictionRaw json.RawMessage `json:"prediction"`
	Title         string          `json:"title"`
}

type Predictions struct {
	RouteTag    string    `json:"routeTag"`
	StopTag     string    `json:"stopTag"`
	RouteTitle  string    `json:"routeTitle"`
	AgencyTitle string    `json:"agencyTitle"`
	StopTitle   string    `json:"stopTitle"`
	Direction   Direction `json:"direction"`
}

func getNextBus(agency string, route string, stopID string) string {
	url := fmt.Sprintf("https://webservices.umoiq.com/service/publicJSONFeed?command=predictions&a=%s&r=%s&stopId=%s", agency, route, stopID)
	resp, err := http.Get(url)
	if err != nil {
		logger.Log(fmt.Sprintf("[nextbus] Error: Failed to make requests to umoiq for next bus predictions. Reason: %s", err))
		return "[nextbus] error calling umoiq"
	}
	body, readErr := io.ReadAll(resp.Body)

	if (readErr) != nil {
		logger.Log(fmt.Sprintf("[nextbus] Error: Failed to parse response body from umoiq. Reason: %s", err))
		return "[nextbus] error parsing response body"
	}

	resp.Body.Close()

	var parsedResponse AutoGenerated

	// partial unmarshal, we don't have r.Predictions.Direction.PredictionRaw <- YET
	unmarshalError := json.Unmarshal(body, &parsedResponse)
	if unmarshalError != nil {
		logger.Log(fmt.Sprintf("[nextbus] Error: Failed to unmarshal body from umoiq. Reason: %s", err))
		return "[nextbus] error unmarshalling response body"
	}

	parseType := ""
	// hacky, try object parse first
	nextBusObject := Prediction{}
	if err := json.Unmarshal(parsedResponse.Predictions.Direction.PredictionRaw, &nextBusObject); err == nil {
		parseType = "object"
	}

	// if object parse failed, it probably means we have an array
	nextBuses := []Prediction{}
	if err := json.Unmarshal(parsedResponse.Predictions.Direction.PredictionRaw, &nextBuses); err == nil {
		parseType = "array"
	}

	// build the final string
	finalStr := ""
	if parseType == "array" {
		for _, nextBus := range nextBuses {
			finalStr += fmt.Sprintf("%s | ETA [%s]\n", parsedResponse.Predictions.RouteTitle, strTimeToInt(nextBus.Minutes, nextBus.Seconds))
		}
	} else {
		finalStr += fmt.Sprintf("%s | ETA [%s]\n", parsedResponse.Predictions.RouteTitle, strTimeToInt(nextBusObject.Minutes, nextBusObject.Seconds))
	}

	return finalStr
}

// takes minutes and seconds from the API, does math to find the remainder seconds
// since the API only gives whole minutes
func strTimeToInt(sourceMinutes string, sourceSeconds string) string {
	min, _ := strconv.Atoi(sourceMinutes)
	sec, _ := strconv.Atoi(sourceSeconds)
	sec = sec % 60
	return fmt.Sprintf("%02d:%02d", min, sec)
}

func (widget *Widget) display() {
	widget.Redraw(func() (string, string, bool) {
		return widget.CommonSettings().Title, widget.content(), false
	})
}
