package main

import (
	"fmt"
	"log"
	"strings"
  "flag"
  "os"
  "bufio"
  "regexp"
  "time"
)

type fileArgs []string

func (f *fileArgs) String() string {
  return fmt.Sprint(*f)
}

func (f *fileArgs) Set(value string) error {
  for _, str := range strings.Split(value, ",") {
    *f = append(*f, str)
  }
  return nil
}

func replaceMonth(in []byte) []byte {
  month := string(in)
  switch month {
  case "Jan":
    return []byte("01")
  case "Feb":
    return []byte("02")
  case "Mar":
    return []byte("03")
  case "Apr":
    return []byte("04")
  case "May":
    return []byte("05")
  case "Jun":
    return []byte("06")
  case "Jul":
    return []byte("07")
  case "Aug":
    return []byte("08")
  case "Sep":
    return []byte("09")
  case "Oct":
    return []byte("10")
  case "Nov":
    return []byte("11")
  case "Dec":
  }
  return []byte("12")
}

const logPattern = `^(?P<ipaddress>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})` +
`.*\[(?P<timestamp>\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2} (?:\+|\-)\d{4})\].*` +
`emailAddress=(?P<email>[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+).*` +
`(?P<ua>".*"\s"-"$)`

const timePattern = `^(?P<day>\d{2})\/(?P<month>\d{2})\/(?P<year>\d{4}):`+
`(?P<hour>\d{2}):(?P<minutes>\d{2}):(?P<seconds>\d{2}) (?:\+|\-)\d{4}`

const monthPattern = `(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)`


func main() {
  var inFlags fileArgs
  var locationFlag string

  flag.Var(&inFlags, "in", "comma-separated list of log file names")
  flag.StringVar(&locationFlag, "locale", "",
    "specify target locale for cardinality calculation")
  flag.Parse()

  if (len(inFlags) == 0 || locationFlag == "") {
    flag.PrintDefaults()
    os.Exit(1)
  }

  location, err := time.LoadLocation(locationFlag)
  if err != nil {
    fmt.Println("Unknown location " + locationFlag)
    os.Exit(1)
  }

  re := regexp.MustCompile(logPattern)

  timeRe := regexp.MustCompile(monthPattern)
  timeConvert := regexp.MustCompile(timePattern)


  result := make(chan int)

  for _,fname := range inFlags {
    file, err := os.Open(fname)
    if err != nil {
      log.Fatal(err)
    }
    defer func() {
      if err = file.Close(); err != nil {
        log.Fatal(err)
      }
    } ()
    go func() {
      scanner := bufio.NewScanner(file)
      scanner.Split(bufio.ScanLines)
      for scanner.Scan() {
        line := scanner.Bytes()
        submatch := re.FindAllSubmatch(line, -1)
        if submatch != nil {
          // normalize month
          timeSlice := timeRe.ReplaceAllFunc(submatch[0][1:][1], replaceMonth)

          fmt.Printf("%q\n", timeSlice)

        	normalizedTime := []byte{}

          normalizedTime = timeConvert.Expand(normalizedTime,
            []byte(`${year}-${month}-${day}T${hour}:${minutes}:${seconds}Z`),
            timeSlice,
            timeConvert.FindAllSubmatchIndex(timeSlice, -1)[0])

          t, err := time.ParseInLocation(time.RFC3339,
            string(normalizedTime), location)
          if err != nil {
            log.Fatal(err)
            os.Exit(1)
          }
          fmt.Printf("%q\n", submatch[0][1:], t.In(location))
        }
      }
      if scanner.Err() != nil {
        log.Fatal(err)
      }
      result <- 1
    } ()
  }

  <-result
}
