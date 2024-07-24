package kubepkg

#Affinity: {
	// Describes node affinity scheduling rules for the pod.
	nodeAffinity?: #NodeAffinity
	// Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone, etc. as some other pod(s)).
	podAffinity?: #PodAffinity
	// Describes pod anti-affinity scheduling rules (e.g. avoid putting this pod in the same node, zone, etc. as some other pod(s)).
	podAntiAffinity?: #PodAntiAffinity
}

#AppArmorProfile: {
	// localhostProfile indicates a profile loaded on the node that should be used. The profile must be preconfigured on the node to work. Must match the loaded name of the profile. Must be set if and only if type is "Localhost".
	localhostProfile?: string
	// type indicates which kind of AppArmor profile will be applied. Valid options are:
	//   Localhost - a profile pre-loaded on the node.
	//   RuntimeDefault - the container runtime's default profile.
	//   Unconfined - no AppArmor enforcement.
	type: #AppArmorProfileType
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

#ClaimSource: {
	// ResourceClaimName is the name of a ResourceClaim object in the same namespace as this pod.
	resourceClaimName?: string
	// ResourceClaimTemplateName is the name of a ResourceClaimTemplate object in the same namespace as this pod.
	// 
	// The template will be used to create a new ResourceClaim, which will be bound to this pod. When this pod is deleted, the ResourceClaim will also be deleted. The pod name and resource name, along with a generated component, will be used to form a unique name for the ResourceClaim, which will be recorded in pod.status.resourceClaimStatuses.
	// 
	// This field is immutable and no changes will be made to the corresponding ResourceClaim by the control plane after creating the ResourceClaim.
	resourceClaimTemplateName?: string
}

#CompletionMode: string

#ConcurrencyPolicy: string

#ConditionStatus: string

#ConfigMapSpec: data: [X=string]: string

#ConfigMapVolumeSource: {
	// defaultMode is optional: mode bits used to set permissions on created files by default. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. Defaults to 0644. Directories within the path are not affected by this setting. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
	defaultMode?: int
	// items if unspecified, each key-value pair in the Data field of the referenced ConfigMap will be projected into the volume as a file whose name is the key and content is the value. If specified, the listed keys will be projected into the specified paths, and unlisted keys will not be present. If a key is specified which is not present in the ConfigMap, the volume setup will error unless it is marked optional. Paths must be relative and may not contain the '..' path or start with '..'.
	items?: [...#KeyToPath]
	name?:     string
	optional?: bool
}

#Container: {
	args?: [...string]
	command?: [...string]
	env?: [X=string]: #EnvVarValueOrFrom
	image:          #Image
	lifecycle?:     #Lifecycle
	livenessProbe?: #Probe
	// Ports: [PortName]: ContainerPort
	ports?: {
		[X=string]: int
	}
	readinessProbe?:           #Probe
	resources?:                #ResourceRequirements
	securityContext?:          #SecurityContext
	startupProbe?:             #Probe
	stdin?:                    bool
	stdinOnce?:                bool
	terminationMessagePath?:   string
	terminationMessagePolicy?: #TerminationMessagePolicy
	tty?:                      bool
	workingDir?:               string
}

#CronJobSpec: {
	// Specifies how to treat concurrent executions of a Job.
	// Valid values are:
	// 
	// - "Allow" (default): allows CronJobs to run concurrently;
	// - "Forbid": forbids concurrent runs, skipping next run if previous run hasn't finished yet;
	// - "Replace": cancels currently running job and replaces it with a new one
	concurrencyPolicy?: #ConcurrencyPolicy
	// The number of failed finished jobs to retain. Value must be non-negative integer.
	// Defaults to 1.
	failedJobsHistoryLimit?: int
	// Specifies the job that will be created when executing a CronJob.
	jobTemplate: #JobTemplateSpec
	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	schedule: string
	// Optional deadline in seconds for starting the job if it misses scheduled
	// time for any reason.  Missed jobs executions will be counted as failed ones.
	startingDeadlineSeconds?: int
	// The number of successful finished jobs to retain. Value must be non-negative integer.
	// Defaults to 3.
	successfulJobsHistoryLimit?: int
	// This flag tells the controller to suspend subsequent executions, it does
	// not apply to already started executions.  Defaults to false.
	suspend?: bool
	// The time zone name for the given schedule, see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones.
	// If not specified, this will default to the time zone of the kube-controller-manager process.
	// The set of valid time zone names and the time zone offset is loaded from the system-wide time zone
	// database by the API server during CronJob validation and the controller manager during execution.
	// If no system-wide time zone database can be found a bundled version of the database is used instead.
	// If the time zone name becomes invalid during the lifetime of a CronJob or due to a change in host
	// configuration, the controller will stop creating new new Jobs and will create a system event with the
	// reason UnknownTimeZone.
	// More information can be found in https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/#time-zones
	timeZone?: string
}

#DNSPolicy: string

#DaemonSetSpec: {
	// The minimum number of seconds for which a newly created DaemonSet pod should
	// be ready without any of its container crashing, for it to be considered
	// available. Defaults to 0 (pod will be considered available as soon as it
	// is ready).
	minReadySeconds?: int
	// The number of old history to retain to allow rollback.
	// This is a pointer to distinguish between explicit zero and not specified.
	// Defaults to 10.
	revisionHistoryLimit?: int
	// An object that describes the pod that will be created.
	// The DaemonSet will create exactly one copy of this pod on every node
	// that matches the template's node selector (or on every node if no node
	// selector is specified).
	// The only allowed template.spec.restartPolicy value is "Always".
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#pod-template
	template: #PodPartialTemplateSpec
	// An update strategy to replace existing DaemonSet pods with new pods.
	updateStrategy?: #DaemonSetUpdateStrategy
}

#DaemonSetUpdateStrategy: {
	// Rolling update config params. Present only if type = "RollingUpdate".
	rollingUpdate?: #RollingUpdateDaemonSet
	// Type of daemon set update. Can be "RollingUpdate" or "OnDelete". Default is RollingUpdate.
	type?: #DaemonSetUpdateStrategyType
}

#DaemonSetUpdateStrategyType: "RollingUpdate" | "OnDelete"

#Deploy: #DeployConfigMap | #DeployCronJob | #DeployDaemonSet | #DeployDeployment | #DeployJob | #DeploySecret | #DeployStatefulSet

#DeployConfigMap: {
	annotations?: [X=string]: string
	kind: "ConfigMap"
}

#DeployCronJob: {
	annotations?: [X=string]: string
	kind:  "CronJob"
	spec?: #CronJobSpec
}

#DeployDaemonSet: {
	annotations?: [X=string]: string
	kind:  "DaemonSet"
	spec?: #DaemonSetSpec
}

#DeployDeployment: {
	annotations?: [X=string]: string
	kind:  "Deployment"
	spec?: #DeploymentSpec
}

#DeployJob: {
	annotations?: [X=string]: string
	kind:  "Job"
	spec?: #JobSpec
}

#DeploySecret: {
	annotations?: [X=string]: string
	kind: "Secret"
}

#DeployStatefulSet: {
	annotations?: [X=string]: string
	kind:  "StatefulSet"
	spec?: #StatefulSetSpec
}

#DeploymentSpec: {
	// Minimum number of seconds for which a newly created pod should be ready
	// without any of its container crashing, for it to be considered available.
	// Defaults to 0 (pod will be considered available as soon as it is ready)
	minReadySeconds?: int
	// Indicates that the deployment is paused.
	paused?: bool
	// The maximum time in seconds for a deployment to make progress before it
	// is considered to be failed. The deployment controller will continue to
	// process failed deployments and a condition with a ProgressDeadlineExceeded
	// reason will be surfaced in the deployment status. Note that progress will
	// not be estimated during the time a deployment is paused. Defaults to 600s.
	progressDeadlineSeconds?: int
	// Number of desired pods. This is a pointer to distinguish between explicit
	// zero and not specified. Defaults to 1.
	replicas?: int
	// The number of old ReplicaSets to retain to allow rollback.
	// This is a pointer to distinguish between explicit zero and not specified.
	// Defaults to 10.
	revisionHistoryLimit?: int
	// The deployment strategy to use to replace existing pods with new ones.
	strategy?: #DeploymentStrategy
	// Template describes the pods that will be created.
	// The only allowed template.spec.restartPolicy value is "Always".
	template: #PodPartialTemplateSpec
}

#DeploymentStrategy: {
	// Rolling update config params. Present only if DeploymentStrategyType = RollingUpdate.
	rollingUpdate?: #RollingUpdateDeployment
	// Type of deployment. Can be "Recreate" or "RollingUpdate". Default is RollingUpdate.
	type?: #DeploymentStrategyType
}

#DeploymentStrategyType: "Recreate" | "RollingUpdate"

#DigestMeta: {
	digest:    string
	name:      string
	platform?: string
	size:      #FileSize
	tag?:      string
	type:      #DigestMetaType
}

#DigestMetaType: "blob" | "manifest"

#EmptyDirVolumeSource: {
	// medium represents what type of storage medium should back this directory. The default is "" which means to use the node's default medium. Must be an empty string (default) or Memory. More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir
	medium?: #StorageMedium
	// sizeLimit is the total amount of local storage required for this EmptyDir volume. The size limit is also applicable for memory medium. The maximum usage on memory medium EmptyDir would be the minimum value between the SizeLimit specified here and the sum of memory limits of all containers in a pod. The default is nil which means that the limit is undefined. More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir
	sizeLimit?: #Quantity
}

#EnvVarValueOrFrom: string

#ExecAction: command?: [...string]

#Expose: {
	gateway?: [...string]
	// Type NodePort | Ingress
	type: string
}

#FieldsV1: {}

#FileSize: int

#GRPCAction: {
	// Port number of the gRPC service. Number must be in the range 1 to 65535.
	port: int
	// Service is the name of the service to place in the gRPC HealthCheckRequest (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
	// 
	// If this is not specified, the default behavior is defined by gRPC.
	service: string
}

#HTTPGetAction: {
	// Host name to connect to, defaults to the pod IP. You probably want to set "Host" in httpHeaders instead.
	host?: string
	// Custom headers to set in the request. HTTP allows repeated headers.
	httpHeaders?: [...#HTTPHeader]
	// Path to access on the HTTP server.
	path?: string
	// Name or number of the port to access on the container. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.
	port: #IntOrString
	// Scheme to use for connecting to the host. Defaults to HTTP.
	scheme?: #URIScheme
}

#HTTPHeader: {
	// The header field name. This will be canonicalized upon output, so case-variant names will be understood as the same header.
	name: string
	// The header field value
	value: string
}

#HostAlias: {
	// Hostnames for the above IP address.
	hostnames?: [...string]
	// IP address of the host file entry.
	ip: string
}

#HostPathType: string

#HostPathVolumeSource: {
	// path of the directory on the host. If the path is a symlink, it will follow the link to the real path. More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
	path:  string
	type?: #HostPathType
}

#Image: {
	digest?: string
	name:    string
	platforms?: [...string]
	pullPolicy?: #PullPolicy
	tag?:        string
}

#IntOrString: int | string

