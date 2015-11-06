package main

import (
	"testing"
)

func Test_gupiao(t *testing.T) {
	t.Log(GuPiao("600360"))
	t.Log(GuPiao("000761"))
	t.Log(GuPiao("300444"))
}
