package doccleaner_test

import "github.com/garnaud/doccleaner"
import "encoding/json"
import "fmt"
import "gopkg.in/mgo.v2/bson"
import "github.com/stretchr/testify/assert"
import "strings"
import "testing"

type constantValueCleaner struct {
	changed []interface{}
}

func (c constantValueCleaner) Clean(value interface{}, args ...interface{}) (changed interface{}, err error) {
	if c.changed == nil {
		c.changed = make([]interface{}, 0)
	}

	fmt.Printf("new change: %+v on existing change %+v\n", value, c.changed)
	c.changed = append(c.changed, value)
	fmt.Printf("changed! %+v\n", c.changed)
	switch value.(type) {
	default:
		return "xxx", nil
	case int, float64:
		return 1234, nil
	}
}

func TestCleanFromToml(t *testing.T) {
	// given
	config := `
["leaf1"]
method="constantTransfo"
args=[]

["node2.leaf2"]
method="constantTransfo"
args=[]

["node2.leaf4"]
method="constantTransfo"
args=[]

["node3.node31.node311.node3111.leaf32"]
method="constantTransfo"
args=[]

["leaf4"]
method="constantTransfo"
args=[]

["node5.node51.node511.leaf5"]
method="constantTransfo"
args=[]

`
	input := `{
	 "leaf1":"value1",
	 "node2":[{"leaf2":"value2","leaf3":"value3","leaf4":"value4"}],
	 "node3":[{"node31":[{"node311":{"node3111":[{"leaf31":"abc"},{"leaf32":"cde"}]}}]}],
	 "leaf4":666,
	 "node5":[{"node51":{"node511":{"leaf5":5111}}}]
  }`
	expected := `{
	 "leaf1":"xxx",
	 "node2":[{"leaf2":"xxx","leaf3":"value3","leaf4":"xxx"}],
	 "node3":[{"node31":[{"node311":{"node3111":[{"leaf31":"abc"},{"leaf32":"xxx"}]}}]}],
	 "leaf4":1234,
	 "node5":[{"node51":{"node511":{"leaf5":1234}}}]
  }`

	cleaners := make(map[string]doccleaner.ValueCleaner)
	cleaners["constantTransfo"] = &constantValueCleaner{}
	docCleaner := doccleaner.NewDocCleaner(strings.NewReader(config), cleaners)

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(input), &obj); err == nil {
		// test
		docCleaner.Clean(obj)

		// check
		output, _ := json.Marshal(obj)
		assert.JSONEq(t, expected, string(output))
	} else {
		t.Error(err)
	}

}

func TestBsonClean(t *testing.T) {
	// given
	config := `
  ["leaf1"]
	"method"="constantTransfo"
	"args" = []
`
	cleaners := make(map[string]doccleaner.ValueCleaner)
	cleaners["constantTransfo"] = &constantValueCleaner{}
	docCleaner := doccleaner.NewDocCleaner(strings.NewReader(config), cleaners)
	doc := bson.M{}
	doc["leaf1"] = "toto"
	expected := `{"leaf1":"xxx"}`

	// test
	if _, err := docCleaner.Clean(doc); err == nil {
		// check
		assert.NotNil(t, doc)
		output, _ := json.Marshal(doc)
		assert.JSONEq(t, expected, string(output))
	} else {
		t.Error(err)
	}
}

func TestSimpleArray(t *testing.T) {
	// given
	config := `
  ["leaf1"]
	"method"="constantTransfo"
	"args" = []
`
	cleaners := make(map[string]doccleaner.ValueCleaner)
	cleaners["constantTransfo"] = &constantValueCleaner{}
	docCleaner := doccleaner.NewDocCleaner(strings.NewReader(config), cleaners)

	input := `["tata","toto","tutu"]`
	expected := `["xxx","xxx","xxx"]`

	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err == nil {
		// test
		docCleaner.Clean(obj)

		// check
		output, _ := json.Marshal(obj)
		assert.JSONEq(t, expected, string(output))
	} else {
		t.Error(err)
	}
}

func TestMapArray(t *testing.T) {
	// given
	config := `
  ["leaf1"]
	"method"="constantTransfo"
	"args" = []
	`
	cleaners := make(map[string]doccleaner.ValueCleaner)
	cleaners["constantTransfo"] = &constantValueCleaner{}
	docCleaner := doccleaner.NewDocCleaner(strings.NewReader(config), cleaners)

	input := `[{"leaf1":"value1","leaf2":"value2"}]`
	expected := `[{"leaf1":"xxx","leaf2":"value2"}]`

	var obj interface{}
	if err := json.Unmarshal([]byte(input), &obj); err == nil {
		// test
		docCleaner.Clean(obj)

		// check
		output, _ := json.Marshal(obj)
		assert.JSONEq(t, expected, string(output))
	} else {
		t.Error(err)
	}
}
