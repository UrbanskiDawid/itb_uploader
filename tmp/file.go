package tmp

import (
	"os"
	"path"

	"github.com/google/uuid"
)

var FOLDER = "TMP"

func OpenTmpFile(postfix string) *os.File {
	fn := GenerateTmpFileName(postfix)
	f, err := os.Open(fn)
	if err != nil {
		return nil
	}
	return f
}

func GenerateTmpFileName(postFix string) string {
	return path.Join(FOLDER, "tmp_"+uuid.New().String()+"_"+postFix)
}
