package vlc

import (
	"strings"
	"unicode"
)

type EncoderDecoder struct{}

func New() EncoderDecoder {
	return EncoderDecoder{}
}

// Функция кодировки текста в hexChunks: My name is Ted -> H7 F12 CC J1;
func (_ EncoderDecoder) Encode(str string) []byte {
	//Приводим строку к нужному нам виду, заглавные буквы конвертируем: M -> !m
	str = PrepareText(str)
	//Конвертируем нашу строку в бинарный формат: !mister !bred -> 01010101001010111001101100101
	bStr := EncodeBin(str)
	//Разбиваем бинарный код на чанки по 8 бит, 10101000 00110000 1111111 00000000
	chunks := SplitBinByChunks(bStr, chunkSize)
	//Конвертируем байты в 16-ричную систему счисления
	//Возвращаем результат из 16-ричной системы
	return chunks.Bytes()
}

// Функция де-кодировки текста: H7 F12 CC J1 -> My name is Ted;
func (_ EncoderDecoder) Decode(encodedData []byte) string {
	res := NewBinChunks(encodedData).Join()
	//Строим decoding tree
	dTree := getEcodingTable().DecodingTree()
	//Конвертируем информацию из decoding tree в text(string)
	return exportText(dTree.Decode(res))
}

// Функция складывания бинарных чанков в строку из 0-ей и 1-ниц: 01110000 1110000 -> 011100001110000
func (bsc BinaryChunks) Join() string {
	var buf strings.Builder

	for _, chunks := range bsc {
		buf.WriteString(string(chunks))
	}

	return buf.String()
}

// Функция для конвертации строки в бинарный код
func EncodeBin(text string) string {
	var buf strings.Builder

	for _, value := range text {
		buf.WriteString(bin(value))
	}
	return buf.String()
}

// Функция конвертирующая rune[] -> string()
func bin(ch rune) string {
	table := getEcodingTable()

	res, ok := table[ch]
	if !ok {
		panic("Символ не распознан: " + string(ch))
	}
	return res
}

// Функция с таблицей кодировки
func getEcodingTable() EncodingTable {
	return EncodingTable{
		' ': "11",
		'e': "101",
		't': "1001",
		'o': "10001",
		'n': "10000",
		'a': "011",
		's': "0101",
		'i': "01001",
		'r': "01000",
		'h': "0011",
		'd': "00101",
		'l': "001001",
		'u': "00011",
		'c': "000101",
		'f': "000100",
		'm': "000011",
		'p': "0000101",
		'g': "0000100",
		'w': "0000011",
		'b': "0000010",
		'y': "0000001",
		'v': "00000001",
		'j': "000000001",
		'k': "0000000001",
		'x': "00000000001",
		'z': "000000000000",
		'q': "000000000001",
		'!': "001000",
	}
}

// Функция орабатывает символы текста к нижнему регистру M -> !m
func PrepareText(text string) string {
	var buf strings.Builder

	for _, value := range text {
		if unicode.IsUpper(value) {
			buf.WriteRune('!')
			buf.WriteRune(unicode.ToLower(value))
		} else {
			buf.WriteRune(value)
		}
	}
	return buf.String()
}

// exportText Возвращает обработанный текст после метода Decode
func exportText(text string) string {
	var buf strings.Builder
	isCapital := false
	for _, ch := range text {
		if isCapital == true {
			buf.WriteRune(unicode.ToUpper(ch))
			isCapital = false

			continue
		}

		if ch == '!' {
			isCapital = true
			continue
		} else {
			buf.WriteRune(ch)
		}
	}

	return buf.String()
}
