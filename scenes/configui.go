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

var isModalOpen = false

func IsModalOpen() bool {
	return isModalOpen
}

func NewGameOptions(debug bool) *GameOptions {
	return &GameOptions{debug: debug}
}

type GameOptionsCallback func(options *GameOptions)

func initUI(gameOptions *GameOptions, callback GameOptionsCallback) *ebitenui.UI {
	ui := &ebitenui.UI{}
	buttonImage := loadButtonImage()
	face := assets.GoFace
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0x00})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(20)))),
	)
	buttonMultiplayer := widget.NewButton(
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
		),

		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Multiplayer Options", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openWindow(ui, gameOptions, callback)
		}),
	)
	rootContainer.AddChild(buttonMultiplayer)

	ui.Container = rootContainer

	return ui
}

func openWindow(ui *ebitenui.UI, gameOptions *GameOptions, callback GameOptionsCallback) {
	var rw widget.RemoveWindowFunc
	var window *widget.Window
	var saveButton *widget.Button
	var radioServer, radioClient *widget.LabeledCheckbox
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
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Stretch([]bool{false, true}, []bool{false, false, false, false}),
				widget.GridLayoutOpts.Padding(padding),
				widget.GridLayoutOpts.Spacing(0, 10),
			),
		),
	)

	ticolor := &widget.TextInputColor{
		Idle:  color.NRGBA{254, 255, 255, 255},
		Caret: color.NRGBA{254, 255, 255, 255},
	}
	tiOpts := []widget.TextInputOpt{
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

	radioServer = widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(10),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(imageBtn)),
			widget.CheckboxOpts.Image(loadCheckboxGraphicImage()),
		),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Server", face, &widget.LabelColor{Idle: color.NRGBA{254, 255, 255, 255}})),
	)
	radioClient = widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(10),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(imageBtn)),
			widget.CheckboxOpts.Image(loadCheckboxGraphicImage()),
		),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Client", face, &widget.LabelColor{Idle: color.NRGBA{254, 255, 255, 255}})),
	)
	textServer := widget.NewTextInput(append(
		tiOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.TextInputOpts.Placeholder("port"))...,
	)
	textServer.SetText(gameOptions.serverPort)

	textClient := widget.NewTextInput(append(
		tiOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.TextInputOpts.Placeholder("Server:port"))...,
	)
	textClient.SetText(gameOptions.clientHostPort)

	windowContainer.AddChild(radioServer)
	windowContainer.AddChild(textServer)
	windowContainer.AddChild(radioClient)
	windowContainer.AddChild(textClient)
	var initialElement *widget.Checkbox
	if serverSelect {
		initialElement = radioServer.Checkbox()
	} else {
		initialElement = radioClient.Checkbox()
	}
	widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(radioServer.Checkbox(), radioClient.Checkbox()),
		widget.RadioGroupOpts.ChangedHandler(func(args *widget.RadioGroupChangedEventArgs) {
			if args.Active == radioServer.Checkbox() {
				textServer.GetWidget().Visibility = widget.Visibility_Show
				textClient.GetWidget().Visibility = widget.Visibility_Hide_Blocking
				textServer.Focus(true)
			} else {
				textServer.GetWidget().Visibility = widget.Visibility_Hide_Blocking
				textClient.GetWidget().Visibility = widget.Visibility_Show
				textClient.Focus(true)
			}
		}),
		widget.RadioGroupOpts.InitialElement(initialElement),
	)

	var initialState widget.WidgetState
	if gameOptions.debug {
		initialState = widget.WidgetChecked
	} else {
		initialState = widget.WidgetUnchecked
	}
	chkDebug := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(10),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(imageBtn)),
			widget.CheckboxOpts.Image(loadCheckboxGraphicImage()),
			widget.CheckboxOpts.InitialState(initialState),
		),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("Debug Info", face, &widget.LabelColor{Idle: color.NRGBA{254, 255, 255, 255}})),
	)
	if gameOptions.gridlines {
		initialState = widget.WidgetChecked
	} else {
		initialState = widget.WidgetUnchecked
	}
	chkGridLines := widget.NewLabeledCheckbox(
		widget.LabeledCheckboxOpts.Spacing(10),
		widget.LabeledCheckboxOpts.CheckboxOpts(
			widget.CheckboxOpts.ButtonOpts(widget.ButtonOpts.Image(imageBtn)),
			widget.CheckboxOpts.Image(loadCheckboxGraphicImage()),
			widget.CheckboxOpts.InitialState(initialState),
		),
		widget.LabeledCheckboxOpts.LabelOpts(widget.LabelOpts.Text("GridLines", face, &widget.LabelColor{Idle: color.NRGBA{254, 255, 255, 255}})),
	)
	windowContainer.AddChild(chkDebug)
	windowContainer.AddChild(chkGridLines)

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(50),
		)),
	)
	windowContainer.AddChild(bc)

	saveButton = widget.NewButton(
		widget.ButtonOpts.Image(imageBtn),
		widget.ButtonOpts.TextPadding(padding),
		widget.ButtonOpts.Text("Save", face, colorBtnTxt),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if radioServer.Checkbox().State() == widget.WidgetChecked {
				gameOptions.serverPort = textServer.GetText()
			} else {
				gameOptions.clientHostPort = textClient.GetText()
			}
			gameOptions.debug = chkDebug.Checkbox().State() == widget.WidgetChecked
			gameOptions.gridlines = chkGridLines.Checkbox().State() == widget.WidgetChecked
			callback(gameOptions)
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

	if serverSelect {
		textServer.GetWidget().Visibility = widget.Visibility_Show
		textClient.GetWidget().Visibility = widget.Visibility_Hide_Blocking
		saveButton.GetWidget().Disabled = !validateServerText(textServer.GetText())
		textServer.Focus(true)
	} else {
		textClient.GetWidget().Visibility = widget.Visibility_Show
		textServer.GetWidget().Visibility = widget.Visibility_Hide_Blocking
		saveButton.GetWidget().Disabled = !validateClientText(textClient.GetText())
		textClient.Focus(true)
	}

	window = widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(windowContainer),
		widget.WindowOpts.TitleBar(titleContainer, 30),
		widget.WindowOpts.MinSize(400, 220),
		widget.WindowOpts.MaxSize(500, 400),
	)
	windowSize := input.GetWindowSize()
	x, y := window.Contents.PreferredSize()
	r := img.Rect(0, 0, x, y)
	r = r.Add(img.Point{windowSize.X/2 - x/2 + 25, windowSize.Y/2 - y/2 + 25})
	window.SetLocation(r)

	window.SetCloseFunction(func() { isModalOpen = false })
	rw = ui.AddWindow(window)
	isModalOpen = true
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
