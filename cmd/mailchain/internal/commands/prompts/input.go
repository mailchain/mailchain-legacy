// Copyright 2019 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package prompts

import (
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
)

// RequiredInput create a prompt and require that an input is supplied
func RequiredInput(label string) (string, error) {
	return RequiredInputWithDefault(label, "")
}

// RequiredInputWithDefault create a prompt and require that an input is supplied, specify a default value
func RequiredInputWithDefault(label, defaultValue string) (string, error) {
	prompt := promptui.Prompt{
		Label:   label,
		Default: defaultValue,
		Validate: func(val string) error {
			if strings.TrimSpace(val) == "" {
				return errors.Errorf("value required")
			}
			return nil
		},
	}
	return prompt.Run()
}
