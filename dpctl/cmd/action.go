// @Author: wangzn04@gmail.com
// @Date: 2018-09-06 22:03:53

package cmd

// ActionRunner defines the interface to run action
type ActionRunner interface {
	Detail() string
	Name() string
	Run()
}
