package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
)

func main () {
  argsWithoutProg := os.Args[1:]
  fmt.Println("Arguments: ", argsWithoutProg)

  if len(argsWithoutProg) != 1 {
    fmt.Println("Nothing to do.")
    os.Exit(1)
  }

  req, reqErr := http.NewRequest("GET", argsWithoutProg[0], nil)

  if reqErr != nil {
    fmt.Println("Error while creating request: ", reqErr)
    os.Exit(10)
  }

  client := &http.Client {}
  resp, respErr := client.Do(req)

  if respErr != nil {
    fmt.Println("There was an error.")
    fmt.Println(respErr)
    os.Exit(20)
  }

  defer resp.Body.Close()

  // fmt.Println(req)

  for k, v := range req.Header {
    fmt.Printf(">> '%s' = '%s'\n", k, v)
  }
  fmt.Println("---")

  // fmt.Println(resp)
  for k, v := range resp.Header {
    fmt.Printf("<< '%s' = '%s'\n", k, v)
  }

  fmt.Println("---")

  bodyBytes, readErr := ioutil.ReadAll(resp.Body)

  if readErr != nil {
    fmt.Println("Error while reading body.", readErr)
    os.Exit(30)
  }

  fmt.Println(string(bodyBytes))

}
