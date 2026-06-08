package main

import (
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/rsa17826/go-input-lib"
	"github.com/rsa17826/input-manager/IMan"
)

var evdevToChar = map[uint16]string{
	input.KEY_1: "1", input.KEY_2: "2", input.KEY_3: "3", input.KEY_4: "4", input.KEY_5: "5",
	input.KEY_6: "6", input.KEY_7: "7", input.KEY_8: "8", input.KEY_9: "9", input.KEY_0: "0",

	input.KEY_Q: "q", input.KEY_W: "w", input.KEY_E: "e", input.KEY_R: "r", input.KEY_T: "t",
	input.KEY_Y: "y", input.KEY_U: "u", input.KEY_I: "i", input.KEY_O: "o", input.KEY_P: "p",

	input.KEY_A: "a", input.KEY_S: "s", input.KEY_D: "d", input.KEY_F: "f", input.KEY_G: "g",
	input.KEY_H: "h", input.KEY_J: "j", input.KEY_K: "k", input.KEY_L: "l",

	input.KEY_Z: "z", input.KEY_X: "x", input.KEY_C: "c", input.KEY_V: "v", input.KEY_B: "b",
	input.KEY_N: "n", input.KEY_M: "m",

	input.KEY_SPACE: " ",

	// Unified Numpad Translations
	input.KEY_KP0: "0", input.KEY_KP1: "1", input.KEY_KP2: "2", input.KEY_KP3: "3", input.KEY_KP4: "4",
	input.KEY_KP5: "5", input.KEY_KP6: "6", input.KEY_KP7: "7", input.KEY_KP8: "8", input.KEY_KP9: "9",
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Background Input Capture")
	myWindow.Resize(fyne.NewSize(400, 100))

	textBox := widget.NewEntry()
	textBox.SetPlaceHolder("Captured background text will appear here...")
	myWindow.SetContent(container.NewVBox(textBox))

	mgr, err := IMan.Connect(IMan.ModeBlocking, IMan.ModeInjection)
	if err != nil {
		log.Fatalf("Failed to initialize input manager: %v", err)
	}
	defer mgr.Close()

	var finalResult string
	ctrlPressed := false

	// Helper to handle safe UI updates directly inside Fyne's main loop thread execution boundary
	updateTextSafe := func(newText string) {
		textBox.SetText(newText)
		textBox.Refresh()
	}

	go func() {
		for {
			routedEvent, err := mgr.ReadNext()
			if err != nil {
				return
			}

			if routedEvent.Event.Type == 1 { // EV_KEY
				code := routedEvent.Event.Code
				val := routedEvent.Event.Value // 0 = Release, 1 = Press, 2 = Repeat

				if code == input.KEY_LEFTCTRL || code == input.KEY_RIGHTCTRL {
					if val == 1 {
						ctrlPressed = true
					} else if val == 0 {
						ctrlPressed = false
					}
				}

				if val == 1 || val == 2 {
					switch code {
					case input.KEY_ESC:
						finalResult = ""
						myApp.Quit()
						return

					case input.KEY_ENTER, input.KEY_KPENTER:
						finalResult = textBox.Text
						myApp.Quit()
						return

					case input.KEY_BACKSPACE:
						currentText := textBox.Text
						if ctrlPressed {
							trimmed := strings.TrimRight(currentText, " ")
							lastSpace := strings.LastIndex(trimmed, " ")
							if lastSpace == -1 {
								updateTextSafe("")
							} else {
								updateTextSafe(trimmed[:lastSpace+1])
							}
						} else {
							if len(currentText) > 0 {
								updateTextSafe(currentText[:len(currentText)-1])
							}
						}

					default:
						if char, found := evdevToChar[code]; found {
							updateTextSafe(textBox.Text + char)
						}
					}
				}
			}

			if routedEvent.From == IMan.ModeBlocking {
				_, _ = mgr.BlockInput(1)
			}
		}
	}()

	myWindow.ShowAndRun()
	fmt.Println(finalResult)
}
