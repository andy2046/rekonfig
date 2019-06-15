package testutil

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var labels = map[string]string{
	"app": "example",
}

var trueValue = true

// ExampleDeployment is an example Deployment object
var ExampleDeployment = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example",
		Namespace: "default",
		Labels:    labels,
	},
	Spec: appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: labels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels,
			},
			Spec: corev1.PodSpec{
				Volumes: []corev1.Volume{
					{
						Name: "secret1",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: "example1",
							},
						},
					},
					{
						Name: "configmap1",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example1",
								},
							},
						},
					},
				},
				Containers: []corev1.Container{
					{
						Name:  "container1",
						Image: "container1",
						Env: []corev1.EnvVar{
							{
								Name: "example1_key1",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example1",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example1_key1_new_name",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example1",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_key1",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_key4",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key:      "key4",
										Optional: &trueValue,
									},
								},
							},
							{
								Name: "example4_key1",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example4",
										},
										Key:      "key1",
										Optional: &trueValue,
									},
								},
							},
							{
								Name: "example1_secret_key1",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example1",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_secret_key1",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_secret_key4",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key:      "key4",
										Optional: &trueValue,
									},
								},
							},
							{
								Name: "example4_secret_key1",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example4",
										},
										Key:      "key1",
										Optional: &trueValue,
									},
								},
							},
						},
						EnvFrom: []corev1.EnvFromSource{
							{
								ConfigMapRef: &corev1.ConfigMapEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example1",
									},
								},
							},
							{
								SecretRef: &corev1.SecretEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example1",
									},
								},
							},
						},
					},
					{
						Name:  "container2",
						Image: "container2",
						Env: []corev1.EnvVar{
							{
								Name: "example3_key2",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key2",
									},
								},
							},
							{
								Name: "example3_secret_key2",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key2",
									},
								},
							},
						},
						EnvFrom: []corev1.EnvFromSource{
							{
								ConfigMapRef: &corev1.ConfigMapEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example2",
									},
								},
							},
							{
								SecretRef: &corev1.SecretEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example2",
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

// ExampleStatefulSet is an example StatefulSet object
var ExampleStatefulSet = &appsv1.StatefulSet{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example",
		Namespace: "default",
		Labels:    labels,
	},
	Spec: appsv1.StatefulSetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: labels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels,
			},
			Spec: corev1.PodSpec{
				Volumes: []corev1.Volume{
					{
						Name: "secret1",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: "example1",
							},
						},
					},
					{
						Name: "configmap1",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example1",
								},
							},
						},
					},
				},
				Containers: []corev1.Container{
					{
						Name:  "container1",
						Image: "container1",
						Env: []corev1.EnvVar{
							{
								Name: "example1_key1",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example1",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example1_key1_new_name",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example1",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_key1",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_key4",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key:      "key4",
										Optional: &trueValue,
									},
								},
							},
							{
								Name: "example4_key1",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example4",
										},
										Key:      "key1",
										Optional: &trueValue,
									},
								},
							},
							{
								Name: "example1_secret_key1",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example1",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_secret_key1",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_secret_key4",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key:      "key4",
										Optional: &trueValue,
									},
								},
							},
							{
								Name: "example4_secret_key1",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example4",
										},
										Key:      "key1",
										Optional: &trueValue,
									},
								},
							},
						},
						EnvFrom: []corev1.EnvFromSource{
							{
								ConfigMapRef: &corev1.ConfigMapEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example1",
									},
								},
							},
							{
								SecretRef: &corev1.SecretEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example1",
									},
								},
							},
						},
					},
					{
						Name:  "container2",
						Image: "container2",
						Env: []corev1.EnvVar{
							{
								Name: "example3_key2",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key2",
									},
								},
							},
							{
								Name: "example3_secret_key2",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key2",
									},
								},
							},
						},
						EnvFrom: []corev1.EnvFromSource{
							{
								ConfigMapRef: &corev1.ConfigMapEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example2",
									},
								},
							},
							{
								SecretRef: &corev1.SecretEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example2",
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

// ExampleDaemonSet is an example DaemonSet object
var ExampleDaemonSet = &appsv1.DaemonSet{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example",
		Namespace: "default",
		Labels:    labels,
	},
	Spec: appsv1.DaemonSetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: labels,
		},
		Template: corev1.PodTemplateSpec{
			ObjectMeta: metav1.ObjectMeta{
				Labels: labels,
			},
			Spec: corev1.PodSpec{
				Volumes: []corev1.Volume{
					{
						Name: "secret1",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: "example1",
							},
						},
					},
					{
						Name: "configmap1",
						VolumeSource: corev1.VolumeSource{
							ConfigMap: &corev1.ConfigMapVolumeSource{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example1",
								},
							},
						},
					},
				},
				Containers: []corev1.Container{
					{
						Name:  "container1",
						Image: "container1",
						Env: []corev1.EnvVar{
							{
								Name: "example1_key1",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example1",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example1_key1_new_name",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example1",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_key1",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_key4",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key:      "key4",
										Optional: &trueValue,
									},
								},
							},
							{
								Name: "example4_key1",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example4",
										},
										Key:      "key1",
										Optional: &trueValue,
									},
								},
							},
							{
								Name: "example1_secret_key1",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example1",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_secret_key1",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key1",
									},
								},
							},
							{
								Name: "example3_secret_key4",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key:      "key4",
										Optional: &trueValue,
									},
								},
							},
							{
								Name: "example4_secret_key1",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example4",
										},
										Key:      "key1",
										Optional: &trueValue,
									},
								},
							},
						},
						EnvFrom: []corev1.EnvFromSource{
							{
								ConfigMapRef: &corev1.ConfigMapEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example1",
									},
								},
							},
							{
								SecretRef: &corev1.SecretEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example1",
									},
								},
							},
						},
					},
					{
						Name:  "container2",
						Image: "container2",
						Env: []corev1.EnvVar{
							{
								Name: "example3_key2",
								ValueFrom: &corev1.EnvVarSource{
									ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key2",
									},
								},
							},
							{
								Name: "example3_secret_key2",
								ValueFrom: &corev1.EnvVarSource{
									SecretKeyRef: &corev1.SecretKeySelector{
										LocalObjectReference: corev1.LocalObjectReference{
											Name: "example3",
										},
										Key: "key2",
									},
								},
							},
						},
						EnvFrom: []corev1.EnvFromSource{
							{
								ConfigMapRef: &corev1.ConfigMapEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example2",
									},
								},
							},
							{
								SecretRef: &corev1.SecretEnvSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: "example2",
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

// ExampleConfigMap1 is an example ConfigMap object
var ExampleConfigMap1 = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example1",
		Namespace: "default",
		Labels:    labels,
	},
	Data: map[string]string{
		"key1": "example1:key1",
		"key2": "example1:key2",
		"key3": "example1:key3",
	},
}

