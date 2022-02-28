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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Version",type=string,JSONPath=`.spec.kubeVersion`
//+kubebuilder:printcolumn:name="State",type=string,JSONPath=`.status.state`
// PackageBundle is the Schema for the packagebundles API
type PackageBundle struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PackageBundleSpec   `json:"spec,omitempty"`
	Status PackageBundleStatus `json:"status,omitempty"`
}

// PackageBundleSpec defines the desired state of PackageBundle
type PackageBundleSpec struct {
	// +kubebuilder:validation:Required
	// KubeVersion is the Kubernetes version to which this bundle is compatible.
	KubeVersion string `json:"kubeVersion,omitempty"`

	// +kubebuilder:validation:Required
	// Packages supported by this bundle.
	Packages []BundlePackage `json:"packages"`
}

// BundlePackage specifies a package within a bundle.
type BundlePackage struct {
	// +kubebuilder:validation:Required
	// Name of the package
	Name string `json:"name,omitempty"`

	// +kubebuilder:validation:Required
	// Source location for the package (probably a helm chart).
	Source BundlePackageSource `json:"source"`
}

// BundlePackageSource identifies the location of a package.
type BundlePackageSource struct {
	// +kubebuilder:validation:Required
	// Registry in which the package is found.
	Registry string `json:"registry"`

	// +kubebuilder:validation:Required
	// Repository within the Registry where the package is found.
	Repository string `json:"repository"`

	// +kubebuilder:validation:MinItems=1
	// Versions of the package supported by this bundle.
	Versions []SourceVersion `json:"versions"`
}

// SourceVersion describes a version of an package within a repository.
type SourceVersion struct {
	// +kubebuilder:validation:Required
	// Name is a human-friendly description of the version, e.g. "v1.0".
	Name string `json:"name"`

	// +kubebuilder:validation:Required
	// Tag is a sha256 checksum value identifying the version of the package and
	// its contents.
	Tag string `json:"tag"`
}

// PackageBundleStatus defines the observed state of PackageBundle
type PackageBundleStatus struct {
	Spec  PackageBundleSpec      `json:"spec,omitempty"`
	State PackageBundleStateEnum `json:"state"`
}

//+kubebuilder:validation:Enum={"inactive","active","active (upgrade available)"}
type PackageBundleStateEnum string

const (
	PackageBundleStateInactive         PackageBundleStateEnum = "inactive"
	PackageBundleStateActive           PackageBundleStateEnum = "active"
	PackageBundleStateUpgradeAvailable PackageBundleStateEnum = "active (upgrade available)"
)

//+kubebuilder:object:root=true
// PackageBundleList contains a list of PackageBundle
type PackageBundleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PackageBundle `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PackageBundle{}, &PackageBundleList{})
}