package main

import (
  "encoding/csv"
	"fmt"
	"io"
	"log"
	"strings"
  "flag"
  "os"
  "encoding/json"
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

type explore map[int][]string

func main() {
  var inFlags fileArgs
  flag.Var(&inFlags, "in", "comma-separated list of log file names")
  flag.Parse()

  if len(inFlags) == 0 {
    flag.PrintDefaults()
    os.Exit(1)
  }

  result := make(chan *explore)

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
      example := make(explore)
      r := csv.NewReader(file)
      r.Comma = ' ';
      r.FieldsPerRecord = -1;
      for {
    		record, err := r.Read()
    		if err == io.EOF {
          result <- &example
    			break
    		}
    		if err != nil {
    			log.Fatal(err)
    		}
        example[len(record)] = record
    	}
    } ()
  }
  json, err := json.MarshalIndent(*<-result, "", "  ")
  if err != nil {
		log.Fatal(err)
	}
  fmt.Println(string(json))
}
