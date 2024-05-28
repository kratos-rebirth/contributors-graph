package main

type ContributorInfo struct {
	Login         string `json:"login"`
	Id            int    `json:"id"`
	AvatarUrl     string `json:"avatar_url"`
	Type          string `json:"type"`
	Contributions int    `json:"contributions"`
}

type ContributorInfoDownload struct {
	Login     string
	Id        int
	AvatarBuf []byte
}
