[![Build Status](https://travis-ci.org/fernandezvara/corbel-go.svg?branch=master)](https://travis-ci.org/fernandezvara/corbel-go)
[![GoDoc](https://godoc.org/github.com/fernandezvara/go-silkroad?status.png)](https://godoc.org/github.com/fernandezvara/corbel-go)
[![Coverage Status](https://coveralls.io/repos/fernandezvara/corbel-go/badge.svg?branch=master)](https://coveralls.io/r/fernandezvara/corbel-go?branch=master)

# **Corbel-Go**
-----

corbel-go is an API library for work with [corbel](http://opensource.bq.com/). It currently supports:

  > - Creation of a new Client
  > - Token workflow
  > - Basic Authentication (username/password)
  > - Resources (get/create/update/delete/search)

Note: library in active development; requires >= Go 1.3


## **Usage/Sample Code**

### **Creating a client**

```Go
// NewClient(http.Client, clientId, clientName, clientSecret, clientScopes,
//           clientDomain, JWTSigningMethod, tokenExpirationTime)
client, _ = NewClient(nil, "someID", "", "someSecret", "", "", "HS256", 3000)
```

### **Authorization**

**Getting token for client app**

Using the client for application operations.

```Go
err = client.IAM.OauthToken()
```

**Getting token for user using basic auth**

Using the client for user operations in the app.

```Go
err = client.IAM.OauthTokenBasicAuth("username", "password")
```
### **Resources**

**Adding resource**

Adds a resource of a defined type. Definitions are JSON parseable structs.

*Important Note:* Avoid using omitempty in the JSON definition if you think you could have a value that could turn _false, 0, empty strings or nil_. In those cases _json.Marshal_ won't export the data. So value won't be updated in the backend.
*Important Note 2:* Is recommended to define the ID on the structs to be able to update them correctly without workarounds.

```Go
type ResourceForTest struct {
  ID   string  `json:"id,omitempty"`
  Key1 string  `json:"key1"`
  Key2 int     `json:"key2"`
  Key3 float64 `json:"key3"`
  Key4 bool    `json:"key4"`
}

test1 := ResourceForTest{
  Key1: "test string",
  Key2: 123456,
  Key3: 1.123456,
  Key4: true,
}

err = client.Resources.AddToCollection("test:GoTestResource", &test1)
```

#### **Search for Resources**

Search allow to browse for the required resources using a simple interface.
All Search conditions are shared in all modules of _corbel_, so it's the same for users, for example.

**Per Page**

```Go
var arrResourceForTest []ResourceForTest

search = client.Resources.SearchCollection("test:GoTestResource")

err = search.Page(0, &arrResourceForTest)
```

**Conditions**

Searching resources by specifying conditions.

```Go
search = client.Resources.SearchCollection("test:GoTestResource")

// all items where firstName == "testName"
search.Query.Eq["firstName"] = "testName"

// sort by firstName
search.Sort.Asc = []string{"firstName"}

// list 20 resources por search page
search.PerPage = 20

err = search.Page(0, &arrResourceForTest)
```

**All allowed search conditions**

```Go
	Eq   map[string]string     // Equal to
	Gt   map[string]int        // Greater than
	Gte  map[string]int        // Greater than or equal
	Lt   map[string]int        // Less than
	Lte  map[string]int        // Less than or equal
	Ne   map[string]string     // Not Equal
	In   map[string][]string   // One of this array
	All  map[string][]string   // All of this array
	Like map[string]string     // Like
```

**Sort conditions**

```Go
  Asc  []string   // Ascendent
  Desc []string   // Descendent
```

**Aggregations**

```Go
func (s *Search) Count(field string) (int, error)
func (s *Search) CountAll() (int, error)
func (s *Search) Average(field string) (float64, error)
func (s *Search) Sum(field string) (float64, error)
```

#### **Get resource**

```Go
test2 := ResourceForTest{}

err = client.Resources.GetFromCollection("test:GoTestResource",
                                         "1234567890abcdef", &test2)
```

#### **Updating resource**

```Go
test2.Key1 = "new string"
err = client.Resources.UpdateInCollection("test:GoTestResource",
                                         "1234567890abcdef", &test2)
```

#### **Delete resource**

```Go
err = client.Resources.DeleteFromCollection("test:GoTestResource",
                                            "1234567890abcdef")
```


----

#### **Contributing**

 - Fork it
 - Create your feature branch (git checkout -b my-new-feature)
 - Commit your changes (git commit -am 'Add some feature')
 - Push to the branch (git push origin my-new-feature)
 - Create new Pull Request
 - If applicable, update the README.md
