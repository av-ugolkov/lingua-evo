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

	DictionaryWord = "/dictionary/word"        //get post
	GetRandomWord  = "/dictionary/word/random" //get

	UserVocabulary   = "/account/vocabulary"   //post, put, delete
	UserVocabularies = "/account/vocabularies" //get

	Vocabularies            = "/vocabularies"           //get
	Vocabulary              = "/vocabulary"             //get
	VocabularyInfo          = "/vocabulary/info"        //get
	VocabularyCopy          = "/vocabulary/copy"        //get
	VocabularyAccessForUser = "/vocabulary/access/user" //get post delete patch
	VocabulariesByUser      = "/vocabularies/user"      //get

	VocabularyWord        = "/vocabulary/word"               //get post delete
	VocabularyWordUpdate  = "/vocabulary/word/update"        //post
	VocabularyRandomWords = "/vocabulary/words/random"       //get
	VocabularyWords       = "/vocabulary/words"              //get
	WordPronunciation     = "/vocabulary/word/pronunciation" //get

	VocabularyTags = "/vocabulary/tag" //get

	Accesses = "/accesses" //get

	CheckSubscriber = "/subscriber/check" //get
	Subscribe       = "/user/subscribe"   //post
	Unsubscribe     = "/user/unsubscribe" //post
)
