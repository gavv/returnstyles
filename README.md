# returnstyles [![GoDev](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/gavv/returnstyles) [![Build](https://github.com/gavv/returnstyles/workflows/build/badge.svg)](https://github.com/gavv/returnstyles/actions) [![GitHub release](https://img.shields.io/github/release/gavv/returnstyles.svg)](https://github.com/gavv/returnstyles/releases)

Go (golang) linter to check function return styles according to rules of your choice.

## Install

```
go install -v github.com/gavv/returnstyles/cmd/returnstyles@latest
```

## Usage

```
returnstyles [options] <package>
```

For example:

```
returnstyles -allow-naked-returns=false .
```

## Options

| option                   | default | description                                                               |
|--------------------------|---------|---------------------------------------------------------------------------|
| `-allow-unnamed`         | true    | allow functions with unnamed return variables                             |
| `-allow-named`           | true    | allow functions with named return variables                               |
| `-allow-partially-named` | false   | allow functions with partially named return variables                     |
| `-allow-normal-returns`  | true    | allow normal (non-naked) returns in functions with named return variables |
| `-allow-naked-returns`   | true    | allow naked returns in functions with named return variables              |
| `-allow-mixing-returns`  | false   | allow mixing normal and naked in functions with named return variables    |
| `-include-cgo`           | false   | include diagnostics for cgo-generated functions                                                       |

## Config

Instead of specifying individual command-line options, you can provide a single configuration file in YAML format:

```
returnstyles -config path/to/config.yaml <package>
```

It should have the following format:

```yaml
returnstyles:
  allow-unnamed: true
  allow-named: true
  allow-partially-named: false
  allow-normal-returns: true
  allow-naked-returns: true
  allow-mixing-returns: false
  include-cgo: false
```

## Package

```go
import "github.com/gavv/returnstyles"
```

Global variable `returnstyles.Analyzer` follows guidelines in the [golang.org/x/tools/go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis) package.

## Samples

#### `-allow-unnamed=false`

```go
func foo() (int, int) { // lint: functions with unnamed return variables not allowed
    return 11, 22
}
```

#### `-allow-named=false`

```go
func foo() (a int, b int) { // lint: functions with named return variables not allowed
    return 11, 22
}
```

#### `-allow-partially-named=false`

```go
func foo() (int, b int) { // lint: functions with partially named return variables not allowed
    return 11, 22
}
```

#### `-allow-normal-returns=false`

```go
func foo() (a int, b int) {
    return 11, 22 // lint: normal (non-naked) returns not allowed
                  // in functions with named return variables
}
```

#### `-allow-naked-returns=false`

```go
func foo() (a int, b int) {
    a = 11
    b = 22
    return // lint: naked returns not allowed
}
```

#### `-allow-mixing-returns=false`

```go
func foo() (a int, b int) {
    if cond() {
        return 11, 22
    } else {
        a = 11
        b = 22
        return // lint: mixing normal and naked returns not allowed
    }
}

func bar() (a int, b int) {
    if cond() {
        a = 11
        b = 22
        return
    } else {
        return 11, 22 // lint: mixing normal and naked returns not allowed
    }
}
```

## Authors

See [here](https://github.com/gavv/returnstyles/graphs/contributors).

## License

[MIT](LICENSE)
