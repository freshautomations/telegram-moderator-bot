package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/freshautomations/telegram-moderator-bot/context"
	"github.com/freshautomations/telegram-moderator-bot/defaults"
	"log"
	"net/http"
)

type User struct {
	Id           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

func (u User) String() string {
	name := u.FirstName
	if u.LastName != "" {
		name = fmt.Sprintf("%s %s", name, u.LastName)
	}
	if u.Username != "" {
		return fmt.Sprintf("%s (%s)", name, u.Username)
	}
	return name
}

type Chat struct {
	Id                          int64      `json:"id"`
	Type                        string     `json:"type"`
	Title                       string     `json:"title"`
	Username                    string     `json:"username"`
	FirstName                   string     `json:"first_name"`
	LastName                    string     `json:"last_name"`
	AllMembersAreAdministrators bool       `json:"all_members_are_administrators"`
	Photo                       *ChatPhoto `json:"photo"`
	Description                 string     `json:"description"`
	InviteLink                  string     `json:"invite_link"`
	PinnedMessage               *Message   `json:"pinned_message"`
	StickerSetName              string     `json:"sticker_set_name"`
	CanSetStickerSet            bool       `json:"can_set_sticker_set"`
}

type ChatPhoto struct {
	SmallFileId string `json:"small_file_id"`
	BigFileId   string `json:"big_file_id"`
}

type MessageEntity struct {
	Type   string `json:"type"`
	Offset int    `json:"offset"`
	Length int    `json:"length"`
	Url    string `json:"url"`
	User   *User  `json:"user"`
}

type Audio struct {
	FileId    string     `json:"file_id"`
	Duration  int        `json:"duration"`
	Performer string     `json:"performer"`
	Title     string     `json:"title"`
	MimeType  string     `json:"mime_type"`
	FileSize  int        `json:"file_size"`
	Thumb     *PhotoSize `json:"thumb"`
}

type PhotoSize struct {
	FileId   string `json:"file_id"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	FileSize int    `json:"file_size"`
}

type Document struct {
	FileId   string     `json:"file_id"`
	Thumb    *PhotoSize `json:"thumb"`
	FileName string     `json:"file_name"`
	MimeType string     `json:"mime_type"`
	FileSize int        `json:"file_size"`
}

type Animation struct {
	FileId   string     `json:"file_id"`
	Width    int        `json:"width"`
	Height   int        `json:"height"`
	Duration int        `json:"duration"`
	Thumb    *PhotoSize `json:"thumb"`
	FileName string     `json:"file_name"`
	MimeType string     `json:"mime_type"`
	FileSize int        `json:"file_size"`
}

type Game struct {
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Photo        []*PhotoSize     `json:"photo"`
	Text         string           `json:"Text"`
	TextEntities []*MessageEntity `json:"text_entities"`
	Animation    *Animation       `json:"animation"`
}

type Sticker struct {
	FileId       string        `json:"file_id"`
	Width        int           `json:"width"`
	Height       int           `json:"height"`
	Thumb        *PhotoSize    `json:"thumb"`
	Emoji        string        `json:"emoji"`
	SetName      string        `json:"set_name"`
	MaskPosition *MaskPosition `json:"mask_position"`
	FileSize     int           `json:"file_size"`
}

type MaskPosition struct {
	Point  string  `json:"point"`
	XShift float32 `json:"x_shift"`
	YShift float32 `json:"y_shift"`
	Scale  float32 `json:"scale"`
}

type Video struct {
	FileId   string     `json:"file_id"`
	Width    int        `json:"width"`
	Height   int        `json:"height"`
	Duration int        `json:"duration"`
	Thumb    *PhotoSize `json:"thumb"`
	MimeType string     `json:"mime_type"`
	FileSize int        `json:"file_size"`
}

type Voice struct {
	FileId   string `json:"file_id"`
	Duration int    `json:"duration"`
	MimeType string `json:"mime_type"`
	FileSize int    `json:"file_size"`
}

type VideoNote struct {
	FileId   string     `json:"file_id"`
	Length   int        `json:"length"`
	Duration int        `json:"duration"`
	Thumb    *PhotoSize `json:"thumb"`
	FileSize int        `json:"file_size"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UserId      int    `json:"user_id"`
	Vcard       string `json:"vcard"`
}

type Location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
}

type Venue struct {
	Location       *Location `json:"location"`
	Title          string    `json:"title"`
	Address        string    `json:"address"`
	FoursquareId   string    `json:"foursquare_id"`
	FoursquareType string    `json:"foursquare_type"`
}

