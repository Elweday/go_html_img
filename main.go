package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/chromedp/chromedp"
	"github.com/elweday/subtitles-go/helpers"
	"github.com/elweday/subtitles-go/styles"
)

const DIGITS = 5;


type Style map[string]func(words []helpers.Item, idx int, perc float64) string

var style = Style{
	"fading": styles.Fading,
	"scrolling-box": styles.ScrollingBox,
}

func RenderWord(words []helpers.Item, index int, output string, perc float64) {
	html := "`" + style["scrolling-box"](words, index, perc) +"`"

	cssBytes, _ := os.ReadFile("style.css")
	css := "`<style>" + string(cssBytes) + "</style>`" 

	ctx, cancel := chromedp.NewContext(context.Background())
    defer cancel()

    chromedp.Run(ctx, chromedp.EmulateViewport(1080, 400));

	
	/* if err := chromedp.Run(ctx, chromedp.Navigate("data:text/html,"+html)); err != nil {
		log.Fatal(err)
    } */


	jsScript := fmt.Sprintf(`document.documentElement.innerHTML = %s; document.head.innerHTML = %s;`, html, css)
	if err := chromedp.Run(ctx, chromedp.Evaluate(jsScript, nil)); err != nil {
		log.Fatal(err)
	}
 	var buf []byte
	if err := chromedp.Run(ctx, chromedp.CaptureScreenshot(&buf)); err != nil {
		log.Fatal(err)
	}



	if err := os.WriteFile(output, buf, 0644); err != nil {
		log.Fatal(err)
	}
}


func RenderWords(words []helpers.Item) {
	newpath := filepath.Join(".", "images")
	os.MkdirAll(newpath, os.ModePerm)


	var wg sync.WaitGroup
	count := 0;
	for i:=1; i<len(words); i++ {
		current := words[i]
		prev := words[i-1]
		frames :=  current.Frames - prev.Frames

		for j := 0; j < frames; j++ {
			count += 1
			output := fmt.Sprintf("./images/img_%0*d.png", DIGITS, count)
			perc := float64(j) / float64(frames)

			wg.Add(1)
			go func(index int, output string, perc float64) {
				defer wg.Done()
				RenderWord(words, index, output, perc)
			}(i, output, perc)
		}
	}
	wg.Wait()

}


func main() {
	fps := 15;
	words, err := helpers.ReadAndConvertToFrames("time_stamps.json", fps);
	if err != nil {
		panic(err)
	}
	RenderWords(words)
	
    cmd := exec.Command("ffmpeg", "-y", "-framerate", fmt.Sprint(fps), "-i", "./images/img_%05d.png", "-c:v", "libx264", "-r", "30", "-pix_fmt", "yuv420p", "output.mp4")

    err = cmd.Run()

    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Video created successfully!")
}
