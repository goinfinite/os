package apiController

import (
	"testing"
)

func TestRuntimeControllerParsePhpModules(t *testing.T) {
	type expectedModule struct {
		moduleName   string
		moduleStatus bool
	}
	type expectedFailure struct {
		moduleIndex     uint
		moduleName      string
		hasModuleName   bool
		moduleStatus    bool
		hasModuleStatus bool
		failureReason   string
	}

	testCases := []struct {
		name             string
		rawPhpModules    any
		expectedModules  []expectedModule
		expectedFailures []expectedFailure
	}{
		{
			name: "ValidPhpModuleList",
			rawPhpModules: []any{
				map[string]any{"name": "curl", "status": true},
				map[string]any{"name": "mysqli", "status": false},
			},
			expectedModules: []expectedModule{
				{moduleName: "curl", moduleStatus: true},
				{moduleName: "mysqli", moduleStatus: false},
			},
		},
		{
			name:          "SinglePhpModuleObject",
			rawPhpModules: map[string]any{"name": "opcache", "status": true},
			expectedModules: []expectedModule{
				{moduleName: "opcache", moduleStatus: true},
			},
		},
		{
			name: "InvalidPhpModuleName",
			rawPhpModules: []any{
				map[string]any{"name": "bad/name", "status": true},
			},
			expectedFailures: []expectedFailure{
				{
					moduleIndex:     0,
					moduleStatus:    true,
					hasModuleStatus: true,
					failureReason:   "InvalidPhpModuleName",
				},
			},
		},
		{
			name: "InvalidPhpModuleStatus",
			rawPhpModules: []any{
				map[string]any{"name": "mysqli", "status": "maybe"},
			},
			expectedFailures: []expectedFailure{
				{
					moduleIndex:   0,
					moduleName:    "mysqli",
					hasModuleName: true,
					failureReason: "InvalidPhpModuleStatus",
				},
			},
		},
		{
			name:          "InvalidPhpModuleStructure",
			rawPhpModules: []any{"invalid-module"},
			expectedFailures: []expectedFailure{
				{
					moduleIndex:   0,
					failureReason: "InvalidPhpModuleStructure",
				},
			},
		},
		{
			name:          "InvalidPhpModulesStructure",
			rawPhpModules: "invalid-modules",
			expectedFailures: []expectedFailure{
				{
					moduleIndex:   0,
					failureReason: "InvalidPhpModulesStructure",
				},
			},
		},
		{
			name: "KeepsValidModulesAndReportsFailures",
			rawPhpModules: []any{
				map[string]any{"name": "curl", "status": true},
				map[string]any{"name": "bad/name", "status": true},
				map[string]any{"name": "mysqli", "status": "maybe"},
				"invalid-module",
			},
			expectedModules: []expectedModule{
				{moduleName: "curl", moduleStatus: true},
			},
			expectedFailures: []expectedFailure{
				{
					moduleIndex:     1,
					moduleStatus:    true,
					hasModuleStatus: true,
					failureReason:   "InvalidPhpModuleName",
				},
				{
					moduleIndex:   2,
					moduleName:    "mysqli",
					hasModuleName: true,
					failureReason: "InvalidPhpModuleStatus",
				},
				{
					moduleIndex:   3,
					failureReason: "InvalidPhpModuleStructure",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			modules, parsingFailures := (&RuntimeController{}).parsePhpModules(
				testCase.rawPhpModules,
			)

			if len(modules) != len(testCase.expectedModules) {
				t.Fatalf(
					"ValidModulesCountMismatch: expected %d, got %d",
					len(testCase.expectedModules), len(modules),
				)
			}
			for moduleIndex, expectedModule := range testCase.expectedModules {
				actualModule := modules[moduleIndex]
				if actualModule.Name.String() != expectedModule.moduleName {
					t.Errorf(
						"ModuleNameMismatch: module %d: expected %q, got %q",
						moduleIndex,
						expectedModule.moduleName, actualModule.Name.String(),
					)
				}
				if actualModule.Status != expectedModule.moduleStatus {
					t.Errorf(
						"ModuleStatusMismatch: module %d: expected %t, got %t",
						moduleIndex,
						expectedModule.moduleStatus, actualModule.Status,
					)
				}
			}

			if len(parsingFailures) != len(testCase.expectedFailures) {
				t.Fatalf(
					"ParsingFailuresCountMismatch: expected %d, got %d",
					len(testCase.expectedFailures), len(parsingFailures),
				)
			}
			for failureIndex, expectedFailure := range testCase.expectedFailures {
				actualFailure := parsingFailures[failureIndex]
				if actualFailure.Index != expectedFailure.moduleIndex {
					t.Errorf(
						"FailureIndexMismatch: failure %d: expected module index %d, got %d",
						failureIndex,
						expectedFailure.moduleIndex, actualFailure.Index,
					)
				}

				hasActualModuleName := actualFailure.Name != nil
				if expectedFailure.hasModuleName != hasActualModuleName {
					t.Errorf(
						"FailureModuleNamePresenceMismatch: failure %d: expected %t, got %t",
						failureIndex, expectedFailure.hasModuleName, hasActualModuleName,
					)
				}
				if expectedFailure.hasModuleName && hasActualModuleName &&
					actualFailure.Name.String() != expectedFailure.moduleName {
					t.Errorf(
						"FailureModuleNameMismatch: failure %d: expected %q, got %q",
						failureIndex,
						expectedFailure.moduleName, actualFailure.Name.String(),
					)
				}

				hasActualModuleStatus := actualFailure.Status != nil
				if expectedFailure.hasModuleStatus != hasActualModuleStatus {
					t.Errorf(
						"FailureModuleStatusPresenceMismatch: failure %d: expected %t, got %t",
						failureIndex, expectedFailure.hasModuleStatus,
						hasActualModuleStatus,
					)
				}
				if expectedFailure.hasModuleStatus && hasActualModuleStatus &&
					*actualFailure.Status != expectedFailure.moduleStatus {
					t.Errorf(
						"FailureModuleStatusMismatch: failure %d: expected %t, got %t",
						failureIndex, expectedFailure.moduleStatus, *actualFailure.Status,
					)
				}
				if actualFailure.Reason.String() != expectedFailure.failureReason {
					t.Errorf(
						"FailureReasonMismatch: failure %d: expected %q, got %q",
						failureIndex,
						expectedFailure.failureReason,
						actualFailure.Reason.String(),
					)
				}
			}
		})
	}
}
