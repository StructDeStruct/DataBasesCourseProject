package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Clip struct {
	Id int
	Uuid string
	Name string
	Description string
	UploadTime time.Time
	OwnerId int
	PublisherId int
	Duration int
}

type User struct {
	Id int
	TelegramId int
	Name string
	Info string
}

type Group struct {
	Id int
	Name string
	Description string
	CreatorId int
}

type GroupMember struct {
	Id int
	GroupId int
	UserId int
}

type Owner struct {
	Id int
	IsGroup bool

}

type Tag struct {
	Id int
	Name string
	ClipId int
}

type Comment struct {
	Id              int
	ParentCommentId sql.NullInt64
	CommentText     string
	ClipId          int
	PublisherId     sql.NullInt64
	Time            time.Time
}

type Subscription struct {
	Id           int
	SubscriberId int
	TargetUserId int

}

func GetClipsPageByOwnerId(sqlDb *sql.DB, ownerId int, firstClipOnPageNumber int) ([]*Clip, error) {
	sqlQuery := "SELECT c.clipid, c.uuid, c.name, c.description, c.publisherid, c.duration " +
		"FROM ( " +
		"SELECT clipid " +
		"FROM ( " +
		"SELECT *, row_number() OVER () AS rnum " +
		"FROM ( " +
		"SELECT clipid, name " +
		"FROM clips AS c " +
		"WHERE ownerid = $1 " +
		"ORDER BY name) AS l) AS l " +
		"WHERE l.rnum >= $2 " +
		"LIMIT $3) AS l JOIN clips AS c ON (c.clipid = l.clipid);"

	rows, err := sqlDb.Query(sqlQuery, ownerId, firstClipOnPageNumber, ClipsPerPage+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clips []*Clip

	for rows.Next() {
		clip := new(Clip)

		err = rows.Scan(&clip.Id, &clip.Uuid, &clip.Name, &clip.Description, &clip.PublisherId, &clip.Duration)
		if err != nil {
			return nil, err
		}

		clips = append(clips, clip)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clips, nil
}

func GetClipsPageByTag(sqlDb *sql.DB, ownerId int, tagName string, firstClipOnPageNumber int) ([]*Clip, error) {
	sqlQuery := "SELECT c.clipid, c.uuid, c.name, c.description, c.publisherid, c.duration " +
		"FROM ( " +
		"SELECT l.clipid " +
		"FROM ( " +
		"SELECT l.clipid, row_number() OVER () AS rnum " +
		"FROM ( " +
		"SELECT c.clipid " +
		"FROM ( " +
		"SELECT clipid, name " +
		"FROM clips " +
		"WHERE ownerid = $1) AS c JOIN tags AS t ON (t.clipid = c.clipid) " +
		"WHERE t.name = $2 " +
		"ORDER BY c.name) AS l) AS l " +
		"WHERE l.rnum >= $3 " +
		"LIMIT $4) AS l JOIN clips AS c ON (c.clipid = l.clipid);"

	rows, err := sqlDb.Query(sqlQuery, ownerId, tagName, firstClipOnPageNumber, ClipsPerPage+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clips []*Clip

	for rows.Next() {
		clip := new(Clip)

		err = rows.Scan(&clip.Id, &clip.Uuid, &clip.Name, &clip.Description, &clip.PublisherId, &clip.Duration)
		if err != nil {
			return nil, err
		}

		clips = append(clips, clip)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return clips, nil
}


func GetOwnerName(sqlDb *sql.DB, ownerId int) (string, bool, error) {
	var userName, groupName sql.NullString
	owner := new(Owner)

	sqlQuery := "SELECT o.is_group, u.username, g.name " +
		"FROM ( " +
		"SELECT * " +
		"FROM owners WHERE ownerid = $1) AS o " +
		"LEFT JOIN groups AS g ON (o.ownerid = g.groupid) " +
		"LEFT JOIN users AS u ON (o.ownerid = u.userid);"

	err := sqlDb.QueryRow(sqlQuery, ownerId).Scan(&owner.IsGroup, &userName, &groupName)
	if err != nil {
		return "", false, err
	}

	var ownerName string
	if owner.IsGroup {
		ownerName = groupName.String
	} else {
		ownerName = userName.String
	}

	return ownerName, owner.IsGroup, nil
}

func GetGroupsPage(sqlDb *sql.DB, userId, firstGroupOnPageNumber int) ([]*Group, []string, error) {
	sqlQuery := "SELECT g.groupid, g.name, g.description, g.creatorid, u.username " +
		"FROM ( " +
		"SELECT g.groupid, g.name, g.description, g.creatorid " +
		"FROM ( " +
		"SELECT l.groupid " +
		"FROM ( " +
		"SELECT l.groupid, row_number() OVER () AS rnum " +
		"FROM ( " +
		"SELECT g.groupid " +
		"FROM ( " +
		"SELECT groupid " +
		"FROM group_members " +
		"WHERE userid = $1) AS gm JOIN groups AS g ON (gm.groupid = g.groupid) " +
		"ORDER BY g.name) AS l) AS l " +
		"WHERE l.rnum >= $2 " +
		"LIMIT $3) AS l JOIN groups AS g ON (g.groupid = l.groupid)) AS g " +
		"JOIN users AS u ON (u.userid = g.creatorid);"

	rows, err := sqlDb.Query(sqlQuery, userId, firstGroupOnPageNumber, GroupsPerPage+1)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var groups []*Group
	var creatorNames []string
	for rows.Next() {
		group := new(Group)
		var userName string

		err = rows.Scan(&group.Id, &group.Name, &group.Description, &group.CreatorId, &userName)
		if err != nil {
			return nil, nil, err
		}

		creatorName := userName

		creatorNames = append(creatorNames, creatorName)
		groups = append(groups, group)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return groups, creatorNames, nil
}

func GetGroupMembersPage(sqlDb *sql.DB, groupId, firstGroupMemberOnPageNumber int) ([]*User, error) {
	sqlQuery := "SELECT l.userid, l.username " +
		"FROM ( " +
		"SELECT *, row_number() OVER () AS rnum " +
		"FROM ( " +
		"SELECT u.userid, u.username " +
		"FROM ( " +
		"SELECT userid " +
		"FROM group_members " +
		"WHERE groupid = $1) AS gm JOIN users AS u ON (gm.userid = u.userid) " +
		"ORDER BY u.username) AS l) AS l " +
		"WHERE l.rnum >= $2 " +
		"LIMIT $3;"

	rows, err := sqlDb.Query(sqlQuery, groupId, firstGroupMemberOnPageNumber, GroupMembersPerPage+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := new(User)

		err = rows.Scan(&user.Id, &user.Name)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetSubscriptionsPage(sqlDb *sql.DB, userId, firstSubscriptionOnPageNumber int) ([]*User, error) {
	sqlQuery := "SELECT l.userid, l.username " +
		"FROM ( " +
		"SELECT *, row_number() OVER () AS rnum " +
		"FROM ( " +
		"SELECT u.userid, u.username " +
		"FROM ( " +
		"SELECT target_userid " +
		"FROM subscriptions " +
		"WHERE subscriberid = $1) AS s JOIN users AS u ON (s.target_userid = u.userid) " +
		"ORDER BY u.username) AS l) AS l " +
		"WHERE l.rnum >= $2 " +
		"LIMIT $3;"

	rows, err := sqlDb.Query(sqlQuery, userId, firstSubscriptionOnPageNumber, UsersPerPage+1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := new(User)

		err = rows.Scan(&user.Id, &user.Name)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func GetCommentsPage(sqlDb *sql.DB, clipId, parentCommentId int, hasParentComment bool, firstCommentOnPageNumber int) (*Comment, *User, []*Comment, []*User, error) {
	var rows *sql.Rows
	var err error
	if hasParentComment {
		sqlQuery := "SELECT c.commentid, c.clipid, c.comment_text, c.publisherid, c.time, u.username " +
			"FROM ( " +
			"SELECT commentid, clipid, comment_text, publisherid, time " +
			"FROM comments " +
			"WHERE commentid = $1) AS c JOIN users AS u ON (c.publisherid = u.userid) " +
			"UNION ALL " +
			"SELECT c.commentid, c.clipid, c.comment_text, c.publisherid, c.time, u.username " +
			"FROM ( " +
			"SELECT c.commentid, c.clipid, c.comment_text, c.publisherid, c.time " +
			"FROM ( " +
			"SELECT l.commentid " +
			"FROM ( " +
			"SELECT *, row_number() OVER () AS rnum " +
			"FROM ( " +
			"SELECT commentid, time " +
			"FROM comments " +
			"WHERE clipid = $2 AND parent_commentid IS NOT NULL AND parent_commentid = $1 " +
			"ORDER BY time) AS l) AS l " +
			"WHERE l.rnum >= $3 " +
			"LIMIT $4) AS l JOIN comments AS c ON (l.commentid = c.commentid)) AS c " +
			"JOIN users AS u ON (c.publisherid = u.userid);"

		//log.Println("\t\t\t\t ParrentCommentid " + strconv.Itoa(parentCommentId))
		rows, err = sqlDb.Query(sqlQuery, parentCommentId, clipId, firstCommentOnPageNumber, CommentsPerPage+1)
	} else {
		sqlQuery := "SELECT c.commentid, c.clipid, c.comment_text, c.publisherid, c.time, u.username " +
			"FROM ( " +
			"SELECT c.commentid, c.clipid, c.comment_text, c.publisherid, c.time " +
			"FROM ( " +
			"SELECT l.commentid " +
			"FROM ( " +
			"SELECT *, row_number() OVER () AS rnum " +
			"FROM ( " +
			"SELECT commentid, time " +
			"FROM comments " +
			"WHERE clipid = $1 AND parent_commentid IS NULL " +
			"ORDER BY time) AS l) AS l " +
			"WHERE l.rnum >= $2 " +
			"LIMIT $3) AS l JOIN comments AS c ON (l.commentid = c.commentid)) AS c " +
			"JOIN users AS u ON (c.publisherid = u.userid);"

		rows, err = sqlDb.Query(sqlQuery, clipId, firstCommentOnPageNumber, CommentsPerPage+1)
	}
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer rows.Close()

	var parentComment *Comment
	var parentCommentUser *User
	if hasParentComment {
		if rows.Next() {
			parentComment = new(Comment)
			parentCommentUser = new(User)

			err = rows.Scan(&parentComment.Id, &parentComment.ClipId, &parentComment.CommentText,
				&parentComment.PublisherId, &parentComment.Time,
				&parentCommentUser.Name)

			if err != nil {
				return nil, nil, nil, nil, err
			}

			log.Println("\t\t\t\t ParrentCommentid " + strconv.Itoa(parentComment.Id))
		}
	}


	var comments []*Comment
	var users []*User
	for rows.Next() {
		comment := new(Comment)
		user := new(User)

		err = rows.Scan(&comment.Id, &comment.ClipId, &comment.CommentText,
			&comment.PublisherId, &comment.Time,
			&user.Name)

		if err != nil {
			return nil, nil, nil, nil, err
		}


		comments = append(comments, comment)
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, nil, nil, err
	}

	return parentComment, parentCommentUser, comments, users, nil
}

func GetTagsPage(sqlDb *sql.DB, ownerId, firstTagOnPageNumber int) ([]*Tag, []int, error) {
	sqlQuery := "SELECT l.name, l.clips_counter " +
		"FROM ( " +
		"SELECT *, row_number() OVER () AS rnum " +
		"FROM ( " +
		"SELECT t.name, count(c.clipid) AS clips_counter " +
		"FROM ( " +
		"SELECT clipid " +
		"FROM clips " +
		"WHERE ownerid = $1) AS c JOIN tags AS t ON (c.clipid = t.clipid) " +
		"GROUP BY t.name " +
		"ORDER BY t.name) AS l) AS l " +
		"WHERE l.rnum >= $2 " +
		"LIMIT $3;"

	rows, err := sqlDb.Query(sqlQuery, ownerId, firstTagOnPageNumber, TagsPerPage+1)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var tags []*Tag
	var clipNumbers []int
	for rows.Next() {
		tag := new(Tag)
		var clipNumber int

		err = rows.Scan(&tag.Name, &clipNumber)
		if err != nil {
			return nil, nil, err
		}

		tags = append(tags, tag)
		clipNumbers = append(clipNumbers, clipNumber)
	}
	if err = rows.Err(); err != nil {
		return nil, nil, err
	}

	return tags, clipNumbers, nil
}

func GetUserDataByUserId(sqlDb *sql.DB, userId, requestTelegramId int) (*User, bool, bool, error) {
	sqlQuery := "SELECT u.userid, u.telegramid, u.username, u.info, request_u.userid, s.subscriptionid " +
		"FROM ( " +
		"SELECT * " +
		"FROM users " +
		"WHERE userid = $1) AS u " +
		"JOIN (SELECT userid FROM users WHERE telegramid = $2) AS request_u ON TRUE " +
		"LEFT JOIN subscriptions AS s ON (u.userid = s.target_userid AND request_u.userid = s.subscriberid);"

	user := new(User)
	var subscriptionId sql.NullInt64
	var requestUserId int

	err := sqlDb.QueryRow(sqlQuery, userId, requestTelegramId).Scan(&user.Id, &user.TelegramId,
		&user.Name, &user.Info, &requestUserId, &subscriptionId)
	if err != nil {
		return nil, false, false, err
	}

	equalsRequestUser := userId == requestUserId

	isRequestUserSubscribed := false
	if subscriptionId.Valid {
		isRequestUserSubscribed = true
	}

	return user, isRequestUserSubscribed, equalsRequestUser, nil
}

func GetUserDataByUserName(sqlDb *sql.DB, userName string, requestTelegramId int) (*User, bool, bool, error) {
	sqlQuery := "SELECT u.userid, u.telegramid, u.username, u.info, request_u.userid, s.subscriptionid " +
		"FROM ( " +
		"SELECT * " +
		"FROM users " +
		"WHERE username = $1) AS u " +
		"JOIN (SELECT userid FROM users WHERE telegramid = $2) AS request_u ON TRUE " +
		"LEFT JOIN subscriptions AS s ON (u.userid = s.target_userid AND request_u.userid = s.subscriberid);"

	user := new(User)
	var subscriptionId sql.NullInt64
	var requestUserId int

	err := sqlDb.QueryRow(sqlQuery, userName, requestTelegramId).Scan(&user.Id, &user.TelegramId,
		&user.Name, &user.Info, &requestUserId, &subscriptionId)
	if err != nil {
		return nil, false, false, err
	}

	equalsRequestUser := user.Id == requestUserId

	isRequestUserSubscribed := false
	if subscriptionId.Valid {
		isRequestUserSubscribed = true
	}

	return user, isRequestUserSubscribed, equalsRequestUser, nil
}

func GetUserByUserName(sqlDb *sql.DB, userName string) (*User, error) {
	sqlQuery := "SELECT u.userid, u.telegramid, u.username, u.info " +
		"FROM users AS u " +
		"WHERE username = $1;"

	user := new(User)

	err := sqlDb.QueryRow(sqlQuery, userName).Scan(&user.Id, &user.TelegramId,
		&user.Name, &user.Info)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserDataByTelegramId(sqlDb *sql.DB, userTelegramId int) (*User, error) {
	sqlQuery := "SELECT u.userid, u.telegramid, u.username, u.info " +
		"FROM users AS u " +
		"WHERE telegramid = $1;"

	user := new(User)

	err := sqlDb.QueryRow(sqlQuery, userTelegramId).Scan(&user.Id, &user.TelegramId,
		&user.Name, &user.Info)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetGroupData(sqlDb *sql.DB, groupId, requestTelegramId int) (*Group, string, bool, bool, error) {
	sqlQuery := "SELECT g.groupid, g.name, g.description, g.creatorid, request_u.userid, gm.group_memberid, u.username " +
		"FROM ( " +
		"SELECT * " +
		"FROM groups " +
		"WHERE groupid = $1) AS g " +
		"JOIN (SELECT userid FROM users WHERE telegramid = $2) AS request_u ON TRUE " +
		"LEFT JOIN group_members AS gm ON (g.groupid = gm.groupid AND request_u.userid = gm.userid) " +
		"JOIN users AS u ON (g.creatorid = u.userid);"

	group := new(Group)
	var groupMemberId sql.NullInt64
	var requestUserId int
	user := new(User)

	err := sqlDb.QueryRow(sqlQuery, groupId, requestTelegramId).Scan(&group.Id, &group.Name,
		&group.Description, &group.CreatorId, &requestUserId, &groupMemberId,
		&user.Name)
	if err != nil {
		return nil, "", false, false, err
	}

	isRequestUserGroupCreator := group.CreatorId == requestUserId

	isRequestUserMember := false
	if groupMemberId.Valid {
		isRequestUserMember = true
	}

	return group, user.Name, isRequestUserMember, isRequestUserGroupCreator, nil
}

func GetGroupByCreatorAndName(sqlDb *sql.DB, creatorId int, groupName string) (*Group, error) {
	sqlQuery := "SELECT g.groupid, g.name, g.description, g.creatorid " +
		"FROM groups AS g " +
		"WHERE g.creatorid = $1 AND g.name = $2;"

	group := new(Group)

	err := sqlDb.QueryRow(sqlQuery, creatorId, groupName).Scan(&group.Id, &group.Name,
		&group.Description, &group.CreatorId)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func GetClipDataByClipId(sqlDb *sql.DB, clipId, requestTelegramId int) (*Clip, bool, error) {
	sqlQuery := "SELECT c.clipid, c.uuid, c.name, c.description, c.upload_time, c.ownerid, c.publisherid, c.duration, u.telegramid " +
		"FROM ( " +
		"SELECT clipid, uuid, name, description, upload_time, ownerid, publisherid, duration " +
		"FROM clips " +
		"WHERE clipid = $1) AS c JOIN users AS u ON (u.userid = c.publisherid);"

	clip := new(Clip)
	var publisherTelegramId int
	err := sqlDb.QueryRow(sqlQuery, clipId).Scan(&clip.Id, &clip.Uuid, &clip.Name,
		&clip.Description, &clip.UploadTime, &clip.OwnerId, &clip.PublisherId, &clip.Duration, &publisherTelegramId)
	if err != nil {
		return nil, false, err
	}

	isRequestUserClipPublisher := requestTelegramId == publisherTelegramId
	return clip, isRequestUserClipPublisher, nil
}

func GetClipDataByUuid(sqlDb *sql.DB, uuid string, requestTelegramId int) (*Clip, bool, error) {
	sqlQuery := "SELECT c.clipid, c.uuid, c.name, c.description, c.upload_time, c.ownerid, c.publisherid, c.duration, u.telegramid " +
		"FROM ( " +
		"SELECT clipid, uuid, name, description, upload_time, ownerid, publisherid, duration " +
		"FROM clips " +
		"WHERE uuid = $1) AS c JOIN users AS u ON (u.userid = c.publisherid);"

	clip := new(Clip)
	var publisherTelegramId int
	err := sqlDb.QueryRow(sqlQuery, uuid).Scan(&clip.Id, &clip.Uuid, &clip.Name,
		&clip.Description, &clip.UploadTime, &clip.OwnerId, &clip.PublisherId, &clip.Duration, &publisherTelegramId)
	if err != nil {
		return nil, false, err
	}

	isRequestUserClipPublisher := requestTelegramId == publisherTelegramId
	return clip, isRequestUserClipPublisher, nil
}

func GetClipUUID(sqlDb *sql.DB, clipId int) (string, error) {
	sqlQuery := "SELECT uuid FROM clips WHERE clipid = $1;"

	var uuid string
	err := sqlDb.QueryRow(sqlQuery, clipId).Scan(&uuid)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func GetClipTags(sqlDb *sql.DB, clipId int) ([]*Tag, error) {
	sqlQuery := "SELECT tagid, name, clipid " +
		"FROM tags " +
		"WHERE clipid = $1;"

	rows, err := sqlDb.Query(sqlQuery, clipId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*Tag
	for rows.Next() {
		tag := new(Tag)

		err = rows.Scan(&tag.Id, &tag.Name, &tag.ClipId)
		if err != nil {
			return nil, err
		}

		tags = append(tags, tag)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

func InsertOwner(sqlDb *sql.DB, isGroup bool) (int, error) {
	sqlQuery := "INSERT INTO owners (is_group) VALUES ($1) RETURNING ownerid;"

	var ownerId int

	err := sqlDb.QueryRow(sqlQuery, isGroup).Scan(&ownerId)
	return ownerId, err
}

func InsertUser(sqlDb *sql.DB, user *User, ownerId int) error {
	sqlQuery := "INSERT INTO users (userid, telegramid, username, info) VALUES ($1, $2, $3, $4);"

	_, err := sqlDb.Exec(sqlQuery, ownerId, user.TelegramId, user.Name, user.Info)
	return err
}

func InsertSubscription(sqlDb *sql.DB, userId, targetUserId int) error {
	sqlQuery := "INSERT INTO subscriptions (subscriberid, target_userid) VALUES ($1, $2);"

	_, err := sqlDb.Exec(sqlQuery, userId, targetUserId)

	return err
}

func InsertGroupMember(sqlDb *sql.DB, groupId, userId int) error {
	sqlQuery := "INSERT INTO group_members (groupid, userid) VALUES ($1, $2);"

	_, err := sqlDb.Exec(sqlQuery, groupId, userId)

	return err
}

func InsertClip(sqlDb *sql.DB, clip *Clip) (int, error) {
	sqlQuery := "INSERT INTO clips (uuid, name, description, upload_time, ownerid, publisherid, duration) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING clipid;"

	var clipId int
	err := sqlDb.QueryRow(sqlQuery, clip.Uuid, clip.Name, clip.Description, clip.UploadTime, clip.OwnerId, clip.PublisherId, clip.Duration).Scan(&clipId)

	return clipId, err
}

func InsertTags(sqlDb *sql.DB, clipid int, tags []*Tag) error {
	var valueStrings []string
	var valueArgs []interface{}
	for i := range tags {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", 2 * i + 1, 2 * i + 2))
		valueArgs = append(valueArgs, clipid)
		valueArgs = append(valueArgs, tags[i].Name)
	}

	sqlQuery := fmt.Sprintf("INSERT INTO tags (clipid, name) VALUES %s;", strings.Join(valueStrings, ","))

	_, err := sqlDb.Exec(sqlQuery, valueArgs...)

	return err
}

func InsertGroup(sqlDb *sql.DB, group *Group, ownerId int) (int, error) {
	sqlQuery := "INSERT INTO groups (groupid, name, description, creatorid) VALUES ($1, $2, $3, $4) RETURNING groupid;"

	var groupId int

	err := sqlDb.QueryRow(sqlQuery, ownerId, group.Name, group.Description, group.CreatorId).Scan(&groupId)

	return groupId, err
}

func InsertComment(sqlDb *sql.DB, comment *Comment) error {
	var err error
	if comment.ParentCommentId.Valid {
		sqlQuery := "INSERT INTO comments (parent_commentid, comment_text, clipid, publisherid, time) VALUES ($1, $2, $3, $4, $5);"

		log.Println("\t\t\t\t\t\t trying to use parentCOmmentId")
		parentCommentId, _ := comment.ParentCommentId.Value()
		_, err = sqlDb.Exec(sqlQuery,parentCommentId, comment.CommentText, comment.ClipId, comment.PublisherId, comment.Time)
	} else {
		sqlQuery := "INSERT INTO comments (comment_text, clipid, publisherid, time) VALUES ($1, $2, $3, $4);"

		_, err = sqlDb.Exec(sqlQuery, comment.CommentText, comment.ClipId, comment.PublisherId, comment.Time)
	}

	return err
}


func DeleteOwner(sqlDb *sql.DB, ownerId int) error {
	sqlQuery := "DELETE FROM owners WHERE ownerid = $1;"

	_, err := sqlDb.Exec(sqlQuery, ownerId)

	return err
}

func DeleteSubscription(sqlDb *sql.DB, userId, targetUserId int) error {
	sqlQuery := "DELETE FROM subscriptions WHERE subscriberid = $1 AND target_userid = $2;"

	_, err := sqlDb.Exec(sqlQuery, userId, targetUserId)

	return err
}

func DeleteGroupMember(sqlDb *sql.DB, groupId, userId int) error {
	sqlQuery := "DELETE FROM group_members WHERE groupid = $1 AND userid = $2;"

	_, err := sqlDb.Exec(sqlQuery, groupId, userId)

	return err
}

func DeleteGroup(sqlDb *sql.DB, groupId int) error {
	sqlQuery := "DELETE FROM groups WHERE groupid = $1;"

	_, err := sqlDb.Exec(sqlQuery, groupId)

	return err
}

func DeleteClip(sqlDb *sql.DB, clipId int) error {
	sqlQuery := "DELETE FROM clips WHERE clipid = $1;"

	_, err := sqlDb.Exec(sqlQuery, clipId)

	return err
}


func UpdateUserDescription(sqlDb *sql.DB, requestTelegramId int, description string) error {
	sqlQuery := "UPDATE users SET info = $1 WHERE telegramid = $2;"

	_, err := sqlDb.Exec(sqlQuery, description, requestTelegramId)

	return err
}

func UpdateClipDescription(sqlDb *sql.DB, clipId int, description string) error {
	sqlQuery := "UPDATE clips SET description = $1 WHERE clipid = $2;"

	_, err := sqlDb.Exec(sqlQuery, description, clipId)

	return err
}

func UpdateGroupDescription(sqlDb *sql.DB, groupId int, description string) error {
	sqlQuery := "UPDATE groups SET description = $1 WHERE groupid = $2;"

	_, err := sqlDb.Exec(sqlQuery, description, groupId)

	return err
}

