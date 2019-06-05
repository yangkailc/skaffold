/*
Copyright 2019 The Skaffold Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package kaniko

import (
	"context"
	"io"
	"time"

	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/constants"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/docker"
	runcontext "github.com/GoogleContainerTools/skaffold/pkg/skaffold/runner/context"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/schema/latest"
	"github.com/GoogleContainerTools/skaffold/pkg/skaffold/util"
	"github.com/pkg/errors"
)

// Builder builds docker artifacts on Kubernetes, using Kaniko.
type Builder struct {
	*latest.ClusterDetails

	timeout            time.Duration
	insecureRegistries map[string]bool
}

// NewBuilder creates a new Builder that builds artifacts with Kaniko.
func NewBuilder(runCtx *runcontext.RunContext) (*Builder, error) {
	timeout, err := time.ParseDuration(runCtx.Cfg.Build.Cluster.Timeout)
	if err != nil {
		return nil, errors.Wrap(err, "parsing timeout")
	}

	return &Builder{
		ClusterDetails: runCtx.Cfg.Build.Cluster,
		timeout:        timeout,
                insecureRegistries: runCtx.InsecureRegistries,
	}, nil
}

// Labels are labels specific to Kaniko builder.
func (b *Builder) Labels() map[string]string {
	return map[string]string{
		constants.Labels.Builder: "kaniko",
	}
}

// DependenciesForArtifact returns the Dockerfile dependencies for this kaniko artifact
func (b *Builder) DependenciesForArtifact(ctx context.Context, a *latest.Artifact) ([]string, error) {
	paths, err := docker.GetDependencies(ctx, a.Workspace, a.KanikoArtifact.DockerfilePath, a.KanikoArtifact.BuildArgs, b.insecureRegistries)
	if err != nil {
		return nil, errors.Wrapf(err, "getting dependencies for %s", a.ImageName)
	}
	return util.AbsolutePaths(a.Workspace, paths), nil
}

func (b *Builder) Prune(ctx context.Context, out io.Writer) error {
	return nil
}
