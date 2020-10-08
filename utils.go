package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ParametersSplitter                    = " "
	ParameterNameValueSplitter            = "="
	ButtonDataSplitter                    = "/"
	ButtonName          				  = "bn"
	OwnerIdParameterName                  = "oid"
	FirstClipOnPageNumberParameterName    = "fclopn"
	UserIdParameterName                   = "uid"
	FirstGroupOnPageNumberParameterName   = "fgopn"
	FirstUserOnPageNumberParameterName    = "fuopn"
	FirstTagOnPageNumberParameterName     = "ftopn"
	ClipIdParameterName                   = "clid"
	ParentCommentIdParameterName          = "pcid"
	FirstCommentOnPageNumberParameterName = "fcopn"
	GroupIdParameterName                  = "gid"
	ClipButtonName                        = "cl"
	GroupButtonName                       = "g"
	UserButtonName                        = "u"
	CommentsButtonName                    = "cs"
	MyClipsButtonName                     = "mcls"
	MyGroupsButtonName                    = "mgs"
	MySubscriptionsButtonName             = "mss"
	MainMenuButtonName                    = "mm"
	AddClipButtonName                     = "acl"
	ClipsButtonName                       = "cls"
	CreateGroupButtonName                 = "crg"
	LeaveCommentButtonName                = "lc"
	UnsubscribeButtonName                 = "usub"
	SubscribeButtonName                   = "sub"
	GroupsButtonName                      = "gs"
	GroupMembersButtonName                = "gms"
	JoinGroupButtonName                   = "jg"
	LeaveGroupButtonName           		  = "lg"
	EditUserProfileButtonName      		  = "eup"
	EditUserDescriptionButtonName  		  = "eud"
	EditGroupProfileButtonName     		  = "egp"
	EditGroupDescriptionButtonName 		  = "egd"
	EditClipProfileButtonName      		  = "eclp"
	EditClipDescriptionButtonName 		  = "ecld"
	DeleteGroupButtonName        		  = "dg"
	SubscriptionsButtonName      		  = "subs"
	AddClipNameButtonName        		  = "acln"
	AddClipTagsButtonName        		  = "aclts"
	AddClipDescriptionButtonName 		  = "acld"
	SendClipVideoFileButtonName  		  = "sclvf"
	AddGroupNameButtonName       		  = "agn"
	AddGroupDescriptionButtonName		  = "agd"
	FinishGroupCreationButtonName		  = "fgcr"
	TagButtonName                		  = "t"
	TagNameParameterName         		  = "tn"
	SearchButtonName             		  = "se"
	SearchClipButtonName         		  = "secl"
	SearchUserButtonName         		  = "seu"
	BackButtonName               		  = "b"
	DeleteClipButtonName         		  = "dcl"
	TagsButtonName               		  = "ts"

	DbErrorMessage = "Извините, бот не справился с задачей."
	DbDuplicationErrorMessage = "Извините, но объект с такими значениями атрибутов уже существует в базе данных."
)

func FormShortClipsInfo(firstClipOnPageNumber int, clips []*Clip, ownerName string) string {
	text := ""
	clipsLen := len(clips)
	if clipsLen > ClipsPerPage {
		clipsLen = ClipsPerPage
	}
	for i := 0; i < clipsLen; i++ {
		text += FormShortClipInfo(firstClipOnPageNumber + i, clips[i], ownerName) + "\n\n"
	}
	return text
}

func FormShortClipInfo(i int, clip *Clip, ownerName string) string {
	return fmt.Sprintf("%d. Название: %s\n" +
		"Владелец: %s\n" +
		"Продолжительность: %02d:%02d\n" +
		"UUID: %s",
		i, clip.Name, ownerName, clip.Duration / 60, clip.Duration % 60, clip.Uuid)
}

func FormShortGroupsInfo(firstGroupOnPageNumber int, groups []*Group, creatorNames []string) string {
	text := ""
	groupsLen := len(groups)
	if groupsLen > GroupsPerPage {
		groupsLen = GroupsPerPage
	}
	for i := 0; i < groupsLen; i++ {
		text += FormShortGroupInfo(firstGroupOnPageNumber + i, groups[i], creatorNames[i]) + "\n\n"
	}
	return text
}

