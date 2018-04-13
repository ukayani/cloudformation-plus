package yaml

import (
	"io"
)

func yaml_insert_token(parser *yaml_parser_t, pos int, token *yaml_token_t) {
	//fmt.Println("yaml_insert_token", "pos:", pos, "typ:", token.typ, "head:", parser.tokens_head, "len:", len(parser.tokens))

	// Check if we can move the queue at the beginning of the buffer.
	if parser.tokens_head > 0 && len(parser.tokens) == cap(parser.tokens) {
		if parser.tokens_head != len(parser.tokens) {
			copy(parser.tokens, parser.tokens[parser.tokens_head:])
		}
		parser.tokens = parser.tokens[:len(parser.tokens)-parser.tokens_head]
		parser.tokens_head = 0
	}
	parser.tokens = append(parser.tokens, *token)
	if pos < 0 {
		return
	}
	copy(parser.tokens[parser.tokens_head+pos+1:], parser.tokens[parser.tokens_head+pos:])
	parser.tokens[parser.tokens_head+pos] = *token
}

// Create a new parser object.
func yaml_parser_initialize(parser *yaml_parser_t) bool {
	*parser = yaml_parser_t{
		raw_buffer: make([]byte, 0, input_raw_buffer_size),
		buffer:     make([]byte, 0, input_buffer_size),
	}
	return true
}

// Destroy a parser object.
func yaml_parser_delete(parser *yaml_parser_t) {
	*parser = yaml_parser_t{}
}

// String read handler.
func yaml_string_read_handler(parser *yaml_parser_t, buffer []byte) (n int, err error) {
	if parser.input_pos == len(parser.input) {
		return 0, io.EOF
	}
	n = copy(buffer, parser.input[parser.input_pos:])
	parser.input_pos += n
	return n, nil
}

// Reader read handler.
func yaml_reader_read_handler(parser *yaml_parser_t, buffer []byte) (n int, err error) {
	return parser.input_reader.Read(buffer)
}

// Set a string input.
func yaml_parser_set_input_string(parser *yaml_parser_t, input []byte) {
	if parser.read_handler != nil {
		panic("must set the input source only once")
	}
	parser.read_handler = yaml_string_read_handler
	parser.input = input
	parser.input_pos = 0
}

// Set a file input.
func yaml_parser_set_input_reader(parser *yaml_parser_t, r io.Reader) {
	if parser.read_handler != nil {
		panic("must set the input source only once")
	}
	parser.read_handler = yaml_reader_read_handler
	parser.input_reader = r
}

// Set the source encoding.
func yaml_parser_set_encoding(parser *yaml_parser_t, encoding yaml_encoding_t) {
	if parser.encoding != yaml_ANY_ENCODING {
		panic("must set the encoding only once")
	}
	parser.encoding = encoding
}

// Create a new emitter object.
func yaml_emitter_initialize(emitter *yaml_emitter_t) {
	*emitter = yaml_emitter_t{
		buffer:     make([]byte, output_buffer_size),
		raw_buffer: make([]byte, 0, output_raw_buffer_size),
		states:     make([]yaml_emitter_state_t, 0, initial_stack_size),
		events:     make([]yaml_event_t, 0, initial_queue_size),
	}
}

// Destroy an emitter object.
func yaml_emitter_delete(emitter *yaml_emitter_t) {
	*emitter = yaml_emitter_t{}
}

// String write handler.
func yaml_string_write_handler(emitter *yaml_emitter_t, buffer []byte) error {
	*emitter.output_buffer = append(*emitter.output_buffer, buffer...)
	return nil
}

// yaml_writer_write_handler uses emitter.output_writer to write the
// emitted text.
func yaml_writer_write_handler(emitter *yaml_emitter_t, buffer []byte) error {
	_, err := emitter.output_writer.Write(buffer)
	return err
}

// Set a string output.
func yaml_emitter_set_output_string(emitter *yaml_emitter_t, output_buffer *[]byte) {
	if emitter.write_handler != nil {
		panic("must set the output target only once")
	}
	emitter.write_handler = yaml_string_write_handler
	emitter.output_buffer = output_buffer
}

