package main

import (
  "log"
  "time"
  "os"
  "os/exec"
  "strconv"
  "strings"
  "math/rand"
  "fmt"
  "path/filepath"
)

func main() {
  args := os.Args[1:]
  var period int64 = 30
  usedImagesDict := make(map[string]string)

  if len(args) >= 2 {
    num, err := strconv.ParseInt(args[1], 10, 64)

    if err != nil {
      log.Printf("Error converting argument into number %v", err)
      period = 30
    } else {
      period = num
    }
  } else {
    panic("Few arguments. You must specify [images-directory] [period]")
  }

  log.Printf("The images will be changed each %v seconds", period)

  t := time.NewTicker(time.Second * time.Duration(period))

  for {
    dir := strings.TrimSpace(args[0])
    images, err := getDirectoryImages(dir)

    if (err != nil) {
      log.Panic(err)
    } else {
      background, newDict := getRandomImage(images, usedImagesDict)

      usedImagesDict = newDict

      changeBackgoundImage(background)
    }

    <-t.C
  }
}

func getDirectoryImages(dir string) ([]string, error) {
  return filepath.Glob( fmt.Sprintf("%v/*.jpg", dir) )
}

func changeBackgoundImage(image string) {
  _, err := exec.Command("/usr/bin/gsettings", "set", "org.gnome.desktop.background", "picture-uri", fmt.Sprintf("\"file://%s\"", image)).Output()

  if (err != nil) {
    log.Printf("Error changing background image %s", image)
    log.Println(err)
  }
}

func getRandomImage(imagesArray []string, usedImagesDict map[string]string) (string, map[string]string) {
  background := strings.TrimSpace(imagesArray[ rand.Intn( len(imagesArray) ) ])
  ok := false

  if (len(imagesArray) == len(usedImagesDict)) {
    usedImagesDict = make(map[string]string)
  }

  _, ok = usedImagesDict[background]

  if (!ok) {
    usedImagesDict[background] = background

    return background, usedImagesDict
  } else {
    return getRandomImage(imagesArray, usedImagesDict)
  }

  return background, usedImagesDict
}
