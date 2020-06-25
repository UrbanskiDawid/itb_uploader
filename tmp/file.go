package tmp

import (
	"os"
	"path"
	"strings"

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

func OpenTmpFile(postfix string) (*os.File, error) {
	fn := GenerateTmpFileName(postfix)
	f, err := os.OpenFile(fn, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func GenerateTmpFileName(postFix string) string {
	return "tmp_" + strings.ReplaceAll(uuid.New().String(), "-", "") + "_" + postFix

}
