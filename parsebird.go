/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : parsebird.go

* Purpose :

* Creation Date : 04-04-2014

* Last Modified : Fri 04 Apr 2014 08:14:29 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/kiyor/parsebird/lib"
)

var (
	conff *string = flag.String("f", "/etc/bird/bird.conf", "bird conf path")
	min   *bool   = flag.Bool("min", false, "min output")
)

func init() {
	flag.Parse()
}

func main() {
	bird := parsebird.ParseConf(*conff)
	var j []byte
	if *min {
		j, _ = json.Marshal(bird)
	} else {
		j, _ = json.MarshalIndent(bird, "", "    ")
	}
	fmt.Println(string(j))
}
