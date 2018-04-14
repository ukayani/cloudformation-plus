package yaml

import "io"

type nodeEncoder struct {
	emitter yaml_emitter_t
	event   yaml_event_t
	out     []byte
	flow    bool
	// doneInit holds whether the initial stream_start_event has been
	// emitted.
	doneInit bool
	removeAliases bool
}

func newNodeEncoder() *nodeEncoder {
	e := &nodeEncoder{}
	yaml_emitter_initialize(&e.emitter)
	yaml_emitter_set_output_string(&e.emitter, &e.out)
	yaml_emitter_set_unicode(&e.emitter, true)
	return e
}

func newNodeEncoderWithWriter(w io.Writer) *nodeEncoder {
	e := &nodeEncoder{}
	yaml_emitter_initialize(&e.emitter)
	yaml_emitter_set_output_writer(&e.emitter, w)
	yaml_emitter_set_unicode(&e.emitter, true)
	return e
}

func (e *nodeEncoder) init() {
	if e.doneInit {
		return
	}
	yaml_stream_start_event_initialize(&e.event, yaml_UTF8_ENCODING)
	e.emit()
	e.doneInit = true
}

func (e *nodeEncoder) finish() {
	e.emitter.open_ended = false
	yaml_stream_end_event_initialize(&e.event)
	e.emit()
}

func (e *nodeEncoder) destroy() {
	yaml_emitter_delete(&e.emitter)
}

func (e *nodeEncoder) emit() {
	// This will internally delete the e.event Value.
	e.must(yaml_emitter_emit(&e.emitter, &e.event))
}

func (e *nodeEncoder) must(ok bool) {
	if !ok {
		msg := e.emitter.problem
		if msg == "" {
			msg = "unknown problem generating YAML content"
		}
		failf("%s", msg)
	}
}

func (e *nodeEncoder) marshalDoc(in *Node, removeAliases bool) {
	e.init()
	e.removeAliases = removeAliases
	e.must(in.Kind == DocumentNode)
	yaml_document_start_event_initialize(&e.event, nil, nil, true)
	e.emit()
	for _,c := range in.Children {
		e.marshal(c)
	}
	yaml_document_end_event_initialize(&e.event, true)
	e.emit()
}

func (e *nodeEncoder) marshal(in *Node) {

	switch in.Kind {
	case MappingNode:
		e.emitMapping(in)
	case SequenceNode:
		e.emitSequence(in)
	case AliasNode:
		e.emitAlias(in)
	case ScalarNode:
		e.emitScalar(in)
	}
}

func (e *nodeEncoder) emitMapping(in *Node) {
	implicit := in.Tag == ""
	anchor := ""

	if !e.removeAliases {
		anchor = in.Anchor
	}

	yaml_mapping_start_event_initialize(&e.event, []byte(anchor), []byte(in.Tag), implicit, yaml_BLOCK_MAPPING_STYLE)
	e.emit()

	// first pass through keys of this map to see if any are special merge keys
	// gather all values and resolve them and then merge them according to the spec with the current node
	// then process the current node as is, with all merge keys resolved


	var l = len(in.Children)

	if e.removeAliases {
		for i := 0; i < l; i +=2 {
			key := in.Children[i]
			if isMerge(key) {
				e.merge(in, in.Children[i + 1])
			}
		}
		// update length since it most likely changed due to merge
		l = len(in.Children)
	}

	for i := 0; i < l; i += 2 {
		key := in.Children[i]
		value := in.Children[i+1]
		if !e.removeAliases || !isMerge(key) {
			e.marshal(key)
			e.marshal(value)
		}
	}

	yaml_mapping_end_event_initialize(&e.event)
	e.emit()
}

func (e *nodeEncoder) merge(a *Node, b *Node) {
	switch b.Kind {
	case MappingNode:
		e.mergeMapping(a, b)
	case SequenceNode:
		e.mergeSequence(a, b)
	case AliasNode:
		e.merge(a, b.Alias)
	default:
		failf("Illegal value type (%d) for merge key", b.Kind)
	}
}

func (e *nodeEncoder) mergeMapping(a *Node, b *Node) {
	var keyMap = make(map[string]bool)
	var la = len(a.Children)
	for i := 0; i < la; i += 2 {
		key := a.Children[i]
		if key.Kind == ScalarNode {
			keyMap[key.Value] = true
		}
	}

	var lb = len(b.Children)
	for i := 0; i < lb; i += 2 {
		key := b.Children[i]
		value := b.Children[i + 1]
		if key.Kind != ScalarNode || !keyMap[key.Value] {
			a.Children = append(a.Children, key, value)
		}
	}
}

func (e *nodeEncoder) mergeSequence(a *Node, b *Node) {
	for _, c := range b.Children {
		switch c.Kind {
		case AliasNode:
			if c.Alias.Kind != MappingNode {
				failf("Illegal value type (%d) in sequence for merge key", c.Kind)
			}
			e.mergeMapping(a, c.Alias)
		case MappingNode:
			e.mergeMapping(a, c)
		default:
			failf("Illegal value type (%d) in sequence for merge key", c.Kind)
		}
	}
}

func (e *nodeEncoder) emitSequence(in *Node) {
	implicit := in.Tag == ""

	anchor := ""

	if !e.removeAliases {
		anchor = in.Anchor
	}

	yaml_sequence_start_event_initialize(&e.event, []byte(anchor), []byte(in.Tag), implicit, yaml_BLOCK_SEQUENCE_STYLE)
	e.emit()

	for _,c := range in.Children {
		e.marshal(c)
	}

	yaml_sequence_end_event_initialize(&e.event)
	e.emit()

}

func (e *nodeEncoder) emitAlias(in *Node) {
	if e.removeAliases {
		e.marshal(in.Alias)
		return
	}

	e.must(yaml_alias_event_initialize(&e.event, []byte(in.Value)))
	e.emit()
}

func (e *nodeEncoder) emitScalar(in *Node) {
	tag := in.Tag

	anchor := ""

	if !e.removeAliases {
		anchor = in.Anchor
	}

	value := in.Value
	implicit := tag == ""
	style := yaml_ANY_SCALAR_STYLE
	if !implicit {
		style = yaml_SINGLE_QUOTED_SCALAR_STYLE
	}
	e.must(yaml_scalar_event_initialize(&e.event, []byte(anchor), []byte(tag), []byte(value), implicit, implicit, style))
	e.emit()
}
