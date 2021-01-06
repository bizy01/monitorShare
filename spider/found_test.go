package spider

import (
	"testing"
)

func TestFundRank(t *testing.T) {
   found := NewFundRank()
   found.GetData()
   t.Log("ok")
}
