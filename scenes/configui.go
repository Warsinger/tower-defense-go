package scenes

import (
	img "image"
	"image/color"
	"net"
	"strconv"
	"tower-defense/assets"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
)

type GameOptions struct {
	serverPort     string
	clientHostPort string
	debug          bool
	gridlines      bool
}

type GameOptionsCallback func(options GameOptions)

// TODO  radio buttons to specify client or server
// TODO other options including checkbox for debug mode and grid lines

func openWindow(ui *ebitenui.UI, gameOptions GameOptions, callback GameOptionsCallback) {
	var rw widget.RemoveWindowFunc
	var window *widget.Window
	var serverText *widget.TextInput
	var clientText *widget.TextInput
	var saveButton *widget.Button
	clr := color.NRGBA{254, 255, 255, 255}
	face := assets.GoFace
	imageBtn := loadButtonImage()
	padding := widget.NewInsetsSimple(5)
	colorBtnTxt := &widget.ButtonTextColor{Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff}}
	serverSelect := len(gameOptions.serverPort) > 0 || len(gameOptions.clientHostPort) == 0

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

	imgBackground := image.NewNineSliceColor(color.NRGBA{20, 100, 200, 255})
	windowContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(imgBackground),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
				widget.GridLayoutOpts.Padding(padding),
				widget.GridLayoutOpts.Spacing(0, 10),
			),
		),
	)

	radioServer := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(10),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(imageBtn)),
			widget.CheckboxOpts.Image(loadCheckboxGraphicImage()),
		),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Server", face, &widget.LabelColor{Idle: color.NRGBA{254, 255, 255, 255}})),
	)
	radioClient := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(10),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(imageBtn)),
			widget.CheckboxOpts.Image(loadCheckboxGraphicImage()),
		),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Client", face, &widget.LabelColor{Idle: color.NRGBA{254, 255, 255, 255}})),
	)

	windowContainer.AddChild(radioServer)
	windowContainer.AddChild(radioClient)
	widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(radioServer.Checkbox(), radioClient.Checkbox()),
		widget.RadioGroupOpts.ChangedHandler(func(args *widget.RadioGroupChangedEventArgs) {
			if args.Active == radioServer.Checkbox() {
				serverText.GetWidget().Visibility = widget.Visibility_Show
				clientText.GetWidget().Visibility = widget.Visibility_Hide_Blocking
				serverText.Focus(true)
			} else {
				serverText.GetWidget().Visibility = widget.Visibility_Hide_Blocking
				clientText.GetWidget().Visibility = widget.Visibility_Show
				clientText.Focus(true)
			}
		}),
	)

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
			if radioServer.Checkbox().State() == widget.WidgetChecked {
				saveButton.GetWidget().Disabled = !validateServerText(args.InputText)
			} else {
				saveButton.GetWidget().Disabled = !validateClientText(args.InputText)
			}
		}),
	}

	textContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewGridLayout(
		widget.GridLayoutOpts.Columns(1),
		widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
		widget.GridLayoutOpts.Padding(padding),
		widget.GridLayoutOpts.Spacing(0, 10),
	)))
	serverText = widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.TextInputOpts.Placeholder("port"))...,
	)
	serverText.SetText(gameOptions.serverPort)

	textContainer.AddChild(serverText)
	clientText = widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.TextInputOpts.Placeholder("Server:port"))...,
	)
	clientText.SetText(gameOptions.clientHostPort)
	textContainer.AddChild(clientText)
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
			gameOpts := GameOptions{}
			if radioServer.Checkbox().State() == widget.WidgetChecked {
				gameOpts.serverPort = serverText.GetText()
			} else {
				gameOpts.clientHostPort = clientText.GetText()
			}
			callback(gameOpts)
			rw()
		}),
	)
	bc.AddChild(saveButton)

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
		widget.WindowOpts.MinSize(300, 175),
		widget.WindowOpts.MaxSize(500, 350),
	)
	windowSize := input.GetWindowSize()
	x, y := window.Contents.PreferredSize()
	r := img.Rect(0, 0, x, y)
	r = r.Add(img.Point{windowSize.X/2 - x + 25, windowSize.Y/2 - y + 25})
	window.SetLocation(r)

	rw = ui.AddWindow(window)

	if serverSelect {
		serverText.Focus(true)
		serverText.GetWidget().Visibility = widget.Visibility_Show
		saveButton.GetWidget().Disabled = !validateServerText(serverText.GetText())
		clientText.GetWidget().Visibility = widget.Visibility_Hide_Blocking
		radioServer.SetState(widget.WidgetChecked)
	} else {
		clientText.Focus(true)
		clientText.GetWidget().Visibility = widget.Visibility_Show
		serverText.GetWidget().Visibility = widget.Visibility_Hide_Blocking
		saveButton.GetWidget().Disabled = !validateClientText(clientText.GetText())
		radioClient.SetState(widget.WidgetChecked)
	}
}

func validateServerText(port string) bool {
	return isNumber(port)
}
func validateClientText(hostport string) bool {
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

func loadCheckboxGraphicImage() *widget.CheckboxGraphicImage {
	size := 17
	radius := float32(size / 2)
	center := float32(radius + 1)
	unchecked := ebiten.NewImage(size, size)
	vector.DrawFilledCircle(unchecked, center, center, radius, color.Black, true)
	checked := ebiten.NewImage(17, 17)
	vector.DrawFilledCircle(checked, center, center, radius, color.White, true)

	return &widget.CheckboxGraphicImage{
		Unchecked: &widget.ButtonImageImage{Idle: unchecked},
		Checked:   &widget.ButtonImageImage{Idle: checked},
	}
}
