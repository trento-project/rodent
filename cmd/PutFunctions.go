/*
Copyright © 2022 mr-stringer
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

var ec ExecuteCatalog

type CallbackListener struct {
	ch chan *Callback
}

func (cbl *CallbackListener) CallbackFunc(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		Response(w, "Content Type is not application/json", http.StatusUnsupportedMediaType)
		return
	}
	log.Debug("Callback received response")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Warn("Error reading response")
		log.Fatal(err.Error())
	}
	r.Body.Close()

	var c1 *Callback
	err = json.Unmarshal(body, &c1)
	if err != nil {
		log.Fatal(err.Error())
	}
	Response(w, "ok", http.StatusAccepted)
	cbl.ch <- c1
}

func (cbl *CallbackListener) HandleResult(CheckID string, wg *sync.WaitGroup) {
	for {
		select {
		case clbk := <-cbl.ch:
			switch {
			case clbk.Event == "execution_started":
				log.Debug("Execution started")
			case clbk.Event == "execution_completed":
				/*We only expect one host, make sure that's all we have*/
				if len(clbk.Payload.Hosts) != 1 {
					log.Fatal("More than 1 host in the returned payload. This is unexpected")
				}
				/*Ensure the host was reachable*/
				if !clbk.Payload.Hosts[0].Reachable {
					/*It's probably OK to call this fatal*/
					log.Fatalf("The runner '%s' could not reach the target host '%s'", runner, hostToCheck)

				}

				ec.Mutex.Lock()
				/*Get the cid out of the map based on execution ID*/
				cid, ok := ec.ExecuteMap[clbk.ExecutionID]
				ec.Mutex.Unlock()
				if !ok {
					log.Fatalf("The executionID %v was not found in the ExecutionMap, cannot continue", clbk.ExecutionID)
				}

				log.Debugf("Found CheckID %s for executionID %v", cid, clbk.ExecutionID)

				var reported = false
				/*Find the check*/
				for k := range clbk.Payload.Hosts[0].Results {
					if clbk.Payload.Hosts[0].Results[k].CheckID == cid {
						/*Report the result*/
						log.Infof("CheckID %s state is '%s'", clbk.Payload.Hosts[0].Results[k].CheckID, clbk.Payload.Hosts[0].Results[k].Result)
						if clbk.Payload.Hosts[0].Results[k].Result != "passing" && clbk.Payload.Hosts[0].Results[k].Message != "" {
							log.Infof("Message: '%s'", clbk.Payload.Hosts[0].Results[k].Message)
						}
						reported = true
					}
				}
				wg.Done()
				if !reported {
					log.Warnf("The CheckID %s was not found in the response to executionID %v", cid, clbk.ExecutionID)
				}
			}
		case <-time.After(time.Second * 60):
			log.Error("The callback was not received within the timeout period")
			wg.Done()
			return
		}
	}
}

//Execute a single check and get the result
func ExecuteCheckFunc() {
	/*Check that we have a valid check ID*/
	if !CheckIsValid() {
		log.Fatalf("The check ID %s was not found in the runner.  For a list of valid IDs use the CheckCatalog subcommand", checkID)
	} else {
		log.Debugf("Check ID %s is valid", checkID)
	}

	/*Create a wait group that only waits for the CheckResults function to finish*/
	wg := new(sync.WaitGroup)
	wg.Add(1)

	cbl := &CallbackListener{ch: make(chan *Callback)}
	go cbl.HandleResult(checkID, wg)

	/*Start the WebServer*/
	go WebServer(cbl)

	/*Wait for WebServer initialisation*/
	time.Sleep(time.Second * 5)

	//Send the execute request
	SingleCheck(checkID)

	wg.Wait()
	log.Debug("Finishing up")
}

/*Execute all checks*/
func ExecuteAllChecksFunc() {
	/*First off, get all the checks*/
	Catalog := GetCatalog()

	/*There is no need to check if the tests are valid, they just came out of the API*/

	/*Create a WaitGroup that will wait for all of the checks to finish*/
	wg := new(sync.WaitGroup)
	wg.Add(len(Catalog))

	cbl := &CallbackListener{ch: make(chan *Callback)}
	go cbl.HandleResult(checkID, wg)

	/*Start the WebServer*/
	go WebServer(cbl)

	/*Wait for WebServer initialisation*/
	time.Sleep(time.Second * 5)

	/*Loop through all checks*/
	for k := range Catalog {
		SingleCheck(Catalog[k].ID)
		time.Sleep(time.Second * time.Duration(checkInterval))
	}

	log.Debug("All checks submitted")

	wg.Wait()
	log.Debug("Finishing up")

}

/*There will be at least two calls.  One that states the call was accepted, a second with the result*/

func WebServer(cbl *CallbackListener) {
	log.Debug("Starting WebServer")
	callback := http.HandlerFunc(cbl.CallbackFunc)
	http.Handle(callbackUrl, callback)
	port := fmt.Sprintf(":%d", callbackPort)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Warn("WebServer Crashed!")
		log.Fatal(err.Error())
	}
}

func Response(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func SingleCheck(CheckID string) {
	/*Define the endpoint*/
	ep := fmt.Sprintf("http://%s:8080/api/execute", runner)

	/*Create the execution event*/
	Event := &ExecutionEvent{
		ExecutionID: uuid.New(), //New UUID for this event
		ClusterID:   uuid.New(), //This should be consistent for each call, but right now thats unimportant
		Provider:    provider,
		Checks:      []string{CheckID},
		Hosts: []*Host{
			{
				HostID:  uuid.New(), //should probably match what we have in the GUI, but we'll see what happens
				Address: hostToCheck,
				User:    "root", //root seems a pretty safe assumption for now
			},
		},
	}

	ec.Lock()
	if ec.ExecuteMap == nil {
		ec.ExecuteMap = make(map[uuid.UUID]string)
	}
	ec.ExecuteMap[Event.ExecutionID] = CheckID
	ec.Unlock()

	/*Marshal the event to json*/
	body, err := json.Marshal(Event)
	if err != nil {
		log.Error("Could not encode json")
		log.Fatal(err.Error())
	}

	log.Debugf("Requesting execution check '%s' on target '%s' via %s", CheckID, hostToCheck, ep)

	res, err := http.Post(ep, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err.Error())
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Debug(string(resData))
}

func CheckIsValid() bool {
	cat := GetCatalog()
	/*search for check in the catalog*/
	var found bool /*false by default*/
	for _, v := range cat {
		if v.ID == checkID {
			log.Debugf("Check %s was found in the catalog", checkID)
			found = true
			break
		}
	}
	return found
}
