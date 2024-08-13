package jungleegames

func AnalyseHasNextOrPreviousCursor(
	afterCursor, beforeCursor *string,
	askedLimit, receivedResultLength int,
) (hasNext bool, hasPrev bool) {
	switch {
	case afterCursor == nil && beforeCursor == nil:
		if receivedResultLength > askedLimit {
			return true, false
		}

		return false, false

	case afterCursor != nil:
		if receivedResultLength > askedLimit {
			return true, true
		}

		return false, true

	case beforeCursor != nil:
		if receivedResultLength > askedLimit {
			return true, true
		}

		return true, false

	default:
		return false, false
	}
}
