package styles

import (
	"fmt"

	"github.com/elweday/subtitles-go/helpers"
)
func ScrollingBox(words []helpers.Item, idx int, perc float64) string {
	scale := helpers.Spring( 0.92, 1.05, helpers.SpringOptions{ Stiffness: 100, Damping: 10, Mass: 1 } )
	html := "<p>"
	for i, word := range words {
		if i != idx {
			html += word.Word + " "
		} else {
			html += fmt.Sprintf("<span style='scale: %f;' class='selected'>%s</span> ", scale(perc), word.Word)
		}
	}
	return html + "</p>"
}
