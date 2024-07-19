package handler

const (
	CurrentLanguage    = "/v0/current_language" //get
	AvailableLanguages = "/v0/languages"        //get

	SignIn   = "/v0/auth/sign_in"   //post
	Refresh  = "/v0/auth/refresh"   //get
	SignOut  = "/v0/auth/sign_out"  //get
	SendCode = "/v0/auth/send_code" //post

	SignUp   = "/v0/user/sign_up" //post
	UserByID = "/v0/user/id"      //get

	DictionaryWord = "/v0/dictionary/word"        //get post
	GetRandomWord  = "/v0/dictionary/word/random" //get

	UserVocabulary   = "/v0/account/vocabulary"   //post, put, delete
	UserVocabularies = "/v0/account/vocabularies" //get

	Vocabularies            = "/v0/vocabularies"           //get
	Vocabulary              = "/v0/vocabulary"             //get
	VocabularyCopy          = "/v0/vocabulary/copy"        //get
	VocabularyAccessForUser = "/v0/vocabulary/access/user" //post delete

	VocabularyWord        = "/v0/vocabulary/word"               //get post delete
	VocabularyWordUpdate  = "/v0/vocabulary/word/update"        //post
	VocabularyRandomWords = "/v0/vocabulary/words/random"       //get
	VocabularyWords       = "/v0/vocabulary/words"              //get
	WordPronunciation     = "/v0/vocabulary/word/pronunciation" //get

	VocabularyTags = "/v0/vocabulary/tag" //get

	Accesses = "/v0/accesses" //get

	Subscribers = "/v0/respondents" //get
)
