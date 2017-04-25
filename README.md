# go-adzerk
go-adzerk is a Go client library for accessing [Adzerk's Native Ads API](http://dev.adzerk.com/reference#native-ads-api-overview).

## Usage
```go
import "github.com/Wattpad/go-adzerk/adzerk"
```

Constructing a new Adzerk client:

```go
client := adzerk.NewClient(nil)
```

Constructing the HTTP request can look something like the following:

```go
placements := []adzerk.Placement{
	{
		DivName:   "div1",
		AdTypes:   []int{4, 5},
		SiteID:    123,
		ZoneIDs:   []int{456},
		NetworkID: 789,
	},
}

httpRequest, err := client.NewRequest(adzerk.RequestData{
	IP:               "10.123.123.123",
	UserID:           "ad39231daeb043f2a9610414f08394b5",
	BlockedCreatives: []int{123, 456},
	Keywords:         []string{"foo", "bar", "baz"},
	Placements:       placements,
})

if err != nil {
	// Handle non-nil error
}
```

Performing the request:

```go
var v map[string]interface{}
httpResponse, err = client.Do(ctx, httpRequest, &v)
```

Modifying the URL which the client will hit – for example, when creating mock service tests:

```go
client.URL = "http://mockserveraddress.com"
```

## License

This library is distributed under the MIT-style license found in the [LICENSE](./LICENSE) file.
