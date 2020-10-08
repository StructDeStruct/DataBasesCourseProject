package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
)


func FormMainMenu(userid int, userInfo *UserInfo) *tgbotapi.InlineKeyboardMarkup {
	userInfo.StackedData = ""

	menu := &tgbotapi.InlineKeyboardMarkup{}
	var buttonData ButtonData

	buttonData = InitButtonData(DataParameter{
		Name: ButtonName,
		Value: MyClipsButtonName,
	},
	DataParameter{
		Name: UserIdParameterName,
		Value: strconv.Itoa(userid),
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Мои клипы", buttonData.String())))

	buttonData = InitButtonData(DataParameter{
		Name: ButtonName,
		Value: MyGroupsButtonName,
	},
		DataParameter{
			Name: UserIdParameterName,
			Value: strconv.Itoa(userid),
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Мои группы", buttonData.String())))

	buttonData = InitButtonData(DataParameter{
		Name:  ButtonName,
		Value: SubscriptionsButtonName,
	},
		DataParameter{
			Name: UserIdParameterName,
			Value: strconv.Itoa(userid),
		},
		DataParameter{
			Name:  FirstUserOnPageNumberParameterName,
			Value: "1",
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Мои подписки", buttonData.String())))

	buttonData = InitButtonData(DataParameter{
		Name: ButtonName,
		Value: UserButtonName,
	},
		DataParameter{
			Name: UserIdParameterName,
			Value: strconv.Itoa(userid),
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Профиль", buttonData.String())))

	buttonData = InitButtonData(DataParameter{
		Name: ButtonName,
		Value: SearchButtonName,
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Поиск", buttonData.String())))

	return menu
}

func FormMyClipsMenu(userId int) *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: AddClipButtonName,
	},
	DataParameter{
		Name: OwnerIdParameterName,
		Value: strconv.Itoa(userId),
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить новый клип", buttonData.String())))


	AddClipsButton(menu, userId, "Клипы")
	AddTagsButton(menu, userId, "Тэги")

	AddMenuButton(menu)
	AddBackButton(menu)

	return menu
}

func FormAddClipMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: AddClipNameButtonName,
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить имя клипа", buttonData.String())))

	buttonData = InitButtonData(DataParameter{
		Name: ButtonName,
		Value: AddClipTagsButtonName,
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить тэги клипа", buttonData.String())))

	buttonData = InitButtonData(DataParameter{
		Name: ButtonName,
		Value: AddClipDescriptionButtonName,
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить описание клипа", buttonData.String())))

	buttonData = InitButtonData(DataParameter{
		Name:  ButtonName,
		Value: SendClipVideoFileButtonName,
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Отпарвить видеофайл", buttonData.String())))

	AddMenuButton(menu)
	AddBackButton(menu)

	return menu
}

func FormCreateGroupMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: AddGroupNameButtonName,
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить имя группы", buttonData.String())))

	buttonData = InitButtonData(DataParameter{
		Name: ButtonName,
		Value: AddGroupDescriptionButtonName,
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Добавить описание группы", buttonData.String())))

	buttonData = InitButtonData(DataParameter{
		Name:  ButtonName,
		Value: FinishGroupCreationButtonName,
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Создать группу", buttonData.String())))

	AddMenuButton(menu)
	AddBackButton(menu)

	return menu
}

func AddClipsButton(menu *tgbotapi.InlineKeyboardMarkup, ownerId int, buttonText string) {
	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: ClipsButtonName,
	},
		DataParameter{
			Name: OwnerIdParameterName,
			Value: strconv.Itoa(ownerId),
		},
		DataParameter{
			Name: FirstClipOnPageNumberParameterName,
			Value: "1",
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonData.String())))
}

func AddGroupsButton(menu *tgbotapi.InlineKeyboardMarkup, userId int, buttonText string) {
	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: GroupsButtonName,
	},
		DataParameter{
			Name: UserIdParameterName,
			Value: strconv.Itoa(userId),
		},
		DataParameter{
			Name: FirstGroupOnPageNumberParameterName,
			Value: "1",
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonData.String())))
}

func AddSubscriptionsButton(menu *tgbotapi.InlineKeyboardMarkup, userId int, buttonText string) {
	buttonData := InitButtonData(DataParameter{
		Name:  ButtonName,
		Value: SubscriptionsButtonName,
	},
		DataParameter{
			Name: UserIdParameterName,
			Value: strconv.Itoa(userId),
		},
		DataParameter{
			Name:  FirstUserOnPageNumberParameterName,
			Value: "1",
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonData.String())))
}

func AddGroupMembersButton(menu *tgbotapi.InlineKeyboardMarkup, groupId int, buttonText string) {
	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: GroupMembersButtonName,
	},
		DataParameter{
			Name: GroupIdParameterName,
			Value: strconv.Itoa(groupId),
		},
		DataParameter{
			Name: FirstUserOnPageNumberParameterName,
			Value: "1",
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonData.String())))
}

func AddTagsButton(menu *tgbotapi.InlineKeyboardMarkup, ownerId int, buttonText string) {
	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: TagsButtonName,
	},
		DataParameter{
			Name: OwnerIdParameterName,
			Value: strconv.Itoa(ownerId),
		},
		DataParameter{
			Name: FirstTagOnPageNumberParameterName,
			Value: "1",
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonData.String())))
}

func AddCommentsButton(menu *tgbotapi.InlineKeyboardMarkup, clipId int, buttonText string) {
	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: CommentsButtonName,
	},
		DataParameter{
			Name: ClipIdParameterName,
			Value: strconv.Itoa(clipId),
		},
		DataParameter{
			Name: FirstCommentOnPageNumberParameterName,
			Value: "1",
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonData.String())))
}

func FormMyGroupsMenu(userId int) * tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	buttonData := InitButtonData(DataParameter{
		Name:  ButtonName,
		Value: CreateGroupButtonName,
	},
		DataParameter{
			Name: UserIdParameterName,
			Value: strconv.Itoa(userId),
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Создать новую группу", buttonData.String())))


	AddGroupsButton(menu, userId, "Группы")
	AddMenuButton(menu)
	AddBackButton(menu)

	return menu
}

func FormMySubscriptionsMenu(userId int) *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	AddSubscriptionsButton(menu, userId, "Подписки")
	AddMenuButton(menu)
	AddBackButton(menu)

	return menu
}

func FormSendTextMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	AddMenuButton(menu)
	AddBackButton(menu)

	return menu
}

