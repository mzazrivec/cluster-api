/*
Copyright 2022 The Kubernetes Authors.

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

package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta2"
	runtimev1 "sigs.k8s.io/cluster-api/exp/runtime/api/v1alpha1"
	runtimeclient "sigs.k8s.io/cluster-api/exp/runtime/client"
	"sigs.k8s.io/cluster-api/util/conditions"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/conditions/deprecated/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/cluster-api/util/paused"
	"sigs.k8s.io/cluster-api/util/predicates"
)

const (
	// tlsCAKey is used as a data key in Secret resources to store a CA certificate.
	tlsCAKey = "ca.crt"
)

// +kubebuilder:rbac:groups=runtime.cluster.x-k8s.io,resources=extensionconfigs;extensionconfigs/status,verbs=get;list;watch;patch;update
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch

// Reconciler reconciles an ExtensionConfig object.
type Reconciler struct {
	Client        client.Client
	APIReader     client.Reader
	RuntimeClient runtimeclient.Client
	// WatchFilterValue is the label value used to filter events prior to reconciliation.
	WatchFilterValue string
}

func (r *Reconciler) SetupWithManager(ctx context.Context, mgr ctrl.Manager, options controller.Options, partialSecretCache cache.Cache) error {
	if r.Client == nil || r.APIReader == nil || r.RuntimeClient == nil {
		return errors.New("Client, APIReader and RuntimeClient must not be nil")
	}

	predicateLog := ctrl.LoggerFrom(ctx).WithValues("controller", "extensionconfig")
	err := ctrl.NewControllerManagedBy(mgr).
		For(&runtimev1.ExtensionConfig{}).
		WatchesRawSource(source.Kind(
			partialSecretCache,
			&metav1.PartialObjectMetadata{
				TypeMeta: metav1.TypeMeta{
					Kind:       "Secret",
					APIVersion: "v1",
				},
			},
			handler.TypedEnqueueRequestsFromMapFunc(
				r.secretToExtensionConfig,
			),
			predicates.TypedResourceIsChanged[*metav1.PartialObjectMetadata](mgr.GetScheme(), predicateLog),
		)).
		WithOptions(options).
		WithEventFilter(predicates.ResourceHasFilterLabel(mgr.GetScheme(), predicateLog, r.WatchFilterValue)).
		Complete(r)
	if err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	if err := indexByExtensionInjectCAFromSecretName(ctx, mgr); err != nil {
		return errors.Wrap(err, "failed setting up with a controller manager")
	}

	// warmupRunnable will attempt to sync the RuntimeSDK registry with existing ExtensionConfig objects to ensure extensions
	// are discovered before controllers begin reconciling.
	err = mgr.Add(&warmupRunnable{
		Client:        r.Client,
		APIReader:     r.APIReader,
		RuntimeClient: r.RuntimeClient,
	})
	if err != nil {
		return errors.Wrap(err, "failed adding warmupRunnable to controller manager")
	}
	return nil
}

func (r *Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var errs []error
	log := ctrl.LoggerFrom(ctx)

	// Requeue events when the registry is not ready.
	// The registry will become ready after it is 'warmed up' by warmupRunnable.
	if !r.RuntimeClient.IsReady() {
		return ctrl.Result{Requeue: true}, nil
	}

	extensionConfig := &runtimev1.ExtensionConfig{}
	err := r.Client.Get(ctx, req.NamespacedName, extensionConfig)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// ExtensionConfig not found. Remove from registry.
			// First we need to add Namespace/Name to empty ExtensionConfig object.
			extensionConfig.Name = req.Name
			extensionConfig.Namespace = req.Namespace
			return r.reconcileDelete(ctx, extensionConfig)
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Copy to avoid modifying the original extensionConfig.
	original := extensionConfig.DeepCopy()

	if isPaused, requeue, err := paused.EnsurePausedCondition(ctx, r.Client, nil, extensionConfig); err != nil || isPaused || requeue {
		return ctrl.Result{}, err
	}

	// Handle deletion reconciliation loop.
	if !extensionConfig.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, extensionConfig)
	}

	// Inject CABundle from secret if annotation is set. Otherwise https calls may fail.
	if err := reconcileCABundle(ctx, r.Client, extensionConfig); err != nil {
		return ctrl.Result{}, err
	}

	// discoverExtensionConfig will return a discovered ExtensionConfig with the appropriate conditions.
	discoveredExtensionConfig, err := discoverExtensionConfig(ctx, r.RuntimeClient, extensionConfig)
	if err != nil {
		errs = append(errs, err)
	}

	// Always patch the ExtensionConfig as it may contain updates in conditions or clientConfig.caBundle.
	if err = patchExtensionConfig(ctx, r.Client, original, discoveredExtensionConfig); err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return ctrl.Result{}, kerrors.NewAggregate(errs)
	}

	// Register the ExtensionConfig if it was found and patched without error.
	log.V(4).Info("Registering ExtensionConfig information into registry")
	if err = r.RuntimeClient.Register(discoveredExtensionConfig); err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to register ExtensionConfig %s/%s", extensionConfig.Namespace, extensionConfig.Name)
	}
	return ctrl.Result{}, nil
}

func patchExtensionConfig(ctx context.Context, client client.Client, original, modified *runtimev1.ExtensionConfig, options ...patch.Option) error {
	patchHelper, err := patch.NewHelper(original, client)
	if err != nil {
		return err
	}

	options = append(options,
		patch.WithOwnedV1beta1Conditions{Conditions: []clusterv1.ConditionType{
			runtimev1.RuntimeExtensionDiscoveredV1Beta1Condition,
		}},
		patch.WithOwnedConditions{Conditions: []string{
			clusterv1.PausedCondition,
			runtimev1.ExtensionConfigDiscoveredCondition,
		}},
	)
	return patchHelper.Patch(ctx, modified, options...)
}

// reconcileDelete will remove the ExtensionConfig from the registry on deletion of the object. Note this is a best
// effort deletion that may not catch all cases.
func (r *Reconciler) reconcileDelete(ctx context.Context, extensionConfig *runtimev1.ExtensionConfig) (ctrl.Result, error) {
	log := ctrl.LoggerFrom(ctx)
	log.Info("Unregistering ExtensionConfig information from registry")
	if err := r.RuntimeClient.Unregister(extensionConfig); err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to unregister ExtensionConfig %s", klog.KObj(extensionConfig))
	}
	return ctrl.Result{}, nil
}

// secretToExtensionConfig maps a secret to ExtensionConfigs with the corresponding InjectCAFromSecretAnnotation
// to reconcile them on updates of the secrets.
func (r *Reconciler) secretToExtensionConfig(ctx context.Context, secret *metav1.PartialObjectMetadata) []reconcile.Request {
	result := []ctrl.Request{}

	extensionConfigs := runtimev1.ExtensionConfigList{}
	indexKey := secret.GetNamespace() + "/" + secret.GetName()

	if err := r.Client.List(
		ctx,
		&extensionConfigs,
		client.MatchingFields{injectCAFromSecretAnnotationField: indexKey},
	); err != nil {
		return nil
	}

	for _, ext := range extensionConfigs.Items {
		result = append(result, ctrl.Request{NamespacedName: client.ObjectKey{Name: ext.Name}})
	}

	return result
}

// discoverExtensionConfig attempts to discover the Handlers for an ExtensionConfig.
// If discovery succeeds it returns the ExtensionConfig with Handlers updated in Status and an updated Condition.
// If discovery fails it returns the ExtensionConfig with no update to Handlers and a Failed Condition.
func discoverExtensionConfig(ctx context.Context, runtimeClient runtimeclient.Client, extensionConfig *runtimev1.ExtensionConfig) (*runtimev1.ExtensionConfig, error) {
	discoveredExtension, err := runtimeClient.Discover(ctx, extensionConfig.DeepCopy())
	if err != nil {
		modifiedExtensionConfig := extensionConfig.DeepCopy()
		v1beta1conditions.MarkFalse(modifiedExtensionConfig, runtimev1.RuntimeExtensionDiscoveredV1Beta1Condition, runtimev1.DiscoveryFailedV1Beta1Reason, clusterv1.ConditionSeverityError, "Error in discovery: %v", err)
		conditions.Set(modifiedExtensionConfig, metav1.Condition{
			Type:    runtimev1.ExtensionConfigDiscoveredCondition,
			Status:  metav1.ConditionFalse,
			Reason:  runtimev1.ExtensionConfigNotDiscoveredReason,
			Message: fmt.Sprintf("Error in discovery: %v", err),
		})
		return modifiedExtensionConfig, errors.Wrapf(err, "failed to discover ExtensionConfig %s", klog.KObj(extensionConfig))
	}

	v1beta1conditions.MarkTrue(discoveredExtension, runtimev1.RuntimeExtensionDiscoveredV1Beta1Condition)
	conditions.Set(discoveredExtension, metav1.Condition{
		Type:   runtimev1.ExtensionConfigDiscoveredCondition,
		Status: metav1.ConditionTrue,
		Reason: runtimev1.ExtensionConfigDiscoveredReason,
	})
	return discoveredExtension, nil
}

// reconcileCABundle reconciles the CA bundle for the ExtensionConfig.
// Note: This was implemented to behave similar to the cert-manager cainjector.
// We couldn't use the cert-manager cainjector because it doesn't work with CustomResources.
func reconcileCABundle(ctx context.Context, client client.Client, config *runtimev1.ExtensionConfig) error {
	log := ctrl.LoggerFrom(ctx)

	secretNameRaw, ok := config.Annotations[runtimev1.InjectCAFromSecretAnnotation]
	if !ok {
		return nil
	}
	secretName := splitNamespacedName(secretNameRaw)

	log.V(4).Info(fmt.Sprintf("Injecting CA Bundle into ExtensionConfig from secret %q", secretNameRaw))

	if secretName.Namespace == "" || secretName.Name == "" {
		return errors.Errorf("failed to reconcile caBundle: secret name %q must be in the form <namespace>/<name>", secretNameRaw)
	}

	var secret corev1.Secret
	// Note: this is an expensive API call because secrets are explicitly not cached.
	if err := client.Get(ctx, secretName, &secret); err != nil {
		return errors.Wrapf(err, "failed to reconcile caBundle: failed to get secret %q", secretNameRaw)
	}

	caData, hasCAData := secret.Data[tlsCAKey]
	if !hasCAData {
		return errors.Errorf("failed to reconcile caBundle: secret %s does not contain a %q entry", secretNameRaw, tlsCAKey)
	}

	config.Spec.ClientConfig.CABundle = caData
	return nil
}

// splitNamespacedName turns the string form of a namespaced name
// (<namespace>/<name>) into a types.NamespacedName.
func splitNamespacedName(nameStr string) types.NamespacedName {
	splitPoint := strings.IndexRune(nameStr, types.Separator)
	if splitPoint == -1 {
		return types.NamespacedName{Name: nameStr}
	}
	return types.NamespacedName{Namespace: nameStr[:splitPoint], Name: nameStr[splitPoint+1:]}
}