#JobSpec: {
	// Specifies the duration in seconds relative to the startTime that the job
	// may be continuously active before the system tries to terminate it; value
	// must be positive integer. If a Job is suspended (at creation or through an
	// update), this timer will effectively be stopped and reset when the Job is
	// resumed again.
	activeDeadlineSeconds?: int
	// Specifies the number of retries before marking this job failed.
	// Defaults to 6
	backoffLimit?: int
	// Specifies the limit for the number of retries within an
	// index before marking this index as failed. When enabled the number of
	// failures per index is kept in the pod's
	// batch.kubernetes.io/job-index-failure-count annotation. It can only
	// be set when Job's completionMode=Indexed, and the Pod's restart
	// policy is Never. The field is immutable.
	// This field is beta-level. It can be used when the `JobBackoffLimitPerIndex`
	// feature gate is enabled (enabled by default).
	backoffLimitPerIndex?: int
	// completionMode specifies how Pod completions are tracked. It can be
	// `NonIndexed` (default) or `Indexed`.
	// 
	// `NonIndexed` means that the Job is considered complete when there have
	// been .spec.completions successfully completed Pods. Each Pod completion is
	// homologous to each other.
	// 
	// `Indexed` means that the Pods of a
	// Job get an associated completion index from 0 to (.spec.completions - 1),
	// available in the annotation batch.kubernetes.io/job-completion-index.
	// The Job is considered complete when there is one successfully completed Pod
	// for each index.
	// When value is `Indexed`, .spec.completions must be specified and
	// `.spec.parallelism` must be less than or equal to 10^5.
	// In addition, The Pod name takes the form
	// `$(job-name)-$(index)-$(random-string)`,
	// the Pod hostname takes the form `$(job-name)-$(index)`.
	// 
	// More completion modes can be added in the future.
	// If the Job controller observes a mode that it doesn't recognize, which
	// is possible during upgrades due to version skew, the controller
	// skips updates for the Job.
	completionMode?: #CompletionMode
	// Specifies the desired number of successfully finished pods the
	// job should be run with.  Setting to null means that the success of any
	// pod signals the success of all pods, and allows parallelism to have any positive
	// value.  Setting to 1 means that parallelism is limited to 1 and the success of that
	// pod signals the success of the job.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	completions?: int
	// ManagedBy field indicates the controller that manages a Job. The k8s Job
	// controller reconciles jobs which don't have this field at all or the field
	// value is the reserved string `kubernetes.io/job-controller`, but skips
	// reconciling Jobs with a custom value for this field.
	// The value must be a valid domain-prefixed path (e.g. acme.io/foo) -
	// all characters before the first "/" must be a valid subdomain as defined
	// by RFC 1123. All characters trailing the first "/" must be valid HTTP Path
	// characters as defined by RFC 3986. The value cannot exceed 64 characters.
	// 
	// This field is alpha-level. The job controller accepts setting the field
	// when the feature gate JobManagedBy is enabled (disabled by default).
	managedBy?: string
	// manualSelector controls generation of pod labels and pod selectors.
	// Leave `manualSelector` unset unless you are certain what you are doing.
	// When false or unset, the system pick labels unique to this job
	// and appends those labels to the pod template.  When true,
	// the user is responsible for picking unique labels and specifying
	// the selector.  Failure to pick a unique label may cause this
	// and other jobs to not function correctly.  However, You may see
	// `manualSelector=true` in jobs that were created with the old `extensions/v1beta1`
	// API.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/#specifying-your-own-pod-selector
	manualSelector?: bool
	// Specifies the maximal number of failed indexes before marking the Job as
	// failed, when backoffLimitPerIndex is set. Once the number of failed
	// indexes exceeds this number the entire Job is marked as Failed and its
	// execution is terminated. When left as null the job continues execution of
	// all of its indexes and is marked with the `Complete` Job condition.
	// It can only be specified when backoffLimitPerIndex is set.
	// It can be null or up to completions. It is required and must be
	// less than or equal to 10^4 when is completions greater than 10^5.
	// This field is beta-level. It can be used when the `JobBackoffLimitPerIndex`
	// feature gate is enabled (enabled by default).
	maxFailedIndexes?: int
	// Specifies the maximum desired number of pods the job should
	// run at any given time. The actual number of pods running in steady state will
	// be less than this number when ((.spec.completions - .status.successful) < .spec.parallelism),
	// i.e. when the work left to do is less than max parallelism.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	parallelism?: int
	// Specifies the policy of handling failed pods. In particular, it allows to
	// specify the set of actions and conditions which need to be
	// satisfied to take the associated action.
	// If empty, the default behaviour applies - the counter of failed pods,
	// represented by the jobs's .status.failed field, is incremented and it is
	// checked against the backoffLimit. This field cannot be used in combination
	// with restartPolicy=OnFailure.
	// 
	// This field is beta-level. It can be used when the `JobPodFailurePolicy`
	// feature gate is enabled (enabled by default).
	podFailurePolicy?: #PodFailurePolicy
	// podReplacementPolicy specifies when to create replacement Pods.
	// Possible values are:
	// - TerminatingOrFailed means that we recreate pods
	// when they are terminating (has a metadata.deletionTimestamp) or failed.
	// - Failed means to wait until a previously created Pod is fully terminated (has phase
	// Failed or Succeeded) before creating a replacement Pod.
	// 
	// When using podFailurePolicy, Failed is the the only allowed value.
	// TerminatingOrFailed and Failed are allowed values when podFailurePolicy is not in use.
	// This is an beta field. To use this, enable the JobPodReplacementPolicy feature toggle.
	// This is on by default.
	podReplacementPolicy?: #PodReplacementPolicy
	// successPolicy specifies the policy when the Job can be declared as succeeded.
	// If empty, the default behavior applies - the Job is declared as succeeded
	// only when the number of succeeded pods equals to the completions.
	// When the field is specified, it must be immutable and works only for the Indexed Jobs.
	// Once the Job meets the SuccessPolicy, the lingering pods are terminated.
	// 
	// This field  is alpha-level. To use this field, you must enable the
	// `JobSuccessPolicy` feature gate (disabled by default).
	successPolicy?: #SuccessPolicy
	// suspend specifies whether the Job controller should create Pods or not. If
	// a Job is created with suspend set to true, no Pods are created by the Job
	// controller. If a Job is suspended after creation (i.e. the flag goes from
	// false to true), the Job controller will delete all active Pods associated
	// with this Job. Users must design their workload to gracefully handle this.
	// Suspending a Job will reset the StartTime field of the Job, effectively
	// resetting the ActiveDeadlineSeconds timer too. Defaults to false.
	suspend?: bool
	// Describes the pod that will be created when executing a job.
	// The only allowed template.spec.restartPolicy values are "Never" or "OnFailure".
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/jobs-run-to-completion/
	template: #PodPartialTemplateSpec
	// ttlSecondsAfterFinished limits the lifetime of a Job that has finished
	// execution (either Complete or Failed). If this field is set,
	// ttlSecondsAfterFinished after the Job finishes, it is eligible to be
	// automatically deleted. When the Job is being deleted, its lifecycle
	// guarantees (e.g. finalizers) will be honored. If this field is unset,
	// the Job won't be automatically deleted. If this field is set to zero,
	// the Job becomes eligible to be deleted immediately after it finishes.
	ttlSecondsAfterFinished?: int
}

#JobTemplateSpec: spec?: #JobSpec

#KeyToPath: {
	// key is the key to project.
	key: string
	// mode is Optional: mode bits used to set permissions on this file. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. If not specified, the volume defaultMode will be used. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
	mode?: int
	// path is the relative path of the file to map the key to. May not be an absolute path. May not contain the path element '..'. May not start with the string '..'.
	path: string
}

#KubePkg: {
	apiVersion?: string
	kind?:       string
	metadata?:   #ObjectMeta
	spec:        #Spec
	status?:     #Status
}

#LabelSelector: {
	// matchExpressions is a list of label selector requirements. The requirements are ANDed.
	matchExpressions?: [...#LabelSelectorRequirement]
	// matchLabels is a map of {key,value} pairs. A single {key,value} in the matchLabels map is equivalent to an element of matchExpressions, whose key field is "key", the operator is "In", and the values array contains only "value". The requirements are ANDed.
	matchLabels?: {
		[X=string]: string
	}
}

#LabelSelectorOperator: string

#LabelSelectorRequirement: {
	// key is the label key that the selector applies to.
	key: string
	// operator represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists and DoesNotExist.
	operator: #LabelSelectorOperator
	// values is an array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. This array is replaced during a strategic merge patch.
	values?: [...string]
}

#Lifecycle: {
	// PostStart is called immediately after a container is created. If the handler fails, the container is terminated and restarted according to its restart policy. Other management of the container blocks until the hook completes. More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
	postStart?: #LifecycleHandler
	// PreStop is called immediately before a container is terminated due to an API request or management event such as liveness/startup probe failure, preemption, resource contention, etc. The handler is not called if the container crashes or exits. The Pod's termination grace period countdown begins before the PreStop hook is executed. Regardless of the outcome of the handler, the container will eventually terminate within the Pod's termination grace period (unless delayed by finalizers). Other management of the container blocks until the hook completes or until the termination grace period is reached. More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
	preStop?: #LifecycleHandler
}

#LifecycleHandler: {
	// Exec specifies the action to take.
	exec?: #ExecAction
	// HTTPGet specifies the http request to perform.
	httpGet?: #HTTPGetAction
	// Sleep represents the duration that the container should sleep before being terminated.
	sleep?: #SleepAction
	// Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept for the backward compatibility. There are no validation of this field and lifecycle hooks will fail in runtime when tcp handler is specified.
	tcpSocket?: #TCPSocketAction
}

#LocalObjectReference: {
	// Name of the referent. This field is effectively required, but due to backwards compatibility is allowed to be empty. Instances of this type with an empty value here are almost certainly wrong. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	name?: string
}

#ManagedFieldsEntry: {
	// APIVersion defines the version of this resource that this field set applies to. The format is "group/version" just like the top-level APIVersion field. It is necessary to track the version of a field set because it cannot be automatically converted.
	apiVersion?: string
	// FieldsType is the discriminator for the different fields format and version. There is currently only one possible value: "FieldsV1"
	fieldsType?: string
	// FieldsV1 holds the first JSON version format as described in the "FieldsV1" type.
	fieldsV1?: #FieldsV1
	// Manager is an identifier of the workflow managing these fields.
	manager?: string
	// Operation is the type of operation which lead to this ManagedFieldsEntry being created. The only valid values for this field are 'Apply' and 'Update'.
	operation?: #ManagedFieldsOperationType
	// Subresource is the name of the subresource used to update that object, or empty string if the object was updated through the main resource. The value of this field is used to distinguish between managers, even if they share the same name. For example, a status update will be distinct from a regular update using the same manager name. Note that the APIVersion field is not related to the Subresource field and it always corresponds to the version of the main resource.
	subresource?: string
	// Time is the timestamp of when the ManagedFields entry was added. The timestamp will also be updated if a field is added, the manager changes any of the owned fields value or removes a field. The timestamp does not update when a field is removed from the entry because another manager took it over.
	time?: #Time
}