type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"`
	TotalAmount    int    `json:"total_amount"`
}

type SuccessfulPayment struct {
	Currency                string     `json:"currency"`
	TotalAmount             int        `json:"total_amount"`
	InvoicePayload          string     `json:"invoice_payload"`
	ShippingOptionId        string     `json:"shipping_option_id"`
	OrderInfo               *OrderInfo `json:"order_info"`
	TelegramPaymentChargeId string     `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeId string     `json:"provider_payment_charge_id"`
}

type OrderInfo struct {
	Name            string           `json:"name"`
	PhoneNumber     string           `json:"phone_number"`
	Email           string           `json:"email"`
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

type PassportData struct {
	data        []*EncryptedPassportElement `json:"data"`
	credentials *EncryptedCredentials       `json:"credentials"`
}

type EncryptedPassportElement struct {
	Type        string          `json:"type"`
	Data        string          `json:"data"`
	PhoneNumber string          `json:"phone_number"`
	Email       string          `json:"email"`
	Files       []*PassportFile `json:"files"`
	FrontSide   PassportFile    `json:"front_side"`
	ReverseSide PassportFile    `json:"reverse_side"`
	Selfie      PassportFile    `json:"selfie"`
	Translation []*PassportFile `json:"translation"`
	Hash        string          `json:"hash"`
}

type PassportFile struct {
	FileId   string `json:"file_id"`
	FileSize int    `json:"file_size"`
	FileDate int    `json:"file_date"`
}

type EncryptedCredentials struct {
	Data   string `json:"data"`
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
}

type Message struct {
	MessageId             int64              `json:"message_id"`
	From                  *User              `json:"from"`
	Date                  int                `json:"date"`
	Chat                  *Chat              `json:"chat"`
	ForwardFrom           *User              `json:"forward_from"`
	ForwardFromChat       *Chat              `json:"forward_from_chat"`
	ForwardFromMessageId  int                `json:"forward_from_message_id"`
	ForwardSignature      string             `json:"forward_signature"`
	ForwardDate           int                `json:"forward_date"`
	ReplyToMessage        *Message           `json:"reply_to_message"`
	EditDate              int                `json:"edit_date"`
	MediaGroupId          string             `json:"media_group_id"`
	AuthorSignature       string             `json:"author_signature"`
	Text                  string             `json:"text"`
	Entities              []*MessageEntity   `json:"entities"`
	CaptionEntities       []*MessageEntity   `json:"caption_entities"`
	Audio                 *Audio             `json:"audio"`
	Document              *Document          `json:"document"`
	Animation             *Animation         `json:"animation"`
	Game                  *Game              `json:"game"`
	Photo                 []*PhotoSize       `json:"photo"`
	Sticker               *Sticker           `json:"sticker"`
	Video                 *Video             `json:"video"`
	Voice                 *Voice             `json:"voice"`
	VideoNote             *VideoNote         `json:"video_note"`
	Caption               string             `json:"caption"`
	Contact               *Contact           `json:"contact"`
	Location              *Location          `json:"location"`
	Venue                 *Venue             `json:"venue"`
	NewChatMembers        []*User            `json:"new_chat_members"`
	LeftChatMember        *User              `json:"left_chat_member"`
	NewChatTitle          string             `json:"new_chat_title"`
	NewChatPhoto          []*PhotoSize       `json:"new_chat_photo"`
	DeleteChatPhoto       bool               `json:"delete_chat_photo"`
	GroupChatCreated      bool               `json:"group_chat_created"`
	SupergroupChatCreated bool               `json:"supergroup_chat_created"`
	ChannelChatCreated    bool               `json:"channel_chat_created"`
	MigrateToChatId       int64              `json:"migrate_to_chat_id"`
	MigrateFromChatId     int64              `json:"migrate_from_chat_id"`
	PinnedMessage         *Message           `json:"pinned_message"`
	Invoice               *Invoice           `json:"invoice"`
	SuccessfulPayment     *SuccessfulPayment `json:"successful_payment"`
	ConnectedWebsite      string             `json:"connected_website"`
	PassportData          *PassportData      `json:"passport_data"`
}

type Update struct {
	UpdateId           int                 `json:"update_id"`
	Message            *Message            `json:"message"`
	EditedMessage      *Message            `json:"edited_message"`
	ChannelPost        *Message            `json:"channel_post"`
	EditedChannelPost  *Message            `json:"edited_channel_post"`
	InlineQuery        *InlineQuery        `json:"inline_query"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *CallbackQuery      `json:"callback_query"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query"`
}

