package torrent

import (
	"bufio"
	"errors"
	"os"

	"github.com/samarth8765/bittorrent-client/bencode"
	"github.com/samarth8765/bittorrent-client/entities"
)

// extracting SHA1 values
func batch(data []byte, batchsz int) []entities.SHAHash {
	var res []entities.SHAHash
	for i := 0; i < len(data); i += batchsz {
		hash := entities.SHAHash{}
		end := i + batchsz
		if end > len(data) {
			end = len(data)
		}
		copy(hash[:], data[i:end])
		res = append(res, hash)
	}
	return res
}

func getTorrent(reader *bufio.Reader) (*entities.Torrent, error) {
	data, err := bencode.Decoder(reader)
	if err != nil {
		return nil, err
	}

	tData, ok := data.(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid Torrent file")
	}

	tInfoData, ok := tData["info"].(map[string]interface{})
	if !ok {
		return nil, errors.New("Invalid Torrent file")
	}

	torrent := &entities.Torrent{}
	torrent.Announce = tData["announce"].(string)
	torrent.Info = entities.FileInfo{
		Name:        tInfoData["name"].(string),
		PieceLength: tInfoData["piece length"].(int64),
		Pieces:      batch([]byte(tInfoData["pieces"].(string)), 20),
		Length:      tInfoData["length"].(int64),
	}
	torrent.InfoRaw = tInfoData
	return torrent, nil

}

func ParseTorrent(filePath string) (*entities.Torrent, error) {
	// file path for torrent file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	// close the file
	defer file.Close()
	return getTorrent(bufio.NewReader(file))

}
