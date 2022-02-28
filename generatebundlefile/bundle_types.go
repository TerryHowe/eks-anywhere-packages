package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/aws/eks-anywhere-packages/api/v1alpha1"
)

// +kubebuilder:object:generate=false
// Same as Bundle except stripped down for generation of yaml file during generate bundleconfig
type BundleGenerate struct {
	// TypeMeta   metav1.TypeMeta `json:",inline"`
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec api.PackageBundleSpec `json:"spec,omitempty"`
}

// Types for input file format

// +kubebuilder:object:root=true
// Input is the schema for the Input file
type Input struct {
	Packages          []Org  `json:"packages,omitempty"`
	KubernetesVersion string `json:"kubernetesVersion,omitempty"`
}

// Projects object containing the input file github org and repo locations
type Org struct {
	Name     string    `json:"name,omitempty"`
	Projects []Project `json:"projects,omitempty"`
}

// Repos is the object containing the project within the github org, and the release tag
type Project struct {
	Name       string `json:"name,omitempty"`
	Registry   string `json:"registry,omitempty"`
	Repository string `json:"repository,omitempty"`
	Versions   []Tag  `json:"versions,omitempty"`
}

// Tag is the release tag
type Tag struct {
	Name string `json:"name,omitempty"`
}

// Matches returns a list of inputs which align with ECR tags that exist
func (project Project) Matches(tag string) []string {
	matchlist := []string{}
	for _, version := range project.Versions {
		if version.Name == tag {
			matchlist = append(matchlist, version.Name)
		}
	}
	return matchlist
}