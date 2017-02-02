package bot

import (
	"reflect"
	"strings"
)

type Dispatcher interface {
	Event(name string, params ...interface{}) (bool, error)
	Handle(name string, handler Handler)
	RemoveHandler(name string, handler Handler)
}

type tNode struct {
	handlers []Handler
	children map[string]*tNode
}

func (node *tNode) Handle(name string, params []interface{}) (bool, error) {
	for _, handler := range node.handlers {
		stop, err := handler(name, params)
		if err != nil {
			return false, err
		}

		if stop {
			return true, nil
		}
	}
	return false, nil
}

type trieDispatcher struct {
	tree tNode
}

func (d *trieDispatcher) Event(name string, params ...interface{}) (bool, error) {
	parts := strings.Split(name, ".")

	node := &d.tree
	for i, part := range parts {
		if wc, ok := node.children["*"]; ok {
			stop, err := wc.Handle(name, params)
			if err != nil {
				return false, err
			}

			if stop {
				return true, nil
			}
		}

		if i == len(parts)-1 {
			if wc, ok := node.children["?"]; ok {
				stop, err := wc.Handle(name, params)
				if err != nil {
					return false, err
				}

				if stop {
					return true, nil
				}
			}
		}

		var ok bool
		node, ok = node.children[part]
		if !ok {
			return false, nil
		}

	}

	stop, err := node.Handle(name, params)
	if err != nil {
		return false, err
	}

	if stop {
		return true, nil
	}

	return false, nil
}

func (d *trieDispatcher) Handle(name string, handler Handler) {
	parts := strings.Split(name, ".")

	node := &d.tree
	for _, part := range parts {
		child, ok := node.children[part]
		if ok {
			node = child
			continue
		}
		child = &tNode{}

		if node.children == nil {
			node.children = make(map[string]*tNode)
		}

		node.children[part] = child
		node = child
	}

	node.handlers = append(node.handlers, handler)
}

func (d *trieDispatcher) RemoveHandler(name string, handler Handler) {
	parts := strings.Split(name, ".")
	parents := make([]*tNode, len(parts))

	node := &d.tree
	for i, part := range parts {
		parents[i] = node

		var ok bool
		node, ok = node.children[part]
		if !ok {
			return
		}
	}

	hp := reflect.ValueOf(handler).Pointer()
	for i, handler := range node.handlers {
		if reflect.ValueOf(handler).Pointer() == hp {
			node.handlers[i] = node.handlers[len(node.handlers)-1]
			node.handlers = node.handlers[:len(node.handlers)-1]
			break
		}
	}

	for i := len(parts) - 1; i >= 0; i -= 1 {
		if len(node.handlers) > 0 || len(node.children) > 0 {
			return
		}
		node = parents[i]

		delete(node.children, parts[i])
	}
}
