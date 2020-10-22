package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/Comdex/imgo"
	"github.com/disintegration/imaging"
	"github.com/gordonklaus/portaudio"
	"github.com/morikuni/aec"
	"github.com/tosone/minimp3"
	"gocv.io/x/gocv"
)

var startNano int64
var mp4 = "q"
var basePath string

func init() {
	flag.StringVar(&mp4, "m", "q", "m")
	flag.Parse()
	dir, _ := os.UserHomeDir()
	basePath = path.Join(dir, mp4)
}

func main() {
	startNano = time.Now().UnixNano()
	var group sync.WaitGroup
	group.Add(2)
	go func() {
		_ = paintingIO(basePath)
		fmt.Println("the end")
		group.Done()
	}()
	go func() {
		_ = play(path.Join(basePath, mp4+".mp3"))
		group.Done()
	}()
	group.Wait()
	os.Exit(0)
}

func painting(filename string, up int) error {
	vc, err := gocv.VideoCaptureFile(filename)
	if err != nil {
		panic(err)
	}
	fps := vc.Get(gocv.VideoCaptureFPS)
	frames := vc.Get(gocv.VideoCaptureFrameCount)
	up2 := aec.Up(uint(up))
	// for up
	for i := 0; i < up; i++ {
		fmt.Println()
	}
	for ; ; {
		//vc.Set(gocv.VideoCapturePosFrames, float64(i))
		//img := gocv.NewMat()
		reNow := time.Now().UnixNano() - startNano
		srcImage, err := GetVideoMoment(vc, reNow, frames, fps)
		if err != nil {
			return err
		}
		reImg := imaging.Resize(srcImage, 0, up, imaging.Lanczos)
		height := imgo.GetImageHeight(reImg) // 获取 图片 高度[height]
		width := imgo.GetImageWidth(reImg)   // 获取 图片 宽度[width]
		imgMatrix, _ := imgo.Read(reImg)     // 读取图片RGBA值
		fmt.Print(up2)
		for starty := 0; starty < height; starty++ {
			for startx := 0; startx < width; startx++ {
				R := imgMatrix[starty][startx][0] // 读取图片R值
				G := imgMatrix[starty][startx][1] // 读取图片G值
				B := imgMatrix[starty][startx][2] // 读取图片B值
				fmt.Print(aec.Color8BitB(aec.NewRGB8Bit(uint8(R), uint8(G), uint8(B))).Apply(" "))
			}
			fmt.Println()
		}
	}
}

func paintingIO(basePath string) error {
	json0 := map[string]interface{}{
		//"height": height,
		//"width":  width,
		//"frames": frames,
		//"fps":    fps,
	}
	file, err := ioutil.ReadFile(path.Join(basePath, "0.json"))
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(file, &json0)
	if err != nil {
		panic(err)
	}
	up := int(json0["height"].(float64))
	frames := json0["frames"].(float64)
	fps := json0["fps"].(float64)

	up2 := aec.Up(uint(up))
	// for up
	for i := 0; i < up; i++ {
		fmt.Println()
	}
	for ; ; {
		//vc.Set(gocv.VideoCapturePosFrames, float64(i))
		//img := gocv.NewMat()
		reNow := time.Now().UnixNano() - startNano
		duration := frames / fps // 一共多少S
		timeS := float64(reNow) / 1000.0 / 1000.0 / 1000.0
		frames = (timeS / duration) * frames
		if int(frames) == 0 {
			continue
		}

		file, err := ioutil.ReadFile(path.Join(basePath, fmt.Sprintf("%d.json", int(frames))))
		if err != nil {
			return err
		}
		var imgMatrix [][][3]int
		err = json.Unmarshal(file, &imgMatrix)
		if err != nil {
			return err
		}
		fmt.Print(up2)
		for i := range imgMatrix {
			for j := range imgMatrix[i] {
				R := imgMatrix[i][j][0]
				G := imgMatrix[i][j][1]
				B := imgMatrix[i][j][2]
				fmt.Print(aec.Color8BitB(aec.NewRGB8Bit(uint8(R), uint8(G), uint8(B))).Apply(" "))
			}
			fmt.Println()
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func GetVideoMoment(vc *gocv.VideoCapture, timeNano int64, frames, fps float64) (i image.Image, err error) {
	duration := frames / fps // 一共多少S
	timeS := float64(timeNano) / 1000.0 / 1000.0 / 1000.0
	frames = (timeS / duration) * frames
	// Set Video frames
	vc.Set(gocv.VideoCapturePosFrames, frames)
	img := gocv.NewMat()
	vc.Read(&img)
	imageObject, err := img.ToImage()
	if err != nil {
		return i, err
	}
	return imageObject, err
}

func play(filename string) error {
	var (
		err  error
		data []byte
		dec  *minimp3.Decoder
	)
	err = portaudio.Initialize()
	if err != nil {
		return err
	}
	//读取文件
	bs, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	//用minimp3 解析 音乐
	if dec, data, err = minimp3.DecodeFull(bs); err != nil {
		return err
	}
	//当退出函数时，释放资源
	defer dec.Close()

	// 定义缓冲，portaudio Write 从缓冲里面获取数据。
	out := make([]int16, 8196)
	//用 portaudio 打开音频流
	stream, err := portaudio.OpenDefaultStream(
		0,
		dec.Channels,
		float64(dec.SampleRate),
		len(out),
		&out,
	)
	if err != nil {
		return err
	}
	//当退出函数，释放资源
	defer stream.Close()
	//开始视频处理
	if err := stream.Start(); err != nil {
		return err
	}
	// 退出函数时，释放资源
	defer stream.Stop()
	byteffer := bytes.NewBuffer(data)

	//循环 来读取minimp3解析出来的数据，并给portalaudio处理
	for {
		audio := make([]byte, 2*len(out))
		_, err := byteffer.Read(audio)
		// 有错误，就退出。io.EOF 说明文件已经结束
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		err = binary.Read(bytes.NewBuffer(audio), binary.LittleEndian, out)
		// 有错误，就退出。io.EOF 说明文件已经结束
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		if err := stream.Write(); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