func FormEditUserProfileMenu(userId int) *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: EditUserDescriptionButtonName,
	},
		DataParameter{
			Name: UserIdParameterName,
			Value: strconv.Itoa(userId),
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Изменить инфо", buttonData.String())))

	AddMenuButton(menu)
	AddBackButton(menu)

	return menu
}

func FormEditGroupProfileMenu(groupId int) *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: EditGroupDescriptionButtonName,
	},
		DataParameter{
			Name: GroupIdParameterName,
			Value: strconv.Itoa(groupId),
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Изменить описание группы", buttonData.String())))

	AddMenuButton(menu)
	AddBackButton(menu)

	return menu
}

func FormEditClipProfileMenu(clipId int) *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: EditClipDescriptionButtonName,
	},
		DataParameter{
			Name: ClipIdParameterName,
			Value: strconv.Itoa(clipId),
		})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Имзменить описание клипа", buttonData.String())))

	AddMenuButton(menu)
	AddBackButton(menu)

	return menu
}

func AddMenuButton(menu *tgbotapi.InlineKeyboardMarkup) {
	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: MainMenuButtonName,
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Главное меню", buttonData.String())))
}

func AddBackButton(menu *tgbotapi.InlineKeyboardMarkup) {
	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: BackButtonName,
	})
	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Назад", buttonData.String())))
}

func FormClipsPageMenu(clips []*Clip, buttonDataHead ButtonData, firstClipOnPageNumber int) *tgbotapi.InlineKeyboardMarkup {
	var prevButt, nextButt bool
	itemsLen := len(clips)
	if len(clips) > ClipsPerPage {
		itemsLen = ClipsPerPage
		nextButt = true
	}
	if firstClipOnPageNumber != 1 {
		prevButt = true
	}

	var clipButtons []tgbotapi.InlineKeyboardButton
	for i := 0; i < itemsLen; i++ {
		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: ClipButtonName,
		},
			DataParameter{
				Name: ClipIdParameterName,
				Value: strconv.Itoa(clips[i].Id),
			})
		clipButtons = append(clipButtons, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. %s", firstClipOnPageNumber + i, clips[i].Name),
			buttonData.String()))

	}

	menu := &tgbotapi.InlineKeyboardMarkup{}
	AddPageMenu(menu, clipButtons, ClipsPerPage, FirstClipOnPageNumberParameterName, buttonDataHead, prevButt, nextButt)
	AddMenuButton(menu)
	AddBackButton(menu)
	return menu
}


