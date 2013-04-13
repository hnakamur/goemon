package dir

import (
    "os"
    "strings"
    "sort"
)

type FileInfoPredicate func(os.FileInfo)bool

func SelectFileInfo(pred FileInfoPredicate, infos[]os.FileInfo) []os.FileInfo {
    var results []os.FileInfo
    for _, info := range infos {
        if pred(info) {
            results = append(results, info)
        }
    }
    return results
}

func FindRrdFiles(topDir string) (fi []os.FileInfo, err error) {
    topDirFile, err := os.Open(topDir)
    if err != nil {
        return
    }
    defer topDirFile.Close()

    infos, err := topDirFile.Readdir(-1)
    if err != nil {
        return
    }

    fi = SelectFileInfo(func(info os.FileInfo) bool {
        return strings.HasSuffix(info.Name(), ".rrd")
    }, infos)
	sort.Sort(ByName{fi})
    return
}

type FileInfos []os.FileInfo
type ByName struct{ FileInfos }

func (fi ByName) Len() int {
    return len(fi.FileInfos)
}
func (fi ByName) Swap(i, j int) {
    fi.FileInfos[i], fi.FileInfos[j] = fi.FileInfos[j], fi.FileInfos[i]
}
func (fi ByName) Less(i, j int) bool {
    return fi.FileInfos[i].Name() < fi.FileInfos[j].Name()
}
