//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by conversion-gen. DO NOT EDIT.

package v1beta1

import (
	unsafe "unsafe"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
	apiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
	apiv1beta2 "sigs.k8s.io/cluster-api/api/v1beta2"
	v1beta2 "sigs.k8s.io/cluster-api/exp/api/v1beta2"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*MachinePool)(nil), (*v1beta2.MachinePool)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_MachinePool_To_v1beta2_MachinePool(a.(*MachinePool), b.(*v1beta2.MachinePool), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.MachinePool)(nil), (*MachinePool)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_MachinePool_To_v1beta1_MachinePool(a.(*v1beta2.MachinePool), b.(*MachinePool), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachinePoolList)(nil), (*v1beta2.MachinePoolList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_MachinePoolList_To_v1beta2_MachinePoolList(a.(*MachinePoolList), b.(*v1beta2.MachinePoolList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.MachinePoolList)(nil), (*MachinePoolList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_MachinePoolList_To_v1beta1_MachinePoolList(a.(*v1beta2.MachinePoolList), b.(*MachinePoolList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachinePoolSpec)(nil), (*v1beta2.MachinePoolSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_MachinePoolSpec_To_v1beta2_MachinePoolSpec(a.(*MachinePoolSpec), b.(*v1beta2.MachinePoolSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*v1beta2.MachinePoolSpec)(nil), (*MachinePoolSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_MachinePoolSpec_To_v1beta1_MachinePoolSpec(a.(*v1beta2.MachinePoolSpec), b.(*MachinePoolSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1.Condition)(nil), (*apiv1beta1.Condition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1_Condition_To_v1beta1_Condition(a.(*v1.Condition), b.(*apiv1beta1.Condition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apiv1beta1.Condition)(nil), (*v1.Condition)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Condition_To_v1_Condition(a.(*apiv1beta1.Condition), b.(*v1.Condition), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*MachinePoolStatus)(nil), (*v1beta2.MachinePoolStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_MachinePoolStatus_To_v1beta2_MachinePoolStatus(a.(*MachinePoolStatus), b.(*v1beta2.MachinePoolStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apiv1beta1.MachineTemplateSpec)(nil), (*apiv1beta2.MachineTemplateSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_MachineTemplateSpec_To_v1beta2_MachineTemplateSpec(a.(*apiv1beta1.MachineTemplateSpec), b.(*apiv1beta2.MachineTemplateSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*v1beta2.MachinePoolStatus)(nil), (*MachinePoolStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_MachinePoolStatus_To_v1beta1_MachinePoolStatus(a.(*v1beta2.MachinePoolStatus), b.(*MachinePoolStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*apiv1beta2.MachineTemplateSpec)(nil), (*apiv1beta1.MachineTemplateSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta2_MachineTemplateSpec_To_v1beta1_MachineTemplateSpec(a.(*apiv1beta2.MachineTemplateSpec), b.(*apiv1beta1.MachineTemplateSpec), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1beta1_MachinePool_To_v1beta2_MachinePool(in *MachinePool, out *v1beta2.MachinePool, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1beta1_MachinePoolSpec_To_v1beta2_MachinePoolSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_MachinePoolStatus_To_v1beta2_MachinePoolStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_MachinePool_To_v1beta2_MachinePool is an autogenerated conversion function.
func Convert_v1beta1_MachinePool_To_v1beta2_MachinePool(in *MachinePool, out *v1beta2.MachinePool, s conversion.Scope) error {
	return autoConvert_v1beta1_MachinePool_To_v1beta2_MachinePool(in, out, s)
}

func autoConvert_v1beta2_MachinePool_To_v1beta1_MachinePool(in *v1beta2.MachinePool, out *MachinePool, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1beta2_MachinePoolSpec_To_v1beta1_MachinePoolSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1beta2_MachinePoolStatus_To_v1beta1_MachinePoolStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta2_MachinePool_To_v1beta1_MachinePool is an autogenerated conversion function.
func Convert_v1beta2_MachinePool_To_v1beta1_MachinePool(in *v1beta2.MachinePool, out *MachinePool, s conversion.Scope) error {
	return autoConvert_v1beta2_MachinePool_To_v1beta1_MachinePool(in, out, s)
}

func autoConvert_v1beta1_MachinePoolList_To_v1beta2_MachinePoolList(in *MachinePoolList, out *v1beta2.MachinePoolList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]v1beta2.MachinePool, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_MachinePool_To_v1beta2_MachinePool(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1beta1_MachinePoolList_To_v1beta2_MachinePoolList is an autogenerated conversion function.
func Convert_v1beta1_MachinePoolList_To_v1beta2_MachinePoolList(in *MachinePoolList, out *v1beta2.MachinePoolList, s conversion.Scope) error {
	return autoConvert_v1beta1_MachinePoolList_To_v1beta2_MachinePoolList(in, out, s)
}

func autoConvert_v1beta2_MachinePoolList_To_v1beta1_MachinePoolList(in *v1beta2.MachinePoolList, out *MachinePoolList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MachinePool, len(*in))
		for i := range *in {
			if err := Convert_v1beta2_MachinePool_To_v1beta1_MachinePool(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1beta2_MachinePoolList_To_v1beta1_MachinePoolList is an autogenerated conversion function.
func Convert_v1beta2_MachinePoolList_To_v1beta1_MachinePoolList(in *v1beta2.MachinePoolList, out *MachinePoolList, s conversion.Scope) error {
	return autoConvert_v1beta2_MachinePoolList_To_v1beta1_MachinePoolList(in, out, s)
}

func autoConvert_v1beta1_MachinePoolSpec_To_v1beta2_MachinePoolSpec(in *MachinePoolSpec, out *v1beta2.MachinePoolSpec, s conversion.Scope) error {
	out.ClusterName = in.ClusterName
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	if err := Convert_v1beta1_MachineTemplateSpec_To_v1beta2_MachineTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	out.MinReadySeconds = (*int32)(unsafe.Pointer(in.MinReadySeconds))
	out.ProviderIDList = *(*[]string)(unsafe.Pointer(&in.ProviderIDList))
	out.FailureDomains = *(*[]string)(unsafe.Pointer(&in.FailureDomains))
	return nil
}

// Convert_v1beta1_MachinePoolSpec_To_v1beta2_MachinePoolSpec is an autogenerated conversion function.
func Convert_v1beta1_MachinePoolSpec_To_v1beta2_MachinePoolSpec(in *MachinePoolSpec, out *v1beta2.MachinePoolSpec, s conversion.Scope) error {
	return autoConvert_v1beta1_MachinePoolSpec_To_v1beta2_MachinePoolSpec(in, out, s)
}

func autoConvert_v1beta2_MachinePoolSpec_To_v1beta1_MachinePoolSpec(in *v1beta2.MachinePoolSpec, out *MachinePoolSpec, s conversion.Scope) error {
	out.ClusterName = in.ClusterName
	out.Replicas = (*int32)(unsafe.Pointer(in.Replicas))
	if err := Convert_v1beta2_MachineTemplateSpec_To_v1beta1_MachineTemplateSpec(&in.Template, &out.Template, s); err != nil {
		return err
	}
	out.MinReadySeconds = (*int32)(unsafe.Pointer(in.MinReadySeconds))
	out.ProviderIDList = *(*[]string)(unsafe.Pointer(&in.ProviderIDList))
	out.FailureDomains = *(*[]string)(unsafe.Pointer(&in.FailureDomains))
	return nil
}

// Convert_v1beta2_MachinePoolSpec_To_v1beta1_MachinePoolSpec is an autogenerated conversion function.
func Convert_v1beta2_MachinePoolSpec_To_v1beta1_MachinePoolSpec(in *v1beta2.MachinePoolSpec, out *MachinePoolSpec, s conversion.Scope) error {
	return autoConvert_v1beta2_MachinePoolSpec_To_v1beta1_MachinePoolSpec(in, out, s)
}

func autoConvert_v1beta1_MachinePoolStatus_To_v1beta2_MachinePoolStatus(in *MachinePoolStatus, out *v1beta2.MachinePoolStatus, s conversion.Scope) error {
	out.NodeRefs = *(*[]corev1.ObjectReference)(unsafe.Pointer(&in.NodeRefs))
	out.Replicas = in.Replicas
	if err := v1.Convert_int32_To_Pointer_int32(&in.ReadyReplicas, &out.ReadyReplicas, s); err != nil {
		return err
	}
	if err := v1.Convert_int32_To_Pointer_int32(&in.AvailableReplicas, &out.AvailableReplicas, s); err != nil {
		return err
	}
	// WARNING: in.UnavailableReplicas requires manual conversion: does not exist in peer-type
	// WARNING: in.FailureReason requires manual conversion: does not exist in peer-type
	// WARNING: in.FailureMessage requires manual conversion: does not exist in peer-type
	out.Phase = in.Phase
	// WARNING: in.BootstrapReady requires manual conversion: does not exist in peer-type
	// WARNING: in.InfrastructureReady requires manual conversion: does not exist in peer-type
	out.ObservedGeneration = in.ObservedGeneration
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]v1.Condition, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_Condition_To_v1_Condition(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Conditions = nil
	}
	// WARNING: in.V1Beta2 requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_v1beta2_MachinePoolStatus_To_v1beta1_MachinePoolStatus(in *v1beta2.MachinePoolStatus, out *MachinePoolStatus, s conversion.Scope) error {
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make(apiv1beta1.Conditions, len(*in))
		for i := range *in {
			if err := Convert_v1_Condition_To_v1beta1_Condition(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Conditions = nil
	}
	// WARNING: in.Initialization requires manual conversion: does not exist in peer-type
	out.NodeRefs = *(*[]corev1.ObjectReference)(unsafe.Pointer(&in.NodeRefs))
	out.Replicas = in.Replicas
	if err := v1.Convert_Pointer_int32_To_int32(&in.ReadyReplicas, &out.ReadyReplicas, s); err != nil {
		return err
	}
	if err := v1.Convert_Pointer_int32_To_int32(&in.AvailableReplicas, &out.AvailableReplicas, s); err != nil {
		return err
	}
	// WARNING: in.UpToDateReplicas requires manual conversion: does not exist in peer-type
	out.Phase = in.Phase
	out.ObservedGeneration = in.ObservedGeneration
	// WARNING: in.Deprecated requires manual conversion: does not exist in peer-type
	return nil
}
