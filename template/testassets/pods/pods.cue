package pods

import (
	corev1 "k8s.io/api/core/v1"
)


defaults: {}
resource: corev1.#Pod
template: corev1.#Pod
