package yaml

func UnmarshalToTree(in []byte, strict bool) (node *Node, err error) {
	defer handleErr(&err)
	p := newParser(in)
	defer p.destroy()
	node = p.parse()
	return
}

func MarshalFromTree(in *Node, removeAliases bool, normalize bool) (out []byte, err error) {
	defer handleErr(&err)
	e := newNodeEncoder()
	defer e.destroy()
	e.marshalDoc(in, removeAliases, normalize)
	e.finish()
	out = e.out
	return
}
