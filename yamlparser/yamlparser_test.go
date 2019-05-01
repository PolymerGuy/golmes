package yamlparser_test

import (
	"testing"
	"github.com/PolymerGuy/golmes/yamlparser"
	"fmt"
)

func TestComparatorsFromArgs(t *testing.T) {
	data :=` DataComparators:
 - type: synced
   referencefile : reffile
   currentfile : curfile
   commonargsfile :
   keywords:
    - key1
    - key2

 Abaqus_settings:
  path: /path/to/abaqus`

  parser := yamlparser.Parse([]byte(data))

  comparator := parser.NewComparator()[0]

  datas:=comparator.GetFields()
  fmt.Println(datas[1])

  if fmt.Sprint(datas[0]) != "{{reffile key1} {reffile key2}}"{
  	t.Errorf("The current files have the wrong formatting!")
	}
  if fmt.Sprint(datas[1]) != "{{curfile key1} {curfile key2}}"{
  	t.Errorf("The current files have the wrong formatting!")
  }
}
