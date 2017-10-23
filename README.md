# Help

## Usage:
```bash
xpath-go <PATH> <TARGET>
```
Query html from stdin via xpath expression and output in json.

## Arguments:
PATH   : expression to match.
TARGET : raw json string of `key`/`value` pairs. `value` should be relative path expression from leaf node matched above.

## Examples:
```bash
    cat test.html | xpath-go "//div[@class=\"seckill-timer\"]\" "{\"id\":\"./@id\"}"
```

