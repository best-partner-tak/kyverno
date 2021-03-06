// A comment just to push the positions out

package a

import (
	"fmt"
	"sync"
)

type A string //@A

func AStuff() { //@AStuff
	x := 5
	Random2(x) //@godef("dom2", Random2)
	Random()   //@godef("()", Random)

	var err error         //@err
	fmt.Printf("%v", err) //@godef("err", err)

	var y string       //@string,hover("string", string)
	_ = make([]int, 0) //@make,hover("make", make)

	var mu sync.Mutex
	mu.Lock() //@Lock,hover("Lock", Lock)
}
