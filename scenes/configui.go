package scenes

import (
	img "image"
	"image/color"
	"net"
	"strconv"
	"tower-defense/assets"
	"tower-defense/config"

	"github.com/ebitenui/ebitenui"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
)

var isModalOpen = false

func IsModalOpen() bool {
	return isModalOpen
}

type GameOptionsCallback func(gameOptions *config.ConfigData)

func initUI(gameOptions *config.ConfigData, newGameCallback NewGameCallback, gameOptionsCallback GameOptionsCallback) *ebitenui.UI {
	ui := &ebitenui.UI{}
	buttonImage := loadButtonImage()
	face := assets.GoFace
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0x00})),
		widget.ContainerOpts.Layout(widget.NewAnchorLayout(
			widget.AnchorLayoutOpts.Padding(widget.NewInsetsSimple(20)))),
	)
	buttonContainer := widget.NewContainer(
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
			}),
		),

		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(50),
		)),
	)
	rootContainer.AddChild(buttonContainer)

	buttonStart := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Start Game", &face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			newGameCallback(true, &controller, gameOptions)
		}),
	)
	buttonContainer.AddChild(buttonStart)
	buttonMultiplayer := widget.NewButton(
		widget.ButtonOpts.Image(buttonImage),
		widget.ButtonOpts.Text("Game Options", &face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.TextPadding(widget.NewInsetsSimple(5)),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			openWindow(ui, gameOptions, gameOptionsCallback)
		}),
	)
	buttonContainer.AddChild(buttonMultiplayer)

	ui.Container = rootContainer

	return ui
}

