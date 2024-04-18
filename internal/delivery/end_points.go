package delivery

const (
	SignIn   = "/auth/sign_in"  //post
	Refresh  = "/auth/refresh"  //get
	SignOut  = "/auth/sign_out" //get
	SendCode = "/auth/send_code"
	SignUp   = "/auth/sign_up" //post

	UserByID = "/account/id" //get

	DictionaryWord = "/dictionary/word"        //get post
	GetRandomWord  = "/dictionary/word/random" //get

	Vocabulary   = "/vocabulary"   //get post, put, delete
	Vocabularies = "/vocabularies" //get

	VocabularyWord         = "/word/vocabulary"         //get post delete
	VocabularyWordUpdate   = "/word/vocabulary/update"  //post
	VocabularySeveralWords = "/word/vocabulary/several" //get
	VocabularyWords        = "/word/vocabulary/all"     //get

	VocabularyTags = "/tag/vocabulary"
)
