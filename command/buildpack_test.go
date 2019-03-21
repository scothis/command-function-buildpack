/*
 * Copyright 2018 the original author or authors.
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

package command_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/test"
	. "github.com/onsi/gomega"
	"github.com/projectriff/command-function-buildpack/command"
	"github.com/projectriff/riff-buildpack/function"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestName(t *testing.T) {
	spec.Run(t, "Id", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		it("has the right id", func() {
			b := command.NewBuildpack()

			g.Expect(b.Id()).To(Equal("command"))
		})
	}, spec.Report(report.Terminal{}))
}

func TestDetect(t *testing.T) {
	spec.Run(t, "Detect", func(t *testing.T, _ spec.G, it spec.S) {

		g := NewGomegaWithT(t)

		var f *test.DetectFactory
		var m function.Metadata
		var b function.Buildpack

		it.Before(func() {
			f = test.NewDetectFactory(t)
			m = function.Metadata{}
			b = command.NewBuildpack()
		})

		it("fails by default", func() {
			plan, err := b.Detect(f.Detect, m)

			g.Expect(err).To(BeNil())
			g.Expect(plan).To(BeNil())
		})

		it("passes if the artifact is executable", func() {
			err := os.MkdirAll(f.Detect.Application.Root, 0755)
			g.Expect(err).To(BeNil())
			defer os.RemoveAll(f.Detect.Application.Root)
			tmpfile, err := ioutil.TempFile(f.Detect.Application.Root, "example")
			g.Expect(err).To(BeNil())
			tmpfile.Chmod(0755)
			artifact, err := filepath.Rel(f.Detect.Application.Root, tmpfile.Name())
			g.Expect(err).To(BeNil())

			m.Artifact = artifact

			plan, err := b.Detect(f.Detect, m)

			g.Expect(err).To(BeNil())
			g.Expect(plan).To(Equal(&buildplan.BuildPlan{
				command.Dependency: buildplan.Dependency{
					Metadata: buildplan.Metadata{"command": artifact},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}

func TestBuild(t *testing.T) {
	spec.Run(t, "Build", func(t *testing.T, _ spec.G, it spec.S) {
		g := NewGomegaWithT(t)

		var f *test.BuildFactory
		var b function.Buildpack

		it.Before(func() {
			f = test.NewBuildFactory(t)
			b = command.NewBuildpack()
		})

		it("won't build unless passed detection", func() {
			err := b.Build(f.Build)

			g.Expect(err).To(MatchError("buildpack passed detection but did not know how to actually build"))
		})

		it.Pend("will build if passed detection", func() {
			f.AddBuildPlan(command.Dependency, buildplan.Dependency{})
			f.AddDependency(command.Dependency, ".")

			err := b.Build(f.Build)

			g.Expect(err).To(BeNil())
		})
	}, spec.Report(report.Terminal{}))
}
