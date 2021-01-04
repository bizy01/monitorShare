package spider

import (
	"testing"
)

func TestGetIndex(t *testing.T) {
   stock := NewStockIndex()
   stock.GetIndex()
   t.Log("ok")
}
