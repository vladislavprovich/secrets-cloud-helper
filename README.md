# secrets-сloud-helper

## Usage

`go-secretshelper` expects a yaml-based configuration file, which it processes. The configuration contains four major elements:

* **Vaults** specify, where secrets are stored. Examples are Azure Key Vault or AWS Secrets Manager
* **Secrets** define, what data is read from which vault.
* **Transformation** describe,how secrets are modified, e.g. to decode base64 or render a template
* **Sinks** specify where and how secrets are written. At present, only files are supported as sinks.

To run a configuration, use: 

```bash
$ go-secretshelper run -c <config file>
```

Sample configuration file:
```yaml
vaults:
  - name: myvault
    type: aws-secretsmanager
    spec:
      region: us-east-2

secrets:
  - type: secret
    vault: myvault
    name: sample

transformations:
  - type: template
    in:
      - sample
    out: sample-ini
    spec:
      template: |
        thesecret={{ .sample }}

sinks:
  - type: file
    var: sample-ini
    spec:
      path: ./sample.ini
      mode: 400
```

The above configuration defines a secret named `sample`, which is read from the AWS Secrets Manager instance in `us-east-2`. The secret is then transformed by the 
template and stored in a new secret named `sample-ini`. The new secret is written to a file named `./sample.ini` with file mode 400. Such a configuration may define
multiple vaults, secrets, multiple transformations and sinks.

See [docs/](docs/README.md) for more details. A configuration file may contain environment variables, which are expanded before processing by using the `-e` switch, e.g.:

```yaml
secrets:
  - type: secret
    vault: ${VAULT_NAME}
    name: sample
```

This will expand the vault name of the environment variable `VAULT_NAME` and continue. This makes it possible to use the same configuration 
file for multiple environments.

## Building

The Makefile's `build` target builds an executable in `dist/`.

```bash
$ make build 
```

To build exectuables for several platforms, the `release` target uses [goreleaser](https://goreleaser.com/):

```bash
$ make release
```


## Testing

### Unit tests

```bash
$ go test -v ./...
```

### CLI tests

CLI tests are shell-based and written using bats. The executable is expected to be present in `dist/`. so `make build` 
is necessary before. To run the tests:

```bash
$ cd tests
$ bats .
```

## Contributing

Pull requests are welcome!