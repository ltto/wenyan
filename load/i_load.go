package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/Comdex/imgo"
	"github.com/disintegration/imaging"
	"gocv.io/x/gocv"
)

var save0 = false
var loadF = "sw"

func main() {
	dir, _ := os.UserHomeDir()
	var mvPath = path.Join(dir, loadF)
	_ = os.Remove(loadF + ".mp3")
	if err := exec.Command("ffmpeg", "-i", loadF+".mp4", "-f", "mp3", "-vn", path.Join(mvPath, loadF+".mp3")).Run(); err != nil {
		panic(err)
	}
	vc, err := gocv.VideoCaptureFile(loadF + ".mp4")
	if err != nil {
		panic(err)
	}
	fps := vc.Get(gocv.VideoCaptureFPS)
	frames := uint64(vc.Get(gocv.VideoCaptureFrameCount))
	up := 80

	_ = os.Mkdir(path.Clean(mvPath), 0777)
	for i := uint64(1); i <= frames; i++ {
		//vc.Set(gocv.VideoCapturePosFrames, float64(i))
		//img := gocv.NewMat()
		vc.Set(gocv.VideoCapturePosFrames, float64(i))
		img := gocv.NewMat()
		vc.Read(&img)
		srcImage, err := img.ToImage()
		if err != nil {
			panic(err)
		}
		reImg := imaging.Resize(srcImage, 0, int(up), imaging.Lanczos)
		//reImg = imaging.Grayscale(reImg)
		height := imgo.GetImageHeight(reImg) // 获取 图片 高度[height]
		width := imgo.GetImageWidth(reImg)   // 获取 图片 宽度[width]
		imgMatrix, _ := imgo.Read(reImg)     // 读取图片RGBA值
		imgMap := make([][][3]int, height)

		if !save0 {
			json0 := map[string]interface{}{
				"height": height,
				"width":  width,
				"frames": frames,
				"fps":    fps,
			}
			bytes, err := json.Marshal(json0)
			if err != nil {
				panic(err)
			}
			err = ioutil.WriteFile(path.Join(mvPath, fmt.Sprintf("%d.json", 0)), bytes, 0777)
			if err != nil {
				panic(err)
			}
			save0 = true
		}

		for starty := 0; starty < height; starty++ {
			imgMap[starty] = make([][3]int, width)
			for startx := 0; startx < width; startx++ {
				R := imgMatrix[starty][startx][0] // 读取图片R值
				G := imgMatrix[starty][startx][1] // 读取图片G值
				B := imgMatrix[starty][startx][2] // 读取图片B值
				newR := int(R)
				imgMap[starty][startx][0] = newR
				newG := int(G)
				imgMap[starty][startx][1] = newG
				newB := int(B)
				imgMap[starty][startx][2] = newB
			}
		}
		marshal, err := json.Marshal(imgMap)
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(path.Join(mvPath, fmt.Sprintf("%d.json", i)), marshal, 0777)
		if err != nil {
			panic(err)
		}
		fmt.Println(i, "/", frames)
	}
}