// Set a file output.
func yaml_emitter_set_output_writer(emitter *yaml_emitter_t, w io.Writer) {
	if emitter.write_handler != nil {
		panic("must set the output target only once")
	}
	emitter.write_handler = yaml_writer_write_handler
	emitter.output_writer = w
}

// Set the output encoding.
func yaml_emitter_set_encoding(emitter *yaml_emitter_t, encoding yaml_encoding_t) {
	if emitter.encoding != yaml_ANY_ENCODING {
		panic("must set the output encoding only once")
	}
	emitter.encoding = encoding
}

// Set the canonical output style.
func yaml_emitter_set_canonical(emitter *yaml_emitter_t, canonical bool) {
	emitter.canonical = canonical
}

//// Set the indentation increment.
func yaml_emitter_set_indent(emitter *yaml_emitter_t, indent int) {
	if indent < 2 || indent > 9 {
		indent = 2
	}
	emitter.best_indent = indent
}

// Set the preferred line width.
func yaml_emitter_set_width(emitter *yaml_emitter_t, width int) {
	if width < 0 {
		width = -1
	}
	emitter.best_width = width
}

// Set if unescaped non-ASCII characters are allowed.
func yaml_emitter_set_unicode(emitter *yaml_emitter_t, unicode bool) {
	emitter.unicode = unicode
}

// Set the preferred line break character.
func yaml_emitter_set_break(emitter *yaml_emitter_t, line_break yaml_break_t) {
	emitter.line_break = line_break
}

///*
// * Destroy a token object.
// */
//
//YAML_DECLARE(void)
//yaml_token_delete(yaml_token_t *token)
//{
//    assert(token);  // Non-NULL token object expected.
//
//    switch (token.type)
//    {
//        case YAML_TAG_DIRECTIVE_TOKEN:
//            yaml_free(token.data.tag_directive.handle);
//            yaml_free(token.data.tag_directive.prefix);
//            break;
//
//        case YAML_ALIAS_TOKEN:
//            yaml_free(token.data.Alias.Value);
//            break;
//
//        case YAML_ANCHOR_TOKEN:
//            yaml_free(token.data.anchor.Value);
//            break;
//
//        case YAML_TAG_TOKEN:
//            yaml_free(token.data.Tag.handle);
//            yaml_free(token.data.Tag.suffix);
//            break;
//
//        case YAML_SCALAR_TOKEN:
//            yaml_free(token.data.scalar.Value);
//            break;
//
//        default:
//            break;
//    }
//
//    memset(token, 0, sizeof(yaml_token_t));
//}
//
///*
// * Check if a string is a valid UTF-8 sequence.
// *
// * Check 'reader.c' for more details on UTF-8 encoding.
// */
//
//static int
//yaml_check_utf8(yaml_char_t *start, size_t length)
//{
//    yaml_char_t *end = start+length;
//    yaml_char_t *pointer = start;
//
//    while (pointer < end) {
//        unsigned char octet;
//        unsigned int width;
//        unsigned int Value;
//        size_t k;
//
//        octet = pointer[0];
//        width = (octet & 0x80) == 0x00 ? 1 :
//                (octet & 0xE0) == 0xC0 ? 2 :
//                (octet & 0xF0) == 0xE0 ? 3 :
//                (octet & 0xF8) == 0xF0 ? 4 : 0;
//        Value = (octet & 0x80) == 0x00 ? octet & 0x7F :
//                (octet & 0xE0) == 0xC0 ? octet & 0x1F :
//                (octet & 0xF0) == 0xE0 ? octet & 0x0F :
//                (octet & 0xF8) == 0xF0 ? octet & 0x07 : 0;
//        if (!width) return 0;
//        if (pointer+width > end) return 0;
//        for (k = 1; k < width; k ++) {
//            octet = pointer[k];
//            if ((octet & 0xC0) != 0x80) return 0;
//            Value = (Value << 6) + (octet & 0x3F);
//        }
//        if (!((width == 1) ||
//            (width == 2 && Value >= 0x80) ||
//            (width == 3 && Value >= 0x800) ||
//            (width == 4 && Value >= 0x10000))) return 0;
//
//        pointer += width;
//    }
//
//    return 1;
//}
//

