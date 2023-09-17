package bencode

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
)

func Decoder(reader *bufio.Reader) (interface{}, error) {
	ch, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}

	if ch == 'i' {
		var buffer []byte
		for {
			ch, err := reader.ReadByte()

			if err != nil {
				return nil, err
			}

			if ch == 'e' {
				value, err := strconv.ParseInt(string(buffer), 10, 64)
				if err != nil {
					return nil, errors.New(fmt.Sprintf("Invalid integer %s", string(buffer)))
				}

				return value, nil
			}

			buffer = append(buffer, ch)
		}

	} else if ch == 'l' {
		var listholder []interface{}
		for {
			ch, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}
			if ch == 'e' {
				return listholder, nil
			}
			reader.UnreadByte()
			data, err := Decoder(reader)
			if err != nil {
				return nil, err
			}

			listholder = append(listholder, data)
		}

	} else if ch == 'd' {
		var dictholder = map[string]interface{}{}
		for {
			ch, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}

			if ch == 'e' {
				return dictholder, nil
			}

			// reading the key
			reader.UnreadByte()
			data, err := Decoder(reader)
			if err != nil {
				return nil, err
			}

			// Bencode keys are always string
			key, check := data.(string)
			if !check {
				return nil, errors.New(fmt.Sprintf("Key value is not string -> %s", key))
			}

			// read the value
			val, err := Decoder(reader)
			if err != nil {
				return nil, err
			}

			dictholder[key] = val

		}

	} else {
		reader.UnreadByte()
		var lengthBuf []byte
		for {
			ch, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}

			if ch == ':' {
				break
			}

			lengthBuf = append(lengthBuf, ch)
		}

		length, err := strconv.Atoi(string(lengthBuf))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Invalid value of length -> %v", length))
		}

		var strBuf []byte
		for i := 0; i < length; i++ {
			ch, err := reader.ReadByte()
			if err != nil {
				return nil, err
			}

			strBuf = append(strBuf, ch)
		}

		return string(strBuf), nil
	}
}

// func main() {
// 	file, _ := os.Open("temp.torrent")
// 	reader := bufio.NewReader(file)
// 	defer file.Close()

// 	decode, err := Decoder(reader)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fmt.Println(decode)

// }
