package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

func main() {
	//GenerateGallery("./../../img/catalina/")
	//GenerateGallery("./../../img/yacht/")
	GenerateGallery("./../../../micro_site_host/public/images/jdboatcovers/20180404/", "micro_site_host/public/")
}

// to do - convert to a stand along package

func findFilesToProcess(rootDirectory string, suffix string) []os.FileInfo {
	result := []os.FileInfo{}

	files, _ := ioutil.ReadDir(rootDirectory)
	for _, f := range files {
		log.Println(f.Name())
		if strings.HasSuffix(strings.ToLower(f.Name()), ".jpg") || strings.HasSuffix(strings.ToLower(f.Name()), ".png") {
			if !strings.HasSuffix(strings.ToLower(f.Name()), suffix+".jpg") && !strings.HasSuffix(strings.ToLower(f.Name()), suffix+".png") {

				result = append(result, f)
			}
		}
	}
	return result
}

func ResizePath(rootDirectory string, suffix string, maxSize uint) {
	// loop through all the files in the path
	files := findFilesToProcess(rootDirectory, suffix)
	for _, f := range files {
		fmt.Println(f.Name())
		resizedFileName := Resize(rootDirectory+f.Name(), "", suffix, maxSize)
		fmt.Println(resizedFileName)
		Base64(resizedFileName)
	}

}

func GenerateGallery(rootDirectory string, removePath string) {
	galleryPath := rootDirectory + "gallery/"
	if _, err := os.Stat(galleryPath); os.IsNotExist(err) {
		os.Mkdir(galleryPath, os.ModePerm)
	}

	files := findFilesToProcess(rootDirectory, "_web")

	galleryJson, err := os.Create(galleryPath + "gallery.json")
	if err != nil {
		log.Fatal(err)
	}
	defer galleryJson.Close()
	galleryJson.WriteString("[\n")
	for i, f := range files {
		fmt.Println(f.Name())
		if i == 0 {
			galleryJson.WriteString("\t{\n")
		} else {
			galleryJson.WriteString(",\n\t{\n")
		}
		galleryJson.WriteString("\t\t\"name\": \"" + f.Name() + "\",\n")
		// base 64 and micro
		resizedBase64 := ResizeTo(rootDirectory+f.Name(), galleryPath+"micro_"+f.Name(), 6)
		base64 := Base64(resizedBase64)
		galleryJson.WriteString("\t\t\"data\": \"" + base64 + "\",\n")
		galleryJson.WriteString("\t\t\"micro\": \"" + WebRoot(resizedBase64, removePath) + "\",\n")
		// web size
		resizedWeb := ResizeTo(rootDirectory+f.Name(), galleryPath+"web_"+f.Name(), 700)
		galleryJson.WriteString("\t\t\"web\": \"" + WebRoot(resizedWeb, removePath) + "\",\n")
		// mid size
		resizedMid := ResizeTo(rootDirectory+f.Name(), galleryPath+"mid_"+f.Name(), 300)
		galleryJson.WriteString("\t\t\"mid\": \"" + WebRoot(resizedMid, removePath) + "\",\n")
		// thumbnail size
		resizedThumb := ResizeTo(rootDirectory+f.Name(), galleryPath+"thumb_"+f.Name(), 100)
		galleryJson.WriteString("\t\t\"thumb\": \"" + WebRoot(resizedThumb, removePath) + "\"\n")
		// ~end
		galleryJson.WriteString("\t}")
	}
	galleryJson.WriteString("\n]")
}

func WebRoot(directory string, removePath string) string {
	result := strings.Replace(directory, "../", "", -1)
	result = strings.Replace(result, "./", "/", -1)
	result = strings.Replace(result, removePath, "", -1)
	return result
}

func Swancraft(rootDirectory string) {
	// loop through all the files in the path
	files := findFilesToProcess(rootDirectory, "_web")
	for _, f := range files {
		fmt.Println(f.Name())
		resizedFileName := Resize(rootDirectory+f.Name(), "", "_web", 700)
		fmt.Println(resizedFileName)
		resizedFileName4 := Resize(rootDirectory+f.Name(), "", "_mid", 300)
		fmt.Println(resizedFileName4)
		resizedFileName2 := Resize(rootDirectory+f.Name(), "", "_thumb", 100)
		fmt.Println(resizedFileName2)
		//resizedFileName3 := Resize(rootDirectory + f.Name(), "", "_data", 6)
		//Base64(resizedFileName3)
	}

}

func openImage(fileName string) (image.Image, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(strings.ToLower(fileName), ".jpg") {
		img, err := jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
		file.Close()
		return img, nil
	}
	if strings.HasSuffix(strings.ToLower(fileName), ".png") {
		img, err := png.Decode(file)
		if err != nil {
			return nil, err
		}
		file.Close()
		return img, nil
	}
	return nil, fmt.Errorf("Format not supported: %s", fileName)
}

func Resize(fileName string, folder string, suffix string, maxSize uint) string {
	img, err := openImage(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// ensure the largest dimension is scaled
	var imgWidth uint = maxSize
	var imgHeight uint = 0
	b := img.Bounds()
	if b.Max.Y > b.Max.X {
		log.Printf("image %s is portrait", fileName)
		imgWidth = 0
		imgHeight = maxSize
	}

	// resize using Lanczos resampling and preserve aspect ratio
	m := resize.Resize(imgWidth, imgHeight, img, resize.Lanczos3)

	resizedFileName := strings.Replace(strings.ToLower(fileName), ".png", suffix+".png", -1)
	resizedFileName = strings.Replace(strings.ToLower(fileName), ".jpg", suffix+".jpg", -1)
	resizedFileName = folder + resizedFileName

	out, err := os.Create(resizedFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	png.Encode(out, m)

	return resizedFileName
}

func ResizeTo(fileName string, toFileName string, maxSize uint) string {
	img, err := openImage(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// ensure the largest dimension is scaled
	var imgWidth uint = maxSize
	var imgHeight uint = 0
	b := img.Bounds()
	if b.Max.Y > b.Max.X {
		log.Printf("image %s is portrait", fileName)
		imgWidth = 0
		imgHeight = maxSize
	}

	// resize using Lanczos resampling and preserve aspect ratio
	m := resize.Resize(imgWidth, imgHeight, img, resize.Lanczos3)

	out, err := os.Create(toFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	png.Encode(out, m)

	return toFileName
}

func Base64(fileName string) string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	//log.Printf(file.Name())

	defer file.Close()

	// get the new file size
	fInfo, _ := file.Stat()
	var size int64 = fInfo.Size()
	buf := make([]byte, size)

	// read the content into the buffer
	fReader := bufio.NewReader(file)
	fReader.Read(buf)

	imageBase64String := base64.StdEncoding.EncodeToString(buf)

	dataSourceString := "data:image/png;base64," + imageBase64String

	return dataSourceString
}
