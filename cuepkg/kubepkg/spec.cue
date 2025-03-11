package kubepkg

#Affinity: {
	// Describes node affinity scheduling rules for the pod
	nodeAffinity?: #NodeAffinity
	// Describes pod affinity scheduling rules (e
	podAffinity?: #PodAffinity
	// Describes pod anti-affinity scheduling rules (e
	podAntiAffinity?: #PodAntiAffinity
}

#AppArmorProfile: {
	// type indicates which kind of AppArmor profile will be applied
	type: #AppArmorProfileType
	// localhostProfile indicates a profile loaded on the node that should be used
	localhostProfile?: string
}

#AppArmorProfileType: string

#Capabilities: {
	// Added capabilities
	add?: [...#Capability]
	// Removed capabilities
	drop?: [...#Capability]
}

#Capability: string

#ClaimResourceStatus: string

#CompletionMode: string

#ConcurrencyPolicy: string

#ConditionStatus: string

#ConfigMapSpec: data: [X=string]: string

#ConfigMapVolumeSource: {
	name?: string
	// items if unspecified, each key-value pair in the Data field of the referenced ConfigMap will be projected into the volume as a file whose name is the key and content is the value
	items?: [...#KeyToPath]
	// defaultMode is optional: mode bits used to set permissions on created files by default
	defaultMode?: int
	// optional specify whether the ConfigMap or its keys must be defined
	optional?: bool
}

#Container: {
	// 镜像
	image: #Image
	// 运行目录
	workingDir?: string
	// 命令
	command?: [...string]
	// 参数
	args?: [...string]
	// 环境变量
	env?: {
		[X=string]: #EnvVarValueOrFrom
	}
	// 暴露端口
	ports?: {
		[X=string]: int
	}
	stdin?:                    bool
	stdinOnce?:                bool
	tty?:                      bool
	resources?:                #ResourceRequirements
	livenessProbe?:            #Probe
	readinessProbe?:           #Probe
	startupProbe?:             #Probe
	lifecycle?:                #Lifecycle
	securityContext?:          #SecurityContext
	terminationMessagePath?:   string
	terminationMessagePolicy?: #TerminationMessagePolicy
}

#CronJobSpec: {
	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	schedule: string
	// The time zone name for the given schedule, see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones.
	timeZone?: string
	// Optional deadline in seconds for starting the job if it misses scheduled
	startingDeadlineSeconds?: int
	// Specifies how to treat concurrent executions of a Job.
	concurrencyPolicy?: #ConcurrencyPolicy
	// This flag tells the controller to suspend subsequent executions, it does
	suspend?: bool
	// Specifies the job that will be created when executing a CronJob.
	template?: #JobTemplateSpec
	// The number of successful finished jobs to retain. Value must be non-negative integer.
	successfulJobsHistoryLimit?: int
	// The number of failed finished jobs to retain. Value must be non-negative integer.
	failedJobsHistoryLimit?: int
}

#DNSPolicy: string

#DaemonSetSpec: {
	// An object that describes the pod that will be created.
	template?: #PodPartialTemplateSpec
	// An update strategy to replace existing DaemonSet pods with new pods.
	updateStrategy?: #DaemonSetUpdateStrategy
	// The minimum number of seconds for which a newly created DaemonSet pod should
	minReadySeconds?: int
	// The number of old history to retain to allow rollback.
	revisionHistoryLimit?: int
}

#DaemonSetUpdateStrategy: {
	// Type of daemon set update
	type?: #DaemonSetUpdateStrategyType
	// Rolling update config params
	rollingUpdate?: #RollingUpdateDaemonSet
}

#DaemonSetUpdateStrategyType: "RollingUpdate" | "OnDelete"

#Deploy: #DeployConfigMap | #DeployCronJob | #DeployDaemonSet | #DeployDeployment | #DeployEndpoints | #DeployJob | #DeploySecret | #DeployStatefulSet

#DeployConfigMap: {
	kind: "ConfigMap"
	labels?: [X=string]: string
	annotations?: [X=string]: string
}

#DeployCronJob: {
	kind: "CronJob"
	labels?: [X=string]: string
	annotations?: [X=string]: string
	spec?: #CronJobSpec
}

#DeployDaemonSet: {
	kind: "DaemonSet"
	labels?: [X=string]: string
	annotations?: [X=string]: string
	spec?: #DaemonSetSpec
}

#DeployDeployment: {
	kind: "Deployment"
	labels?: [X=string]: string
	annotations?: [X=string]: string
	spec?: #DeploymentSpec
}

#DeployEndpoints: {
	kind: "Endpoints"
	labels?: [X=string]: string
	annotations?: [X=string]: string
	ports: [X=string]: int
	addresses?: [...#EndpointAddress]
}

#DeployJob: {
	kind: "Job"
	labels?: [X=string]: string
	annotations?: [X=string]: string
	spec?: #JobSpec
}

#DeploySecret: {
	kind: "Secret"
	labels?: [X=string]: string
	annotations?: [X=string]: string
}

#DeployStatefulSet: {
	kind: "StatefulSet"
	labels?: [X=string]: string
	annotations?: [X=string]: string
	spec?: #StatefulSetSpec
}