func FormGroupsPageMenu(groups []*Group, buttonDataHead ButtonData, firstClipOnPageNumber int) *tgbotapi.InlineKeyboardMarkup {
	var prevButt, nextButt bool
	itemsLen := len(groups)
	if len(groups) > GroupsPerPage {
		itemsLen = GroupsPerPage
		nextButt = true
	}
	if firstClipOnPageNumber != 1 {
		prevButt = true
	}

	var groupButtons []tgbotapi.InlineKeyboardButton
	for i := 0; i < itemsLen; i++ {
		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: GroupButtonName,
		},
			DataParameter{
				Name: GroupIdParameterName,
				Value: strconv.Itoa(groups[i].Id),
			})
		groupButtons = append(groupButtons, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. %s", firstClipOnPageNumber + i, groups[i].Name),
			buttonData.String()))

	}

	menu := &tgbotapi.InlineKeyboardMarkup{}
	AddPageMenu(menu, groupButtons, GroupsPerPage, FirstGroupOnPageNumberParameterName, buttonDataHead, prevButt, nextButt)
	AddMenuButton(menu)
	AddBackButton(menu)
	return menu
}


func FormSubscriptionsPageMenu(users []*User, buttonDataHead ButtonData, firstSubscriptionOnPageNumber int) *tgbotapi.InlineKeyboardMarkup {
	var prevButt, nextButt bool
	itemsLen := len(users)
	if len(users) > UsersPerPage {
		itemsLen = UsersPerPage
		nextButt = true
	}
	if firstSubscriptionOnPageNumber != 1 {
		prevButt = true
	}

	var usersButtons []tgbotapi.InlineKeyboardButton
	for i := 0; i < itemsLen; i++ {
		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: UserButtonName,
		},
			DataParameter{
				Name: UserIdParameterName,
				Value: strconv.Itoa(users[i].Id),
			})
		usersButtons = append(usersButtons, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. %s", firstSubscriptionOnPageNumber+ i, users[i].Name),
			buttonData.String()))

	}

	menu := &tgbotapi.InlineKeyboardMarkup{}
	AddPageMenu(menu, usersButtons, UsersPerPage, FirstUserOnPageNumberParameterName, buttonDataHead, prevButt, nextButt)
	AddMenuButton(menu)
	AddBackButton(menu)
	return menu
}

func FormGroupMembersPageMenu(users []*User, buttonDataHead ButtonData, firstGroupMemberOnPageNumber int) *tgbotapi.InlineKeyboardMarkup {
	var prevButt, nextButt bool
	itemsLen := len(users)
	if len(users) > UsersPerPage {
		itemsLen = UsersPerPage
		nextButt = true
	}
	if firstGroupMemberOnPageNumber != 1 {
		prevButt = true
	}

	var usersButtons []tgbotapi.InlineKeyboardButton
	for i := 0; i < itemsLen; i++ {
		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: UserButtonName,
		},
			DataParameter{
				Name: UserIdParameterName,
				Value: strconv.Itoa(users[i].Id),
			})
		usersButtons = append(usersButtons, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. %s", firstGroupMemberOnPageNumber+ i, users[i].Name),
			buttonData.String()))

	}

	menu := &tgbotapi.InlineKeyboardMarkup{}
	AddPageMenu(menu, usersButtons, UsersPerPage, FirstUserOnPageNumberParameterName, buttonDataHead, prevButt, nextButt)
	AddMenuButton(menu)
	AddBackButton(menu)
	return menu
}

func FormCommentsPageMenu(comments []*Comment, users []*User, clipid int, hasParentComment bool, parentCommentId int,
	buttonDataHead ButtonData, firstCommentOnPageNumber int) *tgbotapi.InlineKeyboardMarkup {
	var prevButt, nextButt bool
	itemsLen := len(users)
	if len(users) > CommentsPerPage {
		itemsLen = CommentsPerPage
		nextButt = true
	}
	if firstCommentOnPageNumber != 1 {
		prevButt = true
	}

	menu := &tgbotapi.InlineKeyboardMarkup{}

	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: LeaveCommentButtonName,
	},
		DataParameter{
			Name: ClipIdParameterName,
			Value: strconv.Itoa(clipid),
		})

	buttonText := "Оставить комментарий"
	if hasParentComment {
		buttonData.WriteParameters(DataParameter{
			Name: ParentCommentIdParameterName,
			Value: strconv.Itoa(parentCommentId),
		})

		buttonText = "Ответить на комментарий"
	}

	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(buttonText, buttonData.String())))

	var commentsButtons []tgbotapi.InlineKeyboardButton
	for i := 0; i < itemsLen; i++ {
		buttonData = InitButtonData(DataParameter{
			Name: ButtonName,
			Value: CommentsButtonName,
		},
			DataParameter{
				Name: ClipIdParameterName,
				Value: strconv.Itoa(comments[i].ClipId),
			},
			DataParameter{
				Name: ParentCommentIdParameterName,
				Value: strconv.Itoa(comments[i].Id),
			},
			DataParameter{
				Name: FirstCommentOnPageNumberParameterName,
				Value: "1",
			})

		commentsButtons = append(commentsButtons, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. %s", firstCommentOnPageNumber + i, "Subcomment"),
			buttonData.String()))
	}


	AddPageMenu(menu, commentsButtons, CommentsPerPage, FirstCommentOnPageNumberParameterName, buttonDataHead, prevButt, nextButt)
	AddMenuButton(menu)
	AddBackButton(menu)
	return menu
}

