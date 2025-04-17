// Copyright 2025 Google LLC
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

package model

import (
	"slices"
	"strings"

	"github.com/GoogleCloudPlatform/khi/pkg/common/filter"
	"github.com/GoogleCloudPlatform/khi/pkg/common/typedmap"
	"github.com/GoogleCloudPlatform/khi/pkg/inspection"
	inspection_task "github.com/GoogleCloudPlatform/khi/pkg/inspection/task"
	"github.com/GoogleCloudPlatform/khi/pkg/inspection/task/label"
	"github.com/GoogleCloudPlatform/khi/pkg/model/enum"
	"github.com/GoogleCloudPlatform/khi/pkg/task"
)

// FeatureDocumentModel is a model type for generating document docs/en/reference/features.md
type FeatureDocumentModel struct {
	// Features are the list of feature tasks defined in KHI.
	Features []FeatureDocumentElement
}

// FeatureDocumentElement is a model type for a feature task used in FeatureDocumentModel.
type FeatureDocumentElement struct {
	// ID is the unique name of the feature task.
	ID string
	// Name is the human readable name of the feature task.
	Name string
	// Description is the string explain the feature task.
	Description string
	// Forms is the list of information about form inputs that is required from the feature task.
	Forms []FeatureDependentFormElement
	// IndirectQueryDependency is the list of query tasks that is required from this feature task but not the target query task.
	IndirectQueryDependency []FeatureIndirectDependentQueryElement
	// TargetQUeryDependency is the main query task used in this feature task.
	TargetQueryDependency FeatureDependentTargetQueryElement
	// OutputTimelines is the list of timelines(=ParentRelationship type) that can be generated by this feature task.
	OutputTimelines []FeatureOutputTimelineElement
	// AvailableInspctionType is the list of InspectionType that supports this feature task.
	AvailableInspectionTypes []FeatureAvailableInspectionType
}

// FeatureDependentFormElement is a model type for a input form required from a feature task.
type FeatureDependentFormElement struct {
	// ID is the unique name of this form element.
	ID string
	// Label is a human readable short name of this input element.
	Label string
	// Description is a string explaining this form input.
	Description string
}

// FeatureIndirectDependentQueryElement is a model type for query tasks required from a feature task but not the target query task.
type FeatureIndirectDependentQueryElement struct {
	// ID is the unique name of this query task.
	ID string
	// LogTypeLabel is a human readable short name of the log type queried by the query task.
	LogTypeLabel string
	// LogTypeColorCode is the hex color code without the `#` prefix for the log type.
	LogTypeColorCode string
}

// FeatureDependentTargetQueryElement is a model type for a target query task of the feature task.
type FeatureDependentTargetQueryElement struct {
	// ID is the unique name of this query task.
	ID string
	// LogTypeLabel is a human readable short name of the log type queried by the query task.
	LogTypeLabel string
	// LogTypeColorCode is the hex color code without the `#` prefix for the log type.
	LogTypeColorCode string
	// SampleQuery is an example query string used in this query task.
	SampleQuery string
}

// FeatureOutputTimelineElement is a model type for one of relationship type of timelines that can be related to this feature.
type FeatureOutputTimelineElement struct {
	// RelationshipID is the unique name of the relationship type.
	RelationshipID string
	// RelationshipColorCode is the hex color code without the `#` prefix for the relationship type.
	RelationshipColorCode string
	// LongName is the human readable name of the relationship
	LongName string
	// Label is the short name of the timeline. This is also used in the chip on the left side of timelines.
	Label string
	// Description is the string explains the relationship.
	Description string
}

// FeatureAvailableInspectionType is a model type for a InspectionType that supports the feature task.
type FeatureAvailableInspectionType struct {
	// ID is the unique name of the InspectionType.
	ID string
	// Name is the human readable name of the InspectionType.
	Name string
}