#DeploymentSpec: {
	// Number of desired pods. This is a pointer to distinguish between explicit
	replicas?: int
	// describes the pods that will be created.
	template?: #PodPartialTemplateSpec
	// The deployment strategy to use to replace existing pods with new ones.
	strategy?: #DeploymentStrategy
	// Minimum number of seconds for which a newly created pod should be ready
	minReadySeconds?: int
	// The number of old ReplicaSets to retain to allow rollback.
	revisionHistoryLimit?: int
	// Indicates that the deployment is paused.
	paused?: bool
	// The maximum time in seconds for a deployment to make progress before it
	progressDeadlineSeconds?: int
}

#DeploymentStrategy: {
	// Type of deployment
	type?: #DeploymentStrategyType
	// Rolling update config params
	rollingUpdate?: #RollingUpdateDeployment
}

#DeploymentStrategyType: "Recreate" | "RollingUpdate"

#DigestMeta: {
	type:      #DigestMetaType
	digest:    string
	name:      string
	size:      #FileSize
	tag?:      string
	platform?: string
}

#DigestMetaType: "blob" | "manifest"

#EmptyDirVolumeSource: {
	// medium represents what type of storage medium should back this directory
	medium?: #StorageMedium
	// sizeLimit is the total amount of local storage required for this EmptyDir volume
	sizeLimit?: #Quantity
}

#EndpointAddress: {
	// The IP of this endpoint
	ip: string
	// The Hostname of this endpoint
	hostname?: string
	// Optional: Node hosting this endpoint
	nodeName?: string
	// Reference to object providing the endpoint
	targetRef?: #ObjectReference
}

#EnvVarValueOrFrom: string

#ExecAction: {
	// 命令
	command?: [...string]
}

#Expose: #ExposeIngress | #ExposeNodePort

#ExposeIngress: {
	type: "Ingress"
	gateway?: [...string]
}

#ExposeNodePort: type: "NodePort"

#FieldsV1: {}

#FileSize: int

#GRPCAction: {
	// Port number of the gRPC service
	port: int
	// Service is the name of the service to place in the gRPC HealthCheckRequest (see https://github
	service: string
}

#HTTPGetAction: {
	// Path to access on the HTTP server
	path?: string
	// Name or number of the port to access on the container
	port: #IntOrString
	// Host name to connect to, defaults to the pod IP
	host?: string
	// Scheme to use for connecting to the host
	scheme?: #URIScheme
	// Custom headers to set in the request
	httpHeaders?: [...#HTTPHeader]
}

#HTTPHeader: {
	// The header field name
	name: string
	// The header field value
	value: string
}

#HostAlias: {
	// IP address of the host file entry
	ip: string
	// Hostnames for the above IP address
	hostnames?: [...string]
}

#HostPathType: string

#HostPathVolumeSource: {
	// path of the directory on the host
	path: string
	// type for HostPath Volume Defaults to "" More info: https://kubernetes
	type?: #HostPathType
}

#Image: {
	// 镜像名
	name: string
	// 镜像标签
	tag?: string
	// 镜像摘要
	digest?: string
	// 镜像支持的平台
	platforms?: [...string]
	// 镜像拉取策略
	pullPolicy?: #PullPolicy
}

#ImageVolumeSource: {
	// Required: Image or artifact reference to be used
	reference?: string
	// Policy for pulling OCI objects
	pullPolicy?: #PullPolicy
}

#IntOrString: int | string

#JobSpec: {
	// Specifies the maximum desired number of pods the job should
	parallelism?: int
	// Specifies the desired number of successfully finished pods the
	completions?: int
	// Specifies the duration in seconds relative to the startTime that the job
	activeDeadlineSeconds?: int
	// Specifies the policy of handling failed pods. In particular, it allows to
	podFailurePolicy?: #PodFailurePolicy
	// successPolicy specifies the policy when the Job can be declared as succeeded.
	successPolicy?: #SuccessPolicy
	// Specifies the number of retries before marking this job failed.
	backoffLimit?: int
	// Specifies the limit for the number of retries within an
	backoffLimitPerIndex?: int
	// Specifies the maximal number of failed indexes before marking the Job as
	maxFailedIndexes?: int
	// manualSelector controls generation of pod labels and pod selectors.
	manualSelector?: bool
	// Describes the pod that will be created when executing a job.
	template?: #PodPartialTemplateSpec
	// ttlSecondsAfterFinished limits the lifetime of a Job that has finished
	ttlSecondsAfterFinished?: int
	// completionMode specifies how Pod completions are tracked. It can be
	completionMode?: #CompletionMode
	// suspend specifies whether the Job controller should create Pods or not. If
	suspend?: bool
	// podReplacementPolicy specifies when to create replacement Pods.
	podReplacementPolicy?: #PodReplacementPolicy
	// field indicates the controller that manages a Job. The k8s Job
	managedBy?: string
}

#JobTemplateSpec: spec?: #JobSpec

#KeyToPath: {
	// key is the key to project
	key: string
	// path is the relative path of the file to map the key to
	path: string
	// mode is Optional: mode bits used to set permissions on this file
	mode?: int
}

#KubePkg: {
	kind?:       string
	apiVersion?: string
	metadata?:   #ObjectMeta
	spec:        #Spec
	status?:     #Status
}

#LabelSelector: {
	// matchLabels is a map of {key,value} pairs
	matchLabels?: {
		[X=string]: string
	}
	// matchExpressions is a list of label selector requirements
	matchExpressions?: [...#LabelSelectorRequirement]
}