#ManagedFieldsOperationType: string

#ModifyVolumeStatus: {
	// status is the status of the ControllerModifyVolume operation. It can be in any of following states:
	//  - Pending
	//    Pending indicates that the PersistentVolumeClaim cannot be modified due to unmet requirements, such as
	//    the specified VolumeAttributesClass not existing.
	//  - InProgress
	//    InProgress indicates that the volume is being modified.
	//  - Infeasible
	//   Infeasible indicates that the request has been rejected as invalid by the CSI driver. To
	// 	  resolve the error, a valid VolumeAttributesClass needs to be specified.
	// Note: New statuses can be added in the future. Consumers should check for unknown statuses and fail appropriately.
	status: #PersistentVolumeClaimModifyVolumeStatus
	// targetVolumeAttributesClassName is the name of the VolumeAttributesClass the PVC currently being reconciled
	targetVolumeAttributesClassName?: string
}

#NodeAffinity: {
	// The scheduler will prefer to schedule pods to nodes that satisfy the affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions. The node that is most preferred is the one with the greatest sum of weights, i.e. for each node that meets all of the scheduling requirements (resource request, requiredDuringScheduling affinity expressions, etc.), compute a sum by iterating through the elements of this field and adding "weight" to the sum if the node matches the corresponding matchExpressions; the node(s) with the highest sum are the most preferred.
	preferredDuringSchedulingIgnoredDuringExecution?: [...#PreferredSchedulingTerm]
	// If the affinity requirements specified by this field are not met at scheduling time, the pod will not be scheduled onto the node. If the affinity requirements specified by this field cease to be met at some point during pod execution (e.g. due to an update), the system may or may not try to eventually evict the pod from its node.
	requiredDuringSchedulingIgnoredDuringExecution?: #NodeSelector
}

#NodeInclusionPolicy: string

#NodeSelector: {
	// Required. A list of node selector terms. The terms are ORed.
	nodeSelectorTerms: [...#NodeSelectorTerm]
}

#NodeSelectorOperator: string

#NodeSelectorRequirement: {
	// The label key that the selector applies to.
	key: string
	// Represents a key's relationship to a set of values. Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
	operator: #NodeSelectorOperator
	// An array of string values. If the operator is In or NotIn, the values array must be non-empty. If the operator is Exists or DoesNotExist, the values array must be empty. If the operator is Gt or Lt, the values array must have a single element, which will be interpreted as an integer. This array is replaced during a strategic merge patch.
	values?: [...string]
}

#NodeSelectorTerm: {
	// A list of node selector requirements by node's labels.
	matchExpressions?: [...#NodeSelectorRequirement]
	// A list of node selector requirements by node's fields.
	matchFields?: [...#NodeSelectorRequirement]
}

#OSName: string

#ObjectMeta: {
	// Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata. They are not queryable and should be preserved when modifying objects. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations
	annotations?: {
		[X=string]: string
	}
	// CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC.
	// 
	// Populated by the system. Read-only. Null for lists. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	creationTimestamp?: #Time
	// Number of seconds allowed for this object to gracefully terminate before it will be removed from the system. Only set when deletionTimestamp is also set. May only be shortened. Read-only.
	deletionGracePeriodSeconds?: int
	// DeletionTimestamp is RFC 3339 date and time at which this resource will be deleted. This field is set by the server when a graceful deletion is requested by the user, and is not directly settable by a client. The resource is expected to be deleted (no longer visible from resource lists, and not reachable by name) after the time in this field, once the finalizers list is empty. As long as the finalizers list contains items, deletion is blocked. Once the deletionTimestamp is set, this value may not be unset or be set further into the future, although it may be shortened or the resource may be deleted prior to this time. For example, a user may request that a pod is deleted in 30 seconds. The Kubelet will react by sending a graceful termination signal to the containers in the pod. After that 30 seconds, the Kubelet will send a hard termination signal (SIGKILL) to the container and after cleanup, remove the pod from the API. In the presence of network partitions, this object may still exist after this timestamp, until an administrator or automated process can determine the resource is fully terminated. If not set, graceful deletion of the object has not been requested.
	// 
	// Populated by the system when a graceful deletion is requested. Read-only. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	deletionTimestamp?: #Time
	// Must be empty before the object is deleted from the registry. Each entry is an identifier for the responsible component that will remove the entry from the list. If the deletionTimestamp of the object is non-nil, entries in this list can only be removed. Finalizers may be processed and removed in any order.  Order is NOT enforced because it introduces significant risk of stuck finalizers. finalizers is a shared field, any actor with permission can reorder it. If the finalizer list is processed in order, then this can lead to a situation in which the component responsible for the first finalizer in the list is waiting for a signal (field value, external system, or other) produced by a component responsible for a finalizer later in the list, resulting in a deadlock. Without enforced ordering finalizers are free to order amongst themselves and are not vulnerable to ordering changes in the list.
	finalizers?: [...string]
	// GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided. If this field is used, the name returned to the client will be different than the name passed. This value will also be combined with a unique suffix. The provided value has the same validation rules as the Name field, and may be truncated by the length of the suffix required to make the value unique on the server.
	// 
	// If this field is specified and the generated name exists, the server will return a 409.
	// 
	// Applied only if Name is not specified. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#idempotency
	generateName?: string
	// A sequence number representing a specific generation of the desired state. Populated by the system. Read-only.
	generation?: int
	// Map of string keys and values that can be used to organize and categorize (scope and select) objects. May match selectors of replication controllers and services. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels
	labels?: {
		[X=string]: string
	}
	// ManagedFields maps workflow-id and version to the set of fields that are managed by that workflow. This is mostly for internal housekeeping, and users typically shouldn't need to set or understand this field. A workflow can be the user's name, a controller's name, or the name of a specific apply path like "ci-cd". The set of fields is always in the version that the workflow used when modifying the object.
	managedFields?: [...#ManagedFieldsEntry]
	// Name must be unique within a namespace. Is required when creating resources, although some resources may allow a client to request the generation of an appropriate name automatically. Name is primarily intended for creation idempotence and configuration definition. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names
	name?: string
	// Namespace defines the space within which each name must be unique. An empty namespace is equivalent to the "default" namespace, but "default" is the canonical representation. Not all objects are required to be scoped to a namespace - the value of this field for those objects will be empty.
	// 
	// Must be a DNS_LABEL. Cannot be updated. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces
	namespace?: string
	// List of objects depended by this object. If ALL objects in the list have been deleted, this object will be garbage collected. If this object is managed by a controller, then an entry in this list will point to this controller, with the controller field set to true. There cannot be more than one managing controller.
	ownerReferences?: [...#OwnerReference]
	// An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed. May be used for optimistic concurrency, change detection, and the watch operation on a resource or set of resources. Clients must treat these values as opaque and passed unmodified back to the server. They may only be valid for a particular resource or set of resources.
	// 
	// Populated by the system. Read-only. Value must be treated as opaque by clients and . More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
	resourceVersion?: string
	// Deprecated: selfLink is a legacy read-only field that is no longer populated by the system.
	selfLink?: string
	// UID is the unique in time and space value for this object. It is typically generated by the server on successful creation of a resource and is not allowed to change on PUT operations.
	// 
	// Populated by the system. Read-only. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids
	uid?: #UID
}

#OwnerReference: {
	// API version of the referent.
	apiVersion: string
	// If true, AND if the owner has the "foregroundDeletion" finalizer, then the owner cannot be deleted from the key-value store until this reference is removed. See https://kubernetes.io/docs/concepts/architecture/garbage-collection/#foreground-deletion for how the garbage collector interacts with this field and enforces the foreground deletion. Defaults to false. To set this field, a user needs "delete" permission of the owner, otherwise 422 (Unprocessable Entity) will be returned.
	blockOwnerDeletion?: bool
	// If true, this reference points to the managing controller.
	controller?: bool
	// Kind of the referent. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	kind: string
	// Name of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#names
	name: string
	// UID of the referent. More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names#uids
	uid: #UID
}

#PersistentVolumeAccessMode: string

#PersistentVolumeClaim: {
	apiVersion?: string
	kind?:       string
	// Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metadata?: #ObjectMeta
	// spec defines the desired characteristics of a volume requested by a pod author. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
	spec?: #PersistentVolumeClaimSpec
	// status represents the current information/status of a persistent volume claim. Read-only. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
	status?: #PersistentVolumeClaimStatus
}

#PersistentVolumeClaimCondition: {
	// lastProbeTime is the time we probed the condition.
	lastProbeTime?: #Time
	// lastTransitionTime is the time the condition transitioned from one status to another.
	lastTransitionTime?: #Time
	// message is the human-readable message indicating details about last transition.
	message?: string
	// reason is a unique, this should be a short, machine understandable string that gives the reason for condition's last transition. If it reports "Resizing" that means the underlying persistent volume is being resized.
	reason?: string
	status:  #ConditionStatus
	type:    #PersistentVolumeClaimConditionType
}

#PersistentVolumeClaimConditionType: string

#PersistentVolumeClaimModifyVolumeStatus: string

#PersistentVolumeClaimPhase: string

#PersistentVolumeClaimRetentionPolicyType: string

#PersistentVolumeClaimSpec: {
	// accessModes contains the desired access modes the volume should have. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
	accessModes?: [...#PersistentVolumeAccessMode]
	// dataSource field can be used to specify either: * An existing VolumeSnapshot object (snapshot.storage.k8s.io/VolumeSnapshot) * An existing PVC (PersistentVolumeClaim) If the provisioner or an external controller can support the specified data source, it will create a new volume based on the contents of the specified data source. When the AnyVolumeDataSource feature gate is enabled, dataSource contents will be copied to dataSourceRef, and dataSourceRef contents will be copied to dataSource when dataSourceRef.namespace is not specified. If the namespace is specified, then dataSourceRef will not be copied to dataSource.
	dataSource?: #TypedLocalObjectReference
	// dataSourceRef specifies the object from which to populate the volume with data, if a non-empty volume is desired. This may be any object from a non-empty API group (non core object) or a PersistentVolumeClaim object. When this field is specified, volume binding will only succeed if the type of the specified object matches some installed volume populator or dynamic provisioner. This field will replace the functionality of the dataSource field and as such if both fields are non-empty, they must have the same value. For backwards compatibility, when namespace isn't specified in dataSourceRef, both fields (dataSource and dataSourceRef) will be set to the same value automatically if one of them is empty and the other is non-empty. When namespace is specified in dataSourceRef, dataSource isn't set to the same value and must be empty. There are three important differences between dataSource and dataSourceRef: * While dataSource only allows two specific types of objects, dataSourceRef
	//   allows any non-core object, as well as PersistentVolumeClaim objects.
	// * While dataSource ignores disallowed values (dropping them), dataSourceRef
	//   preserves all values, and generates an error if a disallowed value is
	//   specified.
	// * While dataSource only allows local objects, dataSourceRef allows objects
	//   in any namespaces.
	// (Beta) Using this field requires the AnyVolumeDataSource feature gate to be enabled. (Alpha) Using the namespace field of dataSourceRef requires the CrossNamespaceVolumeDataSource feature gate to be enabled.
	dataSourceRef?: #TypedObjectReference
	// resources represents the minimum resources the volume should have. If RecoverVolumeExpansionFailure feature is enabled users are allowed to specify resource requirements that are lower than previous value but must still be higher than capacity recorded in the status field of the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources
	resources?: #VolumeResourceRequirements
	// selector is a label query over volumes to consider for binding.
	selector?: #LabelSelector
	// storageClassName is the name of the StorageClass required by the claim. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1
	storageClassName?: string
	// volumeAttributesClassName may be used to set the VolumeAttributesClass used by this claim. If specified, the CSI driver will create or update the volume with the attributes defined in the corresponding VolumeAttributesClass. This has a different purpose than storageClassName, it can be changed after the claim is created. An empty string value means that no VolumeAttributesClass will be applied to the claim but it's not allowed to reset this field to empty string once it is set. If unspecified and the PersistentVolumeClaim is unbound, the default VolumeAttributesClass will be set by the persistentvolume controller if it exists. If the resource referred to by volumeAttributesClass does not exist, this PersistentVolumeClaim will be set to a Pending state, as reflected by the modifyVolumeStatus field, until such as a resource exists. More info: https://kubernetes.io/docs/concepts/storage/volume-attributes-classes/ (Alpha) Using this field requires the VolumeAttributesClass feature gate to be enabled.
	volumeAttributesClassName?: string
	// volumeMode defines what type of volume is required by the claim. Value of Filesystem is implied when not included in claim spec.
	volumeMode?: #PersistentVolumeMode
	// volumeName is the binding reference to the PersistentVolume backing this claim.
	volumeName?: string
}

