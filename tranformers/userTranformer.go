package tranformers

import "github.com/Haleluak/kb-backend/models"

func UserTranformer(user *models.User) (map[string] interface{}) {
	return map[string] interface{}{
		"user_id": user.ID,
		"name": user.Name,
		"profile_pic_url": user.ProfilePicUrl,
		"profile_pic_thumbnail_url": user.ProfilePicThumbnailUrl,
		"exp_from": user.ExpFrom,
		"exp_to": user.ExpTo,
	}
}
