package converter

import (
	// "fmt"
	// "strings"

	"k8s.io/klog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func convertExampleCRD(Object *unstructured.Unstructured, toVersion string) (*unstructured.Unstructured, metav1.Status) {
	klog.V(2).Info("converting crd")

	convertedObject := Object.DeepCopy()
	fromVersion := Object.GetAPIVersion()

	if toVersion == fromVersion {
		return nil, statusErrorWithMessage("conversion from a version to itself should not call the webhook: %s", toVersion)
	}

	// apiVersion: com.cisco.ccc/v2
	// kind: WxccService
	// metadata:
	//   name: wxccservice-sample
	// spec:
	//   global:
	// 	appName: helloapp
	// 	appPrefix: test
	// 	fullname: helloapp

	switch Object.GetAPIVersion() {
	case "stable.example.com/v1":
		switch toVersion {
		case "stable.example.com/v2":
			//spec, ok, _ := unstructured.NestedMap(convertedObject.Object, "spec")
			Fullname, hasFullname, _ := unstructured.NestedString(convertedObject.Object, "spec", "fullname")
			if hasFullname {
				//delete(convertedObject.Object, "spec", "fullname")
				// parts := strings.Split(hostPort, ":")
				// if len(parts) != 2 {
				// 	return nil, statusErrorWithMessage("invalid hostPort value `%v`", hostPort)
				// }
				// host := parts[0]
				// port := parts[1]
				unstructured.SetNestedField(convertedObject.Object, Fullname, "spec", "full")
			}
		default:
			return nil, statusErrorWithMessage("unexpected conversion version %q", toVersion)
		}
	case "stable.example.com/v2":
		switch toVersion {
		case "stable.example.com/v1":
			Full, hasFull, _ := unstructured.NestedString(convertedObject.Object, "spec", "full")
			if hasFull {
				unstructured.SetNestedField(convertedObject.Object, Full, "spec", "fullname")
			}
			// host, hasHost, _ := unstructured.NestedString(convertedObject.Object, "spec", "host")
			// port, hasPort, _ := unstructured.NestedString(convertedObject.Object, "spec", "port")
			// if hasHost || hasPort {
			// 	if !hasHost {
			// 		host = ""
			// 	}
			// 	if !hasPort {
			// 		port = ""
			// 	}
			// 	hostPort := fmt.Sprintf("%s:%s", host, port)
			// 	unstructured.SetNestedField(convertedObject.Object, hostPort, "spec", "hostPort")
			// }
		default:
			return nil, statusErrorWithMessage("unexpected conversion version %q", toVersion)
		}
	default:
		return nil, statusErrorWithMessage("unexpected conversion version %q", fromVersion)
	}
	return convertedObject, statusSucceed()
}
