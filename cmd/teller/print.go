package main

import (
	"fmt"
	"sort"
	"time"
)

func print(data map[time.Month]string) {
	keys := []time.Month{}
	for m := range data {
		keys = append(keys, m)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, m := range keys {
		fmt.Printf("%s\n------------\n", m.String())
		fmt.Println(data[m])
	}
}
