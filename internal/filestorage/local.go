package filestorage

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/lzaxel/zero-manga-backend/internal/logger"
)

type LocalFileStorageConfig struct {
	StoragePath string `yaml:"storagePath" env:"LOCAL_FILE_STORAGE_PATH"`
	Host        string `yaml:"host" env:"LOCAL_FILE_STORAGE_HOST"`
	Port        uint   `yaml:"port" env:"LOCAL_FILE_STORAGE_PORT"`
}

type LocalFileStorage struct {
	logger      logger.Logger
	storagePath string
	host        string
	port        uint
}

func NewLocalFileStorage(config LocalFileStorageConfig, logger logger.Logger) *LocalFileStorage {
	return &LocalFileStorage{
		storagePath: config.StoragePath,
		host:        config.Host,
		port:        config.Port,
		logger:      logger,
	}
}

func (s *LocalFileStorage) SaveFile(bucket string, filename string, data []byte) (uuid.UUID, error) {
	id := uuid.New()
	fileExtansion := filepath.Ext(filename)

	filePath := filepath.Join(s.storagePath, bucket, id.String()+fileExtansion)
	os.MkdirAll(filepath.Dir(filePath), os.ModePerm)

	file, err := os.Create(filePath)
	if err != nil {
		return uuid.Nil, fmt.Errorf("LocalFileStorage.SaveFile: creating file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return uuid.Nil, fmt.Errorf("LocalFileStorage.SaveFile: writing file: %w", err)
	}

	return id, nil
}

func (s *LocalFileStorage) GetFile(bucket string, id uuid.UUID) ([]byte, error) {
	filepath, err := s.findFilePathByUUID(bucket, id)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("LocalFileStorage.GetFile: reading file: %w", err)
	}

	return data, nil
}

func (s *LocalFileStorage) DeleteFile(bucket string, id uuid.UUID) error {
	filepath, err := s.findFilePathByUUID(bucket, id)
	if err != nil {
		return err
	}

	if err := os.Remove(filepath); err != nil {
		return fmt.Errorf("LocalFileStorage.DeleteFile: removing file: %w", err)
	}

	return nil
}

func (s *LocalFileStorage) GetFilePath(bucket string, id uuid.UUID) (string, error) {
	filepath, err := s.findFilePathByUUID(bucket, id)
	if err != nil {
		return "", fmt.Errorf("LocalFileStorage.GetFileURI: %w", err)
	}

	return filepath, nil
}

func (s *LocalFileStorage) GetFileURL(bucket string, id uuid.UUID) (string, error) {
	filepath, err := s.findFilePathByUUID(bucket, id)
	if err != nil {
		return "", fmt.Errorf("LocalFileStorage.GetFileURL: %w", err)
	}

	fileurl := fmt.Sprintf("http://%s:%d/", s.host, s.port) + path.Join(bucket, path.Base(filepath))

	return fileurl, nil
}

func (s *LocalFileStorage) findFilePathByUUID(bucket string, id uuid.UUID) (string, error) {
	pattern := filepath.Join(s.storagePath, bucket, id.String()) + ".*"
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("LocalFileStorage.findFilePathByUUID: %w: paths=%s", err, paths)
	}
	if len(paths) != 1 {
		return "", fmt.Errorf("LocalFileStorage.findFilePathByUUID: %w: pattern=%s", ErrNotFound, pattern)
	}

	return paths[0], nil
}

func (s *LocalFileStorage) Serve() {
	s.logger.Infof("starting local file storage server on %s:%d (path: %s)", s.host, s.port, s.storagePath)
	http.Handle("/", http.FileServer(http.Dir(s.storagePath)))
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port),
		nil,
	)
}
