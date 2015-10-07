package main

//import "bufio"
import "crypto/tls"
import "fmt"
import "flag"
import "log"
import "os"
import "strings"
import "github.com/nsf/termbox-go"
import "github.com/mattn/go-xmpp"
import "golang.org/x/crypto/ssh/terminal"

var width, height int
var lastMessage = ""

var server = flag.String("server", "talk.google.com:443", "server")
var username = flag.String("username", "", "username")
var password string
var status = "xa"
var statusMessage = "I for one welcome our new codebot overlords."
var session = flag.Bool("session", false, "use server session")

func drawWelcomeScreen() {
	welcomeMessage := "Hello world"
	msgLen := len(welcomeMessage)
	xOffset := width/2 - msgLen/2
	yOffset := height / 2

	for _, c := range welcomeMessage {
		termbox.SetCell(xOffset, yOffset, c, termbox.ColorDefault, termbox.ColorDefault)
		xOffset += 1
	}

	if lastMessage != "" {
		xOffset = width/2 - len(lastMessage)/2
		yOffset++
		for _, c := range lastMessage {
			termbox.SetCell(xOffset, yOffset, c, termbox.ColorDefault, termbox.ColorDefault)
			xOffset += 1
		}
	}
}

func main() {
	logFile, err := os.Create("/tmp/gochat.log")
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	fmt.Println("Enter password:")
	passwdBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))

	if err != nil {
		panic(err)
	}

	password = string(passwdBytes)

	flag.Parse()

	if *username == "" || password == "" {
		flag.Usage()
	}

	xmpp.DefaultConfig = tls.Config{
		ServerName:         strings.Split(*server, ":")[0],
		InsecureSkipVerify: true,
	}

	var talk *xmpp.Client

	options := xmpp.Options{Host: *server,
		User:          *username,
		Password:      password,
		NoTLS:         true,
		Debug:         false,
		Session:       *session,
		Status:        status,
		StatusMessage: statusMessage,
	}

	fmt.Println("Attempting to create client")
	talk, err = options.NewClient()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Created client")

	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, height = termbox.Size()

	renderer := &TermBoxRenderer{}
	renderer.OnViewPortResize(width, height)

	eventChannelToModel := GetEventChannelToModel(renderer)

	drawWelcomeScreen()
	termbox.Flush()

	go func() {
		for {
			chat, err := talk.Recv()
			if err != nil {
				log.Fatal(err)
			}
			switch v := chat.(type) {
			case xmpp.Chat:
				lastMessage = "Chat:" + v.Remote + "/" + v.Text
				log.Printf("Received chat update. remote:%v, text:%v, roster:%v, type:%v", v.Remote, v.Text, v.Roster, v.Type)
				eventChannelToModel <- &MessageReceivedEvent{v}
			case xmpp.Presence:
				log.Printf("Received presence update. from:%v, to:%v, show:%v, type:%v", v.From, v.To, v.Show, v.Type)
				lastMessage = "Presence:" + v.From + "/" + v.Show
				eventChannelToModel <- &PresenceUpdateEvent{v}
			}
		}
	}()

	ctrlxpressed := false

	//	for {
	//		in := bufio.NewReader(os.Stdin)
	//		line, err := in.ReadString('\n')
	//		if err != nil {
	//			continue
	//		}
	//		line = strings.TrimRight(line, "\n")
	//
	//		tokens := strings.SplitN(line, " ", 2)
	//		if len(tokens) == 2 {
	//			talk.Send(xmpp.Chat{Remote: tokens[0], Type: "chat", Text: tokens[1]})
	//		}
	//	}

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			//			if ev.Key == termbox.KeyCtrlS && ctrlxpressed {
			//				termbox.Sync()
			//			}
			if ev.Key == termbox.KeyCtrlQ && ctrlxpressed {
				break loop
			}
			if ev.Key == termbox.KeyCtrlC && ctrlxpressed {
				eventChannelToModel <- &SwitchViewEvent{CHAT_WINDOW}
			}
			if ev.Key == termbox.KeyCtrlP && ctrlxpressed {
				eventChannelToModel <- &SwitchViewEvent{CONTACT_WINDOW}
			}
			if ev.Key == termbox.KeyCtrlX {
				ctrlxpressed = true
			} else {
				ctrlxpressed = false
			}

			//			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			//			draw_keyboard()
			//			dispatch_press(&ev)
			//			pretty_print_press(&ev)
			//			termbox.Flush()
		case termbox.EventResize:
			eventChannelToModel <- &ResizeEvent{ev.Width, ev.Height}
			//			pretty_print_resize(&ev)
		case termbox.EventError:
			log.Fatal(ev.Err)
			panic(ev.Err)
		}
	}
}
