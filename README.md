# NAME

**pub** - Publish resources at specified paths

# SYNOPSIS

**pub** [res1 path1]... [resN [pathN]]

# DESCRIPTION

The `pub` command accepts an optional sequence of resource-path pairs. Each pair consists of a resource (`res1`) and a path (`path1`). These pairs can be repeated any number of times, including not at all. For each pair, the path is mandatory except for the last resource (`resN`). When the path is omitted, it defaults to "/".

# EXAMPLES

Publish a directory at the default path:

```
$ pub .
```

Publish a specified port:

```
$ pub :8080
```

Publish a URL at the default path:

```
$ pub https://example.com
```

Publish at a port with a specified path, followed by a URL at the default path:

```
$ pub \
  :9944 /rpc/ws \
  https://polkadot.js.org/apps/
```

# SEE ALSO

- https://pkg.go.dev/github.com/btwiuse/pub
