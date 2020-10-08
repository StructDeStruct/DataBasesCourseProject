DROP TRIGGER IF EXISTS validate_group_members ON group_members;
DROP TRIGGER IF EXISTS validate_groups ON groups;
DROP TRIGGER IF EXISTS validate_clips ON clips;
DROP TRIGGER IF EXISTS validate_subscriptions ON subscriptions;
DROP TRIGGER IF EXISTS validate_tags ON tags;
DROP TRIGGER IF EXISTS validate_users ON users;

DROP FUNCTION IF EXISTS process_validate_group_members;
DROP FUNCTION IF EXISTS process_validate_groups;
DROP FUNCTION IF EXISTS process_validate_clips;
DROP FUNCTION IF EXISTS process_validate_subscriptions;
DROP FUNCTION IF EXISTS process_validate_tags;
DROP FUNCTION IF EXISTS process_validate_users;

DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS clips;
DROP TABLE IF EXISTS group_members;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS owners;


CREATE TABLE IF NOT EXISTS owners (
    ownerid int PRIMARY KEY GENERATED ALWAYS AS IDENTITY NOT NULL,
    is_group bool NOT NULL
);

CREATE TABLE IF NOT EXISTS users(
    userid int PRIMARY KEY NOT NULL,
    telegramid int UNIQUE NOT NULL,
    username text NOT NULL,
    info text NOT NULL,
    CONSTRAINT FK_Users_Userid FOREIGN KEY (userid) REFERENCES owners(ownerid) ON DELETE CASCADE
);

