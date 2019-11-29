package bingosvc

import (
	"github.com/shihtzu-systems/bingo/pkg/bingo"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func NewBingo(boxes bingo.Boxes) bingo.Board {
	if len(boxes) < 25 {
		log.Fatal("not enough boxes to build a bingo board")
	}

	// shuffle
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(boxes), func(i, j int) { boxes[i], boxes[j] = boxes[j], boxes[i] })

	return bingo.Board{
		B: boxes[:5],
		I: boxes[5:10],
		N: boxes[10:15],
		G: boxes[15:20],
		O: boxes[20:25],
	}
}
