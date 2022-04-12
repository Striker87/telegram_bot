package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Striker87/telegram_bot/storage"
)

const defPerm = 0774

type Storage struct {
	basePath string
}

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) error {
	fPath := filepath.Join(s.basePath, page.UserName) // определяем путь до директории в которую будет сохранятся файл

	err := os.MkdirAll(fPath, defPerm) // создаем все нужные директории по этому пути
	if err != nil {
		return fmt.Errorf("failed to create dir %s due error: %v", fPath, err)
	}

	fName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("failed to get fileName %v due error: %v", page, err)
	}

	fPath = filepath.Join(fPath, fName)
	file, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s due error: %v", file.Name(), err)
	}
	defer file.Close()

	err = gob.NewEncoder(file).Encode(page)
	if err != nil {
		return fmt.Errorf("failed to encode file %s via gob due error: %v", file.Name(), err)
	}
	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to readFile %s due error: %v", path, err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNotSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("failed generate file %v due error: %v", p, err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)
	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("failed remove file %s due error: %v", path, err)
	}
	return nil
}

func (s Storage) IsExists(p *storage.Page) (bool, error) { // существует данная страница или нет, сохранял ли ее пользователь ранее?
	fileName, err := fileName(p)
	if err != nil {
		return false, fmt.Errorf("failed generate file %v due error: %v", p, err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("failed to check if file %s exists due error: %v", path, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s due error: %v", filePath, err)
	}
	defer f.Close()

	var p storage.Page
	err = gob.NewDecoder(f).Decode(&p)
	if err != nil {
		return nil, fmt.Errorf("")
	}
	return &p, nil
}
func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
