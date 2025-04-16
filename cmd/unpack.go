package cmd

import (
	"arch/lib/compression"
	"arch/lib/vlc"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var unPackCmd = &cobra.Command{
	Use:   "unpack",
	Short: "Выбери файл для распаковки",
	Run:   unpack,
}

const unpackedExtension = "txt"

func unpack(cmdFlag *cobra.Command, args []string) {
	var decoder compression.Decoder
	method := cmdFlag.Flag("method").Value.String()
	if len(args) == 0 {
		HandleEror(ErrNonPath)
	}

	switch method {
	case "vlc":
		decoder = vlc.New()
	default:
		panic("Неизвестный метод распаковки файла. Fatal")
	}
	filePath := args[0] // Открываем файл и вычитываем его содержимое
	r, err := os.Open(filePath)
	if err != nil {
		HandleEror(err)
	}
	defer r.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		HandleEror(err)
	}

	// data -> Encode(data)
	packed := decoder.Decode(data)
	fmt.Println(string(data), "\n", data) //	TODU: remove
	err = os.WriteFile(unPackedFileName(filePath), []byte(packed), 0644)
	if err != nil {
		HandleEror(err)
	}
}

// Доработать
func unPackedFileName(path string) string {
	// path/to/file/text.txt -> text.txt -> text.vlc
	fileName := filepath.Base(path)               // text.txt
	ext := filepath.Ext(fileName)                 // .txt
	baseName := strings.TrimSuffix(fileName, ext) // text.txt -> text --- .txt => text

	return baseName + "." + unpackedExtension

}

func init() {

	// Добавить решение проблемы имен, т.к ./arch pack vlc ; ./arch unpack ???(vlc); Можно попробовать добавить флаги ./arch pack -p vlc
	rootCmd.AddCommand(unPackCmd)

	unPackCmd.Flags().StringP("method", "m", "", "Выбор метода распаковки: vlc архивация - 'm'.")

	if err := packCmd.MarkFlagRequired("method"); err != nil {
		panic("Произошла ошибка при выборе флага" + err.Error())
	}
}
