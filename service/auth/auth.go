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

/*
type AuthInterface interface {
	CreateAuth(string, *TokenDetails) error
	FetchAuth(string) (string, error)
	DeleteTokens(*AccessDetails) error
}

type AuthService struct {
}

// AuthService implements the AuthInterface
var _ AuthInterface = &AuthService{}

// Save token metadata to Redis
func (tk *AuthService) CreateAuth(userId string, td *TokenDetails) error {
	//at := time.Unix(td.AtExpires, 0) //converting Unix to UTC(to Time object)
	//now := time.Now()

	//atCreated, err := tk.client.Set(td.TokenUuid, userId, at.Sub(now)).Result()
	//if err != nil {
	//	return err
	//}
	//
	//if atCreated == "0"  {
	//	return errors.New("no record inserted")
	//}
	return nil
}

// Check the metadata saved
func (tk *AuthService) FetchAuth(tokenUuid string) (string, error) {
	//userid, err := tk.client.Get(tokenUuid).Result()
	//if err != nil {
	//	return "", err
	//}
	//return userid, nil
	return "", nil
}

// Once a user row in the token table
func (tk *AuthService) DeleteTokens(authD *AccessDetails) error {
	////get the refresh uuid
	//refreshUuid := fmt.Sprintf("%s++%s", authD.TokenUuid, authD.UserId)
	////delete access token
	//deletedAt, err := tk.client.Del(authD.TokenUuid).Result()
	//if err != nil {
	//	return err
	//}
	////delete refresh token
	//deletedRt, err := tk.client.Del(refreshUuid).Result()
	//if err != nil {
	//	return err
	//}
	////When the record is deleted, the return value is 1
	//if deletedAt != 1 || deletedRt != 1 {
	//	return errors.New("something went wrong")
	//}
	return nil
}
*/
