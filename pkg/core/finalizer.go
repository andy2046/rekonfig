package core

// addFinalizer adds the rekonfig finalizer to the given PodController
func addFinalizer(obj podController) {
	finalizers := obj.GetFinalizers()
	for _, finalizer := range finalizers {
		if finalizer == FinalizerString {
			// podController already contains the finalizer
			return
		}
	}

	//podController does not contain the finalizer, add it
	finalizers = append(finalizers, FinalizerString)
	obj.SetFinalizers(finalizers)
}

// removeFinalizer removes the rekonfig finalizer from the given podController
func removeFinalizer(obj podController) {
	finalizers := obj.GetFinalizers()

	// Filter existing finalizers removing any that match the finalizerString
	newFinalizers := []string{}
	for _, finalizer := range finalizers {
		if finalizer != FinalizerString {
			newFinalizers = append(newFinalizers, finalizer)
		}
	}

	// Update the object's finalizers
	obj.SetFinalizers(newFinalizers)
}

// hasFinalizer checks for the presence of the rekonfig finalizer
func hasFinalizer(obj podController) bool {
	finalizers := obj.GetFinalizers()
	for _, finalizer := range finalizers {
		if finalizer == FinalizerString {
			// podController already contains the finalizer
			return true
		}
	}

	return false
}
