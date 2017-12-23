package main

import (
	"fmt"
	"os"
)

func main() {
	var obj Indexer = Indexer{}
	PATH := os.Getenv("CRAWLER_PATH")
	obj.OpenCon()
	obj.ReadFolder(fmt.Sprintf("%s", PATH))
	obj.CloseCon()
}
