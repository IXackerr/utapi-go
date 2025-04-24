# utapi-go (current as of 24.04.2025)

A thin wrapper for the UploadThing API.

## Why?

There was no working UploadThing API client for Go (only outdated and broken ones).

## Setup

You will need a `.env` file with your UploadThing API secret key:

```env
UPLOADTHING_SECRET=sk_*************************
```

## Usage

After adding your import statement as below, run `go mod tidy` to install dependencies.

```go
package main

import (
    "github.com/jesses-code-adventures/utapi-go"
    "os"
    "fmt"
)

func main() {
    // Create API handler
    utApi, err := utapi.NewUtApi()
    if err != nil {
        fmt.Println("Error creating UploadThing API handler:", err)
        os.Exit(1)
    }

    // Example: deleting a file
    // This is the key returned by UploadThing when you upload a file
    keys := []string{"fc8d296b-20f6-4173-bfa5-5d6c32fc9f6b-geat9r.csv"}
    resp, err := utApi.DeleteFiles(keys)
    if err != nil {
        fmt.Println("Error deleting files:", err)
    } else {
        fmt.Println("Successfully deleted file(s):", resp.Success)
    }
}
```

## More examples

See [examples.md](examples.md) for additional usage scenarios.