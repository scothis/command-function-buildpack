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
	"path/filepath"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/buildpack/libbuildpack/detect"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/projectriff/command-function-buildpack/command"
	"github.com/projectriff/libfnbuildpack/function"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestBuildpack(t *testing.T) {
	spec.Run(t, "Buildpack", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var (
			b command.Buildpack
			f *test.DetectFactory
		)

		it.Before(func() {
			b = command.Buildpack{}
			f = test.NewDetectFactory(t)
		})

		when("id", func() {

			it("returns id", func() {
				g.Expect(b.Id()).To(gomega.Equal("command"))
			})
		})

		when("detect", func() {

			it("fails with no artifact", func() {
				g.Expect(b.Detect(f.Detect, function.Metadata{})).To(gomega.Equal(detect.FailStatusCode))
			})

			it("fails with non-existent file", func() {
				g.Expect(b.Detect(f.Detect, function.Metadata{Artifact: "test-file"})).To(gomega.Equal(function.Error_ComponentInternal))
			})

			it("errors with non-executable file", func() {
				test.TouchFile(t, f.Detect.Application.Root, "test-file")

				g.Expect(b.Detect(f.Detect, function.Metadata{Artifact: "test-file"})).To(gomega.Equal(function.Error_ComponentInternal))
			})

			it("passes with executable file", func() {
				test.WriteFileWithPerm(t, filepath.Join(f.Detect.Application.Root, "test-file"), 0700, "")

				g.Expect(b.Detect(f.Detect, function.Metadata{Artifact: "test-file"})).To(gomega.Equal(detect.PassStatusCode))
				g.Expect(f.Plans).To(test.HavePlans(buildplan.Plan{
					Provides: []buildplan.Provided{
						{Name: command.Dependency},
					},
					Requires: []buildplan.Required{
						{
							Name: command.Dependency,
							Metadata: map[string]interface{}{
								command.Command: "test-file",
							},
						},
					},
				}))
			})
		})
	}, spec.Report(report.Terminal{}))
}
