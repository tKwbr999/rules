package handler

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestGetRulesPath(t *testing.T) {
	testCases := []struct {
		name        string
		rulesPath   string
		expectedPath string
		expectError bool
	}{
		{
			name:        "RULES_PATH is set",
			rulesPath:   "/path/to/rules",
			expectedPath: "/path/to/rules",
			expectError: false,
		},
		{
			name:        "RULES_PATH is not set",
			rulesPath:   "",
			expectedPath: "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.rulesPath != "" {
				os.Setenv("RULES_PATH", tc.rulesPath)
				defer os.Unsetenv("RULES_PATH")
			} else {
				os.Unsetenv("RULES_PATH")
			}

			path, err := GetRulesPath()

			if tc.expectError {
				if err == nil {
					t.Errorf("Test Case: %s, GetRulesPath() should return an error", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("Test Case: %s, GetRulesPath() error = %v", tc.name, err)
				}
				if path != tc.expectedPath {
					t.Errorf("Test Case: %s, GetRulesPath() = %v, want %v", tc.name, path, tc.expectedPath)
				}
			}
		})
	}
}

func TestGetMdFiles(t *testing.T) {
	testCases := []struct {
		name          string
		env           string
		testFiles     []string
		expectedFiles []string
		expectError   bool
	}{
		{
			name:          "No environment specified",
			env:           "",
			testFiles:     []string{"test.md", "dev.md", "prod.md", "other.txt"},
			expectedFiles: []string{"test.md", "dev.md", "prod.md"},
			expectError:   false,
		},
		{
			name:          "Environment specified",
			env:           "dev",
			testFiles:     []string{"test.md", "dev.md", "prod.md", "other.txt"},
			expectedFiles: []string{"dev.md"},
			expectError:   false,
		},
		{
			name:          "No files found",
			env:           "",
			testFiles:     []string{"other.txt"},
			expectedFiles: []string{},
			expectError:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// テスト用のディレクトリを作成
			rulesPath := t.TempDir()

			// テスト用のファイルを作成
			for _, file := range tc.testFiles {
				if err := os.WriteFile(filepath.Join(rulesPath, file), []byte("test content"), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			files, err := GetMdFiles(rulesPath, tc.env)

			if tc.expectError {
				if err == nil {
					t.Errorf("Test Case: %s, GetMdFiles() should return an error", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("Test Case: %s, GetMdFiles() error = %v", tc.name, err)
				}

				// 絶対パスに変換
				var expectedFiles []string
				for _, file := range tc.expectedFiles {
					expectedFiles = append(expectedFiles, filepath.Join(rulesPath, file))
				}

				sort.Strings(files)
				sort.Strings(expectedFiles)

				if !reflect.DeepEqual(files, expectedFiles) {
					t.Errorf("Test Case: %s, GetMdFiles() = %v, want %v", tc.name, files, expectedFiles)
				}
			}
		})
	}
}

func TestCombineFiles(t *testing.T) {
	testCases := []struct {
		name            string
		testFiles       map[string]string
		nonExistentFile bool
		expectedContent string
		expectError     bool
	}{
		{
			name: "Files exist",
			testFiles: map[string]string{
				"file1.txt": "content of file1",
				"file2.txt": "content of file2",
			},
			nonExistentFile: false,
			expectedContent: "content of file1\n\ncontent of file2\n\n",
			expectError:     false,
		},
		{
			name:            "File does not exist",
			testFiles:       map[string]string{},
			nonExistentFile: true,
			expectedContent: "",
			expectError:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var files []string
			for name, content := range tc.testFiles {
				tmpFile := filepath.Join(t.TempDir(), name)
				if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				files = append(files, tmpFile)
			}

			if tc.nonExistentFile {
				files = []string{"nonexistent_file.txt"}
			}

			combinedContent, err := CombineFiles(files)

			if tc.expectError {
				if err == nil {
					t.Errorf("Test Case: %s, CombineFiles() should return an error", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("Test Case: %s, CombineFiles() error = %v", tc.name, err)
				}

				if combinedContent != tc.expectedContent {
					t.Errorf("Test Case: %s, CombineFiles() = %v, want %v", tc.name, combinedContent, tc.expectedContent)
				}
			}
		})
	}
}

func TestGetOutputPath(t *testing.T) {
	testCases := []struct {
		name          string
		editor        string
		expectedSuffix string
		expectError   bool
	}{
		{
			name:          "Editor is specified",
			editor:        "vscode",
			expectedSuffix: ".vscoderules",
			expectError:   false,
		},
		{
			name:          "Editor is empty",
			editor:        "",
			expectedSuffix: ".rules",
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outputPath, err := GetOutputPath(tc.editor)

			if err != nil {
				t.Errorf("Test Case: %s, GetOutputPath() error = %v", tc.name, err)
			}
			if !strings.HasSuffix(outputPath, tc.expectedSuffix) {
				t.Errorf("Test Case: %s, GetOutputPath() = %v, want suffix %v", tc.name, outputPath, tc.expectedSuffix)
			}
		})
	}
}
