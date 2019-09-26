/*
 * Copyright 2018-2019 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package command

import (
	"fmt"
	"path/filepath"

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// Command is the key identifying the command executable in the build plan.
const Command = "command"

// Function represents the function to be executed.
type Function struct {
	application application.Application
	executable  string
	layer       layers.Layer
}

// Contributes makes the contribution to the launch layer.
func (f Function) Contribute() error {
	return f.layer.Contribute(marker{"Command", f.executable}, func(layer layers.Layer) error {
		return layer.OverrideLaunchEnv("FUNCTION_URI", filepath.Join(f.application.Root, f.executable))
	}, layers.Launch)
}

// NewFunction creates a new instance returning true if the riff-invoker-command plan exists.
func NewFunction(build build.Build) (Function, bool, error) {
	p, ok, err := build.Plans.GetShallowMerged(Dependency)
	if err != nil {
		return Function{}, false, err
	}
	if !ok {
		return Function{}, false, nil
	}

	exec, ok := p.Metadata[Command].(string)
	if !ok {
		return Function{}, false, fmt.Errorf("command metadata of incorrect type: %v", p.Metadata[Command])
	}

	return Function{
		build.Application,
		exec,
		build.Layers.Layer("command-function"),
	}, true, nil
}

type marker struct {
	Type       string `toml:"type"`
	Executable string `toml:"executable"`
}

func (m marker) Identity() (string, string) {
	return m.Type, m.Executable
}
