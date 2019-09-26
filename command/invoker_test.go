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

package command_test

import (
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/projectriff/command-function-buildpack/command"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestInvoker(t *testing.T) {
	spec.Run(t, "Invoker", func(t *testing.T, _ spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)
		})

		it("returns true if build plan exists", func() {
			f.AddDependency(command.Dependency, filepath.Join("testdata", "stub-invoker.tgz"))
			f.AddPlan(buildpackplan.Plan{Name: command.Dependency})

			_, ok, err := command.NewInvoker(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(ok).To(gomega.BeTrue())
		})

		it("returns false if build plan does not exist", func() {
			_, ok, err := command.NewInvoker(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(ok).To(gomega.BeFalse())
		})

		it("contributes invoker to launch", func() {
			f.AddDependency(command.Dependency, filepath.Join("testdata", "stub-invoker.tgz"))
			f.AddPlan(buildpackplan.Plan{Name: command.Dependency})

			i, _, err := command.NewInvoker(f.Build)
			g.Expect(err).NotTo(gomega.HaveOccurred())

			g.Expect(i.Contribute()).To(gomega.Succeed())

			layer := f.Build.Layers.Layer(command.Dependency)
			g.Expect(layer).To(test.HaveLayerMetadata(false, false, true))
			g.Expect(filepath.Join(layer.Root, "bin", "fixture-marker")).To(gomega.BeARegularFile())

			command := "command-function-invoker"
			g.Expect(f.Build.Layers).To(test.HaveApplicationMetadata(layers.Metadata{
				Processes: []layers.Process{
					{Type: "function", Command: command, Direct: false},
					{Type: "web", Command: command, Direct: false},
				},
			}))
		})
	}, spec.Report(report.Terminal{}))
}
