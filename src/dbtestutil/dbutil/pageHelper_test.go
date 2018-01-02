package dbutil

import (
	"fmt"
	"testing"
)

func TestGetConnNum(t *testing.T) {
	connNum, pageNum := GetConnNum(404)
	fmt.Println(connNum, pageNum)
	//fmt.Printf("%+v %+v \n", connNum,pageNum)
	if connNum > ConnMax {
		t.Errorf("expecting connNum > %s, got %s", ConnMax, connNum)
	}
	if pageNum > RowPageMin {
		t.Errorf("expecting pageNum >%s, got %s", RowPageMin, pageNum)
	}

	connNum, pageNum = GetConnNum(4)
	fmt.Println(connNum, pageNum)
	//fmt.Printf("%+v %+v \n", connNum,pageNum)
	if connNum > ConnMax {
		t.Errorf("expecting connNum >10, got %s", connNum)
	}
	if pageNum != RowPageMin {
		t.Errorf("expecting pageNum ==10, got %s", pageNum)
	}
}
