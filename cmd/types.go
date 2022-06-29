/*
Copyright © 2022 mr-stringer
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"
)

type ReadyResult *struct {
	Ready bool `json:"ready"`
}

type HealthResult *struct {
	Status string `json:"status"`
}

type Catalog []*CatalogCheck

type CatalogCheck struct {
	ID             string `json:"id,omitempty" binding:"required"`
	Name           string `json:"name,omitempty" binding:"required"`
	Group          string `json:"group" binding:"required"`
	Provider       string `json:"provider" binding:"required"`
	Description    string `json:"description,omitempty"`
	Remediation    string `json:"remediation,omitempty"`
	Implementation string `json:"implementation,omitempty"`
	Labels         string `json:"labels,omitempty"`
	Premium        bool   `json:"premium,omitempty"`
}

/*The following two types are used to execute checks*/

type ExecutionEvent struct {
	ExecutionID uuid.UUID `json:"execution_id" binding:"required"`
	ClusterID   uuid.UUID `json:"cluster_id" binding:"required"`
	Provider    string    `json:"provider" binding:"required"`
	Checks      []string  `json:"checks" binding:"required"`
	Hosts       []*Host   `json:"hosts" binding:"required"`
}

type Host struct {
	HostID  uuid.UUID `json:"host_id" binding:"required"`
	Address string    `json:"address" binding:"required"`
	User    string    `json:"user" binding:"required"`
}

/*The below types handle the callback from trento-web*/

type CheckResult struct {
	CheckID string `json:"check_id"`
	Result  string `json:"result"`
	Message string `json:"msg"`
}

type CheckHost struct {
	HostID    uuid.UUID      `json:"host_id"`
	Reachable bool           `json:"reachable"`
	Results   []*CheckResult `json:"results"`
}

type Payload struct {
	ClusterID uuid.UUID    `json:"cluster_id"`
	Hosts     []*CheckHost `json:"hosts"`
}

type Callback struct {
	ExecutionID uuid.UUID `json:"execution_id"`
	Event       string    `json:"event"`
	Payload     *Payload  `json:"payload"`
}

func (c Callback) Print() {
	cbj, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}
	fmt.Print(string(cbj))
}

type ExecuteCatalog struct {
	ExecuteMap map[uuid.UUID]string
	sync.Mutex
}