func openWindow(ui *ebitenui.UI, gameOptions *config.ConfigData, callback GameOptionsCallback) {
	var rw widget.RemoveWindowFunc
	var window *widget.Window
	var saveButton *widget.Button
	var radioServer, radioClient *widget.Checkbox
	clr := color.NRGBA{254, 255, 255, 255}
	face := assets.GoFace
	imageBtn := loadButtonImage()
	padding := widget.NewInsetsSimple(5)
	colorBtnTxt := &widget.ButtonTextColor{Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff}}
	serverSelect := len(gameOptions.ServerPort) > 0 || len(gameOptions.ClientHostPort) == 0

	titleContainer := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{5, 50, 255, 255})),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(widget.GridLayoutOpts.Columns(3),
				widget.GridLayoutOpts.Stretch([]bool{true, false, false}, []bool{true}))))

	titleContainer.AddChild(widget.NewText(
		widget.TextOpts.Text("Multiplayer Configuration", &face, clr),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	titleContainer.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(imageBtn),
		widget.ButtonOpts.TextPadding(padding),
		widget.ButtonOpts.Text("X", &face, colorBtnTxt),
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
		widget.TextInputOpts.Face(&face),
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			saveButton.Click()
		}),
		widget.TextInputOpts.ChangedHandler(func(args *widget.TextInputChangedEventArgs) {
			if radioServer.State() == widget.WidgetChecked {
				saveButton.GetWidget().Disabled = !validateServerText(args.InputText)
			} else {
				saveButton.GetWidget().Disabled = !validateClientText(args.InputText)
			}
		}),
	}

	radioServer = widget.NewCheckbox(
		widget.CheckboxOpts.Spacing(10),
		widget.CheckboxOpts.Image(loadCheckboxGraphicImage()),
		widget.CheckboxOpts.Text("Server", &face, &widget.LabelColor{Idle: color.NRGBA{254, 255, 255, 255}}),
	)
	radioClient = widget.NewCheckbox(
		widget.CheckboxOpts.Spacing(10),
		widget.CheckboxOpts.Image(loadCheckboxGraphicImage()),
		widget.CheckboxOpts.Text("Client", &face, &widget.LabelColor{Idle: color.NRGBA{254, 255, 255, 255}}),
	)
	textServer := widget.NewTextInput(append(
		tiOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.TextInputOpts.Placeholder("port"))...,
	)
	textServer.SetText(gameOptions.ServerPort)

	textClient := widget.NewTextInput(append(
		tiOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.TextInputOpts.Placeholder("Server:port"))...,
	)
	textClient.SetText(gameOptions.ClientHostPort)

	windowContainer.AddChild(radioServer)
	windowContainer.AddChild(textServer)
	windowContainer.AddChild(radioClient)
	windowContainer.AddChild(textClient)
	var initialElement *widget.Checkbox
	if serverSelect {
		initialElement = radioServer
	} else {
		initialElement = radioClient
	}
	widget.NewRadioGroup(
		widget.RadioGroupOpts.Elements(radioServer, radioClient),
		widget.RadioGroupOpts.ChangedHandler(func(args *widget.RadioGroupChangedEventArgs) {
			if args.Active == radioServer {
				textServer.GetWidget().SetVisibility(widget.Visibility_Show)
				textClient.GetWidget().SetVisibility(widget.Visibility_Hide_Blocking)
				textServer.Focus(true)
			} else {
				textServer.GetWidget().SetVisibility(widget.Visibility_Hide_Blocking)
				textClient.GetWidget().SetVisibility(widget.Visibility_Show)
				textClient.Focus(true)
			}
		}),
		widget.RadioGroupOpts.InitialElement(initialElement),
	)

	var initialState widget.WidgetState
	if gameOptions.Debug {
		initialState = widget.WidgetChecked
	} else {
		initialState = widget.WidgetUnchecked
	}
	chkDebug := widget.NewCheckbox(
		widget.CheckboxOpts.Spacing(10),
		widget.CheckboxOpts.Image(loadCheckboxGraphicImage()),
		widget.CheckboxOpts.InitialState(initialState),
		widget.CheckboxOpts.Text("Debug Info", &face, &widget.LabelColor{Idle: color.NRGBA{254, 255, 255, 255}}),
	)
	if gameOptions.GridLines {
		initialState = widget.WidgetChecked
	} else {
		initialState = widget.WidgetUnchecked
	}
	chkGridLines := widget.NewCheckbox(
		widget.CheckboxOpts.Spacing(10),
		widget.CheckboxOpts.Image(loadCheckboxGraphicImage()),
		widget.CheckboxOpts.InitialState(initialState),
		widget.CheckboxOpts.Text("GridLines", &face, &widget.LabelColor{Idle: color.NRGBA{254, 255, 255, 255}}),
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
		widget.ButtonOpts.Text("Save", &face, colorBtnTxt),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			if radioServer.State() == widget.WidgetChecked {
				gameOptions.ServerPort = textServer.GetText()
			} else {
				gameOptions.ClientHostPort = textClient.GetText()
			}
			gameOptions.Debug = chkDebug.State() == widget.WidgetChecked
			gameOptions.GridLines = chkGridLines.State() == widget.WidgetChecked
			callback(gameOptions)
			rw()
		}),
	)
	bc.AddChild(saveButton)

	cancelButton := widget.NewButton(
		widget.ButtonOpts.Image(imageBtn),
		widget.ButtonOpts.TextPadding(padding),
		widget.ButtonOpts.Text("Cancel", &face, colorBtnTxt),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
	)
	bc.AddChild(cancelButton)

	if serverSelect {
		textServer.GetWidget().SetVisibility(widget.Visibility_Show)
		textClient.GetWidget().SetVisibility(widget.Visibility_Hide_Blocking)
		saveButton.GetWidget().Disabled = !validateServerText(textServer.GetText())
		textServer.Focus(true)
	} else {
		textClient.GetWidget().SetVisibility(widget.Visibility_Show)
		textServer.GetWidget().SetVisibility(widget.Visibility_Hide_Blocking)
		saveButton.GetWidget().Disabled = !validateClientText(textClient.GetText())
		textClient.Focus(true)
	}

	window = widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(windowContainer),
		widget.WindowOpts.TitleBar(titleContainer, 30),
		widget.WindowOpts.MinSize(400, 220),
		widget.WindowOpts.MaxSize(500, 400),
		widget.WindowOpts.ClosedHandler(func(args *widget.WindowClosedEventArgs) {
			isModalOpen = false
		}),
	)
	windowSize := input.GetWindowSize()
	x, y := window.Contents.PreferredSize()
	r := img.Rect(0, 0, x, y)
	r = r.Add(img.Point{windowSize.X/2 - x/2 - 50, windowSize.Y/2 - y/2 - 50})
	window.SetLocation(r)
	rw = ui.AddWindow(window)

	isModalOpen = true
}

func validateServerText(port string) bool {
	if len(port) == 0 {
		return true
	}
	return isNumber(port)
}
func validateClientText(hostport string) bool {
	if len(hostport) == 0 {
		return true
	}
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

func loadCheckboxGraphicImage() *widget.CheckboxImage {
	return &widget.CheckboxImage{
		Unchecked: image.NewNineSliceColor(color.NRGBA{R: 0, G: 0, B: 0, A: 255}),
		Checked:   image.NewNineSliceColor(color.NRGBA{R: 255, G: 255, B: 255, A: 255}),
	}
}
