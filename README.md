[![Build Status](https://travis-ci.org/fernandezvara/corbel-go.svg?branch=master)](https://travis-ci.org/fernandezvara/corbel-go)
[![GoDoc](https://godoc.org/github.com/fernandezvara/go-silkroad?status.png)](https://godoc.org/github.com/fernandezvara/corbel-go)
[![Coverage Status](https://coveralls.io/repos/fernandezvara/corbel-go/badge.svg?branch=master)](https://coveralls.io/r/fernandezvara/corbel-go?branch=master)

# **Corbel-Go**

corbel-go is an API library for work with [corbel](http://opensource.bq.com/). It currently supports:

  > - Creation of a new Client
  > - Token workflow
  > - Basic Authentication (username/password)
  > - Resources (get/create/update/delete/search)

*Note:* Library in active development; requires >= Go 1.3

-----

## **Usage/Sample Code**

### **Creating a client**

The client is the key to use the platform. All information here is custom for each client since its generated for every application that needs to use it.

Instancing the client allow to use it everywhere. Since Authorization can be for the application itself or its users on dynamic applications will be necessary to have several clients with several authorizations (one for each user) because you will have specific scopes based on your own permissions.

See the Authorization part to get all the possible variations.

```Go
// NewClient(http.Client, clientId, clientName, clientSecret, clientScopes,
//           clientDomain, JWTSigningMethod, tokenExpirationTime)
client, _ = NewClient(nil, "someID", "", "someSecret", "", "", "HS256", 3000)
```

### **Authorization**

**Getting token for client app**

When you will use the client for application purposes you must ask for a OauthToken to get the specific permissions, that normally are wide open than users.

```Go
err = client.IAM.OauthToken()
```

**Getting token for user using basic auth**

If the client will be used for operations as user you need validate and get the proper user token. User token will have all the permissions applied to that user that are custom for her.

```Go
err = client.IAM.OauthTokenBasicAuth("username", "password")
```

### **User Administration**

All actions over users on the domain can be done if the application/user have the required permissions. All user interactions are done using the IAM (Identity and Authorization Management) endpoint.

**IAM User definition**

All operations for an user must be done using the IAMUser struct.

```Go
type IAMUser struct {
	ID          string                 `json:"id,omitempty"`
	Domain      string                 `json:"domain,omitempty"`
	Username    string                 `json:"username,omitempty"`
	Email       string                 `json:"email,omitempty"`
	FirstName   string                 `json:"firstName,omitempty"`
	LastName    string                 `json:"lastName,omitempty"`
	ProfileURL  string                 `json:"profileUrl,omitempty"`
	PhoneNumber string                 `json:"phoneNumber,omitempty"`
	Scopes      []string               `json:"scopes,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Country     string                 `json:"country,omitempty"`
	CreatedDate int                    `json:"createdDate,omitempty"`
	CreatedBy   string                 `json:"createdBy,omitempty"`
}
```

*NOTES:*
- User properties is a map that allow to add arbitrary information of that user. All JSON serializable types are allowed. _Beware of those properties that can be empty, 0, false or nil, since it won't be exported if omitempty are used_.
- Scopes can be an empty array if are defined default scopes for users on the domain definition.
- CreatedDate and CreatedBy are managed by the platform itself, so any change there will be ignored.

#### **User Creation**

```Go
anUserProperties := make(map[string]interface{})
anUserProperties["string"] = "test string"
anUserProperties["integer"] = 123456
anUserProperties["float"] = 1.23
anUserProperties["date"] = now

anUser := IAMUser{
  Domain:      "corbel-qa",
  Username:    "corbel-go",
  Email:       "corbel-go@corbel.org",
  FirstName:   "Corbel",
  LastName:    "Go",
  ProfileURL:  "http://corbel.org/corbel-go",
  PhoneNumber: "555-555-555",
  Scopes:      []string{},
  Properties:  anUserProperties,
  Country:     "Somewhere",
}

err = client.IAM.Add(&anUser)
```

#### **User Get by ID**

```Go
anUser2 := IAMUser{}
err = client.IAM.Get("sampleId", &anUser2)
```

#### **User Get Current User**
```Go
currentUser := IAMUser{}
err = client.IAM.GetMe(&currentUser)
```


#### **User Update**

```Go
anUser.Country = "Internet"
err = client.IAM.Update("sampleId", &anUser)
```

#### **User Deletion**

```Go
err = client.IAM.Delete("sampleId")
```

#### **User Search**

```Go
search := client.IAM.Search()
search.Query.Eq["username"] = "corbel-go"

var arrUsers []IAMUser

err = search.Page(0, &arrUsers)
```

*NOTE*: Searching uses the same interface defined in detail on the Resources documentation part.


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
func (s *Search) Count(field string) (int, error) {}
func (s *Search) CountAll() (int, error) {}
func (s *Search) Average(field string) (float64, error) {}
func (s *Search) Sum(field string) (float64, error) {}
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


### **Relations between Resources**

Resources can have related resources using collections. As sample think in a Music Group resource that have several Album resources.

### **Add Related Resource to Collection**

Adding a resource to a Collection automatically creates the Colection, so you don't need to be worried on Collection creation.
All relations allow to add custom metadata. To query that metadata is posible to create custom metadata structs to use it later.

```Go
// Sample without metadata
err = client.Resources.AddRelation("test:MusicGroup", "12345",
                                   "test:Albums",
                                   "test:Album", "23456",
                                   nil)

// Sample with custom metadata
type groupAlbumRelation struct {
  RecordLabel string `json:"recordLabel"`
}
metadata := groupAlbumRelation{
  RecordLabel: "Sample Record Label",
}
err = client.Resources.AddRelation("test:MusicGroup", "12345",
                                   "test:Albums",
                                   "test:Album", "23456",
                                   metadata)
```

### **Search Related Resources Information**

To get relations you can use the RelationData struct if you don't added metadata or extend RelationData with your own data.

```Go
// RelationData definition.
type RelationData struct {
	Order float64                  `json:"_order,omitempty"`
	ID    string                   `json:"id,omitempty"`
	Links []map[string]interface{} `json:"links, omitempty"`
}
```

Sample with standard metadata.
```Go
var arrRelationData []corbel.RelationData

search = client.Resources.SearchRelation("test:Group", "12345",
                                         "test:Albums")
err = search.Page(0, &arrRelationData)
```

Sample with _custom_ metadata.
```Go
type customRelationData struct {
	Order       float64                  `json:"_order,omitempty"`
	ID          string                   `json:"id,omitempty"`
	Links       []map[string]interface{} `json:"links, omitempty"`
  RecordLabel string                   `json:"recordLabel"`
}
var arrRelationData []customRelationData

search = client.Resources.SearchRelation("test:Group", "12345",
                                         "test:Albums")
err = search.Page(0, &arrRelationData)
```

### **Get Resource from Response**

Searching for resources information does not return the target object itself, it returns a pointer to to plus the custom metadata (if added).

```Go
type Album struct {
  Title           string `json:"title"`
  PublicationYear int    `json:"publicationYear"`
}

var anAlbum Album

err = client.Resources.GetFromRelationDefinition(arrRelationData[0].ID, &anAlbum)
```

### **Delete a Relation***

```Go
err = client.Resources.DeleteRelation("test:Group", "12345",
                                      "test:Albums",
                                      "test:Album", "23456")
```

### **Delete all Relations**

```Go
err = client.Resources.DeleteAllRelations("test:Group", "12345",
                                          "test:Albums")
```


----

### **Contributing**

 - Fork it
 - Create your feature branch (git checkout -b my-new-feature)
 - Commit your changes (git commit -am 'Add some feature')
 - Push to the branch (git push origin my-new-feature)
 - Create new Pull Request
 - If applicable, update the README.md
