package utils

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/mattn/go-tty"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

var mux sync.RWMutex

const (
	// Unknown message type, maybe sent by a bug
	//UnknownInput = '0'
	// User input typically from a keyboard
	Input = '1'
	// Ping to the server
	Ping = '2'
	// Notify that the browser size has been changed
	ResizeTerminal = '3'
)

const duration = 5

const (
	// Unknown message type, maybe set by a bug
	//UnknownOutput = '0'
	// Normal output to the terminal
	Output = '1'
	// Pong to the browser
	Pong = '2'
	// Set window title of the terminal
	SetWindowTitle = '3'
	// Set terminal preference
	//SetPreferences = '4'
	// Make terminal to reconnect
	//SetReconnect = '5'
)

var logfilePath = GetLogFilePath("sdb.log")

func GetCurrentTime() string {
	t := time.Now()
	currentTime := t.Format("2006-01-02 15:04:05")
	return currentTime
}

func GetLogFilePath(filename string) string {
	homePath := GetStringEnv("HOME", "home-env")
	if homePath == "home-env" {
		log.Error().Err(errors.New("env error")).Msg("Get HOME env error")
		return ""
	}
	logfilePath := fmt.Sprintf("%s/%s", homePath, filename)
	return logfilePath
}

func Connect(parameter []string, cmd *cobra.Command) {
	var origin *string
	var keepAlive bool
METHOD:
	for _, p := range parameter {
		switch {
		case p == "url" || p == "u":
			v, err := GetFlagValueStringP(cmd, "url")
			if err != nil {
				return
			}
			keepAlive, err = GetFlagValueBool(cmd, "alive")
			if err != nil {
				keepAlive = false
				return
			}
			if origin, err = GetFlagValueStringP(cmd, "origin"); err != nil {
				return
			}
			if *origin == "" {
				*origin = "https://www.pypypy.cn"
			}
			// check origin is valid
			if !CheckURLValid(origin) {
				log.Print("origin is not url")
				return
			}
			DebugForURL(*v, *origin, keepAlive)
			break METHOD
		case p == "ep":
			ep, err := GetFlagValueStringP(cmd, p)
			if err != nil {
				return
			}
			t, err := GetFlagValueStringP(cmd, "token")
			if err != nil {
				return
			}
			if origin, err = GetFlagValueStringP(cmd, "origin"); err != nil {
				return
			}
			if *origin == "" {
				*origin = "https://www.pypypy.cn"
			}
			keepAlive, err = GetFlagValueBool(cmd, "alive")
			if err != nil {
				keepAlive = false
				return
			}
			if !CheckURLValid(origin) {
				log.Print("origin is not url")
				return
			}
			DebugForEp(*ep, *t, *origin, keepAlive)
			break METHOD
		default:
			log.Printf("Unknown connect sandbox method %s", p)
			return
		}
	}
	return
}

func DebugForURL(wssURL string, origin string, keepAlive bool) {
	if wssURL == "" {
		log.Print("websocket url can not be empty")
		return
	}
	wssURL = strings.Replace(wssURL, " ", "", -1)
	parsedURL, err := url.Parse(wssURL)
	if err != nil {
		log.Error().Err(err).Msg("parse url error")
		return
	}
	SandboxDebug(parsedURL, origin, keepAlive)
}

func DebugForEp(endpoint string, token string, origin string, keepAlive bool) {
	if endpoint == "" || token == "" {
		log.Print("must specific token when using sandbox endpoint")
		return
	}
	endpoint = strings.Replace(endpoint, " ", "", -1)
	token = strings.Replace(token, " ", "", -1)
	u, err := url.Parse(endpoint)
	if err != nil {
		glog.Fatal("Parse url failed: ", err)
		return
	}
	u.Scheme = "wss"
	u.Path = path.Join(u.Path, "terminal", "ws")
	u.RawQuery = fmt.Sprintf("token=%s", token)
	SandboxDebug(u, origin, keepAlive)
}

