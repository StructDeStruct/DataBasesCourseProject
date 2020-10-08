package main

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lib/pq"
	"log"
	"strconv"
	"time"
)

const (
	DuplicatonInsertionError = "37846"
)


func makeDbErrorResponseData(err error, messageText string, sqlDb *sql.DB, requestTelegramId int, userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	log.Println("sql data base error:", err)
	user, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return "Бот не отвечает", makeEmptyMenu()
	}
	return messageText, FormMainMenu(user.Id, userInfo)
}

func makeEmptyMenu() *tgbotapi.InlineKeyboardMarkup {
	//menu := tgbotapi.NewInlineKeyboardMarkup()
	return &tgbotapi.InlineKeyboardMarkup{[][]tgbotapi.InlineKeyboardButton{}}
}


func HandleClips(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update, userInfo *UserInfo, users Users,
	requestTelegramId int, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	ownerIdParameter, _ := buttonDataHead.ReadParameter(OwnerIdParameterName)
	ownerId, _ := strconv.Atoi(ownerIdParameter.Value)

	firstClipOnPageNumberParameter, _ := buttonDataHead.ReadParameter(FirstClipOnPageNumberParameterName)
	firstClipOnPageNumber, _ := strconv.Atoi(firstClipOnPageNumberParameter.Value)

	clips, err := GetClipsPageByOwnerId(sqlDb, ownerId, firstClipOnPageNumber)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	ownerName, isGroup, err := GetOwnerName(sqlDb, ownerId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	if len(clips) == 0 {
		if isGroup {
			t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
			return fmt.Sprintf("В группе %s нет опубликованных клипов", ownerName) + "\n\n" + t1, menu
		} else {
			t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
			return fmt.Sprintf("У пользователя %s нет опубликованных клипов", ownerName) + "\n\n" + t1, menu
		}
	}

	menu := FormClipsPageMenu(clips, buttonDataHead, firstClipOnPageNumber)
	text := FormShortClipsInfo(firstClipOnPageNumber, clips, ownerName)
	return text, menu
}

func HandleStartCommand(sqlDb *sql.DB, userInfo *UserInfo, userId int) (string, *tgbotapi.InlineKeyboardMarkup) {
	user, err := GetUserDataByTelegramId(sqlDb, userId)
	if err != nil {
		log.Println(err)
		userInfo.UserState = AddUserName
		return "Кажется, Вы здесь в первый раз, выберите уникальное имя пользователя", makeEmptyMenu()
	}

	return FormWelcomeMessage(), FormMainMenu(user.Id, userInfo)
}

func HandleAddUserNameMessage(sqlDb *sql.DB, requestTelegramId int, userName string, userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	user, err := GetUserByUserName(sqlDb, userName)
	if err == nil {
		return "Пользователь с таким именем уже существует. Пожалуйста, выберите другое имя.", makeEmptyMenu()
	}

	user = &User{
		TelegramId: requestTelegramId,
		Name:       userName,
	}

	ownerId, err := InsertOwner(sqlDb, false)
	if err != nil {
		//log.Println("\t\t\t didn't insert owner")
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	err = InsertUser(sqlDb, user, ownerId)
	if err != nil {
		//log.Println("\t\t\t didn't insert user")
		_ = DeleteOwner(sqlDb, ownerId)
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	return FormUserRegistrationMessage(user), FormMainMenu(user.Id, userInfo)
}


func HandleGroups(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update, requestTelegramId int,
	userInfo *UserInfo, users Users, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	userIdParameter, _ := buttonDataHead.ReadParameter(UserIdParameterName)
	userId, _ := strconv.Atoi(userIdParameter.Value)

	firstGroupOnPageNumberParameter, _ := buttonDataHead.ReadParameter(FirstGroupOnPageNumberParameterName)
	firstGroupOnPageNumber, _ := strconv.Atoi(firstGroupOnPageNumberParameter.Value)

	groups, creatorNames, err := GetGroupsPage(sqlDb, userId, firstGroupOnPageNumber)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	if len(groups) == 0 {
		t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
		return "Пользователь не числится ни в одной группе" + "\n\n" + t1, menu
	}

	menu := FormGroupsPageMenu(groups, buttonDataHead, firstGroupOnPageNumber)
	text := FormShortGroupsInfo(firstGroupOnPageNumber, groups, creatorNames)
	return text, menu
}

func HandleSubscriptions(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	requestTelegramId int, userInfo *UserInfo, usersInfos Users, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	userIdParameter, _ := buttonDataHead.ReadParameter(UserIdParameterName)
	userId, _ := strconv.Atoi(userIdParameter.Value)

	firstSubscriptonOnPageNumberParameter, _ := buttonDataHead.ReadParameter(FirstUserOnPageNumberParameterName)
	firstSubscriptionOnPageNumber, _ := strconv.Atoi(firstSubscriptonOnPageNumberParameter.Value)

	users, err := GetSubscriptionsPage(sqlDb, userId, firstSubscriptionOnPageNumber)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	if len(users) == 0 {
		t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, usersInfos)
		return "Пользователь ни на кого не подписан" + "\n\n" + t1, menu
	}

	menu := FormSubscriptionsPageMenu(users, buttonDataHead, firstSubscriptionOnPageNumber)
	text := FormShortUsersInfo(firstSubscriptionOnPageNumber, users)

	return text, menu
}

func HandleGroupMembers(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	requestTelegramId int, userInfo *UserInfo, usersInfos Users, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	groupIdParameter, _ := buttonDataHead.ReadParameter(GroupIdParameterName)
	GroupId, _ := strconv.Atoi(groupIdParameter.Value)

	firstGroupMemberOnPageNumberParameter, _ := buttonDataHead.ReadParameter(FirstUserOnPageNumberParameterName)
	firstGroupMemberOnPageNumber, _ := strconv.Atoi(firstGroupMemberOnPageNumberParameter.Value)

	users, err := GetGroupMembersPage(sqlDb, GroupId, firstGroupMemberOnPageNumber)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	if len(users) == 0 {
		t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, usersInfos)
		return "В группе нет ни одного участника" + "\n\n" + t1, menu
	}

	menu := FormGroupMembersPageMenu(users, buttonDataHead, firstGroupMemberOnPageNumber)
	text := FormShortUsersInfo(firstGroupMemberOnPageNumber, users)

	return text, menu
}

func HandleComments(sqlDb *sql.DB, requestTelegramId int, userInfo *UserInfo, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	clipIdParameter, _ := buttonDataHead.ReadParameter(ClipIdParameterName)
	clipId, _ := strconv.Atoi(clipIdParameter.Value)

	parentCommentIdParameter, err := buttonDataHead.ReadParameter(ParentCommentIdParameterName)
	var hasParrentComment bool
	var parentCommentId int
	if err == nil {
		hasParrentComment = true
		parentCommentId, _ = strconv.Atoi(parentCommentIdParameter.Value)
	}

	firstCommentOnPageNumberParameter, _ := buttonDataHead.ReadParameter(FirstCommentOnPageNumberParameterName)
	firstCommentOnPageNumber, _ := strconv.Atoi(firstCommentOnPageNumberParameter.Value)

	parentComment, parentCommentUser, comments, users, err := GetCommentsPage(sqlDb,
		clipId, parentCommentId, hasParrentComment, firstCommentOnPageNumber)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	menu := FormCommentsPageMenu(comments, users, clipId, hasParrentComment, parentCommentId,
		buttonDataHead, firstCommentOnPageNumber)
	text := FormCommentsInfo(parentComment, parentCommentUser, firstCommentOnPageNumber, comments, users)

	return text, menu
}

func HandleTags(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	requestTelegramId int, userInfo *UserInfo, users Users, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	ownerIdParameter, _ := buttonDataHead.ReadParameter(OwnerIdParameterName)
	ownerId, _ := strconv.Atoi(ownerIdParameter.Value)

	firstTagOnPageNumberParameter, _ := buttonDataHead.ReadParameter(FirstTagOnPageNumberParameterName)
	firstTagOnPageNumber, _ := strconv.Atoi(firstTagOnPageNumberParameter.Value)

	tags, clipNumbers, err := GetTagsPage(sqlDb, ownerId, firstTagOnPageNumber)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	if len(tags) == 0 {
		t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
		return "Пока что нет ни одного тэга" + "\n\n" + t1, menu
	}

	menu := FormTagsPageMenu(tags, ownerId, buttonDataHead, firstTagOnPageNumber)
	text := FormTagsInfo(firstTagOnPageNumber, tags, clipNumbers)

	return text, menu
}


func HandleUser(sqlDb *sql.DB, requestTelegramId int, userInfo *UserInfo, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	userIdParameter, _ := buttonDataHead.ReadParameter(UserIdParameterName)
	userId, _ := strconv.Atoi(userIdParameter.Value)

	user, isRequestUserSubscribed, equalsRequestUser, err := GetUserDataByUserId(sqlDb, userId, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	return HandleUserData(user, isRequestUserSubscribed, equalsRequestUser)
}

func HandleSearchUserByUserNameMessage(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	requestTelegramId int, userName string, userInfo *UserInfo, users Users) (string, *tgbotapi.InlineKeyboardMarkup) {
	user, isRequestUserSubscribed, equalsRequestUser, err := GetUserDataByUserName(sqlDb, userName, requestTelegramId)
	if err != nil {
		if err == sql.ErrNoRows {
			t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
			return "Пользователь не найден" + "\n\n" + t1, menu
		}
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	return HandleUserData(user, isRequestUserSubscribed, equalsRequestUser)
}

func HandleUserData(user *User, isRequestUserSubscribed, equalsRequestUser bool) (string, *tgbotapi.InlineKeyboardMarkup) {

	menu := FormUserMenu(user.Id, isRequestUserSubscribed, equalsRequestUser)
	text := FormFullUserInfo(user)

	return text, menu

}

func HandleGroup(sqlDb *sql.DB, requestTelegramId int, userInfo *UserInfo, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	groupIdParameter, _ := buttonDataHead.ReadParameter(GroupIdParameterName)
	groupId, _ := strconv.Atoi(groupIdParameter.Value)

	group, creatorName, isRequestUserMember, isRequestUserGroupCreator, err := GetGroupData(sqlDb, groupId, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	menu := FormGroupMenu(groupId, isRequestUserMember, isRequestUserGroupCreator)
	text := FormFullGroupInfo(group, creatorName)

	return text, menu
}

func HandleClip(bot *tgbotapi.BotAPI, sqlDb *sql.DB, awsS3 *AwsS3, chatId int64,
	requestTelegramId int, userInfo *UserInfo, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	clipIdParameter, _ := buttonDataHead.ReadParameter(ClipIdParameterName)
	clipId, _ := strconv.Atoi(clipIdParameter.Value)

	clip, isRequestUserClipPublisher, err := GetClipDataByClipId(sqlDb, clipId, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	return HandleClipExtraData(bot, sqlDb, awsS3, userInfo, requestTelegramId, chatId, clip, isRequestUserClipPublisher)
}

func HandleSearchClipByUuidMessage(bot *tgbotapi.BotAPI, sqlDb *sql.DB, awsS3 *AwsS3, chatId int64,
	userInfo *UserInfo, requestTelegramId int, uuid string) (string, *tgbotapi.InlineKeyboardMarkup) {

	clip, isRequestUserClipPublisher, err := GetClipDataByUuid(sqlDb, uuid, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, "Клип с таким идентификатором не найден", sqlDb, requestTelegramId, userInfo)
	}

	return HandleClipExtraData(bot, sqlDb, awsS3, userInfo, requestTelegramId, chatId, clip,isRequestUserClipPublisher)
}

func HandleClipExtraData(bot *tgbotapi.BotAPI, sqlDb *sql.DB, awsS3 *AwsS3, userInfo *UserInfo, requestTelegramId int, chatId int64,
	clip *Clip, isRequestUserClipPublisher bool) (string, *tgbotapi.InlineKeyboardMarkup) {

	ownerName, isGroup, err := GetOwnerName(sqlDb, clip.OwnerId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	tags, err := GetClipTags(sqlDb, clip.Id)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	go HandleGetClipFromAwsS3(bot, sqlDb, awsS3, userInfo, requestTelegramId, chatId, clip, tags, ownerName, isGroup, isRequestUserClipPublisher)

	text := FormGetClipMessage()
	return text, makeEmptyMenu()
}

func HandleGetClipFromAwsS3(bot *tgbotapi.BotAPI, sqlDb *sql.DB, awsS3 *AwsS3, userInfo *UserInfo, requestTelegramId int, chatId int64,
	clip *Clip, tags []*Tag, ownerName string, isGroup, isRequestUserClipPublisher bool) {
	message := tgbotapi.NewMessage(chatId, "")
	message.ParseMode = "markdown"

	file, err := awsS3.DownloadClip(clip.Uuid)
	if err != nil {
		message.Text, message.ReplyMarkup = makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	} else {
		vidMsg := tgbotapi.VideoConfig{
			BaseFile:	tgbotapi.BaseFile{
				BaseChat: tgbotapi.BaseChat{ChatID: chatId},
				File:	tgbotapi.FileBytes{
					Name: clip.Name + "mp4",
					Bytes: file,
				},
			},
		}
		message.Text = FormFullClipInfo(clip, tags, ownerName)
		message.ReplyMarkup = FormClipMenu(clip.Id, clip.OwnerId, isGroup, isRequestUserClipPublisher)
		bot.Send(vidMsg)
	}

	bot.Send(message)
	//bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackQueryId, "Done"))
}

func HandleTag(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	userInfo *UserInfo, users Users, requestTelegramId int, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	ownerIdParameter, _ := buttonDataHead.ReadParameter(OwnerIdParameterName)
	ownerId, _ := strconv.Atoi(ownerIdParameter.Value)

	tagNameParameter, _ := buttonDataHead.ReadParameter(TagNameParameterName)

	firstClipOnPageNumberParameter, _ := buttonDataHead.ReadParameter(FirstClipOnPageNumberParameterName)
	firstClipOnPageNumber, _ := strconv.Atoi(firstClipOnPageNumberParameter.Value)

	clips, err := GetClipsPageByTag(sqlDb, ownerId, tagNameParameter.Value, firstClipOnPageNumber)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	ownerName, isGroup, err := GetOwnerName(sqlDb, ownerId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	if len(clips) == 0 {
		if isGroup {
			t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
			return fmt.Sprintf("В группе %s нет клипов с таким тэгом", ownerName) + "\n\n" + t1, menu
		} else {
			t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
			return fmt.Sprintf("Пользователь %s не имеет клипов с таким тэгом", ownerName) + "\n\n" + t1, menu
		}
	}

	menu := FormClipsPageMenu(clips, buttonDataHead, firstClipOnPageNumber)
	text := FormShortClipsInfo(firstClipOnPageNumber, clips, ownerName)

	return text, menu
}


func HandleMyClips(buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	userIdParameter, _ := buttonDataHead.ReadParameter(UserIdParameterName)
	userId, _ := strconv.Atoi(userIdParameter.Value)

	menu := FormMyClipsMenu(userId)
	text := FormMyClipsMessage()

	return text, menu
}

func HandleMyGroups(buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	userIdParameter, _ := buttonDataHead.ReadParameter(UserIdParameterName)
	userId, _ := strconv.Atoi(userIdParameter.Value)

	menu := FormMyGroupsMenu(userId)
	text := FormMyGroupsMessage()

	return text, menu
}

func HandleMySubscriptions(buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	userIdParameter, _ := buttonDataHead.ReadParameter(UserIdParameterName)
	userId, _ := strconv.Atoi(userIdParameter.Value)

	menu := FormMySubscriptionsMenu(userId)
	text := FormMySubscriptionsMessage()

	return text, menu
}

func HandleAddClip(sqlDb *sql.DB, userInfo *UserInfo, requestTelegramId int, buttonDataHead ButtonData, isBacked bool) (string, *tgbotapi.InlineKeyboardMarkup) {
	ownerIdParameter, _ := buttonDataHead.ReadParameter(OwnerIdParameterName)
	ownerId, _ := strconv.Atoi(ownerIdParameter.Value)

	user, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	if !isBacked {
		userInfo.UserState = AddClip
		userInfo.TempClip = new(Clip)
		userInfo.TempTags = make([]*Tag, 0)
		userInfo.TempClip.PublisherId = user.Id
		userInfo.TempClip.OwnerId = ownerId
	}

	menu := FormAddClipMenu()
	text := FormAddClipMessage()

	return text, menu
}

func HandleAddClipName(userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.UserState = AddClipName

	menu := FormSendTextMenu()
	text := FormAddClipNameMessage()

	return text, menu
}

func HandleAddClipTags(userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.UserState = AddClipTags

	menu := FormSendTextMenu()
	text := FormAddClipTagsMessage()

	return text, menu
}

func HandleAddClipDescription(userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.UserState = AddClipDescrpition

	menu := FormSendTextMenu()
	text := FormAddClipDescriptionMessage()

	return text, menu
}

func HandleSendClipVideoFile(userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.UserState = SendClipVideoFile

	menu := FormSendTextMenu()
	text := FormSendClipVideoFileMessage()

	return text, menu
}

func HandleCreateGroup(sqlDb *sql.DB, userInfo *UserInfo, requestTelegramId int, isBacked bool) (string, *tgbotapi.InlineKeyboardMarkup) {
	user, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	if !isBacked {
		userInfo.UserState = CreateNewGroup
		userInfo.TempGroup = new(Group)
		userInfo.TempGroup.CreatorId = user.Id
	}

	menu := FormCreateGroupMenu()
	text := FormCreateGroupMessage()

	return text, menu
}

func HandleAddGroupName(userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.UserState = AddGroupName

	menu := FormSendTextMenu()
	text := FormAddGroupNameMessage()

	return text, menu
}

func HandleAddGroupDescription(userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.UserState = AddGroupDescription

	menu := FormSendTextMenu()
	text := FormAddGroupDescriptionMessage() + "\n\n" + FormWelcomeMessage()

	return text, menu
}

func HandleFinishGroupCreation(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update, userInfo *UserInfo, users Users, requestTelegramId int) (string, *tgbotapi.InlineKeyboardMarkup) {
	user, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	_, err = GetGroupByCreatorAndName(sqlDb, user.Id, userInfo.TempGroup.Name)
	if err == nil {
		t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
		return "Вы уже создали группу с таким именем. Пожалуйста, выберете другое имя." + "\n\n" + t1, menu
	}

	ownerId, err := InsertOwner(sqlDb, true)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	groupId, err := InsertGroup(sqlDb, userInfo.TempGroup, ownerId)
	if err != nil {
		_ = DeleteOwner(sqlDb, ownerId)
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	err = InsertGroupMember(sqlDb, groupId, user.Id)

	menu := FormMainMenu(user.Id, userInfo)
	text := FormFinishGroupCreationMessage()

	return text, menu
}

func HandleAddComment(userInfo *UserInfo, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	var parentCommentId int
	var hasParentComment bool
	parentCommentIdParameter, err := buttonDataHead.ReadParameter(ParentCommentIdParameterName)
	if err == nil {
		//log.Println("\t\t\t\thave PARENT_PARAM")
		hasParentComment = true
		parentCommentId, err = strconv.Atoi(parentCommentIdParameter.Value)
	}

	clipIdParameter, _ := buttonDataHead.ReadParameter(ClipIdParameterName)
	clipId, _ := strconv.Atoi(clipIdParameter.Value)

	userInfo.UserState = AddComment
	userInfo.TempComment = new(Comment)
	userInfo.TempComment.ClipId = clipId
	if hasParentComment {
		userInfo.TempComment.ParentCommentId = sql.NullInt64{int64(parentCommentId), true}
	} else {
		userInfo.TempComment.ParentCommentId = sql.NullInt64{0, false}
	}

	menu := FormSendTextMenu()
	text := FormAddCommentMessage()

	return text, menu
}


func HandleEditUserProfile(buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	userIdParameter, _ := buttonDataHead.ReadParameter(UserIdParameterName)
	userId, _ := strconv.Atoi(userIdParameter.Value)

	menu := FormEditUserProfileMenu(userId)
	text := FormEditUserProfileMessage()
	return text, menu
}

func HandleEditUserDescription(userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.UserState = EditUserDescription

	menu := FormSendTextMenu()
	text := FormEditUserDescriptionMessage()
	return text, menu
}

func HandleEditGroupProfile(userInfo *UserInfo, buttonDataHead ButtonData, isBacked bool) (string, *tgbotapi.InlineKeyboardMarkup) {
	groupIdParameter, _ := buttonDataHead.ReadParameter(GroupIdParameterName)
	groupId, _ := strconv.Atoi(groupIdParameter.Value)

	if !isBacked {
		userInfo.TempGroup = new(Group)
		userInfo.TempGroup.Id = groupId
	}

	menu := FormEditGroupProfileMenu(groupId)
	text := FormEditGroupProfileMessage()
	return text, menu
}

func HandleEditGroupDescription(userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.UserState = EditGroupDescription

	menu := FormSendTextMenu()
	text := FormEditGroupDescriptionMessage()
	return text, menu
}

func HandleEditClipProfile(userInfo *UserInfo, buttonDataHead ButtonData, isBacked bool) (string, *tgbotapi.InlineKeyboardMarkup) {
	clipIdParameter, _ := buttonDataHead.ReadParameter(ClipIdParameterName)
	clipId, _ := strconv.Atoi(clipIdParameter.Value)

	if !isBacked {
		userInfo.TempClip = new(Clip)
		userInfo.TempClip.Id = clipId
	}

	menu := FormEditClipProfileMenu(clipId)
	text := FormEditClipProfileMessage()
	return text, menu
}

func HandleEditClipDescription(userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.UserState = EditClipDescription

	menu := FormSendTextMenu()
	text := FormEditClipDescriptionMessage()
	return text, menu
}


func HandleSubscribe(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI,
	update *tgbotapi.Update, userInfo *UserInfo, users Users, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {

	requestTelegramId := GetRequestTelegramIdFromUpdate(update)
	userIdParameter, _ := buttonDataHead.ReadParameter(UserIdParameterName)
	userId, _ := strconv.Atoi(userIdParameter.Value)

	requestUser, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	err = InsertSubscription(sqlDb, userId, requestUser.Id)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormInsertSubscribeMessage() + "\n\n" + t1
	return text, menu
}

func HandleUnsubscribe(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI,
	update *tgbotapi.Update, userInfo *UserInfo, users Users, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {

	requestTelegramId := GetRequestTelegramIdFromUpdate(update)
	userIdParameter, _ := buttonDataHead.ReadParameter(UserIdParameterName)
	userId, _ := strconv.Atoi(userIdParameter.Value)

	requestUser, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	err = DeleteSubscription(sqlDb, userId, requestUser.Id)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormDeleteSubscribeMessage() + "\n\n" + t1
	return text, menu
}

func HandleJoin(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI,
	update *tgbotapi.Update, userInfo *UserInfo, users Users, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {

	requestTelegramId := GetRequestTelegramIdFromUpdate(update)
	groupIdParameter, _ := buttonDataHead.ReadParameter(GroupIdParameterName)
	groupId, _ := strconv.Atoi(groupIdParameter.Value)

	requestUser, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	err = InsertGroupMember(sqlDb, groupId, requestUser.Id)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormInsertGroupMemberMessage() + "\n\n" + t1
	return text, menu
}

func HandleDisjoin(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI,
	update *tgbotapi.Update, userInfo *UserInfo, users Users, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {

	requestTelegramId := GetRequestTelegramIdFromUpdate(update)
	groupIdParameter, _ := buttonDataHead.ReadParameter(GroupIdParameterName)
	groupId, _ := strconv.Atoi(groupIdParameter.Value)

	requestUser, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	err = DeleteGroupMember(sqlDb, groupId, requestUser.Id)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormDeleteGroupMemberMessage() + "\n\n" + t1
	return text, menu
}

func HandleDeleteGroup(sqlDb *sql.DB, userInfo *UserInfo, requestTelegramId int, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	groupIdParameter, _ := buttonDataHead.ReadParameter(GroupIdParameterName)
	groupId, _ := strconv.Atoi(groupIdParameter.Value)

	err := DeleteGroup(sqlDb, groupId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	user, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	menu := FormMainMenu(user.Id, userInfo)
	text := FormDeleteGroupMessage()
	return text, menu
}

func HandleDeleteClip(bot *tgbotapi.BotAPI, sqlDb *sql.DB, awsS3 *AwsS3, chatId int64, userInfo *UserInfo,
	requestTelegramId int, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {
	clipIdParameter, _ := buttonDataHead.ReadParameter(ClipIdParameterName)
	clipId, _ := strconv.Atoi(clipIdParameter.Value)

	uuid, err := GetClipUUID(sqlDb, clipId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	user, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	go HandleDeleteClipFromAwsS3(bot, sqlDb, awsS3, user.Id, userInfo, requestTelegramId, chatId, uuid)

	err = DeleteClip(sqlDb, clipId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	text := FormWaitDeleteClipMessage()
	return text, makeEmptyMenu()
}

func HandleDeleteClipFromAwsS3(bot *tgbotapi.BotAPI, sqlDb *sql.DB, awsS3 *AwsS3, userId int, userInfo *UserInfo, requestTelegramId int, chatId int64,
	uuid string) {
	message := tgbotapi.NewMessage(chatId, "")
	message.ParseMode = "markdown"

	err := awsS3.DeleteClip(uuid)
	if err != nil {
		message.Text, message.ReplyMarkup = makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	message.Text = FormDeleteClipMessage()
	message.ReplyMarkup = FormMainMenu(userId, userInfo)
	bot.Send(message)
	//bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackQueryId, "Done"))
}

func HandleUploadClip(bot *tgbotapi.BotAPI, sqlDb *sql.DB, awsS3 *AwsS3, chatId int64, userInfo *UserInfo,
	requestTelegramId int, clip *Clip, tags []*Tag, videoInfo *tgbotapi.Video) (string, *tgbotapi.InlineKeyboardMarkup) {

	user, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	clip.Uuid = NewUuidString()
	clip.UploadTime = time.Now()
	clip.Duration = videoInfo.Duration

	fileId := videoInfo.FileID
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: fileId})

	clipId, err := InsertClip(sqlDb, clip)
	if err != nil {
		if e, ok := err.(*pq.Error); ok {
			if e.Code == DuplicatonInsertionError {
				return makeDbErrorResponseData(err, DbDuplicationErrorMessage, sqlDb, requestTelegramId, userInfo)
			}
		}
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	if len(tags) != 0 {
		err = InsertTags(sqlDb, clipId, tags)
		if err != nil {
			return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
		}
	}


	go HandleUploadClipToAwsS3(bot, sqlDb, awsS3, userInfo, requestTelegramId, chatId, file, clip.Uuid, user.Id)

	text := FormWaitUploadClipMessage()
	return text, makeEmptyMenu()
}

func HandleUploadClipToAwsS3(bot *tgbotapi.BotAPI, sqlDb *sql.DB, awsS3 *AwsS3, userInfo *UserInfo, requestTelegramId int, chatId int64,
	file tgbotapi.File, uuid string, userId int) {

	message := tgbotapi.NewMessage(chatId, "")
	message.ParseMode = "markdown"

	fileLink := file.Link(bot.Token)

	fileBuf, err := DownloadFile(fileLink)
	if err != nil {
		message.Text, message.ReplyMarkup = makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
		bot.Send(message)
		return
	}


	err = awsS3.UploadClip(uuid, fileBuf)
	if err != nil {
		message.Text, message.ReplyMarkup = makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
		bot.Send(message)
		return
	}

	message.Text = FormUploadClipMessage() + "\n\n" + FormWelcomeMessage()
	message.ReplyMarkup = FormMainMenu(userId, userInfo)
	bot.Send(message)
	//bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackQueryId, "Done"))
}

func HandleMainMenu(sqlDb *sql.DB, userInfo *UserInfo, requestTelegramId int, users Users) (string, *tgbotapi.InlineKeyboardMarkup) {
	user, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	users.User(requestTelegramId).StackedData = ""

	text := FormWelcomeMessage()
	menu := FormMainMenu(user.Id, userInfo)
	return text, menu
}

func HandleBack(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI,
	update *tgbotapi.Update, userInfo *UserInfo, users Users, buttonDataHead ButtonData) (string, *tgbotapi.InlineKeyboardMarkup) {

	requestTelegramId := GetRequestTelegramIdFromUpdate(update)
	newStackedData := GetTailFromStackedData(GetTailFromStackedData(userInfo.StackedData))
	if newStackedData == "" {
		return HandleMainMenu(sqlDb, userInfo, requestTelegramId, users)
	}

	newStackedDataHead := GetHeadFromStackedData(newStackedData)
	newButtonDataHead := ParseButtonData(newStackedDataHead)

	userInfo.StackedData = newStackedData

	switch userInfo.UserState {
	case AddClipTags, AddClipDescrpition, AddClipName:
		userInfo.UserState = AddClip
		return HandleAddClip(sqlDb, userInfo, requestTelegramId, newButtonDataHead, true)
	case AddGroupName, AddGroupDescription:
		userInfo.UserState = CreateNewGroup
		return HandleCreateGroup(sqlDb, userInfo, requestTelegramId, true)
	case EditGroupDescription:
		userInfo.UserState = EditGroupProfile
		return HandleEditGroupProfile(userInfo, newButtonDataHead, true)
	case EditClipDescription:
		userInfo.UserState = EditClipDescription
		return HandleEditClipProfile(userInfo, newButtonDataHead, true)
	}

	return ProcessButtonPush(sqlDb, awsS3, bot, update, userInfo, users, newButtonDataHead)
}


func UpdateCurrentMenu(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI,
	update *tgbotapi.Update, userInfo *UserInfo, users Users) (string, *tgbotapi.InlineKeyboardMarkup) {
	requestTelegramId := GetRequestTelegramIdFromUpdate(update)

	var newStackedData string
	if update.CallbackQuery != nil {
		newStackedData = GetTailFromStackedData(userInfo.StackedData)
	} else {
		newStackedData = userInfo.StackedData
	}
	//log.Println(newStackedData)

	if newStackedData == "" {
		return HandleMainMenu(sqlDb, userInfo, requestTelegramId, users)
	}

	newStackedDataHead := GetHeadFromStackedData(newStackedData)
	newButtonDataHead := ParseButtonData(newStackedDataHead)

	//log.Println(newStackedDataHead)
	userInfo.StackedData = newStackedData
	return ProcessButtonPush(sqlDb, awsS3, bot, update, userInfo, users, newButtonDataHead)
}

func GetRequestTelegramIdFromUpdate(update *tgbotapi.Update) int {
	if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	return update.Message.From.ID
}

func HandleAddClipNameMessage(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	userInfo *UserInfo, users Users, name string) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.TempClip.Name = name

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormAddClipNameResponseMessage() + "\n\n" + t1
	return text, menu
}

func HandleAddClipDescriptionMessage(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	userInfo *UserInfo, users Users, description string) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.TempClip.Description = description

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormAddClipDescriptionResponseMessage() + "\n\n" + t1
	return text, menu
}

func HandleAddClipTagsMessage(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	userInfo *UserInfo, users Users, tagsText string) (string, *tgbotapi.InlineKeyboardMarkup) {
	termTags, _ := ParseTags(tagsText, userInfo.TempClip.Id)
	userInfo.TempTags = termTags

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormAddClipTagsResponseMessage() + "\n\n" + t1
	return text, menu
}

func HandleAddGroupNameMessage(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	userInfo *UserInfo, users Users, name string) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.TempGroup.Name = name

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormAddGroupNameResponseMessage() + "\n\n" + t1
	return text, menu
}

func HandleAddGroupDescriptionMessage(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	userInfo *UserInfo, users Users, description string) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.TempGroup.Description = description

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormAddGroupDescriptionResponseMessage() + "\n\n" + t1
	return text, menu
}

func HandleAddCommentMessage(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	requestTelegramId int, userInfo *UserInfo, users Users, commentText string) (string, *tgbotapi.InlineKeyboardMarkup) {
	user, err := GetUserDataByTelegramId(sqlDb, requestTelegramId)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	userInfo.TempComment.PublisherId = sql.NullInt64{int64(user.Id), true}
	userInfo.TempComment.CommentText = commentText
	userInfo.TempComment.Time = time.Now()

	err = InsertComment(sqlDb, userInfo.TempComment)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormAddCommentResponseMessage() + "\n\n" + t1
	return text, menu
}

func HandleEditUserDescriptionMessage(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	userInfo *UserInfo, users Users, requestTelegramId int, description string) (string, *tgbotapi.InlineKeyboardMarkup) {
	err := UpdateUserDescription(sqlDb, requestTelegramId, description)
	if err != nil {
		log.Println("\t\tFail to update user description!")
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormEditUserDescriptionResponseMessage() + "\n\n" + t1
	return text, menu
}

func HandleEditClipDescriptionMessage(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	userInfo *UserInfo, users Users, requestTelegramId int, clipId int, description string) (string, *tgbotapi.InlineKeyboardMarkup) {
	err := UpdateClipDescription(sqlDb, clipId, description)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormEditClipDescriptionResponseMessage() + "\n\n" + t1
	return text, menu
}

func HandleEditGroupDescriptionMessage(sqlDb *sql.DB, awsS3 *AwsS3, bot *tgbotapi.BotAPI, update *tgbotapi.Update,
	userInfo *UserInfo, users Users, requestTelegramId int, groupId int, description string) (string, *tgbotapi.InlineKeyboardMarkup) {
	err := UpdateGroupDescription(sqlDb, groupId, description)
	if err != nil {
		return makeDbErrorResponseData(err, DbErrorMessage, sqlDb, requestTelegramId, userInfo)
	}

	t1, menu := UpdateCurrentMenu(sqlDb, awsS3, bot, update, userInfo, users)
	text := FormEditGroupDescriptionResponseMessage() + "\n\n" + t1
	return text, menu
}

func HandleSearch() (string, *tgbotapi.InlineKeyboardMarkup) {
	text := FormSearchMessage()
	menu := FormSearchMenu()
	return text, menu
}

func HandleSearchUser(userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.UserState = SearchUserByUserName

	text := FormSearchUserMessage()
	menu := FormSendTextMenu()
	return text, menu
}

func HandleSearchClip(userInfo *UserInfo) (string, *tgbotapi.InlineKeyboardMarkup) {
	userInfo.UserState = SearchClipByUuid

	text := FormSearchClipMessage()
	menu := FormSendTextMenu()
	return text, menu
}