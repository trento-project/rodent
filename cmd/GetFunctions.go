package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func CheckHealthFunc() {
	ep := fmt.Sprintf("http://%s:8080/api/health", runner)
	log.Debugf("Attempting to get %s", ep)
	res, err := http.Get(ep)
	if err != nil {
		log.Fatal(err.Error())
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	var j1 HealthResult
	err = json.Unmarshal(resData, &j1)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Debug("Successfully unmarshalled response")

	if j1.Status == "ok" {
		log.Infof("The status of endpoint '%s' is '%s'\n", ep, j1.Status)
	} else {
		log.Warnf("The status of the endpoint '%s' is not '%s\n'", j1.Status)
	}
}

func CheckReadyFunc() {
	ep := fmt.Sprintf("http://%s:8080/api/ready", runner)
	log.Debugf("Attempting to get %s", ep)
	res, err := http.Get(ep)
	if err != nil {
		log.Fatal(err.Error())
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	var j1 ReadyResult
	err = json.Unmarshal(resData, &j1)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Debug("Successfully unmarshalled response")

	if j1.Ready {
		log.Infof("The status of endpoint '%s' is 'true'\n", ep)
	} else {
		log.Warnf("The status of the endpoint '%s' is 'false'\n", ep)
	}
}

func CheckCatalogFunc() {
	if debug {
		log.SetLevel(log.DebugLevel)
	}
	j1 := GetCatalog()
	log.Infof("Found %d checks\n", len(j1))
	log.Info("ID\t\tName")
	for _, v := range j1 {
		log.Infof("%s\t%s", v.ID, v.Name)
	}
}

func GetCatalog() Catalog {
	ep := fmt.Sprintf("http://%s:8080/api/catalog", runner)
	log.Debugf("Attempting to get %s", ep)
	res, err := http.Get(ep)
	if err != nil {
		log.Fatal(err.Error())
	}

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err.Error())
	}

	j1 := Catalog{}
	err = json.Unmarshal(resData, &j1)
	if err != nil {
		log.Fatal(err.Error())
	}
	return j1
}
