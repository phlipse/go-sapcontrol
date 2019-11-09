package sapcontrol

// GetNodePath returns hierarchically path from root to given node.
func (n AlertNodes) GetNodePath(node AlertNode) (nodes AlertNodes) {
	// check if we have a parent
	switch n.GetNodeByID(node.Parent) {
	case AlertNode{}:
		// we have no parent, so we are the root
		nodes = append(nodes, node)
	default:
		// we have a parent, get his parents names
		nodes = append(nodes, n.GetNodePath(n.GetNodeByID(node.Parent))...)
		// append current node
		nodes = append(nodes, node)
	}

	return
}

// NodePathToName returns name from root to last child node as string, delimited by given string.
func (n AlertNodes) NodePathToName(delim string) (name string) {
	for i := 0; i < len(n); i++ {
		if i == len(n)-1 {
			name += n[i].Name
		} else {
			name += n[i].Name + delim
		}
	}

	return
}
