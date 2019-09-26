/*
 * Copyright 2019 The original author or authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package command

import (
	"os"
	"path/filepath"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/projectriff/libfnbuildpack/function"
)

const executable = 0100

type Buildpack struct{}

func (b Buildpack) Build(build build.Build) (int, error) {
	if f, ok, err := NewFunction(build); err != nil {
		return build.Failure(function.Error_ComponentInitialization), err
	} else if ok {
		if err := f.Contribute(); err != nil {
			return build.Failure(function.Error_ComponentContribution), err
		}
	}

	if invoker, ok, err := NewInvoker(build); err != nil {
		return build.Failure(function.Error_ComponentInitialization), err
	} else if ok {
		if err := invoker.Contribute(); err != nil {
			return build.Failure(function.Error_ComponentContribution), err
		}
	}

	return build.Success()
}

func (b Buildpack) Detect(detect detect.Detect, metadata function.Metadata) (int, error) {
	if metadata.Artifact == "" {
		return detect.Fail(), nil
	}

	path := filepath.Join(detect.Application.Root, metadata.Artifact)

	ok, err := helper.FileExists(path)
	if err != nil || !ok {
		return detect.Error(function.Error_ComponentInternal), err
	}

	info, err := os.Stat(path)
	if err != nil {
		return detect.Error(function.Error_ComponentInternal), err
	}

	if !b.executable(info) {
		detect.Logger.Debug("Disregarding %q for the 'command' invoker, as it does not have executable permission", path)
		return detect.Error(function.Error_ComponentInternal), nil
	}

	return detect.Pass(buildplan.Plan{
		Provides: []buildplan.Provided{
			{Name: Dependency},
		},
		Requires: []buildplan.Required{
			{
				Name: Dependency,
				Metadata: map[string]interface{}{
					Command: metadata.Artifact,
				},
			},
		},
	})
}

func (b Buildpack) executable(fileInfo os.FileInfo) bool {
	return fileInfo.Mode().IsRegular() && (fileInfo.Mode().Perm()&executable == executable)
}

func (b Buildpack) Id() string {
	return "command"
}
