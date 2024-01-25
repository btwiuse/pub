# NAME

**pub** - Publish anything to the open web

# SYNOPSIS

**pub** [res1 path1]... [resN [pathN]]

# DESCRIPTION

The `pub` command accepts an optional sequence of resource-path pairs, where resource could be any file, dir, port or url, and path specifies the mountpoint of the resource.

For each pair, the path is mandatory except for the last resource (`resN`).  When the path is omitted, it defaults to "/".

In the special case of N=1, `pub README.md /` could be simplified to `pub README.md`.

# EXAMPLES

```
# Publish a directory at the default path:
$ pub .

# Publish a specified port:
$ pub :8080

# Publish a URL at the default path:
$ pub https://example.com

# Publish at a port with a specified path, followed by a URL at the default path:
$ pub \
  :9944 /rpc/ws \
  https://polkadot.js.org/apps/
```

# SEE ALSO

- https://pkg.go.dev/github.com/btwiuse/pub
