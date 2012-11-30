# mux

An experimental line muxer

## Build

```bash
$ go get
```

## Usage

### Basic

```bash
# in one window
$ mux -r
```

```bash
# in a second
$ mux -r
```

```bash
# in a third
$ echo hi there | mux -w
```

Notice `hi there` appear in the first and second windows.

### Custom Channel

The channel defaults to 'mux'.  You can change this with `-c`.

```bash
$ mux -r -c foobar
```