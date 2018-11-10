# Not affected files

This scenario simulates multiple files in packages but only one of each is
affected with cycle. The output should omit not affected files and imports.

```text
    +-----+     +-----+     +-----+
    |     | --> |     |     |     |
    | BAR |     | BAZ | --> | FOO |
    |     | <-- |     |     |     |
    +-----+     +-----+     +-----+
```
