package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/elweday/subtitles-go/helpers"
	"github.com/elweday/subtitles-go/styles"
	"github.com/elweday/subtitles-go/types"
)

const chromaColor = "00ff00"
const DIGITS = 5;


type Style map[string]func(args types.Args) string

var style = Style{
	"fading": styles.Fading,
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

func RenderWord(words []helpers.Item, index int, perc float64) []byte {
	subtitlesStyle :=  "fading"
	args := types.Args{
		Words: words,
		Index: index,
		Percentage: perc,
	}
	body := style[subtitlesStyle](args)
	
	html := formatHtml(body)

	imageBytes, err := htmlToImage(html)
	if err != nil {
		log.Fatal(err)
	}

/* 	if err := os.WriteFile(output, imageBytes, 0644); err != nil {
		log.Fatal(err)
	}
 */
	return imageBytes
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

type Update struct {
	Index int
	Data  []byte
}

func RenderWords(words []helpers.Item, shared *TypedMap[int, []byte]) (int, error) {
	var wg sync.WaitGroup
	count := 0;
	for i:=1; i<len(words); i++ {
		current := words[i]
		prev := words[i-1]
		frames :=  current.Frames - prev.Frames

		for j := 0; j < frames; j++ {
			count += 1
			// output := fmt.Sprintf("./images/img_%0*d.jpeg", DIGITS, count)
			duration := 50
			nframes := min(frames, duration)
			p := float64(j) / float64(nframes)

			wg.Add(1)
			go func(index int, perc float64) {
				defer wg.Done()
				imageBytes := RenderWord(words, index, perc)
				shared.Store(index, imageBytes)
			}(i, p)
		}
	}
	wg.Wait()
	return count, nil

}


func main() {
	start := time.Now()
	defer func(){
		fmt.Printf("Video created successfully! %s", time.Since(start))
	}()

	fps := 30;
	words, err := helpers.ReadAndConvertToFrames("time_stamps.json", fps);
	if err != nil {
		panic(err)
	}
	fmt.Println("STAGE 0 done ")
	sharedMap := TypedMap[int, []byte]{}

	count, _ := RenderWords(words, &sharedMap)

	cmd := exec.Command("ffmpeg",
		"-y",
		"-f", "rawvideo",
		"-pixel_format", "rgb24",
		"-video_size", "1080x1920",
    	"-framerate", fmt.Sprint(fps), // Adjust frame rate as per your requirement
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"output.mp4",
	)

    pipeReader, pipeWriter := io.Pipe()

    cmd.Stdin = pipeReader
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr


    // Close the writer once ffmpeg is done writing
    defer func() {
        pipeWriter.Close()
        cmd.Wait()
    }()


	for i :=0; i<count; i++ {
		b, ok := sharedMap.Load(i)
	
		if ok  {
			pipeWriter.Write(b);
		}
	}

	cmd.Run()

}


func htmlToImage(html string) ([]byte, error) {
    pr, pw := io.Pipe()
	
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