func FormTagsPageMenu(tags []*Tag, ownerId int, buttonDataHead ButtonData, firstTagOnPageNumber int) *tgbotapi.InlineKeyboardMarkup {
	var prevButt, nextButt bool
	itemsLen := len(tags)
	if len(tags) > TagsPerPage {
		itemsLen = TagsPerPage
		nextButt = true
	}
	if firstTagOnPageNumber != 1 {
		prevButt = true
	}

	var usersButtons []tgbotapi.InlineKeyboardButton
	for i := 0; i < itemsLen; i++ {
		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: TagButtonName,
		},
			DataParameter{
				Name: OwnerIdParameterName,
				Value: strconv.Itoa(ownerId),
			},
			DataParameter{
				Name:  TagNameParameterName,
				Value: tags[i].Name,
		},
			DataParameter{
				Name: FirstClipOnPageNumberParameterName,
				Value: "1",
			})
		usersButtons = append(usersButtons, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%d. %s", firstTagOnPageNumber+ i, tags[i].Name),
			buttonData.String()))
	}

	menu := &tgbotapi.InlineKeyboardMarkup{}
	AddPageMenu(menu, usersButtons, UsersPerPage, FirstTagOnPageNumberParameterName, buttonDataHead, prevButt, nextButt)
	AddMenuButton(menu)
	AddBackButton(menu)
	return menu
}


func FormUserMenu(userid int, isRequestUserSubscribed, equalsRequestUser bool) *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	AddClipsButton(menu, userid, "Клипы")
	AddTagsButton(menu, userid, "Тэги")
	AddGroupsButton(menu, userid, "Группы")
	AddSubscriptionsButton(menu, userid, "Подписки")

	if equalsRequestUser {

		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: EditUserProfileButtonName,
		},
		DataParameter{
			Name: UserIdParameterName,
			Value: strconv.Itoa(userid),
		})
		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить", buttonData.String())))
	} else {
		if isRequestUserSubscribed {
			buttonData := InitButtonData(DataParameter{
				Name: ButtonName,
				Value: UnsubscribeButtonName,
			},
				DataParameter{
					Name: UserIdParameterName,
					Value: strconv.Itoa(userid),
				})

			menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Отписаться", buttonData.String())))
		} else {
			buttonData := InitButtonData(DataParameter{
				Name: ButtonName,
				Value: SubscribeButtonName,
			},
				DataParameter{
					Name: UserIdParameterName,
					Value: strconv.Itoa(userid),
				})

			menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Подписаться", buttonData.String())))
		}
	}

	AddMenuButton(menu)
	AddBackButton(menu)
	return menu
}

func FormGroupMenu(groupId int, isRequestUserMember, isRequestUserGroupCreator bool) *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	AddClipsButton(menu, groupId, "Клипы")
	AddTagsButton(menu, groupId, "Тэги")
	AddGroupMembersButton(menu, groupId, "Участники")

	if isRequestUserGroupCreator {
		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: DeleteGroupButtonName,
		},
		DataParameter{
			Name: GroupIdParameterName,
			Value: strconv.Itoa(groupId),
		})

		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить группу", buttonData.String())))
	}
	if isRequestUserMember  {
		if !isRequestUserGroupCreator {
			buttonData := InitButtonData(DataParameter{
				Name:  ButtonName,
				Value: LeaveGroupButtonName,
			},
				DataParameter{
					Name:  GroupIdParameterName,
					Value: strconv.Itoa(groupId),
				})

			menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("Покинуть", buttonData.String())))
		}

		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: AddClipButtonName,
		},
			DataParameter{
				Name: OwnerIdParameterName,
				Value: strconv.Itoa(groupId),
			})

		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Добавить клип", buttonData.String())))

	} else {
		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: JoinGroupButtonName,
		},
			DataParameter{
				Name: GroupIdParameterName,
				Value: strconv.Itoa(groupId),
			})

		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Вступить", buttonData.String())))
	}



	AddMenuButton(menu)
	AddBackButton(menu)
	return menu
}

