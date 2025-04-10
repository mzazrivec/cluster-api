/*
Copyright 2019 The Kubernetes Authors.

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

package v1alpha3

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta2"
	"sigs.k8s.io/cluster-api/util/conditions"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func (src *Cluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*clusterv1.Cluster)

	if err := Convert_v1alpha3_Cluster_To_v1beta2_Cluster(src, dst, nil); err != nil {
		return err
	}

	// Given this is a bool and there is no timestamp associated with it, when this condition is set, its timestamp
	// will be "now". See https://github.com/kubernetes-sigs/cluster-api/issues/3798#issuecomment-708619826 for more
	// discussion.
	if src.Status.ControlPlaneInitialized {
		conditions.MarkTrue(dst, clusterv1.ControlPlaneInitializedCondition)
	}

	// Manually restore data.
	restored := &clusterv1.Cluster{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.AvailabilityGates = restored.Spec.AvailabilityGates
	if restored.Spec.Topology != nil {
		dst.Spec.Topology = restored.Spec.Topology
	}
	dst.Status.V1Beta2 = restored.Status.V1Beta2

	return nil
}

func (dst *Cluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*clusterv1.Cluster)

	if err := Convert_v1beta2_Cluster_To_v1alpha3_Cluster(src, dst, nil); err != nil {
		return err
	}

	// Set the v1alpha3 boolean status field if the v1alpha4 condition was true
	if conditions.IsTrue(src, clusterv1.ControlPlaneInitializedCondition) {
		dst.Status.ControlPlaneInitialized = true
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

func (src *Machine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*clusterv1.Machine)

	if err := Convert_v1alpha3_Machine_To_v1beta2_Machine(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &clusterv1.Machine{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.ReadinessGates = restored.Spec.ReadinessGates
	dst.Spec.NodeDeletionTimeout = restored.Spec.NodeDeletionTimeout
	dst.Spec.NodeVolumeDetachTimeout = restored.Spec.NodeVolumeDetachTimeout
	dst.Status.NodeInfo = restored.Status.NodeInfo
	dst.Status.CertificatesExpiryDate = restored.Status.CertificatesExpiryDate
	dst.Status.Deletion = restored.Status.Deletion
	dst.Status.V1Beta2 = restored.Status.V1Beta2

	return nil
}

func (dst *Machine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*clusterv1.Machine)

	if err := Convert_v1beta2_Machine_To_v1alpha3_Machine(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

func (src *MachineSet) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*clusterv1.MachineSet)

	if err := Convert_v1alpha3_MachineSet_To_v1beta2_MachineSet(src, dst, nil); err != nil {
		return err
	}
	// Manually restore data.
	restored := &clusterv1.MachineSet{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}
	dst.Spec.Template.Spec.ReadinessGates = restored.Spec.Template.Spec.ReadinessGates
	dst.Spec.Template.Spec.NodeDeletionTimeout = restored.Spec.Template.Spec.NodeDeletionTimeout
	dst.Spec.Template.Spec.NodeVolumeDetachTimeout = restored.Spec.Template.Spec.NodeVolumeDetachTimeout
	dst.Status.Conditions = restored.Status.Conditions
	dst.Status.V1Beta2 = restored.Status.V1Beta2

	if restored.Spec.MachineNamingStrategy != nil {
		dst.Spec.MachineNamingStrategy = restored.Spec.MachineNamingStrategy
	}
	return nil
}

func (dst *MachineSet) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*clusterv1.MachineSet)

	if err := Convert_v1beta2_MachineSet_To_v1alpha3_MachineSet(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}
	return nil
}

func (src *MachineDeployment) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*clusterv1.MachineDeployment)

	if err := Convert_v1alpha3_MachineDeployment_To_v1beta2_MachineDeployment(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &clusterv1.MachineDeployment{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	if restored.Spec.Strategy != nil {
		if dst.Spec.Strategy == nil {
			dst.Spec.Strategy = &clusterv1.MachineDeploymentStrategy{}
		}
		if restored.Spec.Strategy.RollingUpdate != nil {
			if dst.Spec.Strategy.RollingUpdate == nil {
				dst.Spec.Strategy.RollingUpdate = &clusterv1.MachineRollingUpdateDeployment{}
			}
			dst.Spec.Strategy.RollingUpdate.DeletePolicy = restored.Spec.Strategy.RollingUpdate.DeletePolicy
		}
		dst.Spec.Strategy.Remediation = restored.Spec.Strategy.Remediation
	}

	if restored.Spec.MachineNamingStrategy != nil {
		dst.Spec.MachineNamingStrategy = restored.Spec.MachineNamingStrategy
	}

	dst.Spec.Template.Spec.ReadinessGates = restored.Spec.Template.Spec.ReadinessGates
	dst.Spec.Template.Spec.NodeDeletionTimeout = restored.Spec.Template.Spec.NodeDeletionTimeout
	dst.Spec.Template.Spec.NodeVolumeDetachTimeout = restored.Spec.Template.Spec.NodeVolumeDetachTimeout
	dst.Spec.RolloutAfter = restored.Spec.RolloutAfter
	dst.Status.Conditions = restored.Status.Conditions
	dst.Status.V1Beta2 = restored.Status.V1Beta2

	return nil
}

func (dst *MachineDeployment) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*clusterv1.MachineDeployment)

	if err := Convert_v1beta2_MachineDeployment_To_v1alpha3_MachineDeployment(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

func (src *MachineHealthCheck) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*clusterv1.MachineHealthCheck)

	if err := Convert_v1alpha3_MachineHealthCheck_To_v1beta2_MachineHealthCheck(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &clusterv1.MachineHealthCheck{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	if restored.Spec.UnhealthyRange != nil {
		dst.Spec.UnhealthyRange = restored.Spec.UnhealthyRange
	}
	dst.Status.V1Beta2 = restored.Status.V1Beta2

	return nil
}

func (dst *MachineHealthCheck) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*clusterv1.MachineHealthCheck)

	if err := Convert_v1beta2_MachineHealthCheck_To_v1alpha3_MachineHealthCheck(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}

	return nil
}

func Convert_v1beta2_MachineSetStatus_To_v1alpha3_MachineSetStatus(in *clusterv1.MachineSetStatus, out *MachineSetStatus, _ apiconversion.Scope) error {
	// Status.Conditions was introduced in v1alpha4, thus requiring a custom conversion function; the values is going to be preserved in an annotation thus allowing roundtrip without loosing informations
	// V1Beta2 was added in v1beta1.
	return autoConvert_v1beta2_MachineSetStatus_To_v1alpha3_MachineSetStatus(in, out, nil)
}

func Convert_v1beta2_ClusterSpec_To_v1alpha3_ClusterSpec(in *clusterv1.ClusterSpec, out *ClusterSpec, s apiconversion.Scope) error {
	// NOTE: custom conversion func is required because spec.Topology does not exist in v1alpha3
	// AvailabilityGates was added in v1beta1.
	return autoConvert_v1beta2_ClusterSpec_To_v1alpha3_ClusterSpec(in, out, s)
}

func Convert_v1beta2_ClusterStatus_To_v1alpha3_ClusterStatus(in *clusterv1.ClusterStatus, out *ClusterStatus, s apiconversion.Scope) error {
	// V1Beta2 was added in v1beta1.
	return autoConvert_v1beta2_ClusterStatus_To_v1alpha3_ClusterStatus(in, out, s)
}

func Convert_v1alpha3_Bootstrap_To_v1beta2_Bootstrap(in *Bootstrap, out *clusterv1.Bootstrap, s apiconversion.Scope) error {
	return autoConvert_v1alpha3_Bootstrap_To_v1beta2_Bootstrap(in, out, s)
}

func Convert_v1beta2_MachineRollingUpdateDeployment_To_v1alpha3_MachineRollingUpdateDeployment(in *clusterv1.MachineRollingUpdateDeployment, out *MachineRollingUpdateDeployment, s apiconversion.Scope) error {
	return autoConvert_v1beta2_MachineRollingUpdateDeployment_To_v1alpha3_MachineRollingUpdateDeployment(in, out, s)
}

func Convert_v1beta2_MachineHealthCheckSpec_To_v1alpha3_MachineHealthCheckSpec(in *clusterv1.MachineHealthCheckSpec, out *MachineHealthCheckSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta2_MachineHealthCheckSpec_To_v1alpha3_MachineHealthCheckSpec(in, out, s)
}

func Convert_v1beta2_MachineHealthCheckStatus_To_v1alpha3_MachineHealthCheckStatus(in *clusterv1.MachineHealthCheckStatus, out *MachineHealthCheckStatus, s apiconversion.Scope) error {
	// V1Beta2 was added in v1beta1.
	return autoConvert_v1beta2_MachineHealthCheckStatus_To_v1alpha3_MachineHealthCheckStatus(in, out, s)
}

func Convert_v1alpha3_ClusterStatus_To_v1beta2_ClusterStatus(in *ClusterStatus, out *clusterv1.ClusterStatus, s apiconversion.Scope) error {
	return autoConvert_v1alpha3_ClusterStatus_To_v1beta2_ClusterStatus(in, out, s)
}

func Convert_v1alpha3_ObjectMeta_To_v1beta2_ObjectMeta(in *ObjectMeta, out *clusterv1.ObjectMeta, s apiconversion.Scope) error {
	return autoConvert_v1alpha3_ObjectMeta_To_v1beta2_ObjectMeta(in, out, s)
}

func Convert_v1beta2_MachineStatus_To_v1alpha3_MachineStatus(in *clusterv1.MachineStatus, out *MachineStatus, s apiconversion.Scope) error {
	// V1Beta2 was added in v1beta1.
	return autoConvert_v1beta2_MachineStatus_To_v1alpha3_MachineStatus(in, out, s)
}

func Convert_v1beta2_MachineSpec_To_v1alpha3_MachineSpec(in *clusterv1.MachineSpec, out *MachineSpec, s apiconversion.Scope) error {
	// spec.nodeDeletionTimeout was added in v1beta1.
	// ReadinessGates was added in v1beta1.
	return autoConvert_v1beta2_MachineSpec_To_v1alpha3_MachineSpec(in, out, s)
}

func Convert_v1beta2_MachineDeploymentSpec_To_v1alpha3_MachineDeploymentSpec(in *clusterv1.MachineDeploymentSpec, out *MachineDeploymentSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta2_MachineDeploymentSpec_To_v1alpha3_MachineDeploymentSpec(in, out, s)
}

func Convert_v1beta2_MachineDeploymentStatus_To_v1alpha3_MachineDeploymentStatus(in *clusterv1.MachineDeploymentStatus, out *MachineDeploymentStatus, s apiconversion.Scope) error {
	// Status.Conditions was introduced in v1alpha4, thus requiring a custom conversion function; the values is going to be preserved in an annotation thus allowing roundtrip without loosing informations
	// V1Beta2 was added in v1beta1.
	return autoConvert_v1beta2_MachineDeploymentStatus_To_v1alpha3_MachineDeploymentStatus(in, out, s)
}

func Convert_v1alpha3_MachineStatus_To_v1beta2_MachineStatus(in *MachineStatus, out *clusterv1.MachineStatus, s apiconversion.Scope) error {
	// Status.version has been removed in v1beta1, thus requiring custom conversion function. the information will be dropped.
	return autoConvert_v1alpha3_MachineStatus_To_v1beta2_MachineStatus(in, out, s)
}

func Convert_v1beta2_MachineDeploymentStrategy_To_v1alpha3_MachineDeploymentStrategy(in *clusterv1.MachineDeploymentStrategy, out *MachineDeploymentStrategy, s apiconversion.Scope) error {
	return autoConvert_v1beta2_MachineDeploymentStrategy_To_v1alpha3_MachineDeploymentStrategy(in, out, s)
}

func Convert_v1beta2_MachineSetSpec_To_v1alpha3_MachineSetSpec(in *clusterv1.MachineSetSpec, out *MachineSetSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta2_MachineSetSpec_To_v1alpha3_MachineSetSpec(in, out, s)
}