#PersistentVolumeClaimStatus: {
	// accessModes contains the actual access modes the volume backing the PVC has. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
	accessModes?: [...#PersistentVolumeAccessMode]
	// allocatedResourceStatuses stores status of resource being resized for the given PVC. Key names follow standard Kubernetes label syntax. Valid values are either:
	// 	* Un-prefixed keys:
	// 		- storage - the capacity of the volume.
	// 	* Custom resources must use implementation-defined prefixed names such as "example.com/my-custom-resource"
	// Apart from above values - keys that are unprefixed or have kubernetes.io prefix are considered reserved and hence may not be used.
	// 
	// ClaimResourceStatus can be in any of following states:
	// 	- ControllerResizeInProgress:
	// 		State set when resize controller starts resizing the volume in control-plane.
	// 	- ControllerResizeFailed:
	// 		State set when resize has failed in resize controller with a terminal error.
	// 	- NodeResizePending:
	// 		State set when resize controller has finished resizing the volume but further resizing of
	// 		volume is needed on the node.
	// 	- NodeResizeInProgress:
	// 		State set when kubelet starts resizing the volume.
	// 	- NodeResizeFailed:
	// 		State set when resizing has failed in kubelet with a terminal error. Transient errors don't set
	// 		NodeResizeFailed.
	// For example: if expanding a PVC for more capacity - this field can be one of the following states:
	// 	- pvc.status.allocatedResourceStatus['storage'] = "ControllerResizeInProgress"
	//      - pvc.status.allocatedResourceStatus['storage'] = "ControllerResizeFailed"
	//      - pvc.status.allocatedResourceStatus['storage'] = "NodeResizePending"
	//      - pvc.status.allocatedResourceStatus['storage'] = "NodeResizeInProgress"
	//      - pvc.status.allocatedResourceStatus['storage'] = "NodeResizeFailed"
	// When this field is not set, it means that no resize operation is in progress for the given PVC.
	// 
	// A controller that receives PVC update with previously unknown resourceName or ClaimResourceStatus should ignore the update for the purpose it was designed. For example - a controller that only is responsible for resizing capacity of the volume, should ignore PVC updates that change other valid resources associated with PVC.
	// 
	// This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
	allocatedResourceStatuses?: {
		[X=#ResourceName]: #ClaimResourceStatus
	}
	// allocatedResources tracks the resources allocated to a PVC including its capacity. Key names follow standard Kubernetes label syntax. Valid values are either:
	// 	* Un-prefixed keys:
	// 		- storage - the capacity of the volume.
	// 	* Custom resources must use implementation-defined prefixed names such as "example.com/my-custom-resource"
	// Apart from above values - keys that are unprefixed or have kubernetes.io prefix are considered reserved and hence may not be used.
	// 
	// Capacity reported here may be larger than the actual capacity when a volume expansion operation is requested. For storage quota, the larger value from allocatedResources and PVC.spec.resources is used. If allocatedResources is not set, PVC.spec.resources alone is used for quota calculation. If a volume expansion capacity request is lowered, allocatedResources is only lowered if there are no expansion operations in progress and if the actual volume capacity is equal or lower than the requested capacity.
	// 
	// A controller that receives PVC update with previously unknown resourceName should ignore the update for the purpose it was designed. For example - a controller that only is responsible for resizing capacity of the volume, should ignore PVC updates that change other valid resources associated with PVC.
	// 
	// This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
	allocatedResources?: #ResourceList
	// capacity represents the actual resources of the underlying volume.
	capacity?: #ResourceList
	// conditions is the current Condition of persistent volume claim. If underlying persistent volume is being resized then the Condition will be set to 'Resizing'.
	conditions?: [...#PersistentVolumeClaimCondition]
	// currentVolumeAttributesClassName is the current name of the VolumeAttributesClass the PVC is using. When unset, there is no VolumeAttributeClass applied to this PersistentVolumeClaim This is an alpha field and requires enabling VolumeAttributesClass feature.
	currentVolumeAttributesClassName?: string
	// ModifyVolumeStatus represents the status object of ControllerModifyVolume operation. When this is unset, there is no ModifyVolume operation being attempted. This is an alpha field and requires enabling VolumeAttributesClass feature.
	modifyVolumeStatus?: #ModifyVolumeStatus
	// phase represents the current phase of PersistentVolumeClaim.
	phase?: #PersistentVolumeClaimPhase
}

#PersistentVolumeClaimVolumeSource: {
	// claimName is the name of a PersistentVolumeClaim in the same namespace as the pod using this volume. More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
	claimName: string
	// else volumeMounts
	readOnly?: bool
}

#PersistentVolumeMode: string

#PodAffinity: {
	// The scheduler will prefer to schedule pods to nodes that satisfy the affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions. The node that is most preferred is the one with the greatest sum of weights, i.e. for each node that meets all of the scheduling requirements (resource request, requiredDuringScheduling affinity expressions, etc.), compute a sum by iterating through the elements of this field and adding "weight" to the sum if the node has pods which matches the corresponding podAffinityTerm; the node(s) with the highest sum are the most preferred.
	preferredDuringSchedulingIgnoredDuringExecution?: [...#WeightedPodAffinityTerm]
	// If the affinity requirements specified by this field are not met at scheduling time, the pod will not be scheduled onto the node. If the affinity requirements specified by this field cease to be met at some point during pod execution (e.g. due to a pod label update), the system may or may not try to eventually evict the pod from its node. When there are multiple elements, the lists of nodes corresponding to each podAffinityTerm are intersected, i.e. all terms must be satisfied.
	requiredDuringSchedulingIgnoredDuringExecution?: [...#PodAffinityTerm]
}

#PodAffinityTerm: {
	// A label query over a set of resources, in this case pods. If it's null, this PodAffinityTerm matches with no Pods.
	labelSelector?: #LabelSelector
	// MatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key in (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both matchLabelKeys and labelSelector. Also, matchLabelKeys cannot be set when labelSelector isn't set. This is an alpha field and requires enabling MatchLabelKeysInPodAffinity feature gate.
	matchLabelKeys?: [...string]
	// MismatchLabelKeys is a set of pod label keys to select which pods will be taken into consideration. The keys are used to lookup values from the incoming pod labels, those key-value labels are merged with `labelSelector` as `key notin (value)` to select the group of existing pods which pods will be taken into consideration for the incoming pod's pod (anti) affinity. Keys that don't exist in the incoming pod labels will be ignored. The default value is empty. The same key is forbidden to exist in both mismatchLabelKeys and labelSelector. Also, mismatchLabelKeys cannot be set when labelSelector isn't set. This is an alpha field and requires enabling MatchLabelKeysInPodAffinity feature gate.
	mismatchLabelKeys?: [...string]
	// A label query over the set of namespaces that the term applies to. The term is applied to the union of the namespaces selected by this field and the ones listed in the namespaces field. null selector and null or empty namespaces list means "this pod's namespace". An empty selector ({}) matches all namespaces.
	namespaceSelector?: #LabelSelector
	// namespaces specifies a static list of namespace names that the term applies to. The term is applied to the union of the namespaces listed in this field and the ones selected by namespaceSelector. null or empty namespaces list and null namespaceSelector means "this pod's namespace".
	namespaces?: [...string]
	// This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching the labelSelector in the specified namespaces, where co-located is defined as running on a node whose value of the label with key topologyKey matches that of any node on which any of the selected pods is running. Empty topologyKey is not allowed.
	topologyKey: string
}

#PodAntiAffinity: {
	// The scheduler will prefer to schedule pods to nodes that satisfy the anti-affinity expressions specified by this field, but it may choose a node that violates one or more of the expressions. The node that is most preferred is the one with the greatest sum of weights, i.e. for each node that meets all of the scheduling requirements (resource request, requiredDuringScheduling anti-affinity expressions, etc.), compute a sum by iterating through the elements of this field and adding "weight" to the sum if the node has pods which matches the corresponding podAffinityTerm; the node(s) with the highest sum are the most preferred.
	preferredDuringSchedulingIgnoredDuringExecution?: [...#WeightedPodAffinityTerm]
	// If the anti-affinity requirements specified by this field are not met at scheduling time, the pod will not be scheduled onto the node. If the anti-affinity requirements specified by this field cease to be met at some point during pod execution (e.g. due to a pod label update), the system may or may not try to eventually evict the pod from its node. When there are multiple elements, the lists of nodes corresponding to each podAffinityTerm are intersected, i.e. all terms must be satisfied.
	requiredDuringSchedulingIgnoredDuringExecution?: [...#PodAffinityTerm]
}

#PodConditionType: string

#PodDNSConfig: {
	// A list of DNS name server IP addresses. This will be appended to the base nameservers generated from DNSPolicy. Duplicated nameservers will be removed.
	nameservers?: [...string]
	// A list of DNS resolver options. This will be merged with the base options generated from DNSPolicy. Duplicated entries will be removed. Resolution options given in Options will override those that appear in the base DNSPolicy.
	options?: [...#PodDNSConfigOption]
	// A list of DNS search domains for host-name lookup. This will be appended to the base search paths generated from DNSPolicy. Duplicated search paths will be removed.
	searches?: [...string]
}

#PodDNSConfigOption: {
	// Required.
	name?:  string
	value?: string
}

#PodFSGroupChangePolicy: string

#PodFailurePolicy: {
	// A list of pod failure policy rules. The rules are evaluated in order. Once a rule matches a Pod failure, the remaining of the rules are ignored. When no rule matches the Pod failure, the default handling applies - the counter of pod failures is incremented and it is checked against the backoffLimit. At most 20 elements are allowed.
	rules: [...#PodFailurePolicyRule]
}

#PodFailurePolicyAction: string

#PodFailurePolicyOnExitCodesOperator: string

#PodFailurePolicyOnExitCodesRequirement: {
	// Restricts the check for exit codes to the container with the specified name. When null, the rule applies to all containers. When specified, it should match one the container or initContainer names in the pod template.
	containerName: string
	// Represents the relationship between the container exit code(s) and the specified values. Containers completed with success (exit code 0) are excluded from the requirement check. Possible values are:
	// 
	// - In: the requirement is satisfied if at least one container exit code
	//   (might be multiple if there are multiple containers not restricted
	//   by the 'containerName' field) is in the set of specified values.
	// - NotIn: the requirement is satisfied if at least one container exit code
	//   (might be multiple if there are multiple containers not restricted
	//   by the 'containerName' field) is not in the set of specified values.
	// Additional values are considered to be added in the future. Clients should react to an unknown operator by assuming the requirement is not satisfied.
	operator: #PodFailurePolicyOnExitCodesOperator
	// Specifies the set of values. Each returned container exit code (might be multiple in case of multiple containers) is checked against this set of values with respect to the operator. The list of values must be ordered and must not contain duplicates. Value '0' cannot be used for the In operator. At least one element is required. At most 255 elements are allowed.
	values: [...int]
}

