package tmp

import (
	"os"

	"github.com/google/uuid"
)

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
