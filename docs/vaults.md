## Vaults

### AGE files

[age](https://github.com/FiloSottile/age) as an encryption tool can be used as a source for
secrets, e.g. for development purposes. `go-secretshelper` accepts armored, age-encrypted files 
and a matching identity, also from a file. Example:

```yaml
vaults:
- name: kv
  type: age-file
  spec:
    path: ./fixtures/test-agefile
    identity: ./fixtures/test-identity
```

This will retrieve the secrets from the age file under `spec/path` and decrypt them with the identity (from `spec/identity`).
The file has to be json-encoded, mapping the names of secrets to the secrets, e.g. to produce:

```bash
$ age-keygen -o ./fixtures/test-identity
$ echo '{ "test": "s3cr3t" }' | age -e -r <identity-from-previous-step> -a
```

### Azure Key Vault

Secrets can be accessed from an [Azure Key Vault](https://azure.microsoft.com/de-de/services/key-vault/).
Within the `vault` section of a configuration file, add the following to access a vault under
a given URL:

```yaml
  - name: kv
    type: azure-key-vault
    spec:
      url: https://my-sample-vault.vault.azure.net/
```

This will access secrets in `my-sample-vault`. If you want to access secrets in a different
type of vault (e.g. HSM-backed) you can specify the URL accordingly.

In case of default vault service, the url can be omitted. The following snippet does the same.
However this requires using the name of the vault as it is in the `sinks` and `transformations` sections

```yaml
  - name: my-sample-vault
    type: azure-key-vault
```

### AWS Secrets Manager

Secrets can be accessed from an [AWS Secret Manager](https://aws.amazon.com/secrets-manager/?nc1=h_ls).
Specify the vault as follows:

```yaml
vaults:
  - name: mysecrets
    type: aws-secretsmanager
    spec:
      region: us-east-2
```

This will access secrets from the given region.

It uses the default [credentials mechanism of the AWS Go SDK](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html).
I.e. the credentials are read from the environment variables `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`.
The `region` entry above is optional if region info is given by the environment variable `AWS_REGION`.

### GCP Secret Manager

Secrets can be accessed from [GCP's Secret Manager](https://cloud.google.com/secret-manager/). It uses the
default credentials mechanism of the Google Cloud SDK, e.g. from environment variables `PROJECT_ID` and
`ACCESS_TOKEN`. A vault is declared with the type `gcp-secretmanager` as follows:

```yaml
vaults:
  - name: mysecrets
    type: gcp-secretmanager
    spec:
      projectID: fancy-projectid-3746342
```