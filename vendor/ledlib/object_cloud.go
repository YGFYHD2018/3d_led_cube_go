package ledlib

import (
	"ledlib/util"
	"time"
)

type ObjectCloud struct {
	cloud []LedManagedObject
}

func NewObjectCloud(z int, updateRate time.Duration) LedManagedObject {
	obj := ObjectCloud{}
	obj.cloud = []LedManagedObject{
		NewObjectScrolledBitmap(
			"/asset/image/cloud/cloud2.png", z, updateRate, false),
		NewObjectScrolledBitmap(
			"/asset/image/cloud/cloud3.png", z+1, updateRate, false),
		NewObjectScrolledBitmap(
			"/asset/image/cloud/cloud4.png", z+2, updateRate, false),
		NewObjectScrolledBitmap(
			"/asset/image/cloud/cloud4.png", z+3, updateRate, false),
		NewObjectScrolledBitmap(
			"/asset/image/cloud/cloud3.png", z+4, updateRate, false),
		NewObjectScrolledBitmap(
			"/asset/image/cloud/cloud2.png", z+5, updateRate, false),
	}
	return &obj
}

func (o *ObjectCloud) Draw(cube util.Image3D) {
	for _, layer := range o.cloud {
		layer.Draw(cube)
	}
}
func (o *ObjectCloud) IsExpired() bool {
	for _, bmp := range o.cloud {
		if !bmp.IsExpired() {
			return false
		}
	}
	return true
}
