
# Introduce

Maportd is a lightweight port mapping services that enable listening on one port and mapping to one or multiple service addresses.

# Compile

```bash
go build -o bin/maportd apps/maport/main.go
```

# Run

### Arguments

```bash
Usage of maportd:
  -dest string
        The destination address.
  -log string
        Write log messages to this file. the default 'stdout'
  -log-caller
        Whether to log the caller.
  -log-format string
        The format of messages to log.
  -log-level string
        The level of messages to log.
  -port int
        The source address.
  -version
        Show the version.
```

### Example

```bash
bin/maportd -port=3307 -dest=192.168.0.2:3306,192.168.0.3:3306
```