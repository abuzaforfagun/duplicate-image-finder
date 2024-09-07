package main

type ImageFile struct {
	Name string
	Type ImageType
}

type ImageType int

const (
	JPG ImageType = iota
	PNG
)

func NewJpegImage(name string) *ImageFile {
	return &ImageFile{
		Name: name,
		Type: JPG,
	}
}

func NewPngImage(name string) *ImageFile {
	return &ImageFile{
		Name: name,
		Type: PNG,
	}
}
