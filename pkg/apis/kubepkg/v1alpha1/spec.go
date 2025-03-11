package v1alpha1

type Spec struct {
	Version string `json:"version"`

	Deploy         Deploy                       `json:"deploy,omitzero"`
	Config         map[string]EnvVarValueOrFrom `json:"config,omitzero"`
	Containers     map[string]Container         `json:"containers,omitzero"`
	Volumes        map[string]Volume            `json:"volumes,omitzero"`
	Services       map[string]Service           `json:"services,omitzero"`
	ServiceAccount *ServiceAccount              `json:"serviceAccount,omitzero"`
	Manifests      Manifests                    `json:"manifests,omitzero"`
	Images         map[string]Image             `json:"images,omitzero"`
}

type Manifests = map[string]any
