package scenes

import (
	"fmt"
	img "image"
	"image/color"

	"tower-defense/assets"

	"github.com/ebitenui/ebitenui"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/input"
	"github.com/ebitenui/ebitenui/widget"
)

func openWindow(ui *ebitenui.UI) {
	var rw widget.RemoveWindowFunc
	var window *widget.Window
	clr := color.NRGBA{254, 255, 255, 255}
	face := assets.GoFace
	imageBtn, _ := loadButtonImage()
	padding := widget.Insets{
		Left:   30,
		Right:  30,
		Top:    5,
		Bottom: 5,
	}
	colorBtn := &widget.ButtonTextColor{
		Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff}}

	titleBar := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{100, 100, 100, 255})),
		widget.ContainerOpts.Layout(widget.NewGridLayout(widget.GridLayoutOpts.Columns(3), widget.GridLayoutOpts.Stretch([]bool{true, false, false}, []bool{true}), widget.GridLayoutOpts.Padding(widget.Insets{
			Left:   30,
			Right:  5,
			Top:    6,
			Bottom: 5,
		}))))

	titleBar.AddChild(widget.NewText(
		widget.TextOpts.Text("Multiplayer Configuration", face, clr),
		widget.TextOpts.Position(widget.TextPositionStart, widget.TextPositionCenter),
	))

	titleBar.AddChild(widget.NewButton(
		widget.ButtonOpts.Image(imageBtn),
		widget.ButtonOpts.TextPadding(padding),
		widget.ButtonOpts.Text("X", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
		widget.ButtonOpts.TabOrder(99),
	))

	c := widget.NewContainer(
		// widget.ContainerOpts.BackgroundImage(res.panel.image),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(1),
				widget.GridLayoutOpts.Stretch([]bool{true}, []bool{false, true, false}),
				// widget.GridLayoutOpts.Padding(res.panel.padding),
				widget.GridLayoutOpts.Spacing(0, 15),
			),
		),
	)

	c.AddChild(widget.NewText(
		widget.TextOpts.Text("Configure server", face, clr),
	))

	tOpts := []widget.TextInputOpt{
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Stretch: true,
		})),
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			Disabled: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),
		// widget.TextInputOpts.Color(res.textInput.color),
		widget.TextInputOpts.Padding(widget.Insets{
			Left:   13,
			Right:  13,
			Top:    7,
			Bottom: 7,
		}),
		widget.TextInputOpts.Face(face),
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(face, 2),
		),
	}

	serverText := widget.NewTextInput(append(
		tOpts,
		widget.TextInputOpts.WidgetOpts(widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
			StretchHorizontal: true,
		})),
		widget.TextInputOpts.Placeholder("Server:port"))...,
	)
	textContainer := widget.NewContainer(widget.ContainerOpts.Layout(widget.NewAnchorLayout()))
	textContainer.AddChild(serverText)
	c.AddChild(textContainer)

	bc := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Spacing(15),
		)),
	)
	c.AddChild(bc)

	saveButton := widget.NewButton(
		widget.ButtonOpts.Image(imageBtn),
		widget.ButtonOpts.TextPadding(padding),
		widget.ButtonOpts.Text("Save", face, colorBtn),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			fmt.Printf("Save button clicked %v\n", serverText.GetText())
			rw()
		}),
	)
	bc.AddChild(saveButton)

	cancelButton := widget.NewButton(
		widget.ButtonOpts.Image(imageBtn),
		widget.ButtonOpts.TextPadding(padding),
		widget.ButtonOpts.Text("Cancel", face, colorBtn),
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			rw()
		}),
	)
	bc.AddChild(cancelButton)

	window = widget.NewWindow(
		widget.WindowOpts.Modal(),
		widget.WindowOpts.Contents(c),
		widget.WindowOpts.TitleBar(titleBar, 30),
		widget.WindowOpts.Draggable(),
		widget.WindowOpts.Resizeable(),
		widget.WindowOpts.MinSize(500, 200),
		widget.WindowOpts.MaxSize(700, 400),
		widget.WindowOpts.ResizeHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Resize: ", args.Rect)
		}),
		widget.WindowOpts.MoveHandler(func(args *widget.WindowChangedEventArgs) {
			fmt.Println("Move: ", args.Rect)
		}),
	)
	windowSize := input.GetWindowSize()
	r := img.Rect(0, 0, 550, 250)
	r = r.Add(img.Point{windowSize.X / 4 / 2, windowSize.Y * 2 / 3 / 2})
	window.SetLocation(r)

	rw = ui.AddWindow(window)
}
func initMultiplayerConfigUI() *ebitenui.UI {
	buttonImage, _ := loadButtonImage()
	face := assets.GoFace
	rootContainer := widget.NewContainer(
		// the container will use a plain color as its background
		widget.ContainerOpts.BackgroundImage(image.NewNineSliceColor(color.NRGBA{0x13, 0x1a, 0x22, 0x00})),

		// the container will use a row layout to layout the textinput widgets
		widget.ContainerOpts.Layout(widget.NewRowLayout(
			widget.RowLayoutOpts.Direction(widget.DirectionVertical),
			widget.RowLayoutOpts.Spacing(20),
			widget.RowLayoutOpts.Padding(widget.NewInsetsSimple(20)))),
	)

	// construct a standard textinput widget
	standardTextInput := widget.NewTextInput(
		widget.TextInputOpts.WidgetOpts(
			//Set the layout information to center the textbox in the parent
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
		),

		//Set the Idle and Disabled background image for the text input
		//If the NineSlice image has a minimum size, the widget will use that or
		// widget.WidgetOpts.MinSize; whichever is greater
		widget.TextInputOpts.Image(&widget.TextInputImage{
			Idle:     image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
			Disabled: image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 100, A: 255}),
		}),

		//Set the font face and size for the widget
		widget.TextInputOpts.Face(face),

		//Set the colors for the text and caret
		widget.TextInputOpts.Color(&widget.TextInputColor{
			Idle:          color.NRGBA{254, 255, 255, 255},
			Disabled:      color.NRGBA{R: 200, G: 200, B: 200, A: 255},
			Caret:         color.NRGBA{254, 255, 255, 255},
			DisabledCaret: color.NRGBA{R: 200, G: 200, B: 200, A: 255},
		}),

		//Set how much padding there is between the edge of the input and the text
		widget.TextInputOpts.Padding(widget.NewInsetsSimple(5)),

		//Set the font and width of the caret
		widget.TextInputOpts.CaretOpts(
			widget.CaretOpts.Size(face, 2),
		),

		//This text is displayed if the input is empty
		widget.TextInputOpts.Placeholder("Server:port"),

		//This is called when the user hits the "Enter" key.
		//There are other options that can configure this behavior
		widget.TextInputOpts.SubmitHandler(func(args *widget.TextInputChangedEventArgs) {
			fmt.Println("Text Submitted: ", args.InputText)
		}),
	)

	rootContainer.AddChild(standardTextInput)

	buttonServer := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Start server", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			// TODO start the client connection to the server
			fmt.Printf("Starting server on port %v\n", standardTextInput.GetText())
		}),
	)

	rootContainer.AddChild(buttonServer)

	buttonClient := widget.NewButton(
		// set general widget options
		widget.ButtonOpts.WidgetOpts(
			// instruct the container's anchor layout to center the button both horizontally and vertically
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionCenter,
			}),
		),

		// specify the images to use
		widget.ButtonOpts.Image(buttonImage),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text("Connect to server", face, &widget.ButtonTextColor{
			Idle: color.NRGBA{0xdf, 0xf4, 0xff, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Left:   30,
			Right:  30,
			Top:    5,
			Bottom: 5,
		}),

		// add a handler that reacts to clicking the button
		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			// TODO start the client connection to the server
			fmt.Printf("Connecting to the server %v\n", standardTextInput.GetText())
		}),
	)

	// add the button as a child of the container
	rootContainer.AddChild(buttonClient)

	// rootContainer.SetLocation(img.Rect(0, 650, t.width, 800))
	// standardTextInput.SetLocation(img.Rect(0, 650, t.width, 700))
	// buttonServer.SetLocation(img.Rect(0, 700, t.width/2, 750))
	// buttonClient.SetLocation(img.Rect(t.width/2+1, 700, t.width, 750))
	// construct the UI
	ui := ebitenui.UI{
		Container: rootContainer,
	}
	return &ui
}

func loadButtonImage() (*widget.ButtonImage, error) {
	idle := image.NewNineSliceColor(color.NRGBA{R: 170, G: 170, B: 180, A: 255})

	hover := image.NewNineSliceColor(color.NRGBA{R: 130, G: 130, B: 150, A: 255})

	pressed := image.NewNineSliceColor(color.NRGBA{R: 100, G: 100, B: 120, A: 255})

	return &widget.ButtonImage{
		Idle:    idle,
		Hover:   hover,
		Pressed: pressed,
	}, nil
}
