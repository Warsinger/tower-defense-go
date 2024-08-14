package scenes

import (
	img "image"
	"image/color"
	"net"
	"strconv"
	"tower-defense/assets"

	"github.com/ebitenui/ebitenui"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
)

type GameOptions struct {
	connectString string
}

type GameOptionsCallback func(options GameOptions)

// TODO  radio buttons to specify client or server
// TODO other options including checkbox for debug mode and grid lines

func openWindow(ui *ebitenui.UI, gameOptions GameOptions, callback GameOptionsCallback) {
	var rw widget.RemoveWindowFunc
	var window *widget.Window
	clr := color.NRGBA{254, 255, 255, 255}
	face := assets.GoFace
	imageBtn := loadButtonImage()
	padding := widget.NewInsetsSimple(5)
	colorBtnTxt := &widget.ButtonTextColor{Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff}}

	titleContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{5, 50, 255, 255})),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(widget.GridLayoutOpts.Columns(3),
				widget.GridLayoutOpts.Stretch([]bool{true, false, false}, []bool{true}))))

	titleContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Multiplayer Configuration", face, clr),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	titleContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(imageBtn),
		widget.ButtonOpts.TextPadding(padding),
		widget.ButtonOpts.Text("X", face, colorBtnTxt),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
		widget.ButtonOpts.TabOrder(99),
	))

	windowContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{20, 100, 200, 255})),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
				widget.GridLayoutOpts.Padding(padding),
				widget.GridLayoutOpts.Spacing(0, 10),
			),
		),
	)

	windowContainer.AddChild(widget.NewText(widget.TextOpts.Text("Configure server", face, clr)))
	var saveButton *widget.Button
	ticolor := &widget.TextInputColor{
		Idle:  color.NRGBA{254, 255, 255, 255},
		Caret: color.NRGBA{254, 255, 255, 255},
	}
	tOpts := []widget.TextInputOpt{
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),
		widget.TextInputOpts.Color(ticolor),
		widget.TextInputOpts.Padding(padding),
		widget.TextInputOpts.Face(face),
		widget.TextInputOpts.CaretOpts(widget.CaretOpts.Size(face, 2)),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			saveButton.Click()
		}),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			saveButton.GetWidget().Disabled = !validateServerText(args.InputText)
		}),
	}

	serverText := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.TextInputOpts.Placeholder("Server:port"))...,
	)
	serverText.Focus(true)
	serverText.SetText(gameOptions.connectString)

	textContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout()))
	textContainer.AddChild(serverText)
	windowContainer.AddChild(textContainer)

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(15),
		)),
	)
	windowContainer.AddChild(bc)

	saveButton = widget.NewButton(
		widget.ButtonOpts.Image(imageBtn),
		widget.ButtonOpts.TextPadding(padding),
		widget.ButtonOpts.Text("Save", face, colorBtnTxt),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			callback(GameOptions{connectString: serverText.GetText()})
			rw()
		}),
	)
	bc.AddChild(saveButton)
	saveButton.GetWidget().Disabled = !validateServerText(serverText.GetText())

	cancelButton := widget.NewButton(
		widget.ButtonOpts.Image(imageBtn),
		widget.ButtonOpts.TextPadding(padding),
		widget.ButtonOpts.Text("Cancel", face, colorBtnTxt),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
	)
	bc.AddChild(cancelButton)

	window = widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(windowContainer),
		widget.WindowOpts.TitleBar(titleContainer, 30),
		// widget.WindowOpts.Draggable(),
		// widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(300, 175),
		widget.WindowOpts.MaxSize(500, 350),
	)
	windowSize := input.GetWindowSize()
	x, y := window.Contents.PreferredSize()
	//Create a rect with the preferred size of the content
	r := img.Rect(0, 0, x, y)
	r = r.Add(img.Point{windowSize.X/2 - x + 25, windowSize.Y/2 - y + 25})
	window.SetLocation(r)

	rw = ui.AddWindow(window)
}

func validateServerText(hostport string) bool {
	host, port, err := net.SplitHostPort(hostport)
	return err == nil && len(host) > 0 && len(port) > 0 && isNumber(port)
}

func isNumber(port string) bool {
	num, err := strconv.Atoi(port)
	return err == nil && num > 0 && num <= 65535
}

func loadButtonImage() *widget.ButtonImage {
	idle := image.NewNineSliceColor(color.NRGBA{R: 5, G: 50, B: 200, A: 255})
	hover := image.NewNineSliceColor(color.NRGBA{R: 0, G: 0, B: 150, A: 255})
	pressed := image.NewNineSliceColor(color.NRGBA{R: 0, G: 0, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}
}
