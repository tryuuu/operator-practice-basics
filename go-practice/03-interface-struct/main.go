package main

import "fmt"

// 共通のメソッド (GetName) を持つ構造体 (Parent と Child) を一貫してObjectとして操作できるようにするインターフェース
type Object interface {
	GetName() string
}

type Parent struct {
	Id   int
	Name string
}

// p Parentの部分はreceiver
func (p Parent) GetName() string {
	return p.Name
}

// pointer receiver(fieldを変更する場合)
func (p *Parent) SetName(name string) {
	p.Name = name
}

// Parentの持つフィールドをそのまま使うことができる(拡張、embedding)
type Child struct {
	Parent
	OtherField string
}

func printObjectName(obj Object) {
	//interfaceを使うことで、ParentとChildを同じ関数で使用できる
	fmt.Println(obj.GetName())
}

func main() {
	p := Parent{Id: 1, Name: "tryu"}
	fmt.Println(p)
	printObjectName(&p)
	ch := Child{Parent: p}
	printObjectName(&ch)
	ch = Child{Parent{Id: 2, Name: "child"}, "Other"}
	printObjectName(&ch)
}

/*
Object インターフェースが無かったら、型ごとに違う関数を書く必要がある:
func printParentName(p Parent) {
    fmt.Println(p.GetName())
}

func printChildName(c Child) {
    fmt.Println(c.GetName())
}

*/
