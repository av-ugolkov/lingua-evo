package handler

const (
	CurrentLanguage    = "/current_language" //get
	AvailableLanguages = "/languages"        //get

	SignIn     = "/auth/sign_in"   //post
	SignUp     = "/auth/sign_up"   //post
	Refresh    = "/auth/refresh"   //get
	SignOut    = "/auth/sign_out"  //get
	SendCode   = "/auth/send_code" //post
	GoogleAuth = "/auth/google"    //get, post

	UserByID = "/user/id" //get
	Users    = "/users"   //get

	AccountSettingsAccount         = "/account/settings/account"           //get
	AccountSettingsPersonalInfo    = "/account/settings/personal_info"     //get
	AccountSettingsEmailNotif      = "/account/settings/email_notif"       //get
	AccountSettingsUpdatePswCode   = "/account/settings/update_psw_code"   //post
	AccountSettingsUpdatePsw       = "/account/settings/update_psw"        //post
	AccountSettingsUpdateEmailCode = "/account/settings/update_email_code" //post
	AccountSettingsUpdateEmail     = "/account/settings/update_email"      //post
	AccountSettingsUpdateNickname  = "/account/settings/update_nickname"   //post

	Dictionary        = "/dictionary"
	DictionaryWord    = "/dictionary/word"               //get post
	GetRandomWord     = "/dictionary/word/random"        //get
	WordPronunciation = "/dictionary/word/pronunciation" //get

	UserVocabularies = "/account/vocabularies" //get

	Vocabulary              = "/vocabulary"               //post, put, delete
	Vocabularies            = "/vocabularies"             //get
	VocabulariesRecommended = "/vocabularies/recommended" //get
	VocabularyInfo          = "/vocabulary/info"          //get
	VocabularyCopy          = "/vocabulary/copy"          //get
	VocabularyAccessForUser = "/vocabulary/access/user"   //get post delete patch
	VocabulariesByUser      = "/vocabularies/user"        //get

	VocabularyWord              = "/vocabulary/word"               //get post delete
	VocabularyWordText          = "/vocabulary/word/text"          //post
	VocabularyWordPronunciation = "/vocabulary/word/pronunciation" //post
	VocabularyWordDefinition    = "/vocabulary/word/definition"    //post
	VocabularyWordTranslates    = "/vocabulary/word/translates"    //post
	VocabularyWordExamples      = "/vocabulary/word/examples"      //post
	VocabularyWords             = "/vocabulary/words"              //get

	VocabularyTags = "/vocabulary/tag" //get

	Accesses = "/accesses" //get

	CheckSubscriber = "/subscriber/check" //get
	Subscribe       = "/user/subscribe"   //post
	Unsubscribe     = "/user/unsubscribe" //post

	NotificationVocab = "/notifications/vocabulary"

	SupportRequest = "/support/request"

	Events      = "/events"
	CountEvents = "/events/count"
	MarkWatched = "/event/watched"

	ReviseGame = "/game/revise"
)
