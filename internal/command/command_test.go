package command

import (
	"os"
	"testing"
)

func TestParseArgs(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		expectedEditor string
		expectedEnv     string
		expectedError   error
	}{
		{
			name:          "No arguments",
			args:          []string{},
			expectedEditor: "",
			expectedEnv:     "",
			expectedError:   nil,
		},
		{
			name:          "One argument (editor)",
			args:          []string{"vim"},
			expectedEditor: "vim",
			expectedEnv:     "",
			expectedError:   nil,
		},
		{
			name:          "Two arguments (editor and env)",
			args:          []string{"vim", "prod"},
			expectedEditor: "vim",
			expectedEnv:     "prod",
			expectedError:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save and restore original command-line arguments
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set command-line arguments for the test case
			os.Args = append([]string{"rules_cli"}, tc.args...)

			editor, env, err := ParseArgs()

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if editor != tc.expectedEditor {
				t.Errorf("Expected editor %q, got %q", tc.expectedEditor, editor)
			}

			if env != tc.expectedEnv {
				t.Errorf("Expected env %q, got %q", tc.expectedEnv, env)
			}
		})
	}
}
