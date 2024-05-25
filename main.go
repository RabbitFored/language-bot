package main

import (
	"io"
	"net/http"

	"github.com/BRUHItsABunny/ugodict"
	tl "github.com/goTelegramBot/telepher"
	"github.com/goTelegramBot/telepher/types"

	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	translator "github.com/Conight/go-googletrans"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/client9/misspell"
	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/watson-developer-cloud/go-sdk/v2/speechtotextv1"
)

func main() {
	b, err := tl.NewBot(os.Getenv("tkoeko"), nil)

	if err != nil {
		log.Println(err)
		return
	}
	http.HandleFunc("/", handler)
  go http.ListenAndServe(":8080", nil)

	b.Command("start", start)
	b.Command("tr", translate)
	b.Command("urban", define)
	b.Command("spell", spell)
	b.Command("speak", speak)
	b.Command("transcribe", transcript)
	b.Command("help", help)
	b.Command("about", about)
	b.Start()
}

func start(bot tl.Bot, message *types.Message) {
	text := fmt.Sprintf("*வணக்கம் %s*\n\nNot sure what this mean?\nDon't worry I am a language expert.\nI can help you translate texts for you.", message.From.FirstName)

	markup := tl.InlineKeyboardMarkup()
	but1 := types.InlineKeyboardButton{Text: "Channel", Url: "https://t.me/theostrich"}
	row1 := markup.Row(but1)

	keyboard := markup.Parse(row1)

	bot.SendMessage(message.Chat.Id, text, &tl.Options{ReplyMarkup: &keyboard, ParseMode: "Markdown"})

}
func help(bot tl.Bot,message *types.Message){
    text := fmt.Sprintf(`*Hi %s!* 
    Here is a detailed guide on using me.
    
    *Helpful commands:*
    - /start : Check if I am alive! You've probably already used this.
    - /help  : I'll tell you more about myself!
    - /urban : Get urban definition.
    - /tr <lang> : Translates a text.
    - /spell : Spell check a text.
    - /speak : Text to speech
    - /transcribe : Transcribe an audio
    - /about  : Know about me.
    `, message.From.FirstName)
    
    markup := tl.InlineKeyboardMarkup()
but1 := types.InlineKeyboardButton{Text:"Get Help",Url: "https://t.me/ostrichdiscussion"}
but2 := types.InlineKeyboardButton{Text:"🔖Add me in group",Url: "https://t.me/thelanguageBot?startgroup=new"}
    row1 := markup.Row(but1,but2)

  keyboard := markup.Parse(row1)
     
     bot.SendMessage(message.Chat.Id, text,&tl.Options{ReplyMarkup:&keyboard,ParseMode:"Markdown"})
}
func about(bot tl.Bot,message *types.Message){
    text :="<b>About Me :</b>\n\n" +
    "  - <b>Name        :</b> thelanguagebot\n" +
    "  - <b>Creator     :</b> @theostrich\n" +
    "  - <b>Language  :</b> Golang\n" +
    "  - <b>Library      :</b> <a href=\"https://github.com/goTelegramBot/gogram\">Gogram</a>"

    markup := tl.InlineKeyboardMarkup()
but1 := types.InlineKeyboardButton{Text:"Channel",Url: "https://t.me/theostrich"}
but2 := types.InlineKeyboardButton{Text:"Support Group",Url: "https://t.me/ostrichdiscussion"}
    row1 := markup.Row(but1,but2)

  keyboard := markup.Parse(row1)

     bot.SendMessage(message.Chat.Id, text,&tl.Options{ReplyMarkup:&keyboard,ParseMode:"html",DisableWebPagePreview:true})
}
func translate(b tl.Bot, m *types.Message) {
	var message string
	if m.ReplyToMessage != nil {

		message = m.ReplyToMessage.Text
	} else {
		b.SendMessage(m.Chat.Id, "_Reply to any message with_ /tr <lang> _to translate_", &tl.Options{ParseMode: "Markdown"})
		return
	}
	args := m.Args()
  if len(args) == 1 {
    b.SendMessage(m.Chat.Id, "*Specify some language codes.\nExample:* `/tr ta`", &tl.Options{ParseMode: "Markdown"})
		return
  }
	if args[1] == "<lang>" {
		b.SendMessage(m.Chat.Id, "*Specify some language codes.\nExample:* `/tr ta`", &tl.Options{ParseMode: "Markdown"})
		return
	}
	t := translator.New()
	result, err := t.Translate(message, "auto", args[1])
	if err != nil {
		log.Println(err)
	}

	b.SendMessage(m.Chat.Id, result.Text, &tl.Options{ParseMode: "Markdown"})

}

