## Transformations

### Template

A Template Transformation is able to render a text-based template with secrets previously
pulled from vaults. To use a template, add the following section to a configuration:

```yaml
transformations:
  - type: template
    in:
      - inputVar1
      - inputVar2
    out: outputVar
    spec:
      template: |
        other={{ .inputVar1 }}
        value={{ .inputVar2 }}
```

The above snippet renders the following output and stores it into `outputVar` for further processing,
given that `inputVar1` equals `some` and `inputVar2` equals `secret`.

```
other=some
value=secret
```

To add a Content Type, add to the spec part:

```yaml
    spec:
      contentType: text/plain
      template: |
        other={{ .inputVar1 }}
        value={{ .inputVar2 }}
```

### Age encryption

The Age encrypt transformation takes one or more secrets as input and encrypts them
using age, for the specified recipients. Output is rendered as armored age and put
into the output variable:

```yaml
transformations:
  - type: age
    in:
      - test
    out: test-enc
    spec:
      recipient: ${age_recipient}
```

The above part will encrypt the secret `test` and store it in `test-enc`. The recipient
used for age-encryption is taken from the environment variable `age_recipient`.

### JQ

The JQ transformation takes one or more secrets as input and applies a jq filter to them. This
allows for filtering of values in arbitrary json structs, e.g. for a secret named `s1` of content type
application/json and content like

```json
{
  "endpoint": "foo",
  "apikey": "bar"
}
```

this transformation pulls out the `apikey` and stores it in the output variable `s1-apikey`:

```yaml
transformations:
  - type: jq
    in:
      - s1
    out: s1-apikey
    spec:
      q: ".apikey"
      raw: true
```

Setting `raw` to `false` lets you use the actual response as a json structure for later
processing, e.g. for a secret `s2`

```json
{
  "db1": {
    "endpoint": "foo",
    "apikey": "bar"
  },
  "db2": {
    "endpoint": "baz",
    "apikey": "qux"
  }
}
```

this transformation 

```yaml
transformations:
  - type: jq
    in:
      - s2
    out: s2-db2
    spec:
      q: ".db2"
      raw: false
```

yields the contents of `db2` as a json structure into new secret `s2-db2`.

```json
{
    "endpoint": "baz",
    "apikey": "qux"
}
```