package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alpha = "abcdefghijklmnopgrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomServerName() string {
	var sb strings.Builder

	sb.WriteString("VCS")

	l := len(alpha)

	for i := 0; i < 3; i++ {
		sb.WriteByte(alpha[rand.Intn(l)])
	}
	fmt.Print(sb.String())
	return sb.String()
}

func RandomStatus() int {
	return rand.Intn(2)
}

func RandomIP() string {
	cur := []string{"1.1.1.1", "66.220.149.25", "143.166.83.38", "72.30.2.43", "2.2.2.2", "1.2.3.4"}
	return cur[rand.Intn(len(cur))]
}
