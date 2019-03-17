# log-cardinality
Count cardinality of given combination of fields from provided log file

Log file records must strictly match the folowing regexp

```go
const logPattern = `^(?P<ipaddress>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})` +
`.*\[(?P<timestamp>\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2} (?:\+|\-)\d{4})\].*` +
`emailAddress=(?P<email>[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+).*` +
`(?P<ua>".*"\s"-"$)`
```

Will likely depend on
  * [Siphash](https://github.com/dchest/siphash)
  * [Highway](https://github.com/google/highwayhash/)
  * [HyperLogLog++ for Go](https://github.com/lytics/hll)
  * [HyperLogLog axiomhq](https://github.com/axiomhq/hyperloglog)
