package config

import (
	"os"
	"testing"
)

func TestLoadEnv(t *testing.T) {
	testConfig := &Config{
		Env: map[string]string{},
	}

	if err := testConfig.LoadEnv([]string{"EMPTY"}).Err(); err == nil {
		t.Errorf("Empty env variable should return an error, but error is nil")
	}

	// Reset error
	testConfig.err = nil

	env := "EXISTS"
	val := "yes"
	os.Setenv(env, val)
	if err := testConfig.LoadEnv([]string{env}).Err(); err != nil {
		t.Errorf("Existing env variable should not return an error, but got %s", err)
	}

	if testConfig.Env[env] != val {
		t.Errorf("Config should contain %s with value %s", env, val)
	}
}

func TestLoadConfig(t *testing.T) {
	testApp := &Config{}
	os.Setenv("MAX_FILES_PER_USER", "10")
	os.Setenv("MAX_FILE_SIZE", "100000")
	os.Setenv("MAX_TOTAL_SIZE_PER_USER", "1000000")
	os.Setenv("ORIGINAL_FILE_EXTENTION", "mp3")
	os.Setenv("TARGET_FILE_EXTENTION", "ogg,aac")
	os.Setenv("ORIGINAL_FILES_PATH", "/tmp")
	os.Setenv("CONVERTED_FILES_PATH", "/tmp")
	defer os.Unsetenv("MAX_FILES_PER_USER")
	defer os.Unsetenv("MAX_FILE_SIZE")
	defer os.Unsetenv("MAX_TOTAL_SIZE_PER_USER")
	defer os.Unsetenv("ORIGINAL_FILE_EXTENTION")
	defer os.Unsetenv("TARGET_FILE_EXTENTION")
	defer os.Unsetenv("ORIGINAL_FILES_PATH")
	defer os.Unsetenv("CONVERTED_FILES_PATH")

	err := testApp.LoadConfigs().Err()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if testApp.MaxFilesPerUser != 10 {
		t.Errorf("Expected MaxFilesPerUser to be %d, got %d", 10, testApp.MaxFilesPerUser)
	}
	if testApp.MaxFileSize != 100000 {
		t.Errorf("Expected MaxSize to be %d, got %d", 100000, testApp.MaxFileSize)
	}
	if testApp.MaxTotalSizePerUser != 1000000 {
		t.Errorf("Expected MaxSizePerUser to be %d, got %d", 1000000, testApp.MaxTotalSizePerUser)
	}
	if len(testApp.OriginFileExtention) != 1 {
		t.Errorf("Expected OriginFileExtention to have %d items, it has %d", 1, len(testApp.OriginFileExtention))
	}
	if testApp.OriginFileExtention[0] != "mp3" {
		t.Errorf("Expected first OriginFileExtention to be %s, got %s", "mp3", testApp.OriginFileExtention[0])
	}
	if len(testApp.TargetFileExtention) != 2 {
		t.Errorf("Expected TargetFileExtention to bhave %d items, it has %d", 2, len(testApp.TargetFileExtention))
	}

	if testApp.TargetFileExtention[0] != "ogg" {
		t.Errorf("Expected first TargetFileExtention to be %s, got %s", "ogg", testApp.TargetFileExtention[0])
	}
	if testApp.TargetFileExtention[1] != "aac" {
		t.Errorf("Expected first TargetFileExtention to be %s, got %s", "aac", testApp.TargetFileExtention[1])
	}

}
