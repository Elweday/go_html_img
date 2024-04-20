package styles

import (
	"fmt"

	"github.com/elweday/subtitles-go/helpers"
)
func Fading(words []helpers.Item, idx int, perc float64) string {
	html := "<p style='text-align: left;'>"
	for i, word := range words {
		if i < idx{
			html += word.Word + " "
		} else if i==idx  {
			html += fmt.Sprintf("<span style='opacity: %f; display: inline-block; transform: translate(0, %dpx);'>%s</span> ",
				perc,
				int((1-perc)*40),
				word.Word,
			)
		}
	}
	return html + "</p>"
}
