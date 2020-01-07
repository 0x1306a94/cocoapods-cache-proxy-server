package model

type ReposParams struct {
	Name       string `uri:"name" binding:"required"`
	Git        string `form:"git"`
	Tag        string `form:"tag"`
	Submodules bool   `form:"submodules"`
}
