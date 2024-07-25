package json_convert

import (
	"fmt"
	"testing"
)

func TestElementRecursor_ExpandAppendKeys(t *testing.T) {
	//找到第一个[*]出现的位置，并进行展开。展开的关键是根据既有的数据展开，要与既有的数据match上
	//a.b.[*].c.[*].d →
	//		a.b.[0].c.[*].d →
	//				a.b.[0].c.[0].d
	//				a.b.[0].c.[1].d
	//		a.b.[1].c.[*].d、→
	//				a.b.[1].c.[0].d
	//		a.b.[2].c.[*].d→
	//				a.b.[2].c.[0].d
	//				a.b.[2].c.[1].d
	//				a.b.[2].c.[2].d
	elemMapper := new(ElementMapper)
	elemMapper.Init(nil, nil)
	elemMapper.saveMappedData([]string{"a.b.[0].c.[0].d"}, "1")
	elemMapper.saveMappedData([]string{"a.b.[0].c.[1].d"}, "1")
	elemMapper.saveMappedData([]string{"a.b.[1].c.[0].d"}, "1")
	elemMapper.saveMappedData([]string{"a.b.[2].c.[0].d"}, "1")
	elemMapper.saveMappedData([]string{"a.b.[2].c.[1].d"}, "1")
	elemMapper.saveMappedData([]string{"a.b.[2].c.[2].d"}, "1")
	elemMapper.saveMappedData([]string{"a.b.[2].c.[3].d"}, "1")
	elemMapper.saveMappedData([]string{"a.b.[0].c.[0].d.e.[0].f"}, "1")

	tobeExpandKey := "a.b.[*]"
	printExpandKey(tobeExpandKey, elemMapper.expandAppendKey(tobeExpandKey))
	tobeExpandKey = "a.b.[*].x"
	printExpandKey(tobeExpandKey, elemMapper.expandAppendKey(tobeExpandKey))
	tobeExpandKey = "a.b.[*]."
	printExpandKey(tobeExpandKey, elemMapper.expandAppendKey(tobeExpandKey))
	tobeExpandKey = "a.b.[*].c"
	printExpandKey(tobeExpandKey, elemMapper.expandAppendKey(tobeExpandKey))
	tobeExpandKey = "a.b.[*].c.[*]"
	printExpandKey(tobeExpandKey, elemMapper.expandAppendKey(tobeExpandKey))
	tobeExpandKey = "a.b.[*].c.[*].d"
	printExpandKey(tobeExpandKey, elemMapper.expandAppendKey(tobeExpandKey))
	tobeExpandKey = "a.b.[*].c.[*].d.[*].f"
	printExpandKey(tobeExpandKey, elemMapper.expandAppendKey(tobeExpandKey))
}

func printExpandKey(tobeExpandKey string, expandedKeys []string) {
	fmt.Println(tobeExpandKey, ":")
	for _, expandedKey := range expandedKeys {
		fmt.Println("\t\t", expandedKey)
	}
}
