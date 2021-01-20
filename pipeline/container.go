// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package pipeline

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/drone/envsubst"
	"github.com/go-vela/types/constants"
)

type (
	// ContainerSlice is the pipeline representation
	// of the Containers block for a pipeline.
	//
	// swagger:model PipelineContainerSlice
	//
	// swagger:model PipelineContainerSlice
	ContainerSlice []*Container

	// Container is the pipeline representation
	// of a Container in a pipeline.
	//
	// swagger:model PipelineContainer
	// nolint:maligned // suppressing struct optimization, prefer to keep current order
	Container struct {
		ID          string            `json:"id,omitempty"          yaml:"id,omitempty"`
		Commands    []string          `json:"commands,omitempty"    yaml:"commands,omitempty"`
		Detach      bool              `json:"detach,omitempty"      yaml:"detach,omitempty"`
		Directory   string            `json:"directory,omitempty"   yaml:"directory,omitempty"`
		Entrypoint  []string          `json:"entrypoint,omitempty"  yaml:"entrypoint,omitempty"`
		Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`
		ExitCode    int               `json:"exit_code,omitempty"   yaml:"exit_code,omitempty"`
		Image       string            `json:"image,omitempty"       yaml:"image,omitempty"`
		Name        string            `json:"name,omitempty"        yaml:"name,omitempty"`
		Needs       []string          `json:"needs,omitempty"       yaml:"needs,omitempty"`
		Networks    []string          `json:"networks,omitempty"    yaml:"networks,omitempty"`
		Number      int               `json:"number,omitempty"      yaml:"number,omitempty"`
		Ports       []string          `json:"ports,omitempty"       yaml:"ports,omitempty"`
		Privileged  bool              `json:"privileged,omitempty"  yaml:"privileged,omitempty"`
		Pull        string            `json:"pull,omitempty"        yaml:"pull,omitempty"`
		Ruleset     Ruleset           `json:"ruleset,omitempty"     yaml:"ruleset,omitempty"`
		Secrets     StepSecretSlice   `json:"secrets,omitempty"     yaml:"secrets,omitempty"`
		Ulimits     UlimitSlice       `json:"ulimits,omitempty"     yaml:"ulimits,omitempty"`
		Volumes     VolumeSlice       `json:"volumes,omitempty"     yaml:"volumes,omitempty"`
	}
)

// Purge removes the Containers that have a ruleset
// that do not match the provided ruledata.
func (c *ContainerSlice) Purge(r *RuleData) *ContainerSlice {
	counter := 1
	containers := new(ContainerSlice)

	// iterate through each Container in the pipeline
	for _, container := range *c {
		// verify ruleset matches
		if container.Ruleset.Match(r) {
			// overwrite the Container number with the Container counter
			container.Number = counter

			// increment Container counter
			counter = counter + 1

			// append the Container to the new slice of Containers
			*containers = append(*containers, container)
		}
	}

	// return the new slice of Containers
	return containers
}

// Sanitize cleans the fields for every step in the pipeline so they
// can be safely executed on the worker. The fields are sanitized
// based off of the provided runtime driver which is setup on every
// worker. Currently, this function supports the following runtimes:
//
//   * Docker
//   * Kubernetes
func (c *ContainerSlice) Sanitize(driver string) *ContainerSlice {
	containers := new(ContainerSlice)

	// iterate through each Container in the pipeline
	for _, container := range *c {
		// sanitize container
		cont := container.Sanitize(driver)

		// append the Container to the new slice of Containers
		*containers = append(*containers, cont)
	}

	return containers
}

// Empty returns true if the provided container is empty.
func (c *Container) Empty() bool {
	// return true if the container is nil
	if c == nil {
		return true
	}

	// return true if every container field is empty
	if len(c.ID) == 0 &&
		len(c.Commands) == 0 &&
		!c.Detach &&
		len(c.Directory) == 0 &&
		len(c.Entrypoint) == 0 &&
		len(c.Environment) == 0 &&
		c.ExitCode == 0 &&
		len(c.Image) == 0 &&
		len(c.Name) == 0 &&
		len(c.Needs) == 0 &&
		len(c.Networks) == 0 &&
		c.Number == 0 &&
		len(c.Ports) == 0 &&
		!c.Privileged &&
		len(c.Pull) == 0 &&
		reflect.DeepEqual(c.Ruleset, Ruleset{}) &&
		len(c.Secrets) == 0 &&
		len(c.Ulimits) == 0 &&
		len(c.Volumes) == 0 {
		return true
	}

	// return false if any of the ruletype is provided
	return false
}

// Execute returns true when the provided ruledata matches
// the conditions when we should be running the container on the worker.
func (c *Container) Execute(r *RuleData) bool {
	// return false if the container is nil
	if c == nil {
		return false
	}

	// check if the build is in a running state
	if strings.EqualFold(r.Status, constants.StatusRunning) {
		// treat the ruleset status as success
		r.Status = constants.StatusSuccess

		// skip evaluating path in ruleset
		//
		// the compiler is the component responsible for
		// choosing whether a container will run based
		// off the files changed for a build
		//
		// the worker doesn't have any record of
		// what files changed for a build so we
		// should "skip" evaluating what the
		// user provided for the path element
		c.Ruleset.If.Path = []string{}
		c.Ruleset.Unless.Path = []string{}

		// return if the container ruleset matches the conditions
		return c.Ruleset.Match(r)
	}

	// assume you will execute the container
	execute := true

	// capture the build status out of the ruleset
	status := r.Status

	// check if the build status is successful
	if !strings.EqualFold(status, constants.StatusSuccess) {
		// disregard the need to run the container
		execute = false

		// check if you need to run a status failure ruleset
		if !(c.Ruleset.If.Empty() && c.Ruleset.Unless.Empty()) &&
			!(c.Ruleset.If.NoStatus() && c.Ruleset.Unless.NoStatus()) &&
			c.Ruleset.Match(r) {
			// approve the need to run the container
			execute = true
		}
	}

	r.Status = constants.StatusFailure

	// check if you need to skip a status failure ruleset
	if strings.EqualFold(status, constants.StatusSuccess) &&
		!(c.Ruleset.If.NoStatus() && c.Ruleset.Unless.NoStatus()) &&
		!(c.Ruleset.If.Empty() && c.Ruleset.Unless.Empty()) && c.Ruleset.Match(r) {
		r.Status = constants.StatusSuccess

		if !c.Ruleset.Match(r) {
			// disregard the need to run the container
			execute = false
		}
	}

	return execute
}

// Sanitize cleans the fields for every step in the pipeline so they
// can be safely executed on the worker. The fields are sanitized
// based off of the provided runtime driver which is setup on every
// worker. Currently, this function supports the following runtimes:
//
//   * Docker
//   * Kubernetes
func (c *Container) Sanitize(driver string) *Container {
	container := c

	switch driver {
	// sanitize container for Docker
	case constants.DriverDocker:
		if strings.Contains(c.ID, " ") {
			c.ID = strings.ReplaceAll(c.ID, " ", "-")
		}

		return container
	// sanitize container for Kubernetes
	case constants.DriverKubernetes:
		if strings.Contains(c.ID, " ") {
			container.ID = strings.ReplaceAll(c.ID, " ", "-")
		}

		if strings.Contains(c.ID, "_") {
			container.ID = strings.ReplaceAll(c.ID, "_", "-")
		}

		if strings.Contains(c.ID, ".") {
			container.ID = strings.ReplaceAll(c.ID, ".", "-")
		}

		return container
	// unrecognized driver
	default:
		// TODO: add a log message indicating how we got here
		return nil
	}
}

// Substitute replaces every reference (${VAR} or $${VAR}) to an
// environment variable in the container configuration with the
// corresponding value for that environment variable.
func (c *Container) Substitute() error {
	// check if container or container environment are nil
	if c == nil || c.Environment == nil {
		return errors.New("empty container environment provided")
	}

	// marshal container configuration
	body, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// create substitute function
	subFunc := func(name string) string {
		// capture the environment variable value
		value := c.Environment[name]

		// check for a new line in the value
		if strings.Contains(value, "\n") {
			// safely escape the environment variable
			value = fmt.Sprintf("%q", value)
		}

		return value
	}

	// substitute the environment variables
	//
	// https://pkg.go.dev/github.com/drone/envsubst?tab=doc#Eval
	ctn, err := envsubst.Eval(string(body), subFunc)
	if err != nil {
		return err
	}

	// unmarshal container configuration
	err = json.Unmarshal([]byte(ctn), c)
	if err != nil {
		// create a new buffer for encoded JSON
		//
		// will be thrown away after encoding
		b := new(bytes.Buffer)

		// create new JSON encoder attached to buffer
		enc := json.NewEncoder(b)

		// JSON encode container output
		//
		// buffer is thrown away
		err = enc.Encode(c)
		if err != nil {
			return err
		}
	}

	return nil
}
