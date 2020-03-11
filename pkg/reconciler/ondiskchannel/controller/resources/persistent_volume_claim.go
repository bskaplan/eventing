package resources 

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
type PersistentVolumeClaimArgs struct{
	Size string
	Name string
	Namespace string
}

func MakePersistentVolumeClaim( args PersistentVolumeClaimArgs) *corev1.PersistentVolumeClaim {
	return &corev1.PersistentVolumeClaim {
		TypeMeta: metav1.TypeMeta {
			APIVersion: "core/v1",
			Kind: "PersistentVolumeClaim",
		},
		ObjectMeta: metav1.ObjectMeta {
			Name: args.Name,
			Namespace: args.Namespace,
		},


	}
}