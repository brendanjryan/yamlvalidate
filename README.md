`yamlvalidate`
---

A command line tool for validating `yaml` files against a [JSON Schema](https://json-schema.org/).

This utility was originally written for verifying Kubernetes YAML files, though it is generic enough to be used for any valid YAML file and JSON schema. 

## Usage

```console
./yamlvalidate -s example/schema.json example/invalid.yaml 
```

## Examples

### Valid File
```console
$ .yamlvalidate -s example/schema.json example/valid.yaml 
2019/09/17 15:37:28 All files validated successfully!
```

### Invalid File
```console
$ .yamlvalidate -s example/schema.json example/invalid.yaml 
2019/09/17 15:38:42 error validating file example/invalid.yaml - tags.2 Invalid type. Expected: string, given: integer
2019/09/17 15:38:42 error validating file example/invalid.yaml - dimensions.height Invalid type. Expected: integer, given: string
2019/09/17 15:38:42 error validating file example/invalid.yaml - price Invalid type. Expected: number, given: string
2019/09/17 15:38:42 Validation failed
2019/09/17 15:38:42 Please fix the following files: 
2019/09/17 15:38:42 example/invalid.yaml
```

### Mixed Files
```console
$ .yamlvalidate -s example/schema.json example/valid.yaml example/invalid.yaml
2019/09/17 15:38:14 error validating file example/invalid.yaml - price Invalid type. Expected: number, given: string
2019/09/17 15:38:14 error validating file example/invalid.yaml - tags.2 Invalid type. Expected: string, given: integer
2019/09/17 15:38:14 error validating file example/invalid.yaml - dimensions.height Invalid type. Expected: integer, given: string
2019/09/17 15:38:14 Validation failed
2019/09/17 15:38:14 Please fix the following files: 
2019/09/17 15:38:14 example/invalid.yaml
```

## Validating Kubernetes Schemas

By [extending the JSON Schema](https://json-schema.org/understanding-json-schema/structuring.html) of the Kubernetes API, we are able to write our own typesafe validation pipelines on top of the already powerful OpenAPI spec for all Kubernetes objects. Business domain specific rules can enforce constraints on Kubernetes files, such as resource limits or whitelisted registry domains without needing to migrate all of your kubernetes files from `yaml` to a protocol-buffer based framework like [isopod](https://github.com/cruise-automation/isopod) or [skycfg](https://github.com/stripe/skycfg).

For example: the following rule can be used to check that all [`Quantity`](https://github.com/kubernetes/apimachinery/blob/master/pkg/api/resource/quantity.go) declarations use no more than a single (numeric) virtual core. 

```json
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://brendanjryan.com/limited_quantity.json",
  "allOf": [
    { 
      "$ref": "https://kubernetesjsonschema.dev/master/_definitions.json#/definitions/io.k8s.apimachinery.pkg.api.resource.Quantity"
    },
    { 
      "type": "number",
      "minimum": 0,
      "exclusiveMaximum": 1
    }
  ]
}
```

