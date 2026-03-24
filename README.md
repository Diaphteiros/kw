# KubeSwitcher

KubeSwitcher is a small CLI tool for switching between multiple kubeconfig files. It is useful for working with multiple kubernetes clusters at the same time.

Opposed to similar tools, e.g. [kubeswitch](https://github.com/danielfoehrKn/kubeswitch), kubeswitcher is optimized for working with multiple kubeconfig files that each target only a single cluster instead of having big kubeconfig files with multiple contexts for multiple clusters.

## Installation

To install the KubeSwitcher tool, simply run the following command
```shell
go install github.com/Diaphteiros/kw@latest
```
or clone the repository and run
```shell
task install
```

> [!NOTE]
> For technical reasons, `go install github.com/Diaphteiros/kw` works only for `latest` and release versions.
> For specific commits, the tool has to be checked out and built manually via `task install`. Probably, `go work use . ./pluginlib` is required to be run before as well.

> [!NOTE]
> This project uses [task](https://taskfile.dev/) instead of `make`.

### The 'KUBECONFIG' Env Var

KubeSwitcher is inspired heavily from [gardenctl](https://github.com/gardener/gardenctl-v2). It uses the same concept of modifying a session-specific kubeconfig file in temporary folder, which means that it requires the `KUBECONFIG` env var to point to this file for `kubectl` to target the selected cluster without having to specify the `--kubeconfig` argument. The `kw kubectl-env` subcommand can be used to generate a code snippet that sets the `KUBECONFIG` env var accordingly. It is recommended to source the generated snippet in the `.bashrc` file (or its equivalent) to ensure that `kubectl` always uses the correct kubeconfig by default.

#### KubeSwitcher and Gardenctl

As explained above, both, kubeswitcher and gardenctl, expect to 'own' the kubeconfig file, which makes it difficult to use both tools in parallel. To avoid this issue, there is a plugin for kubeswitcher (see [below](#known-plugins)) which basically wraps gardenctl and allows to call it via a kubeswitcher subcommand, integrating it into the kubeswitcher experience.

#### Generating a Prompt

Similarly to the `kw kubectl-env` command, there is a `kw prompt` command which generates a code snippet that can be sourced to make a `_kw_prompt` shell function available, which, when called, returns a short string describing the currently selected kubeconfig's target. This is designed to be used as part of a shell prompt so that users always know which cluster they are currently targeting.

## Configuration

KubeSwitcher requires a configuration file. Any KubeSwitcher command will create it with default values, if it does not exist.
The file is expected at the location returned by `kw config path`.
```yaml
kubeswitcher:
  printInfoOnKubeconfigChange: true
  printChangeInfoToStderr: false
  historyDepth: 5
builtin:
  custom:
    maxIdLength: 30
plugins: []
```
- `kubeswitcher` - General kubeswitcher configuration.
  - `printInfoOnKubeconfigChange` (_optional_, default `true`) - Whether a message should be printed if the kubeconfig changed as a result of a kubeswitcher command.
  - `printChangeInfoToStderr` (_optional_, default `false`) - Wether the message should be printed to stderr instead of stdout.
  - `historyDepth` (_optional_, default `5`) - How many history entries should be stored. `0` will deactivate the history, which renders some commands unusable. The maximum value is `50`.
- `builtin` (_optional_) - Configuration specific to built-in subcommands.
  - `custom` (_optional_) - Configuration for the built-in `custom` subcommand.
    - `maxIdLength` (_optional_, default `0`) - The maximum length of the path contained in the id created by the `custom` subcommand. The generated id usually looks like `custom:<absolute_path_to_kubeconfig>`, which can become a bit too long to be used in a shell prompt (which is one of the purposes of the id). If this value is set to a positive number, the beginning of the path to the kubeconfig will be truncated if the path exceeds the maximum length.
      - Example: `/foo/bar/baz` with a max length of `6` results in the id `custom:…ar/baz`.
- `plugins` (_optional_) - Registered plugins. See [below](#plugin-configuration) for more information.

See also the [`config` subcommand](./docs/reference/kw_config.md) for further information.

## Command Reference

Run with the `--help` flag or look [here](./docs/reference/kw.md).

### Usage Examples

Here are some examples for the more commonly used kubeswitcher commands:

> Tip: Most subcommands have aliases that are significantly shorter. Call the command with the `--help` flag to find out what they are. More often than not, one of the available aliases is simply the first letter of the respective subcommand.

---
```shell
kw custom ./my-local-kubeconfig.yaml
```
The `kw custom` (or short `kw c`) subcommand simply switches to a local kubeconfig file.

---
```shell
kw temporary
```
Use this to switch to a kubeconfig currently contained in the clipboard.

---
```shell
kw bookmark save <name>
```
This stores the currently selected kubeconfig under the specified name. The name can be omitted, in which case the tool will prompt you for it.
```shell
kw bookmark load <name>
```
can be used to load a previously bookmarked kubeconfig. If the name is omitted, kubeswitcher provides a list of all bookmarked kubeconfigs to choose from.

Use
```shell
kw bookmark view
```
can be used to view all bookmarks.

---
```shell
kw history view
```
shows the last used kubeconfigs (if the history depth in the [configuration](#configuration) is greater than 0). It takes a `--global` (or `-g`) flag to show the history across terminal sessions, without it will only show the current terminal session's history.

```shell
kw history load <index>
```
can be used to switch back to a recently used kubeconfig. The index starts at `0`, pointing to the currently selected kubeconfig.

Two convenience aliases for loading kubeconfigs from history are
```shell
kw flip
```
which is an alias for `kw history load 1` and basically flips between the current kubeconfig and the one used before when called repeatedly, and
```shell
kw repeat
```
which is an alias for `kw history load 0 --global`, setting this terminal session's kubeconfig to the same kubeconfig that was selected last in any other terminal session (has no effect if used in the same terminal session that was the last to switch to a kubeconfig).

---
```shell
kw namespace <name>
```
(or `kw ns`) is the only default commands which actually modifies a kubeconfig by injecting the specified namespace (which can be left empty to be prompted for it) as the default namespace in the current context. This means that all `kubectl` commands without the `--namespace`/`-n` option will then target this namespace.

Note that the namespace is not persisted in the history, so if a kubeconfig is loaded from history or bookmark, its default namespace will always be unset, which usually leads to `kubectl` using `default` as the default namespace.


## Extensibility

K8s clusters can be created in many different ways, e.g. GKE, Gardener, kind, and so on. Instead of trying to handle all of these within a single binary, kubeswitcher is designed in an extensible manner and allows plugins to be called for specific subcommands. This way, the main binary will only contain a set of generic subcommands, with any provider-/landscape-specific subcommand being redirected to the corresponding plugin.

The [contract documentation](docs/development/contract.md) describes the contract that is expected from a plugin in more detail.

### Known Plugins

| **Name** | **Description** |
| --- | --- |
| [garden](https://github.com/Diaphteiros/kw_garden) | Switch to clusters of a Gardener landscape. |
| [kind](https://github.com/Diaphteiros/kw_kind) | Switch to local kind clusters. |


