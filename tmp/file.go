package tmp

import (
	"os"
	"path"

	"github.com/google/uuid"
)

var FOLDER = "TMP"

func MoveAppWorkingDirecotryToTmp() {

	tmpPath := path.Join(".", FOLDER)

	_ = os.Mkdir(tmpPath, os.ModePerm)

	err := os.Chdir(tmpPath)
	if err != nil {
		panic(err)
	}
}

func OpenTmpFile(postfix string) *os.File {
	fn := GenerateTmpFileName(postfix)
	f, err := os.Open(fn)
	if err != nil {
		return nil
	}
	return f
}

func GenerateTmpFileName(postFix string) string {
	return "tmp_" + uuid.New().String() + "_"
}