// Create STREAM-START.
func yaml_stream_start_event_initialize(event *yaml_event_t, encoding yaml_encoding_t) {
	*event = yaml_event_t{
		typ:      yaml_STREAM_START_EVENT,
		encoding: encoding,
	}
}

// Create STREAM-END.
func yaml_stream_end_event_initialize(event *yaml_event_t) {
	*event = yaml_event_t{
		typ: yaml_STREAM_END_EVENT,
	}
}

// Create DOCUMENT-START.
func yaml_document_start_event_initialize(
	event *yaml_event_t,
	version_directive *yaml_version_directive_t,
	tag_directives []yaml_tag_directive_t,
	implicit bool,
) {
	*event = yaml_event_t{
		typ:               yaml_DOCUMENT_START_EVENT,
		version_directive: version_directive,
		tag_directives:    tag_directives,
		implicit:          implicit,
	}
}

// Create DOCUMENT-END.
func yaml_document_end_event_initialize(event *yaml_event_t, implicit bool) {
	*event = yaml_event_t{
		typ:      yaml_DOCUMENT_END_EVENT,
		implicit: implicit,
	}
}

///*
// * Create ALIAS.
// */
//
//YAML_DECLARE(int)
//yaml_alias_event_initialize(event *yaml_event_t, anchor *yaml_char_t)
//{
//    mark yaml_mark_t = { 0, 0, 0 }
//    anchor_copy *yaml_char_t = NULL
//
//    assert(event) // Non-NULL event object is expected.
//    assert(anchor) // Non-NULL anchor is expected.
//
//    if (!yaml_check_utf8(anchor, strlen((char *)anchor))) return 0
//
//    anchor_copy = yaml_strdup(anchor)
//    if (!anchor_copy)
//        return 0
//
//    ALIAS_EVENT_INIT(*event, anchor_copy, mark, mark)
//
//    return 1
//}

func yaml_alias_event_initialize(event *yaml_event_t, anchor []byte) bool {
	*event = yaml_event_t{
		typ: yaml_ALIAS_EVENT,
		anchor: anchor,
	}

	return true
}

// Create SCALAR.
func yaml_scalar_event_initialize(event *yaml_event_t, anchor, tag, value []byte, plain_implicit, quoted_implicit bool, style yaml_scalar_style_t) bool {
	*event = yaml_event_t{
		typ:             yaml_SCALAR_EVENT,
		anchor:          anchor,
		tag:             tag,
		value:           value,
		implicit:        plain_implicit,
		quoted_implicit: quoted_implicit,
		style:           yaml_style_t(style),
	}
	return true
}

// Create SEQUENCE-START.
func yaml_sequence_start_event_initialize(event *yaml_event_t, anchor, tag []byte, implicit bool, style yaml_sequence_style_t) bool {
	*event = yaml_event_t{
		typ:      yaml_SEQUENCE_START_EVENT,
		anchor:   anchor,
		tag:      tag,
		implicit: implicit,
		style:    yaml_style_t(style),
	}
	return true
}

// Create SEQUENCE-END.
func yaml_sequence_end_event_initialize(event *yaml_event_t) bool {
	*event = yaml_event_t{
		typ: yaml_SEQUENCE_END_EVENT,
	}
	return true
}

// Create MAPPING-START.
func yaml_mapping_start_event_initialize(event *yaml_event_t, anchor, tag []byte, implicit bool, style yaml_mapping_style_t) {
	*event = yaml_event_t{
		typ:      yaml_MAPPING_START_EVENT,
		anchor:   anchor,
		tag:      tag,
		implicit: implicit,
		style:    yaml_style_t(style),
	}
}

// Create MAPPING-END.
func yaml_mapping_end_event_initialize(event *yaml_event_t) {
	*event = yaml_event_t{
		typ: yaml_MAPPING_END_EVENT,
	}
}

