/*
Copyright 2020 The Kubernetes Authors.

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

package v1alpha4

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta2"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta2"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta2"
	bootstrapv1alpha4 "sigs.k8s.io/cluster-api/internal/apis/bootstrap/kubeadm/v1alpha4"
	clusterv1alpha4 "sigs.k8s.io/cluster-api/internal/apis/core/v1alpha4"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
)

func (src *KubeadmControlPlane) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*controlplanev1.KubeadmControlPlane)

	if err := Convert_v1alpha4_KubeadmControlPlane_To_v1beta2_KubeadmControlPlane(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &controlplanev1.KubeadmControlPlane{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.MachineTemplate.ReadinessGates = restored.Spec.MachineTemplate.ReadinessGates
	dst.Spec.MachineTemplate.NodeDeletionTimeout = restored.Spec.MachineTemplate.NodeDeletionTimeout
	dst.Spec.MachineTemplate.NodeVolumeDetachTimeout = restored.Spec.MachineTemplate.NodeVolumeDetachTimeout
	dst.Spec.RolloutBefore = restored.Spec.RolloutBefore

	if restored.Spec.RemediationStrategy != nil {
		dst.Spec.RemediationStrategy = restored.Spec.RemediationStrategy
	}
	if restored.Status.LastRemediation != nil {
		dst.Status.LastRemediation = restored.Status.LastRemediation
	}

	if restored.Spec.MachineNamingStrategy != nil {
		dst.Spec.MachineNamingStrategy = restored.Spec.MachineNamingStrategy
	}

	bootstrapv1alpha4.MergeRestoredKubeadmConfigSpec(&dst.Spec.KubeadmConfigSpec, &restored.Spec.KubeadmConfigSpec)
	dst.Status.V1Beta2 = restored.Status.V1Beta2

	return nil
}

func (dst *KubeadmControlPlane) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*controlplanev1.KubeadmControlPlane)

	if err := Convert_v1beta2_KubeadmControlPlane_To_v1alpha4_KubeadmControlPlane(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	return utilconversion.MarshalData(src, dst)
}

func (src *KubeadmControlPlaneTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*controlplanev1.KubeadmControlPlaneTemplate)

	if err := Convert_v1alpha4_KubeadmControlPlaneTemplate_To_v1beta2_KubeadmControlPlaneTemplate(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &controlplanev1.KubeadmControlPlaneTemplate{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.Template.Spec.MachineTemplate = restored.Spec.Template.Spec.MachineTemplate

	dst.Spec.Template.ObjectMeta = restored.Spec.Template.ObjectMeta
	if restored.Spec.Template.Spec.MachineTemplate != nil {
		dst.Spec.Template.Spec.MachineTemplate.ObjectMeta = restored.Spec.Template.Spec.MachineTemplate.ObjectMeta
	}

	if dst.Spec.Template.Spec.MachineTemplate == nil {
		dst.Spec.Template.Spec.MachineTemplate = restored.Spec.Template.Spec.MachineTemplate
	} else if restored.Spec.Template.Spec.MachineTemplate != nil {
		dst.Spec.Template.Spec.MachineTemplate.NodeDeletionTimeout = restored.Spec.Template.Spec.MachineTemplate.NodeDeletionTimeout
		dst.Spec.Template.Spec.MachineTemplate.NodeVolumeDetachTimeout = restored.Spec.Template.Spec.MachineTemplate.NodeVolumeDetachTimeout
	}

	dst.Spec.Template.Spec.RolloutBefore = restored.Spec.Template.Spec.RolloutBefore

	if restored.Spec.Template.Spec.RemediationStrategy != nil {
		dst.Spec.Template.Spec.RemediationStrategy = restored.Spec.Template.Spec.RemediationStrategy
	}

	if restored.Spec.Template.Spec.MachineNamingStrategy != nil {
		dst.Spec.Template.Spec.MachineNamingStrategy = restored.Spec.Template.Spec.MachineNamingStrategy
	}

	bootstrapv1alpha4.MergeRestoredKubeadmConfigSpec(&dst.Spec.Template.Spec.KubeadmConfigSpec, &restored.Spec.Template.Spec.KubeadmConfigSpec)

	return nil
}

func (dst *KubeadmControlPlaneTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*controlplanev1.KubeadmControlPlaneTemplate)

	if err := Convert_v1beta2_KubeadmControlPlaneTemplate_To_v1alpha4_KubeadmControlPlaneTemplate(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata
	return utilconversion.MarshalData(src, dst)
}

func Convert_v1alpha4_KubeadmControlPlaneSpec_To_v1beta2_KubeadmControlPlaneTemplateResourceSpec(in *KubeadmControlPlaneSpec, out *controlplanev1.KubeadmControlPlaneTemplateResourceSpec, s apiconversion.Scope) error {
	out.MachineTemplate = &controlplanev1.KubeadmControlPlaneTemplateMachineTemplate{
		NodeDrainTimeout: in.MachineTemplate.NodeDrainTimeout,
	}

	if err := bootstrapv1alpha4.Convert_v1alpha4_KubeadmConfigSpec_To_v1beta2_KubeadmConfigSpec(&in.KubeadmConfigSpec, &out.KubeadmConfigSpec, s); err != nil {
		return err
	}

	out.RolloutAfter = in.RolloutAfter

	if in.RolloutStrategy != nil {
		out.RolloutStrategy = &controlplanev1.RolloutStrategy{}
		if len(in.RolloutStrategy.Type) > 0 {
			out.RolloutStrategy.Type = controlplanev1.RolloutStrategyType(in.RolloutStrategy.Type)
		}
		if in.RolloutStrategy.RollingUpdate != nil {
			out.RolloutStrategy.RollingUpdate = &controlplanev1.RollingUpdate{}

			if in.RolloutStrategy.RollingUpdate.MaxSurge != nil {
				out.RolloutStrategy.RollingUpdate.MaxSurge = in.RolloutStrategy.RollingUpdate.MaxSurge
			}
		}
	}

	return nil
}

func Convert_v1beta2_KubeadmControlPlaneTemplateResourceSpec_To_v1alpha4_KubeadmControlPlaneSpec(in *controlplanev1.KubeadmControlPlaneTemplateResourceSpec, out *KubeadmControlPlaneSpec, s apiconversion.Scope) error {
	if in.MachineTemplate != nil {
		out.MachineTemplate.NodeDrainTimeout = in.MachineTemplate.NodeDrainTimeout
	}

	if err := bootstrapv1alpha4.Convert_v1beta2_KubeadmConfigSpec_To_v1alpha4_KubeadmConfigSpec(&in.KubeadmConfigSpec, &out.KubeadmConfigSpec, s); err != nil {
		return err
	}

	out.RolloutAfter = in.RolloutAfter

	if in.RolloutStrategy != nil {
		out.RolloutStrategy = &RolloutStrategy{}
		if len(in.RolloutStrategy.Type) > 0 {
			out.RolloutStrategy.Type = RolloutStrategyType(in.RolloutStrategy.Type)
		}
		if in.RolloutStrategy.RollingUpdate != nil {
			out.RolloutStrategy.RollingUpdate = &RollingUpdate{}

			if in.RolloutStrategy.RollingUpdate.MaxSurge != nil {
				out.RolloutStrategy.RollingUpdate.MaxSurge = in.RolloutStrategy.RollingUpdate.MaxSurge
			}
		}
	}

	return nil
}

func Convert_v1beta2_KubeadmControlPlaneMachineTemplate_To_v1alpha4_KubeadmControlPlaneMachineTemplate(in *controlplanev1.KubeadmControlPlaneMachineTemplate, out *KubeadmControlPlaneMachineTemplate, s apiconversion.Scope) error {
	// .NodeDrainTimeout was added in v1beta1.
	return autoConvert_v1beta2_KubeadmControlPlaneMachineTemplate_To_v1alpha4_KubeadmControlPlaneMachineTemplate(in, out, s)
}

func Convert_v1beta2_KubeadmControlPlaneSpec_To_v1alpha4_KubeadmControlPlaneSpec(in *controlplanev1.KubeadmControlPlaneSpec, out *KubeadmControlPlaneSpec, scope apiconversion.Scope) error {
	// .RolloutBefore was added in v1beta1.
	// .RemediationStrategy was added in v1beta1.
	return autoConvert_v1beta2_KubeadmControlPlaneSpec_To_v1alpha4_KubeadmControlPlaneSpec(in, out, scope)
}

func Convert_v1beta2_KubeadmControlPlaneStatus_To_v1alpha4_KubeadmControlPlaneStatus(in *controlplanev1.KubeadmControlPlaneStatus, out *KubeadmControlPlaneStatus, scope apiconversion.Scope) error {
	// .LastRemediation was added in v1beta1.
	// .V1Beta2 was added in v1beta1.
	return autoConvert_v1beta2_KubeadmControlPlaneStatus_To_v1alpha4_KubeadmControlPlaneStatus(in, out, scope)
}

func Convert_v1beta2_KubeadmControlPlaneTemplateResource_To_v1alpha4_KubeadmControlPlaneTemplateResource(in *controlplanev1.KubeadmControlPlaneTemplateResource, out *KubeadmControlPlaneTemplateResource, scope apiconversion.Scope) error {
	// .metadata and .spec.machineTemplate.metadata was added in v1beta1.
	return autoConvert_v1beta2_KubeadmControlPlaneTemplateResource_To_v1alpha4_KubeadmControlPlaneTemplateResource(in, out, scope)
}

func Convert_v1beta2_KubeadmConfigSpec_To_v1alpha4_KubeadmConfigSpec(in *bootstrapv1.KubeadmConfigSpec, out *bootstrapv1alpha4.KubeadmConfigSpec, s apiconversion.Scope) error {
	return bootstrapv1alpha4.Convert_v1beta2_KubeadmConfigSpec_To_v1alpha4_KubeadmConfigSpec(in, out, s)
}

func Convert_v1alpha4_KubeadmConfigSpec_To_v1beta2_KubeadmConfigSpec(in *bootstrapv1alpha4.KubeadmConfigSpec, out *bootstrapv1.KubeadmConfigSpec, s apiconversion.Scope) error {
	return bootstrapv1alpha4.Convert_v1alpha4_KubeadmConfigSpec_To_v1beta2_KubeadmConfigSpec(in, out, s)
}

func Convert_v1alpha4_ObjectMeta_To_v1beta2_ObjectMeta(in *clusterv1alpha4.ObjectMeta, out *clusterv1.ObjectMeta, s apiconversion.Scope) error {
	return clusterv1alpha4.Convert_v1alpha4_ObjectMeta_To_v1beta2_ObjectMeta(in, out, s)
}

func Convert_v1beta2_ObjectMeta_To_v1alpha4_ObjectMeta(in *clusterv1.ObjectMeta, out *clusterv1alpha4.ObjectMeta, s apiconversion.Scope) error {
	return clusterv1alpha4.Convert_v1beta2_ObjectMeta_To_v1alpha4_ObjectMeta(in, out, s)
}