#LabelSelectorOperator: string

#LabelSelectorRequirement: {
	// key is the label key that the selector applies to
	key: string
	// operator represents a key's relationship to a set of values
	operator: #LabelSelectorOperator
	// values is an array of string values
	values?: [...string]
}

#Lifecycle: {
	// PostStart is called immediately after a container is created
	postStart?: #LifecycleHandler
	// PreStop is called immediately before a container is terminated due to an API request or management event such as liveness/startup probe failure, preemption, resource contention, etc
	preStop?: #LifecycleHandler
}

#LifecycleHandler: {
	// Exec specifies a command to execute in the container
	exec?: #ExecAction
	// HTTPGet specifies an HTTP GET request to perform
	httpGet?: #HTTPGetAction
	// Deprecated
	tcpSocket?: #TCPSocketAction
	// Sleep represents a duration that the container should sleep
	sleep?: #SleepAction
}

#LocalObjectReference: {
	// Name of the referent
	name?: string
}

#ManagedFieldsEntry: {
	// Manager is an identifier of the workflow managing these fields
	manager?: string
	// Operation is the type of operation which lead to this ManagedFieldsEntry being created
	operation?: #ManagedFieldsOperationType
	// APIVersion defines the version of this resource that this field set applies to
	apiVersion?: string
	// Time is the timestamp of when the ManagedFields entry was added
	time?: #Time
	// FieldsType is the discriminator for the different fields format and version
	fieldsType?: string
	// FieldsV1 holds the first JSON version format as described in the "FieldsV1" type
	fieldsV1?: #FieldsV1
	// Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource
	subresource?: string
}

#ManagedFieldsOperationType: string

#ModifyVolumeStatus: {
	// targetVolumeAttributesClassName is the name of the VolumeAttributesClass the PVC currently being reconciled
	targetVolumeAttributesClassName?: string
	// status is the status of the ControllerModifyVolume operation
	status: #PersistentVolumeClaimModifyVolumeStatus
}

#NodeAffinity: {
	// If the affinity requirements specified by this field are not met at scheduling time, the pod will not be scheduled onto the node
	requiredDuringSchedulingIgnoredDuringExecution?: #NodeSelector
	// The scheduler will prefer to schedule pods to nodes that satisfy the affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions
	preferredDuringSchedulingIgnoredDuringExecution?: [...#PreferredSchedulingTerm]
}

#NodeInclusionPolicy: string

#NodeSelector: {
	// Required
	nodeSelectorTerms: [...#NodeSelectorTerm]
}

#NodeSelectorOperator: string

#NodeSelectorRequirement: {
	// The label key that the selector applies to
	key: string
	// Represents a key's relationship to a set of values
	operator: #NodeSelectorOperator
	// An array of string values
	values?: [...string]
}

#NodeSelectorTerm: {
	// A list of node selector requirements by node's labels
	matchExpressions?: [...#NodeSelectorRequirement]
	// A list of node selector requirements by node's fields
	matchFields?: [...#NodeSelectorRequirement]
}

#OSName: string

#ObjectMeta: {
	// Name must be unique within a namespace
	name?: string
	// GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided
	generateName?: string
	// Namespace defines the space within which each name must be unique
	namespace?: string
	// Deprecated: selfLink is a legacy read-only field that is no longer populated by the system
	selfLink?: string
	// UID is the unique in time and space value for this object
	uid?: #UID
	// An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed
	resourceVersion?: string
	// A sequence number representing a specific generation of the desired state
	generation?: int
	// CreationTimestamp is a timestamp representing the server time when this object was created
	creationTimestamp?: #Time
	// DeletionTimestamp is RFC 3339 date and time at which this resource will be deleted
	deletionTimestamp?: #Time
	// Number of seconds allowed for this object to gracefully terminate before it will be removed from the system
	deletionGracePeriodSeconds?: int
	// Map of string keys and values that can be used to organize and categorize (scope and select) objects
	labels?: {
		[X=string]: string
	}
	// Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata
	annotations?: {
		[X=string]: string
	}
	// List of objects depended by this object
	ownerReferences?: [...#OwnerReference]
	// Must be empty before the object is deleted from the registry
	finalizers?: [...string]
	// ManagedFields maps workflow-id and version to the set of fields that are managed by that workflow
	managedFields?: [...#ManagedFieldsEntry]
}

#ObjectReference: {
	// Kind of the referent
	kind?: string
	// Namespace of the referent
	namespace?: string
	// Name of the referent
	name?: string
	// UID of the referent
	uid?: #UID
	// API version of the referent
	apiVersion?: string
	// Specific resourceVersion to which this reference is made, if any
	resourceVersion?: string
	// If referring to a piece of an object instead of an entire object, this string should contain a valid JSON/Go field access statement, such as desiredState
	fieldPath?: string
}

#OwnerReference: {
	// API version of the referent
	apiVersion: string
	// Kind of the referent
	kind: string
	// Name of the referent
	name: string
	// UID of the referent
	uid: #UID
	// If true, this reference points to the managing controller
	controller?: bool
	// If true, AND if the owner has the "foregroundDeletion" finalizer, then the owner cannot be deleted from the key-value store until this reference is removed
	blockOwnerDeletion?: bool
}