func FormShortGroupInfo(i int, group *Group, creatorName string) string {
	return fmt.Sprintf("%d. Название: %s\n" +
		"Создатель: %s", i, group.Name, creatorName)
}

func FormShortUsersInfo(firstUserOnPageNumber int, users []*User) string {
	text := ""
	usersLen := len(users)
	if usersLen > UsersPerPage {
		usersLen = UsersPerPage
	}
	for i := 0; i < usersLen; i++ {
		text += FormShortUserInfo(firstUserOnPageNumber + i, users[i]) + "\n\n"
	}
	return text
}

func FormShortUserInfo(i int, user *User) string {
	return fmt.Sprintf("%d. Имя: %s", i, user.Name)
}

func FormCommentsInfo(parentComment *Comment, parentCommentUser *User, firstCommentOnPageNumber int, comments []*Comment, users []*User) string {
	text := "Комментарии:\n"
	hasParentComment := parentComment != nil

	if hasParentComment {
		text += FormParentCommentInfo(parentComment, parentCommentUser) + "\n\n"
	}

	commentsLen := len(comments)
	if commentsLen > CommentsPerPage {
		commentsLen = CommentsPerPage
	}

	for i := 0; i < commentsLen; i++ {
		text += FormCommentInfo(firstCommentOnPageNumber + i, comments[i], users[i], hasParentComment) + "\n\n"
	}

	return text
}

func FormParentCommentInfo(comment *Comment, user *User) string {
	return fmt.Sprintf("Текст: %s\n" +
		"Дата публикации: %s\n" +
		"Пользователь: %s", comment.CommentText, comment.Time.String(), user.Name)
}

func FormCommentInfo(i int, comment *Comment, user *User, hasParentComent bool) string {
	if hasParentComent {
		return fmt.Sprintf("\t%d. Текст: %s\n" +
			"\tДата публикации: %s\n" +
			"\tПользователь: %s", i, comment.CommentText, comment.Time.String(), user.Name)
	} else {
		return fmt.Sprintf("%d. Текст: %s\n" +
			"Дата публикации: %s\n" +
			"Пользователь: %s", i, comment.CommentText, comment.Time.String(), user.Name)
	}
}

func FormMyClipsMessage() string {
	return "Опубликуйте новый клип или просматривайте уже добавленные клипы"
}

func FormMyGroupsMessage() string {
	return "Создайте группу или просматривайте группы, в которые Вы уже вступили"
}

func FormMySubscriptionsMessage() string {
	return "Просматривайте пользователей, на которых Вы подписались"
}

func FormAddClipMessage() string {
	return "Добавьте название клипа(обязательно), его описание или тэги"
}

func FormAddClipNameMessage() string {
	return "Отправьте название клипа"
}

func FormAddClipTagsMessage() string {
	return "Перечислите УНИКАЛЬНЫЕ тэги, каждый должен начинаться с символа решетки '#'." +
		"Тэг может содержать буквы и цифры. Тэги должны быть разделены одним пробельным символом." +
		"(пример: #первыйТэг #2Тэг #ТретийТэг)"
}

func FormAddClipDescriptionMessage() string {
	return "Отправьте описание клипа"
}

func FormSendClipVideoFileMessage() string {
	return "Отправьте клип медиа файлом, это завершит процесс публикации клипа."
}

func FormCreateGroupMessage() string {
	return "Добавьте название группы(обязательно) и описание."
}

func FormAddGroupNameMessage() string {
	return "Отправьте название группы."
}

func FormAddGroupDescriptionMessage() string {
	return "Отправьте описание группы."
}

func FormFinishGroupCreationMessage() string {
	return "Группа успешно создана!"
}

func FormAddCommentMessage() string {
	return "Отправьте текст комментария"
}

func FormEditUserProfileMessage() string {
	return "Выберете Инфо для изменения"
}

func FormEditUserDescriptionMessage() string {
	return "Отправьте новый текст для Инфо"
}

func FormEditGroupProfileMessage() string {
	return "Выберите описание группы для изменения"
}

func FormEditGroupDescriptionMessage() string {
	return "Отправьте новое описание группы"
}

func FormEditClipProfileMessage() string {
	return "Выберите описание клипа для его изменения"
}

func FormEditClipDescriptionMessage() string {
	return "Отправьте новое описание клипа"
}

func FormInsertSubscribeMessage() string {
	return "Вы успешно подписались!"
}

