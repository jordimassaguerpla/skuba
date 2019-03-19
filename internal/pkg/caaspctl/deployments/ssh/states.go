/*
 * Copyright (c) 2019 SUSE LLC. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package ssh

import (
	"log"
)

var (
	stateMap = map[string]Runner{}
)

type Runner func(t *Target, data interface{}) error

func (t *Target) Apply(data interface{}, states ...string) error {
	for _, stateName := range states {
		log.Printf("=== applying state %s ===\n", stateName)
		if state, stateExists := stateMap[stateName]; stateExists {
			if err := state(t, data); err != nil {
				log.Printf("=== failed to apply state %s: %v ===\n", stateName, err)
			} else {
				log.Printf("=== state %s applied successfully ===\n", stateName)
			}
		} else {
			log.Fatalf("state does not exist: %s", stateName)
		}
	}
	return nil
}