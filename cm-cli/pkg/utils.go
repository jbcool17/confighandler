package pkg

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

func ensureRootDir(root string) error {
	if root == "" {
		return errors.New("root path empty")
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return fmt.Errorf("failed to create root dir: %w", err)
	}
	return nil
}

func WriteYAMLFile(path string, node interface{}) error {
	yamlBytes, err := yaml.Marshal(node)
	if err != nil {
		return err
	}
	content := "---\n" + string(yamlBytes)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}

	fmt.Printf("Successfully wrote YAML data to %s\n", path)
	fmt.Println("--- Generated YAML ---")
	fmt.Println(string(yamlBytes))
	return nil
}

func parseKeyValuePairs(s string) map[string]string {
	out := map[string]string{}
	s = strings.TrimSpace(s)
	if s == "" {
		return out
	}
	parts := strings.Split(s, ",")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		out[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	return out
}

// parseValue converts a string value to the appropriate YAML node with correct type
func parseValue(s string) *yaml.Node {
	s = strings.TrimSpace(s)
	// check for boolean
	if s == "true" || s == "false" {
		return &yaml.Node{Kind: yaml.ScalarNode, Value: s, Tag: "!!bool"}
	}
	// check for number (int or float)
	if _, err := strconv.Atoi(s); err == nil {
		return &yaml.Node{Kind: yaml.ScalarNode, Value: s, Tag: "!!int"}
	}
	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return &yaml.Node{Kind: yaml.ScalarNode, Value: s, Tag: "!!float"}
	}
	// default to string
	return &yaml.Node{Kind: yaml.ScalarNode, Value: s, Tag: "!!str"}
}

// setYAMLPath finds or creates mapping nodes along path and sets the final scalar value.
func setYAMLPath(node *yaml.Node, path []string, value string) error {
	if len(path) == 0 {
		return errors.New("empty path")
	}
	// Ensure we're at a mapping node
	if node.Kind != yaml.MappingNode {
		return fmt.Errorf("node is not a mapping node")
	}
	// handle segment which may include an index like key[0]
	seg := path[0]
	// parse index if present
	if strings.Contains(seg, "[") && strings.HasSuffix(seg, "]") {
		// indexed sequence under a mapping key
		idxStart := strings.Index(seg, "[")
		key := seg[:idxStart]
		idxStr := seg[idxStart+1 : len(seg)-1]
		idx, err := strconv.Atoi(idxStr)
		if err != nil {
			return fmt.Errorf("invalid index %s: %w", idxStr, err)
		}

		// find or create the mapping key
		for i := 0; i < len(node.Content); i += 2 {
			k := node.Content[i]
			v := node.Content[i+1]
			if k.Value == key {
				// ensure v is a sequence node
				if v.Kind != yaml.SequenceNode {
					v.Kind = yaml.SequenceNode
					v.Content = []*yaml.Node{}
				}
				// extend sequence if necessary
				for len(v.Content) <= idx {
					// append an empty mapping node by default
					v.Content = append(v.Content, &yaml.Node{Kind: yaml.MappingNode})
				}
				elem := v.Content[idx]
				if len(path) == 1 {
					// set the element directly
					v.Content[idx] = parseValue(value)
					return nil
				}
				// descend into element
				if elem.Kind != yaml.MappingNode {
					elem.Kind = yaml.MappingNode
					elem.Content = []*yaml.Node{}
				}
				return setYAMLPath(elem, path[1:], value)
			}
		}

		// key not found: create sequence with enough elements
		kNode := &yaml.Node{Kind: yaml.ScalarNode, Value: key, Tag: "!!str"}
		seq := &yaml.Node{Kind: yaml.SequenceNode}
		for i := 0; i <= idx; i++ {
			seq.Content = append(seq.Content, &yaml.Node{Kind: yaml.MappingNode})
		}
		node.Content = append(node.Content, kNode, seq)
		// now descend into seq.Content[idx]
		if len(path) == 1 {
			seq.Content[idx] = parseValue(value)
			return nil
		}
		return setYAMLPath(seq.Content[idx], path[1:], value)
	}

	// non-indexed mapping key handling
	key := seg
	for i := 0; i < len(node.Content); i += 2 {
		k := node.Content[i]
		v := node.Content[i+1]
		if k.Value == key {
			if len(path) == 1 {
				*v = *parseValue(value)
				return nil
			}
			// descend: if this is a sequence, apply to each element
			if v.Kind == yaml.SequenceNode {
				if len(v.Content) == 0 {
					// create one element to modify
					v.Content = append(v.Content, &yaml.Node{Kind: yaml.MappingNode})
				}
				for _, elem := range v.Content {
					if elem.Kind != yaml.MappingNode {
						elem.Kind = yaml.MappingNode
						elem.Content = []*yaml.Node{}
					}
					if err := setYAMLPath(elem, path[1:], value); err != nil {
						return err
					}
				}
				return nil
			}
			if v.Kind != yaml.MappingNode {
				v.Kind = yaml.MappingNode
				v.Content = []*yaml.Node{}
			}
			return setYAMLPath(v, path[1:], value)
		}
	}

	// key not found: create it
	kNode := &yaml.Node{Kind: yaml.ScalarNode, Value: key, Tag: "!!str"}
	var vNode *yaml.Node
	if len(path) == 1 {
		vNode = parseValue(value)
	} else {
		vNode = &yaml.Node{Kind: yaml.MappingNode}
		if err := setYAMLPath(vNode, path[1:], value); err != nil {
			return err
		}
	}
	node.Content = append(node.Content, kNode, vNode)
	return nil
}

// DEBUG - YAML

func DumpYAMLNode(n *yaml.Node) {
	fmt.Println("YAML Node Dump:")
	dumpNode(n, 0, "")
}

func dumpNode(n *yaml.Node, indent int, label string) {
	if n == nil {
		return
	}

	prefix := strings.Repeat("  ", indent)

	kind := map[yaml.Kind]string{
		yaml.DocumentNode: "Document",
		yaml.MappingNode:  "Mapping",
		yaml.SequenceNode: "Sequence",
		yaml.ScalarNode:   "Scalar",
		yaml.AliasNode:    "Alias",
	}[n.Kind]

	if label != "" {
		label = label + ": "
	}

	fmt.Printf(
		"%s%s[%s] Tag=%q Value=%q Anchor=%q (line=%d col=%d)\n",
		prefix,
		label,
		kind,
		n.Tag,
		n.Value,
		n.Anchor,
		n.Line,
		n.Column,
	)

	switch n.Kind {

	case yaml.DocumentNode:
		for _, c := range n.Content {
			dumpNode(c, indent+1, "document")
		}

	case yaml.MappingNode:
		for i := 0; i < len(n.Content); i += 2 {
			key := n.Content[i]
			val := n.Content[i+1]

			fmt.Printf("%s  Key:\n", prefix)
			dumpNode(key, indent+2, "key")

			fmt.Printf("%s  Value:\n", prefix)
			dumpNode(val, indent+2, "value")
		}

	case yaml.SequenceNode:
		for i, c := range n.Content {
			dumpNode(c, indent+1, fmt.Sprintf("item[%d]", i))
		}

	case yaml.ScalarNode:
		// nothing more to recurse into

	case yaml.AliasNode:
		dumpNode(n.Alias, indent+1, "alias")
	}
}
