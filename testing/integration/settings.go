package integration

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

// Given two maps, recursively merge right into left, NEVER replacing any key that already exists in left
func mergeKeys(left, right map[string]interface{}) map[string]interface{} {
	for key, rightVal := range right {
		if leftVal, present := left[key]; present {
			//then we don't want to replace it - recurse
			left[key] = mergeKeys(leftVal.(map[string]interface{}), rightVal.(map[string]interface{}))
		} else {
			// key not in left so we can just shove it in
			left[key] = rightVal
		}
	}

	return left
}

func createTestCaseSettingsFile(baseSettings map[string]interface{}, testPath string) (string, error) {
	settings := map[string]interface{}{
		"mailboxState": map[string]interface{}{
			"badgerdb": map[string]string{
				"path": strings.Join([]string{testPath, "mailbox"}, "/"),
			},
		},
		"keystore": map[string]interface{}{
			"kind": "nacl-filestore",
			"nacl-filestore": map[string]string{
				"path": strings.Join([]string{testPath, "keystore"}, "/"),
			},
		},
		"cache": map[string]string{
			"path":    strings.Join([]string{testPath, "message-cache"}, "/"),
			"timeout": "1h",
		},
		"fetcher": map[string]interface{}{
			"disabled": true,
		},
	}

	d, err := yaml.Marshal(settings)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	settingsFileName := strings.Join([]string{testPath, "settings.yaml"}, "/")

	return settingsFileName, ioutil.WriteFile(settingsFileName, d, 0644)
}
