package manga

// Uploads provided preview and returns url and filename
// func (m *Manga) uploadPreview(ctx context.Context, mangaID uuid.UUID, preview models.UploadFile) (string, string, error) {
// 	previewBucket := filepath.Join(filestorage.MangaBucket, mangaID.String())
// 	filename, err := m.fileStorage.SaveFile(
// 		previewBucket,
// 		"preview"+filepath.Ext(preview.Filename),
// 		preview.Data,
// 	)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	previewURL, err := m.fileStorage.GetFileURL(previewBucket, filename)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	return previewURL, filename, nil
// }

// func getMangaBucket(mangaID uuid.UUID) string {
// 	return filepath.Join(filestorage.MangaBucket, mangaID.String())
// }