#PodFailurePolicyOnPodConditionsPattern: {
	// Specifies the required Pod condition status. To match a pod condition it is required that the specified status equals the pod condition status. Defaults to True.
	status: #ConditionStatus
	// Specifies the required Pod condition type. To match a pod condition it is required that specified type equals the pod condition type.
	type: #PodConditionType
}

#PodFailurePolicyRule: {
	// Specifies the action taken on a pod failure when the requirements are satisfied. Possible values are:
	// 
	// - FailJob: indicates that the pod's job is marked as Failed and all
	//   running pods are terminated.
	// - FailIndex: indicates that the pod's index is marked as Failed and will
	//   not be restarted.
	//   This value is beta-level. It can be used when the
	//   `JobBackoffLimitPerIndex` feature gate is enabled (enabled by default).
	// - Ignore: indicates that the counter towards the .backoffLimit is not
	//   incremented and a replacement pod is created.
	// - Count: indicates that the pod is handled in the default way - the
	//   counter towards the .backoffLimit is incremented.
	// Additional values are considered to be added in the future. Clients should react to an unknown action by skipping the rule.
	action: #PodFailurePolicyAction
	// Represents the requirement on the container exit codes.
	onExitCodes: #PodFailurePolicyOnExitCodesRequirement
	// Represents the requirement on the pod conditions. The requirement is represented as a list of pod condition patterns. The requirement is satisfied if at least one pattern matches an actual pod condition. At most 20 elements are allowed.
	onPodConditions: [...#PodFailurePolicyOnPodConditionsPattern]
}

#PodManagementPolicyType: string

#PodOS: {
	// Name is the name of the operating system. The currently supported values are linux and windows. Additional value may be defined in future and can be one of: https://github.com/opencontainers/runtime-spec/blob/master/config.md#platform-specific-configuration Clients should expect to handle additional values and treat unrecognized values in this field as os: null
	name: #OSName
}

#PodPartialSpec: {
	// Optional duration in seconds the pod may be active on the node relative to
	// StartTime before the system will actively try to mark it failed and kill associated containers.
	// Value must be a positive integer.
	activeDeadlineSeconds?: int
	// If specified, the pod's scheduling constraints
	affinity?: #Affinity
	// AutomountServiceAccountToken indicates whether a service account token should be automatically mounted.
	automountServiceAccountToken?: bool
	// Specifies the DNS parameters of a pod.
	// Parameters specified here will be merged to the generated DNS
	// configuration based on DNSPolicy.
	dnsConfig?: #PodDNSConfig
	// Set DNS policy for the pod.
	// Defaults to "ClusterFirst".
	// Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'.
	// DNS parameters given in DNSConfig will be merged with the policy selected with DNSPolicy.
	// To have DNS options set along with hostNetwork, you have to specify DNS policy
	// explicitly to 'ClusterFirstWithHostNet'.
	dnsPolicy?: #DNSPolicy
	// EnableServiceLinks indicates whether information about services should be injected into pod's
	// environment variables, matching the syntax of Docker links.
	// Optional: Defaults to true.
	enableServiceLinks?: bool
	// HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts
	// file if specified.
	hostAliases?: [...#HostAlias]
	// Use the host's ipc namespace.
	// Optional: Default to false.
	hostIPC?: bool
	// Host networking requested for this pod. Use the host's network namespace.
	// If this option is set, the ports that will be used must be specified.
	// Default to false.
	hostNetwork?: bool
	// Use the host's pid namespace.
	// Optional: Default to false.
	hostPID?: bool
	// Use the host's user namespace.
	// Optional: Default to true.
	// If set to true or not present, the pod will be run in the host user namespace, useful
	// for when the pod needs a feature only available to the host user namespace, such as
	// loading a kernel module with CAP_SYS_MODULE.
	// When set to false, a new userns is created for the pod. Setting false is useful for
	// mitigating container breakout vulnerabilities even allowing users to run their
	// containers as root without actually having root privileges on the host.
	// This field is alpha-level and is only honored by servers that enable the UserNamespacesSupport feature.
	hostUsers?: bool
	// Specifies the hostname of the Pod
	// If not specified, the pod's hostname will be set to a system-defined value.
	hostname?: string
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	imagePullSecrets?: [...#LocalObjectReference]
	// NodeName is a request to schedule this pod onto a specific node. If it is non-empty,
	// the scheduler simply schedules this pod onto that node, assuming that it fits resource
	// requirements.
	nodeName?: string
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	nodeSelector?: {
		[X=string]: string
	}
	// Specifies the OS of the containers in the pod.
	// Some pod and container fields are restricted if this is set.
	// 
	// If the OS field is set to linux, the following fields must be unset:
	// -securityContext.windowsOptions
	// 
	// If the OS field is set to windows, following fields must be unset:
	// - spec.hostPID
	// - spec.hostIPC
	// - spec.hostUsers
	// - spec.securityContext.appArmorProfile
	// - spec.securityContext.seLinuxOptions
	// - spec.securityContext.seccompProfile
	// - spec.securityContext.fsGroup
	// - spec.securityContext.fsGroupChangePolicy
	// - spec.securityContext.sysctls
	// - spec.shareProcessNamespace
	// - spec.securityContext.runAsUser
	// - spec.securityContext.runAsGroup
	// - spec.securityContext.supplementalGroups
	// - spec.containers[*].securityContext.appArmorProfile
	// - spec.containers[*].securityContext.seLinuxOptions
	// - spec.containers[*].securityContext.seccompProfile
	// - spec.containers[*].securityContext.capabilities
	// - spec.containers[*].securityContext.readOnlyRootFilesystem
	// - spec.containers[*].securityContext.privileged
	// - spec.containers[*].securityContext.allowPrivilegeEscalation
	// - spec.containers[*].securityContext.procMount
	// - spec.containers[*].securityContext.runAsUser
	// - spec.containers[*].securityContext.runAsGroup
	os?: #PodOS
	// Overhead represents the resource overhead associated with running a pod for a given RuntimeClass.
	// This field will be autopopulated at admission time by the RuntimeClass admission controller. If
	// the RuntimeClass admission controller is enabled, overhead must not be set in Pod create requests.
	// The RuntimeClass admission controller will reject Pod create requests which have the overhead already
	// set. If RuntimeClass is configured and selected in the PodSpec, Overhead will be set to the value
	// defined in the corresponding RuntimeClass, otherwise it will remain unset and treated as zero.
	// More info: https://git.k8s.io/enhancements/keps/sig-node/688-pod-overhead/README.md
	overhead?: #ResourceList
	// PreemptionPolicy is the Policy for preempting pods with lower priority.
	// One of Never, PreemptLowerPriority.
	// Defaults to PreemptLowerPriority if unset.
	preemptionPolicy?: #PreemptionPolicy
	// The priority value. Various system components use this field to find the
	// priority of the pod. When Priority Admission Controller is enabled, it
	// prevents users from setting this field. The admission controller populates
	// this field from PriorityClassName.
	// The higher the value, the higher the priority.
	priority?: int
	// If specified, indicates the pod's priority. "system-node-critical" and
	// "system-cluster-critical" are two special keywords which indicate the
	// highest priorities with the former being the highest priority. Any other
	// name must be defined by creating a PriorityClass object with that name.
	// If not specified, the pod priority will be default or zero if there is no
	// default.
	priorityClassName?: string
	// If specified, all readiness gates will be evaluated for pod readiness.
	// A pod is ready when all its containers are ready AND
	// all conditions specified in the readiness gates have status equal to "True"
	// More info: https://git.k8s.io/enhancements/keps/sig-network/580-pod-readiness-gates
	readinessGates?: [...#PodReadinessGate]
	// ResourceClaims defines which ResourceClaims must be allocated
	// and reserved before the Pod is allowed to start. The resources
	// will be made available to those containers which consume them
	// by name.
	// 
	// This is an alpha field and requires enabling the
	// DynamicResourceAllocation feature gate.
	// 
	// This field is immutable.
	resourceClaims?: [...#PodResourceClaim]
	// Restart policy for all containers within the pod.
	// One of Always, OnFailure, Never. In some contexts, only a subset of those values may be permitted.
	// Default to Always.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy
	restartPolicy?: #RestartPolicy
	// RuntimeClassName refers to a RuntimeClass object in the node.k8s.io group, which should be used
	// to run this pod.  If no RuntimeClass resource matches the named class, the pod will not be run.
	// If unset or empty, the "legacy" RuntimeClass will be used, which is an implicit class with an
	// empty definition that uses the default runtime handler.
	// More info: https://git.k8s.io/enhancements/keps/sig-node/585-runtime-class
	runtimeClassName?: string
	// If specified, the pod will be dispatched by specified scheduler.
	// If not specified, the pod will be dispatched by default scheduler.
	schedulerName?: string
	// SchedulingGates is an opaque list of values that if specified will block scheduling the pod.
	// If schedulingGates is not empty, the pod will stay in the SchedulingGated state and the
	// scheduler will not attempt to schedule the pod.
	// 
	// SchedulingGates can only be set at pod creation time, and be removed only afterwards.
	schedulingGates?: [...#PodSchedulingGate]
	// SecurityContext holds pod-level security attributes and common container settings.
	// Optional: Defaults to empty.  See type description for default values of each field.
	securityContext?: #PodSecurityContext
	// If true the pod's hostname will be configured as the pod's FQDN, rather than the leaf name (the default).
	// In Linux containers, this means setting the FQDN in the hostname field of the kernel (the nodename field of struct utsname).
	// In Windows containers, this means setting the registry value of hostname for the registry key HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters to FQDN.
	// If a pod does not have FQDN, this has no effect.
	// Default to false.
	setHostnameAsFQDN?: bool
	// Share a single process namespace between all of the containers in a pod.
	// When this is set containers will be able to view and signal processes from other containers
	// in the same pod, and the first process in each container will not be assigned PID 1.
	// HostPID and ShareProcessNamespace cannot both be set.
	// Optional: Default to false.
	shareProcessNamespace?: bool
	// If specified, the fully qualified Pod hostname will be "<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>".
	// If not specified, the pod will not have a domainname at all.
	subdomain?: string
	// Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request.
	// Value must be non-negative integer. The value zero indicates stop immediately via
	// the kill signal (no opportunity to shut down).
	// If this value is nil, the default grace period will be used instead.
	// The grace period is the duration in seconds after the processes running in the pod are sent
	// a termination signal and the time when the processes are forcibly halted with a kill signal.
	// Set this value longer than the expected cleanup time for your process.
	// Defaults to 30 seconds.
	terminationGracePeriodSeconds?: int
	// If specified, the pod's tolerations.
	tolerations?: [...#Toleration]
	// TopologySpreadConstraints describes how a group of pods ought to spread across topology
	// domains. Scheduler will schedule pods in a way which abides by the constraints.
	// All topologySpreadConstraints are ANDed.
	topologySpreadConstraints?: [...#TopologySpreadConstraint]
}

#PodPartialTemplateSpec: spec?: #PodPartialSpec

#PodReadinessGate: {
	// ConditionType refers to a condition in the pod's condition list with matching type.
	conditionType: #PodConditionType
}

#PodReplacementPolicy: string

#PodResourceClaim: {
	// Name uniquely identifies this resource claim inside the pod. This must be a DNS_LABEL.
	name: string
	// Source describes where to find the ResourceClaim.
	source?: #ClaimSource
}

#PodSchedulingGate: {
	// Name of the scheduling gate. Each scheduling gate must have a unique name field.
	name: string
}

#PodSecurityContext: {
	// appArmorProfile is the AppArmor options to use by the containers in this pod. Note that this field cannot be set when spec.os.name is windows.
	appArmorProfile?: #AppArmorProfile
	// A special supplemental group that applies to all containers in a pod. Some volume types allow the Kubelet to change the ownership of that volume to be owned by the pod:
	// 
	// 1. The owning GID will be the FSGroup 2. The setgid bit is set (new files created in the volume will be owned by FSGroup) 3. The permission bits are OR'd with rw-rw 
	fsGroup?: int
	// fsGroupChangePolicy defines behavior of changing ownership and permission of the volume before being exposed inside Pod. This field will only apply to volume types which support fsGroup based ownership(and permissions). It will have no effect on ephemeral volume types such as: secret, configmaps and emptydir. Valid values are "OnRootMismatch" and "Always". If not specified, "Always" is used. Note that this field cannot be set when spec.os.name is windows.
	fsGroupChangePolicy?: #PodFSGroupChangePolicy
	// The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in SecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container. Note that this field cannot be set when spec.os.name is windows.
	runAsGroup?: int
	// Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in SecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
	runAsNonRoot?: bool
	// The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in SecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container. Note that this field cannot be set when spec.os.name is windows.
	runAsUser?: int
	// The SELinux context to be applied to all containers. If unspecified, the container runtime will allocate a random SELinux context for each container.  May also be set in SecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence for that container. Note that this field cannot be set when spec.os.name is windows.
	seLinuxOptions?: #SELinuxOptions
	// The seccomp options to use by the containers in this pod. Note that this field cannot be set when spec.os.name is windows.
	seccompProfile?: #SeccompProfile
	// A list of groups applied to the first process run in each container, in addition to the container's primary GID, the fsGroup (if specified), and group memberships defined in the container image for the uid of the container process. If unspecified, no additional groups are added to any container. Note that group memberships defined in the container image for the uid of the container process are still effective, even if they are not included in this list. Note that this field cannot be set when spec.os.name is windows.
	supplementalGroups?: [...int]
	// Sysctls hold a list of namespaced sysctls used for the pod. Pods with unsupported sysctls (by the container runtime) might fail to launch. Note that this field cannot be set when spec.os.name is windows.
	sysctls?: [...#Sysctl]
	// The Windows specific settings applied to all containers. If unspecified, the options within a container's SecurityContext will be used. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is linux.
	windowsOptions?: #WindowsSecurityContextOptions
}

#PolicyRule: {
	// APIGroups is the name of the APIGroup that contains the resources.  If multiple API groups are specified, any action requested against one of the enumerated resources in any API group will be allowed. "" represents the core API group and "*" represents all API groups.
	apiGroups?: [...string]
	// NonResourceURLs is a set of partial urls that a user should have access to.  *s are allowed, but only as the full, final step in the path Since non-resource URLs are not namespaced, this field is only applicable for ClusterRoles referenced from a ClusterRoleBinding. Rules can either apply to API resources (such as "pods" or "secrets") or non-resource URL paths (such as "/api"),  but not both.
	nonResourceURLs?: [...string]
	// ResourceNames is an optional white list of names that the rule applies to.  An empty set means that everything is allowed.
	resourceNames?: [...string]
	// Resources is a list of resources this rule applies to. '*' represents all resources.
	resources?: [...string]
	// Verbs is a list of Verbs that apply to ALL the ResourceKinds contained in this rule. '*' represents all verbs.
	verbs: [...string]
}

#PreemptionPolicy: string

#PreferredSchedulingTerm: {
	// A node selector term, associated with the corresponding weight.
	preference: #NodeSelectorTerm
	// Weight associated with matching the corresponding nodeSelectorTerm, in the range 1-100.
	weight: int
}

#Probe: {
	exec?: #ExecAction
	// Minimum consecutive failures for the probe to be considered failed after having succeeded. Defaults to 3. Minimum value is 1.
	failureThreshold?: int
	grpc?:             #GRPCAction
	httpGet?:          #HTTPGetAction
	// Number of seconds after the container has started before liveness probes are initiated. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	initialDelaySeconds?: int
	// How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
	periodSeconds?: int
	// Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
	successThreshold?: int
	tcpSocket?:        #TCPSocketAction
	// Optional duration in seconds the pod needs to terminate gracefully upon probe failure. The grace period is the duration in seconds after the processes running in the pod are sent a termination signal and the time when the processes are forcibly halted with a kill signal. Set this value longer than the expected cleanup time for your process. If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this value overrides the value provided by the pod spec. Value must be non-negative integer. The value zero indicates stop immediately via the kill signal (no opportunity to shut down). This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate. Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
	terminationGracePeriodSeconds?: int
	// Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1. More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	timeoutSeconds?: int
}