#PersistentVolumeAccessMode: string

#PersistentVolumeClaim: {
	kind?:       string
	apiVersion?: string
	// Standard object's metadata
	metadata?: #ObjectMeta
	// spec defines the desired characteristics of a volume requested by a pod author
	spec?: #PersistentVolumeClaimSpec
	// status represents the current information/status of a persistent volume claim
	status?: #PersistentVolumeClaimStatus
}

#PersistentVolumeClaimCondition: {
	// Type is the type of the condition
	type: #PersistentVolumeClaimConditionType
	// Status is the status of the condition
	status: #ConditionStatus
	// lastProbeTime is the time we probed the condition
	lastProbeTime?: #Time
	// lastTransitionTime is the time the condition transitioned from one status to another
	lastTransitionTime?: #Time
	// reason is a unique, this should be a short, machine understandable string that gives the reason for condition's last transition
	reason?: string
	// message is the human-readable message indicating details about last transition
	message?: string
}

#PersistentVolumeClaimConditionType: string

#PersistentVolumeClaimModifyVolumeStatus: string

#PersistentVolumeClaimPhase: string

#PersistentVolumeClaimRetentionPolicyType: string

#PersistentVolumeClaimSpec: {
	// accessModes contains the desired access modes the volume should have
	accessModes?: [...#PersistentVolumeAccessMode]
	// selector is a label query over volumes to consider for binding
	selector?: #LabelSelector
	// resources represents the minimum resources the volume should have
	resources?: #VolumeResourceRequirements
	// volumeName is the binding reference to the PersistentVolume backing this claim
	volumeName?: string
	// storageClassName is the name of the StorageClass required by the claim
	storageClassName?: string
	// volumeMode defines what type of volume is required by the claim
	volumeMode?: #PersistentVolumeMode
	// dataSource field can be used to specify either: * An existing VolumeSnapshot object (snapshot
	dataSource?: #TypedLocalObjectReference
	// dataSourceRef specifies the object from which to populate the volume with data, if a non-empty volume is desired
	dataSourceRef?: #TypedObjectReference
	// volumeAttributesClassName may be used to set the VolumeAttributesClass used by this claim
	volumeAttributesClassName?: string
}

#PersistentVolumeClaimStatus: {
	// phase represents the current phase of PersistentVolumeClaim
	phase?: #PersistentVolumeClaimPhase
	// accessModes contains the actual access modes the volume backing the PVC has
	accessModes?: [...#PersistentVolumeAccessMode]
	// capacity represents the actual resources of the underlying volume
	capacity?: #ResourceList
	// conditions is the current Condition of persistent volume claim
	conditions?: [...#PersistentVolumeClaimCondition]
	// allocatedResources tracks the resources allocated to a PVC including its capacity
	allocatedResources?: #ResourceList
	// allocatedResourceStatuses stores status of resource being resized for the given PVC
	allocatedResourceStatuses?: {
		[X=#ResourceName]: #ClaimResourceStatus
	}
	// currentVolumeAttributesClassName is the current name of the VolumeAttributesClass the PVC is using
	currentVolumeAttributesClassName?: string
	// ModifyVolumeStatus represents the status object of ControllerModifyVolume operation
	modifyVolumeStatus?: #ModifyVolumeStatus
}

#PersistentVolumeClaimVolumeSource: {
	// claimName is the name of a PersistentVolumeClaim in the same namespace as the pod using this volume
	claimName: string
	// else volumeMounts
	readOnly?: bool
}

#PersistentVolumeMode: string

#PodAffinity: {
	// If the affinity requirements specified by this field are not met at scheduling time, the pod will not be scheduled onto the node
	requiredDuringSchedulingIgnoredDuringExecution?: [...#PodAffinityTerm]
	// The scheduler will prefer to schedule pods to nodes that satisfy the affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions
	preferredDuringSchedulingIgnoredDuringExecution?: [...#WeightedPodAffinityTerm]
}

#PodAffinityTerm: {
	// A label query over a set of resources, in this case pods
	labelSelector?: #LabelSelector
	// namespaces specifies a static list of namespace names that the term applies to
	namespaces?: [...string]
	// This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching the labelSelector in the specified namespaces, where co-located is defined as running on a node whose value of the label with key topologyKey matches that of any node on which any of the selected pods is running
	topologyKey: string
	// A label query over the set of namespaces that the term applies to
	namespaceSelector?: #LabelSelector
	// MatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration
	matchLabelKeys?: [...string]
	// MismatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration
	mismatchLabelKeys?: [...string]
}

#PodAntiAffinity: {
	// If the anti-affinity requirements specified by this field are not met at scheduling time, the pod will not be scheduled onto the node
	requiredDuringSchedulingIgnoredDuringExecution?: [...#PodAffinityTerm]
	// The scheduler will prefer to schedule pods to nodes that satisfy the anti-affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions
	preferredDuringSchedulingIgnoredDuringExecution?: [...#WeightedPodAffinityTerm]
}

#PodConditionType: string

#PodDNSConfig: {
	// A list of DNS name server IP addresses
	nameservers?: [...string]
	// A list of DNS search domains for host-name lookup
	searches?: [...string]
	// A list of DNS resolver options
	options?: [...#PodDNSConfigOption]
}

#PodDNSConfigOption: {
	// Name is this DNS resolver option's name
	name?: string
	// Value is this DNS resolver option's value
	value?: string
}

#PodFSGroupChangePolicy: string

#PodFailurePolicy: {
	// A list of pod failure policy rules
	rules: [...#PodFailurePolicyRule]
}

#PodFailurePolicyAction: string

#PodFailurePolicyOnExitCodesOperator: string

#PodFailurePolicyOnExitCodesRequirement: {
	// Restricts the check for exit codes to the container with the specified name
	containerName?: string
	// Represents the relationship between the container exit code(s) and the specified values
	operator: #PodFailurePolicyOnExitCodesOperator
	// Specifies the set of values
	values: [...int]
}

