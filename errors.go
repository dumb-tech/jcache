package jcache

type JCError string

func (e JCError) Is(err error) bool {
	if err == nil {
		return false
	}

	return err.Error() == string(e)
}

func (e JCError) Error() string { return string(e) }

func NewJCError(msg string) JCError {
	return JCError(msg)
}

var ErrorCacheIsFull = NewJCError("items is full")
