package command

import (
	"os"
	"reflect"
	"testing"
)

func TestParseArgs(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expectedEditor string
		expectedFiles  []string
		expectedError  error
	}{
		{
			name:           "No arguments",
			args:           []string{},
			expectedEditor: "",
			expectedFiles:  []string{},
			expectedError:  nil,
		},
		{
			name:           "One argument (editor)",
			args:           []string{"vim"},
			expectedEditor: "vim",
			expectedFiles:  []string{},
			expectedError:  nil,
		},
		{
			name:           "Two arguments (editor and file)",
			args:           []string{"vim", "file1"},
			expectedEditor: "vim",
			expectedFiles:  []string{"file1"},
			expectedError:  nil,
		},
		{
			name:           "Multiple arguments (editor and files)",
			args:           []string{"vim", "file1", "file2", "file3"},
			expectedEditor: "vim",
			expectedFiles:  []string{"file1", "file2", "file3"},
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Save and restore original command-line arguments
			originalArgs := os.Args
			defer func() { os.Args = originalArgs }()

			// Set command-line arguments for the test case
			os.Args = append([]string{"rules_cli"}, tc.args...)

			editor, files, err := ParseArgs()

			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if editor != tc.expectedEditor {
				t.Errorf("Expected editor %q, got %q", tc.expectedEditor, editor)
			}

			if !reflect.DeepEqual(files, tc.expectedFiles) {
				t.Errorf("Expected files %v, got %v", tc.expectedFiles, files)
			}
		})
	}
}
