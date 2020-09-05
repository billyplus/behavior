package bench

import (
	"fmt"
	"testing"

	"github.com/billyplus/behavior"
	"github.com/billyplus/behavior/action"
	"github.com/billyplus/behavior/config"
	_ "github.com/billyplus/behavior/loader"

	"github.com/stretchr/testify/suite"
)

func BenchmarkAction(b *testing.B) {
	behavior.Register("SuccessAfterTimes", &action.SuccessAfterTimes{})
	bh3prj, err := config.ParseB3File(proj)
	if err != nil {
		fmt.Println(err)
		return
	}

	mgr, err := behavior.NewBehaviorManager(&bh3prj.Data)
	if err != nil {
		fmt.Println(err)
		return
	}

	tree := mgr.SelectBehaviorTree("A behavior tree")
	if tree == nil {
		fmt.Println("tree not exist")
		return
	}

	// fmt.Printf("tree is %v\n", tree)

	bb := tree.NewBlackboard()

	round := 0
	status := behavior.StatusRunning
	for i := 0; i < b.N; i++ {

		for {
			status = tree.Tick(bb)
			// fmt.Println("round ", round)
			round++
			if status == behavior.StatusSuccess {
				break
			} else if status == behavior.StatusFailed {
				b.Fatal()
			}
			if round > 100 {
				b.Fatal()
			}
		}
		round = 0
	}

}

type projTestSuite struct {
	suite.Suite
}

func (suite *projTestSuite) SetupSuite() {

}

func (suite *projTestSuite) TestParseProj() {
	behavior.Register("SuccessAfterTimes", &action.SuccessAfterTimes{})
	bh3prj, err := config.ParseB3File(proj)
	suite.Nil(err, "parse ph3 project from []byte")

	mgr, err := behavior.NewBehaviorManager(&bh3prj.Data)
	suite.Nil(err, "create manager from project")

	tree := mgr.SelectBehaviorTree("A behavior tree")
	suite.NotNil(tree, "select tree")

}

// func TestBehavior3Proj(t *testing.T) {
// 	suite.Run(t, &projTestSuite{})
// }

var proj = []byte(`
{
	"name": "y3pp",
	"description": "",
	"data": {
		"version": "0.3.0",
		"scope": "project",
		"selectedTree": "7f321519-38f8-4520-a549-ee74d50cc7fe",
		"trees": [
		{
			"version": "0.3.0",
			"scope": "tree",
			"id": "7f321519-38f8-4520-a549-ee74d50cc7fe",
			"title": "A behavior tree",
			"description": "",
			"root": "8a5c9c13-1036-4c76-b5bc-8948b6ca1b5f",
			"properties": {},
			"nodes": {
			"8a5c9c13-1036-4c76-b5bc-8948b6ca1b5f": {
				"id": "8a5c9c13-1036-4c76-b5bc-8948b6ca1b5f",
				"name": "Selector",
				"category": "composite",
				"title": "Selector",
				"description": "顺序执行子节点，如果有一个子节点返回success或running，那么父节点马上返回success或running",
				"properties": {},
				"display": {
				"x": -348,
				"y": -96
				},
				"children": [
				"a9b4a15d-b5b7-44f2-845c-ddb87678b25b"
				]
			},
			"a9b4a15d-b5b7-44f2-845c-ddb87678b25b": {
				"id": "a9b4a15d-b5b7-44f2-845c-ddb87678b25b",
				"name": "SuccessAfterTimes",
				"category": "action",
				"title": "SuccessAfterTimes",
				"description": "test action",
				"properties": {
				"times": 5
				},
				"display": {
				"x": -48,
				"y": -96
				}
			}
			},
			"display": {
			"camera_x": 960,
			"camera_y": 403,
			"camera_z": 1,
			"x": -576,
			"y": -108
			}
		}
		],
		"custom_nodes": [
		{
			"version": "0.3.0",
			"scope": "node",
			"name": "Selector",
			"category": "composite",
			"title": "Selector",
			"description": "顺序执行子节点，如果有一个子节点返回success或running，那么父节点马上返回success或running",
			"properties": {}
		},
		{
			"version": "0.3.0",
			"scope": "node",
			"name": "Sequenece",
			"category": "composite",
			"title": "Sequenece",
			"description": "顺序执行子节点，只要有一个子节点返回failed或running，当前节点也返回failed或running",
			"properties": {}
		},
		{
			"version": "0.3.0",
			"scope": "node",
			"name": "Repeater",
			"category": "decorator",
			"title": "Repeater",
			"description": "重复执行子节点，在达到maxLoop之前，一直返回running。当达到maxLoop后，返回success。maxLoop=0表示无限重复",
			"properties": {
			"maxLoop": 0
			}
		},
		{
			"version": "0.3.0",
			"scope": "node",
			"name": "Succeeder",
			"category": "decorator",
			"title": "Succeeder",
			"description": "不论子节点返回结果是什么，都返回success",
			"properties": {}
		},
		{
			"version": "0.3.0",
			"scope": "node",
			"name": "UntilFailure",
			"category": "decorator",
			"title": "UntilFailure",
			"description": "重复执行子节点，直到子节点返回failed，此节点会返回success。或者达到maxLoop次后，此节点返回failed",
			"properties": {
			"maxLoop": 0
			}
		},
		{
			"version": "0.3.0",
			"scope": "node",
			"name": "UntilSuccess",
			"category": "decorator",
			"title": "UntilSuccess",
			"description": "重复执行子节点，直到子节点返回success，此节点会返回success。或者达到maxLoop次后，此节点返回failed",
			"properties": {
			"maxLoop": 0
			}
		},
		{
			"version": "0.3.0",
			"scope": "node",
			"name": "Wait",
			"category": "action",
			"title": "Wait@ms",
			"description": "等待一待时间",
			"properties": {
			"duration": 500
			}
		},
		{
			"version": "0.3.0",
			"scope": "node",
			"name": "SuccessAfterTimes",
			"category": "action",
			"title": "SuccessAfterTimes",
			"description": "test action",
			"properties": {
			"times": 5
			}
		}
		]
	}
  }
`)
