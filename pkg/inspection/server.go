// Copyright 2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package inspection

import (
	"fmt"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/common/idgenerator"
	inspectioncontract "github.com/GoogleCloudPlatform/khi/pkg/inspection/contract"
	"github.com/GoogleCloudPlatform/khi/pkg/inspection/inspectiondata"
	"github.com/GoogleCloudPlatform/khi/pkg/task"
	"golang.org/x/exp/slices"
)

type PrepareInspectionServerFunc = func(inspectionServer *InspectionTaskServer) error

type InspectionType struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Priority    int    `json:"-"`

	// Document properties
	DocumentDescription string `json:"-"`
}

type FeatureListItem struct {
	Id          string `json:"id"`
	Label       string `json:"label"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

type InspectionDryRunResult struct {
	Metadata interface{} `json:"metadata"`
}

type InspectionRunResult struct {
	Metadata    interface{}
	ResultStore inspectiondata.Store
}

// InspectionTaskServer manages tasks and provides apis to get task related information in JSON convertible type.
type InspectionTaskServer struct {
	// RootTaskSet is the set of the all tasks in KHI.
	RootTaskSet *task.TaskSet
	// inspectionTypes are kinds of tasks. Users will select this at first to filter togglable feature tasks.
	inspectionTypes []*InspectionType
	// inspections are generated inspection task runers
	inspections           map[string]*InspectionTaskRunner
	inspectionIDGenerator idgenerator.IDGenerator

	ioConfig *inspectioncontract.IOConfig
}

func NewServer(ioConfig *inspectioncontract.IOConfig) (*InspectionTaskServer, error) {
	ns, err := task.NewTaskSet([]task.UntypedTask{})
	if err != nil {
		return nil, err
	}
	return &InspectionTaskServer{
		RootTaskSet:           ns,
		inspectionTypes:       make([]*InspectionType, 0),
		inspections:           map[string]*InspectionTaskRunner{},
		inspectionIDGenerator: idgenerator.NewPrefixIDGenerator("inspection-"),
		ioConfig:              ioConfig,
	}, nil
}

// AddInspectionType register a inspection type.
func (s *InspectionTaskServer) AddInspectionType(newInspectionType InspectionType) error {
	if strings.Contains(newInspectionType.Id, "/") {
		return fmt.Errorf("inspection type must not contain /")
	}
	idMap := map[string]interface{}{}
	for _, inspectionType := range s.inspectionTypes {
		idMap[inspectionType.Id] = struct{}{}
	}
	if _, exist := idMap[newInspectionType.Id]; exist {
		return fmt.Errorf("inspection type id:%s is duplicated. InspectionType ID must be unique", newInspectionType.Id)
	}
	s.inspectionTypes = append(s.inspectionTypes, &newInspectionType)
	slices.SortFunc(s.inspectionTypes, func(a *InspectionType, b *InspectionType) int {
		return b.Priority - a.Priority
	})
	return nil
}

// AddTask register a task usable for the inspection task graph execution.
func (s *InspectionTaskServer) AddTask(task task.UntypedTask) error {
	return s.RootTaskSet.Add(task)
}

// CreateInspection generates an inspection and returns inspection ID
func (s *InspectionTaskServer) CreateInspection(inspectionType string) (string, error) {
	id := s.inspectionIDGenerator.Generate()
	inspectionTask := NewInspectionRunner(s, s.ioConfig, id)
	err := inspectionTask.SetInspectionType(inspectionType)
	if err != nil {
		return "", err
	}
	s.inspections[inspectionTask.ID] = inspectionTask
	return inspectionTask.ID, nil
}

// Inspection returns an instance of an Inspection queried with given inspection ID.
func (s *InspectionTaskServer) GetInspection(inspectionID string) *InspectionTaskRunner {
	return s.inspections[inspectionID]
}

func (s *InspectionTaskServer) GetAllInspectionTypes() []*InspectionType {
	return append([]*InspectionType{}, s.inspectionTypes...)
}

func (s *InspectionTaskServer) GetInspectionType(inspectionTypeId string) *InspectionType {
	for _, registeredType := range s.inspectionTypes {
		if registeredType.Id == inspectionTypeId {
			return registeredType
		}
	}
	return nil
}

func (s *InspectionTaskServer) GetAllRunners() []*InspectionTaskRunner {
	inspections := []*InspectionTaskRunner{}
	for _, value := range s.inspections {
		inspections = append(inspections, value)
	}
	return inspections
}

// GetAllRegisteredTasks returns a cloned list of all tasks registered in this server.
func (s *InspectionTaskServer) GetAllRegisteredTasks() []task.UntypedTask {
	return s.RootTaskSet.GetAll()
}