func FormClipMenu(clipId int, ownerId int, isGroup, isRequestUserClipPublisher bool) *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	if isGroup {
		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: GroupButtonName,
		},
			DataParameter{
				Name: GroupIdParameterName,
				Value: strconv.Itoa(ownerId),
			})

		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Группа", buttonData.String())))
	} else {
		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: UserButtonName,
		},
			DataParameter{
				Name: UserIdParameterName,
				Value: strconv.Itoa(ownerId),
			})

		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Пользователь", buttonData.String())))
	}

	AddCommentsButton(menu, clipId, "Комментарии")

	if isRequestUserClipPublisher {
		buttonData := InitButtonData(DataParameter{
			Name: ButtonName,
			Value: EditClipProfileButtonName,
		},
			DataParameter{
				Name: ClipIdParameterName,
				Value: strconv.Itoa(clipId),
			})

		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить", buttonData.String())))

		buttonData = InitButtonData(DataParameter{
			Name: ButtonName,
			Value: DeleteClipButtonName,
		},
			DataParameter{
				Name: ClipIdParameterName,
				Value: strconv.Itoa(clipId),
			})

		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить клип", buttonData.String())))
	}


	AddMenuButton(menu)
	AddBackButton(menu)
	return menu
}

func FormSearchMenu() *tgbotapi.InlineKeyboardMarkup {
	menu := &tgbotapi.InlineKeyboardMarkup{}

	buttonData := InitButtonData(DataParameter{
		Name: ButtonName,
		Value: SearchUserButtonName,
	})

	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Поиск пользователя", buttonData.String())))

	buttonData = InitButtonData(DataParameter{
		Name: ButtonName,
		Value: SearchClipButtonName,
	})

	menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Поиск клипа", buttonData.String())))


	AddMenuButton(menu)
	AddBackButton(menu)
	return menu
}




func FormPreviousPageButtonData(currentPageButtonData ButtonData, itemsPerPage int, firstItemOnPageNumberParameterName string) ButtonData {
	previousPageButtonData := currentPageButtonData.Copy()

	firstClipOnPageNumberParameter, _ := currentPageButtonData.ReadParameter(firstItemOnPageNumberParameterName)
	firstClipOnPageNumber, _ := strconv.Atoi(firstClipOnPageNumberParameter.Value)

	previousPageButtonData.WriteParameters(DataParameter{
		Name:firstItemOnPageNumberParameterName,
		Value: strconv.Itoa(firstClipOnPageNumber - itemsPerPage)})

	return previousPageButtonData
}

func FormNextPageButtonData(currentPageButtonData ButtonData, itemsPerPage int, firstItemOnPageNumberParameterName string) ButtonData {
	nextPageButtonData := currentPageButtonData.Copy()

	firstClipOnPageNumberParameter, _ := currentPageButtonData.ReadParameter(firstItemOnPageNumberParameterName)
	firstClipOnPageNumber, _ := strconv.Atoi(firstClipOnPageNumberParameter.Value)

	nextPageButtonData.WriteParameters(DataParameter{
		Name: firstItemOnPageNumberParameterName,
		Value: strconv.Itoa(firstClipOnPageNumber + itemsPerPage)})

	return nextPageButtonData
}

func AddPageMenu(menu *tgbotapi.InlineKeyboardMarkup, Items []tgbotapi.InlineKeyboardButton, itemsPerPage int,
	firstItemOnPageNumberParameterName string, ButtonDataHead ButtonData, prevButt, nextButt bool) {
	for i := 0; i < len(Items); i++ {
		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(Items[i]))
	}

	if prevButt && nextButt {
		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущая страница",
				FormPreviousPageButtonData(ButtonDataHead, itemsPerPage, firstItemOnPageNumberParameterName).String()),
			tgbotapi.NewInlineKeyboardButtonData("Следующая страница",
				FormNextPageButtonData(ButtonDataHead, itemsPerPage, firstItemOnPageNumberParameterName).String())))
	} else if prevButt {
		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Предыдущая страница",
				FormPreviousPageButtonData(ButtonDataHead, itemsPerPage, firstItemOnPageNumberParameterName).String())))
	} else if nextButt {
		menu.InlineKeyboard = append(menu.InlineKeyboard, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Следующая страница",
				FormNextPageButtonData(ButtonDataHead, itemsPerPage, firstItemOnPageNumberParameterName).String())))
	}
}
