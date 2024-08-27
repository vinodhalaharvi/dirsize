# dirsize

`dirsize` is a Go utility for calculating directory sizes. It provides a simple command-line interface to measure the size of directories, with options for recursive calculation and human-readable output.

## Features

- Calculate the size of one or more directories
- Recursive mode for including subdirectories
- Human-readable output option
- Concurrent processing for improved performance

## Installation

### Prerequisites

- Go 1.x or higher

### Building from source

1. Clone the repository:
   ```
   git clone https://github.com/vinodhalaharvi/dirsize.git
   cd dirsize
   ```

2. Build the project:
   ```
   make build
   ```

3. (Optional) Install golangci-lint and other tools:
   ```
   make install-tools
   ```

## Usage

After building the project, you can run `dirsize` as follows:

```
./dirsize [flags] [directories...]
```

### Flags

- `-recursive`: Calculate sizes recursively (include subdirectories)
- `-human`: Display sizes in human-readable format

### Examples

Calculate the size of the current directory:
```
./dirsize .
```

Calculate the sizes of multiple directories recursively with human-readable output:
```
./dirsize -recursive -human /path/to/dir1 /path/to/dir2
```

## Development

### Running Tests

To run the test suite:

```
make test
```

### Linting

To run the linter:

```
make lint
```

### Running All Quality Checks

To run all quality checks (lint, vet, and test):

```
make quality
```

## Contributing

Contributions to `dirsize` are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [MIT License](LICENSE).

