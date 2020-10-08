package main

import (
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"strconv"
	//tgbotapi "github.com/Syfaro/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
	"os"
	"time"
)

const (
	ClipsPerPage        = 2
	TagsPerPage 		= 2
	GroupsPerPage       = 2
	UsersPerPage        = 2
	GroupMembersPerPage = 2
	CommentsPerPage     = 3
)


func monitorUsers(ticker *time.Ticker, users Users) {
	for {
		select {
		case <-ticker.C:
			for key, user := range users {
				if time.Since(user.LastSeen).Hours() > 1 {
					users.Delete(key)
				}
			}
		}
	}
}

func ReadPsSqlDsn() string {
	host := flag.String("psSqlHost", os.Getenv("PsSQLHost"), "postgresql db host")
	portString := flag.String("psSqlPort", os.Getenv("PsSQLPort"), "postgresql db port")
	user := flag.String("psSqlUser", os.Getenv("PsSQLUser"), "postgresql db user")
	password := flag.String("psSqlPassword", os.Getenv("PsSQLPassword"), "postgresql db password")
	dbname := flag.String("psSqlDbname", os.Getenv("PsSQLDbname"), "postgresql db name")
	var port int
	withPassword := true

	flag.Parse()

	if len(*host) == 0 {
		log.Fatal("sql host not specified")
	}

	if len(*portString) == 0 {
		log.Fatal("sql port not specified")
	} else {
		var err error
		port, err = strconv.Atoi(*portString)
		if err != nil {
			log.Fatal("unable to parse sql port")
		}
	}

	if len(*user) == 0 {
		log.Fatal("sql user not specified")
	}

	if len(*password) == 0 {
		withPassword = false
		log.Println("sql password not specified, will try to connect without password")
	}

	if len(*dbname) == 0 {
		log.Fatal("sql db name not specified")
	}

	var dsn string
	if withPassword {
		dsn = fmt.Sprintf("host=%s port=%d user=%s "+
			"password=%s dbname=%s sslmode=disable",
			*host, port, *user, *password, *dbname)
	} else {
		dsn = fmt.Sprintf("host=%s port=%d user=%s "+
			"dbname=%s sslmode=disable",
			*host, port, *user, *dbname)
	}

	return dsn
}


func ConnectToSqlDb(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal("unable to use sql data source name", err)

	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("unable to connect to sql database: %v", err)
	}

	log.Println("successfully connected to sql database")
	return db
}

func CloseSqlDb(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("error during closing sql database connection: %v", err)
	} else {
		log.Println("sql database connection closed successfully")
	}
}


func NewUuidString() string {
	id := uuid.New()
	idBytes := id[:]
	return hex.EncodeToString(idBytes)
}

