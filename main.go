// main.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/boltdb/bolt"
	. "github.com/ulvham/helper"
	"golang.org/x/net/proxy"
)

const (
	api         = ""
	telegramUrl = "https://api.telegram.org/bot"
)

type Action struct {
	Upd        *UpdateReturn
	Msg        *SendMessageReturn
	Bolt       *bolt.DB
	ProxyUsage bool
	ProxyUrl   string
}

type ActionDo interface {
	getUpdates()
	sendMessage()
	answerInlineQuery()
	answerCallbackQuery()
}

type UpdateReturn struct {
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"from"`
			Date int `json:"date"`
			Chat struct {
				ID                          int    `json:"id"`
				Type                        string `json:"type"`
				Title                       string `json:"title"`
				Username                    string `json:"username"`
				FirstName                   string `json:"first_name"`
				LastName                    string `json:"last_name"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
				Photo                       struct {
					SmallFileID string `json:"small_file_id"`
					BigFileID   string `json:"big_file_id"`
				} `json:"photo"`
				Description      string `json:"description"`
				InviteLink       string `json:"invite_link"`
				StickerSetName   string `json:"sticker_set_name"`
				CanSetStickerSet bool   `json:"can_set_sticker_set"`
			} `json:"chat"`
			ForwardFrom struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"forward_from"`
			ForwardFromChat struct {
				ID                          int    `json:"id"`
				Type                        string `json:"type"`
				Title                       string `json:"title"`
				Username                    string `json:"username"`
				FirstName                   string `json:"first_name"`
				LastName                    string `json:"last_name"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
				Photo                       struct {
					SmallFileID string `json:"small_file_id"`
					BigFileID   string `json:"big_file_id"`
				} `json:"photo"`
				Description      string `json:"description"`
				InviteLink       string `json:"invite_link"`
				StickerSetName   string `json:"sticker_set_name"`
				CanSetStickerSet bool   `json:"can_set_sticker_set"`
			} `json:"forward_from_chat"`
			ForwardFromMessageID int    `json:"forward_from_message_id"`
			ForwardDate          int    `json:"forward_date"`
			EditDate             int    `json:"edit_date"`
			Text                 string `json:"text"`
			Entities             []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"entities"`
			CaptionEntities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"caption_entities"`
			Audio struct {
				FileID    string `json:"file_id"`
				Duration  int    `json:"duration"`
				Performer string `json:"performer"`
				Title     string `json:"title"`
				MimeType  string `json:"mime_type"`
				FileSize  int    `json:"file_size"`
			} `json:"audio"`
			Document struct {
				FileID string `json:"file_id"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileName string `json:"file_name"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"document"`
			Game struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Photo       []struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"photo"`
				Text         string `json:"text"`
				TextEntities []struct {
					Type   string `json:"type"`
					Offset int    `json:"offset"`
					Length int    `json:"length"`
					URL    string `json:"url"`
					User   struct {
						ID           int    `json:"id"`
						Username     string `json:"username"`
						FirstName    string `json:"first_name"`
						LastName     string `json:"last_name"`
						LanguageCode string `json:"language_code"`
						IsBot        bool   `json:"is_bot"`
					} `json:"user"`
				} `json:"text_entities"`
				Animation struct {
					FileID string `json:"file_id"`
					Thumb  struct {
						FileID   string `json:"file_id"`
						Width    int    `json:"width"`
						Height   int    `json:"height"`
						FileSize int    `json:"file_size"`
					} `json:"thumb"`
					FileName string `json:"file_name"`
					MimeType string `json:"mime_type"`
					FileSize int    `json:"file_size"`
				} `json:"animation"`
			} `json:"game"`
			Photo []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"photo"`
			Sticker struct {
				FileID string `json:"file_id"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				Emoji        string `json:"emoji"`
				SetName      string `json:"set_name"`
				MaskPosition struct {
					Point  string `json:"point"`
					XShift int    `json:"x_shift"`
					YShift int    `json:"y_shift"`
					Zoom   int    `json:"zoom"`
				} `json:"mask_position"`
				FileSize int `json:"file_size"`
			} `json:"sticker"`
			Video struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				Duration int    `json:"duration"`
				Thumb    struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"video"`
			Voice struct {
				FileID   string `json:"file_id"`
				Duration int    `json:"duration"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"voice"`
			VideoNote struct {
				FileID   string `json:"file_id"`
				Length   int    `json:"length"`
				Duration int    `json:"duration"`
				Thumb    struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileSize int `json:"file_size"`
			} `json:"video_note"`
			Caption string `json:"caption"`
			Contact struct {
				PhoneNumber string `json:"phone_number"`
				FirstName   string `json:"first_name"`
				LastName    string `json:"last_name"`
				UserID      int    `json:"user_id"`
			} `json:"contact"`
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Venue struct {
				Location struct {
					Longitude int `json:"longitude"`
					Latitude  int `json:"latitude"`
				} `json:"location"`
				Title        string `json:"title"`
				Address      string `json:"address"`
				FoursquareID string `json:"foursquare_id"`
			} `json:"venue"`
			NewChatMembers []struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"new_chat_members"`
			LeftChatMember struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"left_chat_member"`
			NewChatTitle string `json:"new_chat_title"`
			NewChatPhoto []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"new_chat_photo"`
			DeleteChatPhoto       bool `json:"delete_chat_photo"`
			GroupChatCreated      bool `json:"group_chat_created"`
			SupergroupChatCreated bool `json:"supergroup_chat_created"`
			ChannelChatCreated    bool `json:"channel_chat_created"`
			MigrateToChatID       int  `json:"migrate_to_chat_id"`
			MigrateFromChatID     int  `json:"migrate_from_chat_id"`
			Invoice               struct {
				Title          string `json:"title"`
				Description    string `json:"description"`
				StartParameter string `json:"start_parameter"`
				Currency       string `json:"currency"`
				TotalAmount    int    `json:"total_amount"`
			} `json:"invoice"`
			SuccessfulPayment struct {
				Currency         string `json:"currency"`
				TotalAmount      int    `json:"total_amount"`
				InvoicePayload   string `json:"invoice_payload"`
				ShippingOptionID string `json:"shipping_option_id"`
				OrderInfo        struct {
					Name            string `json:"name"`
					PhoneNumber     string `json:"phone_number"`
					Email           string `json:"email"`
					ShippingAddress struct {
						CountryCode string `json:"country_code"`
						Stat        string `json:"stat"`
						City        string `json:"city"`
						StreetLine1 string `json:"street_line1"`
						StreetLine2 string `json:"street_line2"`
						PostCode    string `json:"post_code"`
					} `json:"shipping_address"`
				} `json:"order_info"`
				TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
				ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
			} `json:"successful_payment"`
			ForwardSignature string `json:"forward_signature"`
			AuthorSignature  string `json:"author_signature"`
			ConnectedWebsite string `json:"connected_website"`
		} `json:"message"`
		EditedMessage struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"from"`
			Date int `json:"date"`
			Chat struct {
				ID                          int    `json:"id"`
				Type                        string `json:"type"`
				Title                       string `json:"title"`
				Username                    string `json:"username"`
				FirstName                   string `json:"first_name"`
				LastName                    string `json:"last_name"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
				Photo                       struct {
					SmallFileID string `json:"small_file_id"`
					BigFileID   string `json:"big_file_id"`
				} `json:"photo"`
				Description      string `json:"description"`
				InviteLink       string `json:"invite_link"`
				StickerSetName   string `json:"sticker_set_name"`
				CanSetStickerSet bool   `json:"can_set_sticker_set"`
			} `json:"chat"`
			ForwardFrom struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"forward_from"`
			ForwardFromChat struct {
				ID                          int    `json:"id"`
				Type                        string `json:"type"`
				Title                       string `json:"title"`
				Username                    string `json:"username"`
				FirstName                   string `json:"first_name"`
				LastName                    string `json:"last_name"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
				Photo                       struct {
					SmallFileID string `json:"small_file_id"`
					BigFileID   string `json:"big_file_id"`
				} `json:"photo"`
				Description      string `json:"description"`
				InviteLink       string `json:"invite_link"`
				StickerSetName   string `json:"sticker_set_name"`
				CanSetStickerSet bool   `json:"can_set_sticker_set"`
			} `json:"forward_from_chat"`
			ForwardFromMessageID int    `json:"forward_from_message_id"`
			ForwardDate          int    `json:"forward_date"`
			EditDate             int    `json:"edit_date"`
			Text                 string `json:"text"`
			Entities             []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"entities"`
			CaptionEntities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"caption_entities"`
			Audio struct {
				FileID    string `json:"file_id"`
				Duration  int    `json:"duration"`
				Performer string `json:"performer"`
				Title     string `json:"title"`
				MimeType  string `json:"mime_type"`
				FileSize  int    `json:"file_size"`
			} `json:"audio"`
			Document struct {
				FileID string `json:"file_id"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileName string `json:"file_name"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"document"`
			Game struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Photo       []struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"photo"`
				Text         string `json:"text"`
				TextEntities []struct {
					Type   string `json:"type"`
					Offset int    `json:"offset"`
					Length int    `json:"length"`
					URL    string `json:"url"`
					User   struct {
						ID           int    `json:"id"`
						Username     string `json:"username"`
						FirstName    string `json:"first_name"`
						LastName     string `json:"last_name"`
						LanguageCode string `json:"language_code"`
						IsBot        bool   `json:"is_bot"`
					} `json:"user"`
				} `json:"text_entities"`
				Animation struct {
					FileID string `json:"file_id"`
					Thumb  struct {
						FileID   string `json:"file_id"`
						Width    int    `json:"width"`
						Height   int    `json:"height"`
						FileSize int    `json:"file_size"`
					} `json:"thumb"`
					FileName string `json:"file_name"`
					MimeType string `json:"mime_type"`
					FileSize int    `json:"file_size"`
				} `json:"animation"`
			} `json:"game"`
			Photo []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"photo"`
			Sticker struct {
				FileID string `json:"file_id"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				Emoji        string `json:"emoji"`
				SetName      string `json:"set_name"`
				MaskPosition struct {
					Point  string `json:"point"`
					XShift int    `json:"x_shift"`
					YShift int    `json:"y_shift"`
					Zoom   int    `json:"zoom"`
				} `json:"mask_position"`
				FileSize int `json:"file_size"`
			} `json:"sticker"`
			Video struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				Duration int    `json:"duration"`
				Thumb    struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"video"`
			Voice struct {
				FileID   string `json:"file_id"`
				Duration int    `json:"duration"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"voice"`
			VideoNote struct {
				FileID   string `json:"file_id"`
				Length   int    `json:"length"`
				Duration int    `json:"duration"`
				Thumb    struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileSize int `json:"file_size"`
			} `json:"video_note"`
			Caption string `json:"caption"`
			Contact struct {
				PhoneNumber string `json:"phone_number"`
				FirstName   string `json:"first_name"`
				LastName    string `json:"last_name"`
				UserID      int    `json:"user_id"`
			} `json:"contact"`
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Venue struct {
				Location struct {
					Longitude int `json:"longitude"`
					Latitude  int `json:"latitude"`
				} `json:"location"`
				Title        string `json:"title"`
				Address      string `json:"address"`
				FoursquareID string `json:"foursquare_id"`
			} `json:"venue"`
			NewChatMembers []struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"new_chat_members"`
			LeftChatMember struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"left_chat_member"`
			NewChatTitle string `json:"new_chat_title"`
			NewChatPhoto []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"new_chat_photo"`
			DeleteChatPhoto       bool `json:"delete_chat_photo"`
			GroupChatCreated      bool `json:"group_chat_created"`
			SupergroupChatCreated bool `json:"supergroup_chat_created"`
			ChannelChatCreated    bool `json:"channel_chat_created"`
			MigrateToChatID       int  `json:"migrate_to_chat_id"`
			MigrateFromChatID     int  `json:"migrate_from_chat_id"`
			Invoice               struct {
				Title          string `json:"title"`
				Description    string `json:"description"`
				StartParameter string `json:"start_parameter"`
				Currency       string `json:"currency"`
				TotalAmount    int    `json:"total_amount"`
			} `json:"invoice"`
			SuccessfulPayment struct {
				Currency         string `json:"currency"`
				TotalAmount      int    `json:"total_amount"`
				InvoicePayload   string `json:"invoice_payload"`
				ShippingOptionID string `json:"shipping_option_id"`
				OrderInfo        struct {
					Name            string `json:"name"`
					PhoneNumber     string `json:"phone_number"`
					Email           string `json:"email"`
					ShippingAddress struct {
						CountryCode string `json:"country_code"`
						Stat        string `json:"stat"`
						City        string `json:"city"`
						StreetLine1 string `json:"street_line1"`
						StreetLine2 string `json:"street_line2"`
						PostCode    string `json:"post_code"`
					} `json:"shipping_address"`
				} `json:"order_info"`
				TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
				ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
			} `json:"successful_payment"`
			ForwardSignature string `json:"forward_signature"`
			AuthorSignature  string `json:"author_signature"`
			ConnectedWebsite string `json:"connected_website"`
		} `json:"edited_message"`
		ChannelPost struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"from"`
			Date int `json:"date"`
			Chat struct {
				ID                          int    `json:"id"`
				Type                        string `json:"type"`
				Title                       string `json:"title"`
				Username                    string `json:"username"`
				FirstName                   string `json:"first_name"`
				LastName                    string `json:"last_name"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
				Photo                       struct {
					SmallFileID string `json:"small_file_id"`
					BigFileID   string `json:"big_file_id"`
				} `json:"photo"`
				Description      string `json:"description"`
				InviteLink       string `json:"invite_link"`
				StickerSetName   string `json:"sticker_set_name"`
				CanSetStickerSet bool   `json:"can_set_sticker_set"`
			} `json:"chat"`
			ForwardFrom struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"forward_from"`
			ForwardFromChat struct {
				ID                          int    `json:"id"`
				Type                        string `json:"type"`
				Title                       string `json:"title"`
				Username                    string `json:"username"`
				FirstName                   string `json:"first_name"`
				LastName                    string `json:"last_name"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
				Photo                       struct {
					SmallFileID string `json:"small_file_id"`
					BigFileID   string `json:"big_file_id"`
				} `json:"photo"`
				Description      string `json:"description"`
				InviteLink       string `json:"invite_link"`
				StickerSetName   string `json:"sticker_set_name"`
				CanSetStickerSet bool   `json:"can_set_sticker_set"`
			} `json:"forward_from_chat"`
			ForwardFromMessageID int    `json:"forward_from_message_id"`
			ForwardDate          int    `json:"forward_date"`
			EditDate             int    `json:"edit_date"`
			Text                 string `json:"text"`
			Entities             []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"entities"`
			CaptionEntities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"caption_entities"`
			Audio struct {
				FileID    string `json:"file_id"`
				Duration  int    `json:"duration"`
				Performer string `json:"performer"`
				Title     string `json:"title"`
				MimeType  string `json:"mime_type"`
				FileSize  int    `json:"file_size"`
			} `json:"audio"`
			Document struct {
				FileID string `json:"file_id"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileName string `json:"file_name"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"document"`
			Game struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Photo       []struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"photo"`
				Text         string `json:"text"`
				TextEntities []struct {
					Type   string `json:"type"`
					Offset int    `json:"offset"`
					Length int    `json:"length"`
					URL    string `json:"url"`
					User   struct {
						ID           int    `json:"id"`
						Username     string `json:"username"`
						FirstName    string `json:"first_name"`
						LastName     string `json:"last_name"`
						LanguageCode string `json:"language_code"`
						IsBot        bool   `json:"is_bot"`
					} `json:"user"`
				} `json:"text_entities"`
				Animation struct {
					FileID string `json:"file_id"`
					Thumb  struct {
						FileID   string `json:"file_id"`
						Width    int    `json:"width"`
						Height   int    `json:"height"`
						FileSize int    `json:"file_size"`
					} `json:"thumb"`
					FileName string `json:"file_name"`
					MimeType string `json:"mime_type"`
					FileSize int    `json:"file_size"`
				} `json:"animation"`
			} `json:"game"`
			Photo []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"photo"`
			Sticker struct {
				FileID string `json:"file_id"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				Emoji        string `json:"emoji"`
				SetName      string `json:"set_name"`
				MaskPosition struct {
					Point  string `json:"point"`
					XShift int    `json:"x_shift"`
					YShift int    `json:"y_shift"`
					Zoom   int    `json:"zoom"`
				} `json:"mask_position"`
				FileSize int `json:"file_size"`
			} `json:"sticker"`
			Video struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				Duration int    `json:"duration"`
				Thumb    struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"video"`
			Voice struct {
				FileID   string `json:"file_id"`
				Duration int    `json:"duration"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"voice"`
			VideoNote struct {
				FileID   string `json:"file_id"`
				Length   int    `json:"length"`
				Duration int    `json:"duration"`
				Thumb    struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileSize int `json:"file_size"`
			} `json:"video_note"`
			Caption string `json:"caption"`
			Contact struct {
				PhoneNumber string `json:"phone_number"`
				FirstName   string `json:"first_name"`
				LastName    string `json:"last_name"`
				UserID      int    `json:"user_id"`
			} `json:"contact"`
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Venue struct {
				Location struct {
					Longitude int `json:"longitude"`
					Latitude  int `json:"latitude"`
				} `json:"location"`
				Title        string `json:"title"`
				Address      string `json:"address"`
				FoursquareID string `json:"foursquare_id"`
			} `json:"venue"`
			NewChatMembers []struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"new_chat_members"`
			LeftChatMember struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"left_chat_member"`
			NewChatTitle string `json:"new_chat_title"`
			NewChatPhoto []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"new_chat_photo"`
			DeleteChatPhoto       bool `json:"delete_chat_photo"`
			GroupChatCreated      bool `json:"group_chat_created"`
			SupergroupChatCreated bool `json:"supergroup_chat_created"`
			ChannelChatCreated    bool `json:"channel_chat_created"`
			MigrateToChatID       int  `json:"migrate_to_chat_id"`
			MigrateFromChatID     int  `json:"migrate_from_chat_id"`
			Invoice               struct {
				Title          string `json:"title"`
				Description    string `json:"description"`
				StartParameter string `json:"start_parameter"`
				Currency       string `json:"currency"`
				TotalAmount    int    `json:"total_amount"`
			} `json:"invoice"`
			SuccessfulPayment struct {
				Currency         string `json:"currency"`
				TotalAmount      int    `json:"total_amount"`
				InvoicePayload   string `json:"invoice_payload"`
				ShippingOptionID string `json:"shipping_option_id"`
				OrderInfo        struct {
					Name            string `json:"name"`
					PhoneNumber     string `json:"phone_number"`
					Email           string `json:"email"`
					ShippingAddress struct {
						CountryCode string `json:"country_code"`
						Stat        string `json:"stat"`
						City        string `json:"city"`
						StreetLine1 string `json:"street_line1"`
						StreetLine2 string `json:"street_line2"`
						PostCode    string `json:"post_code"`
					} `json:"shipping_address"`
				} `json:"order_info"`
				TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
				ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
			} `json:"successful_payment"`
			ForwardSignature string `json:"forward_signature"`
			AuthorSignature  string `json:"author_signature"`
			ConnectedWebsite string `json:"connected_website"`
		} `json:"channel_post"`
		EditedChannelPost struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"from"`
			Date int `json:"date"`
			Chat struct {
				ID                          int    `json:"id"`
				Type                        string `json:"type"`
				Title                       string `json:"title"`
				Username                    string `json:"username"`
				FirstName                   string `json:"first_name"`
				LastName                    string `json:"last_name"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
				Photo                       struct {
					SmallFileID string `json:"small_file_id"`
					BigFileID   string `json:"big_file_id"`
				} `json:"photo"`
				Description      string `json:"description"`
				InviteLink       string `json:"invite_link"`
				StickerSetName   string `json:"sticker_set_name"`
				CanSetStickerSet bool   `json:"can_set_sticker_set"`
			} `json:"chat"`
			ForwardFrom struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"forward_from"`
			ForwardFromChat struct {
				ID                          int    `json:"id"`
				Type                        string `json:"type"`
				Title                       string `json:"title"`
				Username                    string `json:"username"`
				FirstName                   string `json:"first_name"`
				LastName                    string `json:"last_name"`
				AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
				Photo                       struct {
					SmallFileID string `json:"small_file_id"`
					BigFileID   string `json:"big_file_id"`
				} `json:"photo"`
				Description      string `json:"description"`
				InviteLink       string `json:"invite_link"`
				StickerSetName   string `json:"sticker_set_name"`
				CanSetStickerSet bool   `json:"can_set_sticker_set"`
			} `json:"forward_from_chat"`
			ForwardFromMessageID int    `json:"forward_from_message_id"`
			ForwardDate          int    `json:"forward_date"`
			EditDate             int    `json:"edit_date"`
			Text                 string `json:"text"`
			Entities             []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"entities"`
			CaptionEntities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"caption_entities"`
			Audio struct {
				FileID    string `json:"file_id"`
				Duration  int    `json:"duration"`
				Performer string `json:"performer"`
				Title     string `json:"title"`
				MimeType  string `json:"mime_type"`
				FileSize  int    `json:"file_size"`
			} `json:"audio"`
			Document struct {
				FileID string `json:"file_id"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileName string `json:"file_name"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"document"`
			Game struct {
				Title       string `json:"title"`
				Description string `json:"description"`
				Photo       []struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"photo"`
				Text         string `json:"text"`
				TextEntities []struct {
					Type   string `json:"type"`
					Offset int    `json:"offset"`
					Length int    `json:"length"`
					URL    string `json:"url"`
					User   struct {
						ID           int    `json:"id"`
						Username     string `json:"username"`
						FirstName    string `json:"first_name"`
						LastName     string `json:"last_name"`
						LanguageCode string `json:"language_code"`
						IsBot        bool   `json:"is_bot"`
					} `json:"user"`
				} `json:"text_entities"`
				Animation struct {
					FileID string `json:"file_id"`
					Thumb  struct {
						FileID   string `json:"file_id"`
						Width    int    `json:"width"`
						Height   int    `json:"height"`
						FileSize int    `json:"file_size"`
					} `json:"thumb"`
					FileName string `json:"file_name"`
					MimeType string `json:"mime_type"`
					FileSize int    `json:"file_size"`
				} `json:"animation"`
			} `json:"game"`
			Photo []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"photo"`
			Sticker struct {
				FileID string `json:"file_id"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				Emoji        string `json:"emoji"`
				SetName      string `json:"set_name"`
				MaskPosition struct {
					Point  string `json:"point"`
					XShift int    `json:"x_shift"`
					YShift int    `json:"y_shift"`
					Zoom   int    `json:"zoom"`
				} `json:"mask_position"`
				FileSize int `json:"file_size"`
			} `json:"sticker"`
			Video struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				Duration int    `json:"duration"`
				Thumb    struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"video"`
			Voice struct {
				FileID   string `json:"file_id"`
				Duration int    `json:"duration"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"voice"`
			VideoNote struct {
				FileID   string `json:"file_id"`
				Length   int    `json:"length"`
				Duration int    `json:"duration"`
				Thumb    struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileSize int `json:"file_size"`
			} `json:"video_note"`
			Caption string `json:"caption"`
			Contact struct {
				PhoneNumber string `json:"phone_number"`
				FirstName   string `json:"first_name"`
				LastName    string `json:"last_name"`
				UserID      int    `json:"user_id"`
			} `json:"contact"`
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Venue struct {
				Location struct {
					Longitude int `json:"longitude"`
					Latitude  int `json:"latitude"`
				} `json:"location"`
				Title        string `json:"title"`
				Address      string `json:"address"`
				FoursquareID string `json:"foursquare_id"`
			} `json:"venue"`
			NewChatMembers []struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"new_chat_members"`
			LeftChatMember struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"left_chat_member"`
			NewChatTitle string `json:"new_chat_title"`
			NewChatPhoto []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"new_chat_photo"`
			DeleteChatPhoto       bool `json:"delete_chat_photo"`
			GroupChatCreated      bool `json:"group_chat_created"`
			SupergroupChatCreated bool `json:"supergroup_chat_created"`
			ChannelChatCreated    bool `json:"channel_chat_created"`
			MigrateToChatID       int  `json:"migrate_to_chat_id"`
			MigrateFromChatID     int  `json:"migrate_from_chat_id"`
			Invoice               struct {
				Title          string `json:"title"`
				Description    string `json:"description"`
				StartParameter string `json:"start_parameter"`
				Currency       string `json:"currency"`
				TotalAmount    int    `json:"total_amount"`
			} `json:"invoice"`
			SuccessfulPayment struct {
				Currency         string `json:"currency"`
				TotalAmount      int    `json:"total_amount"`
				InvoicePayload   string `json:"invoice_payload"`
				ShippingOptionID string `json:"shipping_option_id"`
				OrderInfo        struct {
					Name            string `json:"name"`
					PhoneNumber     string `json:"phone_number"`
					Email           string `json:"email"`
					ShippingAddress struct {
						CountryCode string `json:"country_code"`
						Stat        string `json:"stat"`
						City        string `json:"city"`
						StreetLine1 string `json:"street_line1"`
						StreetLine2 string `json:"street_line2"`
						PostCode    string `json:"post_code"`
					} `json:"shipping_address"`
				} `json:"order_info"`
				TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
				ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
			} `json:"successful_payment"`
			ForwardSignature string `json:"forward_signature"`
			AuthorSignature  string `json:"author_signature"`
			ConnectedWebsite string `json:"connected_website"`
		} `json:"edited_channel_post"`
		InlineQuery struct {
			ID   string `json:"id"`
			From struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"from"`
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Query  string `json:"query"`
			Offset string `json:"offset"`
		} `json:"inline_query"`
		ChosenInlineResult struct {
			ResultID string `json:"result_id"`
			From     struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"from"`
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			InlineMessageID string `json:"inline_message_id"`
			Query           string `json:"query"`
		} `json:"chosen_inline_result"`
		CallbackQuery struct {
			ID   string `json:"id"`
			From struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"from"`
			Message struct {
				MessageID int `json:"message_id"`
				From      struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"from"`
				Date int `json:"date"`
				Chat struct {
					ID                          int    `json:"id"`
					Type                        string `json:"type"`
					Title                       string `json:"title"`
					Username                    string `json:"username"`
					FirstName                   string `json:"first_name"`
					LastName                    string `json:"last_name"`
					AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
					Photo                       struct {
						SmallFileID string `json:"small_file_id"`
						BigFileID   string `json:"big_file_id"`
					} `json:"photo"`
					Description      string `json:"description"`
					InviteLink       string `json:"invite_link"`
					StickerSetName   string `json:"sticker_set_name"`
					CanSetStickerSet bool   `json:"can_set_sticker_set"`
				} `json:"chat"`
				ForwardFrom struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"forward_from"`
				ForwardFromChat struct {
					ID                          int    `json:"id"`
					Type                        string `json:"type"`
					Title                       string `json:"title"`
					Username                    string `json:"username"`
					FirstName                   string `json:"first_name"`
					LastName                    string `json:"last_name"`
					AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
					Photo                       struct {
						SmallFileID string `json:"small_file_id"`
						BigFileID   string `json:"big_file_id"`
					} `json:"photo"`
					Description      string `json:"description"`
					InviteLink       string `json:"invite_link"`
					StickerSetName   string `json:"sticker_set_name"`
					CanSetStickerSet bool   `json:"can_set_sticker_set"`
				} `json:"forward_from_chat"`
				ForwardFromMessageID int    `json:"forward_from_message_id"`
				ForwardDate          int    `json:"forward_date"`
				EditDate             int    `json:"edit_date"`
				Text                 string `json:"text"`
				Entities             []struct {
					Type   string `json:"type"`
					Offset int    `json:"offset"`
					Length int    `json:"length"`
					URL    string `json:"url"`
					User   struct {
						ID           int    `json:"id"`
						Username     string `json:"username"`
						FirstName    string `json:"first_name"`
						LastName     string `json:"last_name"`
						LanguageCode string `json:"language_code"`
						IsBot        bool   `json:"is_bot"`
					} `json:"user"`
				} `json:"entities"`
				CaptionEntities []struct {
					Type   string `json:"type"`
					Offset int    `json:"offset"`
					Length int    `json:"length"`
					URL    string `json:"url"`
					User   struct {
						ID           int    `json:"id"`
						Username     string `json:"username"`
						FirstName    string `json:"first_name"`
						LastName     string `json:"last_name"`
						LanguageCode string `json:"language_code"`
						IsBot        bool   `json:"is_bot"`
					} `json:"user"`
				} `json:"caption_entities"`
				Audio struct {
					FileID    string `json:"file_id"`
					Duration  int    `json:"duration"`
					Performer string `json:"performer"`
					Title     string `json:"title"`
					MimeType  string `json:"mime_type"`
					FileSize  int    `json:"file_size"`
				} `json:"audio"`
				Document struct {
					FileID string `json:"file_id"`
					Thumb  struct {
						FileID   string `json:"file_id"`
						Width    int    `json:"width"`
						Height   int    `json:"height"`
						FileSize int    `json:"file_size"`
					} `json:"thumb"`
					FileName string `json:"file_name"`
					MimeType string `json:"mime_type"`
					FileSize int    `json:"file_size"`
				} `json:"document"`
				Game struct {
					Title       string `json:"title"`
					Description string `json:"description"`
					Photo       []struct {
						FileID   string `json:"file_id"`
						Width    int    `json:"width"`
						Height   int    `json:"height"`
						FileSize int    `json:"file_size"`
					} `json:"photo"`
					Text         string `json:"text"`
					TextEntities []struct {
						Type   string `json:"type"`
						Offset int    `json:"offset"`
						Length int    `json:"length"`
						URL    string `json:"url"`
						User   struct {
							ID           int    `json:"id"`
							Username     string `json:"username"`
							FirstName    string `json:"first_name"`
							LastName     string `json:"last_name"`
							LanguageCode string `json:"language_code"`
							IsBot        bool   `json:"is_bot"`
						} `json:"user"`
					} `json:"text_entities"`
					Animation struct {
						FileID string `json:"file_id"`
						Thumb  struct {
							FileID   string `json:"file_id"`
							Width    int    `json:"width"`
							Height   int    `json:"height"`
							FileSize int    `json:"file_size"`
						} `json:"thumb"`
						FileName string `json:"file_name"`
						MimeType string `json:"mime_type"`
						FileSize int    `json:"file_size"`
					} `json:"animation"`
				} `json:"game"`
				Photo []struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"photo"`
				Sticker struct {
					FileID string `json:"file_id"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
					Thumb  struct {
						FileID   string `json:"file_id"`
						Width    int    `json:"width"`
						Height   int    `json:"height"`
						FileSize int    `json:"file_size"`
					} `json:"thumb"`
					Emoji        string `json:"emoji"`
					SetName      string `json:"set_name"`
					MaskPosition struct {
						Point  string `json:"point"`
						XShift int    `json:"x_shift"`
						YShift int    `json:"y_shift"`
						Zoom   int    `json:"zoom"`
					} `json:"mask_position"`
					FileSize int `json:"file_size"`
				} `json:"sticker"`
				Video struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					Duration int    `json:"duration"`
					Thumb    struct {
						FileID   string `json:"file_id"`
						Width    int    `json:"width"`
						Height   int    `json:"height"`
						FileSize int    `json:"file_size"`
					} `json:"thumb"`
					MimeType string `json:"mime_type"`
					FileSize int    `json:"file_size"`
				} `json:"video"`
				Voice struct {
					FileID   string `json:"file_id"`
					Duration int    `json:"duration"`
					MimeType string `json:"mime_type"`
					FileSize int    `json:"file_size"`
				} `json:"voice"`
				VideoNote struct {
					FileID   string `json:"file_id"`
					Length   int    `json:"length"`
					Duration int    `json:"duration"`
					Thumb    struct {
						FileID   string `json:"file_id"`
						Width    int    `json:"width"`
						Height   int    `json:"height"`
						FileSize int    `json:"file_size"`
					} `json:"thumb"`
					FileSize int `json:"file_size"`
				} `json:"video_note"`
				Caption string `json:"caption"`
				Contact struct {
					PhoneNumber string `json:"phone_number"`
					FirstName   string `json:"first_name"`
					LastName    string `json:"last_name"`
					UserID      int    `json:"user_id"`
				} `json:"contact"`
				Location struct {
					Longitude int `json:"longitude"`
					Latitude  int `json:"latitude"`
				} `json:"location"`
				Venue struct {
					Location struct {
						Longitude int `json:"longitude"`
						Latitude  int `json:"latitude"`
					} `json:"location"`
					Title        string `json:"title"`
					Address      string `json:"address"`
					FoursquareID string `json:"foursquare_id"`
				} `json:"venue"`
				NewChatMembers []struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"new_chat_members"`
				LeftChatMember struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"left_chat_member"`
				NewChatTitle string `json:"new_chat_title"`
				NewChatPhoto []struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"new_chat_photo"`
				DeleteChatPhoto       bool `json:"delete_chat_photo"`
				GroupChatCreated      bool `json:"group_chat_created"`
				SupergroupChatCreated bool `json:"supergroup_chat_created"`
				ChannelChatCreated    bool `json:"channel_chat_created"`
				MigrateToChatID       int  `json:"migrate_to_chat_id"`
				MigrateFromChatID     int  `json:"migrate_from_chat_id"`
				Invoice               struct {
					Title          string `json:"title"`
					Description    string `json:"description"`
					StartParameter string `json:"start_parameter"`
					Currency       string `json:"currency"`
					TotalAmount    int    `json:"total_amount"`
				} `json:"invoice"`
				SuccessfulPayment struct {
					Currency         string `json:"currency"`
					TotalAmount      int    `json:"total_amount"`
					InvoicePayload   string `json:"invoice_payload"`
					ShippingOptionID string `json:"shipping_option_id"`
					OrderInfo        struct {
						Name            string `json:"name"`
						PhoneNumber     string `json:"phone_number"`
						Email           string `json:"email"`
						ShippingAddress struct {
							CountryCode string `json:"country_code"`
							Stat        string `json:"stat"`
							City        string `json:"city"`
							StreetLine1 string `json:"street_line1"`
							StreetLine2 string `json:"street_line2"`
							PostCode    string `json:"post_code"`
						} `json:"shipping_address"`
					} `json:"order_info"`
					TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
					ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
				} `json:"successful_payment"`
				ForwardSignature string `json:"forward_signature"`
				AuthorSignature  string `json:"author_signature"`
				ConnectedWebsite string `json:"connected_website"`
			} `json:"message"`
			InlineMessageID string `json:"inline_message_id"`
			ChatInstance    string `json:"chat_instance"`
			Data            string `json:"data"`
			GameShortName   string `json:"game_short_name"`
		} `json:"callback_query"`
		ShippingQuery struct {
			ID   string `json:"id"`
			From struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"from"`
			InvoicePayload  string `json:"invoice_payload"`
			ShippingAddress struct {
				CountryCode string `json:"country_code"`
				Stat        string `json:"stat"`
				City        string `json:"city"`
				StreetLine1 string `json:"street_line1"`
				StreetLine2 string `json:"street_line2"`
				PostCode    string `json:"post_code"`
			} `json:"shipping_address"`
		} `json:"shipping_query"`
		PreCheckoutQuery struct {
			ID   string `json:"id"`
			From struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"from"`
			Currency         string `json:"currency"`
			TotalAmount      int    `json:"total_amount"`
			InvoicePayload   string `json:"invoice_payload"`
			ShippingOptionID string `json:"shipping_option_id"`
			OrderInfo        struct {
				Name            string `json:"name"`
				PhoneNumber     string `json:"phone_number"`
				Email           string `json:"email"`
				ShippingAddress struct {
					CountryCode string `json:"country_code"`
					Stat        string `json:"stat"`
					City        string `json:"city"`
					StreetLine1 string `json:"street_line1"`
					StreetLine2 string `json:"street_line2"`
					PostCode    string `json:"post_code"`
				} `json:"shipping_address"`
			} `json:"order_info"`
		} `json:"pre_checkout_query"`
	} `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

type SendMessageReturn struct {
	Result struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"from"`
		Date int `json:"date"`
		Chat struct {
			ID                          int    `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"chat"`
		ForwardFrom struct {
			ID           int    `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"forward_from"`
		ForwardFromChat struct {
			ID                          int    `json:"id"`
			Type                        string `json:"type"`
			Title                       string `json:"title"`
			Username                    string `json:"username"`
			FirstName                   string `json:"first_name"`
			LastName                    string `json:"last_name"`
			AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
			Photo                       struct {
				SmallFileID string `json:"small_file_id"`
				BigFileID   string `json:"big_file_id"`
			} `json:"photo"`
			Description      string `json:"description"`
			InviteLink       string `json:"invite_link"`
			StickerSetName   string `json:"sticker_set_name"`
			CanSetStickerSet bool   `json:"can_set_sticker_set"`
		} `json:"forward_from_chat"`
		ForwardFromMessageID int    `json:"forward_from_message_id"`
		ForwardDate          int    `json:"forward_date"`
		EditDate             int    `json:"edit_date"`
		Text                 string `json:"text"`
		Entities             []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"entities"`
		CaptionEntities []struct {
			Type   string `json:"type"`
			Offset int    `json:"offset"`
			Length int    `json:"length"`
			URL    string `json:"url"`
			User   struct {
				ID           int    `json:"id"`
				Username     string `json:"username"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				LanguageCode string `json:"language_code"`
				IsBot        bool   `json:"is_bot"`
			} `json:"user"`
		} `json:"caption_entities"`
		Audio struct {
			FileID    string `json:"file_id"`
			Duration  int    `json:"duration"`
			Performer string `json:"performer"`
			Title     string `json:"title"`
			MimeType  string `json:"mime_type"`
			FileSize  int    `json:"file_size"`
		} `json:"audio"`
		Document struct {
			FileID string `json:"file_id"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileName string `json:"file_name"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"document"`
		Game struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			Photo       []struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"photo"`
			Text         string `json:"text"`
			TextEntities []struct {
				Type   string `json:"type"`
				Offset int    `json:"offset"`
				Length int    `json:"length"`
				URL    string `json:"url"`
				User   struct {
					ID           int    `json:"id"`
					Username     string `json:"username"`
					FirstName    string `json:"first_name"`
					LastName     string `json:"last_name"`
					LanguageCode string `json:"language_code"`
					IsBot        bool   `json:"is_bot"`
				} `json:"user"`
			} `json:"text_entities"`
			Animation struct {
				FileID string `json:"file_id"`
				Thumb  struct {
					FileID   string `json:"file_id"`
					Width    int    `json:"width"`
					Height   int    `json:"height"`
					FileSize int    `json:"file_size"`
				} `json:"thumb"`
				FileName string `json:"file_name"`
				MimeType string `json:"mime_type"`
				FileSize int    `json:"file_size"`
			} `json:"animation"`
		} `json:"game"`
		Photo []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"photo"`
		Sticker struct {
			FileID string `json:"file_id"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
			Thumb  struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			Emoji        string `json:"emoji"`
			SetName      string `json:"set_name"`
			MaskPosition struct {
				Point  string `json:"point"`
				XShift int    `json:"x_shift"`
				YShift int    `json:"y_shift"`
				Zoom   int    `json:"zoom"`
			} `json:"mask_position"`
			FileSize int `json:"file_size"`
		} `json:"sticker"`
		Video struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"video"`
		Voice struct {
			FileID   string `json:"file_id"`
			Duration int    `json:"duration"`
			MimeType string `json:"mime_type"`
			FileSize int    `json:"file_size"`
		} `json:"voice"`
		VideoNote struct {
			FileID   string `json:"file_id"`
			Length   int    `json:"length"`
			Duration int    `json:"duration"`
			Thumb    struct {
				FileID   string `json:"file_id"`
				Width    int    `json:"width"`
				Height   int    `json:"height"`
				FileSize int    `json:"file_size"`
			} `json:"thumb"`
			FileSize int `json:"file_size"`
		} `json:"video_note"`
		Caption string `json:"caption"`
		Contact struct {
			PhoneNumber string `json:"phone_number"`
			FirstName   string `json:"first_name"`
			LastName    string `json:"last_name"`
			UserID      int    `json:"user_id"`
		} `json:"contact"`
		Location struct {
			Longitude int `json:"longitude"`
			Latitude  int `json:"latitude"`
		} `json:"location"`
		Venue struct {
			Location struct {
				Longitude int `json:"longitude"`
				Latitude  int `json:"latitude"`
			} `json:"location"`
			Title        string `json:"title"`
			Address      string `json:"address"`
			FoursquareID string `json:"foursquare_id"`
		} `json:"venue"`
		NewChatMembers []struct {
			ID           int    `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"new_chat_members"`
		LeftChatMember struct {
			ID           int    `json:"id"`
			Username     string `json:"username"`
			FirstName    string `json:"first_name"`
			LastName     string `json:"last_name"`
			LanguageCode string `json:"language_code"`
			IsBot        bool   `json:"is_bot"`
		} `json:"left_chat_member"`
		NewChatTitle string `json:"new_chat_title"`
		NewChatPhoto []struct {
			FileID   string `json:"file_id"`
			Width    int    `json:"width"`
			Height   int    `json:"height"`
			FileSize int    `json:"file_size"`
		} `json:"new_chat_photo"`
		DeleteChatPhoto       bool `json:"delete_chat_photo"`
		GroupChatCreated      bool `json:"group_chat_created"`
		SupergroupChatCreated bool `json:"supergroup_chat_created"`
		ChannelChatCreated    bool `json:"channel_chat_created"`
		MigrateToChatID       int  `json:"migrate_to_chat_id"`
		MigrateFromChatID     int  `json:"migrate_from_chat_id"`
		Invoice               struct {
			Title          string `json:"title"`
			Description    string `json:"description"`
			StartParameter string `json:"start_parameter"`
			Currency       string `json:"currency"`
			TotalAmount    int    `json:"total_amount"`
		} `json:"invoice"`
		SuccessfulPayment struct {
			Currency         string `json:"currency"`
			TotalAmount      int    `json:"total_amount"`
			InvoicePayload   string `json:"invoice_payload"`
			ShippingOptionID string `json:"shipping_option_id"`
			OrderInfo        struct {
				Name            string `json:"name"`
				PhoneNumber     string `json:"phone_number"`
				Email           string `json:"email"`
				ShippingAddress struct {
					CountryCode string `json:"country_code"`
					Stat        string `json:"stat"`
					City        string `json:"city"`
					StreetLine1 string `json:"street_line1"`
					StreetLine2 string `json:"street_line2"`
					PostCode    string `json:"post_code"`
				} `json:"shipping_address"`
			} `json:"order_info"`
			TelegramPaymentChargeID string `json:"telegram_payment_charge_id"`
			ProviderPaymentChargeID string `json:"provider_payment_charge_id"`
		} `json:"successful_payment"`
		ForwardSignature string `json:"forward_signature"`
		AuthorSignature  string `json:"author_signature"`
		ConnectedWebsite string `json:"connected_website"`
	} `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

type InlineReturn struct {
	Result      bool   `json:"result"`
	ErrorCode   int    `json:"error_code"`
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

type PayloadGetUpdates struct {
	Offset         int      `json:"offset"`
	Limit          int      `json:"limit"`
	Timeout        int      `json:"timeout"`
	AllowedUpdates []string `json:"allowed_updates"`
}

type Button struct {
	InlineKeyboard [][]Button_ `json:"inline_keyboard"`
}

type Button_ struct {
	Text         string `json:"text"`
	Url          string `json:"url"`
	CallbackData string `json:"callback_data"`
}

type PayloadMesageSend struct {
	ChatID                int    `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview"`
	DisableNotification   bool   `json:"disable_notification"`
	ReplyToMessageID      int    `json:"reply_to_message_id"`
	ReplyMarkup           Button `json:"reply_markup"`
}

func (obj *Action) getUpdates() {
	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	if obj.ProxyUsage {
		dialer, err := proxy.SOCKS5("tcp", obj.ProxyUrl, nil, proxy.Direct)
		Dbg(err)
		httpTransport.Dial = dialer.Dial
	}

	data := PayloadGetUpdates{}
	data.Timeout = 1
	data.Limit = 5
	data.Offset = -5
	data.AllowedUpdates = append(data.AllowedUpdates, "message")
	data.AllowedUpdates = append(data.AllowedUpdates, "callback_query")
	data.AllowedUpdates = append(data.AllowedUpdates, "inline_query")

	payloadBytes, err := json.Marshal(data)
	Dbg(err)

	fmt.Println(string(payloadBytes))
	body := strings.NewReader(string(payloadBytes))

	req, err := http.NewRequest("POST", telegramUrl+api+"/getUpdates", body)
	Dbg(err)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	Dbg(err)
	defer resp.Body.Close()

	bodyret, _ := ioutil.ReadAll(resp.Body)
	ret := new(UpdateReturn)
	json.Unmarshal(bodyret, ret)
	obj.Upd = ret
	fmt.Println(string(bodyret))
}

func (obj *Action) sendMessage() {
	var wg sync.WaitGroup

	httpTransport := &http.Transport{}
	httpClient := &http.Client{Transport: httpTransport}
	if obj.ProxyUsage {
		dialer, err := proxy.SOCKS5("tcp", obj.ProxyUrl, nil, proxy.Direct)
		Dbg(err)
		httpTransport.Dial = dialer.Dial
	}
	for _, val := range obj.Upd.Result {
		exists := false
		obj.Bolt.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("Get"))
			v := b.Get([]byte("key" + ToStr(val.Message.MessageID)))
			if v != nil {
				exists = true
			}
			return nil
		})
		if exists {
			continue
		}
		wg.Add(1)
		go func() {
			obj.Bolt.Update(func(tx *bolt.Tx) error {
				b, _ := tx.CreateBucketIfNotExists([]byte("Get"))
				err := b.Put([]byte("key"+ToStr(val.Message.MessageID)), []byte(ToStr(val.Message.Text)))
				return err
			})
			wg.Done()
		}()
		req0, err := http.NewRequest("GET", telegramUrl+api+"/sendChatAction?chat_id="+ToStr(val.Message.Chat.ID)+"&action=typing", nil)
		Dbg(err)
		req0.Header.Set("Accept", "application/json")
		resp0, err := httpClient.Do(req0)
		Dbg(err)
		defer resp0.Body.Close()
		data := PayloadMesageSend{}
		data.ChatID = val.Message.Chat.ID
		//data.ReplyToMessageID = val.Message.MessageID
		data.Text = val.Message.Text

		//type Button struct {
		//	Text         string `json:"text"`
		//	CallbackData string `json:"callback_data"`
		//}

		but := Button{}
		but1 := []Button_{}
		but2 := [][]Button_{}
		but0 := Button_{Text: "yes", CallbackData: "yes my boy!"} //Url: `https://www.myqnapcloud.com/smartshare/6d10h8h5np2m2612251994ya_6ZPL31G`
		but00 := Button_{Text: "no", CallbackData: "no my boy?"}  //Url: `https://www.myqnapcloud.com/smartshare/6d10h8h5np2m2612251994ya_6ZPL31G`
		but1 = append(but1, but0)
		but1 = append(but1, but00)
		but2 = append(but2, but1)

		but.InlineKeyboard = but2

		//		var buts []Button
		//		var buts_ [][]Button

		//		buts = append(buts, but)
		//		buts_ = append(buts_, buts)

		data.ReplyMarkup = but

		fmt.Println(data)

		payloadBytes, err := json.Marshal(data)
		Dbg(err)
		body := bytes.NewReader(payloadBytes)
		req, err := http.NewRequest("POST", telegramUrl+api+"/sendMessage", body)
		Dbg(err)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		Dbg(err)
		defer resp.Body.Close()
		bodyret, _ := ioutil.ReadAll(resp.Body)
		ret := new(SendMessageReturn)
		json.Unmarshal(bodyret, ret)
	}
	wg.Wait()
}

func (obj *Action) answerCallbackQuery() {
	for _, val := range obj.Upd.Result {
		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}
		if obj.ProxyUsage {
			dialer, err := proxy.SOCKS5("tcp", obj.ProxyUrl, nil, proxy.Direct)
			Dbg(err)
			httpTransport.Dial = dialer.Dial
		}

		type Payload struct {
			CallbackQueryId string `json:"callback_query_id"`
			Text            string `json:"text"`
			ShowAlert       bool   `json:"show_alert"`
			Url             string `json:"url"`
			CacheTime       int    `json:"cache_time"`
		}

		data := Payload{
			// fill struct
		}
		//res_ := []string{"test1", "test2"}
		fmt.Println(val.CallbackQuery.ID)
		fmt.Println(val.CallbackQuery.Data)
		data.CallbackQueryId = val.CallbackQuery.ID
		data.Text = val.CallbackQuery.Data
		//data.Results = "1111111111111111111111"

		payloadBytes, err := json.Marshal(data)
		if err != nil {
			// handle err
		}
		body := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", telegramUrl+api+"/answerCallbackQuery", body)
		if err != nil {
			// handle err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			// handle err
		}
		defer resp.Body.Close()

		bodyret, _ := ioutil.ReadAll(resp.Body)
		ret := new(InlineReturn)
		json.Unmarshal(bodyret, ret)

		fmt.Println(string(bodyret))
		//fmt.Println(ret)
	}
}

func (obj *Action) answerInlineQuery() {
	for _, val := range obj.Upd.Result {
		httpTransport := &http.Transport{}
		httpClient := &http.Client{Transport: httpTransport}
		if obj.ProxyUsage {
			dialer, err := proxy.SOCKS5("tcp", obj.ProxyUrl, nil, proxy.Direct)
			Dbg(err)
			httpTransport.Dial = dialer.Dial
		}

		type Payload struct {
			InlineQueryID     string `json:"inline_query_id"`
			Results           string
			CacheTime         int    `json:"cache_time"`
			IsPersonal        bool   `json:"is_personal"`
			NextOffset        string `json:"next_offset"`
			SwitchPmText      string `json:"switch_pm_text"`
			SwitchPmParameter string `json:"switch_pm_parameter"`
		}

		data := Payload{
			// fill struct
		}
		//res_ := []string{"test1", "test2"}
		fmt.Println(val.InlineQuery.ID)
		data.InlineQueryID = val.InlineQuery.ID
		//data.Results = "1111111111111111111111"

		payloadBytes, err := json.Marshal(data)
		if err != nil {
			// handle err
		}
		body := bytes.NewReader(payloadBytes)

		req, err := http.NewRequest("POST", telegramUrl+api+"/answerInlineQuery", body)
		if err != nil {
			// handle err
		}
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			// handle err
		}
		defer resp.Body.Close()

		bodyret, _ := ioutil.ReadAll(resp.Body)
		ret := new(InlineReturn)
		json.Unmarshal(bodyret, ret)

		fmt.Println(string(bodyret))
		//fmt.Println(ret)
	}
}

func main() {
	FlagDbg = true
	obj := new(Action)
	obj.ProxyUsage = false
	obj.ProxyUrl = ""
	obj.Bolt, _ = bolt.Open("telega.db", 0750, &bolt.Options{Timeout: 1 * time.Second})
	obj.getUpdates()
	obj.sendMessage()
	obj.Bolt.View(func(tx *bolt.Tx) error {
		b_ := tx.Bucket([]byte("Get"))
		b_.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})
	//time.Sleep(10 * time.Second)
	//obj.getUpdates()
	//time.Sleep(1 * time.Second)
	obj.answerCallbackQuery()
	defer obj.Bolt.Close()
}
