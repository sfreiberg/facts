Facts
=====

Facts was originally included in [Shepherd](https://github.com/sfreiberg/shepherd) to gather facts about the operating system. It makes sense to break it out into a separate project so that they can move at different speeds.

Install
=======

```
$ go get github.com/sfreiberg/facts
```

Documentation
=============

http://godoc.org/github.com/sfreiberg/facts


Simple Example
==============
```
package main

import (
	"github.com/sfreiberg/facts"

	"fmt"
)

func main() {
	f := factsFindFacts()
	json, err := f.ToPrettyJson()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", json)
}
```


Current Features
================

* Get hostname, domain, cpu, os and interfaces

Planned Features
================

* BSD support
* Additional facts
  * Memory
  * Mounts

Current Limitations
===================

* Only works on Linux, Mac OS X
* Lack of testing