func main() {
	err := godotenv.Load("Database.env")
	if err != nil {
		log.Fatalf("Error loading Database.env file")
	}

	psSqlDsn := ReadPsSqlDsn()
	psSqlDb := ConnectToSqlDb(psSqlDsn)
	defer CloseSqlDb(psSqlDb)

	awsS3Region, awsS3Bucket := ReadAwsS3ConnectionData()
	awsS3 := ConnectToAwsS3(awsS3Region, awsS3Bucket)


	bot, err := tgbotapi.NewBotAPI(os.Getenv("Token"))
	if err != nil {
		log.Fatalf("unable to start bot: %v", err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	users := InitUsers()
	timeTicker := time.NewTicker(30 * time.Minute)

	go monitorUsers(timeTicker, users)
	for update := range updates {
		if update.CallbackQuery != nil {
			chatId := update.CallbackQuery.Message.Chat.ID
			requestTelegramId := update.CallbackQuery.From.ID

			userInfo := users.User(requestTelegramId)
			userInfo.LastSeen = time.Now()

			msg := tgbotapi.NewEditMessageText(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, "")
			msg.ParseMode = "markdown"

			buttonData := ParseButtonData(update.CallbackQuery.Data)
			userInfo.StackedData += buttonData.String()

			log.Printf("[%s u: %d c: %d] %s\n", update.CallbackQuery.From.UserName, requestTelegramId, chatId, update.CallbackQuery.Data)
			log.Println(buttonData)

			msg.Text, msg.ReplyMarkup = ProcessButtonPush(psSqlDb, awsS3, bot, &update, userInfo, users, buttonData)

			bot.Send(msg)
			//bot.AnswerCallbackQuery(tgbotapi.NewCallback(update.CallbackQuery.ID, "Done"))
		} else if update.Message != nil {
			chatId := update.Message.Chat.ID
			requestTelegramId := update.Message.From.ID

			log.Printf("[%s u: %d c: %d] %s\n", update.Message.From.UserName, requestTelegramId, chatId, update.Message.Text)

			userInfo := users.User(requestTelegramId)
			userInfo.LastSeen = time.Now()

			msg := tgbotapi.NewMessage(chatId, "")
			msg.ParseMode = "markdown"

			if update.Message.Video != nil {
				if users.User(requestTelegramId).UserState != SendClipVideoFile {
					msg.Text, msg.ReplyMarkup = "Бот не ожидает от вас текстовых сообщений в данный момент. " +
						"Пожалуйста, используйте кнопки или отправьте команду \"/menu\".", makeEmptyMenu()
					bot.Send(msg)
					continue
				}

				msg.Text, msg.ReplyMarkup = HandleUploadClip(bot, psSqlDb, awsS3, chatId, userInfo,
					requestTelegramId, userInfo.TempClip, userInfo.TempTags, update.Message.Video)
				userInfo.TempClip = nil
				userInfo.TempTags = nil

				bot.Send(msg)
				continue
			} else {
				if update.Message.IsCommand() {
					switch update.Message.Command() {
					case "start":
						msg.Text, msg.ReplyMarkup = HandleStartCommand(psSqlDb, userInfo, requestTelegramId)
					case "help", "menu":
						msg.Text, msg.ReplyMarkup = HandleMainMenu(psSqlDb, userInfo, requestTelegramId, users)
					default:
						msg.Text = "Данная команда не предусмотрена"
					}
					bot.Send(msg)
					continue
				}

				text := update.Message.Text
				switch userInfo.UserState {
				case AddClipName:
					msg.Text, msg.ReplyMarkup = HandleAddClipNameMessage(psSqlDb, awsS3, bot, &update, userInfo, users, text)
				case AddClipDescrpition:
					msg.Text, msg.ReplyMarkup = HandleAddClipDescriptionMessage(psSqlDb, awsS3, bot, &update, userInfo, users, text)
				case AddClipTags:
					msg.Text, msg.ReplyMarkup = HandleAddClipTagsMessage(psSqlDb, awsS3, bot, &update, userInfo, users, text)
				case AddGroupName:
					msg.Text, msg.ReplyMarkup = HandleAddGroupNameMessage(psSqlDb, awsS3, bot, &update, userInfo, users, text)
				case AddGroupDescription:
					msg.Text, msg.ReplyMarkup = HandleAddGroupDescriptionMessage(psSqlDb, awsS3, bot, &update, userInfo, users, text)
				case AddComment:
					msg.Text, msg.ReplyMarkup = HandleAddCommentMessage(psSqlDb, awsS3, bot, &update, requestTelegramId, userInfo, users, text)
				case EditUserDescription:
					msg.Text, msg.ReplyMarkup = HandleEditUserDescriptionMessage(psSqlDb, awsS3, bot, &update, userInfo, users, requestTelegramId, text)
				case EditClipDescription:
					msg.Text, msg.ReplyMarkup = HandleEditClipDescriptionMessage(psSqlDb, awsS3, bot, &update, userInfo, users, requestTelegramId, userInfo.TempClip.Id, text)
				case EditGroupDescription:
					msg.Text, msg.ReplyMarkup = HandleEditGroupDescriptionMessage(psSqlDb, awsS3, bot, &update, userInfo, users, requestTelegramId, userInfo.TempGroup.Id, text)
				case SearchClipByUuid:
					msg.Text, msg.ReplyMarkup = HandleSearchClipByUuidMessage(bot, psSqlDb, awsS3, chatId, userInfo, requestTelegramId, text)
				case SearchUserByUserName:
					msg.Text, msg.ReplyMarkup = HandleSearchUserByUserNameMessage(psSqlDb, awsS3, bot, &update, requestTelegramId, text, userInfo, users)
				case AddUserName:
					msg.Text, msg.ReplyMarkup = HandleAddUserNameMessage(psSqlDb, requestTelegramId, text, userInfo)
				default:
					msg.Text, msg.ReplyMarkup = "Бот не ожидает от вас текстовых сообщений в данный момент. " +
						"Пожалуйста, используйте кнопки или отправьте команду \"/menu\".", makeEmptyMenu()
				}
				bot.Send(msg)
				continue
			}
		}
	}
}



func ProcessButtonPush(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI,
	update *tgbotapi.Update, userInfo *UserInfo, users Users, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	buttonNameParameter, _ := buttonDataHead.ReadParameter(ButtonName)
	var chatId int64
	if update.CallbackQuery != nil {
		chatId = update.CallbackQuery.Message.Chat.ID
	} else {
		chatId = update.Message.Chat.ID
	}

	requestTelegramId := GetRequestTelegramIdFromUpdate(update)

	text := ""
	var menu *tgbotapi.InlineKeyboardMarkup
	switch buttonNameParameter.Value {
	case DeleteClipButtonName:
		text, menu = HandleDeleteClip(bot, sqlDb, awsS3, chatId, userInfo, requestTelegramId, buttonDataHead)
	case MainMenuButtonName:
		text, menu = HandleMainMenu(sqlDb, userInfo, requestTelegramId, users)
	case MyClipsButtonName:
		text, menu = HandleMyClips(buttonDataHead)
	case MyGroupsButtonName:
		text, menu = HandleMyGroups(buttonDataHead)
	case MySubscriptionsButtonName:
		text, menu = HandleMySubscriptions(buttonDataHead)
	case UserButtonName:
		text, menu = HandleUser(sqlDb, requestTelegramId, userInfo, buttonDataHead)
	case SearchButtonName:
		text, menu = HandleSearch()
	case BackButtonName:
		text, menu = HandleBack(sqlDb, awsS3, bot, update, userInfo, users, buttonDataHead)
	case SearchUserButtonName:
		text, menu = HandleSearchUser(userInfo)
	case SearchClipButtonName:
		text, menu = HandleSearchClip(userInfo)
	case AddClipButtonName:
		text, menu = HandleAddClip(sqlDb, userInfo, requestTelegramId, buttonDataHead, false)
	case ClipsButtonName:
		text, menu = HandleClips(sqlDb, awsS3, bot, update, userInfo, users, requestTelegramId, buttonDataHead)
	case ClipButtonName:
		text, menu = HandleClip(bot, sqlDb, awsS3, chatId, requestTelegramId, userInfo, buttonDataHead)
	case AddClipNameButtonName:
		text, menu = HandleAddClipName(userInfo)
	case AddClipDescriptionButtonName:
		text, menu = HandleAddClipDescription(userInfo)
	case AddClipTagsButtonName:
		text, menu = HandleAddClipTags(userInfo)
	case SendClipVideoFileButtonName:
		text, menu = HandleSendClipVideoFile(userInfo)
	case EditClipDescriptionButtonName:
		text, menu = HandleEditClipDescription(userInfo)
	case EditClipProfileButtonName:
		text, menu = HandleEditClipProfile(userInfo, buttonDataHead, false)
	case GroupsButtonName:
		text, menu = HandleGroups(sqlDb, awsS3, bot, update, requestTelegramId, userInfo, users, buttonDataHead)
	case GroupButtonName:
		text, menu = HandleGroup(sqlDb, requestTelegramId, userInfo, buttonDataHead)
	case GroupMembersButtonName:
		text, menu = HandleGroupMembers(sqlDb, awsS3, bot, update, requestTelegramId, userInfo, users, buttonDataHead)
	case JoinGroupButtonName:
		text, menu = HandleJoin(sqlDb, awsS3, bot, update, userInfo, users, buttonDataHead)
	case LeaveGroupButtonName:
		text, menu = HandleDisjoin(sqlDb, awsS3, bot, update, userInfo, users, buttonDataHead)
	case CreateGroupButtonName:
		text, menu = HandleCreateGroup(sqlDb, userInfo, requestTelegramId, false)
	case AddGroupNameButtonName:
		text, menu = HandleAddGroupName(userInfo)
	case AddGroupDescriptionButtonName:
		text, menu = HandleAddGroupDescription(userInfo)
	case FinishGroupCreationButtonName:
		text, menu = HandleFinishGroupCreation(sqlDb, awsS3, bot, update, userInfo, users, requestTelegramId)
	case EditGroupProfileButtonName:
		text, menu = HandleEditGroupProfile(userInfo, buttonDataHead, false)
	case EditGroupDescriptionButtonName:
		text, menu = HandleEditGroupDescription(userInfo)
	case DeleteGroupButtonName:
		text, menu = HandleDeleteGroup(sqlDb, userInfo, requestTelegramId, buttonDataHead)
	case SubscriptionsButtonName:
		text, menu = HandleSubscriptions(sqlDb, awsS3, bot, update, requestTelegramId, userInfo, users, buttonDataHead)
	case SubscribeButtonName:
		text, menu = HandleSubscribe(sqlDb, awsS3, bot, update, userInfo, users, buttonDataHead)
	case UnsubscribeButtonName:
		text, menu = HandleUnsubscribe(sqlDb, awsS3, bot, update, userInfo, users, buttonDataHead)
	case EditUserProfileButtonName:
		text, menu = HandleEditUserProfile(buttonDataHead)
	case EditUserDescriptionButtonName:
		text, menu = HandleEditUserDescription(userInfo)
	case CommentsButtonName:
		text, menu = HandleComments(sqlDb, requestTelegramId, userInfo, buttonDataHead)
	case LeaveCommentButtonName:
		text, menu = HandleAddComment(userInfo, buttonDataHead)
	case TagButtonName:
		text, menu = HandleTag(sqlDb, awsS3, bot, update, userInfo, users, requestTelegramId, buttonDataHead)
	case TagsButtonName:
		text, menu = HandleTags(sqlDb, awsS3, bot, update, requestTelegramId, userInfo, users, buttonDataHead)
	default:
		text, menu = fmt.Sprintf("Неизвестная кнопка с именем %s", buttonNameParameter.Value), nil
	}

	return text, menu
}
