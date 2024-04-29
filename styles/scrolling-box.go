package styles

import (
	"fmt"

	"github.com/elweday/subtitles-go/helpers"
)
func ScrollingBox(words []helpers.Item, idx int, perc float64) string {
	html := "<p style='text-align: left; color: black;  font-size: 80px; padding: 0px; '>"
	for i, word := range words {
		if i != idx {
			html += word.Word + " "
		} else {
			html += fmt.Sprintf(
				`<span style='scale: %f; background: blue; padding: 50px; border-radius: 10px; color: white; display: inline-block; height: fit-content;' >
					%s
				</span> `,
				perc,
				word.Word)
		}
	}
	return html + "</p>"
}
