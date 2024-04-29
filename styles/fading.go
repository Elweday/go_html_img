package styles

import (
	"fmt"

	"github.com/elweday/subtitles-go/helpers"
)
func Fading(words []helpers.Item, idx int, perc float64) string {
	html := "<p style='text-align: left; color: black;  font-size: 80px; padding: 0px; '>"
	for i, word := range words {
		if i < idx{
			html += word.Word + " "
		} else if i==idx  {
			html += fmt.Sprintf(`
			&nbsp;
			<span style='position: relative; padding: 0;'>
				<span style='opacity: %f; position: absolute; top: %dpx;'>%s</span>
			</span>
			`,
				perc,
				int((1-perc)*40),
				word.Word,
			)
		}
	}
	return html + "</p>"
}

