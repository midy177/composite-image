package main

import (
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	wg          sync.WaitGroup
	factoryList = make(chan struct{}, 1000)
)

type FileList [][]interface{}

func main() {
	for i := 0; i < 10; i++ {
		factoryList <- struct{}{}
	}
	GetDir()
	wg.Wait()
}

func Product(sets ...[]interface{}) {
	lens := func(i int) int { return len(sets[i]) }
	//product := [][]interface{}{}
	for ix := make([]int, len(sets)); ix[0] < lens(0); nextIndex(ix, lens) {
		var r []interface{}
		for j, k := range ix {
			r = append(r, sets[j][k])
		}
		<-factoryList
		go ComImg(r)
	}
}

func nextIndex(ix []int, lens func(i int) int) {
	for j := len(ix) - 1; j >= 0; j-- {
		ix[j]++
		if j == 0 || ix[j] < lens(j) {
			return
		}
		ix[j] = 0
	}
}

func GetDir() {
	var tmpList = make(FileList, 7)
	filepath.Walk("./加密骑士/", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			switch strings.Split(path, "/")[1] {
			case "01背景":
				tmpList[0] = append(tmpList[0], path)
			case "02身体":
				tmpList[1] = append(tmpList[1], path)
			case "03眼睛":
				tmpList[2] = append(tmpList[2], path)
			case "04嘴巴":
				tmpList[3] = append(tmpList[3], path)
			case "05帽子":
				tmpList[4] = append(tmpList[4], path)
			case "06衣服":
				tmpList[5] = append(tmpList[5], path)
			case "07特殊配饰":
				tmpList[6] = append(tmpList[6], path)
			}

		}
		return nil
	})
	Product(tmpList[0], tmpList[1], tmpList[2], tmpList[3], tmpList[4], tmpList[5], tmpList[6])
}

func ComImg(item []interface{}) {
	wg.Add(1)
	if len(item) == 7 {
		//背景
		backgroundFile, _ := os.Open(item[0].(string))
		imgBackground, _ := png.Decode(backgroundFile)
		defer backgroundFile.Close()
		//身体
		bodyFile, _ := os.Open(item[1].(string))
		imgBody, _ := png.Decode(bodyFile)
		defer bodyFile.Close()
		//眼睛
		eyeFile, _ := os.Open(item[2].(string))
		imgEye, _ := png.Decode(eyeFile)
		defer eyeFile.Close()
		//mouth
		mouthFile, _ := os.Open(item[3].(string))
		imgMouth, _ := png.Decode(mouthFile)
		defer mouthFile.Close()
		//hat
		hatFile, _ := os.Open(item[4].(string))
		imgHat, _ := png.Decode(hatFile)
		//
		clothesFile, _ := os.Open(item[5].(string))
		imgClothes, _ := png.Decode(clothesFile)
		//special clothes
		specialClothesFile, _ := os.Open(item[6].(string))
		imgSpecialClothes, _ := png.Decode(specialClothesFile)

		offset := image.Pt(0, 0)
		m := image.NewRGBA(imgBody.Bounds())

		draw.Draw(m, imgBackground.Bounds().Add(offset), imgBackground, image.ZP, draw.Src)
		draw.Draw(m, imgBody.Bounds().Add(offset), imgBody, image.ZP, draw.Over)
		draw.Draw(m, imgEye.Bounds().Add(offset), imgEye, image.ZP, draw.Over)
		draw.Draw(m, imgMouth.Bounds().Add(offset), imgMouth, image.ZP, draw.Over)
		draw.Draw(m, imgHat.Bounds().Add(offset), imgHat, image.ZP, draw.Over)
		draw.Draw(m, imgClothes.Bounds().Add(offset), imgClothes, image.ZP, draw.Over)
		draw.Draw(m, imgSpecialClothes.Bounds().Add(offset), imgSpecialClothes, image.ZP, draw.Over)
		var newImgName string
		for k, v := range item {
			tmpName := strings.SplitAfter(v.(string), "/")
			lens := len(tmpName)
			if k != 6 {
				newImgName += strings.Split(tmpName[lens-1], ".")[0] + "_"
			} else {
				newImgName += tmpName[lens-1]
			}
		}
		imgCom, _ := os.Create("tmp/" + newImgName)
		jpeg.Encode(imgCom, m, &jpeg.Options{jpeg.DefaultQuality})
		fmt.Printf("save img: %s\n", newImgName)
		defer imgCom.Close()
	} else {
		fmt.Print("item 长度错误.\n")
	}
	factoryList <- struct{}{}
	wg.Done()
}

// 图片大小调整
func ImageResize(src image.Image, w, h int) image.Image {
	return resize.Resize(uint(w), uint(h), src, resize.Lanczos3)
}