#ProcMountType: string

#PullPolicy: string

#Quantity: string

#ResourceClaim: {
	// Name must match the name of one entry in pod.spec.resourceClaims of the Pod where this field is used. It makes that resource available inside a container.
	name: string
}

#ResourceList: [X=#ResourceName]: #Quantity

#ResourceName: string

#ResourceRequirements: {
	// Claims lists the names of resources, defined in spec.resourceClaims, that are used by this container.
	// 
	// This is an alpha field and requires enabling the DynamicResourceAllocation feature gate.
	// 
	// This field is immutable. It can only be set for containers.
	claims?: [...#ResourceClaim]
	// Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	limits?: #ResourceList
	// Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	requests?: #ResourceList
}

#RestartPolicy: string

#RollingUpdateDaemonSet: {
	// The maximum number of nodes with an existing available DaemonSet pod that can have an updated DaemonSet pod during during an update. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). This can not be 0 if MaxUnavailable is 0. Absolute number is calculated from percentage by rounding up to a minimum of 1. Default value is 0. Example: when this is set to 30%, at most 30% of the total number of nodes that should be running the daemon pod (i.e. status.desiredNumberScheduled) can have their a new pod created before the old pod is marked as deleted. The update starts by launching new pods on 30% of nodes. Once an updated pod is available (Ready for at least minReadySeconds) the old DaemonSet pod on that node is marked deleted. If the old pod becomes unavailable for any reason (Ready transitions to false, is evicted, or is drained) an updated pod is immediatedly created on that node without considering surge limits. Allowing surge implies the possibility that the resources consumed by the daemonset on any given node can double if the readiness check fails, and so resource intensive daemonsets should take into account that they may cause evictions during disruption.
	maxSurge?: #IntOrString
	// The maximum number of DaemonSet pods that can be unavailable during the update. Value can be an absolute number (ex: 5) or a percentage of total number of DaemonSet pods at the start of the update (ex: 10%). Absolute number is calculated from percentage by rounding up. This cannot be 0 if MaxSurge is 0 Default value is 1. Example: when this is set to 30%, at most 30% of the total number of nodes that should be running the daemon pod (i.e. status.desiredNumberScheduled) can have their pods stopped for an update at any given time. The update starts by stopping at most 30% of those DaemonSet pods and then brings up new DaemonSet pods in their place. Once the new pods are available, it then proceeds onto other DaemonSet pods, thus ensuring that at least 70% of original number of DaemonSet pods are available at all times during the update.
	maxUnavailable?: #IntOrString
}

#RollingUpdateDeployment: {
	// The maximum number of pods that can be scheduled above the desired number of pods. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). This can not be 0 if MaxUnavailable is 0. Absolute number is calculated from percentage by rounding up. Defaults to 25%. Example: when this is set to 30%, the new ReplicaSet can be scaled up immediately when the rolling update starts, such that the total number of old and new pods do not exceed 130% of desired pods. Once old pods have been killed, new ReplicaSet can be scaled up further, ensuring that total number of pods running at any time during the update is at most 130% of desired pods.
	maxSurge?: #IntOrString
	// The maximum number of pods that can be unavailable during the update. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). Absolute number is calculated from percentage by rounding down. This can not be 0 if MaxSurge is 0. Defaults to 25%. Example: when this is set to 30%, the old ReplicaSet can be scaled down to 70% of desired pods immediately when the rolling update starts. Once new pods are ready, old ReplicaSet can be scaled down further, followed by scaling up the new ReplicaSet, ensuring that the total number of pods available at all times during the update is at least 70% of desired pods.
	maxUnavailable?: #IntOrString
}

#RollingUpdateStatefulSetStrategy: {
	// The maximum number of pods that can be unavailable during the update. Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%). Absolute number is calculated from percentage by rounding up. This can not be 0. Defaults to 1. This field is alpha-level and is only honored by servers that enable the MaxUnavailableStatefulSet feature. The field applies to all pods in the range 0 to Replicas-1. That means if there is any unavailable pod in the range 0 to Replicas-1, it will be counted towards MaxUnavailable.
	maxUnavailable?: #IntOrString
	// Partition indicates the ordinal at which the StatefulSet should be partitioned for updates. During a rolling update, all pods from ordinal Replicas-1 to Partition are updated. All pods from ordinal Partition-1 to 0 remain untouched. This is helpful in being able to do a canary based deployment. The default value is 0.
	partition?: int
}

#SELinuxOptions: {
	// Level is SELinux level label that applies to the container.
	level?: string
	// Role is a SELinux role label that applies to the container.
	role?: string
	// Type is a SELinux type label that applies to the container.
	type?: string
	// User is a SELinux user label that applies to the container.
	user?: string
}