#PodFailurePolicyOnPodConditionsPattern: {
	// Specifies the required Pod condition type
	type: #PodConditionType
	// Specifies the required Pod condition status
	status: #ConditionStatus
}

#PodFailurePolicyRule: {
	// Specifies the action taken on a pod failure when the requirements are satisfied
	action: #PodFailurePolicyAction
	// Represents the requirement on the container exit codes
	onExitCodes?: #PodFailurePolicyOnExitCodesRequirement
	// Represents the requirement on the pod conditions
	onPodConditions?: [...#PodFailurePolicyOnPodConditionsPattern]
}

#PodManagementPolicyType: string

#PodOS: {
	// Name is the name of the operating system
	name: #OSName
}

#PodPartialSpec: {
	// Restart policy for all containers within the pod.
	restartPolicy?: #RestartPolicy
	// Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request.
	terminationGracePeriodSeconds?: int
	// Optional duration in seconds the pod may be active on the node relative to
	activeDeadlineSeconds?: int
	// Set DNS policy for the pod.
	dnsPolicy?: #DNSPolicy
	// is a selector which must be true for the pod to fit on a node.
	nodeSelector?: {
		[X=string]: string
	}
	// indicates whether a service account token should be automatically mounted.
	automountServiceAccountToken?: bool
	// indicates in which node this pod is scheduled.
	nodeName?: string
	// Host networking requested for this pod. Use the host's network namespace.
	hostNetwork?: bool
	// Use the host's pid namespace.
	hostPID?: bool
	// Use the host's ipc namespace.
	hostIPC?: bool
	// Share a single process namespace between all of the containers in a pod.
	shareProcessNamespace?: bool
	// holds pod-level security attributes and common container settings.
	securityContext?: #PodSecurityContext
	// is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	imagePullSecrets?: [...#LocalObjectReference]
	// Specifies the hostname of the Pod
	hostname?: string
	// If specified, the fully qualified Pod hostname will be "<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>".
	subdomain?: string
	// If specified, the pod's scheduling constraints
	affinity?: #Affinity
	// If specified, the pod will be dispatched by specified scheduler.
	schedulerName?: string
	// If specified, the pod's tolerations.
	tolerations?: [...#Toleration]
	// is an optional list of hosts and IPs that will be injected into the pod's hosts
	hostAliases?: [...#HostAlias]
	// If specified, indicates the pod's priority. "system-node-critical" and
	priorityClassName?: string
	// The priority value. Various system components use this field to find the
	priority?: int
	// Specifies the DNS parameters of a pod.
	dnsConfig?: #PodDNSConfig
	// If specified, all readiness gates will be evaluated for pod readiness.
	readinessGates?: [...#PodReadinessGate]
	// refers to a RuntimeClass object in the node.k8s.io group, which should be used
	runtimeClassName?: string
	// indicates whether information about services should be injected into pod's
	enableServiceLinks?: bool
	// is the Policy for preempting pods with lower priority.
	preemptionPolicy?: #PreemptionPolicy
	// represents the resource overhead associated with running a pod for a given RuntimeClass.
	overhead?: #ResourceList
	// describes how a group of pods ought to spread across topology
	topologySpreadConstraints?: [...#TopologySpreadConstraint]
	// If true the pod's hostname will be configured as the pod's FQDN, rather than the leaf name (the default).
	setHostnameAsFQDN?: bool
	// Specifies the OS of the containers in the pod.
	os?: #PodOS
	// Use the host's user namespace.
	hostUsers?: bool
	// is an opaque list of values that if specified will block scheduling the pod.
	schedulingGates?: [...#PodSchedulingGate]
	// defines which ResourceClaims must be allocated
	resourceClaims?: [...#PodResourceClaim]
	// is the total amount of CPU and Memory resources required by all
	resources?: #ResourceRequirements
}

#PodPartialTemplateSpec: spec?: #PodPartialSpec

#PodReadinessGate: {
	// ConditionType refers to a condition in the pod's condition list with matching type
	conditionType: #PodConditionType
}

#PodReplacementPolicy: string

#PodResourceClaim: {
	// Name uniquely identifies this resource claim inside the pod
	name: string
	// ResourceClaimName is the name of a ResourceClaim object in the same namespace as this pod
	resourceClaimName?: string
	// ResourceClaimTemplateName is the name of a ResourceClaimTemplate object in the same namespace as this pod
	resourceClaimTemplateName?: string
}

#PodSELinuxChangePolicy: string

#PodSchedulingGate: {
	// Name of the scheduling gate
	name: string
}

#PodSecurityContext: {
	// The SELinux context to be applied to all containers
	seLinuxOptions?: #SELinuxOptions
	// The Windows specific settings applied to all containers
	windowsOptions?: #WindowsSecurityContextOptions
	// The UID to run the entrypoint of the container process
	runAsUser?: int
	// The GID to run the entrypoint of the container process
	runAsGroup?: int
	// Indicates that the container must run as a non-root user
	runAsNonRoot?: bool
	// A list of groups applied to the first process run in each container, in addition to the container's primary GID and fsGroup (if specified)
	supplementalGroups?: [...int]
	// Defines how supplemental groups of the first container processes are calculated
	supplementalGroupsPolicy?: #SupplementalGroupsPolicy
	// A special supplemental group that applies to all containers in a pod
	fsGroup?: int
	// Sysctls hold a list of namespaced sysctls used for the pod
	sysctls?: [...#Sysctl]
	// fsGroupChangePolicy defines behavior of changing ownership and permission of the volume before being exposed inside Pod
	fsGroupChangePolicy?: #PodFSGroupChangePolicy
	// The seccomp options to use by the containers in this pod
	seccompProfile?: #SeccompProfile
	// appArmorProfile is the AppArmor options to use by the containers in this pod
	appArmorProfile?: #AppArmorProfile
	// seLinuxChangePolicy defines how the container's SELinux label is applied to all volumes used by the Pod
	seLinuxChangePolicy?: #PodSELinuxChangePolicy
}

#PolicyRule: {
	// Verbs is a list of Verbs that apply to ALL the ResourceKinds contained in this rule
	verbs: [...string]
	// APIGroups is the name of the APIGroup that contains the resources
	apiGroups?: [...string]
	// Resources is a list of resources this rule applies to
	resources?: [...string]
	// ResourceNames is an optional white list of names that the rule applies to
	resourceNames?: [...string]
	// NonResourceURLs is a set of partial urls that a user should have access to
	nonResourceURLs?: [...string]
}

#PreemptionPolicy: string

#PreferredSchedulingTerm: {
	// Weight associated with matching the corresponding nodeSelectorTerm, in the range 1-100
	weight: int
	// A node selector term, associated with the corresponding weight
	preference: #NodeSelectorTerm
}

