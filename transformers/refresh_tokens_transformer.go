package transformers

import "golang_starter/db"

type ResponseRefreshToken struct {
	Token string `json:"token"`
	User  string `json:"user"`
}

func TransformRefreshToken(item db.RefreshToken) ResponseRefreshToken {
	return ResponseRefreshToken{
		User: item.User,
		Token: item.Token,
	}
}

func TransformRefreshTokens(items []db.RefreshToken) []ResponseRefreshToken {
	var listItems []ResponseRefreshToken
	for _, data := range items {
		listItems = append(listItems, TransformRefreshToken(data))
	}
	return listItems
}
