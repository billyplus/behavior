package behavior

import (
	"github.com/billyplus/behavior/config"

	"github.com/pkg/errors"
)

// NewBehaviorManager 加载behavior3格式的project配置
func NewBehaviorManager(cfg *config.BH3Project) (*BehaviorManager, error) {
	mgr := &BehaviorManager{
		treelist: make(map[string]*BehaviorTree, len(cfg.Trees)),
	}
	// nodeid := make(map[string]int)
	// roots := make(map[string]int, len(cfg.Trees))
	trees := make(map[string]config.BH3Tree, len(cfg.Trees))
	for _, tcfg := range cfg.Trees {
		trees[tcfg.ID] = tcfg
	}
	// index := 0
	for _, tcfg := range cfg.Trees {
		// generate index
		guidlist := make(map[string]int)
		nodelist := make([]config.BH3Node, 0)
		subtreelist := make([]string, 0)
		subtreelist = append(subtreelist, tcfg.ID)

		for subindex := 0; ; {
			if subindex >= len(subtreelist) {
				break
			}
			subtreeid := subtreelist[subindex]
			tcfg := trees[subtreeid]

			for guid, cfg := range tcfg.Nodes {
				if _, ok := guidlist[guid]; !ok {
					guidlist[guid] = 0
					nodelist = append(nodelist, cfg)
				}
			}
			for _, node := range tcfg.Nodes {
				if len(node.Children) > 0 {
					for _, childid := range node.Children {
						if _, ok := guidlist[childid]; !ok {
							// 是个subtree
							guidlist[childid] = 0
							subtreelist = append(subtreelist, childid)
						}
					}
				}
				if len(node.Child) > 0 {
					if _, ok := guidlist[node.Child]; !ok {
						// 是个subtree
						guidlist[node.Child] = 0
						subtreelist = append(subtreelist, node.Child)
					}
				}
			}
			subindex++
		}

		tree := &BehaviorTree{}
		mgr.treelist[tcfg.Title] = tree
		wraplist := make([]Wrapper, len(nodelist))

		for i, node := range nodelist {

			// fatory function for node
			treenode, err := createNodeByName(node.Name)
			if err != nil {
				return nil, err
			}
			nw := &wraplist[i]
			nw.Node = treenode
			nw.name = node.Name
			nw.index = i

			// set root node
			if node.ID == tcfg.Root {
				tree.root = nw
			}
			// add child
			if len(node.Child) > 0 {
				found := false
				for childindex, child := range nodelist {
					if child.ID == node.Child {
						childwr := &wraplist[childindex]
						nw.AddChild(childwr)
						found = true
					}
				}
				if !found {
					// sub tree
					subt := trees[node.Child]
					for childindex, child := range nodelist {
						if child.ID == subt.Root {
							childwr := &wraplist[childindex]
							nw.AddChild(childwr)
							found = true
						}
					}
				}
				if !found {
					return nil, errors.Errorf("unknow child guid '%s'", node.Child)

				}
			}
			// add children
			if len(node.Children) > 0 {
				for _, childguid := range node.Children {
					found := false
					for childindex, child := range nodelist {
						if child.ID == childguid {
							childwr := &wraplist[childindex]
							nw.AddChild(childwr)
							found = true
						}
					}
					if !found {
						// sub tree
						subt := trees[childguid]
						for childindex, child := range nodelist {
							if child.ID == subt.Root {
								childwr := &wraplist[childindex]
								nw.AddChild(childwr)
								found = true
							}
						}
					}
					if !found {
						return nil, errors.Errorf("unknow child guid '%s'", childguid)

					}
				}
			}

			// Initialize treenode when finish adding children
			if err = treenode.Initialize(&node); err != nil {
				return nil, err
			}

		}

		tree.nodelist = wraplist
	}

	return mgr, nil
}