type InlineQuery struct {
	Id       string    `json:"id"`
	From     *User     `json:"from"`
	Location *Location `json:"location"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
}

type ChosenInlineResult struct {
	ResultId        string    `json:"result_id"`
	From            *User     `json:"from"`
	Location        *Location `json:"location"`
	InlineMessageId string    `json:"inline_message_id"`
	Query           string    `json:"query"`
}

type CallbackQuery struct {
	Id              string   `json:"id"`
	From            *User    `json:"from"`
	Message         *Message `json:"message"`
	InlineMessageId string   `json:"inline_message_id"`
	ChatInstance    string   `json:"chat_instance"`
	Data            string   `json:"data"`
	GameShortName   string   `json:"game_short_name"`
}

type ShippingQuery struct {
	Id              string           `json:"id"`
	From            *User            `json:"from"`
	InvoicePayload  string           `json:"invoice_payload"`
	ShippingAddress *ShippingAddress `json:"shipping_address"`
}

type PreCheckoutQuery struct {
	Id               string     `json:"id"`
	From             *User      `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionId string     `json:"shipping_option_id"`
	OrderInfo        *OrderInfo `json:"order_info"`
}

type ChatMember struct {
	User                  *User  `json:"user"`
	Status                string `json:"status"` //Can be creator, administrator, member, restricted, left or kicked
	UntilDate             int    `json:"until_date,omitempty"`
	CanBeEdited           bool   `json:"can_be_edited,omitempty"`
	CanChangeInfo         bool   `json:"can_change_info,omitempty"`
	CanPostMessages       bool   `json:"can_post_messages,omitempty"`
	CanEditMessages       bool   `json:"can_edit_messages,omitempty"`
	CanDeleteMessages     bool   `json:"can_delete_messages,omitempty"`
	CanInviteUsers        bool   `json:"can_invite_users,omitempty"`
	CanRestrictMembers    bool   `json:"can_restrict_members,omitempty"`
	CanPinMessages        bool   `json:"can_pin_messages,omitempty"`
	CanPromoteMembers     bool   `json:"can_promote_members,omitempty"`
	CanSendMessages       bool   `json:"can_send_messages,omitempty"`
	CanSendMediaMessages  bool   `json:"can_send_media_messages,omitempty"`
	CanSendOtherMessages  bool   `json:"can_send_other_messages,omitempty"`
	CanAddWebPagePreviews bool   `json:"can_add_web_page_previews,omitempty"`
}

type SendMessageRequest struct {
	ChatId                int64  `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool   `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool   `json:"disable_notification,omitempty"`
	ReplyToMessageId      int64  `json:"reply_to_message_id,omitempty"`
	//Todo: implement reply_markup
}

type SendMessageResponse struct {
	Ok          bool    `json:"ok"`
	Result      Message `json:"result"`
	ErrorCode   int     `json:"error_code,omitempty"`
	Description string  `json:"description,omitempty"`
}

type GetChatAdministratorsRequest struct {
	ChatId int64 `json:"chat_id"`
}

type GetChatAdministratorsResponse struct {
	Ok          bool          `json:"ok"`
	Result      []*ChatMember `json:"result"`
	ErrorCode   int           `json:"error_code,omitempty"`
	Description string        `json:"description,omitempty"`
}

type GetChatMemberRequest struct {
	ChatId int64 `json:"chat_id"`
	UserId int   `json:"user_id"`
}

type GetChatMemberResponse struct {
	Ok          bool        `json:"ok"`
	Result      *ChatMember `json:"result"`
	ErrorCode   int         `json:"error_code,omitempty"`
	Description string      `json:"description,omitempty"`
}

type PromoteChatMemberRequest struct {
	ChatId             int64 `json:"chat_id"`
	UserId             int   `json:"user_id"`
	CanChangeInfo      bool  `json:"can_change_info,omitempty"`
	CanPostMessages    bool  `json:"can_post_messages,omitempty"`
	CanEditMessages    bool  `json:"can_edit_messages,omitempty"`
	CanDeleteMessages  bool  `json:"can_delete_messages,omitempty"`
	CanInviteUsers     bool  `json:"can_invite_users,omitempty"`
	CanRestrictMembers bool  `json:"can_restrict_members,omitempty"`
	CanPinMessages     bool  `json:"can_pin_messages,omitempty"`
	CanPromoteMembers  bool  `json:"can_promote_members,omitempty"`
}

type PromoteChatMemberResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

