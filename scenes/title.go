package scenes

import (
	"fmt"
	"image/color"

	"tower-defense/assets"
	comp "tower-defense/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"github.com/ebitenui/ebitenui"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

type TitleScene struct {
	width           int
	height          int
	gameStats       *GameStats
	newGameCallback NewGameCallback

	// ui elements
	ui *ebitenui.UI
}
type NewGameCallback func(broadcast bool) error

func NewTitleScene(width, height int, gameStats *GameStats, newGameCallback NewGameCallback) (*TitleScene, error) {
	title := &TitleScene{width: width, height: height, gameStats: gameStats, newGameCallback: newGameCallback}
	title.ui = initUI()
	return title, nil
}

func initUI() *ebitenui.UI {
	ui := &ebitenui.UI{}
	buttonImage, _ := loadButtonImage()
	face := assets.GoFace
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0x00})),

		// the container will use a row layout to layout the textinput widgets
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(20)))),
	)
	buttonMultiplayer := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Multiplayer Options", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 50,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			// TODO start the client connection to the server
			fmt.Printf("Loading Multiplayer options\n")
			openWindow2(ui)
		}),
	)

	// add the button as a child of the container
	rootContainer.AddChild(buttonMultiplayer)

	ui.Container = rootContainer

	return ui
}

func (t *TitleScene) Update() error {
	t.ui.Update()

	// if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) || ebiten.IsKeyPressed(ebiten.KeySpace) {
	// 	return t.newGameCallback(true)
	// }

	return nil
}

func (t *TitleScene) Draw(screen *ebiten.Image) {
	screen.Clear()

	backgroundImage := assets.GetImage("backgroundV")
	opts := &ebiten.DrawImageOptions{}
	screen.DrawImage(backgroundImage, opts)
	width := float64(t.width)
	halfWidth := width / 2

	str := fmt.Sprintf("HIGH %05d", t.gameStats.HighScore)
	nextY := comp.DrawTextLines(screen, assets.ScoreFace, str, width, comp.TextBorder, text.AlignEnd, text.AlignStart)
	str = fmt.Sprintf("High Creep Level %d\nHigh Tower Level %d\n", t.gameStats.HighCreepLevel, t.gameStats.HighTowerLevel)
	_ = comp.DrawTextLines(screen, assets.InfoFace, str, width, nextY, text.AlignEnd, text.AlignStart)

	str = "TOWER DEFENSE"
	_ = comp.DrawTextLines(screen, assets.ScoreFace, str, width, 100, text.AlignCenter, text.AlignStart)

	str = "Click to place towers Cost: $50"
	nextY = comp.DrawTextLines(screen, assets.ScoreFace, str, width, 250, text.AlignCenter, text.AlignStart)

	str = "Mouse over a tower\nPress H to heal to full Cost: $25\nPress U to upgrade and heal to full Cost: $50\nMax upgrade level is 5 (+1 for every 20 upgrades)"
	nextY = comp.DrawTextLines(screen, assets.InfoFace, str, width, nextY, text.AlignCenter, text.AlignStart)

	towerImage := assets.GetImage("tower")
	opts = &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(halfWidth-float64(towerImage.Bounds().Dx()/2), nextY)
	screen.DrawImage(towerImage, opts)
	nextY += float64(towerImage.Bounds().Dy()) + 10

	str = "Protect your base from aliens and earn $$"
	nextY = comp.DrawTextLines(screen, assets.ScoreFace, str, width, nextY, text.AlignCenter, text.AlignStart)

	const creepCount = 4
	const creepSize = 48
	x := halfWidth - creepSize*creepCount/2
	for i := 1; i <= creepCount; i++ {
		creepImage := assets.GetImage(fmt.Sprintf("creep%v", i))
		opts = &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(x, nextY)
		screen.DrawImage(creepImage, opts)
		x += float64(creepImage.Bounds().Dx())
	}

	nextY = 600
	str = "Click or press space to start"
	_ = comp.DrawTextLines(screen, assets.ScoreFace, str, width, nextY, text.AlignCenter, text.AlignStart)

	// draw UI elements
	t.ui.Draw(screen)
}
