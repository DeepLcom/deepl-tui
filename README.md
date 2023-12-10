# DeepL TUI

`deepl-tui` is a terminal user interface for the DeepL language translation API.

![](./assets/demo/demo.gif)

## Installation

To install `deepl-tui` you can download a [prebuilt binary][prebuilt-binaries]
that matches your system and place it in a directory that's part of your
system's search path, e.g.
```shell
# system information
OS=linux
ARCH=amd64

# install dir (must exist)
BIN_DIR=$HOME/.local/bin

# download latest release into install dir
RELEASE_TAG=$(curl -sSfL https://api.github.com/repos/DeepLcom/deepl-tui/releases/latest | jq -r '.tag_name')
curl \
    -L https://github.com/DeepLcom/deepl-tui/releases/download/${RELEASE_TAG}/deepl-tui_${RELEASE_TAG}_${OS}_${ARCH} \
    -o ${BIN_DIR}/deepl-tui

# make it executable
chmod +x ${BIN_DIR}/deepl-tui
```

Alternatively, if you have the [Go tools][go-install] installed, you can use
```shell
go install github.com/DeepLcom/deepl-tui@latest
```

## Usage

Since `deepl-tui` uses the DeepL API you'll need an API authentication key.
To get a key, please [create an account here][create-account]. With a DeepL API 
Free account you can translate up to 500,000 characters/month.

You can either pass the authentication key as an environment variable or via
the `--auth-key` option, i.e.
```shell
$ export DEEPL_AUTH_KEY="f63c02c5-f056..."  # replace with your key
$ deepl-tui
```
or
```shell
$ deepl-tui --auth-key=f63c02c5-f056...
```

### Key bindings

| Action                         | Keys     | Comment                     |
| ---                            | ---      | ---                         |
| Focus input text area          | `alt-i`  |                             |
| Focus source language dropdown | `alt-s`  | Hit `enter` to list options |
| Focus target language dropdown | `alt-t`  | Hit `enter` to list options |
| Quit the application           | `ctrl-q` |                             |

## Changelog

Notable changes to this project will be documented in the [CHANGELOG.md](./CHANGELOG.md)

## License

This project is released under the [MIT License](./LICENSE).

<!-- Links -->
[prebuilt-binaries]: https://github.com/DeepLcom/deepl-tui/releases/latest
[go-install]: https://go.dev/doc/install
[create-account]: https://www.deepl.com/pro#developer