#Probe: {
	exec?:      #ExecAction
	httpGet?:   #HTTPGetAction
	tcpSocket?: #TCPSocketAction
	grpc?:      #GRPCAction
	// Number of seconds after the container has started before liveness probes are initiated
	initialDelaySeconds?: int
	// Number of seconds after which the probe times out
	timeoutSeconds?: int
	// How often (in seconds) to perform the probe
	periodSeconds?: int
	// Minimum consecutive successes for the probe to be considered successful after having failed
	successThreshold?: int
	// Minimum consecutive failures for the probe to be considered failed after having succeeded
	failureThreshold?: int
	// Optional duration in seconds the pod needs to terminate gracefully upon probe failure
	terminationGracePeriodSeconds?: int
}

#ProcMountType: string

#PullPolicy: string

#Quantity: string

#ResourceClaim: {
	// Name must match the name of one entry in pod
	name: string
	// Request is the name chosen for a request in the referenced claim
	request?: string
}

#ResourceList: [X=#ResourceName]: #Quantity

#ResourceName: string

#ResourceRequirements: {
	// Limits describes the maximum amount of compute resources allowed
	limits?: #ResourceList
	// Requests describes the minimum amount of compute resources required
	requests?: #ResourceList
	// Claims lists the names of resources, defined in spec
	claims?: [...#ResourceClaim]
}

#RestartPolicy: string

#RollingUpdateDaemonSet: {
	// The maximum number of DaemonSet pods that can be unavailable during the update
	maxUnavailable?: #IntOrString
	// The maximum number of nodes with an existing available DaemonSet pod that can have an updated DaemonSet pod during during an update
	maxSurge?: #IntOrString
}

#RollingUpdateDeployment: {
	// The maximum number of pods that can be unavailable during the update
	maxUnavailable?: #IntOrString
	// The maximum number of pods that can be scheduled above the desired number of pods
	maxSurge?: #IntOrString
}

#RollingUpdateStatefulSetStrategy: {
	// Partition indicates the ordinal at which the StatefulSet should be partitioned for updates
	partition?: int
	// The maximum number of pods that can be unavailable during the update
	maxUnavailable?: #IntOrString
}

#SELinuxOptions: {
	// User is a SELinux user label that applies to the container
	user?: string
	// Role is a SELinux role label that applies to the container
	role?: string
	// Type is a SELinux type label that applies to the container
	type?: string
	// Level is SELinux level label that applies to the container
	level?: string
}

#ScopeType: "Cluster" | "Namespace"

#SeccompProfile: {
	// type indicates which kind of seccomp profile will be applied
	type: #SeccompProfileType
	// localhostProfile indicates a profile defined in a file on the node should be used
	localhostProfile?: string
}

#SeccompProfileType: string

#SecretVolumeSource: {
	// secretName is the name of the secret in the pod's namespace to use
	secretName?: string
	// items If unspecified, each key-value pair in the Data field of the referenced Secret will be projected into the volume as a file whose name is the key and content is the value
	items?: [...#KeyToPath]
	// defaultMode is Optional: mode bits used to set permissions on created files by default
	defaultMode?: int
	// optional field specify whether the Secret or its keys must be defined
	optional?: bool
}

