package spider

import (
	"testing"
)

func TestGetTopData(t *testing.T) {
   top := NewNorthTop()
   top.GetData()
   t.Log("ok")
}