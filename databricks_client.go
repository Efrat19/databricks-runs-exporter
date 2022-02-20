package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

var databricksHost string
var databricksToken string
var runsScrapeLimit int
var runsScrapeTimespanSeconds int

func init() {
	var err error
	databricksHost = getEnv("DATABRICKS_HOST", "")
	databricksToken = getEnv("DATABRICKS_TOKEN", "")
	runsScrapeLimit, err = strconv.Atoi(getEnv("RUNS_SCRAPE_LIMIT", "50"))
	runsScrapeTimespanSeconds, err = strconv.Atoi(getEnv("RUNS_SCRAPE_TIMESPAN_SECONDS", "10"))
	if err != nil {
		panic(err)
	}
	if databricksHost == "" || databricksToken == "" {
		panic(errors.New("Failed to Initialize exporter: env vars DATABRICKS_HOST and DATABRICKS_TOKEN must be specified"))
	}
	if !strings.HasPrefix(databricksHost, "https://") {
		databricksHost = fmt.Sprintf("https://%s", databricksHost)
	}

}

func getScrapeWindowEdges() (int64, int64) {
	to := time.Now()
	safetyMargin := 10
	span := time.Duration((-runsScrapeTimespanSeconds - safetyMargin) * int(time.Second))
	from := to.Add(span)
	return to.Unix(), from.Unix()
}

func (r *Run) ToLabelValues(labels []string) ([]string, error) {
	jsonRun, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}
	labelVals := []string{}
	for _, label := range labels {
		value := gjson.Get(string(jsonRun), label)
		labelVals = append(labelVals, value.String())
	}
	return labelVals, nil
}

func formatRuns(rruns *[]rawRun) *[]Run {
	formated := make([]Run, len(*rruns))
	for i, rr := range *rruns {
		formated[i] = Run{
			JobID:                   rr.JobID,
			RunID:                   rr.RunID,
			NumberInJob:             rr.NumberInJob,
			OriginalAttemptRunID:    rr.OriginalAttemptRunID,
			LifeCycleState:          rr.State.LifeCycleState,
			ResultState:             rr.State.ResultState,
			StateMessage:            rr.State.StateMessage,
			UserCancelledOrTimedout: rr.State.UserCancelledOrTimedout,
			StartTime:               rr.StartTime,
			SetupDuration:           rr.SetupDuration,
			ExecutionDuration:       rr.ExecutionDuration,
			CleanupDuration:         rr.CleanupDuration,
			EndTime:                 rr.EndTime,
			Trigger:                 rr.Trigger,
			CreatorUserName:         rr.CreatorUserName,
			RunName:                 rr.RunName,
			RunType:                 rr.RunType,
			AttemptNumber:           rr.AttemptNumber,
			Format:                  rr.Format,
		}
	}
	return &formated
}

func dedupByRunID(runs *[]Run) *[]Run {
	unique := []Run{}
	set := map[int]Run{}
	for _, run := range *runs {
		set[run.RunID] = run
	}
	for _, run := range set {
		unique = append(unique, run)
	}
	return &unique
}

func sendApiRequest(url string, method string) ([]byte, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", databricksToken))
	resp, err := client.Do(req)
	//Handle Error
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetRuns() (*[]Run, error) {
	startTimeFrom, startTimeTo := getScrapeWindowEdges()
	url := fmt.Sprintf("%s/api/2.0/jobs/runs/list?limit=%d&start_time_from=%d&start_time_to=%d", databricksHost, runsScrapeLimit, startTimeFrom, startTimeTo)
	body, err := sendApiRequest(url, "GET")
	if err != nil {
		return nil, err
	}
	var apiResponse *ApiResponse
	json.Unmarshal(body, &apiResponse)
	log.Infof("Collected %d runs started between %d and %d\n", runsScrapeLimit, startTimeFrom, startTimeTo)
	return dedupByRunID(formatRuns(&apiResponse.Runs)), nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
