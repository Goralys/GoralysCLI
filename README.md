## Goralys CLI

This CLI tool is made only for the Goralys project. To
learn more about this project you can check out the
[monorepo](https://github.com/SAMSAM-55/Goralys).

### Presentation

This CLI was first made to avoid the hassle of writing
bash/batch scripts for both Linux and Windows separately.

The **Go** (or **GoLang**) programming language was chosen for its ease
of use and its ability to cross-compile binaries for multiple
platforms. It was also chosen because of the **Cobra** library
that allows to easily make CLI tools in **Go**.

### Setup

To set up the project, simply run the following commands:

```bash
# downloads the dependencies
go mod tidy

# builds the project (add the .exe extension for windows)
go build -o build/goralys-cli
```

Then you can simply run the tool:

```bash
# linux
./build/goralys-cli --help

# windows
.\build\goralys-cli --help
```

### Contact

For any request about this project please use the following
address:

Sami Saubion <sami.saubion@dev.goralys.fr>

– Main developer of the Goralys project 