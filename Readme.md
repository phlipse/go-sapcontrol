# SAPcontrol
SAPcontrol provides access to the output of CCMS function calls on a SAP system.

## Usage
First call the **Read(io.Reader)** function to read in the output of a sapcontrol function call. Then you can work on that returned object. Use the objects **GetProcessList()** or **GetAlertNodes()** method to extract the content to useful structs.

## Construction of AlertNode methods
*GetAlertTree* returns multiple CSV records representing the state of a SAP system in CCMS. There is no fixed data structure behind it. Each AlertNode references its parent, that's all. You have to follow the references in reverse order to get the data you want. There are multiple methods provided by the *sapcontrol* package which support you by doing this - look at *node_arrange/traverse.go* in this repository to get the entire list.

### Traversing through AlertNodes
All methods which support you by traversing through AlertNodes are constructed as follows:

**"Get" + hierarchy + set + "By" + attribute + recursion**

* hierarchie: where to search for nodes, child means exactly one hierarchy below, last means the last node in hierarchy
* set: how many nodes get returned
* attribute: by which attribute the search performs
* recursion: recursive search for child nodes

e.g.:
* GetChildNodeByName
* GetNodesByParentID
* GetNodesByNameRecursive
* GetLastNodesByName

### Arranging AlertNodes
All methods which support you by rearranging the AlertNodes are listet below:

* GetNodePath: build up Path from given node to farest parent node
* NodePathToName: convert node hierarchy to delimited names of all nodes, could be used with GetNodePath

## Examples
For examples look at *example_test.go* file in this repository.

## License
[Apache License 2.0](https://github.com/phlipse/go-sapcontrol/blob/master/LICENSE)