func WriteLog(filepath string, content string) error {
	if filepath == "" {
		return errors.New("filepath can not be empty")
	}
	var fl *os.File
	var err error
	fl, err = os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		if !Exists(filepath) {
			fl, err = os.Create(filepath)
			if err != nil {
				log.Error().Err(err).Msgf("create file %s error %v ", filepath, err)
				return err
			}

		}
		log.Error().Err(err).Msgf("open file error: %v", err)
		return err
	}
	defer func() {
		err = fl.Close()
		if err != nil {
			log.Error().Err(err).Msg("file close error")
		}
	}()
	write := bufio.NewWriter(fl)
	content = fmt.Sprintf("[%s]  %s \n", GetCurrentTime(), content)
	_, err = write.WriteString(content)
	if err != nil {
		log.Error().Err(err).Msgf("Writer write to buffer error %v", err)
		return err
	}
	err = write.Flush()
	if err != nil {
		log.Error().Err(err).Msgf("Writer flush to file error %v", err)
		return err
	}
	return nil
}

func SandboxDebug(wssUrl *url.URL, wsOrigin string, keepAlive bool) {
	if wssUrl.String() == "" || wsOrigin == "" {
		return
	}
	if logfilePath == "" {
		return
	}
	content := fmt.Sprintf("connecting to: %s", wssUrl.String())
	if err := WriteLog(logfilePath, content); err != nil {
		return
	}
	header := http.Header{
		"Origin":                 []string{wsOrigin},
		"Sec-WebSocket-Protocol": []string{"webtty"},
		"Sec-WebSocket-Version":  []string{"13"},
	}
	// create ws connection ，parse sandbox_token from wss url
	c, _, err := websocket.DefaultDialer.Dial(wssUrl.String(), header)
	token := wssUrl.Query().Get("token")
	_token := &token
	if err != nil {
		log.Error().Err(err).Msg("connect to sandbox failed. check the websocket url and origin again")
		fmt.Println(wssUrl.String())
		return
	}
	defer func() {
		if err := c.Close(); err != nil {
			log.Error().Err(err).Msg("close connection error")
		}
	}()

	done := make(chan int)

	// set  operation when receiving close msg
	c.SetCloseHandler(func(code int, text string) error {
		message := websocket.FormatCloseMessage(code, "")
		if err := c.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second)); err != nil {
			log.Error().Err(err).Msg("write control message error")
			return err
		}
		done <- 1
		return nil
	})

	go CheckConnection(c, done, keepAlive)

	// init ws connection
	err = WriteMsg(c, websocket.TextMessage, []byte(fmt.Sprintf(`{"Arguments":"","AuthToken":"%s"}`, *_token)))
	//err = c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"Arguments":"","AuthToken":"%s"}`, *_token)))
	if err != nil {
		log.Fatal().Msgf("write message: %v", err)
	}

	t, err := tty.Open()
	if err != nil {
		log.Error().Err(err).Msg("Open tty error")
	}
	defer func() {
		if err := t.Close(); err != nil {
			fmt.Println("close tty error")
			return
		}
	}()

	w, h, err := t.Size()
	if err != nil {
		log.Error().Err(err).Msg("Get tty size error")
		return
	}
	if logfilePath == "" {
		return
	}
	content = fmt.Sprintf("Resized: %d -- %d", w, h)
	if err := WriteLog(logfilePath, content); err != nil {
		return
	}

	// use `...` syntax for divide second bytes array into single byte to append to first byte array
	err = WriteMsg(c, websocket.TextMessage, append([]byte{ResizeTerminal}, []byte(fmt.Sprintf(`{"colmuns":%d,"rows":%d}`, w, h))...))
	//err = c.WriteMessage(websocket.TextMessage, append([]byte{ResizeTerminal}, []byte(fmt.Sprintf(`{"colmuns":%d,"rows":%d}`, w, h))...))
	if err != nil {
		log.Error().Err(err).Msgf("Resized window Write message error %v", err)
		return
	}

	go TerminalResize(c, t)

	// read raw from tty . that is to say : when you input something to tty and tap enter, this will
	// get tty raw and send it to server by websocket connection
	_, err = t.Raw()
	if err != nil {
		log.Fatal().Err(err)
	}
	go ReadMessage(c, t, done)
	<-done
}

