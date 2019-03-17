package main

import (
	"fmt"
	"log"
	"strings"
  "flag"
  "os"
  "bufio"
  "regexp"
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

const logPattern = `^(?P<ipaddress>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})` +
`.*\[(?P<timestamp>\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2} (?:\+|\-)\d{4})\].*` +
`emailAddress=(?P<email>[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+).*` +
`(?P<ua>".*"\s"-"$)`

//type explore map[int][]string

//type record struct {
//  ip []byte
//  dateTime []byte
//  email []byte
//  userAgent []byte
//}

func main() {
  var inFlags fileArgs

  flag.Var(&inFlags, "in", "comma-separated list of log file names")
  flag.Parse()

  if len(inFlags) == 0 {
    flag.PrintDefaults()
    os.Exit(1)
  }

  re := regexp.MustCompile(logPattern)
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
          //recordVal := record{ submatch[0][1],
          //  submatch[0][2], submatch[0][3], submatch[0][3] }
          fmt.Printf("%q\n", submatch[0][1:])
        }
      }
      if scanner.Err() != nil {
        log.Fatal(err)
      }
      result <- 1

    // go func() {
    //   example := make(explore)
    //   r := csv.NewReader(file)
    //   r.Comma = ' ';
    //   r.FieldsPerRecord = -1;
    //   for {
    // 		record, err := r.Read()
    // 		if err == io.EOF {
    //       result <- &example
    // 			break
    // 		}
    // 		if err != nil {
    // 			log.Fatal(err)
    // 		}
    //     example[len(record)] = record
    // 	}
    } ()
  }

  <-result
//  json, err := json.MarshalIndent(*<-result, "", "  ")
//  if err != nil {
//		log.Fatal(err)
//	}
//  fmt.Println(string(json))
}