// NewBehaviorManager 加载behavior3格式的project配置
// func NewBehaviorManager(cfg *config.BH3Project) (*BehaviorManager, error) {
// 	mgr := &BehaviorManager{
// 		treelist: make(map[string]*BehaviorTree, len(cfg.Trees)),
// 	}
// 	nodeid := make(map[string]int)
// 	roots := make(map[string]int, len(cfg.Trees))
// 	index := 0
// 	for _, tcfg := range cfg.Trees {
// 		// generate index
// 		for guid := range tcfg.Nodes {
// 			if _, ok := nodeid[guid]; !ok {
// 				nodeid[guid] = index
// 				index++
// 			}
// 		}
// 		roots[tcfg.ID] = nodeid[tcfg.Root]
// 	}

// 	// generate node
// 	mgr.nodelist = make([]Wrapper, index)

// 	for _, tcfg := range cfg.Trees {
// 		// log.Logger.Debug().Msg("NewBehaviorManager new tree")
// 		tree := BehaviorTree{}
// 		// guid
// 		tree.id = tcfg.ID
// 		tree.title = tcfg.Title
// 		// root node
// 		rootid := roots[tcfg.ID]
// 		tree.root = &mgr.nodelist[rootid]

// 		mgr.treelist[tree.title] = &tree

// 		for guid, c := range tcfg.Nodes {
// 			id := nodeid[guid]
// 			name := c.Name
// 			// fatory function for node
// 			treenode, err := createNodeByName(name)
// 			if err != nil {
// 				return nil, err
// 			}
// 			// log.Logger.Debug().Msgf("NewBehaviorManager new node %+v", treenode)

// 			nw := &mgr.nodelist[id]
// 			nw.Node = treenode
// 			nw.index = id
// 			// log.Logger.Debug().Msgf("NewBehaviorManager nw is %+v", nw)
// 			// add child
// 			if len(c.Child) > 0 {
// 				if childid, ok := nodeid[c.Child]; ok {
// 					// log.Logger.Debug().Msgf("NewBehaviorManager child index is %d", childid)
// 					childnode := &mgr.nodelist[childid]
// 					nw.AddChild(childnode)
// 				} else if childid, ok := roots[c.Child]; ok {
// 					// sub tree
// 					childnode := &mgr.nodelist[childid]
// 					nw.AddChild(childnode)
// 				} else {
// 					return nil, errors.Errorf("unknow child guid '%s'", c.Child)
// 				}
// 			}
// 			// add children
// 			if len(c.Children) > 0 {
// 				for _, child := range c.Children {
// 					if childid, ok := nodeid[child]; ok {
// 						// log.Logger.Debug().Msgf("NewBehaviorManager child index is %d", childid)
// 						childnode := &mgr.nodelist[childid]
// 						nw.AddChild(childnode)
// 					} else {
// 						return nil, errors.Errorf("unknow child guid '%s'", child)
// 					}
// 				}
// 			}

// 			// Initialize treenode when finish adding children
// 			if err = treenode.Initialize(&c); err != nil {
// 				return nil, err
// 			}
// 		}
// 	}

// 	// log.Logger.Debug().Msgf("tree list %+v", mgr.treelist["A behavior tree"].root.Node)
// 	return mgr, nil
// }

type BehaviorManager struct {
	treelist map[string]*BehaviorTree
	// nodelist []Wrapper
}

func (mgr *BehaviorManager) SelectBehaviorTree(name string) *BehaviorTree {
	return mgr.treelist[name]
}

// func (mgr *BehaviorManager) NodeCount() int {
// 	return len(mgr.nodelist)
// }

// func (mgr *BehaviorManager) NewBlackboard() *Blackboard {
// 	bb := newBlackboard(mgr.NodeCount())
// 	for i := 0; i < mgr.NodeCount(); i++ {
// 		bb.nodeMemo[i] = mgr.nodelist[i].Node.CreateMemo()
// 	}
// 	return bb
// }