#ScopeType: "Cluster" | "Namespace"

#SeccompProfile: {
	// localhostProfile indicates a profile defined in a file on the node should be used. The profile must be preconfigured on the node to work. Must be a descending path, relative to the kubelet's configured seccomp profile location. Must be set if type is "Localhost". Must NOT be set for any other type.
	localhostProfile?: string
	// type indicates which kind of seccomp profile will be applied. Valid options are:
	// 
	// Localhost - a profile defined in a file on the node should be used. RuntimeDefault - the container runtime default profile should be used. Unconfined - no profile should be applied.
	type: #SeccompProfileType
}

#SeccompProfileType: string

#SecretVolumeSource: {
	// defaultMode is Optional: mode bits used to set permissions on created files by default. Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511. YAML accepts both octal and decimal values, JSON requires decimal values for mode bits. Defaults to 0644. Directories within the path are not affected by this setting. This might be in conflict with other options that affect the file mode, like fsGroup, and the result can be other mode bits set.
	defaultMode?: int
	// items If unspecified, each key-value pair in the Data field of the referenced Secret will be projected into the volume as a file whose name is the key and content is the value. If specified, the listed keys will be projected into the specified paths, and unlisted keys will not be present. If a key is specified which is not present in the Secret, the volume setup will error unless it is marked optional. Paths must be relative and may not contain the '..' path or start with '..'.
	items?: [...#KeyToPath]
	optional?: bool
	// secretName is the name of the secret in the pod's namespace to use. More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
	secretName?: string
}

#SecurityContext: {
	// AllowPrivilegeEscalation controls whether a process can gain more privileges than its parent process. This bool directly controls if the no_new_privs flag will be set on the container process. AllowPrivilegeEscalation is true always when the container is: 1) run as Privileged 2) has CAP_SYS_ADMIN Note that this field cannot be set when spec.os.name is windows.
	allowPrivilegeEscalation?: bool
	// appArmorProfile is the AppArmor options to use by this container. If set, this profile overrides the pod's appArmorProfile. Note that this field cannot be set when spec.os.name is windows.
	appArmorProfile?: #AppArmorProfile
	// The capabilities to add/drop when running containers. Defaults to the default set of capabilities granted by the container runtime. Note that this field cannot be set when spec.os.name is windows.
	capabilities?: #Capabilities
	// Run container in privileged mode. Processes in privileged containers are essentially equivalent to root on the host. Defaults to false. Note that this field cannot be set when spec.os.name is windows.
	privileged?: bool
	// procMount denotes the type of proc mount to use for the containers. The default is DefaultProcMount which uses the container runtime defaults for readonly paths and masked paths. This requires the ProcMountType feature flag to be enabled. Note that this field cannot be set when spec.os.name is windows.
	procMount?: #ProcMountType
	// Whether this container has a read-only root filesystem. Default is false. Note that this field cannot be set when spec.os.name is windows.
	readOnlyRootFilesystem?: bool
	// The GID to run the entrypoint of the container process. Uses runtime default if unset. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.
	runAsGroup?: int
	// Indicates that the container must run as a non-root user. If true, the Kubelet will validate the image at runtime to ensure that it does not run as UID 0 (root) and fail to start the container if it does. If unset or false, no such validation will be performed. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
	runAsNonRoot?: bool
	// The UID to run the entrypoint of the container process. Defaults to user specified in image metadata if unspecified. May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.
	runAsUser?: int
	// The SELinux context to be applied to the container. If unspecified, the container runtime will allocate a random SELinux context for each container.  May also be set in PodSecurityContext.  If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is windows.
	seLinuxOptions?: #SELinuxOptions
	// The seccomp options to use by this container. If seccomp options are provided at both the pod & container level, the container options override the pod options. Note that this field cannot be set when spec.os.name is windows.
	seccompProfile?: #SeccompProfile
	// The Windows specific settings applied to all containers. If unspecified, the options from the PodSecurityContext will be used. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence. Note that this field cannot be set when spec.os.name is linux.
	windowsOptions?: #WindowsSecurityContextOptions
}

#Service: {
	clusterIP?: string
	expose?:    #Expose
	// Paths [PortName]BashPath
	paths?: {
		[X=string]: string
	}
	// Ports [PortName]servicePort
	ports?: {
		[X=string]: int
	}
}

#ServiceAccount: {
	rules: [...#PolicyRule]
	scope?: #ScopeType
}

#SleepAction: {
	// Seconds is the number of seconds to sleep.
	seconds: int
}

#Spec: {
	config?: [X=string]: #EnvVarValueOrFrom
	containers?: [X=string]: #Container
	deploy?: #Deploy
	manifests?: [X=string]: _
	serviceAccount?: #ServiceAccount
	services?: [X=string]: #Service
	version: string
	volumes?: [X=string]: #Volume
}

#StatefulSetOrdinals: {
	// start is the number representing the first replica's index. It may be used to number replicas from an alternate index (eg: 1-indexed) over the default 0-indexed names, or to orchestrate progressive movement of replicas from one StatefulSet to another. If set, replica indices will be in the range:
	//   [.spec.ordinals.start, .spec.ordinals.start + .spec.replicas).
	// If unset, defaults to 0. Replica indices will be in the range:
	//   [0, .spec.replicas).
	start: int
}

#StatefulSetPersistentVolumeClaimRetentionPolicy: {
	// WhenDeleted specifies what happens to PVCs created from StatefulSet VolumeClaimTemplates when the StatefulSet is deleted. The default policy of `Retain` causes PVCs to not be affected by StatefulSet deletion. The `Delete` policy causes those PVCs to be deleted.
	whenDeleted?: #PersistentVolumeClaimRetentionPolicyType
	// WhenScaled specifies what happens to PVCs created from StatefulSet VolumeClaimTemplates when the StatefulSet is scaled down. The default policy of `Retain` causes PVCs to not be affected by a scaledown. The `Delete` policy causes the associated PVCs for any excess pods above the replica count to be deleted.
	whenScaled?: #PersistentVolumeClaimRetentionPolicyType
}

#StatefulSetSpec: {
	// Minimum number of seconds for which a newly created pod should be ready
	// without any of its container crashing for it to be considered available.
	// Defaults to 0 (pod will be considered available as soon as it is ready)
	minReadySeconds?: int
	// ordinals controls the numbering of replica indices in a StatefulSet. The
	// default ordinals behavior assigns a "0" index to the first replica and
	// increments the index by one for each additional replica requested. Using
	// the ordinals field requires the StatefulSetStartOrdinal feature gate to be
	// enabled, which is beta.
	ordinals?: #StatefulSetOrdinals
	// persistentVolumeClaimRetentionPolicy describes the lifecycle of persistent
	// volume claims created from volumeClaimTemplates. By default, all persistent
	// volume claims are created as needed and retained until manually deleted. This
	// policy allows the lifecycle to be altered, for example by deleting persistent
	// volume claims when their stateful set is deleted, or when their pod is scaled
	// down. This requires the StatefulSetAutoDeletePVC feature gate to be enabled,
	// which is alpha.  +optional
	persistentVolumeClaimRetentionPolicy?: #StatefulSetPersistentVolumeClaimRetentionPolicy
	// podManagementPolicy controls how pods are created during initial scale up,
	// when replacing pods on nodes, or when scaling down. The default policy is
	// `OrderedReady`, where pods are created in increasing order (pod-0, then
	// pod-1, etc) and the controller will wait until each pod is ready before
	// continuing. When scaling down, the pods are removed in the opposite order.
	// The alternative policy is `Parallel` which will create pods in parallel
	// to match the desired scale without waiting, and on scale down will delete
	// all pods at once.
	podManagementPolicy?: #PodManagementPolicyType
	// replicas is the desired number of replicas of the given Template.
	// These are replicas in the sense that they are instantiations of the
	// same Template, but individual replicas also have a consistent identity.
	// If unspecified, defaults to 1.
	// TODO: Consider a rename of this field.
	replicas?: int
	// revisionHistoryLimit is the maximum number of revisions that will
	// be maintained in the StatefulSet's revision history. The revision history
	// consists of all revisions not represented by a currently applied
	// StatefulSetSpec version. The default value is 10.
	revisionHistoryLimit?: int
	// serviceName is the name of the service that governs this StatefulSet.
	// This service must exist before the StatefulSet, and is responsible for
	// the network identity of the set. Pods get DNS/hostnames that follow the
	// pattern: pod-specific-string.serviceName.default.svc.cluster.local
	// where "pod-specific-string" is managed by the StatefulSet controller.
	serviceName: string
	// template is the object that describes the pod that will be created if
	// insufficient replicas are detected. Each pod stamped out by the StatefulSet
	// will fulfill this Template, but have a unique identity from the rest
	// of the StatefulSet. Each pod will be named with the format
	// <statefulsetname>-<podindex>. For example, a pod in a StatefulSet named
	// "web" with index number "3" would be named "web-3".
	// The only allowed template.spec.restartPolicy value is "Always".
	template: #PodPartialTemplateSpec
	// updateStrategy indicates the StatefulSetUpdateStrategy that will be
	// employed to update Pods in the StatefulSet when a revision is made to
	// Template.
	updateStrategy?: #StatefulSetUpdateStrategy
	// volumeClaimTemplates is a list of claims that pods are allowed to reference.
	// The StatefulSet controller is responsible for mapping network identities to
	// claims in a way that maintains the identity of a pod. Every claim in
	// this list must have at least one matching (by name) volumeMount in one
	// container in the template. A claim in this list takes precedence over
	// any volumes in the template, with the same name.
	// TODO: Define the behavior if a claim already exists with the same name.
	volumeClaimTemplates?: [...#PersistentVolumeClaim]
}

#StatefulSetUpdateStrategy: {
	// RollingUpdate is used to communicate parameters when Type is RollingUpdateStatefulSetStrategyType.
	rollingUpdate?: #RollingUpdateStatefulSetStrategy
	// Type indicates the type of the StatefulSetUpdateStrategy. Default is RollingUpdate.
	type?: #StatefulSetUpdateStrategyType
}

#StatefulSetUpdateStrategyType: string

#Status: {
	digests?: [...#DigestMeta]
	endpoint?: [X=string]: string
	images?: [X=string]: string
	resources?: [...{
		[X=string]: _
	}]
}

#StorageMedium: string

#SuccessPolicy: {
	// rules represents the list of alternative rules for the declaring the Jobs as successful before `.status.succeeded >= .spec.completions`. Once any of the rules are met, the "SucceededCriteriaMet" condition is added, and the lingering pods are removed. The terminal state for such a Job has the "Complete" condition. Additionally, these rules are evaluated in order; Once the Job meets one of the rules, other rules are ignored. At most 20 elements are allowed.
	rules: [...#SuccessPolicyRule]
}

