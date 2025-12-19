package usecases

import (
	"errors"
	"testing"
)

// TestSecretDetectorService_ExecuteHomePathDetection_Normal はホームパス検知実行の正常系テスト
func TestSecretDetectorService_ExecuteHomePathDetection_Normal(t *testing.T) {
	testCases := []struct {
		name          string
		configFile    string
		mockOutput    string
		mockError     error
		expectedExit  int
		expectedError bool
	}{
		{
			name:          "No files found",
			configFile:    "",
			mockOutput:    "",
			mockError:     nil,
			expectedExit:  0,
			expectedError: false,
		},
		{
			name:          "Git command error",
			configFile:    "",
			mockOutput:    "",
			mockError:     errors.New("git command failed"),
			expectedExit:  1,
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockExecutor := &MockCommandExecutor{
				ExecuteFunc: func(name string, args ...string) ([]byte, error) {
					if tc.mockError != nil {
						return nil, tc.mockError
					}
					return []byte(tc.mockOutput), nil
				},
			}
			mockOutputWriter := &MockOutputWriter{}
			service := NewSecretDetectorService(false, false, tc.configFile, mockExecutor, mockOutputWriter)

			exitCode, err := service.ExecuteHomePathDetection()

			if tc.expectedError {
				if err == nil {
					t.Errorf("ExecuteHomePathDetection() expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("ExecuteHomePathDetection() unexpected error: %v", err)
				return
			}

			if exitCode != tc.expectedExit {
				t.Errorf("ExecuteHomePathDetection() returned exit code %d, expected %d", exitCode, tc.expectedExit)
			}
		})
	}
}
