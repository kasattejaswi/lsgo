package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

/*
TODO : Add sorting capabilities from command line
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
		readableFlag, _ := cmd.Flags().GetBool("readable")
		if pathFlag == "" {
			pathFlag = "."
		}
		if longState {
			printLongList(listFiles(pathFlag, readableFlag), allFlag)
		} else {
			printShortList(listFiles(pathFlag, readableFlag), allFlag)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("path", "p", "", "See contents of a path")
	listCmd.Flags().BoolP("all", "a", false, "List all files including hidden files")
	listCmd.Flags().BoolP("long", "l", false, "List all files including hidden files")
	listCmd.Flags().BoolP("readable", "r", false, "Prints output in human readable format. Works with long lists")
	listCmd.Flags().StringP("sort", "s", "NAME", "Sorts rows by column names. Works with long lists")
}

type fileItem struct {
	name        string
	fileType    string
	permissions string
	lastMod     string
	isHidden    bool
	owner       string
	group       string
	fileSize    string
}

var colorReset = "\033[0m"
var colorGreen = "\033[32;1m"
var colorBlue = "\033[34;1m"

func getTerminalWidth() int {
	width, _, err := terminal.GetSize(0)
	if err != nil {
		fmt.Println("Warning: Unable to read terminal width. Output will be truncated")
	} else {
		return width
	}
	return 0
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
	terminalWidth := getTerminalWidth()
	maxLength := 0
	tab := 2
	for _, myFile := range fileList {
		if maxLength < len(myFile.name) {
			maxLength = len(myFile.name)
		}
	}
	numC := terminalWidth / (maxLength + tab)
	printCounter := 1
	for _, myFile := range fileList {
		if allFlag {
			if printCounter == numC {
				switch myFile.fileType {
				case "directory":
					fmt.Println(string(colorBlue) + myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)) + string(colorReset))
				case "symlink":
					fmt.Println(string(colorGreen) + myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)) + string(colorReset))
				default:
					fmt.Println(myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)))
				}
				printCounter = 1
			} else {
				switch myFile.fileType {
				case "directory":
					fmt.Print(string(colorBlue) + myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)) + string(colorReset))
				case "symlink":
					fmt.Print(string(colorGreen) + myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)) + string(colorReset))
				default:
					fmt.Print(myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)))
				}
				printCounter += 1
			}
		} else {
			if !myFile.isHidden {
				if printCounter == numC {
					switch myFile.fileType {
					case "directory":
						fmt.Println(string(colorBlue) + myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)) + string(colorReset))
					case "symlink":
						fmt.Println(string(colorGreen) + myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)) + string(colorReset))
					default:
						fmt.Println(myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)))
					}
					printCounter = 1
				} else {
					switch myFile.fileType {
					case "directory":
						fmt.Print(string(colorBlue) + myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)) + string(colorReset))
					case "symlink":
						fmt.Print(string(colorGreen) + myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)) + string(colorReset))
					default:
						fmt.Print(myFile.name + strings.Repeat(" ", tab+maxLength-len(myFile.name)))
					}
					printCounter += 1
				}
			}
		}

	}
	fmt.Println()
}

func printLongList(fileList []fileItem, allFlag bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Name", "Permissions", "Size", "Owner", "Group", "Modified"})
	for _, myFile := range fileList {
		if allFlag {
			switch myFile.fileType {
			case "directory":
				t.AppendRow([]interface{}{string(colorBlue) + myFile.name + string(colorReset), myFile.permissions, myFile.fileSize, myFile.owner, myFile.group, myFile.lastMod})
			case "symlink":
				t.AppendRow([]interface{}{string(colorGreen) + myFile.name + string(colorReset), myFile.permissions, myFile.fileSize, myFile.owner, myFile.group, myFile.lastMod})
			default:
				t.AppendRow([]interface{}{myFile.name, myFile.permissions, myFile.fileSize, myFile.owner, myFile.group, myFile.lastMod})
			}

		} else {
			if !myFile.isHidden {
				switch myFile.fileType {
				case "directory":
					t.AppendRow([]interface{}{string(colorBlue) + myFile.name + string(colorReset), myFile.permissions, myFile.fileSize, myFile.owner, myFile.group, myFile.lastMod})
				case "symlink":
					t.AppendRow([]interface{}{string(colorGreen) + myFile.name + string(colorReset), myFile.permissions, myFile.fileSize, myFile.owner, myFile.group, myFile.lastMod})
				default:
					t.AppendRow([]interface{}{myFile.name, myFile.permissions, myFile.fileSize, myFile.owner, myFile.group, myFile.lastMod})
				}
			}
		}
	}
	t.AppendSeparator()
	t.SetStyle(table.StyleBold)
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

func listFiles(path string, readableFlag bool) []fileItem {
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
		fileName := file.Name()
		var fileSize string
		if readableFlag {
			if stat.Size < 1000 {
				fileSize = strconv.Itoa(int(stat.Size)) + " B"
			} else if stat.Size > 1000 && stat.Size < 1e6 {
				fileSize = strconv.Itoa(int(stat.Size)/1000) + " K"
			} else if stat.Size > 1e6 && stat.Size < 10e9 {
				fileSize = strconv.Itoa(int(stat.Size)/1e6) + " M"
			} else if stat.Size > 1e9 {
				fileSize = strconv.Itoa(int(stat.Size)/1e9) + " G"
			}
		} else {
			fileSize = strconv.Itoa(int(stat.Size))
		}
		if file.Mode()&os.ModeSymlink == os.ModeSymlink {
			followLink, linkErr := os.Readlink(filepath.Join(path, file.Name()))
			if linkErr == nil {
				fileName += " -> " + followLink
			}
			isSymLink = true
		}

		fileList = append(fileList, fileItem{
			name:        fileName,
			fileType:    getFileType(isSymLink, file.IsDir()),
			permissions: file.Mode().String(),
			lastMod:     file.ModTime().Format("Jan 1, 2006 - 3:04pm IST"),
			isHidden:    isFileHidden(file.Name()),
			owner:       owner,
			group:       group,
			fileSize:    fileSize,
		})
	}
	return fileList
}
