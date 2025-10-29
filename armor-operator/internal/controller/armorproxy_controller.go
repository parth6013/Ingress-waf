/*
Copyright 2025.

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

package controller

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"

	appsv1 "armor.io/operator/api/v1"
	appsdv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/yaml"
)

// ArmorProxyReconciler reconciles a ArmorProxy object
type ArmorProxyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.armor.io,resources=armorproxies,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.armor.io,resources=armorproxies/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.armor.io,resources=armorproxies/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ArmorProxy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *ArmorProxyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("Starting reconciliation", "namespace", req.Namespace, "name", req.Name)

	// Fetch the ArmorProxy instance
	armorProxy := &appsv1.ArmorProxy{}
	if err := r.Get(ctx, req.NamespacedName, armorProxy); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("ArmorProxy resource not found. Resources already cleaned up by finalizer.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Unable to fetch ArmorProxy")
		return ctrl.Result{}, err
	}

	// Add finalizer logic for proper cleanup
	myFinalizerName := "armor.io/finalizer"

	// Check if the object is being deleted
	if !armorProxy.ObjectMeta.DeletionTimestamp.IsZero() {
		// Object is being deleted
		if containsString(armorProxy.GetFinalizers(), myFinalizerName) {
			// Run cleanup logic - now we can still access the Spec!
			if err := r.cleanupResources(ctx, armorProxy); err != nil {
				log.Error(err, "Failed to clean up resources")
				return ctrl.Result{}, err
			}

			// Remove the finalizer once cleanup is done
			armorProxy.SetFinalizers(removeString(armorProxy.GetFinalizers(), myFinalizerName))
			if err := r.Update(ctx, armorProxy); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{}, nil
		}
		// If there are no finalizers, just return
		return ctrl.Result{}, nil
	}

	// Add finalizer if it doesn't exist (this happens during resource creation)
	if !containsString(armorProxy.GetFinalizers(), myFinalizerName) {
		log.Info("Adding finalizer to ArmorProxy", "name", armorProxy.Name)
		armorProxy.SetFinalizers(append(armorProxy.GetFinalizers(), myFinalizerName))
		if err := r.Update(ctx, armorProxy); err != nil {
			log.Error(err, "Failed to add finalizer")
			return ctrl.Result{}, err
		}
		// No need to requeue here - the update will trigger another reconciliation
		return ctrl.Result{}, nil
	}

	// Continue with normal reconciliation - create/update resources...

	// Fetch the associated Service
	service := &corev1.Service{}
	serviceKey := client.ObjectKey{
		Namespace: armorProxy.Spec.ServiceNamespace,
		Name:      armorProxy.Spec.ServiceName,
	}
	if err := r.Get(ctx, serviceKey, service); err != nil {
		log.Error(err, "Unable to fetch associated Service")
		return ctrl.Result{}, err
	}

	// Extract fields from the Service
	targetHost := service.Name
	targetPort := service.Spec.Ports[0].Port // Assuming the first port is used

	// Render and Create/Update ConfigMap
	configMap := &corev1.ConfigMap{}

	// Prepare Elasticsearch configuration
	elasticsearchEnable := armorProxy.Spec.Elasticsearch.Enable
	elasticsearchURL := armorProxy.Spec.Elasticsearch.URL
	if len(elasticsearchURL) == 0 {
		elasticsearchURL = []string{"http://localhost:9200"}
	}

	// Prepare CRS rules path
	crsRulesPath := armorProxy.Spec.CrsRulesPath
	if crsRulesPath == "" {
		crsRulesPath = "/crs4"
	}

	// Prepare OIDC configuration
	oidcClientID := armorProxy.Spec.OidcConfig.ClientID
	if oidcClientID == "" {
		oidcClientID = ""
	}
	oidcProviderUrl := armorProxy.Spec.OidcConfig.ProviderUrl
	if oidcProviderUrl == "" {
		oidcProviderUrl = ""
	}
	oidcRedirectApplication := armorProxy.Spec.OidcConfig.RedirectApplication
	if oidcRedirectApplication == "" {
		oidcRedirectApplication = ""
	}

	configMapData := map[string]interface{}{
		"ServiceName":             armorProxy.Spec.ServiceName,
		"Waf":                     armorProxy.Spec.Waf,
		"Oidc":                    armorProxy.Spec.Oidc,
		"TargetHost":              targetHost,
		"TargetPort":              targetPort,
		"CrsRulesPath":            crsRulesPath,
		"ElasticsearchEnable":     elasticsearchEnable,
		"ElasticsearchURL":        elasticsearchURL,
		"OidcClientID":            oidcClientID,
		"OidcProviderUrl":         oidcProviderUrl,
		"OidcRedirectApplication": oidcRedirectApplication,
	}
	if err := r.renderK8sResourceTemplate("templates/configmap.yaml", configMapData, configMap); err != nil {
		log.Error(err, "Failed to render ConfigMap template")
		return ctrl.Result{}, err
	}

	if err := r.createOrUpdateResource(ctx, req, configMap); err != nil {
		log.Error(err, "Failed to create or update ConfigMap")
		return ctrl.Result{}, err
	}

	// Render and Create/Update Deployment
	deployment := &appsdv1.Deployment{}
	deploymentData := map[string]interface{}{
		"ServiceName":      armorProxy.Spec.ServiceName,
		"ServiceNamespace": armorProxy.Spec.ServiceNamespace,
	}
	if err := r.renderK8sResourceTemplate("templates/deployment.yaml", deploymentData, deployment); err != nil {
		log.Error(err, "Failed to render Deployment template")
		return ctrl.Result{}, err
	}
	if err := r.createOrUpdateResource(ctx, req, deployment); err != nil {
		log.Error(err, "Failed to create or update Deployment")
		return ctrl.Result{}, err
	}

	// Render and Create/Update Service
	service = &corev1.Service{}
	serviceData := map[string]interface{}{
		"ServiceName":      armorProxy.Spec.ServiceName,
		"ServiceNamespace": armorProxy.Spec.ServiceNamespace,
	}
	if err := r.renderK8sResourceTemplate("templates/service.yaml", serviceData, service); err != nil {
		log.Error(err, "Failed to render Service template")
		return ctrl.Result{}, err
	}
	if err := r.createOrUpdateResource(ctx, req, service); err != nil {
		log.Error(err, "Failed to create or update Service")
		return ctrl.Result{}, err
	}

	// Render and Create/Update Ingress (if IngressHost is provided)
	if armorProxy.Spec.IngressHost != "" {
		ingress := &networkingv1.Ingress{}
		ingressData := map[string]interface{}{
			"ServiceName":      armorProxy.Spec.ServiceName,
			"ServiceNamespace": armorProxy.Spec.ServiceNamespace,
			"IngressHost":      armorProxy.Spec.IngressHost,
			"IngressClassName": armorProxy.Spec.IngressClassName,
		}
		if err := r.renderK8sResourceTemplate("templates/ingress.yaml", ingressData, ingress); err != nil {
			log.Error(err, "Failed to render Ingress template")
			return ctrl.Result{}, err
		}
		if err := r.createOrUpdateResource(ctx, req, ingress); err != nil {
			log.Error(err, "Failed to create or update Ingress")
			return ctrl.Result{}, err
		}
	}

	log.Info("Reconciliation complete")
	return ctrl.Result{}, nil
}

func (r *ArmorProxyReconciler) renderK8sResourceTemplate(templatePath string, data interface{}, obj client.Object) error {
	// Render the template
	rendered, err := r.renderTemplate(templatePath, data)
	if err != nil {
		return err
	}

	// Use Kubernetes serializer for proper object unmarshaling
	deserializer := serializer.NewCodecFactory(r.Scheme).UniversalDeserializer()

	// Convert YAML to JSON (as the deserializer works better with JSON)
	jsonData, err := yaml.YAMLToJSON([]byte(rendered))
	if err != nil {
		fmt.Printf("Error converting YAML to JSON: %v\n", err)
		return err
	}

	// Deserialize the JSON into the object
	_, _, err = deserializer.Decode(jsonData, nil, obj)
	if err != nil {
		fmt.Printf("Error decoding JSON: %v\n", err)
		return err
	}

	return nil
}

func (r *ArmorProxyReconciler) renderTemplate(templatePath string, data interface{}) (string, error) {
	// Use absolute path for templates
	fullPath := "/templates/" + templatePath

	// For local development, fall back to relative path if absolute doesn't exist
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		fullPath = templatePath
	}

	tmpl, err := template.ParseFiles(fullPath)
	if err != nil {
		return "", fmt.Errorf("error parsing template %s: %v", fullPath, err)
	}

	var rendered bytes.Buffer
	if err := tmpl.Execute(&rendered, data); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return rendered.String(), nil
}

func (r *ArmorProxyReconciler) createOrUpdateResource(ctx context.Context, req ctrl.Request, obj client.Object) error {
	// Add owner labels and annotations
	labels := obj.GetLabels()
	if labels == nil {
		labels = make(map[string]string)
	}
	labels["app.kubernetes.io/managed-by"] = "armor-operator"
	obj.SetLabels(labels)

	// Add owner reference annotation
	annotations := obj.GetAnnotations()
	if annotations == nil {
		annotations = make(map[string]string)
	}

	// Try to get owner reference from context or namespace
	ns := req.Namespace
	name := req.Name
	annotations["armor.io/armorproxy-name"] = name
	annotations["armor.io/armorproxy-namespace"] = ns
	obj.SetAnnotations(annotations)

	// Continue with create or update
	existing := obj.DeepCopyObject().(client.Object)
	err := r.Get(ctx, client.ObjectKeyFromObject(obj), existing)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return r.Create(ctx, obj)
		}
		return err
	}
	obj.SetResourceVersion(existing.GetResourceVersion())
	return r.Update(ctx, obj)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ArmorProxyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.ArmorProxy{}).
		Complete(r)
}

func (r *ArmorProxyReconciler) cleanupResources(ctx context.Context, armorProxy *appsv1.ArmorProxy) error {
	log := log.FromContext(ctx)
	log.Info("Cleaning up resources for ArmorProxy", "namespace", armorProxy.Namespace, "name", armorProxy.Name)

	// Find resources with label selector
	configMaps := &corev1.ConfigMapList{}
	listOpts := []client.ListOption{
		client.InNamespace("armor"),
		client.MatchingLabels(map[string]string{"app.kubernetes.io/managed-by": "armor-operator"}),
	}

	if err := r.List(ctx, configMaps, listOpts...); err != nil {
		log.Error(err, "Failed to list ConfigMaps")
	} else {
		for _, cm := range configMaps.Items {
			// Check if this ConfigMap belongs to the ArmorProxy
			if cm.GetName() == "armor-config-"+armorProxy.Spec.ServiceName ||
				cm.GetAnnotations()["armor.io/armorproxy-name"] == armorProxy.Name {
				if err := r.Delete(ctx, &cm); client.IgnoreNotFound(err) != nil {
					log.Error(err, "Failed to delete ConfigMap", "name", cm.GetName())
				} else {
					log.Info("Deleted ConfigMap", "name", cm.GetName())
				}
			}
		}
	}

	// Similarly, clean up deployments
	deployments := &appsdv1.DeploymentList{}
	if err := r.List(ctx, deployments, listOpts...); err != nil {
		log.Error(err, "Failed to list Deployments")
	} else {
		for _, deploy := range deployments.Items {
			if deploy.GetName() == "armor-"+armorProxy.Spec.ServiceName ||
				deploy.GetAnnotations()["armor.io/armorproxy-name"] == armorProxy.Name {
				if err := r.Delete(ctx, &deploy); client.IgnoreNotFound(err) != nil {
					log.Error(err, "Failed to delete Deployment", "name", deploy.GetName())
				} else {
					log.Info("Deleted Deployment", "name", deploy.GetName())
				}
			}
		}
	}

	// Clean up services
	services := &corev1.ServiceList{}
	if err := r.List(ctx, services, listOpts...); err != nil {
		log.Error(err, "Failed to list Services")
	} else {
		for _, svc := range services.Items {
			if svc.GetName() == "armor-"+armorProxy.Spec.ServiceName ||
				svc.GetAnnotations()["armor.io/armorproxy-name"] == armorProxy.Name {
				if err := r.Delete(ctx, &svc); client.IgnoreNotFound(err) != nil {
					log.Error(err, "Failed to delete Service", "name", svc.GetName())
				} else {
					log.Info("Deleted Service", "name", svc.GetName())
				}
			}
		}
	}

	// Clean up ingresses
	ingresses := &networkingv1.IngressList{}
	if err := r.List(ctx, ingresses, listOpts...); err != nil {
		log.Error(err, "Failed to list Ingresses")
	} else {
		for _, ing := range ingresses.Items {
			if ing.GetName() == "armor-"+armorProxy.Spec.ServiceName ||
				ing.GetAnnotations()["armor.io/armorproxy-name"] == armorProxy.Name {
				if err := r.Delete(ctx, &ing); client.IgnoreNotFound(err) != nil {
					log.Error(err, "Failed to delete Ingress", "name", ing.GetName())
				} else {
					log.Info("Deleted Ingress", "name", ing.GetName())
				}
			}
		}
	}

	return nil
}

func containsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func removeString(slice []string, str string) []string {
	var result []string
	for _, item := range slice {
		if item != str {
			result = append(result, item)
		}
	}
	return result
}
