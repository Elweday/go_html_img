package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/elweday/subtitles-go/helpers"
	"github.com/elweday/subtitles-go/styles"
)

const chromaColor = "00ff00"
const DIGITS = 5;


type Style map[string]func(words []helpers.Item, idx int, perc float64) string

var style = Style{
	"fading": styles.Fading,
	"scrolling-box": styles.ScrollingBox,
}

func formatHtml(body string) string {
		html :=fmt.Sprintf( `
<html lang="en" style="background-color: #%s; width: 1080; height: 1920; padding: 0; margin: 0">
	<body>
		%s
	</body>
</html>
	`, chromaColor, body)
	return html


}

func RenderWord(words []helpers.Item, index int, output string, perc float64) {
	subtitlesStyle :=  "scrolling-box"
	body := "`" + style[subtitlesStyle](words, index, perc) +"`"

	
	html := formatHtml(body)

	imageBytes, err := htmlToImage(html)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile(output, imageBytes, 0644); err != nil {
		log.Fatal(err)
	}
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func RenderWords(words []helpers.Item) error {
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
			output := fmt.Sprintf("./images/img_%0*d.jpeg", DIGITS, count)
			duration := 50
			nframes := min(frames, duration)
			perc := float64(j) / float64(nframes)

			wg.Add(1)
			go func(index int, output string, perc float64) {
				defer wg.Done()
				RenderWord(words, index, output, perc)
			}(i, output, perc)
		}
	}
	wg.Wait()
	return nil

}


func main() {
	start := time.Now()
	defer func(){
		fmt.Printf("Video created successfully! %s", time.Since(start))
	}()

	fps := 15;
	words, err := helpers.ReadAndConvertToFrames("time_stamps.json", fps);
	if err != nil {
		panic(err)
	}
	err = RenderWords(words)
	   if err != nil {
        log.Fatal(err)
        return
    }
	fmt.Println("frames rendered ")
	
	cmd := exec.Command("ffmpeg",
		"-y",
        "-f", "image2",
        "-framerate", fmt.Sprint(fps), // Adjust frame rate as per your requirement
        "-i", "./images/img_%05d.jpeg", // Input PNG files path pattern
        "-c:v", "libx264",
        "-pix_fmt", "yuv420p",
        "-vf", fmt.Sprintf("colorkey=0x%s:0.1:0.1", chromaColor), // Scaling and Chroma key filter
        "output.mp4",
    )
	// cmd.Stderr = os.Stderr


    err = cmd.Run()

    if err != nil {
        log.Fatal(err)
        return
    }

    fmt.Println("Video created successfully!")
}


func htmlToImage(html string) ([]byte, error) {
    // Create a pipe
    pr, pw := io.Pipe()

    // Execute wkhtmltoimage command
    cmd := exec.Command("wkhtmltoimage", "-", "-")
    cmd.Stdin = pr // Set the pipe as stdin
    var out bytes.Buffer
    cmd.Stdout = &out

    // Start the command
    if err := cmd.Start(); err != nil {
        return nil, err
    }

    // Write the HTML content to the pipe asynchronously
    go func() {
        defer pw.Close()
        _, _ = io.WriteString(pw, html)
    }()

    // Wait for the command to finish
    if err := cmd.Wait(); err != nil {
        return nil, err
    }

    // Read the image data from the buffer
    imageBytes := out.Bytes()
    return imageBytes, nil
}