// GetFeatureDocumentModel returns the document model for feature tasks from the task server.
func GetFeatureDocumentModel(taskServer *inspection.InspectionTaskServer) (*FeatureDocumentModel, error) {
	result := FeatureDocumentModel{}
	features := task.Subset(taskServer.RootTaskSet, filter.NewEnabledFilter(inspection_task.LabelKeyInspectionFeatureFlag, false))
	for _, feature := range features.GetAll() {
		indirectQueryDependencyElement := []FeatureIndirectDependentQueryElement{}
		targetQueryDependencyElement := FeatureDependentTargetQueryElement{}
		targetLogTypeKey := typedmap.GetOrDefault(feature.Labels(), inspection_task.LabelKeyFeatureTaskTargetLogType, enum.LogTypeUnknown)

		// Get query related tasks in the dependency of this feature.
		queryTasksInDependency, err := getDependentQueryTasks(taskServer, feature)
		if err != nil {
			return nil, err
		}
		for _, queryTask := range queryTasksInDependency {
			logTypeKey := typedmap.GetOrDefault(queryTask.Labels(), label.TaskLabelKeyQueryTaskTargetLogType, enum.LogTypeUnknown)
			if targetLogTypeKey != logTypeKey {
				logType := enum.LogTypes[logTypeKey]
				indirectQueryDependencyElement = append(indirectQueryDependencyElement, FeatureIndirectDependentQueryElement{
					ID:               queryTask.UntypedID().String(),
					LogTypeLabel:     logType.Label,
					LogTypeColorCode: strings.TrimLeft(logType.LabelBackgroundColor, "#"),
				})
			} else {
				targetQueryDependencyElement = FeatureDependentTargetQueryElement{
					ID:               queryTask.UntypedID().String(),
					LogTypeLabel:     enum.LogTypes[targetLogTypeKey].Label,
					LogTypeColorCode: strings.TrimLeft(enum.LogTypes[targetLogTypeKey].LabelBackgroundColor, "#"),
					SampleQuery:      typedmap.GetOrDefault(queryTask.Labels(), label.TaskLabelKeyQueryTaskSampleQuery, ""),
				}
			}
		}

		formElements := []FeatureDependentFormElement{}
		formTasks, err := getDependentFormTasks(taskServer, feature)
		if err != nil {
			return nil, err
		}
		for _, formTask := range formTasks {
			formElements = append(formElements, FeatureDependentFormElement{
				ID:          formTask.UntypedID().String(),
				Label:       typedmap.GetOrDefault(formTask.Labels(), label.TaskLabelKeyFormFieldLabel, ""),
				Description: typedmap.GetOrDefault(formTask.Labels(), label.TaskLabelKeyFormFieldDescription, ""),
			})
		}

		outputTimelines := []FeatureOutputTimelineElement{}
		for i := 0; i < enum.EnumParentRelationshipLength; i++ {
			relationshipKey := enum.ParentRelationship(i)
			relationship := enum.ParentRelationships[relationshipKey]

			isRelated := false
			for _, event := range relationship.GeneratableEvents {
				if event.SourceLogType == targetLogTypeKey {
					isRelated = true
					break
				}
			}
			for _, revision := range relationship.GeneratableRevisions {
				if revision.SourceLogType == targetLogTypeKey {
					isRelated = true
					break
				}
			}
			for _, alias := range relationship.GeneratableAliasTimelineInfo {
				if alias.SourceLogType == targetLogTypeKey {
					isRelated = true
					break
				}
			}
			if isRelated {
				outputTimelines = append(outputTimelines, FeatureOutputTimelineElement{
					RelationshipID:        relationship.EnumKeyName,
					RelationshipColorCode: strings.TrimLeft(relationship.LabelBackgroundColor, "#"),
					LongName:              relationship.LongName,
					Label:                 relationship.Label,
					Description:           relationship.Description,
				})
			}
		}

		result.Features = append(result.Features, FeatureDocumentElement{
			ID:                       feature.UntypedID().String(),
			Name:                     typedmap.GetOrDefault(feature.Labels(), inspection_task.LabelKeyFeatureTaskTitle, ""),
			Description:              typedmap.GetOrDefault(feature.Labels(), inspection_task.LabelKeyFeatureTaskDescription, ""),
			IndirectQueryDependency:  indirectQueryDependencyElement,
			TargetQueryDependency:    targetQueryDependencyElement,
			Forms:                    formElements,
			OutputTimelines:          outputTimelines,
			AvailableInspectionTypes: getAvailableInspectionTypes(taskServer, feature),
		})

	}
	return &result, nil
}

// getDependentQueryTasks returns the list of query tasks required by the feature task.
func getDependentQueryTasks(taskServer *inspection.InspectionTaskServer, featureTask task.UntypedTask) ([]task.UntypedTask, error) {
	resolveSource, err := task.NewTaskSet([]task.UntypedTask{featureTask})
	if err != nil {
		return nil, err
	}
	resolved, err := resolveSource.ResolveTask(taskServer.RootTaskSet)
	if err != nil {
		return nil, err
	}
	return task.Subset(resolved, filter.NewEnabledFilter(label.TaskLabelKeyIsQueryTask, false)).GetAll(), nil
}

// getDependentFormTasks returns the list of form tasks required by the feature task.
func getDependentFormTasks(taskServer *inspection.InspectionTaskServer, featureTask task.UntypedTask) ([]task.UntypedTask, error) {
	resolveSource, err := task.NewTaskSet([]task.UntypedTask{featureTask})
	if err != nil {
		return nil, err
	}
	resolved, err := resolveSource.ResolveTask(taskServer.RootTaskSet)
	if err != nil {
		return nil, err
	}
	return task.Subset(resolved, filter.NewEnabledFilter(label.TaskLabelKeyIsFormTask, false)).GetAll(), nil
}

// getAvailableInspectionTypes returns the list of information about inspection type that supports this feature.
func getAvailableInspectionTypes(taskServer *inspection.InspectionTaskServer, featureTask task.UntypedTask) []FeatureAvailableInspectionType {
	result := []FeatureAvailableInspectionType{}
	inspectionTypes := taskServer.GetAllInspectionTypes()
	for _, inspectionType := range inspectionTypes {
		labels, found := typedmap.Get(featureTask.Labels(), inspection_task.LabelKeyInspectionTypes)
		if found && slices.Contains(labels, inspectionType.Id) {
			result = append(result, FeatureAvailableInspectionType{
				ID:   inspectionType.Id,
				Name: inspectionType.Name,
			})
		}
	}
	return result
}
