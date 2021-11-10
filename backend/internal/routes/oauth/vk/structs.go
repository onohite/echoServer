package vk

type VkStruct struct {
	Response []struct {
		FirstName       string `json:"first_name"`
		ID              int    `json:"id"`
		LastName        string `json:"last_name"`
		CanAccessClosed bool   `json:"can_access_closed"`
		IsClosed        bool   `json:"is_closed"`
		Sex             int    `json:"sex"`
		Bdate           string `json:"bdate"`
		Photo400Orig    string `json:"photo_400_orig"`
		Games           string `json:"games"`
	} `json:"response"`
}
