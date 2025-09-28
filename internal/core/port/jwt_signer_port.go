package port

type JWTSigner interface {
	Sign(userID int64) (string, error)
	Verify(tokenStr string) (int64, error) // returns userID
}
