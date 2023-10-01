package new_bot

import (
	_ "embed"
	"log"
	"strings"

	"github.com/NicoNex/echotron/v3"
	"github.com/cohere-ai/cohere-go"
)

type stateFn func(*echotron.Update) stateFn

type bot struct {
	chatID int64
    API echotron.API
    cohereClient *cohere.Client
    state        stateFn
}

var (
	//go:embed tgtoken
	TelegramToken	string
	//go:embed chtoken
	cohereAPIKey	string
	Commands = []echotron.BotCommand{
		{Command: "/start", Description: "Activate the bot."},
		{Command: "/generate", Description: "Generate an answer."},
	}
)

func NewBot (chatID int64) echotron.Bot {
	cohereClient, err := cohere.CreateClient(cohereAPIKey)
	if err != nil {
		log.Fatalln(err)
	}

	b := &bot{
		chatID: chatID,
		API: echotron.NewAPI(TelegramToken),
		cohereClient: cohereClient,
	}
	b.state = b.handleMessage
	return b
}

func (b *bot) handlePrompt(update *echotron.Update) stateFn {
	_, err := b.API.SendChatAction(echotron.Typing, b.chatID, nil)
	if err != nil {
		log.Fatalln(err)
	}

	response, err := b.generateText(message(update))
	if err != nil {
        log.Println("handlePrompt", err)
        b.API.SendMessage("An error occurred!", b.chatID, nil)
        return b.handleMessage
    }

	_, err = b.API.SendMessage(response, b.chatID, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return b.handleMessage
}

func (b *bot) handleMessage(update *echotron.Update) stateFn {
	switch m := message(update); {
	case strings.HasPrefix(m, "/start"):
		_, err := b.API.SendMessage("Hello world", b.chatID, nil)
		if err != nil {
			log.Fatalln(err)
		}
	case strings.HasPrefix(m, "/generate"):
		_, err := b.API.SendMessage("Please enter a prompt:", b.chatID, nil)
		if err != nil {
			log.Fatalln(err)
		}
		return b.handlePrompt
	}
	return b.handleMessage
}

func (b *bot) Update(update *echotron.Update) {
	b.state = b.state(update)
}

func (b *bot) generateText(prompt string) (string, error) {
	var valToken uint
	maxToken := &valToken
	*maxToken = 300

	var valTemp float64
	temp := &valTemp
	*temp = 0.9

	var valK int
	k := &valK
	*k = 0

	options := cohere.GenerateOptions{
        Model:             "command",
        Prompt:            prompt,
        MaxTokens:         maxToken,
        Temperature:       temp,
        K:                 k,
        StopSequences:     []string{},
        ReturnLikelihoods: "NONE",
    }

	response, err := b.cohereClient.Generate(options)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return response.Generations[0].Text, nil
}

func message(update *echotron.Update) string {
	if update.Message != nil {
		return update.Message.Text
	} else if update.EditedMessage != nil {
		return update.EditedMessage.Text
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Data
	}
	return ""
}
