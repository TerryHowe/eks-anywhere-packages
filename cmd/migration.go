// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
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

package cmd

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/aws/eks-anywhere-packages/api/v1alpha1"
	"github.com/spf13/cobra"
)

type migrateContext struct {
}

var migrateCommandContext = &migrateContext{}

func init() {
	rootCmd.AddCommand(migrateCommand)
}

func getKubeConfig() string {
	kubeconfig := os.Getenv("KUBECONFIG")
	if len(kubeconfig) == 0 {
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	}
	return kubeconfig
}

func migrate() error {
	ctx := context.Background()

	clusterName := os.Getenv("CLUSTER_NAME")
	if len(clusterName) == 0 {
		return fmt.Errorf("CLUSTER_NAME environment variable is not set or exported")
	}
	clusterNamespace := "eksa-packages-" + clusterName

	config, err := rest.InClusterConfig()
	if err != nil {
		config, err = clientcmd.BuildConfigFromFlags("", getKubeConfig())
		if err != nil {
			return fmt.Errorf("building Kubernetes configuration: %v", err)
		}
	}

	var scheme = runtime.NewScheme()
	err = v1alpha1.AddToScheme(scheme)
	if err != nil {
		return fmt.Errorf("add schema: %v", err)
	}

	clientOptions := client.Options{
		Scheme: scheme,
	}
	rtClient, err := client.New(config, clientOptions)
	if err != nil {
		return fmt.Errorf("creating Kubernetes runtime client: %v", err)
	}

	packageList := v1alpha1.PackageList{}
	err = rtClient.List(ctx, &packageList, &client.ListOptions{Namespace: "eksa-packages"})
	if err != nil {
		return fmt.Errorf("reading: %v", err)
	}

	for _, pkg := range packageList.Items {
		newPackage := v1alpha1.Package{}
		err = rtClient.Get(ctx, client.ObjectKey{Namespace: clusterNamespace, Name: pkg.Name}, &newPackage)
		if err != nil {
			pkg.Namespace = clusterNamespace
			pkg.ResourceVersion = ""
			pkg.UID = ""
			err = rtClient.Create(ctx, &pkg, &client.CreateOptions{})
			if err != nil {
				return fmt.Errorf("update error: %v", err)
			}
		} else {
			pkg.Namespace = clusterNamespace
			pkg.ResourceVersion = newPackage.ResourceVersion
			pkg.UID = newPackage.UID
			err = rtClient.Update(ctx, &pkg, &client.UpdateOptions{})
			if err != nil {
				return fmt.Errorf("update error: %v", err)
			}
		}
	}
	return nil
}

func runMigration(_ *cobra.Command, _ []string) {
	err := migrate()
	if err != nil {
		packageLog.Error(err, "migration")
	}
}

var migrateCommand = &cobra.Command{
	Use:   "migrate",
	Short: "Run package migrations",
	Long:  "Run package migrations",
	Run:   runMigration,
}