#SuccessPolicyRule: {
	// succeededCount specifies the minimal required size of the actual set of the succeeded indexes for the Job. When succeededCount is used along with succeededIndexes, the check is constrained only to the set of indexes specified by succeededIndexes. For example, given that succeededIndexes is "1-4", succeededCount is "3", and completed indexes are "1", "3", and "5", the Job isn't declared as succeeded because only "1" and "3" indexes are considered in that rules. When this field is null, this doesn't default to any value and is never evaluated at any time. When specified it needs to be a positive integer.
	succeededCount?: int
	// succeededIndexes specifies the set of indexes which need to be contained in the actual set of the succeeded indexes for the Job. The list of indexes must be within 0 to ".spec.completions-1" and must not contain duplicates. At least one element is required. The indexes are represented as intervals separated by commas. The intervals can be a decimal integer or a pair of decimal integers separated by a hyphen. The number are listed in represented by the first and last element of the series, separated by a hyphen. For example, if the completed indexes are 1, 3, 4, 5 and 7, they are represented as "1,3-5,7". When this field is null, this field doesn't default to any value and is never evaluated at any time.
	succeededIndexes?: string
}

#Sysctl: {
	// Name of a property to set
	name: string
	// Value of a property to set
	value: string
}

#TCPSocketAction: {
	// Optional: Host name to connect to, defaults to the pod IP.
	host?: string
	// Number or name of the port to access on the container. Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.
	port: #IntOrString
}

#TaintEffect: string

#TerminationMessagePolicy: string

#Time: string

#Toleration: {
	// Effect indicates the taint effect to match. Empty means match all taint effects. When specified, allowed values are NoSchedule, PreferNoSchedule and NoExecute.
	effect?: #TaintEffect
	// Key is the taint key that the toleration applies to. Empty means match all taint keys. If the key is empty, operator must be Exists; this combination means to match all values and all keys.
	key?: string
	// Operator represents a key's relationship to the value. Valid operators are Exists and Equal. Defaults to Equal. Exists is equivalent to wildcard for value, so that a pod can tolerate all taints of a particular category.
	operator?: #TolerationOperator
	// TolerationSeconds represents the period of time the toleration (which must be of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default, it is not set, which means tolerate the taint forever (do not evict). Zero and negative values will be treated as 0 (evict immediately) by the system.
	tolerationSeconds?: int
	// Value is the taint value the toleration matches to. If the operator is Exists, the value should be empty, otherwise just a regular string.
	value?: string
}

#TolerationOperator: string

#TopologySpreadConstraint: {
	// LabelSelector is used to find matching pods. Pods that match this label selector are counted to determine the number of pods in their corresponding topology domain.
	labelSelector?: #LabelSelector
	// MatchLabelKeys is a set of pod label keys to select the pods over which spreading will be calculated. The keys are used to lookup values from the incoming pod labels, those key-value labels are ANDed with labelSelector to select the group of existing pods over which spreading will be calculated for the incoming pod. The same key is forbidden to exist in both MatchLabelKeys and LabelSelector. MatchLabelKeys cannot be set when LabelSelector isn't set. Keys that don't exist in the incoming pod labels will be ignored. A null or empty list means only match against labelSelector.
	// 
	// This is a beta field and requires the MatchLabelKeysInPodTopologySpread feature gate to be enabled (enabled by default).
	matchLabelKeys?: [...string]
	// MaxSkew describes the degree to which pods may be unevenly distributed. When `whenUnsatisfiable=DoNotSchedule`, it is the maximum permitted difference between the number of matching pods in the target topology and the global minimum. The global minimum is the minimum number of matching pods in an eligible domain or zero if the number of eligible domains is less than MinDomains. For example, in a 3-zone cluster, MaxSkew is set to 1, and pods with the same labelSelector spread as 2/2/1: In this case, the global minimum is 1. 
	maxSkew: int
	// MinDomains indicates a minimum number of eligible domains. When the number of eligible domains with matching topology keys is less than minDomains, Pod Topology Spread treats "global minimum" as 0, and then the calculation of Skew is performed. And when the number of eligible domains with matching topology keys equals or greater than minDomains, this value has no effect on scheduling. As a result, when the number of eligible domains is less than minDomains, scheduler won't schedule more than maxSkew Pods to those domains. If value is nil, the constraint behaves as if MinDomains is equal to 1. Valid values are integers greater than 0. When value is not nil, WhenUnsatisfiable must be DoNotSchedule.
	// 
	// For example, in a 3-zone cluster, MaxSkew is set to 2, MinDomains is set to 5 and pods with the same labelSelector spread as 2/2/2: 
	minDomains?: int
	// NodeAffinityPolicy indicates how we will treat Pod's nodeAffinity/nodeSelector when calculating pod topology spread skew. Options are: - Honor: only nodes matching nodeAffinity/nodeSelector are included in the calculations. - Ignore: nodeAffinity/nodeSelector are ignored. All nodes are included in the calculations.
	// 
	// If this value is nil, the behavior is equivalent to the Honor policy. This is a beta-level feature default enabled by the NodeInclusionPolicyInPodTopologySpread feature flag.
	nodeAffinityPolicy?: #NodeInclusionPolicy
	// NodeTaintsPolicy indicates how we will treat node taints when calculating pod topology spread skew. Options are: - Honor: nodes without taints, along with tainted nodes for which the incoming pod has a toleration, are included. - Ignore: node taints are ignored. All nodes are included.
	// 
	// If this value is nil, the behavior is equivalent to the Ignore policy. This is a beta-level feature default enabled by the NodeInclusionPolicyInPodTopologySpread feature flag.
	nodeTaintsPolicy?: #NodeInclusionPolicy
	// TopologyKey is the key of node labels. Nodes that have a label with this key and identical values are considered to be in the same topology. We consider each <key, value> as a "bucket", and try to put balanced number of pods into each bucket. We define a domain as a particular instance of a topology. Also, we define an eligible domain as a domain whose nodes meet the requirements of nodeAffinityPolicy and nodeTaintsPolicy. e.g. If TopologyKey is "kubernetes.io/hostname", each Node is a domain of that topology. And, if TopologyKey is "topology.kubernetes.io/zone", each zone is a domain of that topology. It's a required field.
	topologyKey: string
	// WhenUnsatisfiable indicates how to deal with a pod if it doesn't satisfy the spread constraint. - DoNotSchedule (default) tells the scheduler not to schedule it. - ScheduleAnyway tells the scheduler to schedule the pod in any location,
	//   but giving higher precedence to topologies that would help reduce the
	//   skew.
	// A constraint is considered "Unsatisfiable" for an incoming pod if and only if every possible node assignment for that pod would violate "MaxSkew" on some topology. For example, in a 3-zone cluster, MaxSkew is set to 1, and pods with the same labelSelector spread as 3/1/1: 
	whenUnsatisfiable: #UnsatisfiableConstraintAction
}

#TypedLocalObjectReference: {
	// APIGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must be in the core API group. For any other third-party types, APIGroup is required.
	apiGroup: string
	// Kind is the type of resource being referenced
	kind: string
	// Name is the name of resource being referenced
	name: string
}

#TypedObjectReference: {
	// APIGroup is the group for the resource being referenced. If APIGroup is not specified, the specified Kind must be in the core API group. For any other third-party types, APIGroup is required.
	apiGroup: string
	// Kind is the type of resource being referenced
	kind: string
	// Name is the name of resource being referenced
	name: string
	// Namespace is the namespace of resource being referenced Note that when a namespace is specified, a gateway.networking.k8s.io/ReferenceGrant object is required in the referent namespace to allow that namespace's owner to accept the reference. See the ReferenceGrant documentation for details. (Alpha) This field requires the CrossNamespaceVolumeDataSource feature gate to be enabled.
	namespace?: string
}

#UID: string

#URIScheme: string

#UnsatisfiableConstraintAction: string

#Volume: #VolumeConfigMap | #VolumeEmptyDir | #VolumeHostPath | #VolumePersistentVolumeClaim | #VolumeSecret

#VolumeConfigMap: {
	mountPath:         string
	mountPropagation?: "Bidirectional" | "HostToContainer"
	opt?:              #ConfigMapVolumeSource
	optional?:         bool
	// Prefix mountPath == export, use as envFrom
	prefix?: string
	// else volumeMounts
	readOnly?: bool
	spec?:     #ConfigMapSpec
	subPath?:  string
	type:      "ConfigMap"
}

#VolumeEmptyDir: {
	mountPath:         string
	mountPropagation?: "Bidirectional" | "HostToContainer"
	opt?:              #EmptyDirVolumeSource
	optional?:         bool
	// Prefix mountPath == export, use as envFrom
	prefix?: string
	// else volumeMounts
	readOnly?: bool
	subPath?:  string
	type:      "EmptyDir"
}

#VolumeHostPath: {
	mountPath:         string
	mountPropagation?: "Bidirectional" | "HostToContainer"
	opt?:              #HostPathVolumeSource
	optional?:         bool
	// Prefix mountPath == export, use as envFrom
	prefix?: string
	// else volumeMounts
	readOnly?: bool
	subPath?:  string
	type:      "HostPath"
}

#VolumePersistentVolumeClaim: {
	mountPath:         string
	mountPropagation?: "Bidirectional" | "HostToContainer"
	opt?:              #PersistentVolumeClaimVolumeSource
	optional?:         bool
	// Prefix mountPath == export, use as envFrom
	prefix?: string
	// else volumeMounts
	readOnly?: bool
	spec:      #PersistentVolumeClaimSpec
	subPath?:  string
	type:      "PersistentVolumeClaim"
}

#VolumeResourceRequirements: {
	// Limits describes the maximum amount of compute resources allowed. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	limits?: #ResourceList
	// Requests describes the minimum amount of compute resources required. If Requests is omitted for a container, it defaults to Limits if that is explicitly specified, otherwise to an implementation-defined value. Requests cannot exceed Limits. More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	requests?: #ResourceList
}

#VolumeSecret: {
	mountPath:         string
	mountPropagation?: "Bidirectional" | "HostToContainer"
	opt?:              #SecretVolumeSource
	optional?:         bool
	// Prefix mountPath == export, use as envFrom
	prefix?: string
	// else volumeMounts
	readOnly?: bool
	spec?:     #ConfigMapSpec
	subPath?:  string
	type:      "Secret"
}

#WeightedPodAffinityTerm: {
	// Required. A pod affinity term, associated with the corresponding weight.
	podAffinityTerm: #PodAffinityTerm
	// weight associated with matching the corresponding podAffinityTerm, in the range 1-100.
	weight: int
}

#WindowsSecurityContextOptions: {
	// GMSACredentialSpec is where the GMSA admission webhook (https://github.com/kubernetes-sigs/windows-gmsa) inlines the contents of the GMSA credential spec named by the GMSACredentialSpecName field.
	gmsaCredentialSpec?: string
	// GMSACredentialSpecName is the name of the GMSA credential spec to use.
	gmsaCredentialSpecName?: string
	// HostProcess determines if a container should be run as a 'Host Process' container. All of a Pod's containers must have the same effective HostProcess value (it is not allowed to have a mix of HostProcess containers and non-HostProcess containers). In addition, if HostProcess is true then HostNetwork must also be set to true.
	hostProcess?: bool
	// The UserName in Windows to run the entrypoint of the container process. Defaults to the user specified in image metadata if unspecified. May also be set in PodSecurityContext. If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
	runAsUserName?: string
}
