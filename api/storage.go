package api

//Storage : utils.RedisCli have the method ,so is the interface
type Storage interface {
	Shorten(url string, exp int64) (string, error)
	ShortLinkInfo(eid string) (interface{}, error)
	Unshorten(eid string) (string, error)
}
