package temp

import (
	"github.com/GuoYuefei/DOStorage1/distributed/config"
	"net/http"
	"os"
	"path"
	"strings"
)

func delete(_ http.ResponseWriter, r *http.Request) {
	uu := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := path.Join(config.TempRoot, uu)
	datFile := infoFile + ".dat"
	os.Remove(infoFile)
	os.Remove(datFile)
}
