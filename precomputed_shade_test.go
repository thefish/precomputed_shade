package precomputed_shade

import (
	"fmt"
	"github.com/thefish/precomputed_shade/types"
	"strings"
	"testing"
)

func TestPrecompShade(t *testing.T) {
	ppFov := NewPrecomputedShade(15)
	_ = ppFov

	t.Log("ok")

	//level := gamemap.NewLevel(util.ClientCtx{}, "test", 1)

	level := &types.Level{
		Name:  "test1",
		Depth: 1,
		Rect:  types.NewRect(0, 0, 20, 20),
	}

	level.Tiles = make([]*types.Tile, level.W*level.H)

	var tile func() *types.Tile

	for x := 0; x < level.W; x++ {
		for y := 0; y < level.H; y++ {
			if x == 0 || y == 0 || x == (level.W-1) || y == (level.H-1) {
				tile = types.NewWall
			} else {
				tile = types.NewFloor
			}
			level.SetTileByXY(x, y, tile())
		}
	}

	playerCoords := types.Coords{10, 10}

	level.SetTileByXY(8, 12, types.NewWall())
	level.SetTileByXY(10, 8, types.NewWall())

	level.SetTileByXY(7, 9, types.NewWall())
	level.SetTileByXY(7, 11, types.NewWall())
	level.SetTileByXY(5, 10, types.NewWall())

	level.SetTileByXY(10, 11, types.NewWall())
	level.SetTileByXY(10, 12, types.NewWall())
	level.SetTileByXY(10, 13, types.NewWall())

	level.SetTileByXY(11, 10, types.NewWall())

	ppFov.ComputeFov(level, playerCoords, 12)

	fmt.Printf("\n\n")

	var render = func(x, y int) string {
		if playerCoords.X == x && playerCoords.Y == y {
			return "@"
		}
		result := level.GetTileByXY(x, y).Char
		if !level.GetTileByXY(x, y).Visible {
			result = "?"
		}
		return result
	}
	result := ""
	for y := 0; y < level.H; y++ {
		for x := 0; x < level.W; x++ {

			fmt.Printf("%s", render(x, y))
			result = result + fmt.Sprintf("%s", render(x, y))
		}
		fmt.Printf("\n")
		result = result + fmt.Sprintf("\n")
	}
	expected := strings.Join([]string{
		"???######???######??",
		"???......???......??",
		"??.......???.......#",
		"#........???.......#",
		"#........???.......#",
		"#.........?........#",
		"???.......?.......??",
		"?????.....?.....????",
		"???????...#...??????",
		"#......#.....???????",
		"?????#....@#????????",
		"#......#..#..???????",
		"???????.#.#...??????",
		"?????..?.???....????",
		"???..??.?????.....??",
		"#..????.?????......#",
		"#.????.???????.....#",
		"??????.???????.....#",
		"?????.?????????....#",
		"?????#?????????####?",
	}, "\n") + "\n"

	if result != expected {
		t.Fail()
	}
}
