package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"text/tabwriter"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

/*
TODO : Add sorting capabilities from command line
TODO : Fix lining of all data without long lists
TODO : Add human readable format
*/

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Shows the files and folders in the current directory",
	Long: `list command lists down all the files and folders present in the current directory. 
By default, it hides the files which starts with .
Use appropriate flags in order to see such kind of hidden files`,
	Run: func(cmd *cobra.Command, args []string) {
		longState, _ := cmd.Flags().GetBool("long")
		allFlag, _ := cmd.Flags().GetBool("all")
		pathFlag, _ := cmd.Flags().GetString("path")
		if pathFlag == "" {
			pathFlag = "."
		}
		if longState {
			printLongList(listFiles(pathFlag), allFlag)
		} else {
			printShortList(listFiles(pathFlag), allFlag)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("path", "p", "", "See contents of a path")
	listCmd.Flags().BoolP("all", "a", false, "List all files including hidden files")
	listCmd.Flags().BoolP("long", "l", false, "List all files including hidden files")
}

type fileItem struct {
	name        string
	fileType    string
	permissions string
	lastMod     string
	isHidden    bool
	owner       string
	group       string
	fileSize    int64
}

func getFileType(isSymLink bool, isDir bool) string {
	if isSymLink {
		return "symlink"
	} else if isDir {
		return "directory"
	} else {
		return "file"
	}
}

func isFileHidden(fileName string) bool {
	if strings.HasPrefix(fileName, ".") {
		return true
	} else {
		return false
	}
}

func printShortList(fileList []fileItem, allFlag bool) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	for _, myFile := range fileList {
		if allFlag {
			fmt.Fprint(w, myFile.name, "\t\t\t")
		} else {
			if !myFile.isHidden {
				fmt.Fprint(w, myFile.name, "\t\t\t")
			}
		}
	}
	w.Flush()
	fmt.Println()
}

func printLongList(fileList []fileItem, allFlag bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Permissions", "Type", "Size", "Owner", "Group", "Modified"})
	for _, myFile := range fileList {
		if allFlag {
			t.AppendRow([]interface{}{myFile.name, myFile.permissions, myFile.fileType, myFile.fileSize, myFile.owner, myFile.group, myFile.lastMod})
		} else {
			if !myFile.isHidden {
				t.AppendRow([]interface{}{myFile.name, myFile.permissions, myFile.fileType, myFile.fileSize, myFile.owner, myFile.group, myFile.lastMod})
			}
		}
	}
	t.AppendSeparator()
	t.SetStyle(table.StyleColoredBlackOnRedWhite)
	t.Render()
}

func getUserFromUid(uid uint32) string {
	u := strconv.FormatUint(uint64(uid), 10)
	usr, err := user.LookupId(u)
	if err == nil {
		return usr.Name
	}
	return "unknown"

}

func getGroupFromGid(gid uint32) string {
	g := strconv.FormatUint(uint64(gid), 10)
	grp, err := user.LookupGroupId(g)
	if err == nil {
		return grp.Name
	}
	return "unknown"
}

func listFiles(path string) []fileItem {
	files, err := ioutil.ReadDir(path)
	var fileList []fileItem
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		stat := file.Sys().(*syscall.Stat_t)
		owner := getUserFromUid(stat.Gid)
		group := getGroupFromGid(stat.Gid)
		isSymLink := false
		if file.Mode()&os.ModeSymlink == os.ModeSymlink {
			isSymLink = true
		}

		fileList = append(fileList, fileItem{
			name:        file.Name(),
			fileType:    getFileType(isSymLink, file.IsDir()),
			permissions: file.Mode().String(),
			lastMod:     file.ModTime().Format("Jan 1, 2006 - 3:04pm IST"),
			isHidden:    isFileHidden(file.Name()),
			owner:       owner,
			group:       group,
			fileSize:    stat.Size,
		})
	}
	return fileList
}
