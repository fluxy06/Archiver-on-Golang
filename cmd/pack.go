package cmd

import (
	"arch/lib/compression"
	"arch/lib/vlc"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "Выбери файл для архивации",
	Run:   pack,
}

var ErrNonPath = errors.New("Укажите путь файла.")

const packedExtension = "vlc"

func pack(cmdFlag *cobra.Command, args []string) {
	var encoder compression.Encoder
	if len(args) == 0 {
		HandleEror(ErrNonPath)
	}

	method := cmdFlag.Flag("method").Value.String()

	switch method {
	case "vlc":
		encoder = vlc.New()
	default:
		panic("Неизвестный метод распаковки файла. Fatal")
	}
	filePath := args[0]
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
	packed := encoder.Encode(string(data))
	fmt.Println(string(data), "\n", data) //	TODU: remove
	err = os.WriteFile(packedFileName(filePath), packed, 0644)
	if err != nil {
		HandleEror(err)
	}
}

func packedFileName(path string) string {
	// path/to/file/text.txt -> text.txt -> text.vlc
	fileName := filepath.Base(path)               // text.txt
	ext := filepath.Ext(fileName)                 // .txt
	baseName := strings.TrimSuffix(fileName, ext) // text.txt -> text --- .txt => text

	return baseName + "." + packedExtension
}

func init() {
	rootCmd.AddCommand(packCmd)

	packCmd.Flags().StringP("method", "m", "", "Выбор метода сжатия: vlc архивация - 'm'.")

	if err := packCmd.MarkFlagRequired("method"); err != nil {
		panic("Произошла ошибка при выборе флага" + err.Error())
	}
}
