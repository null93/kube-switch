package prompt

import (
	"fmt"
	"unicode"
	"github.com/fatih/color"
	term "github.com/nsf/termbox-go"
)

type Options struct {
	Choices [] string
	DefaultValue string
	Header string
	SelectionPrefix string
}

func Pick ( options Options ) ( string, bool ) {
	search := ""
	running := true
	filtered := options.Choices
	colorActive := color.New ( color.BgYellow ).Add ( color.FgBlack )
	index := FindIndex ( filtered, options.DefaultValue )
	if error := term.Init (); error != nil {
		panic ( error )
	}
	defer term.Close ()
	for running {
		term.SetCursor ( 0, 0 )
		fmt.Printf ( "Filter: %s%s\n", search, color.YellowString ("_") )
		fmt.Printf ( "Type to filter, UP/DOWN move, ENTER select, ESC exit\n" )
		fmt.Printf ( "\n" )
		fmt.Printf (
			"%s: %s%s\n",
			options.Header,
			options.SelectionPrefix,
			GetIndex ( filtered, index, color.RedString ("Nothing To Display") ),
		)
		fmt.Printf ( "\n" )
		for i, choice := range filtered {
			line := fmt.Sprintf ( " + %s \n", choice )
			if i == index {
				colorActive.Print ( line )
			} else {
				fmt.Print ( line )
			}
		}
		fmt.Printf ( "\n" )
		term.HideCursor ()
		switch event := term.PollEvent (); event.Type {
			case term.EventKey:
				switch true {
					case event.Key == term.KeyEsc:
						running = false
						break
					case event.Key == term.KeyArrowUp:
						index = Max ( index - 1, 0 )
					case event.Key == term.KeyArrowDown:
						index = Min ( index + 1, Max ( len ( filtered ) - 1, 0 ) )
					case event.Key == term.KeyEnter:
						if len ( filtered ) > 0 && index != -1 {
							return filtered [ index ], false
						}
					case event.Key == term.KeyBackspace || event.Key == term.KeyBackspace2:
						if len ( search ) > 0 {
							search = search [:len ( search ) - 1]
						}
					default:
						if unicode.IsGraphic ( event.Ch ) && unicode.IsPrint ( event.Ch ) {
							search = fmt.Sprintf ( "%s%c", search, event.Ch )
						}
				}
			case term.EventError:
				panic ( event.Err )
		}
		filtered = Filter ( options.Choices, search )
		if index > len ( filtered ) - 1 {
			index = 0
		}
		term.Sync ()
	}
	return "", true
}
