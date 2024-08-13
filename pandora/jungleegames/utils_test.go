package jungleegames_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	jungleegames "bitbucket.org/junglee_games/getsetgo/pandora/jungleegames"
)

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(utilsSuite))
}

type utilsSuite struct {
	suite.Suite
}

func (s *utilsSuite) TestHasNextHasPrevWithInputAfterCursor() {
	var (
		before = ""
		after  = "sahil"
	)

	hasNext, hasPrevious := jungleegames.AnalyseHasNextOrPreviousCursor(&after, &before, 1, 2)

	s.Assert().Equal(hasNext, true)
	s.Assert().Equal(hasPrevious, true)
}

func (s *utilsSuite) TestHasNextHasPrevWithInputNoCursor() {
	hasNext, hasPrevious := jungleegames.AnalyseHasNextOrPreviousCursor(nil, nil, 1, 2)
	s.Assert().Equal(hasNext, true)
	s.Assert().Equal(hasPrevious, false)
}

func (s *utilsSuite) TestHasNextHasPrevWithInputBeforeCursor() {
	before := "sahil"

	hasNext, hasPrevious := jungleegames.AnalyseHasNextOrPreviousCursor(nil, &before, 1, 1)

	s.Assert().Equal(hasNext, true)
	s.Assert().Equal(hasPrevious, false)
}
