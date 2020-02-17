package main

import (
	"testing"
)

func TestPortMap(t *testing.T) {
	cases := []struct {
		Want   string
		Expect int
	}{
		{"abc.net", 9090},
		{"blah.org", 9999},
		{"okokok.org", 9002},
	}

	LoadConfig("./whatever.config")

	for _, c := range cases {
		if p := PortMap(c.Want); p != c.Expect {
			t.Errorf("PortMap(%s) equals %d, but it should equal %d", c.Want, p, c.Expect)
		}
	}
}

//
func TestConcurrency(t *testing.T) {

	LoadConfig("./whatever.config")

	possible_urls := []string{"abc.net", "okokok.org"}

	for i := 0; i < 1000; i++ {
		go PortMap(possible_urls[i%len(possible_urls)])
	}
}
