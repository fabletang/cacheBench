package dbutil

import (
	//"fmt"
	"fmt"
)

const (
	RowPageMin = 5
	//ConnMax    = 300
	ConnMax    = 200
	//CpuTimes=4
)

//func GetConnNum(cpuNum int,tableCounts int) (connNum int,rowPageNum int){
func GetConnNum(tableCounts int) (connNum int, rowPageNum int) {
	tmpNum := tableCounts / RowPageMin
	fmt.Println("tableCounts:", tableCounts)
	if tmpNum <= ConnMax-1 {
		left := tableCounts % RowPageMin
		connNum = tmpNum
		if left != 0 {
			connNum += 1
		}
		rowPageNum = RowPageMin
		return
		//return connNum, rowPageNum
	}
	if tmpNum > ConnMax-1 {
		rowPageNum = tableCounts / (ConnMax-1)
		connNum = ConnMax-1
		if tableCounts%(ConnMax-1) != 0 {
			connNum += 1
		}
		return
	}

	//return connNum, rowPageNum
	return
}
