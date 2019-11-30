package bingosvc

import "github.com/shihtzu-systems/bingo/pkg/bingo"

var bingos = []bingo.Bingo{
	// vertical
	{
		Type: "vertical",
		Id:   "b",

		B: []bool{
			// B-0 - X
			true,
			// B-1 - X
			true,
			// B-2 - X
			true,
			// B-3 - X
			true,
			// B-4 - X
			true,
		},
	},
	{
		Type: "vertical",
		Id:   "i",

		I: []bool{
			// I-0 - X
			true,
			// I-1 - X
			true,
			// I-2 - X
			true,
			// I-3 - X
			true,
			// I-4 - X
			true,
		},
	},
	{
		Type: "vertical",
		Id:   "n",

		N: []bool{
			// N-0 - X
			true,
			// N-1 - X
			true,
			// N-2 - X
			true,
			// N-3 - X
			true,
			// N-4 - X
			true,
		},
	},
	{
		Type: "vertical",
		Id:   "g",

		G: []bool{
			// G-0 - X
			true,
			// G-1 - X
			true,
			// G-2 - X
			true,
			// G-3 - X
			true,
			// G-4 - X
			true,
		},
	},
	{
		Type: "vertical",
		Id:   "o",
		O: []bool{
			// O-0 - X
			true,
			// O-1 - X
			true,
			// O-2 - X
			true,
			// O-3 - X
			true,
			// O-4 - X
			true,
		},
	},
	// horizontal
	{
		Type: "horizontal",
		Id:   "0",

		B: []bool{
			// B-0 - X
			true,
		},
		I: []bool{
			// I-0 - X
			true,
		},
		N: []bool{
			// N-0 - X
			true,
		},
		G: []bool{
			// G-0 - X
			true,
		},
		O: []bool{
			// O-0 - X
			true,
		},
	},
	{
		Type: "horizontal",
		Id:   "1",

		B: []bool{
			// B-0
			false,
			// B-1 - X
			true,
		},
		I: []bool{
			// I-0
			false,
			// I-1 - X
			true,
		},
		N: []bool{
			// N-0
			false,
			// N-1 - X
			true,
		},
		G: []bool{
			// G-0
			false,
			// G-1 - X
			true,
		},
		O: []bool{
			// O-0
			false,
			// O-1 - X
			true,
		},
	},
	{
		Type: "horizontal",
		Id:   "2",

		B: []bool{
			// B-0
			false,
			// B-1
			false,
			// B-2 - X
			true,
		},
		I: []bool{
			// I-0
			false,
			// I-1
			false,
			// I-2 - X
			true,
		},
		N: []bool{
			// N-0
			false,
			// N-1
			false,
			// N-2 - X
			true,
		},
		G: []bool{
			// G-0
			false,
			// G-1
			false,
			// G-2 - X
			true,
		},
		O: []bool{
			// O-0
			false,
			// O-1
			false,
			// O-2 - X
			true,
		},
	},
	{
		Type: "horizontal",
		Id:   "3",

		B: []bool{
			// B-0
			false,
			// B-1
			false,
			// B-2
			false,
			// B-3 - X
			true,
		},
		I: []bool{
			// I-0
			false,
			// I-1
			false,
			// I-2
			false,
			// I-3 - X
			true,
		},
		N: []bool{
			// N-0
			false,
			// N-1
			false,
			// N-2
			false,
			// N-3 - X
			true,
		},
		G: []bool{
			// G-0
			false,
			// G-1
			false,
			// G-2
			false,
			// G-3 - X
			true,
		},
		O: []bool{
			// O-0
			false,
			// O-1
			false,
			// O-2
			false,
			// O-3 - X
			true,
		},
	},
	{
		Type: "horizontal",
		Id:   "4",

		B: []bool{
			// B-0
			false,
			// B-1
			false,
			// B-2
			false,
			// B-3
			false,
			// B-4 - X
			true,
		},
		I: []bool{
			// I-0
			false,
			// I-1
			false,
			// I-2
			false,
			// I-3
			false,
			// I-4 - X
			true,
		},
		N: []bool{
			// N-0
			false,
			// N-1
			false,
			// N-2
			false,
			// N-3
			false,
			// N-4 - X
			true,
		},
		G: []bool{
			// G-0
			false,
			// G-1
			false,
			// G-2
			false,
			// G-3
			false,
			// G-4 - X
			true,
		},
		O: []bool{
			// O-0
			false,
			// O-1
			false,
			// O-2
			false,
			// O-3
			false,
			// O-4 - X
			true,
		},
	},

	// diagonal
	{
		Type: "diagonal",
		Id:   "B0-O4",

		B: []bool{
			// B-0 - X
			true,
			// B-1
			false,
			// B-2
			false,
			// B-3
			false,
			// B-4
			false,
		},
		I: []bool{
			// I-0
			false,
			// I-1 - X
			true,
			// I-2
			false,
			// I-3
			false,
			// I-4
			false,
		},
		N: []bool{
			// N-0
			false,
			// N-1
			false,
			// N-2 - X
			true,
			// N-3
			false,
			// N-4
			false,
		},
		G: []bool{
			// G-0
			false,
			// G-1
			false,
			// G-2
			false,
			// G-3 - X
			true,
			// G-4
			false,
		},
		O: []bool{
			// O-0
			false,
			// O-1
			false,
			// O-2
			false,
			// O-3
			false,
			// O-4 - X
			true,
		},
	},
	{
		Type: "diagonal",
		Id:   "B4-O0",

		B: []bool{
			// B-0
			false,
			// B-1
			false,
			// B-2
			false,
			// B-3
			false,
			// B-4 - X
			true,
		},
		I: []bool{
			// I-0
			false,
			// I-1
			false,
			// I-2
			false,
			// I-3 - X
			true,
			// I-4
			false,
		},
		N: []bool{
			// N-0
			false,
			// N-1
			false,
			// N-2 - X
			true,
			// N-3
			false,
			// N-4
			false,
		},
		G: []bool{
			// G-0
			false,
			// G-1 - X
			true,
			// G-2
			false,
			// G-3
			false,
			// G-4
			false,
		},
		O: []bool{
			// O-0 - X
			true,
			// O-1
			false,
			// O-2
			false,
			// O-3
			false,
			// O-4
			false,
		},
	},
}
