package entities

type FileInfo struct {
	Name        string
	Pieces      []SHAHash
	PieceLength int64
	Length      int64
}

type Torrent struct {
	Announce string
	Info     FileInfo
	InfoRaw  map[string]interface{}
}
