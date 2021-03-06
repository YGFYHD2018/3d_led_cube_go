package ledlib

import (
	"3d_led_cube_go/ledlib/servicegateway"
	"3d_led_cube_go/ledlib/util"
	"math/rand"
	"sort"
	"sync"
	"time"
)

func red(ix int) uint8 {
	i := ix % 90
	switch {
	case i < 30:
		return uint8(i * 255 / 30)
	case i < 60:
		return uint8((60 - i) * 255 / 30)
	default:
		return 0
	}
}

func rgb(ix float64) util.Color32 {
	n := util.FloorToInt(ix * 1 * 90)
	return util.NewColorFromRGB(red(n), red(n+30), red(n+60))
}

type ObjectFireworks struct {
	cube        util.Image3D
	poss        []util.PointC
	vs          []util.Point
	ix          int
	addTimer    Timer
	updateTimer Timer
	sounds      []string
}

func localNewObjectFireworks() *ObjectFireworks {
	obj := ObjectFireworks{}
	obj.cube = NewLedImage3D()
	obj.poss = make([]util.PointC, 0)
	obj.vs = make([]util.Point, 0)
	obj.ix = 0
	obj.addTimer = NewTimer(2500 * time.Millisecond)
	obj.updateTimer = NewTimer(80 * time.Millisecond)
	obj.sounds = []string{"se_fireworks.wav", "se_fireworks2.wav", "se_fireworks3.wav"}

	return &obj
}

func NewObjectFireworks() LedObject {
	return localNewObjectFireworks()
}
func NewManagedObjectFireworks() LedManagedObject {
	return localNewObjectFireworks()
}

func (b *ObjectFireworks) IsExpired() bool {
	return false
}
func (b *ObjectFireworks) Draw(cube util.Image3D) {
	mux := &sync.Mutex{}
	if len(b.poss) < 250 {
		sound := b.sounds[rand.Intn(len(b.sounds))]
		servicegateway.GetAudigoSeriveGateway().Play(sound, false, false)

		cx := LedWidth * rand.Float64()
		cy := LedHeight * rand.Float64()
		cz := LedDepth * rand.Float64()

		util.ConcurrentEnum(0, 1000, func(i int) {
			sf := util.GetSphereFace()
			mux.Lock()
			b.vs = append(b.vs, sf)
			b.poss = append(b.poss, util.NewPointC(cx, cy, cz, rgb(sf.Len())))
			mux.Unlock()
		})
	}

	dIdx := make([]int, 0)

	isPast := b.updateTimer.IsPast()

	util.ConcurrentEnum(0, len(b.poss), func(i int) {
		p := b.poss[i]
		v := b.vs[i]
		if util.CanShow(p, LedWidth, LedHeight, LedDepth) {
			cube.SetAt(util.RoundToInt(p.X()),
				util.RoundToInt(p.Y()),
				util.RoundToInt(p.Z()),
				p.Color())
		} else {
			mux.Lock()
			dIdx = append(dIdx, i)
			mux.Unlock()
		}
		if isPast {
			p.Add(v)
			p.SetColor(util.Darken(p.Color()))
		}
	})
	if len(dIdx) > 0 {
		sort.Slice(dIdx, func(lhs, rhs int) bool { return dIdx[lhs] > dIdx[rhs] })
	}
	for i := 0; i < len(dIdx); i++ {
		b.vs = append(b.vs[:dIdx[i]], b.vs[dIdx[i]+1:]...)
		b.poss = append(b.poss[:dIdx[i]], b.poss[dIdx[i]+1:]...)
	}

}

func (b *ObjectFireworks) GetImage3D(param LedCanvasParam) util.ImmutableImage3D {
	b.cube.Clear()
	b.Draw(b.cube)
	return b.cube
}
