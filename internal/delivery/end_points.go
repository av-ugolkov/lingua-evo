package delivery

const (
	//auth
	SignIn  = "/auth/signin"  //post
	Refresh = "/auth/refresh" //get
	Logout  = "/auth/logout"  //get

	//dictionary
	DictionaryWord = "/dictionary/word"        //get post
	GetRandomWord  = "/dictionary/word/random" //get

	//vocabulary
	Vocabulary   = "/account/vocabulary"   //get post, put, delete
	Vocabularies = "/account/vocabularies" //get

	//word
	VocabularyWord         = "/vocabulary/word"         //get post delete patch
	VocabularySeveralWords = "/vocabulary/word/several" //get
	VocabularyWords        = "/vocabulary/word/all"     //get
)
