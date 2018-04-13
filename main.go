package main

import (
	"github.com/ukayani/cloudformation-plus/yaml"
	"fmt"
)

func walk(node *yaml.Node) {
	switch kind := node.Kind; kind {
	case yaml.AliasNode:
		fmt.Println(node.Alias.Value)
	case yaml.ScalarNode:
		fmt.Println(node.Value)
	}
	for _, c := range node.Children {
		walk(c)
	}
}
func main() {
	var data = `
a: &test !Ref 'a string from struct A'
b: !Ref 'a string from struct B'
c: &hi
 c1: test
 c2: another
cP: &hi2
 inner:
  <<: *hi
 c1: testP
 c3: testP3
d: &howdy
- Testing: !Ref 'Hello'
- Another: {"test":"hello"}
e: *hi
f:
 <<: *hi2
 in: hello
`

	var node = yaml.GetTree([]byte(data))

	walk(node)

	node.Children[0].Children[1].Tag = "!Hello"
	var enc = yaml.NewNodeEncoder()
	defer enc.Destroy()
	enc.MarshalDoc(node, true)
	enc.Finish()
	println("Dumping node")
	print(string(enc.Out))

}