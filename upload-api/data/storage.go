package data

import "io"

type Storage interface {
	Save(filePath string, body io.Reader) error
}