// Destroy an event object.
func yaml_event_delete(event *yaml_event_t) {
	*event = yaml_event_t{}
}

///*
// * Create a document object.
// */
//
//YAML_DECLARE(int)
//yaml_document_initialize(document *yaml_document_t,
//        version_directive *yaml_version_directive_t,
//        tag_directives_start *yaml_tag_directive_t,
//        tag_directives_end *yaml_tag_directive_t,
//        start_implicit int, end_implicit int)
//{
//    struct {
//        error yaml_error_type_t
//    } context
//    struct {
//        start *yaml_node_t
//        end *yaml_node_t
//        top *yaml_node_t
//    } nodes = { NULL, NULL, NULL }
//    version_directive_copy *yaml_version_directive_t = NULL
//    struct {
//        start *yaml_tag_directive_t
//        end *yaml_tag_directive_t
//        top *yaml_tag_directive_t
//    } tag_directives_copy = { NULL, NULL, NULL }
//    Value yaml_tag_directive_t = { NULL, NULL }
//    mark yaml_mark_t = { 0, 0, 0 }
//
//    assert(document) // Non-NULL document object is expected.
//    assert((tag_directives_start && tag_directives_end) ||
//            (tag_directives_start == tag_directives_end))
//                            // Valid Tag directives are expected.
//
//    if (!STACK_INIT(&context, nodes, INITIAL_STACK_SIZE)) goto error
//
//    if (version_directive) {
//        version_directive_copy = yaml_malloc(sizeof(yaml_version_directive_t))
//        if (!version_directive_copy) goto error
//        version_directive_copy.major = version_directive.major
//        version_directive_copy.minor = version_directive.minor
//    }
//
//    if (tag_directives_start != tag_directives_end) {
//        tag_directive *yaml_tag_directive_t
//        if (!STACK_INIT(&context, tag_directives_copy, INITIAL_STACK_SIZE))
//            goto error
//        for (tag_directive = tag_directives_start
//                tag_directive != tag_directives_end; tag_directive ++) {
//            assert(tag_directive.handle)
//            assert(tag_directive.prefix)
//            if (!yaml_check_utf8(tag_directive.handle,
//                        strlen((char *)tag_directive.handle)))
//                goto error
//            if (!yaml_check_utf8(tag_directive.prefix,
//                        strlen((char *)tag_directive.prefix)))
//                goto error
//            Value.handle = yaml_strdup(tag_directive.handle)
//            Value.prefix = yaml_strdup(tag_directive.prefix)
//            if (!Value.handle || !Value.prefix) goto error
//            if (!PUSH(&context, tag_directives_copy, Value))
//                goto error
//            Value.handle = NULL
//            Value.prefix = NULL
//        }
//    }
//
//    DOCUMENT_INIT(*document, nodes.start, nodes.end, version_directive_copy,
//            tag_directives_copy.start, tag_directives_copy.top,
//            start_implicit, end_implicit, mark, mark)
//
//    return 1
//
//error:
//    STACK_DEL(&context, nodes)
//    yaml_free(version_directive_copy)
//    while (!STACK_EMPTY(&context, tag_directives_copy)) {
//        Value yaml_tag_directive_t = POP(&context, tag_directives_copy)
//        yaml_free(Value.handle)
//        yaml_free(Value.prefix)
//    }
//    STACK_DEL(&context, tag_directives_copy)
//    yaml_free(Value.handle)
//    yaml_free(Value.prefix)
//
//    return 0
//}
//
///*
// * Destroy a document object.
// */
//
//YAML_DECLARE(void)
//yaml_document_delete(document *yaml_document_t)
//{
//    struct {
//        error yaml_error_type_t
//    } context
//    tag_directive *yaml_tag_directive_t
//
//    context.error = YAML_NO_ERROR // Eliminate a compiler warning.
//
//    assert(document) // Non-NULL document object is expected.
//
//    while (!STACK_EMPTY(&context, document.nodes)) {
//        Node yaml_node_t = POP(&context, document.nodes)
//        yaml_free(Node.Tag)
//        switch (Node.type) {
//            case YAML_SCALAR_NODE:
//                yaml_free(Node.data.scalar.Value)
//                break
//            case YAML_SEQUENCE_NODE:
//                STACK_DEL(&context, Node.data.sequence.items)
//                break
//            case YAML_MAPPING_NODE:
//                STACK_DEL(&context, Node.data.mapping.pairs)
//                break
//            default:
//                assert(0) // Should not happen.
//        }
//    }
//    STACK_DEL(&context, document.nodes)
//
//    yaml_free(document.version_directive)
//    for (tag_directive = document.tag_directives.start
//            tag_directive != document.tag_directives.end
//            tag_directive++) {
//        yaml_free(tag_directive.handle)
//        yaml_free(tag_directive.prefix)
//    }
//    yaml_free(document.tag_directives.start)
//
//    memset(document, 0, sizeof(yaml_document_t))
//}
//
///**
// * Get a document Node.
// */
//
//YAML_DECLARE(yaml_node_t *)
//yaml_document_get_node(document *yaml_document_t, index int)
//{
//    assert(document) // Non-NULL document object is expected.
//
//    if (index > 0 && document.nodes.start + index <= document.nodes.top) {
//        return document.nodes.start + index - 1
//    }
//    return NULL
//}
//
///**
// * Get the root object.
// */
//
//YAML_DECLARE(yaml_node_t *)
//yaml_document_get_root_node(document *yaml_document_t)
//{
//    assert(document) // Non-NULL document object is expected.
//
//    if (document.nodes.top != document.nodes.start) {
//        return document.nodes.start
//    }
//    return NULL
//}
//
///*
// * Add a scalar Node to a document.
// */
//
//YAML_DECLARE(int)
//yaml_document_add_scalar(document *yaml_document_t,
//        Tag *yaml_char_t, Value *yaml_char_t, length int,
//        style yaml_scalar_style_t)
//{
//    struct {
//        error yaml_error_type_t
//    } context
//    mark yaml_mark_t = { 0, 0, 0 }
//    tag_copy *yaml_char_t = NULL
//    value_copy *yaml_char_t = NULL
//    Node yaml_node_t
//
//    assert(document) // Non-NULL document object is expected.
//    assert(Value) // Non-NULL Value is expected.
//
//    if (!Tag) {
//        Tag = (yaml_char_t *)YAML_DEFAULT_SCALAR_TAG
//    }
//
//    if (!yaml_check_utf8(Tag, strlen((char *)Tag))) goto error
//    tag_copy = yaml_strdup(Tag)
//    if (!tag_copy) goto error
//
//    if (length < 0) {
//        length = strlen((char *)Value)
//    }
//
//    if (!yaml_check_utf8(Value, length)) goto error
//    value_copy = yaml_malloc(length+1)
//    if (!value_copy) goto error
//    memcpy(value_copy, Value, length)
//    value_copy[length] = '\0'
//
//    SCALAR_NODE_INIT(Node, tag_copy, value_copy, length, style, mark, mark)
//    if (!PUSH(&context, document.nodes, Node)) goto error
//
//    return document.nodes.top - document.nodes.start
//
//error:
//    yaml_free(tag_copy)
//    yaml_free(value_copy)
//
//    return 0
//}
//
///*
// * Add a sequence Node to a document.
// */
//
//YAML_DECLARE(int)
//yaml_document_add_sequence(document *yaml_document_t,
//        Tag *yaml_char_t, style yaml_sequence_style_t)
//{
//    struct {
//        error yaml_error_type_t
//    } context
//    mark yaml_mark_t = { 0, 0, 0 }
//    tag_copy *yaml_char_t = NULL
//    struct {
//        start *yaml_node_item_t
//        end *yaml_node_item_t
//        top *yaml_node_item_t
//    } items = { NULL, NULL, NULL }
//    Node yaml_node_t
//
//    assert(document) // Non-NULL document object is expected.
//
//    if (!Tag) {
//        Tag = (yaml_char_t *)YAML_DEFAULT_SEQUENCE_TAG
//    }
//
//    if (!yaml_check_utf8(Tag, strlen((char *)Tag))) goto error
//    tag_copy = yaml_strdup(Tag)
//    if (!tag_copy) goto error
//
//    if (!STACK_INIT(&context, items, INITIAL_STACK_SIZE)) goto error
//
//    SEQUENCE_NODE_INIT(Node, tag_copy, items.start, items.end,
//            style, mark, mark)
//    if (!PUSH(&context, document.nodes, Node)) goto error
//
//    return document.nodes.top - document.nodes.start
//
//error:
//    STACK_DEL(&context, items)
//    yaml_free(tag_copy)
//
//    return 0
//}
//
///*
// * Add a mapping Node to a document.
// */
//
//YAML_DECLARE(int)
//yaml_document_add_mapping(document *yaml_document_t,
//        Tag *yaml_char_t, style yaml_mapping_style_t)
//{
//    struct {
//        error yaml_error_type_t
//    } context
//    mark yaml_mark_t = { 0, 0, 0 }
//    tag_copy *yaml_char_t = NULL
//    struct {
//        start *yaml_node_pair_t
//        end *yaml_node_pair_t
//        top *yaml_node_pair_t
//    } pairs = { NULL, NULL, NULL }
//    Node yaml_node_t
//
//    assert(document) // Non-NULL document object is expected.
//
//    if (!Tag) {
//        Tag = (yaml_char_t *)YAML_DEFAULT_MAPPING_TAG
//    }
//
//    if (!yaml_check_utf8(Tag, strlen((char *)Tag))) goto error
//    tag_copy = yaml_strdup(Tag)
//    if (!tag_copy) goto error
//
//    if (!STACK_INIT(&context, pairs, INITIAL_STACK_SIZE)) goto error
//
//    MAPPING_NODE_INIT(Node, tag_copy, pairs.start, pairs.end,
//            style, mark, mark)
//    if (!PUSH(&context, document.nodes, Node)) goto error
//
//    return document.nodes.top - document.nodes.start
//
//error:
//    STACK_DEL(&context, pairs)
//    yaml_free(tag_copy)
//
//    return 0
//}
//
///*
// * Append an item to a sequence Node.
// */
//
//YAML_DECLARE(int)
//yaml_document_append_sequence_item(document *yaml_document_t,
//        sequence int, item int)
//{
//    struct {
//        error yaml_error_type_t
//    } context
//
//    assert(document) // Non-NULL document is required.
//    assert(sequence > 0
//            && document.nodes.start + sequence <= document.nodes.top)
//                            // Valid sequence id is required.
//    assert(document.nodes.start[sequence-1].type == YAML_SEQUENCE_NODE)
//                            // A sequence Node is required.
//    assert(item > 0 && document.nodes.start + item <= document.nodes.top)
//                            // Valid item id is required.
//
//    if (!PUSH(&context,
//                document.nodes.start[sequence-1].data.sequence.items, item))
//        return 0
//
//    return 1
//}
//
///*
// * Append a pair of a key and a Value to a mapping Node.
// */
//
//YAML_DECLARE(int)
//yaml_document_append_mapping_pair(document *yaml_document_t,
//        mapping int, key int, Value int)
//{
//    struct {
//        error yaml_error_type_t
//    } context
//
//    pair yaml_node_pair_t
//
//    assert(document) // Non-NULL document is required.
//    assert(mapping > 0
//            && document.nodes.start + mapping <= document.nodes.top)
//                            // Valid mapping id is required.
//    assert(document.nodes.start[mapping-1].type == YAML_MAPPING_NODE)
//                            // A mapping Node is required.
//    assert(key > 0 && document.nodes.start + key <= document.nodes.top)
//                            // Valid key id is required.
//    assert(Value > 0 && document.nodes.start + Value <= document.nodes.top)
//                            // Valid Value id is required.
//
//    pair.key = key
//    pair.Value = Value
//
//    if (!PUSH(&context,
//                document.nodes.start[mapping-1].data.mapping.pairs, pair))
//        return 0
//
//    return 1
//}
//
//
