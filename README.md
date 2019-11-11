# kv2struct
golang key values convert to struct
ref [youkale/go-querystruct](https://github.com/youkale/go-querystruct)

### Usage
```
package main

import (
	"./kv2struct"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strconv"
)

type Person struct {
	Dt int64 `param:"dt,656"`
}
type Student struct {
	SC int64 `param:"sc,130"`
}

type User struct {
	Student

	UserId  int64   `param:"user_id,100"`
	StoreId int     `param:"store_id"`
	Page    float64 `param:"page"`
	Name    string  `param:"name"`
	Age     uint8   `json:"xx"`
	Enable  bool    `param:"enable,false"`
}

type UrlValues url.Values

func (v UrlValues) GetString(key string, def ...string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func main() {
	o := User{}

	sc := rand.Int63()
	userId := rand.Int63()
	storeId := rand.Int()
	page := rand.Float64()
	age := 26
	want := UrlValues{
		"sc":       {fmt.Sprintf("%v", sc)},
		"store_id": {fmt.Sprintf("%v", storeId)},
		"user_id":  {strconv.FormatInt(userId, 10)},
		"page":     {fmt.Sprintf("%v", page)},
		"name":     {"sdfdsfs"},
		"Age":      {fmt.Sprintf("%v", age)},
	}

	kv2struct.SetTagKey("param")

	e := kv2struct.Unmarshal(want, &o)

	if e == nil {
		fmt.Println(o)
		if o.StoreId != storeId || o.UserId != userId || o.Page != page {
			log.Println("has error ")
		} else {
			fmt.Println("ok")
		}
	} else {
		log.Println(e)
	}

	data, _ := json.Marshal(o)
	fmt.Println(string(data))
}

```
