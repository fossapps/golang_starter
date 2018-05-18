package transformers

import "golang_starter/db"

type ResponsePermission struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

func TransformPermission(permission db.Permission) ResponsePermission {
	return ResponsePermission{
		Key: permission.Key,
		Description: permission.Description,
	}
}

func TransformPermissions(items []db.Permission) []ResponsePermission {
	var listItems []ResponsePermission
	for _, data := range items {
		listItems = append(listItems, TransformPermission(data))
	}
	return listItems
}
