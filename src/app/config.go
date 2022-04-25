package app

import "sync"

type Project struct {
	Version     string    `yaml:"version"`
	LogLocation string    `yaml:"log_location,omitempty"`
	LogLevel    string    `yaml:"log_level,omitempty"`
	Processes   Processes `yaml:"processes"`
	Environment []string  `yaml:"environment,omitempty"`

	runningProcesses map[string]*Process
	processStates    map[string]*ProcessState
	mapMutex         sync.Mutex
}

type Processes map[string]ProcessConfig
type ProcessConfig struct {
	Name          string
	Disabled      bool                   `yaml:"disabled,omitempty"`
	Command       string                 `yaml:"command"`
	LogLocation   string                 `yaml:"log_location,omitempty"`
	Environment   []string               `yaml:"environment,omitempty"`
	RestartPolicy RestartPolicyConfig    `yaml:"availability,omitempty"`
	DependsOn     DependsOnConfig        `yaml:"depends_on,omitempty"`
	Extensions    map[string]interface{} `yaml:",inline"`
}

type ProcessState struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Restarts int    `json:"restarts"`
	ExitCode int    `json:"exit_code"`
}

func (p ProcessConfig) GetDependencies() []string {
	dependencies := make([]string, len(p.DependsOn))

	i := 0
	for k := range p.DependsOn {
		dependencies[i] = k
		i++
	}
	return dependencies
}

const (
	RestartPolicyAlways    = "always"
	RestartPolicyOnFailure = "on-failure"
	RestartPolicyNo        = "no"
)

const (
	ProcessStatePending    = "Pending"
	ProcessStateRunning    = "Running"
	ProcessStateRestarting = "Restarting"
	ProcessStateCompleted  = "Completed"
)

type RestartPolicyConfig struct {
	Restart        string `yaml:",omitempty"`
	BackoffSeconds int    `yaml:"backoff_seconds,omitempty"`
	MaxRestarts    int    `yaml:"max_restarts,omitempty"`
}

const (
	// ProcessConditionCompleted is the type for waiting until a process has completed (any exit code).
	ProcessConditionCompleted = "process_completed"

	// ProcessConditionCompletedSuccessfully is the type for waiting until a process has completed successfully (exit code 0).
	ProcessConditionCompletedSuccessfully = "process_completed_successfully"

	// ProcessConditionHealthy is the type for waiting until a process is healthy.
	ProcessConditionHealthy = "process_healthy"

	// ProcessConditionStarted is the type for waiting until a process has started (default).
	ProcessConditionStarted = "process_started"
)

type DependsOnConfig map[string]ProcessDependency

type ProcessDependency struct {
	Condition  string                 `yaml:",omitempty"`
	Extensions map[string]interface{} `yaml:",inline"`
}