#SecurityContext: {
	// The capabilities to add/drop when running containers
	capabilities?: #Capabilities
	// Run container in privileged mode
	privileged?: bool
	// The SELinux context to be applied to the container
	seLinuxOptions?: #SELinuxOptions
	// The Windows specific settings applied to all containers
	windowsOptions?: #WindowsSecurityContextOptions
	// The UID to run the entrypoint of the container process
	runAsUser?: int
	// The GID to run the entrypoint of the container process
	runAsGroup?: int
	// Indicates that the container must run as a non-root user
	runAsNonRoot?: bool
	// Whether this container has a read-only root filesystem
	readOnlyRootFilesystem?: bool
	// AllowPrivilegeEscalation controls whether a process can gain more privileges than its parent process
	allowPrivilegeEscalation?: bool
	// procMount denotes the type of proc mount to use for the containers
	procMount?: #ProcMountType
	// The seccomp options to use by this container
	seccompProfile?: #SeccompProfile
	// appArmorProfile is the AppArmor options to use by this container
	appArmorProfile?: #AppArmorProfile
}

#Service: {
	// [PortName]servicePort
	ports?: {
		[X=string]: int
	}
	// [PortName]BashPath
	paths?: {
		[X=string]: string
	}
	clusterIP?: string
	expose?:    #Expose
}

#ServiceAccount: {
	scope?: #ScopeType
	rules: [...#PolicyRule]
}

#SleepAction: {
	// Seconds is the number of seconds to sleep
	seconds: int
}

#Spec: {
	version: string
	deploy?: #Deploy
	config?: [X=string]: #EnvVarValueOrFrom
	containers?: [X=string]: #Container
	volumes?: [X=string]: #Volume
	services?: [X=string]: #Service
	serviceAccount?: #ServiceAccount
	manifests?: [X=string]: _
	images?: [X=string]: #Image
}

#StatefulSetOrdinals: {
	// start is the number representing the first replica's index
	start: int
}

#StatefulSetPersistentVolumeClaimRetentionPolicy: {
	// WhenDeleted specifies what happens to PVCs created from StatefulSet VolumeClaimTemplates when the StatefulSet is deleted
	whenDeleted?: #PersistentVolumeClaimRetentionPolicyType
	// WhenScaled specifies what happens to PVCs created from StatefulSet VolumeClaimTemplates when the StatefulSet is scaled down
	whenScaled?: #PersistentVolumeClaimRetentionPolicyType
}

#StatefulSetSpec: {
	// replicas is the desired number of replicas of the given Template.
	replicas?: int
	// template is the object that describes the pod that will be created if
	template?: #PodPartialTemplateSpec
	// volumeClaimTemplates is a list of claims that pods are allowed to reference.
	volumeClaimTemplates?: [...#PersistentVolumeClaim]
	// serviceName is the name of the service that governs this StatefulSet.
	serviceName: string
	// podManagementPolicy controls how pods are created during initial scale up,
	podManagementPolicy?: #PodManagementPolicyType
	// updateStrategy indicates the StatefulSetUpdateStrategy that will be
	updateStrategy?: #StatefulSetUpdateStrategy
	// revisionHistoryLimit is the maximum number of revisions that will
	revisionHistoryLimit?: int
	// Minimum number of seconds for which a newly created pod should be ready
	minReadySeconds?: int
	// persistentVolumeClaimRetentionPolicy describes the lifecycle of persistent
	persistentVolumeClaimRetentionPolicy?: #StatefulSetPersistentVolumeClaimRetentionPolicy
	// ordinals controls the numbering of replica indices in a StatefulSet. The
	ordinals?: #StatefulSetOrdinals
}

#StatefulSetUpdateStrategy: {
	// Type indicates the type of the StatefulSetUpdateStrategy
	type?: #StatefulSetUpdateStrategyType
	// RollingUpdate is used to communicate parameters when Type is RollingUpdateStatefulSetStrategyType
	rollingUpdate?: #RollingUpdateStatefulSetStrategy
}

#StatefulSetUpdateStrategyType: string

#Status: {
	endpoint?: [X=string]: string
	resources?: [...{
		[X=string]: _
	}]
	images?: [X=string]: string
	digests?: [...#DigestMeta]
}

#StorageMedium: string

#SuccessPolicy: {
	// rules represents the list of alternative rules for the declaring the Jobs as successful before `
	rules: [...#SuccessPolicyRule]
}

#SuccessPolicyRule: {
	// succeededIndexes specifies the set of indexes which need to be contained in the actual set of the succeeded indexes for the Job
	succeededIndexes?: string
	// succeededCount specifies the minimal required size of the actual set of the succeeded indexes for the Job
	succeededCount?: int
}

#SupplementalGroupsPolicy: string

#Sysctl: {
	// Name of a property to set
	name: string
	// Value of a property to set
	value: string
}

#TCPSocketAction: {
	// Number or name of the port to access on the container
	port: #IntOrString
	// Optional: Host name to connect to, defaults to the pod IP
	host?: string
}

#TaintEffect: string

#TerminationMessagePolicy: string

#Time: string

