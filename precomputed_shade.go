package precomputed_shade

import (
	"bytes"
	"errors"
	"github.com/thefish/precomputed_shade/types"
	"math"
	"sort"
)

var errNotFoundCell = errors.New("Cell not found")
var errOutOfBounds = errors.New("Cell out of bounds")

type Cell struct {
	types.Coords
	distance       float64
	occludedAngles []int //angles occluded by this cell
	lit            int   //light "amount"
}

type CellList []*Cell

type DistanceSorter CellList

func (a DistanceSorter) Len() int           { return len(a) }
func (a DistanceSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a DistanceSorter) Less(i, j int) bool { return a[i].distance < a[j].distance }

type precomputedShade struct {
	originCoords   types.Coords
	MaxTorchRadius int
	CellList       CellList
	LightWalls     bool
}

func NewPrecomputedShade(maxTorchRadius int) *precomputedShade {
	result := &precomputedShade{MaxTorchRadius: maxTorchRadius, LightWalls: true}
	result.PrecomputeFovMap()
	return result
}

func (ps *precomputedShade) FindByCoords(c types.Coords) (int, *Cell, error) {
	for i := range ps.CellList {
		if ps.CellList[i].Coords == c {
			// Found!
			return i, ps.CellList[i], nil
		}
	}
	return 0, &Cell{}, errNotFoundCell
}

func (ps *precomputedShade) IsInFov(coords types.Coords) bool {
	rc := ps.fromLevelCoords(coords)
	_, cell, err := ps.FindByCoords(rc)
	if err != nil {
		return false
	}
	return cell.lit > 0
}

func (ps *precomputedShade) SetLightWalls(value bool) {
	ps.LightWalls = value
}

func (ps *precomputedShade) Init() {
	ps.PrecomputeFovMap()
}

func (ps *precomputedShade) PrecomputeFovMap() {
	max := ps.MaxTorchRadius
	minusMax := (-1) * max
	zeroCoords := types.Coords{0, 0}
	var x, y int
	//fill list
	for x = minusMax; x < max+1; x++ {
		for y = minusMax; y < max+1; y++ {
			if x == 0 && y == 0 {
				continue
			}
			iterCoords := types.Coords{x, y}
			distance := zeroCoords.DistanceTo(iterCoords)
			if distance <= float64(max) {
				ps.CellList = append(ps.CellList, &Cell{iterCoords, distance, nil, 0})
			}
		}
	}
	//Do not change cell order after this!
	sort.Sort(DistanceSorter(ps.CellList))
	//debug
	//for _, cell := range ps.CellList {
	//	fmt.Printf("\n coords: %v, distance: %f, len_occl: %d", cell.Coords, cell.distance, len(cell.occludedAngles))
	//}

	//Bresanham lines / Raycast
	var lineX, lineY float64
	for i := 0; i < 360; i++ {
		dx := math.Sin(float64(i) / (float64(180) / math.Pi))
		dy := math.Cos(float64(i) / (float64(180) / math.Pi))

		lineX = 0
		lineY = 0
		for j := 0; j < max; j++ {
			lineX -= dx
			lineY -= dy

			roundedX := int(round(lineX))
			roundedY := int(round(lineY))

			_, cell, err := ps.FindByCoords(types.Coords{roundedX, roundedY})

			if err != nil {
				//inexistent coord found
				break
			}
			cell.occludedAngles = unique(append(cell.occludedAngles, i))
		}

	}

	//for _, cell := range ps.CellList {
	//	fmt.Printf("\n coords: %v, distance: %f, len_occl: %d", cell.Coords, cell.distance, len(cell.occludedAngles))
	//}
}

func (ps *precomputedShade) recalc(level *types.Level, initCoords types.Coords, radius int) {
	for i, _ := range ps.CellList {
		ps.CellList[i].lit = 0
	}
	ps.originCoords = initCoords

	if radius > ps.MaxTorchRadius {
		radius = ps.MaxTorchRadius //fixme
	}

	level.GetTile(initCoords).Visible = true

	var fullShade = make([]byte, 360)
	for i := range fullShade {
		fullShade[i] = 1
	}
	var emptyShade = make([]byte, 360)
	currentShade := emptyShade
	nextShade := emptyShade

	i := 0
	prevDistance := 0.0
	for !bytes.Equal(currentShade, fullShade) {
		if i == len(ps.CellList)-1 {
			break
		}
		cell := ps.CellList[i]
		i++
		if cell.distance != prevDistance {
			currentShade = nextShade
		}

		if cell.distance > float64(radius) {
			break
		}

		lc, err := ps.toLevelCoords(level, initCoords, cell.Coords)
		if err != nil {
			continue
		}

		//fmt.Printf("\n level coords: %v", lc)
		for _, angle := range cell.occludedAngles {

			if level.GetTile(lc).BlocksSight {
				nextShade[angle] = 1
			}

			if currentShade[angle] == 0 {
				cell.lit = cell.lit + 1
			}

		}
	}
}

func (ps *precomputedShade) ComputeFov(level *types.Level, initCoords types.Coords, radius int) {

	ps.recalc(level, initCoords, radius)

	for _, cell := range ps.CellList {
		//fmt.Printf("\n coords: %v, distance: %f, lit: %d", cell.Coords, cell.distance, cell.lit)
		cs, err := ps.toLevelCoords(level, initCoords, cell.Coords)
		if cell.lit > 0 {
			if err != nil {
				continue
			}
			level.GetTile(cs).Visible = true
		}

		//light walls, crutch
		if level.GetTile(cs).BlocksSight && ps.LightWalls {
			if cell.IsAdjacentTo(&types.Coords{0, 0}) {
				level.GetTile(cs).Visible = true
			} else {
				for _, maybeNb := range ps.CellList {
					if //int(maybeNb.distance) == int(cell.distance-1) &&
					maybeNb.IsAdjacentTo(&cell.Coords) &&
						(maybeNb.X == cell.X || maybeNb.Y == cell.Y) &&
						maybeNb.lit > 0 { //magic constant!
						level.GetTile(cs).Visible = true
					}
				}
			}
		}
	}
}

func (ps *precomputedShade) toLevelCoords(level *types.Level, initCoords, relativeCoords types.Coords) (types.Coords, error) {
	realCoords := types.Coords{initCoords.X + relativeCoords.X, initCoords.Y + relativeCoords.Y}
	if !level.InBounds(realCoords) {
		return types.Coords{}, errOutOfBounds
	}
	return realCoords, nil
}

func (ps *precomputedShade) fromLevelCoords(lc types.Coords) types.Coords {
	relativeCoords := types.Coords{lc.X - ps.originCoords.X, lc.Y - ps.originCoords.Y}
	return relativeCoords
}

func unique(intSlice []int) []int {
	keys := make(map[int]bool)
	list := []int{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func round(f float64) float64 {
	return math.Floor(f + .5)
}
