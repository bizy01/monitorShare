package spider

import (
	"testing"
)

func TestGetData(t *testing.T) {
   north := NewNorthFund()
   north.GetData()
   t.Log("ok")
}