package main

import "time"

type State int

const (
	DefaultState State = iota
	AddComment

	AddClip
	AddClipTags
	AddClipDescrpition
	AddClipName
	SendClipVideoFile

	CreateNewGroup
	AddGroupName
	AddGroupDescription

	EditGroupProfile
	EditGroupDescription

	EditUserDescription

	EditClipDescription

	SearchClipByUuid
	SearchUserByUserName

	AddUserName
)

type UserInfo struct {
	LastSeen  time.Time
	UserState State
	TempUser *User
	TempClip *Clip
	TempGroup *Group
	TempTags []*Tag
	TempComment *Comment

	StackedData string
}

type Users map[int]*UserInfo

func InitUsers() Users {
	return make(Users)
}

func (users Users) Delete(userid int) {
	delete(users, userid)
}

func (users Users) User(userid int) *UserInfo {
	user, ok := users[userid]
	if !ok {
		user = &UserInfo{UserState: DefaultState}
		user.LastSeen = time.Now()
		users[userid] = user
	}
	return user
}