# mux

experimental line muxer.

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