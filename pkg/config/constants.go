package config

import "regexp"

const (
	// ENV_KW_CONFIG_REPO is the name of the environment variable which contains the path to the kubeswitcher config repo
	ENV_KW_CONFIG_REPO = "KW_CONFIG_REPO"
	// KW_CONFIG_REPO_DEFAULT_NAME is the default name of the folder containing the kubeswitcher config files
	KW_CONFIG_REPO_DEFAULT_NAME = ".kubeswitcher_config"
	// ENV_KW_SESSION_ID is the name of the environment variable which contains the kubeswitcher-specific session id
	ENV_KW_SESSION_ID = "KW_SESSION_ID"
	// ENV_DEFAULT_SESSION_ID is the name of the environment variable which is used if the env var from ENV_KW_SESSION_ID is not set
	ENV_DEFAULT_SESSION_ID = "TERM_SESSION_ID"
)

var (
	SIDRegex  = regexp.MustCompile(`^[\w-]{1,128}$`)
	UUIDRegex = regexp.MustCompile(`([a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12})`)
)
