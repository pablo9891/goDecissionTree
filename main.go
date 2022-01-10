package main

import (
	"decissionTree/analysis"
	"decissionTree/datahandler"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func showNode(n *analysis.Node) {
	fmt.Println(strings.Repeat(" ", n.Level), "+ ", n.Val)
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("error in parameters")
		return
	}

	if filepath.Ext(os.Args[1]) != ".csv" {
		fmt.Println("Not valid file extension. Should be csv")
		return
	}

	class := os.Args[2]
	df, err := datahandler.Parse(os.Args[1])

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	root, err := analysis.ID3(df, class)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	analysis.Dfs(root, showNode)
}
