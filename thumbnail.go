package thumbnail

import (
	"github.com/gographics/imagick/imagick"
	"strings"
)

func Thumbnail(filename string, targetWidth, targetHeight uint) {
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	var err error

	// 画像読み込み
	err = mw.ReadImage(filename)
	if err != nil {
		panic(err)
	}

	// 縦横サイズを取得
	width, height := mw.GetImageWidth(), mw.GetImageHeight()

	// 縦横の小さい方に合わせてリサイズしたときのサイズを取得
	resizedWidth, resizedHeight := getResizedWH(width, height, targetWidth, targetHeight)

	// サムネイルを作成
	err = mw.ThumbnailImage(resizedWidth, resizedHeight)
	if err != nil {
		panic(err)
	}

	// 切り抜き開始位置
	startX, startY := getStartPointXY(width, height, resizedWidth, resizedHeight)

	// 最終的なサイズで切り抜く
	err = mw.ExtentImage(targetWidth, targetHeight, startX, startY)
	if err != nil {
		panic(err)
	}

	// 画像のクオリティ
	err = mw.SetImageCompressionQuality(95)
	if err != nil {
		panic(err)
	}

	// ファイルに保存
	// TODO(hfunai): ファイル名にドットがあると死亡
	err = mw.WriteImage(strings.Replace(filename, ".", "_thumb.", 1))
	if err != nil {
		panic(err)
	}
}

func getResizedWH(width, height, targetWidth, targetHeight uint) (resizedWidth, resizedHeight uint) {
	if width < height {
		ratio := float32(width) / float32(height)
		targetHeight = uint(float32(targetWidth) / ratio)
	} else {
		ratio := float32(height) / float32(width)
		targetWidth = uint(float32(targetHeight) / ratio)
	}
	return targetWidth, targetHeight
}

func getStartPointXY(width, height, resizedWidth, resizedHeight uint) (x, y int) {
	startX, startY := 0, 0
	if width < height {
		startY = int((float32(resizedHeight) - float32(resizedWidth)) / 2.0)
	} else {
		startX = int((float32(resizedWidth) - float32(resizedHeight)) / 2.0)
	}
	return startX, startY
}
