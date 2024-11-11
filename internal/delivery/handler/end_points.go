package handler

const (
	CurrentLanguage    = "/current_language" //get
	AvailableLanguages = "/languages"        //get

	SignIn   = "/auth/sign_in"   //post
	Refresh  = "/auth/refresh"   //get
	SignOut  = "/auth/sign_out"  //get
	SendCode = "/auth/send_code" //post

	SignUp   = "/user/sign_up" //post
	UserByID = "/user/id"      //get
	Users    = "/users"        //get

	UserUpdatePswCode   = "/account/settings/update_psw_code"   //post
	UserUpdatePsw       = "/account/settings/update_psw"        //post
	UserUpdateEmailCode = "/account/settings/update_email_code" //post
	UserUpdateEmail     = "/account/settings/update_email"

	DictionaryWord = "/dictionary/word"        //get post
	GetRandomWord  = "/dictionary/word/random" //get

	UserVocabularies = "/account/vocabularies" //get

	Vocabulary              = "/vocabulary"               //post, put, delete
	Vocabularies            = "/vocabularies"             //get
	VocabulariesRecommended = "/vocabularies/recommended" //get
	VocabularyInfo          = "/vocabulary/info"          //get
	VocabularyCopy          = "/vocabulary/copy"          //get
	VocabularyAccessForUser = "/vocabulary/access/user"   //get post delete patch
	VocabulariesByUser      = "/vocabularies/user"        //get

	VocabularyWord       = "/vocabulary/word"               //get post delete
	VocabularyWordUpdate = "/vocabulary/word/update"        //post
	VocabularyWords      = "/vocabulary/words"              //get
	WordPronunciation    = "/vocabulary/word/pronunciation" //get

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
)