#Toleration: {
	// Key is the taint key that the toleration applies to
	key?: string
	// Operator represents a key's relationship to the value
	operator?: #TolerationOperator
	// Value is the taint value the toleration matches to
	value?: string
	// Effect indicates the taint effect to match
	effect?: #TaintEffect
	// TolerationSeconds represents the period of time the toleration (which must be of effect NoExecute, otherwise this field is ignored) tolerates the taint
	tolerationSeconds?: int
}

#TolerationOperator: string

#TopologySpreadConstraint: {
	// MaxSkew describes the degree to which pods may be unevenly distributed
	maxSkew: int
	// TopologyKey is the key of node labels
	topologyKey: string
	// WhenUnsatisfiable indicates how to deal with a pod if it doesn't satisfy the spread constraint
	whenUnsatisfiable: #UnsatisfiableConstraintAction
	// LabelSelector is used to find matching pods
	labelSelector?: #LabelSelector
	// MinDomains indicates a minimum number of eligible domains
	minDomains?: int
	// NodeAffinityPolicy indicates how we will treat Pod's nodeAffinity/nodeSelector when calculating pod topology spread skew
	nodeAffinityPolicy?: #NodeInclusionPolicy
	// NodeTaintsPolicy indicates how we will treat node taints when calculating pod topology spread skew
	nodeTaintsPolicy?: #NodeInclusionPolicy
	// MatchLabelKeys is a set of pod label keys to select the pods over which spreading will be calculated
	matchLabelKeys?: [...string]
}

#TypedLocalObjectReference: {
	// APIGroup is the group for the resource being referenced
	apiGroup: string
	// Kind is the type of resource being referenced
	kind: string
	// Name is the name of resource being referenced
	name: string
}

#TypedObjectReference: {
	// APIGroup is the group for the resource being referenced
	apiGroup: string
	// Kind is the type of resource being referenced
	kind: string
	// Name is the name of resource being referenced
	name: string
	// Namespace is the namespace of resource being referenced Note that when a namespace is specified, a gateway
	namespace?: string
}

#UID: string

#URIScheme: string

#UnsatisfiableConstraintAction: string

#Volume: #VolumeConfigMap | #VolumeEmptyDir | #VolumeHostPath | #VolumeImage | #VolumePersistentVolumeClaim | #VolumeSecret

#VolumeConfigMap: {
	type:              "ConfigMap"
	opt?:              #ConfigMapVolumeSource
	spec?:             #ConfigMapSpec
	mountPath:         string
	mountPropagation?: "Bidirectional" | "HostToContainer"
	// mountPath == export, use as envFrom
	prefix?:   string
	optional?: bool
	// else volumeMounts
	readOnly?: bool
	subPath?:  string
}

#VolumeEmptyDir: {
	type:              "EmptyDir"
	opt?:              #EmptyDirVolumeSource
	mountPath:         string
	mountPropagation?: "Bidirectional" | "HostToContainer"
	// mountPath == export, use as envFrom
	prefix?:   string
	optional?: bool
	// else volumeMounts
	readOnly?: bool
	subPath?:  string
}

#VolumeHostPath: {
	type:              "HostPath"
	opt?:              #HostPathVolumeSource
	mountPath:         string
	mountPropagation?: "Bidirectional" | "HostToContainer"
	// mountPath == export, use as envFrom
	prefix?:   string
	optional?: bool
	// else volumeMounts
	readOnly?: bool
	subPath?:  string
}

#VolumeImage: {
	type:              "Image"
	opt?:              #ImageVolumeSource
	mountPath:         string
	mountPropagation?: "Bidirectional" | "HostToContainer"
	// mountPath == export, use as envFrom
	prefix?:   string
	optional?: bool
	// else volumeMounts
	readOnly?: bool
	subPath?:  string
}

#VolumePersistentVolumeClaim: {
	type:              "PersistentVolumeClaim"
	opt?:              #PersistentVolumeClaimVolumeSource
	spec:              #PersistentVolumeClaimSpec
	mountPath:         string
	mountPropagation?: "Bidirectional" | "HostToContainer"
	// mountPath == export, use as envFrom
	prefix?:   string
	optional?: bool
	// else volumeMounts
	readOnly?: bool
	subPath?:  string
}

#VolumeResourceRequirements: {
	// Limits describes the maximum amount of compute resources allowed
	limits?: #ResourceList
	// Requests describes the minimum amount of compute resources required
	requests?: #ResourceList
}

#VolumeSecret: {
	type:              "Secret"
	opt?:              #SecretVolumeSource
	spec?:             #ConfigMapSpec
	mountPath:         string
	mountPropagation?: "Bidirectional" | "HostToContainer"
	// mountPath == export, use as envFrom
	prefix?:   string
	optional?: bool
	// else volumeMounts
	readOnly?: bool
	subPath?:  string
}

#WeightedPodAffinityTerm: {
	// weight associated with matching the corresponding podAffinityTerm, in the range 1-100
	weight: int
	// Required
	podAffinityTerm: #PodAffinityTerm
}

#WindowsSecurityContextOptions: {
	// GMSACredentialSpecName is the name of the GMSA credential spec to use
	gmsaCredentialSpecName?: string
	// GMSACredentialSpec is where the GMSA admission webhook (https://github
	gmsaCredentialSpec?: string
	// The UserName in Windows to run the entrypoint of the container process
	runAsUserName?: string
	// HostProcess determines if a container should be run as a 'Host Process' container
	hostProcess?: bool
}
