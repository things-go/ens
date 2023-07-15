package sts

import (
	"fmt"
	"testing"

	"github.com/things-go/ens"
)

func Test(t *testing.T) {
	v := ens.Entity{}
	fmt.Println(Parse(v))
}
