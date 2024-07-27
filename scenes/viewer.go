package scenes

import (
	"image/color"
	"tower-defense/assets"
	comp "tower-defense/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

type ViewerScene struct {
	world  donburi.World
	width  int
	height int
	image  *ebiten.Image
}

func NewViewerScene(world donburi.World, width, height int) (*ViewerScene, error) {
	return &ViewerScene{world: world, width: width, height: height}, nil
}

func (v *ViewerScene) Update() error {
	srvsync.UseEsync()
	return nil
}

func (v *ViewerScene) Draw(screen *ebiten.Image) {
	if v.image == nil {
		v.image = ebiten.NewImage(v.width, v.height)
	}
	v.image.Clear()

	img := assets.GetImage("backgroundV")
	opts := &ebiten.DrawImageOptions{}
	v.image.DrawImage(img, opts)

	vector.StrokeLine(v.image, 0, 0, 0, float32(v.height), 3, color.White, true)

	// query for all entities
	query := donburi.NewQuery(
		filter.And(
			filter.Contains(comp.Position, comp.Render),
		),
	)

	// draw all entities
	query.Each(v.world, func(entry *donburi.Entry) {
		r := comp.Render.Get(entry)
		r.Draw(v.image, entry)
	})

	v.DrawText(v.image)

	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(v.width), 0)
	screen.DrawImage(v.image, opts)
}

func (v *ViewerScene) DrawText(image *ebiten.Image) {
	be := comp.Board.MustFirst(v.world)
	board := comp.Board.Get(be)

	// draw high score
	str := "Viewer Mode"
	op := &text.DrawOptions{}
	x, _ := text.Measure(str, assets.ScoreFace, op.LineSpacing)
	op.GeoM.Translate(float64(board.Width)-x-comp.TextBorder, comp.TextBorder)
	text.Draw(image, str, assets.InfoFace, op)

}
