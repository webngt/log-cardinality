# log-cardinality
Counts cardinality of daily active users grouped by application (user agent) type

Source log file records must strictly match the folowing regexp

```go
const logPattern = `^(?P<ipaddress>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})` +
`.*\[(?P<timestamp>\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2} (?:\+|\-)\d{4})\].*` +
`emailAddress=(?P<email>[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+).*` +
`(?P<ua>".*"\s"-"$)`
```

Usage
```bash
log-cardinality -in "`find $LOG_PATH_PATTERN/access.log -printf '%p '`" -locale Europe/Moscow
```

# Sample Output

```json
{
  "2019-03-14": {
    "Electron": 719,
    "Web": 148,
    "grpc-java": 182,
    "grpc-objc": 41
  },
  "2019-03-15": {
    "Electron": 742,
    "Web": 143,
    "grpc-java": 187,
    "grpc-objc": 47
  }
}
```


# Depends on
  * [Highway](https://github.com/google/highwayhash/) for hash calculation
  * [HyperLogLog++ for Go](https://github.com/lytics/hll) to calculate cardinality
