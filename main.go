package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/CyanCipher/discogem-go/gemini"
	"github.com/CyanCipher/discogem-go/pygon"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dg, err := discordgo.New("Bot " + os.Getenv("TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMessageReactions | discordgo.IntentsGuilds | discordgo.IntentsMessageContent

	dg.AddHandler(onMessageCreate)
	dg.AddHandler(onSlashCommand)

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Bot is now active...")

	registerSlashCommands(dg)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	dg.Close()
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if strings.HasPrefix(m.Content, "AryaInt") {
		prompt := strings.Trim(m.Content, "AryaInt ")
		response, err := gemini.AskGemini(prompt)
		if err != nil {
			log.Println("Error occured! ", err)
		}
		s.ChannelMessageSend(m.ChannelID, response)
	}
}

func onSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "ask":
		if len(i.ApplicationCommandData().Options) > 0 {
			prompt := i.ApplicationCommandData().Options[0].StringValue()
			response, err := gemini.AskGemini(prompt)
			if err != nil {
				log.Println("Error occured! ", err)
			}
			mention := "<@" + i.Member.User.ID + "> "
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: mention + response,
					AllowedMentions: &discordgo.MessageAllowedMentions{
						Parse: []discordgo.AllowedMentionType{
							discordgo.AllowedMentionTypeUsers,
						},
					},
				},
			})
		}

	case "imagine":
		if len(i.ApplicationCommandData().Options) > 0 {
			prompt := i.ApplicationCommandData().Options[0].StringValue()

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
			})
			_, err := pygon.GenImage(prompt)
			if err != nil {
				log.Println("Error occured! ", err)
				mention := "<@" + i.Member.User.ID + "> "
				s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
					Content: mention + "Couldn't generate image for that prompt!",
					AllowedMentions: &discordgo.MessageAllowedMentions{
						Parse: []discordgo.AllowedMentionType{
							discordgo.AllowedMentionTypeUsers,
						},
					},
				})
				return
			}
			filereader, err := os.Open("./Media/image.png")
			if err != nil {
				log.Println("Error occured! ", err)
			}

			mention := "<@" + i.Member.User.ID + "> "
			s.FollowupMessageCreate(i.Interaction, false, &discordgo.WebhookParams{
				Files: []*discordgo.File{
					{
						Name:   "image.png",
						Reader: filereader,
					},
				},
				Content: mention + "Here is generated image ✨",
				AllowedMentions: &discordgo.MessageAllowedMentions{
					Parse: []discordgo.AllowedMentionType{
						discordgo.AllowedMentionTypeUsers,
					},
				},
			})

		}
	}
}

func registerSlashCommands(s *discordgo.Session) {
	_, err := s.ApplicationCommandCreate(s.State.User.ID, "938811985565474877", &discordgo.ApplicationCommand{
		Name:        "ask",
		Description: "Ask the AryaInt",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "What is your question?",
				Required:    true,
			},
		},
	})
	if err != nil {
		log.Println("Error occured! ", err)
	}

	_, errx := s.ApplicationCommandCreate(s.State.User.ID, "938811985565474877", &discordgo.ApplicationCommand{
		Name:        "imagine",
		Description: "Generate an image",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "prompt",
				Description: "Describe image",
				Required:    true,
			},
		},
	})

	if errx != nil {
		log.Println("Error occured! ", errx)
	}
}