func FormInsertGroupMemberMessage() string {
	return "Теперь вы новый участник группы!"
}

func FormDeleteSubscribeMessage() string {
	return "Вы успешно отписались!"
}

func FormDeleteGroupMemberMessage() string {
	return "Вы более не участник группы!"
}

func FormDeleteGroupMessage() string {
	return "Группа успешно удалена!"
}

func FormWaitDeleteClipMessage() string {
	return "Пожалуйста, дождитесь завершения процесса удаления клипа"
}

func FormDeleteClipMessage() string {
	return "Клип успешно удален!"
}

func FormWaitUploadClipMessage() string {
	return "Пожалуйста, дождитесь завершения процесса публикации клипа"
}

func FormUploadClipMessage() string {
	return "Клип успешно опубликован!"
}

func FormGetClipMessage() string {
	return "Пожалуйста, дождитесь завершения процесса скачивания клипа"
}

func FormAddClipNameResponseMessage() string {
	return "Название клипа успешно добавлено!"
}

func FormAddClipDescriptionResponseMessage() string {
	return "Описание клипа успешно добавлено!"
}

func FormAddClipTagsResponseMessage() string {
	return "Тэги клипа успешно добавлены!"
}

func FormAddGroupNameResponseMessage() string {
	return "Название группы успешно добавлено!"
}

func FormAddGroupDescriptionResponseMessage() string {
	return "Описание группы успешно добавлено!"
}

func FormAddCommentResponseMessage() string {
	return "Коммент успешно добавлен!\n" +
		"Пожалуйста, теперь нажмите кнопку \"Назад\""
}

func FormEditUserDescriptionResponseMessage() string {
	return "Инфо пользователя успешно изменено!"
}

func FormEditGroupDescriptionResponseMessage() string {
	return "Описание группы успешно изменено!"
}

func FormEditClipDescriptionResponseMessage() string {
	return "Описание клипа успешно изменено!"
}

func FormSearchMessage() string {
	return "Поиск пользователей по имени или клипов по идентификатору."
}

func FormSearchUserMessage() string {
	return "Отправьте имя пользователя для поиска."
}

func FormSearchClipMessage() string {
	return "Отправьте идентификатор клипа для поиска."
}


func FormFullUserInfo(user *User) string {
	return fmt.Sprintf("Имя пользователя: %s\n" +
		"Инфо: %s\n", user.Name, user.Info)
}

func FormFullClipInfo(clip *Clip, tags []*Tag, ownerName string) string {
	if len(tags) > 0 {
		tagsInfo := ""
		for i := 0; i < len(tags)-1; i++ {
			tagsInfo += fmt.Sprintf("#%s, ", tags[i].Name)
		}
		tagsInfo += fmt.Sprintf("#%s", tags[len(tags)-1].Name)

		return fmt.Sprintf("Название: %s\n" +
			"Тэги: %s\n" +
			"Описание: %s\n" +
			"Владелец: %s\n" +
			"Продолжительность: %02d:%02d\n" +
			"UUID: %s",
			clip.Name, tagsInfo, clip.Description, ownerName, clip.Duration / 60, clip.Duration % 60, clip.Uuid)
	}

	return fmt.Sprintf("Название: %s\n" +
		"Описание: %s\n" +
		"Владелецr: %s\n" +
		"Продолительность: %02d:%02d\n" +
		"UUID: %s",
		clip.Name, clip.Description, ownerName, clip.Duration / 60, clip.Duration % 60, clip.Uuid)
}

func FormTagsInfo(firstTagOnPage int, tags []*Tag, clipNumbers [] int) string {
	text := ""
	tagsLen := len(tags)
	if tagsLen > TagsPerPage {
		tagsLen = TagsPerPage
	}
	for i := 0; i < tagsLen; i++ {
		text += FormTagInfo(firstTagOnPage + i, tags[i], clipNumbers[i]) + "\n\n"
	}
	return text
}

func FormTagInfo(i int, tag *Tag, clipNumber int) string {
	return fmt.Sprintf("%d. %s --- %d клипов", i, tag.Name, clipNumber)
}

func FormFullGroupInfo(group *Group, creatorName string) string {
	return fmt.Sprintf("Название: %s\n" +
		"Описание: %s\n" +
		"Создатель: %s", group.Name, group.Description, creatorName)
}

