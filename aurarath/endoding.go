package aurarath

import "bytes"

func decode(s *bytes.Buffer, t interface {}, codec string)
func encode(t interface {}, *bytes.Buffer)
