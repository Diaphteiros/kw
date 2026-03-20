package context

const (
	// Name of the env var pointing to the kubectl binary
	ENV_VAR_KUBECTL_PATH = "KUBESWITCHER_KUBECTL_PATH"

	// Name of the env var pointing to the kubeconfig file
	ENV_VAR_KUBECONFIG_PATH = "KUBESWITCHER_KUBECONFIG_PATH"

	// Name of the env var containing the name of the currently executed plugin
	ENV_VAR_CURRENT_PLUGIN_NAME = "KUBESWITCHER_CURRENT_PLUGIN_NAME"

	// Name of the env var containing the path to the generic state file
	ENV_VAR_GENERIC_STATE_PATH = "KUBESWITCHER_GENERIC_STATE_PATH"

	// Name of the env var containing the path to the plugin state file
	ENV_VAR_PLUGIN_STATE_PATH = "KUBESWITCHER_PLUGIN_STATE_PATH"

	// Name of the env var containing the path to the notification message file
	ENV_VAR_NOTIFICATION_MESSAGE_PATH = "KUBESWITCHER_NOTIFICATION_MESSAGE_PATH"

	// Name of the env var containing the path to the file with the id string.
	ENV_VAR_ID_PATH = "KUBESWITCHER_ID_PATH"

	// Name of the env var containing the path to the internal call file.
	ENV_VAR_INTERNAL_CALL_PATH = "KUBESWITCHER_INTERNAL_CALL_PATH"

	// Name of the env var containing the path to the internal callback request file.
	ENV_VAR_INTERNAL_CALLBACK_REQUEST_PATH = "KUBESWITCHER_INTERNAL_CALLBACK_REQUEST_PATH"

	// Name of the env var containing the path to the internal callback state file.
	ENV_VAR_INTERNAL_CALLBACK_STATE_PATH = "KUBESWITCHER_INTERNAL_CALLBACK_STATE_PATH"

	// Name of the env var containing the statically defined plugin configuration.
	ENV_VAR_PLUGIN_CONFIG = "KUBESWITCHER_PLUGIN_CONFIG"

	// Name of the env var containing the current session id.
	ENV_VAR_SESSION_ID = "KUBESWITCHER_SESSION_ID"

	// Name of the env var containing the current session directory.
	ENV_VAR_SESSION_CONFIG_DIR = "KUBESWITCHER_SESSION_CONFIG_DIR"

	// Name of the env var containing the path to the kubeswitcher config directory.
	ENV_VAR_CONFIG_DIR = "KUBESWITCHER_CONFIG_DIR"

	// Name of the env var containing the debug flag
	ENV_VAR_DEBUG = "KUBESWITCHER_FLAG_DEBUG"
)
