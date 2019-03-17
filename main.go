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

//type explore map[int][]string

func main() {
  var inFlags fileArgs
  var regexpFlag string

  flag.Var(&inFlags, "in", "comma-separated list of log file names")
  flag.StringVar(&regexpFlag, "regexp", "",  "regexp")
  flag.Parse()

  if len(inFlags) == 0 {
    flag.PrintDefaults()
    os.Exit(1)
  }

  re := regexp.MustCompile(regexpFlag)
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
        fmt.Printf("%q\n", re.FindAllSubmatch(line, -1))
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
