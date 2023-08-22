# go-printify

A fork from <https://github.com/omrikiei/go-printify>. I included work from <https://github.com/grantok/go-printify> to add the uploads.

This project should be considered unstable.

## How to use

A basic main.go file that returns a list of your shops is below. Make sure to set your [access token](https://developers.printify.com/#authentication) to a env variable named `PRINTIFY_ACCESSTOKEN`

```
package main

import (
    "os"
    "github.com/brandonmcclure/go-printify"
)

func main() {
	accesstoken := os.Getenv("PRINTIFY_ACCESSTOKEN")
	client := go_printify.NewClient(accesstoken)
	client.ListShops()
}
```