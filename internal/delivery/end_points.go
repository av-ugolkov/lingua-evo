package delivery

const (
	CurrentLanguage    = "/current_language" //get
	AvailableLanguages = "/languages"        //get

	SignIn   = "/auth/sign_in"   //post
	Refresh  = "/auth/refresh"   //get
	SignOut  = "/auth/sign_out"  //get
	SendCode = "/auth/send_code" //post

	SignUp           = "/user/sign_up"       //post
	UserByID         = "/user/id"            //get
	UserEditPassword = "/user/edit/password" //post
	UserEditEmail    = "/user/edit/email"    //post
	UserEditName     = "/user/edit/name"     //post

	DictionaryWord = "/dictionary/word"        //get post
	GetRandomWord  = "/dictionary/word/random" //get

	Vocabulary   = "/account/vocabulary"   //get post, put, delete
	Vocabularies = "/account/vocabularies" //get

	VocabularyWord         = "/vocabulary/word"               //get post delete
	VocabularyWordUpdate   = "/vocabulary/word/update"        //post
	VocabularySeveralWords = "/vocabulary/word/several"       //get
	VocabularyWords        = "/vocabulary/word/all"           //get
	GetPronunciation       = "/vocabulary/word/pronunciation" //get

	VocabularyTags = "/vocabulary/tag" //get
)