func FormUserRegistrationMessage(user *User) string {
	return fmt.Sprintf("Привет, %s, кажется, Вы здесь первый раз. Вас только что зарегестрировали как нового пользователя.\n" +
		"Наслаждайтесь использованием ClipSpace бота!\n\n", user.Name) + FormWelcomeMessage()
}

func FormWelcomeMessage() string {
	return "Добро пожаловать в ClipSpace бота!\n\nЗдесь Вы сможете публиковать свои клипы," +
		"чтобы хранить их удаленно и делиться с кем Вам захочется посредством короткого идентификатора." +
		"Вашему другу достаточно отправить идентификатор этому боту, чтобы увидеть видео!\n\n" +
		"Вы также можете подписываться на других пользователей, чтобы с удобством просматривать их клипы и группы.\n\n" +
		"В группе же Вы можете найти клипы различных пользователей, которые являются её участниками.\n\n" +
		"Добавляйте к клипам тэги, раскрывающие их содержание, и обсуждайте их в комментариях!"
	return "Welcome to the ClipSpace bot!\n\nHere you can cut a video to be just on spot, " +
		"than store it remotely and share with whoever you want via simple identifier, " +
		"that just needed to be send to this bot by the person to see it!\n\n" +
		"You can also subscribe other users to scroll through their clips and be notified " +
		"of new clips they post here.\n\nGive clips suitable tags to give a hint of it's content " +
		"and discuss them in comments!"
}


func DownloadFile(fileLink string) ([]byte, error) {
	//log.Println(fileLink)
	resp, err := http.Get(fileLink)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", resp.Status)
	}

	var out []byte
	out, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func ParseTags(text string, clipId int) ([]*Tag, error) {
	parts := strings.Split(text, " ")

	tags := make([]*Tag, 0)
	for i := range parts {
		if len(parts[i]) > 2 && parts[i][0] == '#' {
			tags = append(tags, &Tag{ClipId: clipId, Name: parts[i]})
		}
	}

	return tags, nil
}



type DataParameter struct {
	Name string
	Value string
}

type ButtonData map[string]string

func InitButtonData(parameters ... DataParameter) ButtonData {
	bd := make(ButtonData)
	bd.WriteParameters(parameters...)
	return bd
}

func (dp DataParameter) String() string {
	return dp.Name + "=" + dp.Value
}

func SplitStackedData(stackedData string) (string, string) {
	splits := strings.Split(stackedData, "/")
	head, tail := splits[len(splits)-1], strings.Join(splits[:len(splits)-1], "/")
	return head, tail
}

func GetHeadFromStackedData(stackedData string) string {
	splits := strings.Split(stackedData, "/")
	head := splits[len(splits)-1]
	head = "/" + head
	return head
}

func GetTailFromStackedData(stackedData string) string {
	splits := strings.Split(stackedData, "/")
	tail := strings.Join(splits[:len(splits)-1], "/")
	return tail
}


func (bd ButtonData) Copy() ButtonData {
	var newBd = make(ButtonData)
	for k, v := range bd {
		newBd[k] = v
	}
	return newBd
}

func ParseButtonData(data string) ButtonData {
	bd := make(ButtonData)
	parameters := strings.Split(data[1:], ParametersSplitter)
	for i := 0; i < len(parameters); i++ {
		parts := strings.Split(parameters[i], ParameterNameValueSplitter)
		if len(parts) == 2 {
			bd[parts[0]] = parts[1]
		}
	}
	return bd
}

func (bd ButtonData) ReadParameter(parameterName string) (DataParameter, error) {
	if value, ok := bd[parameterName]; ok {
		return DataParameter{parameterName, value}, nil
	}

	return DataParameter{parameterName, ""}, fmt.Errorf("no parameter named %s in data string", parameterName)
}

func (bd ButtonData) WriteParameters(parameters ... DataParameter) {
	for _, parameter := range parameters {
		bd[parameter.Name] = parameter.Value
	}
}

func (bd ButtonData) String() string {
	str := ButtonDataSplitter
	for k, v := range bd {
		str += fmt.Sprintf("%s%s%s ", k, ParameterNameValueSplitter, v)
	}

	if str[len(str)-1] == ' ' {
		str = str[:len(str)-1]
	}
	return str
}