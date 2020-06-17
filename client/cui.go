package client

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"os"
)

var (
	OutputView *gocui.View
	LogView *gocui.View
	FilesView *gocui.View
	Gui *gocui.Gui
)

func layout(g *gocui.Gui) error {
	Gui = g
	maxX, maxY := g.Size()
	fv, err := g.SetView("side", 0, 0, int(0.2*float32(maxX)), maxY-20)
	if err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}

	fv.Title = "Changed files"
	fv.Autoscroll = true
	fv.FgColor = gocui.ColorYellow
	FilesView = fv

	ov, err := g.SetView("main", int(0.2*float32(maxX)), 0, maxX, maxY-20)

	if err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}

	ov.Title = "Remote output"
	ov.Autoscroll = true
	OutputView = ov

	lv, err := g.SetView("cmdline", -1, maxY-20, maxX, maxY)
	if err != nil &&
		err != gocui.ErrUnknownView {
		return err
	}

	lv.Title = "Piper log"
	lv.Autoscroll = true
	LogView = lv
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	os.Exit(1)
	return gocui.ErrQuit
}

func SetupCui(callback func()) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	go callback()

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		PushChanges()
		return nil
	}); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func Log(text string)  {
	Gui.Update(func(gui *gocui.Gui) error {
		fmt.Fprintln(LogView, text + "\033[0m")
		return nil
	})
}

func PrintRemote(text string)  {
	Gui.Update(func(gui *gocui.Gui) error {
		fmt.Fprintln(OutputView, text + "\033[0m")
		return nil
	})
}

func PrintFiles(text string)  {
	Gui.Update(func(gui *gocui.Gui) error {
		fmt.Fprintln(FilesView, text + "\033[0m")
		return nil
	})
}