CREATE TABLE subscriptions (
    subscriptionid int PRIMARY KEY GENERATED ALWAYS AS IDENTITY NOT NULL,
    subscriberid int NOT NULL,
    target_userid int NOT NULL,
    CONSTRAINT FK_Subscribers_Subscriberid FOREIGN KEY (subscriberid) REFERENCES users(userid) ON DELETE CASCADE,
    CONSTRAINT FK_Subscribers_TargetUserid FOREIGN KEY (target_userid) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE TABLE groups (
    groupid int PRIMARY KEY NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    creatorid int NOT NULL,
    CONSTRAINT FK_Groups_Groupid FOREIGN KEY (groupid) REFERENCES owners(ownerid) ON DELETE CASCADE,
    CONSTRAINT FK_Groups_Creatorid FOREIGN KEY (creatorid) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE TABLE group_members (
    group_memberid int PRIMARY KEY GENERATED ALWAYS AS IDENTITY NOT NULL,
    groupid int NOT NULL,
    userid int NOT NULL,
    CONSTRAINT FK_GroupMembers_Groupid FOREIGN KEY (groupid) REFERENCES groups(groupid) ON DELETE CASCADE,
    CONSTRAINT FK_GroupMembers_Userid FOREIGN KEY (userid) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS clips (
    clipid int PRIMARY KEY GENERATED ALWAYS AS IDENTITY NOT NULL,
    upload_time timestamp DEFAULT current_timestamp NOT NULL,
    name text NOT NULL,
    description text NOT NULL,
    uuid text NOT NULL,
    ownerid int NOT NULL,
    publisherid int NOT NULL,
    duration int CHECK (duration > 0) NOT NULL,
    CONSTRAINT FK_Clips_Ownerid FOREIGN KEY (ownerid) REFERENCES owners(ownerid) ON DELETE CASCADE,
    CONSTRAINT FK_Clips_Publisherid FOREIGN KEY (publisherid) REFERENCES users(userid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tags(
    tagid int PRIMARY KEY GENERATED ALWAYS AS IDENTITY NOT NULL,
    clipid int NOT NULL,
    name text NOT NULL,
    CONSTRAINT FK_Tags_Clipid FOREIGN KEY (clipid) REFERENCES clips(clipid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comments (
    commentid int PRIMARY KEY GENERATED ALWAYS AS IDENTITY NOT NULL,
    parent_commentid int NULL,
    comment_text text NOT NULL,
    clipid int NOT NULL,
    publisherid int NULL,
    time timestamp DEFAULT now() NOT NULL,
    CONSTRAINT FK_Comments_Clipid FOREIGN KEY (clipid) REFERENCES clips(clipid),
    CONSTRAINT FK_Comments_ParentCommentid FOREIGN KEY (parent_commentid) REFERENCES comments(commentid),
    CONSTRAINT FK_Comments_Publisherid FOREIGN KEY (publisherid) REFERENCES users(userid)
);

CREATE OR REPLACE FUNCTION process_validate_group_members() RETURNS TRIGGER AS $validate_group_members$
    BEGIN
        -- edit comparison here
        IF EXISTS(SELECT * FROM group_members WHERE groupid = NEW.groupid AND userid = NEW.userid) THEN
            RAISE EXCEPTION 'subscription already exists in database' USING ERRCODE = '37846';
        end if;

        RETURN NEW;
    END
$validate_group_members$ LANGUAGE plpgsql;

CREATE TRIGGER validate_group_members BEFORE INSERT On group_members
    FOR EACH ROW EXECUTE PROCEDURE process_validate_group_members();


CREATE OR REPLACE FUNCTION process_validate_groups() RETURNS TRIGGER AS $validate_groups$
    BEGIN
        -- edit comparison here
        IF EXISTS(SELECT * FROM groups WHERE creatorid = NEW.creatorid AND name = NEW.name) THEN
            RAISE EXCEPTION 'group already exists in database' USING ERRCODE = '37846';
        end if;

        RETURN NEW;
    END
$validate_groups$ LANGUAGE plpgsql;

CREATE TRIGGER validate_groups BEFORE INSERT ON groups
    FOR EACH ROW EXECUTE PROCEDURE process_validate_groups();


CREATE OR REPLACE FUNCTION process_validate_clips() RETURNS TRIGGER AS $validate_clips$
    BEGIN
        -- edit comparison here
        IF EXISTS(SELECT * FROM clips WHERE ownerid = NEW.ownerid AND name = NEW.name) THEN
            RAISE EXCEPTION 'clip already exists in database' USING ERRCODE = '37846';
        end if;

        RETURN NEW;
    END
$validate_clips$ LANGUAGE plpgsql;

CREATE TRIGGER validate_clips BEFORE INSERT ON clips
    FOR EACH ROW EXECUTE PROCEDURE process_validate_clips();


CREATE OR REPLACE FUNCTION process_validate_subscriptions() RETURNS TRIGGER AS $validate_subscription$
    BEGIN
        -- edit comparison here
        IF EXISTS(SELECT * FROM subscriptions WHERE subscriberid = NEW.subscriberid AND target_userid = NEW.target_userid) THEN
            RAISE EXCEPTION 'subscription already exists in database' USING ERRCODE = '37846';
        end if;

        RETURN NEW;
    END
$validate_subscription$ LANGUAGE plpgsql;

CREATE TRIGGER validate_subscriptions BEFORE INSERT On subscriptions
    FOR EACH ROW EXECUTE PROCEDURE process_validate_subscriptions();


CREATE OR REPLACE FUNCTION process_validate_tags() RETURNS TRIGGER AS $validate_tags$
    BEGIN
        -- edit comparison here
        IF EXISTS(SELECT * FROM tags WHERE clipid = NEW.clipid AND name = NEW.name) THEN
            RAISE EXCEPTION 'subscription already exists in database' USING ERRCODE = '37846';
        end if;

        RETURN NEW;
    END
$validate_tags$ LANGUAGE plpgsql;

CREATE TRIGGER validate_tags BEFORE INSERT On tags
    FOR EACH ROW EXECUTE PROCEDURE process_validate_tags();


CREATE OR REPLACE FUNCTION process_validate_users() RETURNS TRIGGER AS $validate_users$
    BEGIN
        -- edit comparison here
        IF EXISTS(SELECT * FROM users WHERE telegramid = NEW.telegramid OR username = NEW.username) THEN
            RAISE EXCEPTION 'user already exists in database' USING ERRCODE = '37846';
        end if;
        RETURN NEW;
    END
$validate_users$ LANGUAGE plpgsql;

CREATE TRIGGER validate_users BEFORE INSERT On users
    FOR EACH ROW EXECUTE PROCEDURE process_validate_users();