func spell(b tl.Bot, m *types.Message) {

	var message string
	if m.ReplyToMessage != nil {

		message = m.ReplyToMessage.Text
	} else {
		b.SendMessage(m.Chat.Id, "_Reply to any message with_ /spell _to check misspells in texts_", &tl.Options{ParseMode: "Markdown"})
		return
	}

	r := misspell.Replacer{
		Replacements: misspell.DictMain,
	}
	r.Compile()
	var updated string
	var changes []misspell.Diff

	updated, changes = r.Replace(message)

	var changing string

	for _, diff := range changes {

		li := fmt.Sprintf("*%d:%d :* `%s` misspelled as `%s`\n", diff.Line, diff.Column, diff.Original, diff.Corrected)
		changing = changing + li
	}
	if len(changes) == 0 {
		changing = "_None_"
	}
	length := strings.Count(message, "\n")
	var text string
	if length < 7 {
		text = fmt.Sprintf("*Found %d misspellings*\n\n*Original Text:*\n`%s`\n\n*Corrected Text:* \n`%s`\n\n*Changes:*\n%s", len(changes), message, updated, changing)
	} else {
		text = fmt.Sprintf("*Found %d misspellings*\n\n*Corrected Text:*\n`%s`\n\n*Changes:*\n%s", len(changes), updated, changing)
	}
	if len(text) < 4000 {
		b.SendMessage(m.Chat.Id, text, &tl.Options{ParseMode: "Markdown"})
	} else {
		b.SendMessage(m.Chat.Id, "*Text is too long*", &tl.Options{ParseMode: "Markdown"})
	}

}
func define(b tl.Bot, m *types.Message) {
	args := m.Args()
	var word string
	if m.ReplyToMessage != nil {

		word = m.ReplyToMessage.Text
	}
	if len(args) != 1 {
		word = strings.Join(args[1:], " ")
	}

	client := ugodict.GetClient()
	results, err := client.DefineByTerm(word)

	if err == nil {
		// Select definition
		definition := results[0]

		text := fmt.Sprintf("*Word:* %s\n*Definition:* %s\n*Example:*\n%s", word, definition.Definition, definition.Example)
		b.SendMessage(m.Chat.Id, text, &tl.Options{ParseMode: "Markdown"})
	} else {

		log.Println(err)

		if strings.Contains(err.Error(), "no results found") {
			b.SendMessage(m.Chat.Id, "*No results found*", &tl.Options{ParseMode: "Markdown"})
		}

	}

}
func speak(b tl.Bot, m *types.Message) {

	var message string
	args := m.Args()
	if m.ReplyToMessage != nil {

		message = m.ReplyToMessage.Text
	} else {
		b.SendMessage(m.Chat.Id, "_Reply to any message with_ /speak _to make me speak_", &tl.Options{ParseMode: "Markdown"})
		return
	}
	lang := "en"
	if len(args) != 1 {
		lang = args[1]
	}

	speech := htgotts.Speech{Folder: "audio", Language: lang}
	speech.Speak(message)

	send_options := make(map[string]string)
	send_options["title"] = "Voice - @thelanguagebot"
	send_options["performer"] = "theostrich"

	files, err := listFile("./audio")
	if err != nil {
		log.Println(err)
		b.SendMessage(m.Chat.Id, "_Error: Something bad happened. Contact my support team - @ostrichdiscussion_", &tl.Options{ParseMode: "Markdown"})

		return

	}
	for _, file := range files {
		path := fmt.Sprintf("audio/" + file)

		document := types.InputFile{
			FilePath: path,
		}
		b.SendAudio(m.Chat.Id, document, send_options)
		os.Remove(path)
	}

}

func listFile(directory string) ([]string, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}
	var file []string
	for _, f := range files {
		file = append(file, f.Name())
		if f.Size() < 2222 {
			os.Remove("audio/" + f.Name())
			return nil, fmt.Errorf("ERR %d", f.Size())
		}

	}
	return file, nil
}
func transcript(b tl.Bot, m *types.Message) {
	var FileId string
	if m.ReplyToMessage != nil {
		if m.ReplyToMessage.Audio == nil {
			b.SendMessage(m.Chat.Id, "_Reply to audio file with_ /transcribe _to generate its text_", &tl.Options{ParseMode: "Markdown"})
			return
		}
		FileId = m.ReplyToMessage.Audio.FileId
		duration := m.ReplyToMessage.Audio.Duration
		if duration > 60 {
			b.SendMessage(m.Chat.Id, "*Audio files longer than 1 minute cannot be transcribed for free users. To Upgrade to paid plans contact us via:*\n - @ostrichdiscussion\n - @contactOstrichBot", &tl.Options{ParseMode: "Markdown"})
			return
		}

	} else {
		b.SendMessage(m.Chat.Id, "_Reply to audio file with_ /transcribe _to generate its text_", &tl.Options{ParseMode: "Markdown"})
		return
	}

	file, err := b.GetFile(FileId)
	if err != nil {
		fmt.Println(err)

	}
	url := "https://api.telegram.org/file/bot{token}/" + file.FilePath
	path := "download/ostrich.mp3"
	err = DownloadFile(path, url)
	words, err := watson(path)
	if err != nil {
		log.Println(err)
		if err.Error() == "Empty file" {
			b.SendMessage(m.Chat.Id, "*Empty or Invalid file provided*", &tl.Options{ParseMode: "Markdown"})
		}
		return
	}
	text := *words

	b.SendMessage(m.Chat.Id, text, &tl.Options{ParseMode: "Markdown"})
	os.Remove(path)

}

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func watson(path string) (*string, error) {
	files, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	fi, err := files.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Size() < 2222 {
		return nil, fmt.Errorf("Empty file")
	}
	authenticator := &core.IamAuthenticator{
		ApiKey: "watsonapikey",
	}

	options := &speechtotextv1.SpeechToTextV1Options{
		Authenticator: authenticator,
	}

	speechToText, speechToTextErr := speechtotextv1.NewSpeechToTextV1(options)

	if speechToTextErr != nil {
		panic(speechToTextErr)
	}

	speechToText.SetServiceURL(os.Getenv("sp"))
	file := path

	var audioFile io.ReadCloser
	var audioFileErr error
	audioFile, audioFileErr = os.Open(file)
	if audioFileErr != nil {
		fmt.Println(audioFileErr)
	}
	result, _, responseErr := speechToText.Recognize(
		&speechtotextv1.RecognizeOptions{
			Audio:                     audioFile,
			Timestamps:                core.BoolPtr(true),
			WordAlternativesThreshold: core.Float32Ptr(0.9),
		},
	)
	if responseErr != nil {
		panic(responseErr)
	}

	return result.Results[0].Alternatives[0].Transcript, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}