// ExampleConfigMap2 is an example ConfigMap object
var ExampleConfigMap2 = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example2",
		Namespace: "default",
		Labels:    labels,
	},
	Data: map[string]string{
		"key1": "example2:key1",
		"key2": "example2:key2",
		"key3": "example2:key3",
	},
}

// ExampleConfigMap3 is an example ConfigMap object
var ExampleConfigMap3 = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example3",
		Namespace: "default",
		Labels:    labels,
	},
	Data: map[string]string{
		"key1": "example3:key1",
		"key2": "example3:key2",
		"key3": "example3:key3",
	},
}

// ExampleConfigMap4 is an example ConfigMap object
var ExampleConfigMap4 = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example4",
		Namespace: "default",
		Labels:    labels,
	},
	Data: map[string]string{
		"key1": "example4:key1",
		"key2": "example4:key2",
		"key3": "example4:key3",
	},
}

// ExampleSecret1 is an example Secret object
var ExampleSecret1 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example1",
		Namespace: "default",
		Labels:    labels,
	},
	StringData: map[string]string{
		"key1": "example1:key1",
		"key2": "example1:key2",
		"key3": "example1:key3",
	},
	Data: map[string][]byte{},
}

// ExampleSecret2 is an example Secret object
var ExampleSecret2 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example2",
		Namespace: "default",
		Labels:    labels,
	},
	StringData: map[string]string{
		"key1": "example2:key1",
		"key2": "example2:key2",
		"key3": "example2:key3",
	},
	Data: map[string][]byte{},
}

// ExampleSecret3 is an example Secret object
var ExampleSecret3 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example3",
		Namespace: "default",
		Labels:    labels,
	},
	StringData: map[string]string{
		"key1": "example3:key1",
		"key2": "example3:key2",
		"key3": "example3:key3",
	},
	Data: map[string][]byte{},
}

// ExampleSecret4 is an example Secret object
var ExampleSecret4 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example4",
		Namespace: "default",
		Labels:    labels,
	},
	StringData: map[string]string{
		"key1": "example4:key1",
		"key2": "example4:key2",
		"key3": "example4:key3",
	},
	Data: map[string][]byte{},
}
