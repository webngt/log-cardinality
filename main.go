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
  highway "github.com/dgryski/go-highway"
  "github.com/lytics/hll"
  "strconv"
  "encoding/json"
  "github.com/webngt/timeutils"
)

type cardinality map[string]*hll.Hll

const logPattern = `^(?P<ipaddress>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})` +
`.*\[(?P<timestamp>\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2} (?:\+|\-)\d{4})\].*` +
`emailAddress=(?P<email>[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+).*` +
`(?P<ua>".*"\s"-"$)`

const uaPattern = `(grpc-java|grpc-objc|Electron)`


func main() {
  var monthDecode = [...]string{"01", "02", "03", "04", "05", "06", "07", "08", "09", "10", "11", "12"}
  var in string
  var locationFlag string
  key := highway.Lanes{0x0706050403020100, 0x0F0E0D0C0B0A0908, 0x1716151413121110, 0x1F1E1D1C1B1A1918}

  flag.StringVar(&in, "in", "", "space separated list of log file names")
  flag.StringVar(&locationFlag, "locale", "",
    "specify target locale for cardinality calculation")
  flag.Parse()

  inFlags := strings.Fields(in)

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

  appPattern := regexp.MustCompile(uaPattern)

  result := make(chan map[string]cardinality)
  dayBucket := make(chan map[string]cardinality)


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
      dayBuckets := make(map[string]cardinality)
      currYear, currMonth, currDay := 0,time.January,0
      dateKey := ""
      for scanner.Scan() {
        line := scanner.Bytes()
        submatch := re.FindAllSubmatch(line, -1)
        if submatch != nil {
          // normalize month
          date, err := timeutils.ParseDateString(string(submatch[0][1:][1]))

          if err != nil {
            log.Fatal(err)
            os.Exit(1)
          }
          // swap date to target location
          date = date.In(location)

          // create date key
          year, month, day := date.Date()
          if (currYear != year || currMonth != month || currDay != day) {
            dateKey = "" + strconv.Itoa(year) + "-" +
            monthDecode[month - 1] + "-" +strconv.Itoa(day)
            //fmt.Println(dateKey)
            currYear, currMonth, currDay = year, month, day
          }
          dayCardinality := dayBuckets[dateKey]
          if dayCardinality == nil {
            dayCardinality = make(cardinality)
            dayBuckets[dateKey] = dayCardinality
          }

          // app type
          appTypeSlice := appPattern.FindAllSubmatch(submatch[0][1:][3], -1)
          appType := "Web"
          if appTypeSlice != nil {
              appType = string(appTypeSlice[0][1:][0])
          }

          counter := dayCardinality[appType]

          if counter == nil {
            counter = hll.NewHll(14, 25)
            dayCardinality[appType] = counter
          }

          // user hash
          user := highway.Hash(key, submatch[0][1:][2])

          counter.Add(user)

        }
      }
      if scanner.Err() != nil {
        log.Fatal(err)
      }
      dayBucket <- dayBuckets
    } ()
  }

  // merge hll++
  go func() {
    merge := make(map[string]cardinality)
    for range inFlags {
      for date, apps := range <-dayBucket {
        if merge[date] != nil {
          for app, hll := range apps {
            if merge[date][app] != nil {
              merge[date][app].Combine(hll)
            } else {
              merge[date][app] = hll
            }
          }
        } else {
          merge[date] = apps
        }
      }
    }
    result <- merge
  } ()




  out := make(map[string]map[string]uint64)

  for date, apps := range <-result {
    out[date] = make(map[string]uint64)
    for app, hll := range apps {
      out[date][app] = hll.Cardinality()
    }
  }


  json, err := json.MarshalIndent(out, "", "  ")
  if err != nil {
    log.Fatal(err)
  }

  fmt.Println(string(json))
}
