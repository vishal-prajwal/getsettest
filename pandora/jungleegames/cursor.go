package jungleegames

type CursorInfo struct {
	NextCursor string
	HasNext    *bool
	PrevCursor string
	HasPrev    *bool
}

type CursorPaginationInput struct {
	BeforeCursor *string
	AfterCursor  *string
	Limit        int
}
