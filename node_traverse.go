package sapcontrol

import "strings"

// GetNodeByID returns an AlertNode with specific ID.
func (n AlertNodes) GetNodeByID(id int) (res AlertNode) {
	for _, node := range n {
		if node.ID == id {
			res = node
			break
		}
	}

	return res
}

// GetChildNodeByName returns hierarchically below parentID the next node with given name.
func (n AlertNodes) GetChildNodeByName(name string, parentID int) (res AlertNode) {
	for _, node := range n {
		if strings.ToLower(node.Name) == strings.ToLower(name) && node.Parent == parentID {
			res = node
			break
		}
	}

	return res
}

// GetNodesByName returns all AlertNodes with given name.
func (n AlertNodes) GetNodesByName(name string) (res AlertNodes) {
	for _, node := range n {
		if strings.ToLower(node.Name) == strings.ToLower(name) {
			res = append(res, node)
		}
	}

	return res
}

// GetNodesByParentID returns all AlertNodes which belong to given parent.
func (n AlertNodes) GetNodesByParentID(parentID int) (res AlertNodes) {
	for _, node := range n {
		if node.Parent == parentID {
			res = append(res, node)
		}
	}

	return res
}

// GetNodesByNameRecursive returns all AlertNodes and theit childs which belong to parent node with given name.
func (n AlertNodes) GetNodesByNameRecursive(name string) (res AlertNodes) {
	nodes := n.GetNodesByName(name)

	for _, node := range nodes {
		subNodes := n.GetNodesByParentIDRecursive(node.ID)
		if subNodes != nil {
			res = append(res, subNodes...)
		}
	}

	return res
}

// GetNodesByParentIDRecursive returns all AlertNodes and their childs which blong to parent with given ID.
func (n AlertNodes) GetNodesByParentIDRecursive(parentID int) (res AlertNodes) {
	nodes := n.GetNodesByParentID(parentID)

	res = append(res, nodes...)
	// traverse through nodes and find more childs
	for _, node := range nodes {
		subNodes := n.GetNodesByParentIDRecursive(node.ID)
		if subNodes != nil {
			res = append(res, subNodes...)
		}
	}

	return res
}

// GetLastNodesByName returns hierarchically last AlertNodes which blong to parent node with given name.
func (n AlertNodes) GetLastNodesByName(name string) (res AlertNodes) {
	nodes := n.GetNodesByName(name)

	// traverse through nodes and find last childs
	for _, node := range nodes {
		// check if we have childs
		switch n.GetNodesByParentID(node.ID) {
		case nil:
			// we have no childs, so we are the last node
			res = append(res, node)
		default:
			// we have childs, get their last nodes
			resChilds := n.GetLastNodesByParentID(node.ID)
			if resChilds != nil {
				res = append(res, resChilds...)
			}
		}
	}

	return res
}

// GetLastNodesByParentID returns hierarchically last AlertNodes which blong to parent with given ID.
func (n AlertNodes) GetLastNodesByParentID(parentID int) (res AlertNodes) {
	nodes := n.GetNodesByParentID(parentID)

	// traverse through nodes and find last childs
	for _, node := range nodes {
		// check if we have childs
		switch n.GetNodesByParentID(node.ID) {
		case nil:
			// we have no childs, so we are the last node
			res = append(res, node)
		default:
			// we have childs, get their last nodes
			resChilds := n.GetLastNodesByParentID(node.ID)
			if resChilds != nil {
				res = append(res, resChilds...)
			}
		}
	}

	return res
}
