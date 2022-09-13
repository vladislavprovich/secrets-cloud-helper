## Sinks

### File sink

`file` is the only available type of sink at present. It emits a single variable to a file, optionally 
setting file mode and user. A configuration can have multiple sinks.

Example:

```yaml
sinks:
  - type: file
    var: inputVar1
    spec:
      path: /mnt/secret/sample.dat
      mode: 400
```

This will write the content of `inputVar1` to a file `/mnt/secret/sample.dat` with file mode 400.