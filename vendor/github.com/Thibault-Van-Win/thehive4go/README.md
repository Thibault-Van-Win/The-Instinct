# TheHive4Go

TheHive4Go aims to implement an The Hive API client written in Go. This project contains a private module which can be imported across different projects to improve code reusability.

## How to use

Some examples on how to use the client can be found below

### Initializing the client

The client can be initialized via:

```go
package main

import (
    "github.com/Thibault-Van-Win/thehive4go"
)

func main() {
    config := &thehive4go.Config{
        URL: "https://your-hive-instance",
        APIKey: "your-hive-token",
        SkipTLSVerification: true,
    }

    client := thehive4go.NewAPIClient(*config)
}
```

Next, make sure that `go get` can find our private repository by setting the go environment variable:

```sh
go env -w GOPRIVATE=git.nias.one
```

Finally fetch the module by running:

```sh
go mod tidy
```

### Fetching alerts

Once a client is instantiated, build-in services can be used to fetch different artifacts

```go
package main

import (
    "log"

    "github.com/Thibault-Van-Win/thehive4go"
)

func main() {
    config := &thehive4go.Config{
        URL: "https://your-hive-instance",
        APIKey: "your-hive-token",
        SkipTLSVerification: true,
    }

    client := thehive4go.NewAPIClient(*config)

    alerts, err := client.Alerts.List()
    if err != nil {
        log.Fatalf("Failed to fetch alerts: %v", err)
    }
}
```

### Using filters on listings

The `List` endpoints can be used to query The Hive for a collection of data. An options pattern is in pace for further filtering. Example:

```go
package main

import (
    "log"

    "github.com/Thibault-Van-Win/thehive4go"
    "github.com/Thibault-Van-Win/thehive4go/query"
)

func main() {
    config := &thehive4go.Config{
        URL: "https://your-hive-instance",
        APIKey: "your-hive-token",
        SkipTLSVerification: true,
    }

    client := thehive4go.NewAPIClient(*config)

    alerts, err := client.Alerts.List(
        query.WithCompany("Able bv"),
    )
    if err != nil {
        log.Fatalf("Failed to fetch alerts: %v", err)
    }
}
```

## Dependencies

- Go v1.22
