package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type ImagesRepository interface {
	SaveImage(data []byte, filename string) error
	RemoveImage(imageName string) error
	GetPath() string
}

type imagesRepository struct {
	imagePath string
}

func NewImagesRepository(imagePath string) ImagesRepository {
	return &imagesRepository{
		imagePath: imagePath,
	}
}

func (r *imagesRepository) SaveImage(data []byte, filename string) error {
	err := os.MkdirAll(r.imagePath, os.ModePerm)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(r.imagePath, filename), data, 0644)
}

func (r *imagesRepository) RemoveImage(imageName string) error {

	imagePath := filepath.Join(r.imagePath, imageName)
	fmt.Println("remove image ", imagePath)
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		fmt.Println(fmt.Errorf("file does not exist: %s", imageName))
		return fmt.Errorf("file does not exist: %s", imageName)
	}

	err := os.Remove(imagePath)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to remove file: %w", err))
		return fmt.Errorf("failed to remove file: %w", err)
	}
	fmt.Println("success remove image")
	return nil
}

func (r *imagesRepository) GetPath() string {
	return r.imagePath
}
