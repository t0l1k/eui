package eui

type SceneDefault struct {
	Rect      *Rect
	Container []Drawable
}

func (sc *SceneDefault) Add(d Drawable) {
	sc.Container = append(sc.Container, d)
}
