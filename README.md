# `command-function-buildpack`

The Command Function Buildpack is a Cloud Native Buildpack V3 that provides riff [Command Function Invoker](https://github.com/projectriff/command-function-invoker) to functions.

This buildpack is designed to work in collaboration with other buildpacks, which are tailored to
support (and know how to build / run) languages supported by riff.

## In Plain English

In a nutshell, when combined with the other buildpacks present in the [riff builder](https://github.com/projectriff/builder) what this means (and especially when dealing with the riff CLI which takes care of the creation of the `riff.toml` file for you):

- The fact that the `--artifact` flag points to a file with the execute permission will result in the execution as a command function, thanks to the [command invoker](https://github.com/projectriff/command-function-invoker)
- Ambiguity in the detection process will result in a build failure
- The presence of the `--invoker` flag will entirely bypass the detection mechanism and force a given language/invoker

## Detailed Buildpack Behavior

### Detection Phase

Detection passes if

- a `$APPLICATION_ROOT/riff.toml` exists and
- the file pointed to by the `artifact` value in `riff.toml` exists and is executable

If detection passes, the buildpack will add a `riff-invoker-command` key and `command`
metadata extracted from the riff metadata.

If several languages are detected simultaneously, the detect phase errors out.
The `override` key in `riff.toml` can be used to bypass detection and force the use of a particular invoker.

### Build Phase

If a command function has been selected

- Contributes the riff Command Invoker to a launch layer, set as the main executable with `FUNCTION_URI = <artifact>` set as an environment variable.

The function behavior is exposed _via_ standard buildpack [process types](https://github.com/buildpack/spec/blob/master/buildpack.md#launch):

- Contributes `web` process
- Contributes `function` process

## How to Build

### Prerequisites
To build the command-function-buildpack you'll need

- Go 1.13+
- to run acceptance tests:
  - a running local docker daemon
- for local function builds
  - the [`pack`](https://github.com/buildpacks/pack) command line tool


You can build the buildpack by running

```bash
make
```

This will package (with pre-downloaded cache layers) the buildpack in the
`artifactory/io/projectriff/command/io.projectriff.command/latest` directory. That can be used as a `uri` in a `builder.toml`
file of a builder (see https://github.com/projectriff/builder)

## License

This buildpack is released under version 2.0 of the [Apache License](https://www.apache.org/licenses/LICENSE-2.0).
