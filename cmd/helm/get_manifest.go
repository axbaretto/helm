/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package main

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"k8s.io/helm/pkg/helm"
)

var getManifestHelp = `
This command fetches the generated manifest for a given release.

A manifest is a YAML-encoded representation of the Kubernetes resources that
were generated from this release's chart(s). If a chart is dependent on other
charts, those resources will also be included in the manifest.
`

type getManifestCmd struct {
	release string
	out     io.Writer
	client  helm.Interface
}

func newGetManifestCmd(client helm.Interface, out io.Writer) *cobra.Command {
	get := &getManifestCmd{
		out:    out,
		client: client,
	}
	cmd := &cobra.Command{
		Use:   "manifest [flags] RELEASE_NAME",
		Short: "download the manifest for a named release",
		Long:  getManifestHelp,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errReleaseRequired
			}
			get.release = args[0]
			if get.client == nil {
				get.client = helm.NewClient(helm.Host(tillerHost))
			}
			return get.run()
		},
	}
	return cmd
}

// getManifest implements 'helm get manifest'
func (g *getManifestCmd) run() error {
	res, err := g.client.ReleaseContent(g.release)
	if err != nil {
		return prettyError(err)
	}
	fmt.Fprintln(g.out, res.Release.Manifest)
	return nil
}
