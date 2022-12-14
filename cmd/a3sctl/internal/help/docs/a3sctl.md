## a3sctl

a3sctl is a command line interface that allows to easily and quickly interact
with a running instance of an a3s server. It allows:

### Commands

a3sctl provides the following commands:

* `auth`: retrieve an identity token.
* `api`: interact with the api to create/get/update/delete resources.
* `completion`: dumps completion definitions for various shell.

### Global flags

The following flags are propagated to all all sub commands:

* `--config`: path to an optional configuration file
* `--log-level`: debug, info, warn or error.
* `--help` or `-h`: display help for the current command.

Every flag can be read from multiple places, following the process:

* use the posix flag if provided, otherwise,
* read the associated env var if set, otherwise,
* read from the config file if set, otherwise,
* read from the flag default if set, otherwise,
* error.

The rule to translate a flag to an environment variables is the following:

* remove `--` suffix
* replace `-` by `_`
* uppercase everything
* add prefix `A3SCTL_`

For instance, `--namespace` will become `$A3SCTL_NAMESPACE` and `--source-name`
will become `$A3SCTL_SOURCE_NAME`.

### Completions

a3sctl supports shell autocompletions for a range of shells.

To enable autocompletion for bash:

    . <(a3sctl completion bash)

For zsh:

    compdef _a3sctl a3sctl
    . <(a3sctl completion zsh)

For Fish:

    . <(a3sctl completion fish)

For a more permanent completion definition information, refer to your shell
documentation.

### Configuration file

a3sctl can use user defined default from a configuration file. It will look for
a configuration file in the following folders, using the first one it finds:

* `~/.config/a3sctl/default.yaml`
* `/usr/local/etc/default.yaml`
* `/etc/default.yaml`

The configuration file can be written in various markup language. Just use
appropriate extentensions, like `json` or `toml`.

You can have more than one configuration file in your conf folder. To control
which one to use, you can use `--config-name` or set `A3SCTL_CONFIG_NAME`.

For instance if you have `dev.yaml` and `prod.toml`, to use `dev.yaml`, use
either:

    a3sctl --config-name dev ...
    A3SCTL_CONFIG_NAME=dev a3sctl ...

And to use `prod.toml`, use either:

    a3sctl --config-name prod ...
    A3SCTL_CONFIG_NAME=prod a3sctl ...
