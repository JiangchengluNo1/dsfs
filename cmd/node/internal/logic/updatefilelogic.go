package logic

func UpdateFile(filemeta string, up_data []byte) error {
	filemeta += ".up"
	file, err := OpenOrCreateFile(filemeta)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(up_data)
	if err != nil {
		return err
	}
	return nil
}
