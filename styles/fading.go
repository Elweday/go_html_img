package styles

import (
	"github.com/elweday/subtitles-go/types"
)

var fadingTemplate = `
<p style='text-align: left; color: black; font-size: 80px; padding: 0px; display: flex; flex-wrap: wrap; '>
{{ range $index, $word := .Words }}
	{{ if lt $index $.Index }}
		{{ $word.Word }} 
	{{ else if eq $index $.Index }}
		&nbsp;
		<span style='position: relative; display: inline-block; padding: 0; white-space: nowrap;'>
			<span style='opacity: {{ printf "%.2f" $.Percentage }}; position: absolute; top: {{ offset $.Percentage }}px;'>{{ $word.Word }}</span>
		</span>
	{{ end }}
{{ end }}
</p>
`

var Fading = types.NewStyle("fading", fadingTemplate, map[string]any{
		"offset": func (perc float64) int {
			 return int((1-perc) * 40)
		},
})
