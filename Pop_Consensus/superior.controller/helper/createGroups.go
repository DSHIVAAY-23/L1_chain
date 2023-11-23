package superiorhelper

import "math/rand"

func GroupGenerator(ips *[]string) {
	n := len(*ips)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(100000000) % n
		(*ips)[i], (*ips)[j] = (*ips)[j], (*ips)[i]
	}
}
