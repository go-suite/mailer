package auth

type AccessDetails struct {
	TokenUuid string
	UserName  string
}

type TokenDetails struct {
	AccessToken string
	TokenUuid   string
	AtExpires   int64
}
