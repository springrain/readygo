/*
 * @Author: your name
 * @Date: 2020-03-12 10:13:46
 * @LastEditTime: 2020-03-12 10:14:47
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: \readygo\apistruct\ResponseBodyModel.go
 */

package apistruct

// Api返回标准封装
type ResponseBodyModel struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
