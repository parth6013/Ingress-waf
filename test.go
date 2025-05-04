// import (
// 	"context"
// 	"fmt"
// 	"testing"

// 	"github.com/cisco-it-cloud-infrastructure/openstack/openstack"
// 	"github.com/stretchr/testify/assert"
// 	v1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
// 	"k8s.io/apimachinery/pkg/runtime"
// 	dynamicfake "k8s.io/client-go/dynamic/fake"
// 	"k8s.io/client-go/kubernetes/fake"
// 	// "k8s.io/client-go/testing"
// )

// // TestCheckApplicationHealth tests the CheckApplicationHealth function
// func TestCheckApplicationHealth(t *testing.T) {
// 	// Create a fake dynamic client
// 	dynamicClient := dynamicfake.NewSimpleDynamicClient(runtime.NewScheme(), &unstructured.Unstructured{
// 		Object: map[string]interface{}{
// 			"apiVersion": "argoproj.io/v1alpha1",
// 			"kind":       "Application",
// 			"metadata": map[string]interface{}{
// 				"name":      "test-app",
// 				"namespace": "argocd",
// 			},
// 			"status": map[string]interface{}{
// 				"health": map[string]interface{}{
// 					"status": "Healthy",
// 				},
// 				"sync": map[string]interface{}{
// 					"status": "Synced",
// 				},
// 			},
// 		},
// 	})

// 	namespace := "argocd"
// 	appName := "test-app"

// 	health, err := CheckApplicationHealth(dynamicClient, namespace, appName)
// 	if err != nil {
// 		t.Fatalf("unexpected error: %v", err)
// 	}

// 	if health != "Healthy" {
// 		t.Fatalf("expected Healthy, got %s", health)
// 	}
// }

// // func TestCheckApplicationHealth_NotHealthy(t *testing.T) {
// // 	gvr := schema.GroupVersionResource{Group: "argoproj.io", Version: "v1alpha1", Resource: "applications"}
// // 	fakeClient := fake.NewSimpleDynamicClient()

// // 	app := &unstructured.Unstructured{
// // 		Object: map[string]interface{}{
// // 			"apiVersion": "argoproj.io/v1alpha1",
// // 			"kind":       "Application",
// // 			"metadata": map[string]interface{}{
// // 				"name": "test-app",
// // 			},
// // 			"status": map[string]interface{}{
// // 				"health": map[string]interface{}{"status": "Degraded"},
// // 				"sync":   map[string]interface{}{"status": "OutOfSync"},
// // 			},
// // 		},
// // 	}
// // 	_, _ = fakeClient.Resource(gvr).Namespace("argocd").Create(context.TODO(), app, metav1.CreateOptions{})

// // 	health, err := deployment.CheckApplicationHealth(fakeClient, "argocd", "test-app")
// // 	assert.NoError(t, err)
// // 	assert.Equal(t, "Not Healthy", health)
// // }

// func TestUpdateSecretDeploymen(t *testing.T) {
// 	fakeClient := fake.NewSimpleClientset()
// 	secretName := "test-secret"

// 	// Create a secret without the label
// 	secret := &v1.Secret{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      secretName,
// 			Namespace: "default",
// 			Labels:    map[string]string{},
// 		},
// 	}
// 	_, _ = fakeClient.CoreV1().Secrets("default").Create(context.TODO(), secret, metav1.CreateOptions{})

// 	err := UpdateSecretDeployment(fakeClient, "default", secretName, "test-label")
// 	assert.NoError(t, err)

// 	updatedSecret, _ := fakeClient.CoreV1().Secrets("default").Get(context.TODO(), secretName, metav1.GetOptions{})
// 	assert.Equal(t, "true", updatedSecret.Labels["test-label"])
// }

// func TestGetDeploymentConfig(t *testing.T) {
// 	fakeClient := fake.NewSimpleClientset()
// 	configMapData := `phases:
//   enable_phase_parth_1:
//     wait_time: "2"
//     applications:
//       - name: "app1"
// `

// 	configMap := &v1.ConfigMap{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "addon-deployment-config",
// 			Namespace: "argocd",
// 		},
// 		Data: map[string]string{"config.yaml": configMapData},
// 	}
// 	_, _ = fakeClient.CoreV1().ConfigMaps("argocd").Create(context.TODO(), configMap, metav1.CreateOptions{})

// 	config, err := GetDeploymentConfig(fakeClient, "argocd", "addon-deployment-config")
// 	assert.NoError(t, err)
// 	assert.Equal(t, "2", config.Phases["enable_phase_parth_1"].WaitTime)
// 	assert.Equal(t, "app1", config.Phases["enable_phase_parth_1"].Applications[0].Name)
// }

// func TestValues(t *testing.T) {
// 	// // Values()
// 	// projectId, _ := SetProjectDetails()
// 	// // SetNetworkDetails()
// 	// publicNetworkID, tenantPrivateSubNetworkID, floatingIPNetwork, _, _ := SetNetworkDetails()
// 	// Values("caas-vrl-alln-dev21", publicNetworkID, tenantPrivateSubNetworkID, floatingIPNetwork, projectId)
// 	// PerformGitPush()
// 	// Values("caas-vrl-alln-dev21", "publicNetworkID", "tenantPrivateSubNetworkID", "floatingIPNetwork", "projectId")
// 	// PerformGitPush()
// 	// Values("caas-vrl-alln-dev21", "publicNetworkID", "tenantPrivateSubNetworkID", "floatingIPNetwork", "projectId", "alln", "xyz", "xyz")
// 	// deleteNetgroup("caas-vrl-alln-dev21", "multi-cloud-it.gen", "5^q1lQ8r2eA3")
// 	// deleteNetgroup("caas-vrl-alln-dev21", "mc-it-idev.gen", "BGT%6yhnMJU7")
// 	// err := createNetgroup("caas-vrl-alln-dev21", "multi-cloud-it.gen", "5^q1lQ8r2eA3")
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// _, err = SetProjectDetails("multi-cloud-it.gen", "5^q1lQ8r2eA3", "alln", "caas-vrl-alln-dev21_osp")
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }

// 	networkDetails, err := openstack.GetNetworkDetail("caekube-dev.gen", "CAh67d4%Max#1!", "alln", "CAE-EPS-JENKINS-ALLN")
// 	if err != nil {

// 		fmt.Println("error in getting network details: ")
// 		fmt.Println(err)
// 	}

// 	// projectDetatls, err := openstack.GetProjectDetail("caekube-dev.gen", "CAh67d4%Max#1!", "alln", "CAE-EPS-JENKINS-ALLN")
// 	// if err != nil {
// 	// 	fmt.Println("error in getting project details: ")
// 	// 	fmt.Println(err)
// 	// }
// 	// fmt.Println(projectDetatls)

// 	fmt.Println(networkDetails)

// }
