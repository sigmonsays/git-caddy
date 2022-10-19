package main

import "fmt"

// any type of error during an operation
func NewRepoError(etype, name string) *RepoError {

	ret := &RepoError{}
	ret.Type = etype
	ret.Name = name
	return ret
}

// printf variant to set Err
func NewRepoErrorf(etype, name string, s string, args ...interface{}) *RepoError {
	ret := NewRepoError(etype, name)
	ret.Err = fmt.Errorf(s, args...)
	return ret
}

type RepoError struct {
	Type string // update, pull, status, etc
	Name string
	Err  error
}

func (me *RepoError) Error() string {
	return fmt.Sprintf("repo:%s error:%q", me.Err)
}

func (me *RepoError) WithError(e error) *RepoError {
	me.Err = e
	return me
}
