package transformers

import "github.com/fossapps/starter/db"

// ResponseRefreshToken is response friendly db.RefreshToken
type ResponseRefreshToken struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

// TransformRefreshToken takes db data and returns response friendly data
func TransformRefreshToken(item db.RefreshToken) ResponseRefreshToken {
	return ResponseRefreshToken{
		User:  item.User,
		Token: item.Token,
	}
}

// TransformRefreshTokens takes slice of db data and returns response friendly data
func TransformRefreshTokens(items []db.RefreshToken) []ResponseRefreshToken {
	var listItems []ResponseRefreshToken
	for _, data := range items {
		listItems = append(listItems, TransformRefreshToken(data))
	}
	return listItems
}
