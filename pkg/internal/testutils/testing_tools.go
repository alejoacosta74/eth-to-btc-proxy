package testutils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alejoacosta74/rpc-proxy/pkg/types"
)

const DATA_PATH = "../internal/testdata/"
const DATA_FILE_INPUT_SUFFIX = "input"
const DATA_FILE_OUTPUT_SUFFIX = "golden"
const MOCKQTUMD_PREFIX = "mockqtumd_"
const MOCKQTUMD_SUFFIX = "_data"

var jsonReq = new(types.JSONRPCRequest)
var jsonResp = new(types.JSONRPCResponse)

func HandleFatalError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Fatal error: %+v", err)
	}
}

// readTestNames returns reads the folder names in the testdata folder
func ReadTestNames(t *testing.T) []string {
	t.Helper()
	testNames := make([]string, 0)
	entries, err := os.ReadDir(DATA_PATH)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() {
			testNames = append(testNames, entry.Name())
		}
	}
	return testNames
}

// getTestFilePath returns the path to the test file for
// the given test directory, method and name space (prefix) and suffix
func GetTestFilePath(testdir string, prefix string, method string, suffix string) string {
	filename := prefix + method + suffix + ".json"
	source := filepath.Join(DATA_PATH, testdir, filename)
	return source
}

// decodeRawParams unmarshals the given JSON rawparams into the given params interfaces
func DecodeRawParams(t *testing.T, rawParams []json.RawMessage, params ...interface{}) {
	for i, p := range rawParams {
		err := json.Unmarshal(p, params[i])
		if err != nil {
			t.Fatal(err)
		}
	}
}

// decodeRawResult unmarshals the given JSON rawResult into the given result interface
func DecodeRawResult(t *testing.T, rawResult json.RawMessage, result interface{}) {
	err := json.Unmarshal(rawResult, result)
	if err != nil {
		t.Fatal(err)
	}
}

// LoadJSON reads the given JSON data file and unmarshals the
// content into the given interface
func LoadJSON(t *testing.T, path string, jsonMsg interface{}) {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(content, jsonMsg)
	if err != nil {
		t.Fatal(err)
	}
}

// LoadFromFile reads the given JSON data file and unmarshals the
// content into the given interface to be used as test input or output data
func LoadFromFile(t *testing.T, testname string, prefix string, method string, suffix string, data interface{}) {
	path := GetTestFilePath(testname, prefix, method, suffix)
	LoadJSON(t, path, &data)

}

// LoadDataFromFile reads the given JSON data file and unmarshals the
// content into the given interface to be used as test input or output data
func LoadDataFromFile(t *testing.T, path string, data interface{}) {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(content, data)
	if err != nil {
		t.Fatal(err)
	}

}

// GetTestFilePaths returns the list of input and
// output (golden) files for a given test name, module and method to run the test
//
// Params:
//
// - testdir: the name of the test directory
//
// - module: the name of the module to test (e.g. rpc, qtum, wallet)
//
// - method: the name of the method to test (e.g. listunspent)
func GetTestFilePaths(t *testing.T, testdir string, module string, method string) (inputPaths []string, goldenPaths []string) {
	inputPaths, err := getTestFilePaths(testdir, module, method, DATA_FILE_INPUT_SUFFIX)
	HandleFatalError(t, err)
	goldenPaths, err = getTestFilePaths(testdir, module, method, DATA_FILE_OUTPUT_SUFFIX)
	HandleFatalError(t, err)
	return inputPaths, goldenPaths

}

func getTestFilePaths(testdir string, module string, method string, suffix string) ([]string, error) {
	filenames := module + "_" + method + "_" + suffix + "*.json"
	source := filepath.Join(DATA_PATH, testdir, filenames)
	paths, err := filepath.Glob(source)
	if err != nil {
		return nil, err
	}
	return paths, nil
}

type TestingFn = func(t *testing.T, testname string, i int, inputPath string, goldenPath []string, input, want interface{})

// RunTests is a helper function that runs the received testing function
// in a loop against the available testing data files
func RunTests(t *testing.T, fn TestingFn, module string, method string, input interface{}, want interface{}) {

	testnames := ReadTestNames(t)
	for _, testname := range testnames {

		// Get the input and golden test data file
		inputPaths, goldenPaths := GetTestFilePaths(t, testname, module, method)

		for i, inputPath := range inputPaths {
			// Load the input data for this particular test
			LoadDataFromFile(t, inputPath, input)
			// Load the golden data for this particular test
			LoadDataFromFile(t, goldenPaths[i], want)
			fn(t, testname, i, inputPath, goldenPaths, input, want)
		}
	}

}
