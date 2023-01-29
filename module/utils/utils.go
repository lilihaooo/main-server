package utils

/**
*@authoer:singham<chenxiao.zhao>
*@createDate:2023/1/19
*@description:
 */
// Must _
func Must(i interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return i
}
