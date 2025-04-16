package vlc

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Тип для таблицы со значениями кодировок
type EncodingTable map[rune]string

// Тип для кастомного типа передаваемого чанка
type BinaryChunk string

// Слайс со всеми чанками
type BinaryChunks []BinaryChunk

// Интересующая нас размерность одного чанка
const chunkSize = 8

// SplitBinByChunks Функция разбития на чанки
func SplitBinByChunks(bStr string, chunkSize int) BinaryChunks {
	//Точное количество символов в строке
	lenBstr := utf8.RuneCount([]byte(bStr))
	// Вычисляем размерность деля длинну строки на необходимое количество нам чанков
	chunksCount := lenBstr / chunkSize
	// chunkSize = 8 1111000011110000 -> 11110000 11110000
	// В случае когда не делится нацело
	if lenBstr/chunkSize != 0 {
		chunksCount++
	}
	//создаем слайс типа BinaryChunks с размерностью chunksCount , её мы вычисляем выше.
	res := make(BinaryChunks, 0, chunksCount)
	//Создаем объект типа стрингс билдер
	var buf strings.Builder
	//Бежим по символу, записываем в буфер, конвертируем и прибавляем один пока не будет равен размерности необходимого чанка
	for i, ch := range bStr {
		buf.WriteString(string(ch))
		if (i+1)%chunkSize == 0 {
			res = append(res, BinaryChunk(buf.String())) // записываем в слайс данные с буфера
			buf.Reset()                                  // сбрасываем буфер, т.к переходим к другому символу
		}
	}
	if buf.Len() != 0 { // если в буфере осталось что-то, то мы создаем последний чанк и заполняем его нулям, помимо оставшихся прочих чисел
		lastChunk := buf.String()

		lastChunk += strings.Repeat("0", chunkSize-len(lastChunk)) // тут заполнение нулями, пока опять же не будет достигнута нужная
		res = append(res, BinaryChunk(lastChunk))                  // размерность чанка
	}
	return res // возвращаем обработанную строку в виде слайса чанков
}

func NewBinChunk(code byte) BinaryChunk {
	return BinaryChunk(fmt.Sprintf("%08b", code))
}

func NewBinChunks(data []byte) BinaryChunks {

	res := make(BinaryChunks, 0, len(data))

	for _, code := range data {
		res = append(res, NewBinChunk(code))
	}

	return res
}

// функция для создания слайса байтов кодированной информации
func (bcs BinaryChunks) Bytes() []byte {
	res := make([]byte, 0, len(bcs))
	for _, bc := range bcs {
		res = append(res, bc.Byte())
	}

	return res
}

// функция для конвертации в байт
func (bc BinaryChunk) Byte() byte {
	num, err := strconv.ParseUint(string(bc), 2, chunkSize)
	if err != nil {
		panic("Произошла ошибка:" + err.Error())
	}
	return byte(num)
}
