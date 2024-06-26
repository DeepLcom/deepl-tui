# DeepL TUI

`deepl-tui` is a terminal user interface for the DeepL language translation API.

![](./assets/demo/demo.gif)

## Changelog

Notable changes to this project will be documented in the [CHANGELOG.md](./CHANGELOG.md)

## Installation

To install `deepl-tui` you can download a [prebuilt binary][prebuilt-binaries]
that matches your system and place it in a directory that's part of your
system's search path, e.g.
```shell
# download latest release archive
RELEASE_TAG=$(curl -sSfL https://api.github.com/repos/DeepLcom/deepl-tui/releases/latest | jq -r '.tag_name')
curl -sSfL -o /tmp/deepl-tui.tar.gz \
    https://github.com/DeepLcom/deepl-tui/releases/download/${RELEASE_TAG}/deepl-tui_${RELEASE_TAG}_linux_amd64.tar.gz

# extract executable binary into install dir (must exist)
INSTALL_DIR=$HOME/.local/bin
tar -C ${INSTALL_DIR} -zxof /tmp/deepl-tui.tar.gz deepl-tui
```

## Usage

Since `deepl-tui` uses the DeepL API you'll need an API authentication key.
To get a key, please [create an account here][create-account]. With a DeepL API 
Free account you can translate up to 500,000 characters per month.

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

#### Global

| Action               | Keys      | Comment |
| ---                  | ---       | ---     |
| Cycle through pages  | `alt-tab` |         |
| Quit the application | `ctrl-q`  |         |

#### Translate Page

| Action                          | Keys    | Comment                     |
| ---                             | ---     | ---                         |
| Focus input text area           | `alt-i` |                             |
| Focus source language dropdown  | `alt-s` | Hit `enter` to list options |
| Focus target language dropdown  | `alt-t` | Hit `enter` to list options |
| Focus formality option dropdown | `alt-f` | Hit `enter` to list options |
| Focus glossary option button    | `alt-g` | Hit `enter` to open dialog  |

#### Glossaries Page

| Action                       | Keys    | Comment |
| ---                          | ---     | ---     |
| Focus glossary entry form    | `alt-e` |         |
| Focus glossary info form     | `alt-i` |         |
| Focus glossaries list        | `alt-l` |         |
| Focus glossary entries table | `alt-t` |         |

## License

This project is released under the [MIT License](./LICENSE).

<!-- Links -->
[prebuilt-binaries]: https://github.com/DeepLcom/deepl-tui/releases/latest
[create-account]: https://www.deepl.com/pro#developer