func TerminalResize(c *websocket.Conn, t *tty.TTY) {
	for ws := range t.SIGWINCH() {
		if logfilePath == "" {
			return
		}
		content := fmt.Sprintf("Resized: %d -- %d", ws.W, ws.H)
		if err := WriteLog(logfilePath, content); err != nil {
			return
		}
		err := WriteMsg(c, websocket.TextMessage, append([]byte{ResizeTerminal}, []byte(fmt.Sprintf(`{"colmuns":%d,"rows":%d}`, ws.W, ws.H))...))
		//err := c.WriteMessage(websocket.TextMessage, append([]byte{ResizeTerminal}, []byte(fmt.Sprintf(`{"colmuns":%d,"rows":%d}`, ws.W, ws.H))...))
		if err != nil {
			log.Error().Err(err).Msgf("Resized window Write message error %v ", err)
			return
		}
	}
}

func WriteMsg(c *websocket.Conn, msgType int, data []byte) error {
	defer func() {
		mux.Unlock()
	}()
	mux.Lock()
	err := c.WriteMessage(msgType, data)
	if err != nil {
		return err
	}
	return nil
}

func ReadMessage(c *websocket.Conn, t *tty.TTY, done chan int) {
	reader := bufio.NewReader(t.Input())
	for {
		b, err := reader.ReadByte()
		if err != nil {
			log.Error().Err(err).Msg("Read byte from tty error")
		}

		if b == 0 {
			continue
		}
		// read cmd operation from tty  and send it to sandbox server
		err = WriteMsg(c, websocket.TextMessage, []byte{Input, b})
		//err = c.WriteMessage(websocket.TextMessage, []byte{Input, b})
		if err != nil {
			log.Fatal().Msgf("write message: %v", err)
			done <- 1
			break
		}
	}
}
func PingMessage(conn *websocket.Conn, data []byte, done chan int) {
	time.Sleep(duration)
	err := WriteMsg(conn, websocket.TextMessage, data)
	if err != nil {
		log.Fatal().Msgf("write ping  message: %v", err)
		done <- 1
		panic(err)
	}
}
func CheckConnection(c *websocket.Conn, done chan int, keepAlive bool) {
CHECK:
	for {
		if keepAlive {
			PingMessage(c, []byte{Ping}, done)
		}
		msgType, rawMsg, err := c.ReadMessage()
		if err != nil {
			if !websocket.IsCloseError(err, 1006) {
				log.Error().Err(err).Msgf("read msg error: %v", err)
			} else {
				fmt.Println("\nConnection closed.")
			}
			done <- 1
			return
		}
		if msgType == websocket.CloseMessage {
			fmt.Println("connection closed")
			return
		}
		if msgType != websocket.TextMessage {
			log.Error().Err(err).Msgf("message type note text %v", err)
			continue
		}

		message := string(rawMsg)
		if logfilePath == "" {
			return
		}
		switch message[0] {
		case Pong:
			continue
		case SetWindowTitle:
			content := fmt.Sprintf("[DEBUG] set window title: %s", message[1:])
			if err := WriteLog(logfilePath, content); err != nil {
				break CHECK
			}
			continue
		case Output:
			output, err := base64.StdEncoding.DecodeString(message[1:])
			if err != nil {
				log.Error().Err(err).Msgf("decode: %v", err)
				continue
			}
			content := fmt.Sprintf("[DEBUG] type: Output, raw: %s, output: %s", message[1:], string(output))
			if err := WriteLog(logfilePath, content); err != nil {
				break CHECK
			}
			_, err = fmt.Fprint(os.Stdout, string(output))
			if err != nil {
				log.Error().Err(err).Msgf("write stdout: %v", err)
				continue
			}
		}
	}
}

func CheckURLValid(u *string) bool {
	if u == nil {
		return false
	}
	reg, err := regexp.Compile("((ht|f)tps?):\\/\\/[\\w\\-]+(\\.[\\w\\-]+)+([\\w\\-.,@?^=%&:\\/~+#]*[\\w\\-@?^=%&\\/~+#])?")
	if err != nil {
		log.Error().Err(err).Msg("compile http regex error")
		return false
	}
	matchRS := reg.FindAllString(*u, -1)
	if len(matchRS) == 0 {
		domainReg, err := regexp.Compile("[\\w\\-]+(\\.[\\w\\-]+)+([\\w\\-.,@?^=%&:\\/~+#]*[\\w\\-@?^=%&\\/~+#])")
		if err != nil {
			log.Error().Err(err).Msg("compile domain regex error")
			return false
		}
		matchRS = domainReg.FindAllString(*u, -1)
		if len(matchRS) == 0 {
			return false
		}
		*u = fmt.Sprintf("https://%s", *u)
		return true
	}
	return true
}
