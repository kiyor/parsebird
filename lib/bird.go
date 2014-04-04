/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : bird.go

* Purpose :

* Creation Date : 04-04-2014

* Last Modified : Fri 04 Apr 2014 07:52:36 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package parsebird

import (
	"github.com/kiyor/subnettool"
	"io/ioutil"
	"net"
	"regexp"
	"strings"
)

type Bird map[string]Edge
type Edge map[string][]Route

type Route struct {
	Ip     net.IP
	active bool
}

var (
	reProtocol     = regexp.MustCompile(`protocol static (.*) {`)
	reEdge         = regexp.MustCompile(`Edge(\d\d)`)
	reRoute        = regexp.MustCompile(`route (.*) reject;`)
	reRouteComment = regexp.MustCompile(`#.*route (.*) reject;`)
	reEndEdge      = regexp.MustCompile(`EndEdge`)
)

func chkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func ParseConf(f string) Bird {
	b, err := ioutil.ReadFile(f)
	chkErr(err)
	strs := strings.Split(string(b), "\n")
	strs = strs[:len(strs)-1]
	bird := make(Bird)
	var providerKey, edgeKey string
	var endEdge bool
	var rs []Route
	for _, v := range strs {
		if reProtocol.MatchString(v) {
			providerKey = reProtocol.FindStringSubmatch(v)[1]
			if _, ok := bird[providerKey]; !ok {
				edge := make(Edge)
				bird[providerKey] = edge
			}
			endEdge = false
		}
		if reEdge.MatchString(v) {
			edgeKey = "edge" + reEdge.FindStringSubmatch(v)[1]
			rs = []Route{}
			bird[providerKey][edgeKey] = rs
		}
		if reRoute.MatchString(v) && !endEdge {
			block := reRoute.FindStringSubmatch(v)[1]
			ips := subnettool.GetAllIP(block)
			for _, ip := range ips {
				var r Route
				r.Ip = ip
				r.active = true
				if reRouteComment.MatchString(v) {
					r.active = false
				}
				token := subnettool.ParseIPInt(r.Ip)
				if token[3] != 255 && token[3] != 0 && token[3] != 1 && token[3] != 2 { // if ip is boardcast or swith
					rs = append(rs, r)
				}
			}
			rs = RemoveDup(rs)
			bird[providerKey][edgeKey] = rs
		}
		if reEndEdge.MatchString(v) {
			endEdge = true
		}
	}
	// 	fmt.Println(bird)
	for _, v := range bird { // delete empty key map
		delete(v, "")
	}
	// 	j, _ := json.MarshalIndent(bird, "", "    ")
	return bird
}

func RemoveDup(rs []Route) []Route {
	result := []Route{}
	seen := map[string]Route{}
	for _, r := range rs {
		if _, ok := seen[r.Ip.String()]; !ok {
			if seen[r.Ip.String()].active || r.active { // if any ip is active then is active
				result = append(result, r)
				seen[r.Ip.String()] = r
			}
		}
	}
	return result
}
