package bingo

import (
	"bytes"
	"encoding/json"
	"log"
	"strings"
)

type Board struct {
	Id string `json:"id"`

	B Boxes `json:"b"`
	I Boxes `json:"i"`
	N Boxes `json:"n"`
	G Boxes `json:"g"`
	O Boxes `json:"o"`
}

func (b Board) Print() {
	b.B.Print("B")
	b.I.Print("I")
	b.N.Print("N")
	b.G.Print("G")
	b.O.Print("O")
}

func (b *Board) Mark(letter string, index int) {
	switch strings.ToLower(letter) {
	case "b":
		b.B[index].Marked = !b.B[index].Marked
	case "i":
		b.I[index].Marked = !b.I[index].Marked
	case "n":
		b.N[index].Marked = !b.N[index].Marked
	case "g":
		b.G[index].Marked = !b.G[index].Marked
	case "o":
		b.O[index].Marked = !b.O[index].Marked
	default:
		log.Fatal("unknown letter: ", letter)
	}
}

func (b Board) PrettyJson() []byte {
	jout, err := json.Marshal(b)
	if err != nil {
		log.Fatal(err)
	}
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, jout, "", "  "); err != nil {
		log.Fatal(err)
	}
	return pretty.Bytes()
}

type Boxes []Box

func (b Boxes) Print(letter string) {
	for i, box := range b {
		log.Printf("%s-%d [ %s ] %v", strings.ToUpper(letter), i, box.Content, box.Marked)
	}
}

type Box struct {
	Content string `json:"content"`
	Marked  bool   `json:"marked"`
}