type KickChatMemberRequest struct {
	ChatId    int64 `json:"chat_id"`
	UserId    int   `json:"user_id"`
	UntilDate int64 `json:"until_date,omitempty"`
}

type KickChatMemberResponse struct {
	Ok          bool   `json:"ok"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

// Sends a message to a supergroup that responds to a message.
func SendMessage(ctx *context.Context, ChatId int64, ReplyToMessageId int64, Text string) error {
	jsonValue, _ := json.Marshal(SendMessageRequest{
		ChatId:              ChatId,
		Text:                Text,
		ReplyToMessageId:    ReplyToMessageId,
		DisableNotification: true,
	})

	m, err := http.Post(defaults.TelegramAPIBase+ctx.Cfg.TelegramToken+"/sendMessage", defaults.ContentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("[error] Telegram API response: %v", err)
		return err
	}

	incoming := &SendMessageResponse{}
	err = json.NewDecoder(m.Body).Decode(incoming)
	if err != nil {
		log.Printf("[error] SendMessage decoder: %v", err)
		return err
	}

	if incoming.Ok {
		if defaults.Debug {
			log.Printf("[debug] SendMessage: %s", incoming.Result.Text)
		}
	} else {
		log.Printf("[error] SendMessage %d %s.", incoming.ErrorCode, incoming.Description)
		return errors.New(incoming.Description)
	}

	return nil
}

// Check a user's privileges. Returns IsAdministrator, IsModerator, error.
func GetPrivileges(ctx *context.Context, ChatId int64, UserId int) (bool, bool, error) {
	jsonValue, _ := json.Marshal(GetChatAdministratorsRequest{
		ChatId: ChatId,
	})

	m, err := http.Post(defaults.TelegramAPIBase+ctx.Cfg.TelegramToken+"/getChatAdministrators", defaults.ContentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("[error] Telegram API response: %v", err)
		return false, false, err
	}

	incoming := &GetChatAdministratorsResponse{}
	err = json.NewDecoder(m.Body).Decode(incoming)
	if err != nil {
		log.Printf("[error] GetPrivileges decoder: %v", err)
		return false, false, err
	}

	for _, member := range incoming.Result {
		if member.User.Id == UserId {
			return member.CanPromoteMembers || member.Status == "creator", true, nil
		}
	}

	return false, false, nil
}

// Retrieves the user details of a member of a supergroup based on user ID.
func GetChatMember(ctx *context.Context, ChatId int64, UserId int) (*ChatMember, error) {
	jsonValue, _ := json.Marshal(GetChatMemberRequest{
		ChatId: ChatId,
		UserId: UserId,
	})

	m, err := http.Post(defaults.TelegramAPIBase+ctx.Cfg.TelegramToken+"/getChatMember", defaults.ContentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("[error] Telegram API response: %v", err)
		return nil, err
	}

	incoming := &GetChatMemberResponse{}
	err = json.NewDecoder(m.Body).Decode(incoming)
	if err != nil {
		log.Printf("[error] GetChatMember decoder: %v", err)
		return nil, err
	}

	if incoming.Ok {
		return incoming.Result, nil
	}

	return nil, errors.New(fmt.Sprintf("(%d) %s", incoming.ErrorCode, incoming.Description))
}

// Add moderators to a supergroup.
func AddModerator(ctx *context.Context, ChatId int64, Users []*User) (list []string, errors []string) {
	for _, user := range Users {
		jsonValue, _ := json.Marshal(PromoteChatMemberRequest{
			ChatId:             ChatId,
			UserId:             user.Id,
			CanChangeInfo:      false,
			CanPostMessages:    false, //Channels only
			CanEditMessages:    false, //Channels only
			CanDeleteMessages:  true,
			CanInviteUsers:     false,
			CanRestrictMembers: false,
			CanPinMessages:     true,
			CanPromoteMembers:  false,
		})

		m, err := http.Post(defaults.TelegramAPIBase+ctx.Cfg.TelegramToken+"/promoteChatMember", defaults.ContentType, bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Printf("[error] Telegram API response: %+v, %+v", user, err)
			continue
		}

		incoming := &PromoteChatMemberResponse{}
		err = json.NewDecoder(m.Body).Decode(incoming)
		if err != nil {
			log.Printf("[error] AddModerator decoder: %+v, %+v", user, err)
			continue
		}

		if incoming.Ok {
			list = append(list, user.String())
		} else {
			log.Printf("[error] AddModerator response: %d, %s, %+v", incoming.ErrorCode, incoming.Description, user)
			errors = append(errors, fmt.Sprintf("%s (%s %s): %d: %s", user.Username, user.FirstName, user.LastName, incoming.ErrorCode, incoming.Description))
		}
	}

	return
}

// Remove moderators from a supergroup.
func RemoveModerator(ctx *context.Context, ChatId int64, Users []*User) (list []string, errors []string) {
	for _, user := range Users {
		jsonValue, _ := json.Marshal(PromoteChatMemberRequest{
			ChatId:             ChatId,
			UserId:             user.Id,
			CanChangeInfo:      false,
			CanPostMessages:    false, //Channels only
			CanEditMessages:    false, //Channels only
			CanDeleteMessages:  false,
			CanInviteUsers:     false,
			CanRestrictMembers: false,
			CanPinMessages:     false,
			CanPromoteMembers:  false,
		})

		m, err := http.Post(defaults.TelegramAPIBase+ctx.Cfg.TelegramToken+"/promoteChatMember", defaults.ContentType, bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Printf("[error] Telegram API response: %+v, %+v", user, err)
			continue
		}

		incoming := &PromoteChatMemberResponse{}
		err = json.NewDecoder(m.Body).Decode(incoming)
		if err != nil {
			log.Printf("[error] RemoveModerator decoder: %+v, %+v", user, err)
			continue
		}

		if incoming.Ok {
			list = append(list, user.String())
		} else {
			log.Printf("[error] RemoveModerator response: %d, %s, %+v", incoming.ErrorCode, incoming.Description, user)
			errors = append(errors, fmt.Sprintf("%s (%s %s): %d: %s", user.Username, user.FirstName, user.LastName, incoming.ErrorCode, incoming.Description))
		}
	}

	return
}

// Ban regular members of a supergroup.
func BanMember(ctx *context.Context, ChatId int64, Users []*User) (result []string) {
	for _, user := range Users {
		jsonValue, _ := json.Marshal(KickChatMemberRequest{
			ChatId: ChatId,
			UserId: user.Id,
		})

		m, err := http.Post(defaults.TelegramAPIBase+ctx.Cfg.TelegramToken+"/kickChatMember", defaults.ContentType, bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Printf("[error] Telegram API response: %+v, %+v", user, err)
			continue
		}

		incoming := &KickChatMemberResponse{}
		err = json.NewDecoder(m.Body).Decode(incoming)
		if err != nil {
			log.Printf("[error] BanMember decoder: %+v, %+v", user, err)
			continue
		}

		if incoming.Ok {
			result = append(result, user.String())
		} else {
			log.Printf("[error] BanMember response: %d, %s, %+v", incoming.ErrorCode, incoming.Description, user)
		}
	}

	return
}

// Unban regular members from a supergroup..
func UnbanMember(ctx *context.Context, ChatId int64, Users []*User) (result []string) {
	for _, user := range Users {
		jsonValue, _ := json.Marshal(KickChatMemberRequest{
			ChatId: ChatId,
			UserId: user.Id,
		})

		m, err := http.Post(defaults.TelegramAPIBase+ctx.Cfg.TelegramToken+"/unbanChatMember", defaults.ContentType, bytes.NewBuffer(jsonValue))
		if err != nil {
			log.Printf("[error] Telegram API response: %+v, %+v", user, err)
			continue
		}

		incoming := &KickChatMemberResponse{}
		err = json.NewDecoder(m.Body).Decode(incoming)
		if err != nil {
			log.Printf("[error] UnbanMember decoder: %+v, %+v", user, err)
			continue
		}

		if incoming.Ok {
			result = append(result, user.String())
		} else {
			log.Printf("[error] UnbanMember response: %d, %s, %+v", incoming.ErrorCode, incoming.Description, user)
		}
	}

	return
}

// List moderators in a supergroup.
func ListModerators(ctx *context.Context, ChatId int64) (result []string) {
	jsonValue, _ := json.Marshal(GetChatAdministratorsRequest{
		ChatId: ChatId,
	})

	m, err := http.Post(defaults.TelegramAPIBase+ctx.Cfg.TelegramToken+"/getChatAdministrators", defaults.ContentType, bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("[error] Telegram API response: %v", err)
		return
	}

	incoming := &GetChatAdministratorsResponse{}
	err = json.NewDecoder(m.Body).Decode(incoming)
	if err != nil {
		log.Printf("[error] ListModerators decoder: %v", err)
		return
	}

	for _, member := range incoming.Result {
		if member.CanPromoteMembers || member.Status != "administrator" {
			continue
		}
		result = append(result, member.User.String())
	}

	return
}
