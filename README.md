# yamlkeysort

```
Usage of yamlkeysort:
  -f string
        read and write from the same file
  -h    help text
  -i string
        input filename (default stdin)
  -o string
        output filename (default stdout)

Example:

  yamlkeysort -i input.yml -o output.yml first_key second_key third_key

```


## Example

### Before

```yaml
---
c: 3
dont_care_first: ""
a: 1
dont_care_second: ""
embeded:
  b: ""
  a: ""
slice:
- b: ""
  a: ""
b: 2
```

### Execute

`yamlkeysort -f test_data/example.yml a b c`

### After
```yaml
---
a: 1
b: 2
c: 3
dont_care_first: ""
dont_care_second: ""
embeded:
  a: ""
  b: ""
slice:
- a: ""
  b: ""

```
