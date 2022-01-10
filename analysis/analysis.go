package analysis

import (
	"decissionTree/datahandler"
	"decissionTree/queue"
	"decissionTree/stack"
	"fmt"
	"math"
)

const (
	FieldNode      = iota // Contains a field
	FieldValueNode        // Contains values for a field
	FinishNode            // Finish node
)

type Node struct {
	Children     []*Node
	NodeType     int
	Val          interface{}
	DataType     int
	ColumnHeader string
	Data         *datahandler.Data
	Level        int
}

func GetFrequency(c *datahandler.Column) map[string]int {
	freq := make(map[string]int)

	for _, v := range c.Fields {
		if _, ok := freq[fmt.Sprint(v.Val)]; ok {
			freq[fmt.Sprint(v.Val)]++
		} else {
			freq[fmt.Sprint(v.Val)] = 1
		}
	}
	return freq
}

func GetFrequencyIntersection(c *datahandler.Column, c2 *datahandler.Column) map[string]map[string]int {
	freq := make(map[string]map[string]int)

	i := 0
	for _, v := range c.Fields {
		if _, ok := freq[fmt.Sprint(v.Val)]; ok {
			if _, k := freq[fmt.Sprint(v.Val)][fmt.Sprint(c2.Fields[i].Val)]; k {
				freq[fmt.Sprint(v.Val)][fmt.Sprint(c2.Fields[i].Val)]++
			} else {
				freq[fmt.Sprint(v.Val)][fmt.Sprint(c2.Fields[i].Val)] = 1
			}
		} else {
			freq[fmt.Sprint(v.Val)] = make(map[string]int)
			freq[fmt.Sprint(v.Val)][fmt.Sprint(c2.Fields[i].Val)] = 1
		}
		i++
	}
	return freq
}

func GetEntropy(freq map[string]int, nrows int) float64 {
	acc := 0
	prob := 0.0
	for _, v := range freq {
		acc = acc + v
	}

	for _, v := range freq {
		p := float64(v) / float64(acc)
		prob = prob + (-1)*p*math.Log2(p)
	}
	return prob * (float64(acc) / float64(nrows))
}

func GetInformationGain(entropyMap map[string]float64, labelEntropy float64) map[string]float64 {
	ig := make(map[string]float64)

	for k, v := range entropyMap {
		ig[k] = labelEntropy - v
	}
	return ig
}

func GetMaxInformationGain(igMap map[string]float64, class string) (string, float64) {
	var maxNumber float64
	var key string
	var k string

	filterMap := make(map[string]float64)
	for k, v := range igMap {
		if k != class {
			filterMap[k] = v
		}
	}

	for k, maxNumber = range filterMap {
		key = k
		break
	}
	for k, n := range filterMap {
		if n > maxNumber {
			key = k
			maxNumber = n
		}
	}
	return key, maxNumber
}

func GetMaxInformationGainAttribute(data *datahandler.Data, class string) *Node {
	entropyMap := make(map[string]float64)
	for _, c := range data.Columns {
		fr := GetFrequencyIntersection(c, data.Columns[class])
		for k := range fr {
			entropyMap[c.Name] += GetEntropy(fr[k], data.NumberOfRows)
		}
	}

	labelEntropy := GetEntropy(GetFrequency(data.Columns[class]), data.NumberOfRows)
	igMap := GetInformationGain(entropyMap, labelEntropy)
	k, _ := GetMaxInformationGain(igMap, class)

	return &Node{nil, FieldNode, k, data.Columns[k].DataType, k, nil, -1}
}

func ID3(data *datahandler.Data, class string) (*Node, error) {
	if _, ok := data.Columns[class]; !ok {
		return nil, fmt.Errorf("no class column was found with name %s", class)
	}

	q := queue.Create()
	root := GetMaxInformationGainAttribute(data, class)
	queue.Add(q, root)

	root.Data = data
	root.Level = 0
	for !queue.IsEmpty(q) {
		curr := queue.Get(q).(*Node)
		switch curr.NodeType {
		case FieldNode:
			values := GetFrequency(curr.Data.Columns[curr.ColumnHeader])
			for v := range values {
				split := datahandler.SplitData(curr.Data, curr.Data.Columns[curr.ColumnHeader], v)
				newNode := Node{nil, FieldValueNode, v, curr.DataType, curr.Data.Columns[curr.ColumnHeader].Name, split, curr.Level + 1}
				curr.Children = append(curr.Children, &newNode)
				queue.Add(q, &newNode)
			}
		case FieldValueNode:
			classEntropy := GetEntropy(GetFrequency(curr.Data.Columns[class]), curr.Data.NumberOfRows)
			if classEntropy == 0 {
				newNode := &Node{nil, FinishNode, curr.Data.Columns[class].Fields[0].Val, curr.Data.Columns[class].DataType, curr.Data.Columns[class].Name, curr.Data, curr.Level + 1}
				curr.Children = append(curr.Children, newNode)
			} else {
				if curr.Data.NumberOfCol != 0 {
					newAttrNode := GetMaxInformationGainAttribute(curr.Data, class)
					newAttrNode.Data = curr.Data
					newAttrNode.Level = curr.Level + 1
					curr.Children = append(curr.Children, newAttrNode)
					queue.Add(q, newAttrNode)
				}
			}
		}
	}
	return root, nil
}

func Bfs(root *Node, f func(*Node)) {
	q := queue.Create()
	queue.Add(q, root)
	for !queue.IsEmpty(q) {
		curr := queue.Get(q).(*Node)
		f(curr)
		for _, v := range curr.Children {
			queue.Add(q, v)
		}
	}
}

func Dfs(root *Node, f func(*Node)) {
	s := stack.Create()
	stack.Push(s, root)
	for !stack.IsEmpty(s) {
		curr := stack.Pop(s).(*Node)
		f(curr)
		for _, v := range curr.Children {
			stack.Push(s, v)
		}
	}
}
