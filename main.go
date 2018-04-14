package main

import (
	"github.com/ukayani/cloudformation-plus/yaml"
	"flag"
	"fmt"
	"os"
	"io/ioutil"
)


func failf(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func printUsage() {
	flag.Usage()
	os.Exit(1)
}

func main() {
	var removeAliases = flag.Bool("resolve-aliases", false, "Resolve all aliases to their target nodes")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of: %s [source]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) < 1 {
		printUsage()
	}

	path := flag.Arg(0)

	data,err := ioutil.ReadFile(path)

	failf(err)

	node, err := yaml.UnmarshalToTree([]byte(data), false)

	failf(err)

	out, err := yaml.MarshalFromTree(node, *removeAliases)

	failf(err)
	print(string(out))

}