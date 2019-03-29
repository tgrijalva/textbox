package main

import (
	"fmt"

	"github.com/tgrijalva/textbox"
)

func main() {

	// Create a new Textbox
	fmt.Println("Creating new Textbox.")
	box := textbox.NewTextbox(50, 20)

	// Fill it with char
	fmt.Println("Filling box.")
	box.Fill('x')

	// Print the box
	fmt.Println("Printing box.")
	fmt.Println(box, "\n")

	// Create a second box
	fmt.Println("Creating another Textbox.")
	b2 := textbox.NewTextbox(25, 6)

	// Fill that one
	fmt.Println("Filling box2 with junk.")
	b2.WriteString("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!XX!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	// Whe the box is full the write will stop

	// Print box2
	fmt.Println("Printing box2.")
	fmt.Println(b2, "\n")

	// Now draw box2 into box1 and print it out
	fmt.Println("Drawing box2 into box1.")
	box.Draw(b2, textbox.Point{}, nil)
	fmt.Println(box, "\n")

	// You can also change the drawing location of box2
	fmt.Println("Drawing box2 into box1 again, but in a different location.")
	box.Draw(b2, textbox.Point{35, 15}, nil)
	fmt.Println(box, "\n")

	// you can also swap characters
	fmt.Println("Here, we're swapping out X's for O's in box2, usint the 'Tl' command.")
	b2.Replace('X', 'O')
	fmt.Println(b2, "\n")
	fmt.Println("And of course the O's will show up the next time we draw box2 into box1.")
	box.Draw(b2, textbox.Point{-7, 15}, nil)
	fmt.Println(box, "\n")

	// Textboxes also have cursors
	fmt.Println("Now lets write the character 'Z' at the current location of the cursor in box1")
	box.WriteString("Z")
	fmt.Println(box)
	fmt.Println("When we print it out we can see it along the lower edge.")
	fmt.Println("It's just to the right of the place where the most recient 'Draw' command into box1 ended.", "\n")

	// Cursors increment along to the right when writing. They also wordwrap.
	fmt.Println("That's no coincidence. See what happens when we write another Z.")
	box.WriteString("Z")
	fmt.Println(box)
	fmt.Println("This time the second 'Z' is just to the right of the previous Z.")
	fmt.Println("Cursors increment to the right when writing.", "\n")

	// Set the cursor manually. Anywhere inside the boxes boundry.
	fmt.Println("You can manually set the location of the cursor before writing.")
	fmt.Println("Lets set the cursor then write a long string.")
	p := textbox.Point{}
	p.X, p.Y = box.Cursor()
	box.SetCursor(p.X-5, p.Y-3)
	box.WriteString("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.")
	fmt.Println(box, "\n")

	// You can't overwrite the box boundry
	fmt.Println("You can't overrun the box boundry.")
	fmt.Println("when the cursor reaches the end, the Write stops.", "\n")

	// More interesting uses
	fmt.Println("Now for something more exciting. How about dragons?")
	tbox := textbox.FromStrings(trogdor...)
	fmt.Println(tbox)
	fmt.Println("Here is a dragon inside of a new Textbox.")
	fmt.Println("This Textbox was created from a []string using the FromStrings method.", "\n")

	// Targeted replacement
	fmt.Println("Now lets replaceReplace will null values.")
	tbox.Replace(0, ' ')
	fmt.Println("Then apply the dragon box to our old friend box1.")
	box.Draw(tbox, textbox.Point{0, -1}, []rune{' '})
	fmt.Println(box, "\n")
	fmt.Println("Notice how the null characters don't get drawn into box1.")
	fmt.Println("This lets null values work as an alpha-mask.", "\n")

	fmt.Println("Now we set box1's cursor back to the origin, and then we write some junk to it.")
	box.SetCursor(0, 0)
	for i := 0; i < 20; i++ {
		box.WriteString("///////////////////////////")
	}
	fmt.Println(box, "\n")
	fmt.Println("After that we can re-apply the dragon if we like, to end up with a different background.")
	box.Draw(tbox, textbox.Point{0, -1}, nil)
	fmt.Println(box, "\n")

	fmt.Println("Finally, we trade out the null values for spaces again, and then draw the it one last time.")
	tbox.Replace(0, ' ')
	box.Draw(tbox, textbox.Point{0, -1}, nil)
	fmt.Println(box, "\n")

	fmt.Println("Enjoy.")

	// Tile
	wb := textbox.NewTextbox(100, 50)
	wb.Tile(tbox, nil)
	fmt.Println(wb, "\n")

	// Crop and Copy
	/*	rb, _ := box.Crop(textbox.Rect{image.Pt(10,10),10,10})
		fmt.Println(rb, "\n")
		cb := rb.Copy()
		fmt.Println(cb, "\n")
		box.Tl(' ', 'x')
		fmt.Println(box, "\n")
		fmt.Println(rb, "\n")
		fmt.Println(cb, "\n")*/
}

var trogdor = []string{
	"                                                 :::",
	"                                             :: :::.",
	"                       \\/,                    .:::::",
	"           \\),          \\`-._                 :::888",
	"           /\\            \\   `-.             ::88888",
	"          /  \\            | .(                ::88",
	"         /,.  \\           ; ( `              .:8888",
	"            ), \\         / ;``               :::888",
	"           /_   \\     __/_(_                  :88",
	"             `. ,`..-'      `-._    \\  /      :8",
	"               )__ `.           `._ .\\/.",
	"              /   `. `             `-._______m         _,",
	"  ,-=====-.-;'                 ,  ___________/ _,-_,'\"`/__,-.",
	" C   =--   ;                   `.`._    V V V       -=-'\"#==-._",
	":,  \\     ,|      UuUu _,......__   `-.__A_A_ -. ._ ,--._ \",`` `-",
	"||  |`---' :    uUuUu,'          `'--...____/   `\" `\".   `",
	"|`  :       \\   UuUu:",
	":  /         \\   UuUu`-._",
	" \\(_          `._  uUuUu `-.",
	" (_3             `._  uUu   `._",
	"                    ``-._      `.",
	"                         `-._    `.",
	"                             `.    \\",
	"                               )   ;",
	"                              /   /",
	"               `.        |\\ ,'   /",
	"                 \",_A_/\\-| `   ,'",
	"                   `--..,_|_,-'\\",
	"                          |     \\",
	"                          |      \\__",
	"                          |__"}
